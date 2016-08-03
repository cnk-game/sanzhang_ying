package admin

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	domainPay "game/domain/pay"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
)

type ChangePayType struct {
	Key  interface{} `json:"key"`
	Open interface{} `json:"open"`
}

func ChangeIpsPayTypeHandler(w http.ResponseWriter, r *http.Request) {
	glog.Info("ChangeIpsPayType in")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		glog.Error(err)
		return
	}

	jsonData, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		glog.Error(err)
		return
	}

	res := &ChangePayType{}
	err = json.Unmarshal(jsonData, res)
	glog.Info("ChangeIpsPayTypeHandler:", string(jsonData))
	if err != nil {
		glog.Error("ChangeIpsPayTypeHandler:", err)
		return
	}

	key := fmt.Sprintf("%v", res.Key)
	isOpen := fmt.Sprintf("%v", res.Open)
	glog.Info("ChangeIpsPayTypeHandler open :", isOpen)
	glog.Info("ChangeIpsPayTypeHandler key :", key)
	if key == "qifanthreekey123" {
		w.Write([]byte(isOpen))

		if isOpen == "1" {
			domainPay.IPS_PAY_CHANGE = true
		} else {
			domainPay.IPS_PAY_CHANGE = false
		}
	} else {
		w.Write([]byte(`fail`))
	}
}
