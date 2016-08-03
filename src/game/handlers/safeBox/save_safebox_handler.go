package safeBox

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
)

func SaveSafeBoxHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	glog.Info("SaveSafeBoxHandler in")
	msg := &pb.Msg_SaveGoldSafeBoxReq{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	userId := msg.GetUserId()
	pwd := msg.GetPwd()
	gold := msg.GetGold()

	res := &pb.Msg_SaveGoldSafeBoxRes{}
	_, _, ok := domainUser.GetUserFortuneManager().ConsumeGoldNoMsg(userId, gold, false, "存入保管箱")
	if ok {
		code, reason, savings, currGold, log := domainUser.GetUserFortuneManager().UpdateSavings(userId, pwd, gold, "存款", "")
		if code != 0 {
			domainUser.GetUserFortuneManager().EarnGoldNoMsg(userId, gold, "存入保管箱错误，系统退还")
			res.Code = pb.Msg_SaveGoldSafeBoxRes_FAILED.Enum()
			res.Reason = proto.String(reason)
		} else {
			res.Code = pb.Msg_SaveGoldSafeBoxRes_OK.Enum()
			res.Savings = proto.Int64(savings)
			res.CurrGold = proto.Int64(currGold)
			logmsg := &pb.BoxLogsDef{}
			logmsg.Reason = proto.String(log.Reason)
			logmsg.Gold = proto.Int64(log.Gold)
			logmsg.Savings = proto.Int64(log.Savings)
			logmsg.LogId = proto.String(log.LogId)
			logmsg.Datetime = proto.String(fmt.Sprintf("%d-%d-%d %02d:%02d:%02d", log.DateTime.Year(), log.DateTime.Month(), log.DateTime.Day(), log.DateTime.Hour(), log.DateTime.Minute(), log.DateTime.Second()))
			res.Log = logmsg
		}
	} else {
		res.Code = pb.Msg_SaveGoldSafeBoxRes_FAILED.Enum()
		res.Reason = proto.String("用户信息错误")
	}

	domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_SAVE_GOLD_SAFEBOX), res)
	domainUser.GetUserFortuneManager().SaveUserFortune(userId)
	domainUser.GetUserFortuneManager().UpdateUserFortune2(userId, 0)

	return nil
}
