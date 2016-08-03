package game

import (
	domainGame "game/domain/game"
	domainUser "game/domain/user"
	"game/server"
	"pb"
)

func LeaveWaitQueueHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	domainGame.GetDeskManager().LeaveWaitQueue(player.User.UserId)

	return nil
}
