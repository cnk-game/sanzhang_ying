package config

import (
	"code.google.com/p/goprotobuf/proto"
	"game/domain/config"	
	"game/server"
	"pb"	
	"github.com/golang/glog"
)

func GetMatchConfigListHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	glog.Info("GetMatchConfigListHandler in.")
	
	res := &pb.MsgGetMatchConfigAck{}		
	
	
	datas := config.GetMatchConfigManager().GetData()
	glog.Infof("%v",datas)
	
	for _, data:= range datas{				
		msg := &pb.MatchConfigDef{}		
		
		msg.MatchType =  proto.Int(data.MatchType)		
		msg.Name =  proto.String(data.Name)
		msg.ChipList =  proto.String(data.ChipList)
		msg.DefaultBet =  proto.Int(data.DefaultBet)
		msg.AutoQuitLimit =  proto.Int(data.AutoQuitLimit)
		msg.EnterLimit =  proto.Int(data.EnterLimit)
		msg.MaxLimit = proto.Int(data.MaxLimit)
		msg.Reward =  proto.Int(data.Reward)		
		msg.QuickLimit =  proto.Int(data.QuickLimit)
		
		res.MatchConfig = append(res.MatchConfig, msg)
	}
	return server.BuildClientMsg(m.GetMsgId(), res)
	
}
