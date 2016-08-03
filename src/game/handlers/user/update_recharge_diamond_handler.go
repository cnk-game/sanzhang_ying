package user

import (
	"code.google.com/p/goprotobuf/proto"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
	"util"
)

func UpdateRechargeDiamondHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	msg := &pb.MsgUpdateRechargeDiamond{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	glog.Info("===>updateRechargeDiamond userId:", player.User.UserId, " msg:", msg)

	if msg.GetDiamond() <= 0 {
		return nil
	}

	player.UserTasks.AccomplishTask(util.TaskAccomplishType_RECHARGE_TOTAL_DIAMOND, int64(msg.GetDiamond()), player.SendToClientFunc)

	return nil
}
