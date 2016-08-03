package fish

import (
    "github.com/golang/glog"
    "net/http"
	domainUser "game/domain/user"
    "encoding/json"
	"time"
	mgo "gopkg.in/mgo.v2"
)

type JsonGiftInfo struct {
    FromUid         string
    ToUid           string
    RecordId        string
    Count           int
}

type GiftRetInfo struct {
    IsOk      string
    Reason    string
}


func FishGiftHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    w.Header().Set("Access-Control-Allow-Origin", "*")

    token := r.FormValue("access_token")
    userId := r.FormValue("user_id")
    json_obj := r.PostFormValue("obj")

    retInfo := GiftRetInfo{}
    retInfo.IsOk = "failed"

	if userId == "" || token == "" || json_obj == "" {
	    glog.Error("FishBuyHandler Param Error.")
	    retInfo.Reason = "接口参数错误"
        ret_str, _ := json.Marshal(retInfo)
		w.Write([]byte(ret_str))
		return
	}

    obj := JsonGiftInfo{}
    err := json.Unmarshal([]byte(json_obj), &obj)
    if err != nil {
        glog.Error(err)
	    retInfo.Reason = "接口参数错误"
        ret_str, _ := json.Marshal(retInfo)
        w.Write([]byte(ret_str))
        return
    }

	from_f, f_ok := domainUser.GetUserFortuneManager().GetUserFortune(obj.FromUid)
	_, t_ok := domainUser.GetUserFortuneManager().GetUserFortune(obj.ToUid)
	if !f_ok {
        glog.Error("GetUserFortune error.")
	    retInfo.Reason = "获取用户信息失败"
        ret_str, _ := json.Marshal(retInfo)
        w.Write([]byte(ret_str))
        return
    }

    if from_f.FishToken.Token != token || int(time.Now().Unix()) - from_f.FishToken.TimeStamp > 7200 {
        glog.Infof("Token error. %s, %d", from_f.FishToken.Token, from_f.FishToken.TimeStamp)
	    retInfo.Reason = "token已过期"
        ret_str, _ := json.Marshal(retInfo)
        w.Write([]byte(ret_str))
        return
    }

    from_fishInfo, _ok := from_f.FishInfo[obj.RecordId]
    if !_ok {
        glog.Infof("recordId error. %d", obj.RecordId)
	    retInfo.Reason = "获取用户信息失败"
        ret_str, _ := json.Marshal(retInfo)
        w.Write([]byte(ret_str))
        return
    }

    if from_fishInfo.Count < obj.Count {
        glog.Infof("count error. %d, %d", from_fishInfo.Count, obj.Count)
	    retInfo.Reason = "鱼的数量不足，无法赠送"
        ret_str, _ := json.Marshal(retInfo)
        w.Write([]byte(ret_str))
        return
    }

    FishType := from_fishInfo.FishType

    if t_ok {
        domainUser.GetUserFortuneManager().AddFish(obj.ToUid, obj.FromUid, FishType, obj.Count)
        domainUser.GetUserFortuneManager().SaveUserFortune(obj.ToUid)
    } else {
        err := domainUser.GetUserFortuneManager().AddFish2Mongo(obj.ToUid, obj.FromUid, FishType, obj.Count)
        if err != nil && err == mgo.ErrNotFound {
            retInfo.Reason = "请确认被赠送者用户ID"
            ret_str, _ := json.Marshal(retInfo)
            w.Write([]byte(ret_str))
            return
        } else if err != nil {
            retInfo.Reason = "系统错误，请稍候再试！"
            ret_str, _ := json.Marshal(retInfo)
            w.Write([]byte(ret_str))
            return
        }
    }

    domainUser.GetUserFortuneManager().DelFish(obj.FromUid, obj.RecordId, obj.Count)
    domainUser.GetUserFortuneManager().SaveUserFortune(obj.FromUid)

    domainUser.SaveGiftFishLog(obj.FromUid, obj.ToUid, FishType, obj.Count)
    retInfo.IsOk = "success"
    ret_str, _ := json.Marshal(retInfo)

    w.Write([]byte(ret_str))
}