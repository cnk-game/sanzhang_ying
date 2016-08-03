package prize

import (
	"code.google.com/p/goprotobuf/proto"
	domainPrize "game/domain/prize"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
)

func GainTaskPrizeHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	msg := &pb.MsgGainTaskPrizeReq{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	taskId := int(msg.GetTaskId())

	if !player.UserTasks.IsTaskAccomplished(taskId) {
		return nil
	}

	if player.UserTasks.IsTaskGained(taskId) {
		return nil
	}

	// 发放奖励
	prize, ok := domainPrize.GetTaskPrizeManager().GetTaskPrize(taskId)
	if !ok {
		glog.Error("找不到任务模板taskId:", taskId)
		return nil
	}

	player.UserTasks.SetGainTask(taskId)

	if prize.PrizeGold > 0 || prize.PrizeDiamond > 0 || prize.PrizeScore > 0 {
		domainUser.GetUserFortuneManager().EarnFortune(player.User.UserId, int64(prize.PrizeGold), prize.PrizeDiamond, prize.PrizeScore, false, "任务奖励")
	}

	if prize.PrizeExp > 0 {
		domainUser.GetUserFortuneManager().AddExp(player.User.UserId, prize.PrizeExp)
	}

	if prize.PrizeItemType > 0 && prize.PrizeItemCount > 0 {
		if prize.PrizeItemType == int(pb.MagicItemType_FOURFOLD_GOLD) {
			domainUser.GetUserFortuneManager().BuyDoubleCard(player.User.UserId, 0, prize.PrizeItemCount)
		} else if prize.PrizeItemType == int(pb.MagicItemType_PROHIBIT_COMPARE) {
			domainUser.GetUserFortuneManager().BuyForbidCard(player.User.UserId, 0, prize.PrizeItemCount)
		} else if prize.PrizeItemType == int(pb.MagicItemType_REPLACE_CARD) {
			domainUser.GetUserFortuneManager().BuyChangeCard(player.User.UserId, 0, prize.PrizeItemCount)
		}
	}

	domainUser.GetUserFortuneManager().UpdateUserFortune(player.User.UserId)

	res := &pb.MsgGainTaskPrizeRes{}
	res.Code = pb.MsgGainTaskPrizeRes_OK.Enum()
	res.TaskId = msg.TaskId

	return server.BuildClientMsg(m.GetMsgId(), res)
}
