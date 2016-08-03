package littlegame

import (
	"code.google.com/p/goprotobuf/proto"
	"game/domain/littlegame"
	"game/server"
	"github.com/golang/glog"
	"pb"
)

func GetLigameConfHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	glog.Info("GetLigameConfHandler in.")

	res := &pb.MsgGetLittleGameConfigRes{}

	chipdatas := littlegame.GetCardConfigManager().GetChipData()
	glog.Infof("%v", chipdatas)

	for _, chipdata := range chipdatas {
		msg := &pb.LigameChipDef{}

		msg.LevelId = proto.Int(chipdata.LevelId)
		msg.Chips = []int32{}
		msg.Chips = append(msg.Chips, int32(chipdata.Chip1))
		msg.Chips = append(msg.Chips, int32(chipdata.Chip2))
		msg.Chips = append(msg.Chips, int32(chipdata.Chip3))

		res.LittlegameChip = append(res.LittlegameChip, msg)
	}

	multipledata := littlegame.GetCardConfigManager().GetMultipeData()
	glog.Infof("%v", multipledata)

	res.LittlegameMultiple = []int32{}
	//	res.LittlegameMultiple = append(res.LittlegameMultiple,multipledata.Single)
	res.LittlegameMultiple = append(res.LittlegameMultiple, int32(multipledata.Special))
	res.LittlegameMultiple = append(res.LittlegameMultiple, int32(multipledata.BaoA))
	res.LittlegameMultiple = append(res.LittlegameMultiple, int32(multipledata.BaoZi))
	res.LittlegameMultiple = append(res.LittlegameMultiple, int32(multipledata.ShunJin))
	res.LittlegameMultiple = append(res.LittlegameMultiple, int32(multipledata.JinHua))
	res.LittlegameMultiple = append(res.LittlegameMultiple, int32(multipledata.ShunZi))
	res.LittlegameMultiple = append(res.LittlegameMultiple, int32(multipledata.Double))

	return server.BuildClientMsg(m.GetMsgId(), res)

}
