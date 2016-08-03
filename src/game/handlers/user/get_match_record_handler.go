package user

import (
	domainUser "game/domain/user"
	"game/server"
	"pb"
)

func GetMatchRecordHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	msg := &pb.MsgGetMatchRecordRes{}
	msg.MatchRecord = player.MatchRecord.BuildMessage()

	return server.BuildClientMsg(m.GetMsgId(), msg)
}
