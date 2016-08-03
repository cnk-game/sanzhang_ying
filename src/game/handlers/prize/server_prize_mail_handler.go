package prize

import (
	"code.google.com/p/goprotobuf/proto"
	domainPrize "game/domain/prize"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
)

func ServerPrizeMailHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	if m.GetClient() {
		return nil
	}

	player := domainUser.GetPlayer(sess.Data)

	msg := &pb.PrizeMailDef{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	prizeMail := &domainPrize.PrizeMail{}
	prizeMail.UserId = player.User.UserId
	prizeMail.MailId = msg.GetMailId()
	prizeMail.Content = msg.GetContent()
	prizeMail.Gold = int(msg.GetPrize().GetGold())
	prizeMail.Diamond = int(msg.GetPrize().GetDiamond())
	prizeMail.Score = int(msg.GetPrize().GetScore())
	prizeMail.ItemType = int(msg.GetPrize().GetItemType())
	prizeMail.ItemCount = int(msg.GetPrize().GetItemCount())

	player.PrizeMails.AddPrizeMail(prizeMail)

	mailMsg := &pb.MsgGetPrizeMailListRes{}
	mailMsg.Mails = append(mailMsg.Mails, prizeMail.BuildMessage())

	return server.BuildClientMsg(int32(pb.MessageId_GET_PRIZE_MAIL_LIST), mailMsg)
}
