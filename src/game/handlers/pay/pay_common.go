package pay

import (
	"code.google.com/p/goprotobuf/proto"
	"config"
	"fmt"
	activeUser "game/domain/active"
	newUserTask "game/domain/newusertask"
	domainPay "game/domain/pay"
	domainUser "game/domain/user"
	"github.com/golang/glog"
	"pb"
	"strconv"
	"sync"
	"time"
	"util"
)

const (
	PRODUCT_VIP1     = "100559"
	PRODUCT_VIP2     = "100560"
	PRODUCT_VIP3     = "100561"
	PRODUCT_VIP4     = "100562"
	QUICK_PAY10      = "100563"
	QUICK_PAY50      = "100564"
	QUICK_PAY500     = "100565"
	QUICK_PAY20      = "100577"
	FIRST_PAY30      = "100576"
	IOS_QUICK_12     = "100583"
	IOS_QUICK_518    = "100584"
	IOS_QUICK_18     = "100587"
	IOS_VIP_25_VIP1  = "100585"
	IOS_VIP_898_VIP4 = "100586"
	ACTIVE_10        = "100612"
	ACTIVE_30        = "100613"
)

var commonPayMu sync.RWMutex
var (
	G_pruduct = make(map[string]int)
)

var IsInit = false

func InitPayCode() {
	G_pruduct["Q_100559"] = 28
	G_pruduct["Q_100560"] = 188
	G_pruduct["Q_100561"] = 588
	G_pruduct["Q_100562"] = 888
	G_pruduct["Q_100563"] = 10
	G_pruduct["Q_100564"] = 50
	G_pruduct["Q_100565"] = 500
	G_pruduct["Q_100577"] = 20
	G_pruduct["Q_100576"] = 30
	G_pruduct["Q_100583"] = 12
	G_pruduct["Q_100584"] = 518
	G_pruduct["Q_100587"] = 18
	G_pruduct["Q_100585"] = 25
	G_pruduct["Q_100586"] = 898
	G_pruduct["Q_100612"] = 10
	G_pruduct["Q_100613"] = 30
	G_pruduct["Q_100519"] = 10
	G_pruduct["Q_100520"] = 50
	G_pruduct["Q_100521"] = 100
	G_pruduct["Q_100522"] = 300
	G_pruduct["Q_100523"] = 500
	G_pruduct["Q_100524"] = 1000

	IsInit = true
}

func commonPay(orderId string, channel string, price string, payType string, state string, thirdOrder string) {
	glog.Info("comActivePay in")
	cok, userId, payCode := domainPay.GetTokenManager().CommonCheckToken(orderId)
	if !cok {
		glog.Info("==>commonPay token error:", orderId, "忽略返回!")
		return
	}

	PayCodeStr := "Q_" + payCode
	priceServer := G_pruduct[PayCodeStr]
	amountInt, _ := strconv.ParseFloat(price, 64)
	glog.Info("==>commonPay token priceServer ", priceServer)
	glog.Info("==>commonPay token amountInt ", amountInt)
	if int(amountInt) != int(priceServer) {
		return
	}

	ok := func() bool {
		commonPayMu.Lock()
		defer commonPayMu.Unlock()

		// 充值成功，记录充值日志
		l := &domainPay.CommonPayLog{}
		l, err := domainPay.FindCommonPayLog(orderId)
		if err == nil {
			glog.Info("==>订单已处理orderId:", orderId)
			return false
		}

		l.UserId = userId
		l.Order = orderId
		l.Price = price
		l.PayType = payType
		l.PayCode = payCode
		l.State = state
		l.ThirdOrder = thirdOrder
		l.Channel = channel

		err = domainPay.SaveCommonPayLog(l)
		if err != nil {
			glog.Error("保存充值记录失败err:", err)
		}

		// 支付日志
		payLog := &domainPay.PayLog{}
		payLog.OrderId = orderId
		payLog.UserId = userId
		payLog.PayCode = payCode

		amount, _ := strconv.ParseFloat(price, 64)

		payLog.Amount = int(amount)

		u, err := domainUser.FindByUserId(userId)
		if err == nil && u != nil {
			payLog.Channel = u.ChannelId
		}
		payLog.PayChannel = channel
		payLog.PayType = payType
		domainPay.SavePayLog(payLog)

		return true
	}()
	glog.Info("ok=", ok)
	if !ok {
		glog.Info("==>订单已处理orderId:", orderId, "忽略返回!")
		return
	}

	if payCode == PRODUCT_VIP1 || payCode == PRODUCT_VIP2 || payCode == PRODUCT_VIP3 || payCode == PRODUCT_VIP4 || payCode == IOS_VIP_25_VIP1 || payCode == IOS_VIP_898_VIP4 {
		commonBuyVIP(userId, payCode, price)
	} else if payCode == QUICK_PAY10 || payCode == QUICK_PAY20 || payCode == QUICK_PAY50 || payCode == QUICK_PAY500 || payCode == FIRST_PAY30 || payCode == IOS_QUICK_12 || payCode == IOS_QUICK_18 || payCode == IOS_QUICK_518 {
		commonQuickPay(userId, payCode, price)
	} else if payCode == ACTIVE_10 || payCode == ACTIVE_30 {
		commonActivePay(userId, payCode, price)
	} else {
		commonBuyDiamond(userId, payCode, price)
	}

	result := newUserTask.GetNewUserTaskManager().CheckUserRechargeTask(userId)
	if result > 0 {
		msgT := &pb.MsgNewbeTaskCompletedNotify{}
		msgT.Id = proto.Int(result)
		domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_NOTIFY_NEWBETASK_COMP), msgT)
	}
}

