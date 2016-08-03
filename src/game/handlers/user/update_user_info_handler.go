package user

import (
	"code.google.com/p/goprotobuf/proto"
	"game/domain/forbidWords"
	newUserTask "game/domain/newusertask"
	domainPrize "game/domain/prize"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2/bson"
	"pb"
	"util"
)

func UpdateUserInfoHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	msg := &pb.MsgUpdateUserInfoReq{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		return nil
	}

	res := &pb.MsgUpdateUserInfoRes{}

	if msg.GetUsername() != "" && msg.GetUserpwd() != "" {
		// 绑定账户
		if player.User.IsBind {
			// 已经绑定过
			res.Code = pb.MsgUpdateUserInfoRes_FAILED.Enum()
			res.Reason = proto.String("账号已被绑定")
			return server.BuildClientMsg(m.GetMsgId(), res)
			return nil
		}
		oldUserName := player.User.UserName
		// 更新账号
		player.User.UserName = msg.GetUsername()

		phone := msg.GetPhone()
		if phone == "" {
			phone = msg.GetUsername()
		}

		glog.Info("++--绑定手机 phone ", msg.GetPhone())
		glog.Info("++--绑定手机 GetUsername ", msg.GetUsername())

		userId, _ := domainUser.GetPhoneIsBind(phone)
		if userId != "" {
			res.Code = pb.MsgUpdateUserInfoRes_FAILED.Enum()
			res.Reason = proto.String("账号已被绑定")
			return server.BuildClientMsg(m.GetMsgId(), res)
			return nil
		}

		if msg.GetPhone() != "" {
			phone := msg.GetPhone()
			verify := int(msg.GetCode())
			ok := util.CheckVerify(phone, verify)
			if !ok {
				player.User.UserName = oldUserName
				res.Code = pb.MsgUpdateUserInfoRes_FAILED.Enum()
				res.Reason = proto.String("短信验证码错误")
				return server.BuildClientMsg(m.GetMsgId(), res)
			}
		}

		if err := domainUser.SaveUser(player.User); err != nil {
			// 账号已存在
			player.User.UserName = oldUserName
			res.Code = pb.MsgUpdateUserInfoRes_FAILED.Enum()
			res.Reason = proto.String("账号已被绑定")
			return server.BuildClientMsg(m.GetMsgId(), res)
		} else {
			domainUser.SaveUserNameIdByUserId(player.User.UserId, msg.GetUsername())
		}

		player.User.IsBind = true
		player.User.Password = genPassword(msg.GetUserpwd())

		player.UserTasks.AccomplishTask(util.TaskAccomplishType_REGISTER, 1, player.SendToClientFunc)

		/*if player.User.IsChangedNickname && player.User.IsChangedPhotoUrl {
			sendPrize(player)
		}*/
		domainUser.SaveUser(player.User)
		glog.Info("++--绑定手机号成功 SaveUser ")

		if msg.GetPhone() != "" {
			userId := player.User.UserId
			err := domainUser.SavePhoneUser(phone, userId)
			glog.Info("++--绑定手机号成功 phone ", phone, " userId ", userId, " err ", err)
		}

		CheckNewUserCompleteUserInfoTask(player.User.UserId, player.User.ChannelId)

		res.Code = pb.MsgUpdateUserInfoRes_BIND_OK.Enum()
		res.Reason = proto.String("操作已成功")
		res.User = player.User.BuildMessage(player.MatchRecord.BuildMessage())
		return server.BuildClientMsg(m.GetMsgId(), res)
	}

	if forbidWords.IsForbid(msg.GetNickName()) || forbidWords.IsForbid(msg.GetSigniture()) {
		glog.V(2).Info("==>包含敏感词,更新用户资料失败nickName:", msg.GetNickName(), " signiture:", msg.GetSigniture())
		res.Code = pb.MsgUpdateUserInfoRes_FAILED.Enum()
		res.Reason = proto.String("您的昵称或签名包含敏感词，请修改")
		return server.BuildClientMsg(m.GetMsgId(), res)
	}

	/*checkSendPrize := false
	if !player.User.IsChangedNickname || !player.User.IsChangedPhotoUrl {
		checkSendPrize = true
	}*/

	if msg.GetNickName() != "" {
		player.User.Nickname = msg.GetNickName()
		player.User.IsChangedNickname = true

		CheckNewUserTask(player.User.UserId)
	}

	if msg.Gender != nil {
		if msg.GetGender() == pb.Gender_BOY {
			player.User.Gender = int(pb.Gender_BOY)
		} else {
			player.User.Gender = int(pb.Gender_GIRL)
		}
	}

	if msg.Signiture != nil {
		player.User.Signiture = msg.GetSigniture()
	}

	if msg.PhotoUrl != nil {
		player.User.PhotoUrl = msg.GetPhotoUrl()
		player.User.IsChangedPhotoUrl = true
	}

	player.UserTasks.AccomplishTask(util.TaskAccomplishType_UPDATE_USER_INFO, 1, player.SendToClientFunc)

	res.Code = pb.MsgUpdateUserInfoRes_UPDATE_OK.Enum()
	res.Reason = proto.String("操作已成功")
	res.User = player.User.BuildMessage(player.MatchRecord.BuildMessage())
	return server.BuildClientMsg(m.GetMsgId(), res)
}

func CheckNewUserTask(userId string) {
	glog.Info("==>CheckNewUserTask in")
	result := newUserTask.GetNewUserTaskManager().CheckUserChangeNicknameTask(userId)
	if result > 0 {
		msgT := &pb.MsgNewbeTaskCompletedNotify{}
		msgT.Id = proto.Int(result)
		domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_NOTIFY_NEWBETASK_COMP), msgT)
	}
}

func CheckNewUserCompleteUserInfoTask(userId string, channel string) {
	glog.Info("==>CheckNewUserCompleteUserInfoTask in", userId, ", ", channel)
	result := newUserTask.GetNewUserTaskManager().CheckUserCompleteInfoTask(userId, channel)
	if result > 0 {
		msgT := &pb.MsgNewbeTaskCompletedNotify{}
		msgT.Id = proto.Int(result)
		domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_NOTIFY_NEWBETASK_COMP), msgT)
	}
}

func sendPrize(p *domainUser.GamePlayer) {
	prizeMail := &domainPrize.PrizeMail{}
	prizeMail.UserId = p.User.UserId
	prizeMail.MailId = bson.NewObjectId().Hex()
	prizeMail.Content = "恭喜您获得完善资料的奖励！"
	prizeMail.Gold = 5000

	p.PrizeMails.AddPrizeMail(prizeMail)

	mailMsg := &pb.MsgGetPrizeMailListRes{}
	mailMsg.Mails = append(mailMsg.Mails, prizeMail.BuildMessage())

	p.SendToClient(int32(pb.MessageId_GET_PRIZE_MAIL_LIST), mailMsg)
}
