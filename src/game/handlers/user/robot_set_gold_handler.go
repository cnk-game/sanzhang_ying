package user

import (
	"code.google.com/p/goprotobuf/proto"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
)

func RobotSetGoldHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	if !player.User.IsRobot {
		glog.V(2).Info("===>非机器人修改金币:", sess.IP)
		return nil
	}

	msg := &pb.MsgRobotSetGold{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	glog.V(2).Info("==>机器人修改金币userId:", player.User.Nickname, " gold:", msg.GetGold())

	domainUser.GetUserFortuneManager().EarnGold(player.User.UserId, int64(msg.GetGold()), "机器人设置金币")

	return nil
}
