package user

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	domainPrize "game/domain/prize"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"math/rand"
	"pb"
	"time"
	"util"
)

func RegisterHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	glog.Info("RegisterHandler in.")
	msg := &pb.MsgRegisterReq{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	res := &pb.MsgRegisterRes{}

	userName := msg.GetUsername()
	pwd := msg.GetUserpwd()
	phone := msg.GetPhone()
	code := int(msg.GetCode())

	nickName := msg.GetNickname()
	if nickName == "" {
		nickName = msg.GetModel()
	}
	if nickName == "" {
		nickName = "请修改昵称"
	}

	if userName == "" || pwd == "" {
		glog.Info("===>RegisterHandler ,username is nil, sess:", sess)
		res.Code = pb.MsgRegisterRes_FAILED.Enum()
		return server.BuildClientMsg(m.GetMsgId(), res)
	}

	userId, _ := domainUser.GetPhoneIsBind(phone)
	if userId != "" {
		glog.Info("===>RegisterHandler ,GetPhoneIsBind phone:", phone, " userId ", userId)
		res.Code = pb.MsgRegisterRes_FAILED.Enum()
		return server.BuildClientMsg(m.GetMsgId(), res)
	}

	ok := util.CheckVerify(phone, code)
	if !ok {
		glog.Info("===>RegisterHandler ,CheckVerify error, phone:", phone, "|code=", code)
		res.Code = pb.MsgRegisterRes_FAILED.Enum()
		return server.BuildClientMsg(m.GetMsgId(), res)
	}

	u := &domainUser.User{}

	if rand.Float64() < 0.5 {
		u.Gender = int(pb.Gender_BOY)
		u.PhotoUrl = fmt.Sprintf("%v", 6+rand.Int()%3)
	} else {
		u.Gender = int(pb.Gender_GIRL)
		u.PhotoUrl = fmt.Sprintf("%v", 1+rand.Int()%5)
	}

	// 用户创建
	u.UserId = domainUser.GetNewUserId()
	u.UserName = msg.GetUsername()
	u.Password = genPassword(msg.GetUserpwd())
	u.Nickname = nickName
	u.CreateTime = time.Now()
	u.ChannelId = msg.GetChannelId()
	u.Model = msg.GetModel()
	u.IsBind = true

	err = domainUser.SaveUser(u)
	if err != nil {
		glog.Info("保存用户失败err:", err, " user:", u)
		res.Code = pb.MsgRegisterRes_FAILED.Enum()
		return server.BuildClientMsg(m.GetMsgId(), res)
	} else {
		domainUser.SaveUserNameIdByUserName(u.UserId, u.UserName)
	}

	err1 := domainUser.SavePhoneUser(phone, u.UserId)
	glog.Info("++--绑定手机号成功 phone ", phone, " userId ", u.UserId, " err ", err1)

	domainUser.GetUserFortuneManager().EarnGold(u.UserId, 10000, "新账号") //新账号给10000 add by yelong
	domainUser.GetUserFortuneManager().SaveUserFortune(u.UserId)

	//tt := &domainPrize.UserTasks{}
	tt := domainPrize.NewUserTasks(u.UserId)
	tt.AccomplishTask(util.TaskAccomplishType_REGISTER, 1, nil)
	tt.SaveTasks()

	//player.UserTasks.AccomplishTask(util.TaskAccomplishType_REGISTER, 1, nil)
	//player.UserTasks.SaveTasks()

	res.Code = pb.MsgRegisterRes_OK.Enum()
	return server.BuildClientMsg(m.GetMsgId(), res)
}
