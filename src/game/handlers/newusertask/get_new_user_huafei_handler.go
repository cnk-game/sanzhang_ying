package newUserTask

import (
	"code.google.com/p/goprotobuf/proto"
	domainUserTask "game/domain/newusertask"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
)

func GetHuafeiHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	msg := &pb.MsgExchangeNewbeTaskHf{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		return nil
	}

	phone := msg.GetMobile()
	userId := player.User.UserId
	res := &pb.MsgExchangeNewbeTaskHfRes{}

	result := domainUserTask.GetNewUserTaskManager().GetUserTaskHuafei(userId, phone)
	glog.Info("GetHuafeiHandler result ", result)
	if result == 1 {
		glog.Info("GetHuafeiHandler success ")
		res.Code = pb.MsgExchangeNewbeTaskHfRes_OK.Enum()
		res.Reason = proto.String("领取成功")
	} else if result == 2 {
		res.Code = pb.MsgExchangeNewbeTaskHfRes_FAILED_NO_ENOUGH_HF.Enum()
		res.Reason = proto.String("花费券不足")
	} else if result == 3 {
		res.Code = pb.MsgExchangeNewbeTaskHfRes_FAILED_MOBILE_EXCHANGED.Enum()
		res.Reason = proto.String("该手机已经领过")
	} else {
		res.Code = pb.MsgExchangeNewbeTaskHfRes_FAILED.Enum()
		res.Reason = proto.String("信息错误")
	}

	return server.BuildClientMsg(m.GetMsgId(), res)
}
