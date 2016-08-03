package user

import (
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
)

func GetShopLogsHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	logs, err := domainUser.FindShopLogs(player.User.UserId)
	if err != nil {
		glog.Error(err)
		return nil
	}
	glog.Info("GetShopLogsHandler, logs=", logs)

	msg := &pb.MsgGetShopLogRes{}
	for _, log := range logs {
		msg.LogList = append(msg.LogList, log.BuildMessage())
	}

	return server.BuildClientMsg(m.GetMsgId(), msg)
}
