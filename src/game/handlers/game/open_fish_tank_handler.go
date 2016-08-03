package game

import (
	"code.google.com/p/goprotobuf/proto"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
	"config"
	"crypto/md5"
    "encoding/hex"
)


func OpenFishTankHandler(m *pb.ServerMsg, sess *server.Session) []byte {

    msg := &pb.MsgOpenFishTankReq{}
    err := proto.Unmarshal(m.GetMsgBody(), msg)
    if err != nil {
        glog.Error(err)
        return nil
    }

    userId := msg.GetUserId()

    res := &pb.MsgOpenFishTankRes{}

    if config.FIshFuncClose {
        res.Code = pb.MsgOpenFishTankRes_FAILED.Enum()
        res.Reason = proto.String("FuncClose")
    } else {
        res.Code = pb.MsgOpenFishTankRes_OK.Enum()
        res.Reason = proto.String("OK")
        sourceStr := userId + config.FishMD5Key
        h := md5.New()
        h.Write([]byte(sourceStr))
        token := hex.EncodeToString(h.Sum(nil))
        res.Token = proto.String(token)
        domainUser.GetUserFortuneManager().AddToken(userId, token)
    }
    domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_OPEN_FISH_TANK), res)

    return nil
}