package game

import (
	domainGame "game/domain/game"
	domainUser "game/domain/user"
	"game/server"
	"pb"
)

func JoinWaitQueueHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	domainGame.GetDeskManager().JoinWaitQueue(player.User.UserId)

	return nil
}
