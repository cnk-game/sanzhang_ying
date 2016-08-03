package pay

import (
	"code.google.com/p/goprotobuf/proto"
	"config"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	activeUser "game/domain/active"
	newUserTask "game/domain/newusertask"
	domainPay "game/domain/pay"
	domainUser "game/domain/user"
	"github.com/golang/glog"
	"io"
	"io/ioutil"
	"net/http"
	"pb"
	"strconv"
	"strings"
	"sync"
	"time"
	"util"
)

const (
	qfAppKey = "5af5e15900b90320d5d199cae0eee87b"
)

var qfPayMu sync.RWMutex

type QfRes struct {
	AppId     interface{} `json:"appId"`
	UserId    interface{} `json:"userId"`
	Order     interface{} `json:"order"`
	Price     interface{} `json:"price"`
	PayType   interface{} `json:"payType"`
	PayCode   interface{} `json:"payCode"`
	State     interface{} `json:"state"`
	Time      interface{} `json:"time"`
	GameOrder interface{} `json:"gameOrder"`
	Sign      interface{} `json:"sign"`
}

func QfPayHandler(w http.ResponseWriter, r *http.Request) {
	glog.Info("QfPayHandler in")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		glog.Error(err)
		return
	}

	jsonData, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		glog.Error(err)
		return
	}

	res := &QfRes{}
	err = json.Unmarshal(jsonData, res)
	glog.Info("起凡充值:", string(jsonData))
	if err != nil {
		glog.Error("解析起凡充值数据失败err:", err)
		return
	}

	w.Write([]byte(`success`))

	order := fmt.Sprintf("%v", res.Order)
	price := fmt.Sprintf("%v", res.Price)

	_, err = strconv.ParseFloat(price, 64)
	if err != nil {
		glog.Info("===>起凡支付price参数解析失败:", r.Form)
		return
	}

	payCode := fmt.Sprintf("%v", res.PayCode)
	state := fmt.Sprintf("%v", res.State)
	gameOrder := fmt.Sprintf("%v", res.GameOrder)
	if len(gameOrder) == 0 {
		glog.Info("===>起凡充值gameOrder无效:", gameOrder)
		return
	}
	sign := fmt.Sprintf("%v", res.Sign)

	qfUserId := fmt.Sprintf("%v", res.UserId)

	md5String := qfUserId + payCode + order + qfAppKey
	calcSign := calcMd5(md5String)
	if calcSign != sign {
		glog.Info("====>签名检验失败sign:", sign, " calcSign:", calcSign)
		return
	}

	if state != "success" {
		glog.Info("===>失败交易:", r.Form)
		return
	}

	if payCode == PRODUCT_VIP1 || payCode == PRODUCT_VIP2 || payCode == PRODUCT_VIP3 || payCode == PRODUCT_VIP4 || payCode == IOS_VIP_25_VIP1 || payCode == IOS_VIP_898_VIP4 {
		QfBuyVIP(res, 0)
	} else if payCode == QUICK_PAY10 || payCode == QUICK_PAY20 || payCode == QUICK_PAY50 || payCode == QUICK_PAY500 || payCode == FIRST_PAY30 || payCode == IOS_QUICK_12 || payCode == IOS_QUICK_18 || payCode == IOS_QUICK_518 {
		QfQuickPay(res)
	} else if payCode == ACTIVE_10 || payCode == ACTIVE_30 {
		ActivePay(res)
	} else {
		QfBuyDiamond(res)
	}
}

