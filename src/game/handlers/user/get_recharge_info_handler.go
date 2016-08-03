package user

import (
	"code.google.com/p/goprotobuf/proto"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
)

func GetRechargeInfoHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	msg := &pb.MsgGetRechargeInfoRes{}

	f, ok := domainUser.GetUserFortuneManager().GetUserFortune(player.User.UserId)
	if !ok {
		glog.Error("获取用户财富信息失败userId:", msg.GetUserId())
		return nil
	}

	msg.Price1Count = proto.Int(0)
	if f.FirstRecharge10 {
		msg.Price1Count = proto.Int(1)
	}
	if f.FirstRecharge20 {
		msg.Price2Count = proto.Int(1)
	}
	if f.FirstRecharge30 {
		msg.Price3Count = proto.Int(1)
	}
	if f.FirstRecharge50 {
		msg.Price4Count = proto.Int(1)
	}
	if f.FirstRecharge100 {
		msg.Price5Count = proto.Int(1)
	}
	if f.FirstRecharge500 {
		msg.Price6Count = proto.Int(1)
	}
	msg.FirstRecharge = proto.Bool(f.FirstRecharge)

	return server.BuildClientMsg(m.GetMsgId(), msg)
}
