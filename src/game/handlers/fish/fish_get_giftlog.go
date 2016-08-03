package fish

import (
    "github.com/golang/glog"
    "net/http"
    "encoding/json"
	domainUser "game/domain/user"
	mgo "gopkg.in/mgo.v2"
)

type JsonFishLog struct {
    Code     int
    Logs     []*domainUser.GiftFishLog
}

func FishGetGiftLogHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()

    token := r.FormValue("access_token")
    userId := r.FormValue("user_id")

    ret := JsonFishLog{}
    ret.Code = 0

	if userId == "" || token == "" {
	    glog.Error("FishUserInfoHandler Param Error.")
	    ret.Code = 1
	    ret_str, _ := json.Marshal(ret)
		w.Write([]byte(ret_str))
		return
	}

    logs, err := domainUser.LoadGiftFishLog(userId)
    ret.Logs = logs
    if err != nil && err != mgo.ErrNotFound {
	    glog.Error("LoadGiftFishLog Error.userId=", userId)
	    ret.Code = 2
	    ret_str, _ := json.Marshal(ret)
		w.Write([]byte(ret_str))
		return
    }

    ret_str, _ := json.Marshal(ret)

    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(ret_str))
}