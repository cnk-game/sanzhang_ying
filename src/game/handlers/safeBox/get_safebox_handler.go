package safeBox

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
)

func GetSafeBoxHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	glog.Info("GetSafeBoxHandler in")
	msg := &pb.Msg_GetGoldSafeBoxReq{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	userId := msg.GetUserId()
	pwd := msg.GetPwd()
	gold := msg.GetGold()

	res := &pb.Msg_GetGoldSafeBoxRes{}
	_, ok := domainUser.GetUserFortuneManager().EarnGoldNoMsg(userId, gold, "取款")
	if ok {
		code, reason, savings, currGold, log := domainUser.GetUserFortuneManager().UpdateSavings(userId, pwd, -gold, "取款", "")
		if code != 0 {
			domainUser.GetUserFortuneManager().ConsumeGoldNoMsg(userId, gold, false, "取款错误，系统回收")
			res.Code = pb.Msg_GetGoldSafeBoxRes_FAILED.Enum()
			res.Reason = proto.String(reason)
		} else {
			res.Code = pb.Msg_GetGoldSafeBoxRes_OK.Enum()
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
		res.Code = pb.Msg_GetGoldSafeBoxRes_FAILED.Enum()
		res.Reason = proto.String("用户信息错误")
	}

	domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_GET_GOLD_SAFEBOX), res)
	domainUser.GetUserFortuneManager().SaveUserFortune(userId)
	domainUser.GetUserFortuneManager().UpdateUserFortune2(userId, 0)

	return nil
}
