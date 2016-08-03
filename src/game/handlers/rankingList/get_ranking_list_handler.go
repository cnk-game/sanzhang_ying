package rankingList

import (
	"code.google.com/p/goprotobuf/proto"
	domainRankingList "game/domain/rankingList"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
	"time"
	"util"
)

func GetRankingListHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	msg := &pb.MsgGetRankingListReq{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	res := &pb.MsgGetRankingListRes{}

	for _, t := range msg.GetTypes() {
		glog.Info("rankingList t", t)
		m := domainRankingList.GetRankingList().BuildMessage(t)
		if m != nil {
			res.Items = append(res.Items, m)
		}
	}

	//glog.Info("rankingList req", msg)

	player := domainUser.GetPlayer(sess.Data)

	now := time.Now()

	for _, item := range res.GetItems() {
		for _, rankingItem := range item.GetItems() {
			if rankingItem.GetOrder() == 1 && rankingItem.GetUser().GetUserId() == player.User.UserId {
				if item.GetType() == pb.RankingType_RECHARGE_YESTERDAY {
					if !util.CompareDate(now, player.User.YesterdayRechargeOrder1Time) {
						player.User.YesterdayRechargeOrder1 = true
						player.User.YesterdayRechargeOrder1Time = now
						player.UserTasks.AccomplishTask(util.TaskAccomplishType_YESTERDAY_RECHARGE_ORDER_1_X_TIMES, 1, player.SendToClientFunc)
						break
					}

					if !player.User.YesterdayRechargeOrder1 {
						player.User.YesterdayRechargeOrder1 = true
						player.UserTasks.AccomplishTask(util.TaskAccomplishType_YESTERDAY_RECHARGE_ORDER_1_X_TIMES, 1, player.SendToClientFunc)
					}
				} else if item.GetType() == pb.RankingType_EARNINGS_YESTERDAY {
					if !util.CompareDate(now, player.User.YesterdayEarningOrder1Time) {
						player.User.YesterdayEarningOrder1 = true
						player.User.YesterdayEarningOrder1Time = now
						player.UserTasks.AccomplishTask(util.TaskAccomplishType_YESTERDAY_EARNING_ORDER_1_X_TIMES, 1, player.SendToClientFunc)
						break
					}
					if !player.User.YesterdayEarningOrder1 {
						player.User.YesterdayEarningOrder1 = true
						player.UserTasks.AccomplishTask(util.TaskAccomplishType_YESTERDAY_EARNING_ORDER_1_X_TIMES, 1, player.SendToClientFunc)
					}
				} else if item.GetType() == pb.RankingType_RECHARGE_LAST_WEEK {
					nowYear, nowWeek := now.ISOWeek()
					oldYear, oldWeek := player.User.LastWeekRechargeOrder1Time.ISOWeek()
					if nowYear == oldYear && nowWeek == oldWeek {
						break
					}
					player.User.LastWeekRechargeOrder1 = true
					player.User.LastWeekRechargeOrder1Time = now
					player.UserTasks.AccomplishTask(util.TaskAccomplishType_LAST_WEEK_RECHARGE_ORDER_1_X_TIMES, 1, player.SendToClientFunc)
				}
				break
			}
		}
	}

	//glog.Info("rankingList ", res)

	return server.BuildClientMsg(m.GetMsgId(), res)
}
