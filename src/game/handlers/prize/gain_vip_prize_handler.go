package prize

import (
    "code.google.com/p/goprotobuf/proto"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
	"config"
	"strconv"
	"time"
	"util"
)

func GainVipPrizeHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	msg := &pb.MsgGainVipPrizeReq{}
    err := proto.Unmarshal(m.GetMsgBody(), msg)
    if err != nil {
        glog.Error(err)
        return nil
    }

    level := int(msg.GetGainLevel())

    f, _ := domainUser.GetUserFortuneManager().GetUserFortune(player.User.UserId)
    state, ok := f.VipTaskStates[strconv.Itoa(level)]
    if !ok {
        glog.Error("Vip level error", level)
        return nil
    }

    vip_configs := config.GetVipPriceConfigManager().GetVipConfig()
    prize, ok := vip_configs[level]
    if !ok {
        glog.Error("Vip level error", level)
        return nil
    }

    if state.StartTime + int64(prize.PrizeDays * 86400) < time.Now().Unix() {
        glog.Error("StartTime error", state.StartTime, level)
        return nil
    }

    if state.LastGainTime == util.GetDayZero() {
        glog.Error("LastGainTime error", state.LastGainTime, level)
        return nil
    }

    if prize.PrizeGold > 0 {
        domainUser.GetUserFortuneManager().EarnFortune(player.User.UserId, int64(prize.PrizeGold), 0, 0, true, "特权每日赠送")
        domainUser.GetUserFortuneManager().GainVipTask(player.User.UserId, level)
        domainUser.GetUserFortuneManager().UpdateUserFortune(player.User.UserId)
    }

	res := &pb.MsgGainVipPrizeRes{}
	res.Code = pb.MsgGainVipPrizeRes_OK.Enum()
	res.Level = proto.Int(level)

	return server.BuildClientMsg(m.GetMsgId(), res)
}

