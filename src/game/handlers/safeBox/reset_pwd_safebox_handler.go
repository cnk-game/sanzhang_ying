package safeBox

import (
	"code.google.com/p/goprotobuf/proto"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
	"util"
)

func ResetPwdSafeBoxHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	msg := &pb.Msg_ResetPwdSafeBoxReq{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

    userId := msg.GetUserId()
    newpwd := msg.GetNewpwd()
    phone := msg.GetPhone()
    verify := int(msg.GetCode())

    res := &pb.Msg_ResetPwdSafeBoxRes{}
    res.Code = pb.Msg_ResetPwdSafeBoxRes_FAILED.Enum()

	userIdTemp, _ := domainUser.GetPhoneIsBindSafeBox(phone)
	if userIdTemp != "" {
		if userIdTemp != userId {
			glog.Info("ResetPwdSafeBoxHandler error, phone ", verify, "have been bind ", userIdTemp)
			domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_RESET_PWD_SAFEBOX), res)
			return nil
		}
	}

	ok := util.CheckVerify(phone, verify)
	if !ok {
		glog.Info("CheckVerify error, verify=", verify)
		domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_RESET_PWD_SAFEBOX), res)
		return nil
	}

	code := domainUser.GetUserFortuneManager().ResetPwdSafeBox(userId, newpwd)
	if code == 0 {
		res.Code = pb.Msg_ResetPwdSafeBoxRes_OK.Enum()
		domainUser.SaveSafeBoxUser(phone, userId)
		domainUser.GetUserFortuneManager().SaveUserFortune(userId)
	}
	domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_RESET_PWD_SAFEBOX), res)

	return nil
}
