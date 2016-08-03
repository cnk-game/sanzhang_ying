package active

import (
	"code.google.com/p/goprotobuf/proto"
	activeUser "game/domain/active"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
	"time"
)

func GetActiveStatusHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data) //game_player.go
	glog.Info("GetActiveTokenHandler in.", player)

	res := &pb.Msg_GetAciveStatusAck{}
	userId := player.User.UserId

	isOpen := false

	now := time.Now()
	nowDay := now.Day()
	nowMon := now.Month()
	if nowMon == 2 && nowDay > 5 {
		isOpen = true
	}

	if nowMon == 2 && nowDay < 23 {
		isOpen = true
	}

	isCanBuy := true
	status := activeUser.GetActiveManager().GetStatus(userId)
	if status != "0" {
		isCanBuy = false
	}

	tm2 := int64(1454688000) //2016-2-6
	tm3 := int64(1456070400) //2016-2-22

	channel := player.User.ChannelId
	glog.Info("GetActiveTokenHandler channel", channel)

	if channel != "178" {
		isOpen = false
		tm2 = 1420045261
		tm3 = 1420045262
	}

	msg1 := &pb.ActiveStatusDef{}
	msg1.Id = proto.String("100612")
	msg1.IsCanBuy = proto.Bool(true)
	msg1.IsOpen = proto.Bool(isOpen)
	msg1.BeginTime = proto.Int64(tm2)
	msg1.EndTime = proto.Int64(tm3)

	res.ActiveStatus = append(res.ActiveStatus, msg1)

	msg2 := &pb.ActiveStatusDef{}
	msg2.Id = proto.String("100613")
	msg2.IsCanBuy = proto.Bool(isCanBuy)
	msg2.IsOpen = proto.Bool(isOpen)
	msg2.BeginTime = proto.Int64(tm2)
	msg2.EndTime = proto.Int64(tm3)

	glog.Info("GetActiveTokenHandler out.", res)

	res.ActiveStatus = append(res.ActiveStatus, msg2)
	domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_GET_ACTIVE_STATUS), res)
	return nil
}
