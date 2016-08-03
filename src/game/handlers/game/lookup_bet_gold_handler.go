package game

import (
	"code.google.com/p/goprotobuf/proto"
	domainGame "game/domain/game"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
)

func LookupBetGoldHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	msg := &pb.MsgLookupBetGold{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	domainGame.GetDeskManager().LookupBetGold(player.User.UserId, msg.GetBetUserId())

	return nil
}
