package admin

import (
	"config"
	"game/domain/user"
	"github.com/golang/glog"
	"net/http"
	"strconv"
)

func SetUserFortuneHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("key") != config.ControlKey {
		glog.Error("key不符Addr:", r.RemoteAddr)
		return
	}

	userId := r.FormValue("userId")
	if userId == "" {
		w.Write([]byte(`0`))
		return
	}

	diamond, _ := strconv.Atoi(r.FormValue("diamond"))
	gold, _ := strconv.Atoi(r.FormValue("gold"))

	if !user.GetUserFortuneManager().EarnFortune(userId, int64(gold), diamond, 0, true, "管理员设置") {
		w.Write([]byte(`0`))
		return
	}
	if diamond > 0 {
		user.GetUserFortuneManager().UpdateRechargeRankingList(userId, diamond)
	}
	user.GetUserFortuneManager().UpdateUserFortune2(userId, diamond)

	w.Write([]byte(`1`))
}
