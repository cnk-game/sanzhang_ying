package user

import (
	"code.google.com/p/goprotobuf/proto"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
	"time"
)

const (
    EXCHANGE_GOLD = 1
    EXCHANGE_HORN = 2
)

func ExchangeGoldHandler(m *pb.ServerMsg, sess *server.Session) []byte {
    glog.Info("switch EXCHANGE_HORN in!!!!!!!!")
	player := domainUser.GetPlayer(sess.Data)

	msg := &pb.MsgExchangeGoldReq{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	res := &pb.MsgExchangeGoldRes{}

    f, ok := domainUser.GetUserFortuneManager().GetUserFortune(player.User.UserId)
    if !ok {
        res.Code = pb.MsgExchangeGoldRes_FAILED.Enum()
        return server.BuildClientMsg(m.GetMsgId(), res)
    }

    glog.Info("ExchangeGoldHandler in >>>", msg.GetExchangeType())
	if msg.GetExchangeType() == EXCHANGE_GOLD {
        oldGold := f.Gold

        ok, curGold, curDiamond := domainUser.GetUserFortuneManager().ExchangeGold(player.User.UserId, int(msg.GetDiamond()))
        if !ok {
            res.Code = pb.MsgExchangeGoldRes_FAILED.Enum()
            return server.BuildClientMsg(m.GetMsgId(), res)
        }

        res.Code = pb.MsgExchangeGoldRes_OK.Enum()
        res.ExchangeGold = proto.Int(int(curGold - int64(oldGold)))
        res.CurGold = proto.Int64(curGold)
        res.CurDiamond = proto.Int(curDiamond)

        domainUser.GetUserFortuneManager().UpdateUserFortune(player.User.UserId)

        // 兑换成功
        log := &domainUser.UserShopLog{}
        log.UserId = player.User.UserId
        log.ExchangeGold = int(res.GetExchangeGold())
        log.ExchangeHorn = 0
        log.Time = time.Now()
        domainUser.SaveShopLog(log)

        logs, _ := domainUser.FindShopLogs(player.User.UserId)
        logMsg := &pb.MsgGetShopLogRes{}
        for _, log := range logs {
            logMsg.LogList = append(logMsg.LogList, log.BuildMessage())
        }
        player.SendToClient(int32(pb.MessageId_GET_SHOP_LOG), logMsg)
	} else if msg.GetExchangeType() == EXCHANGE_HORN {
	    exDiamond := int(msg.GetDiamond())
	    glog.Info("switch EXCHANGE_HORN >>>", exDiamond)
        horn := 0
        if exDiamond == 2 {
            horn = 50
        } else if exDiamond == 10 {
            horn = 300
		} else if exDiamond == 6 {
			horn = 150
		} else if exDiamond == 12 {
			horn = 400
        }
	    glog.Info("switch EXCHANGE_HORN, horn=", horn, "|exDiamond=", exDiamond, "|f.Diamond=", f.Diamond)
        if f.Diamond >= exDiamond && horn != 0 {
            domainUser.GetUserFortuneManager().EarnFortune(player.User.UserId, 0, -exDiamond, 0, false, "兑换喇叭")
            domainUser.GetUserFortuneManager().EarnHorn(player.User.UserId, horn)
            domainUser.GetUserFortuneManager().SaveUserFortune(player.User.UserId)

            log := &domainUser.UserShopLog{}
            log.UserId = player.User.UserId
            log.ExchangeGold = 0
            log.ExchangeHorn = horn
            log.Time = time.Now()
            domainUser.SaveShopLog(log)

            logs, _ := domainUser.FindShopLogs(player.User.UserId)
            logMsg := &pb.MsgGetShopLogRes{}
            for _, log := range logs {
                logMsg.LogList = append(logMsg.LogList, log.BuildMessage())
            }
            player.SendToClient(int32(pb.MessageId_GET_SHOP_LOG), logMsg)

            res.Code = pb.MsgExchangeGoldRes_OK.Enum()
            res.CurDiamond = proto.Int(f.Diamond - exDiamond)
            res.CurHorn = proto.Int(f.Horn + horn)
            res.CurGold = proto.Int64(f.Gold)
            res.ExchangeHorn = proto.Int(horn)
            //domainUser.GetUserFortuneManager().UpdateUserFortune(player.User.UserId)
        } else {
            res.Code = pb.MsgExchangeGoldRes_FAILED.Enum()
        }
	}

    return server.BuildClientMsg(m.GetMsgId(), res)
}
