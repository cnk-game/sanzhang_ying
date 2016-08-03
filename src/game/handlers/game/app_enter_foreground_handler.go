package game

import (
	domainGame "game/domain/game"
	domainUser "game/domain/user"
	"game/server"
	"pb"
)

func AppEnterForegroundHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	domainUser.GetBackgroundUserManager().DelUser(player.User.UserId)
	domainGame.GetDeskManager().AppEnterForeground(player.User.UserId)

	return nil
}
