package fish

import (
    "github.com/golang/glog"
    "net/http"
    "encoding/json"
	domainUser "game/domain/user"
	"time"
	"config"
)

type JsonFishInfo struct {
    RecordId        string
    FishType        int
    Url             string
    FishName        string
    Count           int
    GrowUp          int
    FromUid         string
    GetTime         time.Time
    Price           int
}

func FishUserInfoHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()

	token := r.FormValue("access_token")
	userId := r.FormValue("user_id")

	w.Write([]byte(`0`))
	return

	if userId == "" || token == "" {
	    glog.Error("FishUserInfoHandler Param Error.")
		w.Write([]byte(`0`))
		return
	}

	f, ok := domainUser.GetUserFortuneManager().GetUserFortune(userId)
	if !ok {
        glog.Error("GetUserFortune error.")
        w.Write([]byte(`0`))
        return
    }

    if f.FishToken.Token != token || int(time.Now().Unix()) - f.FishToken.TimeStamp > 7200 {
        glog.Infof("Token error. %s, %d", f.FishToken.Token, f.FishToken.TimeStamp)
        w.Write([]byte(`0`))
        return
    }

    ret := []JsonFishInfo{}
    for _, value := range f.FishInfo {
        Jfish := JsonFishInfo{}
        Cfish, _ := config.GetFishConfigManager().GetFishConfig(value.FishType)

        Jfish.RecordId = value.RecordId
        Jfish.FishType = value.FishType
        Jfish.Url = Cfish.Url
        Jfish.FishName = Cfish.FishName
        Jfish.Count = value.Count
        Jfish.GrowUp = Cfish.GrowUp
        Jfish.FromUid = value.FromUid
        Jfish.GetTime = value.GetTime
        Jfish.Price = Cfish.Price
        ret = append(ret, Jfish)
    }

    ret_str, err := json.Marshal(ret)
    if err != nil {
        glog.Error(err)
    } else {
        glog.V(2).Info("ret_str = %s", ret_str)
    }

    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(ret_str))
}
