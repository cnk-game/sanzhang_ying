package littlegame

import (
	"code.google.com/p/goprotobuf/proto"
	"game/domain/littlegame"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"

	"pb"
)

func LigameChipHandler(m *pb.ServerMsg, sess *server.Session) []byte {

	glog.Info("LigameChipHandler in.")

	msg := &pb.MsgLittleGameReq{}

	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	gametype := msg.GetType()
	chip := msg.GetChip()

	res := &pb.MsgLittleGameRes{}

	glog.Info("gametype=", gametype)
	glog.Info("chip=", chip)

	player := domainUser.GetPlayer(sess.Data)

	rescode, carType, coin, cardPoints := littlegame.GetLigameLogicManager().ShuffleCards(int(gametype), int(chip), player.User.UserId, player.User.Nickname)

	res.PaiType = proto.Int(carType)
	res.Coin = proto.Int(coin)
	res.PaiPoint = []int32{}
	for _, cardPoint := range cardPoints {
		res.PaiPoint = append(res.PaiPoint, cardPoint)
	}
	if rescode == 0 {
		res.Code = pb.MsgLittleGameRes_OK.Enum()
	} else if rescode == 1 {
		res.Code = pb.MsgLittleGameRes_FAILED_NO_ENOUGH_YE.Enum()
	} else {
		res.Code = pb.MsgLittleGameRes_FAILED.Enum()
	}

	glog.Info("res.Code=", res.Code)

	return server.BuildClientMsg(m.GetMsgId(), res)

}
