package pay

import (
	"config"
	"encoding/json"
	"fmt"
	domainPay "game/domain/pay"
	domainUser "game/domain/user"
	"github.com/golang/glog"
	"net/http"
	"pb"
	"sync"
	"time"
)

type IAppPayTransdata struct {
	Transtype int     `json:"transtype"`
	Cporderid string  `json:"cporderid"`
	Transid   string  `json:"transid"`
	Appuserid string  `json:"appuserid"`
	Appid     string  `json:"appid"`
	Waresid   int     `json:"waresid"`
	Feetype   int     `json:"feetype"`
	Money     float32 `json:"money"`
	Currency  string  `json:"currency"`
	Result    int     `json:"result"`
	Transtime string  `json:"transtime"`
	Cpprivate string  `json:"cpprivate"`
	Paytype   int     `json:"paytype"`
}

var iappPayMu sync.RWMutex

func IAppPayHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	glog.Info("爱贝充值:", r.Form)

	if r.FormValue("key") != config.ControlKey {
		glog.Error("key不符Addr:", r.RemoteAddr)
		return
	}

	data := r.FormValue("transdata")
	transdata := &IAppPayTransdata{}
	err := json.Unmarshal([]byte(data), transdata)
	if err != nil {
		glog.Error(err)
		return
	}

	glog.Infof("transdata:%#v\n", transdata)

	if transdata.Result != 0 {
		// 交易失败
		return
	}

	if transdata.Appuserid == "" || transdata.Money < 1 {
		glog.Info("====>爱贝充值失败appuserid为空或money小于1")
		return
	}

	ok := func() bool {
		iappPayMu.Lock()
		defer iappPayMu.Unlock()

		l, err := domainPay.FindIAppPayLog(transdata.Transid)
		glog.Info("===>查询充值Log:", l, " err:", err)
		if err == nil {
			if l.Result == 0 {
				// 此交易已处理
				glog.V(2).Info("transId:", transdata.Transid, "已处理，重复通知!")
				return false
			}
		}
		l.Transtype = transdata.Transtype
		l.Cporderid = transdata.Cporderid
		l.Transid = transdata.Transid
		l.Appuserid = transdata.Appuserid
		l.Appid = transdata.Appid
		l.Waresid = transdata.Waresid
		l.Feetype = transdata.Feetype
		l.Money = transdata.Money
		l.Currency = transdata.Currency
		l.Result = transdata.Result
		l.Transtime = transdata.Transtime
		l.Cpprivate = transdata.Cpprivate
		l.Paytype = transdata.Paytype

		err = domainPay.SaveIAppPayLog(l)
		if err != nil {
			glog.Info("保存爱贝交易日志失败err:", err, " log:", l)
		}

		// 支付日志
		payLog := &domainPay.PayLog{}
		payLog.OrderId = transdata.Transid
		payLog.UserId = transdata.Appuserid
		payLog.Amount = int(transdata.Money)

		u, err := domainUser.FindByUserId(transdata.Appuserid)
		if err == nil && u != nil {
			payLog.Channel = u.ChannelId
		}
		payLog.PayChannel = "爱贝"
		payLog.PayType = fmt.Sprintf("%v", transdata.Paytype)
		domainPay.SavePayLog(payLog)

		return true
	}()
	if !ok {
		glog.Info("===>无效交易，忽略不处理transId:", transdata.Transid)
		return
	}

	domainUser.GetUserFortuneManager().EarnFortune(transdata.Appuserid, 0, int(transdata.Money), 0, true, "爱贝充值")
	domainUser.GetUserFortuneManager().SaveUserFortune(transdata.Appuserid)

	domainUser.GetUserFortuneManager().UpdateUserFortune2(transdata.Appuserid, int(transdata.Money))

	domainUser.GetUserFortuneManager().UpdateRechargeRankingList(transdata.Appuserid, int(transdata.Money))
	//domainUser.GetUserFortuneManager().UpdateVipLevel(transdata.Appuserid, int(transdata.Money))

	// 商场日志
	shopLogMsg := &pb.MsgGetShopLogRes{}
	l := &domainUser.UserShopLog{}

	l.UserId = transdata.Appuserid
	l.RechargeDiamond = int(transdata.Money)
	l.Time = time.Now()
	domainUser.SaveShopLog(l)

	shopLogMsg.LogList = append(shopLogMsg.LogList, l.BuildMessage())

	domainUser.GetPlayerManager().SendClientMsg(transdata.Appuserid, int32(pb.MessageId_GET_SHOP_LOG), shopLogMsg)

	updateRechargeDiamond(transdata.Appuserid, int(transdata.Money))
}
