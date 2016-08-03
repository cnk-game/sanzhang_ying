package admin

import (
	"config"
	"github.com/golang/glog"
	"net/http"
	"fmt"
	domainUser "game/domain/user"
	"strconv"
)

func GetUserCountGoldLimitHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.FormValue("key") != config.ControlKey {
		glog.Error("key不符Addr:", r.RemoteAddr)
		w.Write([]byte(fmt.Sprintf("%v", -1)))
		return
	}

    if r.FormValue("limit_low") == "" || r.FormValue("limit_high") == "" {
        glog.Error("GetUserCountGoldLimitHandler Param Error.")
		w.Write([]byte(fmt.Sprintf("%v", -2)))
        return
    }

    limitLow, _ := strconv.Atoi(r.FormValue("limit_low"))
    limitHigh, _ := strconv.Atoi(r.FormValue("limit_high"))

	count, err := domainUser.FindGoldLimitUserCount(limitLow, limitHigh)
	if err != nil {
        glog.Error("FindGoldLimitUserCount Error.err=", err)
		w.Write([]byte(fmt.Sprintf("%v", -3)))
        return
    }

    w.Write([]byte(fmt.Sprintf("%v", count)))
    return
}

