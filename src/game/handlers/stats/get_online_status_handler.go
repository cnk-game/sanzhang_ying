package stats

import (
	"code.google.com/p/goprotobuf/proto"
	//"game/domain/user"
	"game/server"
	"pb"
	"math/rand"
)

func GetOnlineStatusHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	msg := &pb.MsgGetOnlineStatusRes{}
	ran := rand.Int() % 3000 + 2000
	msg.PlayerCount = proto.Int(ran)

	return server.BuildClientMsg(m.GetMsgId(), msg)
}
