package user

import (
	"code.google.com/p/goprotobuf/proto"
	"game/domain/offlineMsg"
	"game/domain/rankingList"
	"game/server"
	"github.com/golang/glog"
	"pb"
	"sync"
	"config"
	"util"
	"time"
	domainPrize "game/domain/prize"
)

type PlayerManager struct {
	sync.RWMutex
	items  map[string]*server.Session
	online int
}

var playerManager *PlayerManager

func UserHeartBeat() {
    lastNotify := util.GetDayZero()
    for {
        time.Sleep(5 * time.Second)
        lastNotify = NotifyVipTask(lastNotify)
    }
}

func NotifyVipTask(lastNotify int64) int64 {
    if util.GetDayZero() != lastNotify {
        for _, v := range playerManager.items {
            p := GetPlayer(v.Data)
            msg := &pb.MsgUpdateVipTask{}
            f, ok := GetUserFortuneManager().GetUserFortune(p.User.UserId)
            if ok {
                vip_configs := config.GetVipPriceConfigManager().GetVipConfig()
                for _, value := range f.VipTaskStates {
                    sub_msg := &pb.UserVipTaskDef{}
                    sub_msg.Level = proto.Int(value.VipTaskId)
                    if value.StartTime != 0 {
                        edTime := value.StartTime + int64(vip_configs[value.VipTaskId].PrizeDays * 86400)
                        if edTime < util.GetDayZero() {
                            sub_msg.State = proto.Int(0)
                            sub_msg.StartTime = proto.Int64(0)
                            sub_msg.EndTime = proto.Int64(0)
                        } else if value.LastGainTime != util.GetDayZero() {
                            sub_msg.State = proto.Int(2)
                            sub_msg.StartTime = proto.Int64(value.StartTime)
                            sub_msg.EndTime = proto.Int64(value.StartTime + int64(vip_configs[value.VipTaskId].PrizeDays * 86400))
                        } else {
                            sub_msg.State = proto.Int(1)
                            sub_msg.StartTime = proto.Int64(value.StartTime)
                            sub_msg.EndTime = proto.Int64(value.StartTime + int64(vip_configs[value.VipTaskId].PrizeDays * 86400))
                        }
                    } else {
                        sub_msg.State = proto.Int(0)
                        sub_msg.StartTime = proto.Int64(0)
                        sub_msg.EndTime = proto.Int64(0)
                    }
                    sub_msg.PrizeGold = proto.Int(vip_configs[value.VipTaskId].PrizeGold)
                    msg.VipTaskList = append(msg.VipTaskList, sub_msg)
                }
                playerManager.SendClientMsg(p.User.UserId, int32(pb.MessageId_UPDATE_VIP_TASK), msg)
            }
        }
        return util.GetDayZero()
    }
    return lastNotify
}

func (m *PlayerManager)FindPlayerById(userId string) *GamePlayer {
    sess, ok := m.items[userId]
    if !ok {
        return nil
    }

    return GetPlayer(sess.Data)
}

func (m *PlayerManager)prizeHeartBeat() {
    activePrizeId := -1
    for {
        time.Sleep(5 * time.Second)
        now := time.Now()
        prizes := domainPrize.GetOnlinePrizeManager().GetOnlineAllPrize()
        if activePrizeId == -1 {
            inTimeId := -1
            for _, v := range prizes {
                beginTime := time.Date(now.Year(), now.Month(), now.Day(), int(v.BeginTime/60), int(v.BeginTime%60), 0, 0, time.Local)
                endTime := time.Date(now.Year(), now.Month(), now.Day(), int(v.EndTime/60), int(v.EndTime%60), 0, 0, time.Local)
                if time.Since(beginTime).Seconds() < 0 || time.Since(endTime).Seconds() > 0 {
                    continue
                }
                inTimeId = v.PrizeID
                break
            }
            if inTimeId != -1 {
                activePrizeId = inTimeId
                m.UpdateOnlinePrizeState(activePrizeId, 2, "")
            }
        } else {
            inTimeId := -1
            for _, v := range prizes {
                beginTime := time.Date(now.Year(), now.Month(), now.Day(), int(v.BeginTime/60), int(v.BeginTime%60), 0, 0, time.Local)
                endTime := time.Date(now.Year(), now.Month(), now.Day(), int(v.EndTime/60), int(v.EndTime%60), 0, 0, time.Local)
                if time.Since(beginTime).Seconds() < 0 || time.Since(endTime).Seconds() > 0 {
                    continue
                }
                inTimeId = v.PrizeID
                break
            }
            if inTimeId == -1 {
                m.UpdateOnlinePrizeState(activePrizeId, 0, "")
                activePrizeId = -1
            }
        }
    }
}

