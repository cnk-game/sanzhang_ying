package newUserTask

import (
	domainUserTask "game/domain/newusertask"
	domainUser "game/domain/user"
	"game/server"
	"pb"
)

func GetTaskHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	userId := player.User.UserId
	mT := domainUserTask.GetNewUserTaskManager().BuildUserTask(userId)

	return server.BuildClientMsg(m.GetMsgId(), mT)
}
