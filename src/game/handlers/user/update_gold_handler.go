package user

import (
	"code.google.com/p/goprotobuf/proto"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
	"util"
)

func UpdateGoldHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	if m.GetClient() {
		return nil
	}

	msg := &pb.MsgUpdateGold{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	player.UserTasks.AccomplishTask(util.TaskAccomplishType_GOLD, msg.GetGold(), player.SendToClientFunc)

	if msg.GetRechargeDiamond() > 0 {
		player.UserTasks.AccomplishTask(util.TaskAccomplishType_RECHARGE_TOTAL_DIAMOND, int64(msg.GetRechargeDiamond()), player.SendToClientFunc)
	}

	if player.LastExp != int(msg.GetExp()) {
		player.LastExp = int(msg.GetExp())
		player.UserTasks.AccomplishTask(util.TaskAccomplishType_LEVEL, int64(msg.GetExp()), player.SendToClientFunc)
	}
	domainUser.GetRankingListUpdater().UpdateUser(player.User.UserId, 0)

	return server.BuildClientMsg3(m.GetMsgId(), m.GetMsgBody())
}
