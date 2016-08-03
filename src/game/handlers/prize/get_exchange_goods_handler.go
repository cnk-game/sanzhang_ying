package prize

import (
	"code.google.com/p/goprotobuf/proto"
	iosActive "game/domain/iosActive"
	"game/server"
	"github.com/golang/glog"
	"pb"
)

func GetExchangeGoodsHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	res := &pb.MsgGetExchangeGoodsRes{}
	glog.Info("+++++++++++++GetExchangeGoodsHandler in")
	actives := iosActive.GetIosActiveManager().GetIosActives()
	for _, active := range actives {
		item := &pb.ExchangeGoodsInfo{}
		item.Id = proto.Int(active.Id)
		item.Remains = proto.Int(active.Remains)
		item.Total = proto.Int(active.Total)
		item.Name = proto.String(active.Name)
		item.Price = proto.Int(active.Price)
		item.Desc = proto.String(active.Desc)
		item.IconUrl = proto.String(active.IconUrl)
		res.ExchangeGoods = append(res.ExchangeGoods, item)
	}

	//glog.Info("+++++++++++++GetExchangeGoodsHandler in res", res)

	return server.BuildClientMsg(m.GetMsgId(), res)
}
