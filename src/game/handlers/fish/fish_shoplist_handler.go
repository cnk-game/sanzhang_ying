package fish

import (
    "github.com/golang/glog"
    "net/http"
    "config"
    "encoding/json"
)

func FishShopListHandler(w http.ResponseWriter, r *http.Request) {
    fish_infos := config.GetFishConfigManager().GetFishAll()
    ret_str, err := json.Marshal(fish_infos)
    if err != nil {
        glog.Error(err)
    } else {
        glog.V(2).Info("ret_str = %s", ret_str)
    }
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Write([]byte(ret_str))
}