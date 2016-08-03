package user

import (
	"code.google.com/p/goprotobuf/proto"
	domainUser "game/domain/user"
	"game/server"
	"pb"
)

func GetShippingAddressHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	msg := &pb.MsgGetPrizeAddressRes{}
	msg.Name = proto.String(player.User.ShippingAddressName)
	msg.Phone = proto.String(player.User.ShippingAddressPhone)
	msg.Address = proto.String(player.User.ShippingAddressAddress)
	msg.ZipCode = proto.String(player.User.ShippingAddressZipCode)
	msg.Qq = proto.String(player.User.ShippingAddressQQ)

	return server.BuildClientMsg(m.GetMsgId(), msg)
}