func commonActivePay(userId string, payCode string, price string) {
	glog.Info("commonActivePay in")
	amount, _ := strconv.ParseFloat(price, 64)

	gold := 0
	charm := 0
	horn := 0
	if payCode == ACTIVE_10 {
		gold = 110000
		charm = 100
		horn = 20
	} else if payCode == ACTIVE_30 {
		gold = 580000
	}
	domainUser.GetUserFortuneManager().EarnFortune(userId, int64(gold), 0, 0, true, "活动充值")

	if payCode == ACTIVE_30 {
		domainUser.GetUserFortuneManager().OpenSafeBox(userId, "")
		activeUser.GetActiveManager().AddItem(userId, payCode)
		activeUser.GetActiveManager().SaveStatus(userId, payCode)
	}
	if payCode == ACTIVE_10 {
		domainUser.GetUserFortuneManager().EarnCharm(userId, charm)
		domainUser.GetUserFortuneManager().EarnHorn(userId, horn)
	}
	domainUser.GetUserFortuneManager().SaveUserFortune(userId)
	domainUser.GetUserFortuneManager().UpdateRechargeRankingList(userId, int(amount))
	domainUser.GetUserFortuneManager().UpdateUserFortune2(userId, 0)
}

func commonQuickPay(userId string, payCode string, price string) {
	glog.Info("commonQuickPay in")

	amount, _ := strconv.ParseFloat(price, 64)

	gold := 0
	if payCode == QUICK_PAY10 {
		gold = 100000
	} else if payCode == QUICK_PAY50 {
		gold = 500000
	} else if payCode == QUICK_PAY500 {
		gold = 5000000
	} else if payCode == FIRST_PAY30 {
		gold = 580000
	} else if payCode == IOS_QUICK_12 {
		gold = 120000
	} else if payCode == IOS_QUICK_18 {
		gold = 180000
	} else if payCode == IOS_QUICK_518 {
		gold = 5180000
	} else if payCode == QUICK_PAY20 {
		gold = 200000
	}
	domainUser.GetUserFortuneManager().EarnFortune(userId, int64(gold), 0, 0, true, "快速充值")
	domainUser.GetUserFortuneManager().SaveUserFortune(userId)
	domainUser.GetUserFortuneManager().UpdateRechargeRankingList(userId, int(amount))
	if payCode == FIRST_PAY30 {
		domainUser.GetUserFortuneManager().OpenSafeBox(userId, "")
		domainUser.GetUserFortuneManager().UpdateFirstRecharge(userId)
	}

	domainUser.GetUserFortuneManager().UpdateUserFortune2(userId, 0)
}

func commonBuyDiamond(userId string, payCode string, price string) {
	glog.Info("commonBuyDiamond in")
	amount, _ := strconv.ParseFloat(price, 64)

	domainUser.GetUserFortuneManager().EarnFortune(userId, 0, int(amount), 0, true, "钻石充值")
	domainUser.GetUserFortuneManager().SaveUserFortune(userId)

	domainUser.GetUserFortuneManager().UpdateUserFortune2(userId, int(amount))

	domainUser.GetUserFortuneManager().UpdateRechargeRankingList(userId, int(amount))

	// 商场日志
	shopLogMsg := &pb.MsgGetShopLogRes{}

	log := &domainUser.UserShopLog{}
	log.UserId = userId
	log.RechargeDiamond = int(amount)
	log.Time = time.Now()
	domainUser.SaveShopLog(log)

	shopLogMsg.LogList = append(shopLogMsg.LogList, log.BuildMessage())

	domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_GET_SHOP_LOG), shopLogMsg)

	updateRechargeDiamond(userId, int(amount))
}

