package game

import (
	domainGame "game/domain/game"
	domainUser "game/domain/user"
	"game/server"
	"pb"
)

func AppEnterBackgroundHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	domainUser.GetBackgroundUserManager().SetUser(player.User.UserId, false)
	domainGame.GetDeskManager().AppEnterBackground(player.User.UserId)

	return nil
}
