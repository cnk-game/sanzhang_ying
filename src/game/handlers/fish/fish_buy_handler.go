package fish

import (
	"encoding/json"
	"net/http"
)


type JsonBuyObj struct {
    FishType        int
    Count           int
}

type BuyRetInfo struct {
    IsOk      string
    Reason    string
}

func FishBuyHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	w.Header().Set("Access-Control-Allow-Origin", "*")

	retInfo := BuyRetInfo{}
	retInfo.IsOk = "failed"
	retInfo.Reason = "获取用户信息失败"
	ret_str, _ := json.Marshal(retInfo)
	w.Write([]byte(ret_str))
	return
}
