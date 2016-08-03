package pay

import (
	"code.google.com/p/goprotobuf/proto"
	"crypto/md5"
	"fmt"
	"game/domain/offlineMsg"

	domainUser "game/domain/user"
	"github.com/golang/glog"
	"io"
	"net/http"
	"net/url"
	"pb"
	"strconv"
	"sync"
)

const (
	haimaAppId   = "d577e975d50234da1b409ef77a4038ad"
	haimaAppKey  = "852027b9f928a69b201e265df73fb050"
	exchangeRate = 1
)

var hmPayMu sync.RWMutex

func HmPayHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	glog.Info("====>海马:", r.Form)

	notify_time := r.FormValue("notify_time")
	appid := r.FormValue("appid")
	out_trade_no := r.FormValue("out_trade_no")
	total_fee := r.FormValue("total_fee")
	subject := r.FormValue("subject")
	body := r.FormValue("body")
	trade_status := r.FormValue("trade_status")
	sign := r.FormValue("sign")
	//userId := r.FormValue("user_param")

	if appid != haimaAppId {
		glog.Error("海马appid不存在appid:", appid)
		return
	}

	signStr := fmt.Sprintf("notify_time=%v&appid=%v&out_trade_no=%v&total_fee=%v&subject=%v&body=%v&trade_status=%v%v",
		url.QueryEscape(notify_time), url.QueryEscape(appid), url.QueryEscape(out_trade_no), url.QueryEscape(total_fee), url.QueryEscape(subject),
		url.QueryEscape(body), url.QueryEscape(trade_status), haimaAppKey)

	h := md5.New()
	io.WriteString(h, signStr)
	localSign := fmt.Sprintf("%x", h.Sum(nil))
	if localSign != sign {
		glog.Info("==>签名验证失败local:", localSign, " sign:", sign)
		return
	}

	if trade_status != "1" {
		glog.Info("==>海马失败订单form:", r.Form)
		return
	}

	coins, err := strconv.ParseFloat(total_fee, 64)
	if err != nil {
		glog.Info("total_fee无效:", total_fee, " signStr:", signStr)
		return
	}

	amount := int(coins * exchangeRate)

	payType := fmt.Sprintf("%v", "海马平台")
	money := fmt.Sprintf("%v", amount)
	commonPay(out_trade_no, "hmios", money, payType, "success", body)

	//	ok := func() bool {
	//		hmPayMu.Lock()
	//		defer hmPayMu.Unlock()

	//		// 充值成功，记录充值日志
	//		l := &domainPay.HmPayLog{}
	//		l, err := domainPay.FindHmPayLog(out_trade_no)
	//		if err == nil {
	//			glog.Info("==>订单已处理orderId:", l.OutTradeNo)
	//			return false
	//		}

	//		l.NotifyTime = notify_time
	//		l.AppId = appid
	//		l.UserId = userId
	//		l.OutTradeNo = out_trade_no
	//		l.TotalFee = total_fee
	//		l.Subject = subject
	//		l.Body = body
	//		l.TradeStatus = trade_status

	//		domainPay.SaveHmPayLog(l)

	//		// 支付日志
	//		payLog := &domainPay.PayLog{}
	//		payLog.OrderId = l.OutTradeNo
	//		payLog.UserId = userId
	//		payLog.Amount = int(coins)

	//		u, err := domainUser.FindByUserId(userId)
	//		if err == nil && u != nil {
	//			payLog.Channel = u.ChannelId
	//		}
	//		payLog.PayChannel = "海马平台"
	//		payLog.PayType = "海马平台"
	//		domainPay.SavePayLog(payLog)

	//		return true
	//	}()
	//	if !ok {
	//		glog.Info("==>订单已处理orderId:", out_trade_no, "忽略返回!")
	//		return
	//	}

	//	fmt.Fprintf(w, `success`)

	//	domainUser.GetUserFortuneManager().EarnFortune(userId, 0, int(amount), 0, true, "海马充值")
	//	domainUser.GetUserFortuneManager().SaveUserFortune(userId)

	//	domainUser.GetUserFortuneManager().UpdateUserFortune2(userId, int(amount))

	//	domainUser.GetUserFortuneManager().UpdateRechargeRankingList(userId, int(amount))
	//	//domainUser.GetUserFortuneManager().UpdateVipLevel(userId, int(amount))

	//	// 商场日志
	//	shopLogMsg := &pb.MsgGetShopLogRes{}

	//	log := &domainUser.UserShopLog{}
	//	log.UserId = userId
	//	log.RechargeDiamond = int(amount)
	//	log.Time = time.Now()
	//	domainUser.SaveShopLog(log)

	//	shopLogMsg.LogList = append(shopLogMsg.LogList, log.BuildMessage())

	//	domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_GET_SHOP_LOG), shopLogMsg)

	//	updateRechargeDiamond(userId, int(amount))
}

func updateRechargeDiamond(userId string, diamond int) {
	if domainUser.GetPlayerManager().IsOnline(userId) {
		return
	}

	msg := &pb.MsgUpdateRechargeDiamond{}
	msg.Diamond = proto.Int(diamond)

	offlineMsg.PutOfflineMsg(userId, int32(pb.ServerMsgId_MQ_UPDATE_RECHAGE_DIAMOND), msg)
}
