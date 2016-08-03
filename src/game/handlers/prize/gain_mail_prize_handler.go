package prize

import (
	"code.google.com/p/goprotobuf/proto"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
)

func GainMailPrizeHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	msg := &pb.MsgGainMailPrizeReq{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	res := &pb.MsgGainMailPrizeRes{}
	res.MailId = msg.MailId

	mailPrize := player.PrizeMails.GetPrizeMail(msg.GetMailId())
	if mailPrize == nil {
		glog.V(2).Info("===>找不到奖励邮件")
		res.Code = pb.MsgGainMailPrizeRes_NOT_EXIST.Enum()
		return server.BuildClientMsg(m.GetMsgId(), res)
	}

	if mailPrize.Gold > 0 || mailPrize.Diamond > 0 || mailPrize.Score > 0 {
		domainUser.GetUserFortuneManager().EarnFortune(player.User.UserId, int64(mailPrize.Gold), mailPrize.Diamond, mailPrize.Score, false, "邮件奖励")
	}

	if mailPrize.ItemType > 0 && mailPrize.ItemCount > 0 {
		if mailPrize.ItemType == int(pb.MagicItemType_FOURFOLD_GOLD) {
			domainUser.GetUserFortuneManager().BuyDoubleCard(player.User.UserId, 0, mailPrize.ItemCount)
		} else if mailPrize.ItemType == int(pb.MagicItemType_PROHIBIT_COMPARE) {
			domainUser.GetUserFortuneManager().BuyForbidCard(player.User.UserId, 0, mailPrize.ItemCount)
		} else if mailPrize.ItemType == int(pb.MagicItemType_REPLACE_CARD) {
			domainUser.GetUserFortuneManager().BuyChangeCard(player.User.UserId, 0, mailPrize.ItemCount)
		}
	}
	domainUser.GetUserFortuneManager().AddExp(player.User.UserId, mailPrize.Exp)

	player.PrizeMails.RemoveMail(msg.GetMailId())

	domainUser.GetUserFortuneManager().UpdateUserFortune(player.User.UserId)

	// 领取成功
	res.Code = pb.MsgGainMailPrizeRes_OK.Enum()
	return server.BuildClientMsg(m.GetMsgId(), res)
}
