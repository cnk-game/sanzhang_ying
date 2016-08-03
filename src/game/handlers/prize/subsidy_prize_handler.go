package prize

import (
	"code.google.com/p/goprotobuf/proto"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
)

func SubsidyPrizeHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	f, _ := domainUser.GetUserFortuneManager().GetUserFortune(player.User.UserId)
	if f.Gold >= 1999 {
		return nil
	}

	if f.SafeBox.Savings > 0 {
		return nil
	}

	player.User.ResetSubsidyPrizeTime()

	msg := &pb.MsgSubsidyPrizeRes{}
	if player.User.SubsidyPrizeTimes >= 3 {
		msg.Code = pb.MsgSubsidyPrizeRes_FAILED.Enum()
		return server.BuildClientMsg(m.GetMsgId(), msg)
	}

	player.User.SubsidyPrizeTimes++

	msg.Code = pb.MsgSubsidyPrizeRes_OK.Enum()
	msg.RemainCount = proto.Int(3 - player.User.SubsidyPrizeTimes)
	msg.GainGold = proto.Int(1000)

	ok := domainUser.GetUserFortuneManager().EarnFortune(player.User.UserId, 1000, 0, 0, false, "津贴")
	if !ok {
		glog.V(2).Info("发放津贴失败userId:", player.User.UserId, " 1000")
	}

	domainUser.GetUserFortuneManager().UpdateUserFortune(player.User.UserId)

	player.User.LuckyValue += 100

	return server.BuildClientMsg(m.GetMsgId(), msg)
}
