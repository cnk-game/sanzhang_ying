package pay

import (
	"code.google.com/p/goprotobuf/proto"
	activeUser "game/domain/active"
	domainPay "game/domain/pay"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
	"time"
)

func GetPayTokenHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	msg := &pb.Msg_GetPayTokenReq{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	res := &pb.Msg_GetPayTokenRes{}
	res.Token = proto.String("0")

	userId := msg.GetUserId()
	productId := msg.GetProductId()

	if productId == "100612" || productId == "100613" {
		tm2 := int64(1454688000)
		tm3 := int64(1456070400)

		cur := int64(time.Now().Unix())

		if cur > tm3 || cur < tm2 {
			res.Token = proto.String("0")
			res.Result = proto.Int(2)
			res.Reason = proto.String("对不起！该商品现在不能购买")
			domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_GET_PAY_TOKEN), res)
			return nil
		}

		if productId == "100613" {
			status := activeUser.GetActiveManager().GetStatus(userId)
			if status != "0" {
				res.Token = proto.String("0")
				res.Result = proto.Int(2)
				res.Reason = proto.String("对不起！该商品已经达到购买次数的上限")
				domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_GET_PAY_TOKEN), res)
				return nil
			}
		}
	}

	if productId == "100576" {
		bFisrt := domainUser.GetUserFortuneManager().GetFirstRecharge(userId)
		if bFisrt == true {
			res.Token = proto.String("0")
			res.Result = proto.Int(2)
			res.Reason = proto.String("对不起！该商品已经达到购买次数的上限")
			domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_GET_PAY_TOKEN), res)
			return nil
		}
	}

	token := domainPay.GetTokenManager().GetToken(userId, productId)

	//wjs 金立通道要从服务器下单 2016年1月20日11:31:14
	player := domainUser.GetPlayer(sess.Data)
	if player.User.ChannelId == "212" {

		submittime := GetJinliOrder(player.User.UserName, productId, token)
		res.Submittime = proto.String(submittime)

	}

	res.Token = proto.String(token)
	res.Result = proto.Int(1)

	domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_GET_PAY_TOKEN), res)

	return nil
}
