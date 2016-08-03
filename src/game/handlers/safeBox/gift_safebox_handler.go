package safeBox

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
)

func GiftSafeBoxHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	glog.Info("GiftSafeBoxHandler in")
	msg := &pb.Msg_GiftGoldSafeBoxReq{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	userId := msg.GetUserId()
	pwd := msg.GetPwd()
	gold := msg.GetGold()
	toUid := msg.GetToUid()

	glog.Info("==safeBoxgive 赠送金币 userId ", userId, " gold ", gold, " toUid ", toUid)

	res := &pb.Msg_GiftGoldSafeBoxRes{}
	if gold < 120000 {
		res.Code = pb.Msg_GiftGoldSafeBoxRes_FAILED.Enum()
		res.Reason = proto.String("赠送失败，赠送好友最少12万金币")
		domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_GIFT_GOLD_SAFEBOX), res)
		return nil
	}

	_, err1 := domainUser.FindByUserId(toUid) //user.go
	if err1 != nil {
		glog.Info("==safeBoxgive 获取用户账户信息失败 userId ", userId)
		res.Code = pb.Msg_GiftGoldSafeBoxRes_FAILED.Enum()
		res.Reason = proto.String("赠送失败 用户不存在")
		domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_GIFT_GOLD_SAFEBOX), res)
		return nil
	}

	code, reason, savings, _, log := domainUser.GetUserFortuneManager().UpdateSavings(userId, pwd, -gold, "赠送扣款", toUid)
	if code != 0 {
		glog.Info("==safeBoxgive 赠送扣款失败 userId ", userId)
		res.Code = pb.Msg_GiftGoldSafeBoxRes_FAILED.Enum()
		res.Reason = proto.String(reason)

		domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_GIFT_GOLD_SAFEBOX), res)
		return nil
	} else {
		glog.Info("==safeBoxgive 赠送扣款成功 userId ", userId)
		_, tok := domainUser.GetUserFortuneManager().GetUserFortune(toUid)
		if tok {
			glog.Info("==safeBoxgive 获取用户账户信息成功 userId ", toUid)
			ok := domainUser.GetUserFortuneManager().UpdateSavingsAdd(userId, toUid, gold)
			if !ok {
				glog.Info("==safeBoxgive 赠送失败 userId ", toUid)
				res.Code = pb.Msg_GiftGoldSafeBoxRes_FAILED.Enum()
				res.Reason = proto.String("赠送失败，请稍候尝试")
				domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_GIFT_GOLD_SAFEBOX), res)
				domainUser.GetUserFortuneManager().UpdateSavings(userId, pwd, gold, "赠送扣款失败返还", toUid)

				return nil
			} else {
				glog.Info("==safeBoxgive 赠送成功 userId ", toUid)
			}
		} else {
			glog.Info("==safeBoxgive 获取用户账户信息失败 userId ", toUid)
			ok := domainUser.GetUserFortuneManager().UpdateSavingsAddOffLine(userId, toUid, gold)
			if !ok {
				glog.Info("==safeBoxgive 赠送失败 userId ", toUid)
				res.Code = pb.Msg_GiftGoldSafeBoxRes_FAILED.Enum()
				res.Reason = proto.String("赠送失败，请稍候尝试")
				domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_GIFT_GOLD_SAFEBOX), res)
				domainUser.GetUserFortuneManager().UpdateSavings(userId, pwd, gold, "赠送扣款失败返还", toUid)
				return nil
			} else {
				glog.Info("==safeBoxgive 赠送成功 userId ", toUid)
			}
		}
	}

	domainUser.GetUserFortuneManager().SaveUserFortune(toUid)
	domainUser.GetUserFortuneManager().UpdateUserFortune2(toUid, 0)
	res.Code = pb.Msg_GiftGoldSafeBoxRes_OK.Enum()
	res.Savings = proto.Int64(savings)
	logmsg := &pb.BoxLogsDef{}
	logmsg.Reason = proto.String(log.Reason)
	logmsg.Gold = proto.Int64(log.Gold)
	logmsg.Savings = proto.Int64(log.Savings)
	logmsg.Datetime = proto.String(fmt.Sprintf("%d-%d-%d %02d:%02d:%02d", log.DateTime.Year(), log.DateTime.Month(), log.DateTime.Day(), log.DateTime.Hour(), log.DateTime.Minute(), log.DateTime.Second()))
	logmsg.ToUid = proto.String(log.ToUid)
	logmsg.LogId = proto.String(log.LogId)
	res.Log = logmsg

	domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_GIFT_GOLD_SAFEBOX), res)
	domainUser.GetUserFortuneManager().SaveUserFortune(userId)
	domainUser.GetUserFortuneManager().UpdateUserFortune2(userId, 0)
	return nil
}
