package game

import (
	"code.google.com/p/goprotobuf/proto"
	domainGame "game/domain/game"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
)

func OpCardHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	msg := &pb.MsgOpCardReq{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	domainGame.GetDeskManager().OnOpCards(player.User.UserId, msg)

	return nil
}
