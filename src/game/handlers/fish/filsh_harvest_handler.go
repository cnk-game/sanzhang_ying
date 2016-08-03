package fish

import (
    "github.com/golang/glog"
    "net/http"
	domainUser "game/domain/user"
    "encoding/json"
	"time"
	"config"
)

type JsonHarvestInfo struct {
    RecordId        string
    Count           int
}

type HarvestRetInfo struct {
    IsOk      string
    Reason    string
}

func FishHarvestHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    w.Header().Set("Access-Control-Allow-Origin", "*")

    token := r.FormValue("access_token")
    userId := r.FormValue("user_id")
    json_obj := r.PostFormValue("obj")

    retInfo := HarvestRetInfo{}
    retInfo.IsOk = "failed"

	if userId == "" || token == "" || json_obj == "" {
	    glog.Error("FishHarvestHandler Param Error.")
	    retInfo.Reason = "接口参数错误"
        ret_str, _ := json.Marshal(retInfo)
		w.Write([]byte(ret_str))
		return
	}

    glog.Info("FishHarvestHandler in. json_obj=", json_obj)
	obj := JsonHarvestInfo{}
	err := json.Unmarshal([]byte(json_obj), &obj)
    if err != nil {
        glog.Error(err)
	    retInfo.Reason = "接口参数格式错误"
        ret_str, _ := json.Marshal(retInfo)
        w.Write([]byte(ret_str))
        return
    }

	f, ok := domainUser.GetUserFortuneManager().GetUserFortune(userId)
	if !ok {
        glog.Error("GetUserFortune error.")
	    retInfo.Reason = "获取用户信息失败"
        ret_str, _ := json.Marshal(retInfo)
        w.Write([]byte(ret_str))
        return
    }

    if f.FishToken.Token != token || int(time.Now().Unix()) - f.FishToken.TimeStamp > 7200 {
        glog.Infof("Token error. %s, %d", f.FishToken.Token, f.FishToken.TimeStamp)
	    retInfo.Reason = "token已过期"
        ret_str, _ := json.Marshal(retInfo)
        w.Write([]byte(ret_str))
        return
    }

    fishInfo, _ok := f.FishInfo[obj.RecordId]
    if !_ok {
        glog.Infof("recordId error. %d", obj.RecordId)
	    retInfo.Reason = "获取用户信息失败"
        ret_str, _ := json.Marshal(retInfo)
        w.Write([]byte(ret_str))
        return
    }

    if fishInfo.Count < obj.Count {
        glog.Infof("count error. %d, %d", fishInfo.Count, obj.Count)
	    retInfo.Reason = "鱼数量不足，无法收获"
        ret_str, _ := json.Marshal(retInfo)
        w.Write([]byte(ret_str))
        return
    }

    Cfish, _ := config.GetFishConfigManager().GetFishConfig(fishInfo.FishType)
    price := int64(Cfish.Price * obj.Count)

    ok = domainUser.GetUserFortuneManager().DelFish(userId, obj.RecordId, obj.Count)
    if ok {
        _, ok = domainUser.GetUserFortuneManager().EarnGold(userId, price, "收获")
        if !ok {
            glog.Infof("EarnGold error.")
            retInfo.Reason = "系统错误，请稍候尝试"
            ret_str, _ := json.Marshal(retInfo)
            w.Write([]byte(ret_str))
            return
        }
        domainUser.GetUserFortuneManager().SaveUserFortune(userId)
    } else {
        glog.Infof("del fish error.")
        retInfo.Reason = "系统错误，请稍候尝试"
        ret_str, _ := json.Marshal(retInfo)
        w.Write([]byte(ret_str))
        return
    }

    retInfo.IsOk = "success"
    ret_str, _ := json.Marshal(retInfo)
    w.Write([]byte(ret_str))
}