package safeBox

import (
	"code.google.com/p/goprotobuf/proto"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
)

func ChangePwdSafeBoxHandler(m *pb.ServerMsg, sess *server.Session) []byte {
    glog.Info("ChangePwdSafeBoxHandler in")
	msg := &pb.Msg_ChangePwdSafeBoxReq{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

    userId := msg.GetUserId()
    oldpwd := msg.GetOldpwd()
    newpwd := msg.GetNewpwd()

    res := &pb.Msg_ChangePwdSafeBoxRes{}
	code, reason := domainUser.GetUserFortuneManager().ChangePwdSafeBox(userId, oldpwd, newpwd)
    if code != 0 {
        res.Code = pb.Msg_ChangePwdSafeBoxRes_FAILED.Enum()
        res.Reason = proto.String(reason)
	} else {
		res.Code = pb.Msg_ChangePwdSafeBoxRes_OK.Enum()
		domainUser.GetUserFortuneManager().SaveUserFortune(userId)
	}
	glog.Info("ChangePwdSafeBoxHandler, code=", code)
	domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_CHANGE_PWD_SAFEBOX), res)

	return nil
}
