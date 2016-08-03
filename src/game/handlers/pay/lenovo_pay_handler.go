package pay

import (
	"crypto/md5"
	"fmt"
	domainPay "game/domain/pay"
	domainUser "game/domain/user"
	"github.com/golang/glog"
	"io"
	"net/http"
	"pb"
	"sort"
	"strconv"
	"sync"
	"time"
)

const (
	lenovoKey = "681b1bf54d716c677e2638d57cc3d361"
)

var lenovoPayMu sync.RWMutex

func LenovoPayHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	glog.Info("联想充值:", r.Form)

	params := []string{}
	for key := range r.Form {
		if r.FormValue(key) == "null" {
			continue
		}

		if key != "sign" {
			params = append(params, key)
		}
	}

	sort.Strings(params)

	signStr := ""
	for _, key := range params {
		signStr += fmt.Sprintf("%v=%v&", key, r.FormValue(key))
	}
	signStr += fmt.Sprintf("key=%v", lenovoKey)

	h := md5.New()
	io.WriteString(h, signStr)
	localSign := fmt.Sprintf("%x", h.Sum(nil))

	if localSign != r.FormValue("sign") {
		glog.Info("===>联想充值签名验证失败localSign:", localSign, " sign:", r.FormValue("sign"), " addr:", r.RemoteAddr)
		return
	}

	w.Write([]byte(`success`))

	if r.FormValue("status") != "1" {
		// 失败订单
		glog.Error("失败订单status:", r.Form)
		return
	}

	userId := r.FormValue("attach")
	if userId == "" || userId == "null" {
		glog.Error("联想充值失败attach userId为空!")
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err != nil {
		glog.Error("联想充值失败，解析支付金额失败userId:", userId, " amount:", r.FormValue("amount"))
		return
	}

	ok := func() bool {
		lenovoPayMu.Lock()
		defer lenovoPayMu.Unlock()

		// 充值成功，记录充值日志
		l, err := domainPay.FindLenovoPayLog(r.FormValue("order_id"))
		if err == nil && l.Status == "1" {
			glog.Info("==>订单已处理orderId:", l.OrderId)
			return false
		}
		l.OrderId = r.FormValue("order_id")
		l.MerchantOrderId = r.FormValue("merchant_order_id")
		l.Amount = int(amount)
		l.AppId = r.FormValue("app_id")
		l.PayTime = r.FormValue("pay_time")
		l.Attach = userId
		l.Status = r.FormValue("status")
		l.Sign = r.FormValue("sign")
		domainPay.SaveLenovoPayLog(l)

		// 支付日志
		payLog := &domainPay.PayLog{}
		payLog.OrderId = l.OrderId
		payLog.UserId = userId
		payLog.Amount = int(amount)

		u, err := domainUser.FindByUserId(userId)
		if err == nil && u != nil {
			payLog.Channel = u.ChannelId
		}
		payLog.PayChannel = "联想"
		payLog.PayType = "联想"
		domainPay.SavePayLog(payLog)

		return true
	}()
	if !ok {
		glog.Info("==>订单已处理orderId:", r.FormValue("order_id"), "忽略返回!")
		return
	}

	domainUser.GetUserFortuneManager().EarnFortune(userId, 0, int(amount), 0, true, "联想充值")
	domainUser.GetUserFortuneManager().SaveUserFortune(userId)

	domainUser.GetUserFortuneManager().UpdateUserFortune2(userId, int(amount))

	domainUser.GetUserFortuneManager().UpdateRechargeRankingList(userId, int(amount))
	//domainUser.GetUserFortuneManager().UpdateVipLevel(userId, int(amount))

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
