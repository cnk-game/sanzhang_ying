package user

import (
	"code.google.com/p/goprotobuf/proto"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
)

func ExchangeGameGoodsHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	msg := &pb.MsgExchangeGameGoodsReq{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	res := &pb.MsgExchangeGameGoodsRes{}

	defer func() {
		if res.GetCode() == pb.MsgExchangeGameGoodsRes_OK {
			domainUser.GetUserFortuneManager().UpdateUserFortune(player.User.UserId)
		}
	}()

	if msg.GetItemId() == 1 {
		if !domainUser.GetUserFortuneManager().BuyDoubleCard(player.User.UserId, 1, 1) {
			res.Code = pb.MsgExchangeGameGoodsRes_FAILED.Enum()
			return server.BuildClientMsg(m.GetMsgId(), res)
		}

		res.Code = pb.MsgExchangeGameGoodsRes_OK.Enum()
		res.BuyItemType = pb.MagicItemType_FOURFOLD_GOLD.Enum()
		res.BuyItemCount = proto.Int(1)
		return server.BuildClientMsg(m.GetMsgId(), res)
	} else if msg.GetItemId() == 2 {
		if !domainUser.GetUserFortuneManager().BuyDoubleCard(player.User.UserId, 10, 12) {
			res.Code = pb.MsgExchangeGameGoodsRes_FAILED.Enum()
			return server.BuildClientMsg(m.GetMsgId(), res)
		}

		res.Code = pb.MsgExchangeGameGoodsRes_OK.Enum()
		res.BuyItemType = pb.MagicItemType_FOURFOLD_GOLD.Enum()
		res.BuyItemCount = proto.Int(12)
		return server.BuildClientMsg(m.GetMsgId(), res)
	} else if msg.GetItemId() == 3 {
		if !domainUser.GetUserFortuneManager().BuyForbidCard(player.User.UserId, 1, 1) {
			res.Code = pb.MsgExchangeGameGoodsRes_FAILED.Enum()
			return server.BuildClientMsg(m.GetMsgId(), res)
		}

		res.Code = pb.MsgExchangeGameGoodsRes_OK.Enum()
		res.BuyItemType = pb.MagicItemType_PROHIBIT_COMPARE.Enum()
		res.BuyItemCount = proto.Int(1)
		return server.BuildClientMsg(m.GetMsgId(), res)
	} else if msg.GetItemId() == 4 {
		if !domainUser.GetUserFortuneManager().BuyForbidCard(player.User.UserId, 10, 12) {
			res.Code = pb.MsgExchangeGameGoodsRes_FAILED.Enum()
			return server.BuildClientMsg(m.GetMsgId(), res)
		}

		res.Code = pb.MsgExchangeGameGoodsRes_OK.Enum()
		res.BuyItemType = pb.MagicItemType_PROHIBIT_COMPARE.Enum()
		res.BuyItemCount = proto.Int(12)
		return server.BuildClientMsg(m.GetMsgId(), res)
	} else if msg.GetItemId() == 5 {
		if !domainUser.GetUserFortuneManager().BuyChangeCard(player.User.UserId, 1, 1) {
			res.Code = pb.MsgExchangeGameGoodsRes_FAILED.Enum()
			return server.BuildClientMsg(m.GetMsgId(), res)
		}

		res.Code = pb.MsgExchangeGameGoodsRes_OK.Enum()
		res.BuyItemType = pb.MagicItemType_REPLACE_CARD.Enum()
		res.BuyItemCount = proto.Int(1)

		return server.BuildClientMsg(m.GetMsgId(), res)
	} else if msg.GetItemId() == 6 {
		if !domainUser.GetUserFortuneManager().BuyChangeCard(player.User.UserId, 10, 12) {
			res.Code = pb.MsgExchangeGameGoodsRes_FAILED.Enum()
			return server.BuildClientMsg(m.GetMsgId(), res)
		}

		res.Code = pb.MsgExchangeGameGoodsRes_OK.Enum()
		res.BuyItemType = pb.MagicItemType_REPLACE_CARD.Enum()
		res.BuyItemCount = proto.Int(12)

		return server.BuildClientMsg(m.GetMsgId(), res)
	}

	return nil
}
