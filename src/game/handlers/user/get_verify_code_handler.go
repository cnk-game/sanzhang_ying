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

func GetVerifyCodeHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	msg := &pb.Msg_GetVerifyCodeReq{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	res := &pb.Msg_GetVerifyCodeRes{}
	phone := msg.GetPhone()
	codeType := msg.GetCodeType()
	glog.Info("GetVerifyCodeHandler codeType = ", codeType)
	if codeType == pb.CodeType_USER_REGISTER {
		userId, _ := domainUser.GetPhoneIsBind(phone)
		_, err := domainUser.FindUserNameIdByUserName(phone)
		if userId != "" || err != mgo.ErrNotFound {
			res.IsOk = proto.Bool(false)
			res.Reason = proto.String("该手机号已经被注册")
			return server.BuildClientMsg(m.GetMsgId(), res)
		}
	}
	if codeType == pb.CodeType_USER_RESET_PWD {
		userId, _ := domainUser.GetPhoneIsBind(phone)
		if userId == "" {
			res.IsOk = proto.Bool(false)
			res.Reason = proto.String("该手机号没有绑定的账号")
			return server.BuildClientMsg(m.GetMsgId(), res)
		}
	}
	if codeType == pb.CodeType_BOX_SET_PWD {
		userId, _ := domainUser.GetPhoneIsBindSafeBox(phone)
		if userId != "" {
			res.IsOk = proto.Bool(false)
			res.Reason = proto.String("该手机号已经其他保管箱绑定")
			return server.BuildClientMsg(m.GetMsgId(), res)
		}
	}
	if codeType == pb.CodeType_BOX_RESET_PWD {
		userId, _ := domainUser.GetPhoneIsBindSafeBox(phone)
		if userId == "" {
			res.IsOk = proto.Bool(false)
			res.Reason = proto.String("该手机号没有绑定任何保管箱")
			return server.BuildClientMsg(m.GetMsgId(), res)
		}
	}

	ok := util.GetVerify(phone)
	if !ok {
		glog.Info("GetVerifyCodeHandler error, phone=", phone)
		res.Reason = proto.String("获取验证码错误")
	}

	res.IsOk = proto.Bool(ok)
	res.Reason = proto.String("获取验证码成功")
	return server.BuildClientMsg(m.GetMsgId(), res)
}