func ActivePay(res *QfRes) {
	glog.Info("ActivePay in")
	appId := fmt.Sprintf("%v", res.AppId)
	order := fmt.Sprintf("%v", res.Order)
	price := fmt.Sprintf("%v", res.Price)

	amount, _ := strconv.ParseFloat(price, 64)

	payType := fmt.Sprintf("%v", res.PayType)
	payCode := fmt.Sprintf("%v", res.PayCode)
	state := fmt.Sprintf("%v", res.State)
	qfTime := fmt.Sprintf("%v", res.Time)
	gameOrder := fmt.Sprintf("%v", res.GameOrder)
	cok, userId := domainPay.GetTokenManager().CheckToken(gameOrder)
	if !cok {
		if strings.Contains(gameOrder, "QF") {
			userId = gameOrder
		} else {
			glog.Info("==>token error:", gameOrder, "忽略返回!")
			return
		}
	}

	ok := func() bool {
		qfPayMu.Lock()
		defer qfPayMu.Unlock()

		// 充值成功，记录充值日志
		l := &domainPay.QfPayLog{}
		l, err := domainPay.FindQfPayLog(order)
		if err == nil {
			glog.Info("==>订单已处理orderId:", l.Order)
			return false
		}

		l.AppId = appId
		l.UserId = userId
		l.Order = order
		l.Price = price
		l.PayType = payType
		l.PayCode = payCode
		l.State = state
		l.Time = qfTime
		l.GameOrder = gameOrder

		err = domainPay.SaveQfPayLog(l)
		if err != nil {
			glog.Error("保存充值记录失败err:", err)
		}

		// 支付日志
		payLog := &domainPay.PayLog{}
		payLog.OrderId = l.Order
		payLog.UserId = userId
		payLog.Amount = int(amount)
		payLog.PayCode = payCode

		u, err := domainUser.FindByUserId(userId)
		if err == nil && u != nil {
			payLog.Channel = u.ChannelId
		}
		payLog.PayChannel = "起凡互娱"
		payLog.PayType = payType
		domainPay.SavePayLog(payLog)

		return true
	}()
	glog.Info("ok=", ok)
	if !ok {
		glog.Info("==>订单已处理orderId:", order, "忽略返回!")
		return
	}

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
	CheckNewUserTask(userId)
}

// add by wangsq start
func QfQuickPay(res *QfRes) {
	glog.Info("QfQuickPay in")
	appId := fmt.Sprintf("%v", res.AppId)
	order := fmt.Sprintf("%v", res.Order)
	price := fmt.Sprintf("%v", res.Price)

	amount, _ := strconv.ParseFloat(price, 64)

	payType := fmt.Sprintf("%v", res.PayType)
	payCode := fmt.Sprintf("%v", res.PayCode)
	state := fmt.Sprintf("%v", res.State)
	qfTime := fmt.Sprintf("%v", res.Time)
	gameOrder := fmt.Sprintf("%v", res.GameOrder)
	cok, userId := domainPay.GetTokenManager().CheckToken(gameOrder)
	if !cok {
		if strings.Contains(gameOrder, "QF") {
			userId = gameOrder
		} else {
			glog.Info("==>token error:", gameOrder, "忽略返回!")
			return
		}
	}

	ok := func() bool {
		qfPayMu.Lock()
		defer qfPayMu.Unlock()

		// 充值成功，记录充值日志
		l := &domainPay.QfPayLog{}
		l, err := domainPay.FindQfPayLog(order)
		if err == nil {
			glog.Info("==>订单已处理orderId:", l.Order)
			return false
		}

		l.AppId = appId
		l.UserId = userId
		l.Order = order
		l.Price = price
		l.PayType = payType
		l.PayCode = payCode
		l.State = state
		l.Time = qfTime
		l.GameOrder = gameOrder

		err = domainPay.SaveQfPayLog(l)
		if err != nil {
			glog.Error("保存充值记录失败err:", err)
		}

		// 支付日志
		payLog := &domainPay.PayLog{}
		payLog.OrderId = l.Order
		payLog.UserId = userId
		payLog.Amount = int(amount)
		payLog.PayCode = payCode

		u, err := domainUser.FindByUserId(userId)
		if err == nil && u != nil {
			payLog.Channel = u.ChannelId
		}
		payLog.PayChannel = "起凡互娱"
		payLog.PayType = payType
		domainPay.SavePayLog(payLog)

		return true
	}()
	glog.Info("ok=", ok)
	if !ok {
		glog.Info("==>订单已处理orderId:", order, "忽略返回!")
		return
	}

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
	domainUser.GetUserFortuneManager().EarnFortune(userId, int64(gold), 0, 0, true, "起凡充值")
	domainUser.GetUserFortuneManager().SaveUserFortune(userId)
	domainUser.GetUserFortuneManager().UpdateRechargeRankingList(userId, int(amount))
	if payCode == FIRST_PAY30 {
		domainUser.GetUserFortuneManager().OpenSafeBox(userId, "")
		domainUser.GetUserFortuneManager().UpdateFirstRecharge(userId)
	}

	domainUser.GetUserFortuneManager().UpdateUserFortune2(userId, 0)

	CheckNewUserTask(userId)
}

