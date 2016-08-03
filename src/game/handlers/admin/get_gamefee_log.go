package admin

import (
	"github.com/golang/glog"
	"net/http"
	"util"
    "encoding/json"
)

func GetGameFeeLogHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
	start := r.FormValue("start")
	end := r.FormValue("end")
	glog.Info("GetGameFeeLogHandler in,", start, end)

    results := util.MongoLog_GetGameFeeLog(start, end)
    ret_str, _ := json.Marshal(results)

    w.Write([]byte(ret_str))
}
