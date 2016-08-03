package safeBox

import (
	"code.google.com/p/goprotobuf/proto"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
)

func LoginSafeBoxHandler(m *pb.ServerMsg, sess *server.Session) []byte {
    glog.Info("LoginSafeBoxHandler in")
	msg := &pb.Msg_LoginSafeBoxReq{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

    userId := msg.GetUserId()
    pwd := msg.GetPwd()
	code := domainUser.GetUserFortuneManager().LoginSafeBox(userId, pwd)

	res := &pb.Msg_LoginSafeBoxRes{}
	res.Code = code
	glog.Info("LoginSafeBoxHandler, code=", code)
	domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_LOGIN_SAFEBOX), res)

	return nil
}