func QfBuyDiamond(res *QfRes) {
	glog.Info("QfBuyDiamond in")
	appId := fmt.Sprintf("%v", res.AppId)
	order := fmt.Sprintf("%v", res.Order)
	price := fmt.Sprintf("%v", res.Price)

	amount, _ := strconv.ParseFloat(price, 64)

	payType := fmt.Sprintf("%v", res.PayType)
	payCode := fmt.Sprintf("%v", res.PayCode)
	state := fmt.Sprintf("%v", res.State)
	qfTime := fmt.Sprintf("%v", res.Time)
	gameOrder := fmt.Sprintf("%v", res.GameOrder)
	cok, userId := domainPay.GetTokenManager().CheckToken(gameOrder)
	if !cok {
		if strings.Contains(gameOrder, "QF") {
			userId = gameOrder
		} else {
			glog.Info("==>token error:", gameOrder, "忽略返回!")
			return
		}
	}

	ok := func() bool {
		qfPayMu.Lock()
		defer qfPayMu.Unlock()

		// 充值成功，记录充值日志
		l := &domainPay.QfPayLog{}
		l, err := domainPay.FindQfPayLog(order)
		if err == nil {
			glog.Info("==>订单已处理orderId:", l.Order)
			return false
		}

		l.AppId = appId
		l.UserId = userId
		l.Order = order
		l.Price = price
		l.PayType = payType
		l.PayCode = payCode
		l.State = state
		l.Time = qfTime
		l.GameOrder = gameOrder

		err = domainPay.SaveQfPayLog(l)
		if err != nil {
			glog.Error("保存充值记录失败err:", err)
		}

		// 支付日志
		payLog := &domainPay.PayLog{}
		payLog.OrderId = l.Order
		payLog.UserId = userId
		payLog.Amount = int(amount)
		payLog.PayCode = payCode

		u, err := domainUser.FindByUserId(userId)
		if err == nil && u != nil {
			payLog.Channel = u.ChannelId
		}
		payLog.PayChannel = "起凡互娱"
		payLog.PayType = payType
		domainPay.SavePayLog(payLog)

		return true
	}()

	if !ok {
		glog.Info("==>订单已处理orderId:", order, "忽略返回!")
		return
	}

	domainUser.GetUserFortuneManager().EarnFortune(userId, 0, int(amount), 0, true, "起凡充值")
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
	CheckNewUserTask(userId)
}

