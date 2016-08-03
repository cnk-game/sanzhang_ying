package slots

import (
	"code.google.com/p/goprotobuf/proto"
	domainGame "game/domain/game"
	domainUser "game/domain/user"
	"game/server"
	"pb"
	"config"
	domainSlots "game/domain/slots"
)

func GainSlotMachinePrizeHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	if player.SlotMachine.GetCardsLen() < 3 {
		return nil
	}

	times := 0

	cardType := domainGame.GetCardType(player.SlotMachine.GetCards())
	switch cardType {
	case domainGame.CARD_TYPE_BAO_ZI:
		times = config.Set_times
	case domainGame.CARD_TYPE_SHUN_JIN:
		times = config.FlushStr_times
	case domainGame.CARD_TYPE_SPECIAL:
		times = config.Special_times
	case domainGame.CARD_TYPE_SHUN_ZI:
		times = config.Str_times
	case domainGame.CARD_TYPE_JIN_HUA:
		times = config.Flush_times
	case domainGame.CARD_TYPE_DOUBLE:
		times = config.Pair_times
	}

	if times > 0 {
		domainUser.GetUserFortuneManager().EarnFortune(player.User.UserId, int64(player.SlotMachine.Coin*times), 0, 0, false, "老虎机")
		domainUser.GetUserFortuneManager().UpdateUserFortune(player.User.UserId)
		domainSlots.GetSlotGlobal().InputPool(-player.SlotMachine.Coin*times)
	}

	player.SlotMachine.Reset()

	res := &pb.Msg_SlotMachinesPrizeRes{}
	res.WinFold = proto.Int(times)

	return server.BuildClientMsg(m.GetMsgId(), res)
}
