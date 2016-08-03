package game

import (
	"code.google.com/p/goprotobuf/proto"
	domainGame "game/domain/game"
	"game/server"
	"github.com/golang/glog"
	"pb"
)


func KickPlayerHandler(m *pb.ServerMsg, sess *server.Session) []byte {
    msg := &pb.MsgKickPlayerReq{}
    err := proto.Unmarshal(m.GetMsgBody(), msg)
    if err != nil {
        glog.Error(err)
        return nil
    }

    from_userId := msg.GetFromUserId()
    target_userId := msg.GetTargetUserId()

    domainGame.GetDeskManager().KickPlayer(from_userId, target_userId)

    return nil
}