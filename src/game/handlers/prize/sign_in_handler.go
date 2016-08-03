package prize

import (
	"code.google.com/p/goprotobuf/proto"
	domainUser "game/domain/user"
	"game/server"
	"pb"
)

func SignInHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	if player.SignInRecord.IsSignIn() {
		return nil
	}

	player.SignInRecord.ResetSignRecord()
	day := player.SignInRecord.SetSignIn()
	domainUser.GetUserFortuneManager().EarnFortune(player.User.UserId, int64(day*100), 0, 0, false, "连续签到")

	domainUser.GetUserFortuneManager().UpdateUserFortune(player.User.UserId)

	res := &pb.MsgSignInRes{}
	res.Ok = proto.Bool(true)

	return server.BuildClientMsg(m.GetMsgId(), res)
}
