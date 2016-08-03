package prize

import (
	domainUser "game/domain/user"
	"game/server"
	"pb"
)

func BuyDailyGiftBagHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	player.User.ResetSubsidyPrizeTime()

	if !domainUser.GetUserFortuneManager().SetBuyDailyGiftBag(player.User.UserId) {
		return nil
	}

	return server.BuildClientMsg(m.GetMsgId(), nil)
}
