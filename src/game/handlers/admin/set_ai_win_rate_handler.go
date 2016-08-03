package admin

import (
	"config"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
)

type CardConfigData struct {
	GameType    int `json:"gameType"`
	Single      int `json:"single"`
	Double      int `json:"double"`
	ShunZi      int `json:"shuZi"`
	JinHua      int `json:"jinHua"`
	ShunJin     int `json:"shunJin"`
	BaoZi       int `json:"baoZi"`
	WinGold     int `json:"winGold"`
	LoseGold    int `json:"loseGold"`
	WinRateHigh int `json:"winRateHigh"`
	WinRateLow  int `json:"winRateLow"`
}

func SetCardConfigDataHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("key====>>:", r.FormValue("key"))
	if r.URL.Query().Get("key") != config.ControlKey {
		glog.Error("key不符Addr:", r.RemoteAddr)
		return
	}

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		glog.Error(err)
		return
	}

	fmt.Println("===>AI牌型胜率:", string(b))

	data := &CardConfigData{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		glog.Error(err)
		w.Write([]byte(`0`))
		return
	}

	if data.Single+data.Double+data.ShunZi+data.JinHua+data.ShunJin+data.BaoZi == 0 {
		glog.Info("===>错误配置，概率和为0:", data)
		w.Write([]byte(`0`))
		return
	}

	config.GetCardConfigManager().SetCardConfig(data.GameType, data.Single, data.Double, data.ShunZi, data.JinHua, data.ShunJin, data.BaoZi, data.WinGold, data.LoseGold, data.WinRateHigh, data.WinRateLow)

	w.Write([]byte(`1`))
}
