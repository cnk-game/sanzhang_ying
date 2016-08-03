package user

import (
	"code.google.com/p/goprotobuf/proto"
	domainUser "game/domain/user"
	domainGame "game/domain/game"
	"game/server"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2/bson"
	"pb"
	"util"
	"time"
)


func ReConnectHandler(m *pb.ServerMsg, sess *server.Session) []byte {
    glog.Info("ReConnectHandler in.")
	msg := &pb.Msg_ReConnectReq{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	res := &pb.Msg_ReConnectRes{}
	userId := msg.GetUserId()
	glog.Info("===>User ReConnect,userId=", userId)

    old_sess, ok := domainUser.GetPlayerManager().FindSessById(userId)
    if !ok {
		glog.Info("User not Reconnect, id=", userId)
		res.Code = pb.Msg_ReConnectRes_FAILED.Enum()
		return server.BuildClientMsg(m.GetMsgId(), res)
    }

    u, err := domainUser.FindByUserId(userId)
    if u.IsLocked {
        glog.Info("==>userId:", u.UserId, "账号被锁定")
        res.Code = pb.Msg_ReConnectRes_FAILED.Enum()
        return server.BuildClientMsg(m.GetMsgId(), res)
    }

    player := domainUser.GetPlayer(old_sess.Data)
    sess.Data = old_sess.Data
    if old_sess.OnLogout != nil {
        old_sess.OnLogout = nil
    }
    player.SessKey = bson.NewObjectId().Hex()

    sess.LoggedIn = true
    sess.OnLogout = player.OnLogout
    sess.Data = player
    player.LoginIP = sess.IP

    player.SendToClientFunc = func(msgId int32, body proto.Message) {
        sess.SendToClient(server.BuildClientMsg(int32(msgId), body))
    }

    player.SetUserCache()

    domainUser.GetPlayerManager().ChangeItem(userId, sess)

	// 重连成功
	res.InGame = proto.Bool(false)
	game := domainGame.GetGameManager().FindGame(player.LastGameId)
	if game != nil {
	    _, ok := game.Players[userId]
	    if ok {
	        res.InGame = proto.Bool(true)
	        go func(userId string, game *domainGame.GameItem) {
                timer := time.NewTimer(time.Second * 2)
                <-timer.C
                msg := &pb.MsgGetPokerDeskInfoRes{}
                msg.DeskInfo = game.BuildMessage()
                if len(game.Players) <= 2 {
                    msg.DeskInfo.LastEndTime = proto.Int64(0)
                }
                msg.Type = util.ToMatchType(game.GameType)
                domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_GET_POKER_DESK_INFO), msg)
            }(userId, game)
	    }
	}
	res.Code = pb.Msg_ReConnectRes_OK.Enum()
	return server.BuildClientMsg(m.GetMsgId(), res)
}