func commonBuyVIP(userId string, payCode string, price string) {
	glog.Info("commonBuyVIP in")
	//amount, _ := strconv.ParseFloat(price, 64)
	vip_configs := config.GetVipPriceConfigManager().GetVipConfig()
	charm := 0
	goldNow := 0
	level := 0
	name := ""
	if payCode == PRODUCT_VIP1 {
		charm = vip_configs[1].PrizeCharm
		name = vip_configs[1].Name
		goldNow = vip_configs[1].PrizeGoldNow
		level = 1
	} else if payCode == PRODUCT_VIP2 {
		charm = vip_configs[2].PrizeCharm
		name = vip_configs[2].Name
		goldNow = vip_configs[2].PrizeGoldNow
		level = 2
	} else if payCode == PRODUCT_VIP3 {
		charm = vip_configs[3].PrizeCharm
		name = vip_configs[3].Name
		goldNow = vip_configs[3].PrizeGoldNow
		level = 3
	} else if payCode == PRODUCT_VIP4 {
		charm = vip_configs[4].PrizeCharm
		name = vip_configs[4].Name
		goldNow = vip_configs[4].PrizeGoldNow
		level = 4
	} else if payCode == IOS_VIP_25_VIP1 {
		charm = vip_configs[11].PrizeCharm
		name = vip_configs[11].Name
		goldNow = vip_configs[11].PrizeGoldNow
		level = 1
	} else if payCode == IOS_VIP_898_VIP4 {
		charm = vip_configs[12].PrizeCharm
		name = vip_configs[12].Name
		goldNow = vip_configs[12].PrizeGoldNow
		level = 4
	}

	domainUser.GetUserFortuneManager().EarnFortune(userId, int64(goldNow), 0, 0, true, "VIP购买")
	domainUser.GetUserFortuneManager().EarnCharm(userId, charm)
	domainUser.GetUserFortuneManager().UpdateVipTaskState(userId, level)
	domainUser.GetUserFortuneManager().UpdateVipLevel(userId, level)
	domainUser.GetUserFortuneManager().OpenSafeBox(userId, "")

	domainUser.GetUserFortuneManager().SaveUserFortune(userId)
	domainUser.GetUserFortuneManager().UpdateUserFortune2(userId, 0)

	player := domainUser.GetPlayerManager().FindPlayerById(userId)
	str_s := "特权"
	switch name {
	case "VIP1":
		str_s = "特权1"
	case "VIP2":
		str_s = "特权2"
	case "VIP3":
		str_s = "特权3"
	case "VIP4":
		str_s = "特权4"
	case "ISOVIP1":
		str_s = "特权1"
	case "ISOVIP4":
		str_s = "特权4"
	}
	domainUser.GetPlayerManager().BroadcastClientMsg(int32(pb.MessageId_CHAT), util.BuildSysBugle(fmt.Sprintf("恭喜%v获得了%v权限！领取魅力值到商场兑换心仪的礼品！", player.User.Nickname, str_s)))

	msg := &pb.MsgUpdateVipTask{}
	f, ok := domainUser.GetUserFortuneManager().GetUserFortune(userId)
	if ok {
		vip_configs := config.GetVipPriceConfigManager().GetVipConfig()
		for _, value := range f.VipTaskStates {
			sub_msg := &pb.UserVipTaskDef{}
			sub_msg.Level = proto.Int(value.VipTaskId)
			if value.StartTime != 0 {
				edTime := value.StartTime + int64(vip_configs[value.VipTaskId].PrizeDays*86400)
				if edTime < util.GetDayZero() {
					sub_msg.State = proto.Int(0)
					sub_msg.StartTime = proto.Int64(0)
					sub_msg.EndTime = proto.Int64(0)
				} else if value.LastGainTime != util.GetDayZero() {
					sub_msg.State = proto.Int(2)
					sub_msg.StartTime = proto.Int64(value.StartTime)
					sub_msg.EndTime = proto.Int64(value.StartTime + int64(vip_configs[value.VipTaskId].PrizeDays*86400))
				} else {
					sub_msg.State = proto.Int(1)
					sub_msg.StartTime = proto.Int64(value.StartTime)
					sub_msg.EndTime = proto.Int64(value.StartTime + int64(vip_configs[value.VipTaskId].PrizeDays*86400))
				}
			} else {
				sub_msg.State = proto.Int(0)
				sub_msg.StartTime = proto.Int64(0)
				sub_msg.EndTime = proto.Int64(0)
			}
			sub_msg.PrizeGold = proto.Int(vip_configs[value.VipTaskId].PrizeGold)
			msg.VipTaskList = append(msg.VipTaskList, sub_msg)
		}
		domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_UPDATE_VIP_TASK), msg)
	}
}
