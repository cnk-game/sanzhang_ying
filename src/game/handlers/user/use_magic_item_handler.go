package user

import (
	"code.google.com/p/goprotobuf/proto"
	domainGame "game/domain/game"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
)

func UseMagicItemHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	msg := &pb.MsgUseMagicItemReq{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	glog.V(2).Info("====>UseMagicItemHandler msg:", msg)

	res := &pb.MsgUseMagicItemRes{}
	res.ItemType = msg.ItemType

	if !domainGame.GetDeskManager().IsPlayingProps(player.User.UserId) {
		res.Code = pb.MsgUseMagicItemRes_FAILED.Enum()
		return server.BuildClientMsg(m.GetMsgId(), res)
	}

	if msg.GetItemType() == pb.MagicItemType_FOURFOLD_GOLD {
		// 翻倍卡
		if !player.User.IsRobot {
			if !domainUser.GetUserFortuneManager().ConsumeDoubleCard(player.User.UserId, 1) {
				res.Code = pb.MsgUseMagicItemRes_FAILED.Enum()
				return server.BuildClientMsg(m.GetMsgId(), res)
			}
		}
	} else if msg.GetItemType() == pb.MagicItemType_PROHIBIT_COMPARE {
		// 禁比卡
		if !player.User.IsRobot {
			if !domainUser.GetUserFortuneManager().ConsumeForbidCard(player.User.UserId, 1) {
				res.Code = pb.MsgUseMagicItemRes_FAILED.Enum()
				return server.BuildClientMsg(m.GetMsgId(), res)
			}
		}
	} else if msg.GetItemType() == pb.MagicItemType_REPLACE_CARD {
		// 换牌卡
		if !player.User.IsRobot {
			if !domainUser.GetUserFortuneManager().ConsumeChangeCard(player.User.UserId, 1) {
				res.Code = pb.MsgUseMagicItemRes_FAILED.Enum()
				return server.BuildClientMsg(m.GetMsgId(), res)
			}
		}
	} else {
		glog.V(2).Info("===>道具类型错误")
		res.Code = pb.MsgUseMagicItemRes_FAILED.Enum()
		return server.BuildClientMsg(m.GetMsgId(), res)
	}

	domainUser.GetUserFortuneManager().UpdateUserFortune(player.User.UserId)

	domainGame.GetDeskManager().OnConsumeProps(player.User.UserId, msg.GetItemType(), int(msg.GetReplaceCard()))

	return nil
}
