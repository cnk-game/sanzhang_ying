package admin

import (
	"encoding/json"
	domainPay "game/domain/pay"
	"github.com/golang/glog"
	"net/http"
)

type ResCount struct {
	UserId string
	Money  int
}

func GetRechargeCountHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	glog.Info("GetRechargeCountHandler:", r.Form)
	userId := r.FormValue("userId")
	err, money := domainPay.GetPayCountBetweenTime(userId)
	if err != nil {
		w.Write([]byte("error"))
		return
	}
	res := ResCount{UserId: userId, Money: money}
	b, err1 := json.Marshal(res)
	if err1 != nil {
		glog.Info("GetRechargeCountHandler Marshal err :", err1)
	}

	w.Write(b)
}
