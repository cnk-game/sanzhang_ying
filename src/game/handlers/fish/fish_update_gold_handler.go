package fish

import (
    "github.com/golang/glog"
    "net/http"
    "encoding/json"
	domainUser "game/domain/user"
	"time"
)

type UpdateGoldInfo struct {
    Gold    int64
    Reason  string
}

func FishUpdateGoldHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()

    token := r.FormValue("access_token")
    userId := r.FormValue("user_id")

    retInfo := UpdateGoldInfo{}

	if userId == "" || token == "" {
	    glog.Error("FishUpdateGoldHandler Param Error.", userId, token)
	    retInfo.Reason = "接口参数错误"
	    retInfo.Gold = -1
        ret_str, _ := json.Marshal(retInfo)
		w.Write([]byte(ret_str))
		return
	}

	f, ok := domainUser.GetUserFortuneManager().GetUserFortune(userId)
	if !ok {
        glog.Error("GetUserFortune error.")
	    retInfo.Reason = "获取用户信息失败"
	    retInfo.Gold = -1
        ret_str, _ := json.Marshal(retInfo)
        w.Write([]byte(ret_str))
        return
    }

    if f.FishToken.Token != token || int(time.Now().Unix()) - f.FishToken.TimeStamp > 7200 {
        glog.Infof("Token error. %s, %d", f.FishToken.Token, f.FishToken.TimeStamp)
	    retInfo.Reason = "token已过期"
	    retInfo.Gold = -1
        ret_str, _ := json.Marshal(retInfo)
        w.Write([]byte(ret_str))
        return
    }

    retInfo.Gold = f.Gold
    ret_str, _ := json.Marshal(retInfo)
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(ret_str))
}