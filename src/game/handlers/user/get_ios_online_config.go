package user

import (
	"code.google.com/p/goprotobuf/proto"
	"config"
	"game/server"
	"github.com/golang/glog"
	"pb"
)

func GetIosOnlineConfigHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	msg := &pb.MsgGetOnlineConfig{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	glog.Info("GetIosOnlineConfigHandler in.", msg)

	res := &pb.MsgGetOnlineConfigRes{}

	versionId := msg.GetVersionName()
	config, bRet := config.GetIosOnlineConfigManager().GetConfig(versionId)
	glog.Info("GetIosOnlineConfigHandler config.", config)
	glog.Info("GetIosOnlineConfigHandler bRet.", bRet)
	if !bRet {
		glog.Info("GetIosOnlineConfigHandler false.")
		res.ReviewStatus = proto.Bool(false)
		return server.BuildClientMsg(m.GetMsgId(), res)
	}

	glog.Info("GetIosOnlineConfigHandler true.")
	res.ReviewStatus = proto.Bool(config.IsOpen)
	return server.BuildClientMsg(m.GetMsgId(), res)
}
