package report

import (
	"code.google.com/p/goprotobuf/proto"
	domainReport "game/domain/report"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
)

func ReportUserHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	msg := &pb.MsgReportUserHeadIllegalReq{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	_, err = domainUser.FindByUserId(msg.GetIllegalUserId())
	if err != nil {
		glog.Error("不存在此举报对象userId:", player.User.UserId, " reportUserId:", msg.GetIllegalUserId())
		return nil
	}

	l := &domainReport.ReportLog{}
	l.UserId = player.User.UserId
	l.ReportUserId = msg.GetIllegalUserId()
	domainReport.SaveReportLog(l)

	return nil
}
