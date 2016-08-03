package active

import (
	"code.google.com/p/goprotobuf/proto"
	"crypto/md5"
	"fmt"
	"game/server"
	"github.com/golang/glog"
	"io"
	"pb"
)

func Md5func(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func GetActiveTokenHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	glog.Info("GetActiveTokenHandler in.")

	msg := &pb.MsgGetActiveTokenReq{}

	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	userid := msg.GetUserId()
	activeid := msg.GetActiveId()

	res := &pb.MsgGetActiveTokenRes{}

	glog.Info("userid=", userid)
	glog.Info("activeid=", activeid)

	sign := "QF" + userid + activeid + "123456"
	token := Md5func(sign)
	token = string(token[0:30])

	res.ActiveToken = proto.String(token)

	return server.BuildClientMsg(m.GetMsgId(), res)

}
