package pay

import (
	"code.google.com/p/goprotobuf/proto"
	//domainPay "game/domain/pay"
	"game/server"
	"pb"
)

func GetPayChangeHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	res := &pb.MsgUseThreePartyPayRes{}
	res.UseThreePartyPay = proto.Bool(false)
	//res.UseThreePartyPay = proto.Bool(domainPay.IPS_PAY_CHANGE)
	return server.BuildClientMsg(m.GetMsgId(), res)
}