func QfBuyVIP(res *QfRes, channel int) {
	glog.Info("QfBuyVIP in")
	appId := fmt.Sprintf("%v", res.AppId)
	order := fmt.Sprintf("%v", res.Order)
	price := fmt.Sprintf("%v", res.Price)

	amount, _ := strconv.ParseFloat(price, 64)

	payType := fmt.Sprintf("%v", res.PayType)
	payCode := fmt.Sprintf("%v", res.PayCode)
	state := fmt.Sprintf("%v", res.State)
	qfTime := fmt.Sprintf("%v", res.Time)
	gameOrder := fmt.Sprintf("%v", res.GameOrder)
	cok, userId := domainPay.GetTokenManager().CheckToken(gameOrder)
	if !cok {
		if strings.Contains(gameOrder, "QF") {
			userId = gameOrder
		} else {
			glog.Info("==>token error:", gameOrder, "忽略返回!")
			return
		}
	}

	ok := func() bool {
		qfPayMu.Lock()
		defer qfPayMu.Unlock()

		// 充值成功，记录充值日志
		l := &domainPay.QfPayLog{}
		l, err := domainPay.FindQfPayLog(order)
		if err == nil {
			glog.Info("==>订单已处理orderId:", l.Order)
			return false
		}

		l.AppId = appId
		l.UserId = userId
		l.Order = order
		l.Price = price
		l.PayType = payType
		l.PayCode = payCode
		l.State = state
		l.Time = qfTime
		l.GameOrder = gameOrder

		err = domainPay.SaveQfPayLog(l)
		if err != nil {
			glog.Error("保存充值记录失败err:", err)
		}

		// 支付日志
		payLog := &domainPay.PayLog{}
		payLog.OrderId = l.Order
		payLog.UserId = userId
		payLog.Amount = int(amount)
		payLog.PayCode = payCode

		u, err := domainUser.FindByUserId(userId)
		if err == nil && u != nil {
			payLog.Channel = u.ChannelId
		}
		payLog.PayChannel = "起凡互娱"
		payLog.PayType = payType
		domainPay.SavePayLog(payLog)

		return true
	}()
	if !ok {
		glog.Info("==>订单已处理orderId:", order, "忽略返回!")
		return
	}

	vip_configs := config.GetVipPriceConfigManager().GetVipConfig()
	charm := 0
	goldNow := 0
	level := 0
	name := ""
	if payCode == PRODUCT_VIP1 {
		charm = vip_configs[1].PrizeCharm
		if channel != 0 {
			charm = 2
		}
		name = vip_configs[1].Name
		goldNow = vip_configs[1].PrizeGoldNow
		level = 1
	} else if payCode == PRODUCT_VIP2 {
		charm = vip_configs[2].PrizeCharm
		if channel != 0 {
			charm = 15
		}
		name = vip_configs[2].Name
		goldNow = vip_configs[2].PrizeGoldNow
		level = 2
	} else if payCode == PRODUCT_VIP3 {
		charm = vip_configs[3].PrizeCharm
		if channel != 0 {
			charm = 25
		}
		name = vip_configs[3].Name
		goldNow = vip_configs[3].PrizeGoldNow
		level = 3
	} else if payCode == PRODUCT_VIP4 {
		charm = vip_configs[4].PrizeCharm
		if channel != 0 {
			charm = 50
		}
		name = vip_configs[4].Name
		goldNow = vip_configs[4].PrizeGoldNow
		level = 4
	} else if payCode == IOS_VIP_25_VIP1 {
		charm = vip_configs[11].PrizeCharm
		if channel != 0 {
			charm = 2
		}
		name = vip_configs[11].Name
		goldNow = vip_configs[11].PrizeGoldNow
		level = 1
	} else if payCode == IOS_VIP_898_VIP4 {
		charm = vip_configs[12].PrizeCharm
		if channel != 0 {
			charm = 50
		}
		name = vip_configs[12].Name
		goldNow = vip_configs[12].PrizeGoldNow
		level = 4
	}

	domainUser.GetUserFortuneManager().EarnFortune(userId, int64(goldNow), 0, 0, true, "起凡VIP购买")
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
	CheckNewUserTask(userId)
}

// add by wangsq end

func calcMd5(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func CheckNewUserTask(userId string) {
	result := newUserTask.GetNewUserTaskManager().CheckUserRechargeTask(userId)
	if result > 0 {
		msgT := &pb.MsgNewbeTaskCompletedNotify{}
		msgT.Id = proto.Int(result)
		domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_NOTIFY_NEWBETASK_COMP), msgT)
	}
}
