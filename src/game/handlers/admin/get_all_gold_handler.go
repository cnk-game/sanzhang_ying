package admin

import (
	"config"
	"github.com/golang/glog"
	"net/http"
	"fmt"
	domainUser "game/domain/user"
)

func GetAllGoldHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.FormValue("key") != config.ControlKey {
		glog.Error("key不符Addr:", r.RemoteAddr)
		w.Write([]byte(fmt.Sprintf("%v", -1)))
		return
	}

	all, err := domainUser.SumAllGolds()
	if err != nil {
        glog.Error("SumAllGolds Error.err=", err)
		w.Write([]byte(fmt.Sprintf("%v", -2)))
        return
    }

    w.Write([]byte(fmt.Sprintf("%v", all)))
    return
}

