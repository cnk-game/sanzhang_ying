package pay

import (
	"crypto/md5"

	"fmt"
	"github.com/golang/glog"
	"io"
	"net/http"
	"strconv"
)

const (
	appkey = "60bdf4f43fdebfa1e28d1d196220b810"
	PayKey = "f5384c405662068e62eb9c7eb6db5223"
)

func ios51doCheck(content string, sign string) bool {

	h := md5.New()
	io.WriteString(h, content)
	localSign := fmt.Sprintf("%x", h.Sum(nil))
	//localSign = strings.ToUpper(localSign)

	if localSign != sign {
		return false
	} else {
		glog.Info("ios51 docheck success")
		return true
	}

}

func Ios51PayHandler(w http.ResponseWriter, r *http.Request) {

	glog.Info("Ios51PayHandler:")

	r.ParseForm()

	glog.Info("Ios51:", r.Form)

	//post方式
	order_no := r.PostFormValue("order_no")
	uid := r.PostFormValue("uid")
	amount := r.PostFormValue("amount")
	serverid := r.PostFormValue("serverid")
	extra := r.PostFormValue("extra")
	notify_time := r.PostFormValue("notify_time")
	sign := r.PostFormValue("sign")
	paysign := r.PostFormValue("paysign")

	glog.Info("sign=", sign)
	glog.Info("paysign=", paysign)

	resultdata := appkey
	resultdata += "amount=" + amount
	resultdata += "&extra=" + extra
	resultdata += "&order_no" + order_no
	resultdata += "&serverid=" + serverid
	resultdata += "&notify_time=" + notify_time
	resultdata += "&uid=" + uid

	glog.Info("resultdata=", resultdata)
	result := ios51doCheck(resultdata, sign)

	w.Write([]byte(`success`))

	if result == false {
		glog.Info("===>51ios充值签名验证失败", " sign:", sign, " addr:", r.RemoteAddr)
		glog.Info("发货失败")
		return
	}

	priceReal, err := strconv.ParseFloat(amount, 64)
	glog.Info("===>51ios充值金额priceReal:", priceReal)
	if err != nil {
		fmt.Println(err)
		return
	}
	if priceReal == 0 {
		return
	}

	money := fmt.Sprintf("%v", amount)
	commonPay(extra, "51ios", money, "", "success", order_no)
}
