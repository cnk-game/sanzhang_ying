package slots

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	domainSlots "game/domain/slots"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"math/rand"
	"pb"
	"time"
	"util"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func PlayLuckyWheelHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	msg := &pb.MsgPlayLuckWheelReq{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	if msg.GetAwardFold() <= 0 {
		msg.AwardFold = proto.Int(1)
	}

	msg.UserId = proto.String(player.User.UserId)
	msg.NickName = proto.String(player.User.Nickname)

	p, ok := domainSlots.GetWheelConfigManager().RandomPrize()
	if !ok {
		glog.V(2).Info("===>大转盘随机奖励失败")
		return nil
	}

	res := &pb.MsgPlayLuckWheelRes{}

	costGold := int(msg.GetAwardFold()) * 10000
	_, _, ok = domainUser.GetUserFortuneManager().ConsumeGold(msg.GetUserId(), int64(costGold) , false, "幸运大转盘")
	if !ok {
		glog.V(2).Info("===>幸运大转盘，金币不足userId:", msg.GetUserId(), " costGold:", costGold)
		res.Code = pb.MsgPlayLuckWheelRes_LACK_GOLD.Enum()
		return server.BuildClientMsg(m.GetMsgId(), res)
	}

	times := int(msg.GetAwardFold())

	res.Code = pb.MsgPlayLuckWheelRes_OK.Enum()

	res.SelectPosId = proto.Int(p.WheelSerialId)
	if p.PrizeGold > 0 || p.PrizeDiamond > 0 || p.PrizeScore > 0 {
		domainUser.GetUserFortuneManager().EarnFortune(msg.GetUserId(), int64(p.PrizeGold*times), p.PrizeDiamond*times, p.PrizeScore*times, false, "幸运大转盘")
	}
	if p.PrizeExp > 0 {
		domainUser.GetUserFortuneManager().AddExp(msg.GetUserId(), p.PrizeExp*times)
	}

	if p.PrizeItemType == int(pb.MagicItemType_FOURFOLD_GOLD) {
		domainUser.GetUserFortuneManager().BuyDoubleCard(msg.GetUserId(), 0, times*p.PrizeItemCount)
	} else if p.PrizeItemType == int(pb.MagicItemType_PROHIBIT_COMPARE) {
		domainUser.GetUserFortuneManager().BuyForbidCard(msg.GetUserId(), 0, times*p.PrizeItemCount)
	} else if p.PrizeItemType == int(pb.MagicItemType_REPLACE_CARD) {
		domainUser.GetUserFortuneManager().BuyChangeCard(msg.GetUserId(), 0, times*p.PrizeItemCount)
	}

	domainUser.GetUserFortuneManager().UpdateUserFortune(msg.GetUserId())

	if res.GetSelectPosId() == 4 {
		// 100积分
		domainUser.GetPlayerManager().BroadcastClientMsg(int32(pb.MessageId_CHAT), util.BuildSysBugle(fmt.Sprintf("恭喜%v在幸运大转盘中转到了100积分！", msg.GetNickName())))
	}

	if res.GetSelectPosId() >= 1 && res.GetSelectPosId() <= 8 {
		logMsg := &pb.MsgGetLuckWheelLogRes{}
		logItem := &pb.MsgGetLuckWheelLogRes_LuckWheelLogDef{}
		logItem.UserId = proto.String(msg.GetUserId())
		logItem.NickName = proto.String(msg.GetNickName())
		if times > 1 {
			logItem.PrizeName = proto.String(fmt.Sprintf("%v x %v", p.PrizeName, times))
		} else {
			logItem.PrizeName = proto.String(p.PrizeName)
		}
		logItem.Time = proto.Int64(time.Now().Unix())
		logMsg.LogList = append(logMsg.LogList, logItem)

		glog.V(2).Info("===>大转盘日志logMsg:", logMsg)

		domainUser.GetPlayerManager().BroadcastClientMsg(int32(pb.MessageId_GET_LUCK_WHEEL_LOG), logMsg)
	}

	return server.BuildClientMsg(m.GetMsgId(), res)
}
