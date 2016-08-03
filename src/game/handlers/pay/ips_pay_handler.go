package pay

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"strconv"
)

func IPSPayHandler(w http.ResponseWriter, r *http.Request) {
	glog.Info("IPSPayHandler in")
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

	res := &QfRes{}
	err = json.Unmarshal(jsonData, res)
	glog.Info("IPS:", string(jsonData))
	if err != nil {
		glog.Error("解析IPS充值数据失败err:", err)
		return
	}

	w.Write([]byte(`success`))

	order := fmt.Sprintf("%v", res.Order)
	price := fmt.Sprintf("%v", res.Price)

	_, err = strconv.ParseFloat(price, 64)
	if err != nil {
		glog.Info("===>IPS支付price参数解析失败:", r.Form)
		return
	}

	payCode := fmt.Sprintf("%v", res.PayCode)
	state := fmt.Sprintf("%v", res.State)
	gameOrder := fmt.Sprintf("%v", res.GameOrder)
	if len(gameOrder) == 0 {
		glog.Info("===>IPS充值gameOrder无效:", gameOrder)
		return
	}
	sign := fmt.Sprintf("%v", res.Sign)

	qfUserId := fmt.Sprintf("%v", res.UserId)

	md5String := qfUserId + payCode + order + qfAppKey
	calcSign := calcMd5(md5String)
	if calcSign != sign {
		glog.Info("====>签名检验失败sign:", sign, " calcSign:", calcSign)
		return
	}

	if state != "success" {
		glog.Info("===>失败交易:", r.Form)
		return
	}

	if payCode == PRODUCT_VIP1 || payCode == PRODUCT_VIP2 || payCode == PRODUCT_VIP3 || payCode == PRODUCT_VIP4 || payCode == IOS_VIP_25_VIP1 || payCode == IOS_VIP_898_VIP4 {
		QfBuyVIP(res, 183)
	} else if payCode == QUICK_PAY10 || payCode == QUICK_PAY20 || payCode == QUICK_PAY50 || payCode == QUICK_PAY500 || payCode == FIRST_PAY30 || payCode == IOS_QUICK_12 || payCode == IOS_QUICK_18 || payCode == IOS_QUICK_518 {
		QfQuickPay(res)
	} else {
		QfBuyDiamond(res)
	}
}
