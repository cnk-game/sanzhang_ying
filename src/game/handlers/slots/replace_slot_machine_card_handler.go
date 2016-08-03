package slots

import (
	"code.google.com/p/goprotobuf/proto"
	domainGame "game/domain/game"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
	"util"
	domainSlots "game/domain/slots"
)

func ReplaceSlotMachineCardHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	msg := &pb.Msg_SlotMachinesReplaceReq{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	glog.V(2).Info("===>换牌msg:", msg)

	if msg.GetReplacePos() < 1 || msg.GetReplacePos() > 3 {
		return nil
	}

	if player.SlotMachine.GetCardsLen() < 3 {
		glog.V(2).Info("====>没有牌userId:", player.User.UserId, " cards:", player.SlotMachine.GetCards())
		return nil
	}

	res := &pb.Msg_SlotMachinesReplaceRes{}

	if player.SlotMachine.ReplaceCardTimes >= 3 {
		glog.Error("老虎机超过换牌次数userId:", player.User.UserId)
		res.Code = pb.Msg_SlotMachinesReplaceRes_LIMIT_COUNT.Enum()
		return server.BuildClientMsg(m.GetMsgId(), res)
	}

	card, ok := GetRandomCard(int(msg.GetReplacePos()-1), player.SlotMachine)
	if !ok {
		glog.V(2).Info("==>老虎机换牌失败!")
		res.Code = pb.Msg_SlotMachinesReplaceRes_FAILED.Enum()
		return server.BuildClientMsg(m.GetMsgId(), res)
	}

	_, _, ok = domainUser.GetUserFortuneManager().ConsumeGold(player.User.UserId, int64(player.SlotMachine.Coin), false, "老虎机")
	if !ok {
		glog.V(2).Info("老虎机换牌扣钱失败userId:", player.User.UserId, " consumeGold:", player.SlotMachine.Coin)
		res.Code = pb.Msg_SlotMachinesReplaceRes_LACK_GOLD.Enum()
		return server.BuildClientMsg(m.GetMsgId(), res)
	}

	domainSlots.GetSlotGlobal().InputPool(player.SlotMachine.Coin)

	player.SlotMachine.ReplaceCardTimes++

	player.SlotMachine.ReplaceCard(int(msg.GetReplacePos()-1), card)

	glog.V(2).Info("===>老虎机换牌userId:", player.User.UserId, " cards:", player.SlotMachine.GetCards(), " 换牌位置:", msg.GetReplacePos(), " card:", card)

	cardType := domainGame.GetCardType(player.SlotMachine.GetCards())
	switch cardType {
	case domainGame.CARD_TYPE_BAO_ZI:
		player.UserTasks.AccomplishTask(util.TaskAccomplishType_SLOT_MACHINE_X_BAO_ZI, 1, player.SendToClientFunc)
	case domainGame.CARD_TYPE_SHUN_JIN:
		player.UserTasks.AccomplishTask(util.TaskAccomplishType_SLOT_MACHINE_X_TONG_HUA_SHUN, 1, player.SendToClientFunc)
	case domainGame.CARD_TYPE_JIN_HUA:
		player.UserTasks.AccomplishTask(util.TaskAccomplishType_SLOT_MACHINE_X_TONG_HUA, 1, player.SendToClientFunc)
	case domainGame.CARD_TYPE_SPECIAL:
		player.UserTasks.AccomplishTask(util.TaskAccomplishType_SLOT_MACHINE_X_SPECIAL, 1, player.SendToClientFunc)
	}

	res.Code = pb.Msg_SlotMachinesReplaceRes_OK.Enum()
	res.ReplacePos = msg.ReplacePos
	res.ReplaceCard = proto.Int(int(card))

	return server.BuildClientMsg(m.GetMsgId(), res)
}