func (m *PlayerManager)UpdateOnlinePrizeState(activePrizeId int, state int, userId string) {
    glog.Info("UpdateOnlinePrizeState in. activePrizeId=", activePrizeId, "|state=", state)
    msg := &pb.Msg_UpdateOnlinePrize{}
    prizes := domainPrize.GetOnlinePrizeManager().GetOnlineAllPrize()
    for _, v := range prizes {
        if activePrizeId == v.PrizeID {
            p := &pb.PrizeOnline{}
            p.PrizeId = proto.Int(v.PrizeID)
            p.PrizeTitle = proto.String(v.PrizeTitle)
            p.IcoRes = proto.String(v.IcoRes)
            p.PrizeGold = proto.Int(v.PrizeGold)
            p.State = proto.Int(state)
            msg.Prize = append(msg.Prize, p)
        }
    }
    if userId == "" {
        m.BroadcastClientMsg(int32(pb.MessageId_UPDATE_PRIZE_ONLINE), msg)
    } else {
        playerManager.SendClientMsg(userId, int32(pb.MessageId_UPDATE_PRIZE_ONLINE), msg)
    }

}

func init() {
	playerManager = &PlayerManager{}
	playerManager.items = make(map[string]*server.Session)
	rankingList.GetRankingList().BroadcastClientFunc = playerManager.BroadcastClientMsg
	rankingList.GetRankingList().SendToUser = playerManager.SendServerMsg

	go UserHeartBeat()
	go playerManager.prizeHeartBeat()
}

func GetPlayerManager() *PlayerManager {
	return playerManager
}

func (m *PlayerManager) AddItem(userId string, isRobot bool, sess *server.Session) bool {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.items[userId]; ok {
		// 已在线
		return m.items[userId] == sess
	}

	m.items[userId] = sess

	if !isRobot {
		m.online++
	}

	return true
}

func (m *PlayerManager)ChangeItem(userId string, sess *server.Session) bool {
	m.Lock()
	defer m.Unlock()

	if _, ok := m.items[userId]; ok {
		m.items[userId] = sess
		return true
	}

	return false
}

func (m *PlayerManager) DelItem(userId string, isRobot bool) {
	m.Lock()
	defer m.Unlock()

	delete(m.items, userId)

	if !isRobot {
		m.online--
	}
}

func (m *PlayerManager) Kickout(userId string) {
	sess := m.getSess(userId)
	if sess != nil {
		sess.Kickout()
	}
}

func (m *PlayerManager) FindSessById(userId string) (*server.Session, bool) {
	m.RLock()
	defer m.RUnlock()

	sess, ok := m.items[userId]
	return sess, ok
}

func (m *PlayerManager) getSess(userId string) *server.Session {
	m.RLock()
	defer m.RUnlock()

	return m.items[userId]
}

func (m *PlayerManager) IsOnline(userId string) bool {
	return m.getSess(userId) != nil
}

func (m *PlayerManager) SendServerMsg(srcId string, dstIds []string, msgId int32, body proto.Message) {
	if body != nil {
		b, err := proto.Marshal(body)
		if err != nil {
			glog.Error(err)
			return
		}
		m.SendServerMsg2(srcId, dstIds, msgId, b)
		return
	}
	m.SendServerMsg2(srcId, dstIds, msgId, nil)
}

func (m *PlayerManager) SendServerMsg2(srcId string, dstIds []string, msgId int32, body []byte) {
	go func() {
		msg := &pb.ServerMsg{}
		msg.Client = proto.Bool(false)
		msg.SrcId = proto.String(srcId)
		msg.MsgId = proto.Int32(msgId)
		msg.MsgBody = body

		for _, dstId := range dstIds {
			sess := m.getSess(dstId)
			if sess != nil {
				sess.SendMQ(msg)
			} else {
				if msgId == int32(pb.ServerMsgId_MQ_PRIZE_MAIL) {
					offlineMsg.PutOfflineMsg(dstId, int32(pb.ServerMsgId_MQ_PRIZE_MAIL), msg)
				}
			}
		}
	}()
}

func (m *PlayerManager) SendClientMsg(userId string, msgId int32, body proto.Message) {
	sess := m.getSess(userId)
	if sess == nil {
		return
	}
	if GetBackgroundUserManager().Filter(userId, msgId) {
		return
	}
	sess.SendToClient(server.BuildClientMsg(msgId, body))
}

func (m *PlayerManager) SendClientMsg2(dstIds []string, msgId int32, body proto.Message) {
	if len(dstIds) <= 0 {
		return
	}

	b := server.BuildClientMsg(msgId, body)

	for _, userId := range dstIds {
		if GetBackgroundUserManager().Filter(userId, msgId) {
			continue
		}
		sess := m.getSess(userId)
		if sess != nil {
			sess.SendToClient(b)
		}
	}
}

func (m *PlayerManager) BroadcastClientMsg(msgId int32, body proto.Message) {
	go func() {
		m.RLock()
		items := []*server.Session{}
		for userId, item := range m.items {
			if GetBackgroundUserManager().Filter(userId, msgId) {
				continue
			}
			items = append(items, item)
		}
		m.RUnlock()

		b := server.BuildClientMsg(msgId, body)

		for _, item := range items {
			if item != nil {
				item.SendToClient(b)
			}
		}
	}()
}

func (m *PlayerManager) GetOnlineCount() int {
	m.RLock()
	defer m.RUnlock()

	return m.online
}

func (m *PlayerManager) GetOnlineCountWithRobot() int {
	m.RLock()
	defer m.RUnlock()

	return len(m.items) * 5
}
