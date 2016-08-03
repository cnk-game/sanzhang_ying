package exchange

import (
	"code.google.com/p/goprotobuf/proto"
	"game/domain/exchange"	
	"game/server"
	"pb"	
	"github.com/golang/glog"
)

func GetGoodListHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	glog.Info("GetGoodListHandler in.")
	
	res := &pb.MsgGetGoodsListAck{}		
	
	exchange.GetGoodManager().Init();	//read data from db
	datas := exchange.GetGoodManager().GetData()
	glog.Infof("%v",datas)
	
	for _, data:= range datas{				
		msg := &pb.GoodsListDef{}		
		
		msg.GoodsId =  proto.Int(data.GoodsId)		
		msg.Name =  proto.String(data.Name)
		msg.Desc =  proto.String(data.Desc)
		msg.IconRes =  proto.String(data.IconRes)
		msg.NeedScore =  proto.Int(data.NeedScore)
		msg.TotalCount = proto.Int(data.TotalCount)
		msg.RemainderCount = proto.Int(data.RemainderCount)
	
		
		res.GoodsList = append(res.GoodsList, msg)
	}
	return server.BuildClientMsg(m.GetMsgId(), res)
	
	
	
	
}
