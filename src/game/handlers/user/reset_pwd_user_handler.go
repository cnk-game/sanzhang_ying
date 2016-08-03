package user

import (
	"code.google.com/p/goprotobuf/proto"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"pb"
	"util"
)

func ResetPwdUserHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	msg := &pb.Msg_ResetPwdUserReq{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	res := &pb.Msg_ResetPwdUserRes{}

	userId := msg.GetUserId()
	pwd := msg.GetNewpwd()
	phone := msg.GetPhone()
	code := int(msg.GetCode())

	if pwd == "" || userId == "" {
		glog.Info("===>ResetPwdUserHandler param error., sess:", sess)
		res.Code = pb.Msg_ResetPwdUserRes_FAILED.Enum()
		return server.BuildClientMsg(m.GetMsgId(), res)
	}

	ok := util.CheckVerify(phone, code)
	if !ok {
		glog.Info("===>ResetPwdUserHandler ,CheckVerify error, phone:", phone)
		res.Code = pb.Msg_ResetPwdUserRes_FAILED.Enum()
		return server.BuildClientMsg(m.GetMsgId(), res)
	}
	userIdTemp := ""

	uTemp, err := domainUser.FindUserNameIdByUserName(phone)

	if err == mgo.ErrNotFound {
		userId, _ := domainUser.GetPhoneIsBind(phone)
		if userId == "" {
			glog.Info("===>RegisterHandler ,GetPhoneIsBind error, phone:", phone)
			res.Code = pb.Msg_ResetPwdUserRes_FAILED.Enum()
			return server.BuildClientMsg(m.GetMsgId(), res)
		} else {
			userIdTemp = userId
		}
	} else {
		userIdTemp = uTemp.UserId
	}

	glog.Info("===>ResetPwdUserHandler userIdTemp:", userIdTemp)

	u, err := domainUser.FindByUserId(userIdTemp)
	if err != nil {
		glog.Info("===>RegisterHandler ,FindByUserId error, userId:", userId)
		res.Code = pb.Msg_ResetPwdUserRes_FAILED.Enum()
		return server.BuildClientMsg(m.GetMsgId(), res)
	}

	u.Password = genPassword(pwd)
	if player.User.UserName == phone {
		player.User.Password = u.Password
	}
	glog.Info("ResetPwdUserHandler:", " user:", u)
	err = domainUser.SaveUser(u)
	if err != nil {
		glog.Info("保存用户失败err:", err, " user:", u)
		res.Code = pb.Msg_ResetPwdUserRes_FAILED.Enum()
		return server.BuildClientMsg(m.GetMsgId(), res)
	}

	res.Code = pb.Msg_ResetPwdUserRes_OK.Enum()
	return server.BuildClientMsg(m.GetMsgId(), res)
}
