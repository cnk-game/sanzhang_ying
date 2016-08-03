package admin

import (
	"config"
	"fmt"
	"game/domain/user"
	"github.com/golang/glog"
	"net/http"
)

func GetOnlineCountHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.FormValue("key") != config.ControlKey {
		glog.Error("key不符Addr:", r.RemoteAddr)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte(fmt.Sprintf("%v", user.GetPlayerManager().GetOnlineCount())))
}
