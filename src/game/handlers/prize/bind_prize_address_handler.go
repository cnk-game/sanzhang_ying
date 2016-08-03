package prize

import (
	"code.google.com/p/goprotobuf/proto"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
)

func BindPrizeAddressHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	msg := &pb.MsgBindPrizeAddressReq{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	player.User.ShippingAddressName = msg.GetName()
	player.User.ShippingAddressPhone = msg.GetPhone()
	player.User.ShippingAddressAddress = msg.GetAddress()
	player.User.ShippingAddressZipCode = msg.GetZipCode()
	player.User.ShippingAddressQQ = msg.GetQq()

	player.User.IsBindShippingAddress = true

	res := &pb.MsgBindPrizeAddressRes{}
	res.Code = pb.MsgBindPrizeAddressRes_OK.Enum()

	return server.BuildClientMsg(m.GetMsgId(), res)
}
