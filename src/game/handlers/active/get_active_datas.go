package active

import (
	"code.google.com/p/goprotobuf/proto"
	activeUser "game/domain/activedata"	
	"game/server"
	"pb"	
	"github.com/golang/glog"
	domainUser "game/domain/user"
)

func GetActiveDatasHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	glog.Info("GetActiveDatasHandler in.")
	
	player := domainUser.GetPlayer(sess.Data)
	
	
	res := &pb.MsgGetActiveDataAck{}	

	datas := activeUser.GetActiveManager().GetData(player.User.ChannelId)
	glog.Infof("%v",datas)
	glog.Infof("%v",player.User.ChannelId)
	
	for _, data:= range datas{				
		msg := &pb.ActiveDataDef{}
		msg.ActivityId =  proto.String(data.ActivityId)		
		msg.Name =  proto.String(data.Name)
		msg.IconRes =  proto.String(data.IconRes)
		msg.Desc =  proto.String(data.Desc)
		msg.ButtonTitle =  proto.String(data.ButtonTitle)
		msg.ShowViewId =  proto.Int(data.ShowViewId)
		msg.DailyBeginShowSecond = proto.Int(data.DailyBeginShowSecond);
		msg.DailyEndShowSecond = proto.Int(data.DailyEndShowSecond);
		msg.IsCompleteClose = proto.Int(data.IsCompleteClose);
		msg.StartDate = proto.String(data.StartDate);
		msg.ShowLeftButton = proto.Bool(data.ShowLeftButton);
		msg.LeftButtonTitle =  proto.String(data.LeftButtonTitle);
		msg.OpenUrl = proto.String(data.OpenUrl);
		res.ActiveData = append(res.ActiveData, msg)
	}
	return server.BuildClientMsg(m.GetMsgId(), res)
	
	
	
	
}
