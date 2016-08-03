package admin

import (
	"config"
	"github.com/golang/glog"
	"net/http"
	"util"
    "encoding/json"
)

type OnlineTypeInfo struct {
    GameType    int    `json:"gameType"`
    Count       int    `json:"count"`
}

type OnlineTypeResp struct {
    Infos    []OnlineTypeInfo    `json:"infos"`
}

func GetOnlineTypeHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.FormValue("key") != config.ControlKey {
		glog.Error("key不符Addr:", r.RemoteAddr)
		return
	}

    counts, err := util.MongoLog_GetPlayerCountByGType()
    if err != nil {
		glog.Error("MongoLog_GetPlayerCountByGType error. err=", err)
		return
    }

    ret := OnlineTypeResp{}
    for _, v := range counts {
        info := OnlineTypeInfo{}
        info.GameType = v.GameType
        info.Count = v.PlayerCount
        ret.Infos = append(ret.Infos, info)
    }
    ret_str, _ := json.Marshal(ret)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write([]byte(ret_str))
}
