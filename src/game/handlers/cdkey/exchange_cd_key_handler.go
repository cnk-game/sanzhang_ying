package cdkey

import (
	"code.google.com/p/goprotobuf/proto"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
	"util"
)

func ExchangeCDKeyHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	msg := &pb.MsgExchangeCodePrizeReq{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	glog.V(2).Info("==>兑换CDKey:", msg)
	res := &pb.MsgExchangeCodePrizeRes{}

	code := msg.GetCode()
	userChannel := msg.GetUserChannel()
	imei := msg.GetImei()
	if userChannel == 0 || imei == "" {
		glog.Error("对话码请求错误 msg:", msg)
		res.Code = pb.MsgExchangeCodePrizeRes_FAILED.Enum()
		return server.BuildClientMsg(m.GetMsgId(), res)
	}
	userId := player.User.UserId
	result, count, desc := util.CheckCDKey(10122, code, int(userChannel), imei, userId)
	if result == 1 {
		if desc == 1 {
			domainUser.GetUserFortuneManager().EarnFortune(userId, int64(count), 0, 0, false, "cdkey兑换金币")
			res.Gold = proto.Int(count)
			res.Diamond = proto.Int(0)
		}

		if desc == 2 {
			domainUser.GetUserFortuneManager().EarnFortune(userId, 0, count, 0, false, "cdkey兑换金币")
			res.Gold = proto.Int(0)
			res.Diamond = proto.Int(count)
		}

		domainUser.GetUserFortuneManager().UpdateUserFortune(userId)

		res.Code = pb.MsgExchangeCodePrizeRes_OK.Enum()
		res.GameScore = proto.Int(0)
		res.ItemType = proto.Int(0)
		res.ItemCount = proto.Int(0)

		return server.BuildClientMsg(m.GetMsgId(), res)
	} else {
		if result == 10 {
			res.Code = pb.MsgExchangeCodePrizeRes_CODE_ALREADY_USE.Enum()
			return server.BuildClientMsg(m.GetMsgId(), res)
		} else {
			res.Code = pb.MsgExchangeCodePrizeRes_FAILED.Enum()
			return server.BuildClientMsg(m.GetMsgId(), res)
		}
	}
}
