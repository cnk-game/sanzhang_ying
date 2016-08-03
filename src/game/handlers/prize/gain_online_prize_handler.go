package prize

import (
	"code.google.com/p/goprotobuf/proto"
	domainPrize "game/domain/prize"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
	"time"
)

func GainOnlinePrizeHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	msg := &pb.MsgGainOnlinePrizeReq{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

    res := &pb.MsgGainOnlinePrizeRes{}
    res.Code = pb.MsgGainOnlinePrizeRes_FAILED.Enum()

	prize, ok := domainPrize.GetOnlinePrizeManager().GetOnlinePrize(int(msg.GetPrizeId()))
	if !ok {
		glog.Error("找不到在线奖励Id:", msg.GetPrizeId())
		return server.BuildClientMsg(m.GetMsgId(), res)
	}

	if player.OnlinePrizeGainRecords.IsGained(int(msg.GetPrizeId())) {
		glog.Error("当天在线奖励已经领取userId:", player.User.UserId, " prizeId:", msg.GetPrizeId())
		res.Code = pb.MsgGainOnlinePrizeRes_TODAY_ALREADY_GAIN.Enum()
		return server.BuildClientMsg(m.GetMsgId(), res)
	}

	now := time.Now()
	beginTime := time.Date(now.Year(), now.Month(), now.Day(), int(prize.BeginTime/60), int(prize.BeginTime%60), 0, 0, time.Local)
	endTime := time.Date(now.Year(), now.Month(), now.Day(), int(prize.EndTime/60), int(prize.EndTime%60), 0, 0, time.Local)
	glog.V(2).Info("奖励时间beginTime:", beginTime, " endTime:", endTime)

	if time.Since(beginTime).Seconds() < 0 {
		// 时间未到
		glog.V(2).Info("在线奖励时间未到")
		return server.BuildClientMsg(m.GetMsgId(), res)
	}
	if time.Since(endTime).Seconds() > 0 {
		// 时间已过
		glog.V(2).Info("在线奖励时间已过")
		return server.BuildClientMsg(m.GetMsgId(), res)
	}

	player.OnlinePrizeGainRecords.SetGained(int(msg.GetPrizeId()), now)

	if prize.PrizeGold > 0 || prize.PrizeDiamond > 0 {
		domainUser.GetUserFortuneManager().EarnFortune(player.User.UserId, int64(prize.PrizeGold), prize.PrizeDiamond, 0, false, "在线奖励")
	}

	if prize.PrizeExp > 0 {
		domainUser.GetUserFortuneManager().AddExp(player.User.UserId, prize.PrizeExp)
	}

	domainUser.GetUserFortuneManager().UpdateUserFortune(player.User.UserId)

	domainUser.GetPlayerManager().UpdateOnlinePrizeState(int(msg.GetPrizeId()), 1, player.User.UserId)

	res.Code = pb.MsgGainOnlinePrizeRes_OK.Enum()
	res.PrizeId = msg.PrizeId

	return server.BuildClientMsg(m.GetMsgId(), res)
}
