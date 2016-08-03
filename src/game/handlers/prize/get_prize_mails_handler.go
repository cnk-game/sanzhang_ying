package prize

import (
	domainUser "game/domain/user"
	"game/server"
	"pb"
)

func GetPrizeMailsHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	return server.BuildClientMsg(m.GetMsgId(), player.PrizeMails.BuildMessage())
}
