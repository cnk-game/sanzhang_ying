package admin

import (
	"config"
	"encoding/json"
	"fmt"
	domainPay "game/domain/pay"
	domainUser "game/domain/user"
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"io/ioutil"
	"net/http"
	"strconv"
)

type PayActiveReq struct {
	UserId interface{} `json:"userId"`
	Coins  interface{} `json:"coins"`
}

func PayActiveAddCoinsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("key") != config.ActiveKey {
		glog.Error("key不符Addr:", r.RemoteAddr)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		glog.Error(err)
		return
	}

	res := &PayActiveReq{}
	err = json.Unmarshal(b, res)

	userId := fmt.Sprintf("%v", res.UserId)
	coins := fmt.Sprintf("%v", res.Coins)
	glog.Error("PayActiveAddCoinsHandler userId ", userId)
	glog.Error("PayActiveAddCoinsHandler coins ", coins)

	err1, money := domainPay.GetPayCountBetweenTime(userId)
	if err1 != nil {
		w.Write([]byte("用户数据错误"))
		return
	}

	coinsInt, _ := strconv.ParseFloat(coins, 64)

	ttt := money / 50
	if ttt == 1 || ttt == 2 || ttt == 4 || ttt == 10 || ttt == 20 || ttt == 40 || ttt == 100 {
		errFind := domainPay.GetPayAcitveAddCoinsLog(userId, ttt*50)
		if errFind != mgo.ErrNotFound {
			w.Write([]byte("重复领取"))
			return
		}
	} else {
		w.Write([]byte("用户奖励数据错误"))
		return
	}

	coinsTemp := ttt * 50000
	if coinsTemp > 5000000 || coinsTemp < 50000 || coinsTemp != int(coinsInt) {
		w.Write([]byte("用户奖励数据错误"))
		return
	}

	domainPay.SavePayAcitveAddCoinsLog(userId, ttt*50, coinsTemp)

	domainUser.GetUserFortuneManager().EarnFortune(userId, int64(coinsTemp), 0, 0, true, "充值活动奖励")
	domainUser.GetUserFortuneManager().SaveUserFortune(userId)
	domainUser.GetUserFortuneManager().UpdateUserFortune2(userId, 0)

	w.Write([]byte(`success`))
}
