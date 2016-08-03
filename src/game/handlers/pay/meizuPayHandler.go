package pay

import (
	"crypto/md5"
	"strconv"

	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"io"
	"net/http"
)

const (
	appSecret = "xIF0GIyisU3E1j3GeyKDNI864friYVQ6"
)

type meizuRes struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	Value    string `json:"value"`
	Redirect string `json:"redirect"`
}

func meizudoCheck(content string, sign string) bool {

	h := md5.New()
	io.WriteString(h, content)
	localSign := fmt.Sprintf("%x", h.Sum(nil))
	//localSign = strings.ToUpper(localSign)

	if localSign != sign {
		return false
	} else {
		glog.Info("meizu docheck success")
		return true
	}

}

func MeizuPayHandler(w http.ResponseWriter, r *http.Request) {

	glog.Info("meizuPayHandler:")

	r.ParseForm()

	glog.Info("meizu:", r.Form)

	//post方式
	notify_time := r.PostFormValue("notify_time")
	notify_id := r.PostFormValue("notify_id")
	order_id := r.PostFormValue("order_id")
	app_id := r.PostFormValue("app_id")
	uid := r.PostFormValue("uid")
	partner_id := r.PostFormValue("partner_id")
	cp_order_id := r.PostFormValue("cp_order_id")
	product_id := r.PostFormValue("product_id")
	total_price := r.PostFormValue("total_price")
	trade_status := r.PostFormValue("trade_status")
	create_time := r.PostFormValue("create_time")
	pay_time := r.PostFormValue("pay_time")
	sign := r.PostFormValue("sign")
	//sign_type := r.PostFormValue("sign_type")
	buy_amount := r.PostFormValue("buy_amount")
	pay_type := r.PostFormValue("pay_type")
	product_per_price := r.PostFormValue("product_per_price")
	product_unit := r.PostFormValue("product_unit")
	user_info := r.PostFormValue("user_info")

	glog.Info("sign=", sign)

	resultdata := ""
	resultdata += "app_id=" + app_id
	resultdata += "&buy_amount=" + buy_amount
	resultdata += "&cp_order_id=" + cp_order_id
	resultdata += "&create_time=" + create_time
	resultdata += "&create_time=" + create_time
	resultdata += "&notify_id=" + notify_id
	resultdata += "&notify_time=" + notify_time
	resultdata += "&order_id=" + order_id
	resultdata += "&partner_id=" + partner_id
	resultdata += "&pay_time=" + pay_time
	resultdata += "&pay_type=" + pay_type
	resultdata += "&pay_type=" + pay_type
	resultdata += "&product_id" + product_id
	resultdata += "&product_per_price=" + product_per_price
	resultdata += "&product_unit=" + product_unit
	resultdata += "&total_price=" + total_price
	resultdata += "&trade_status=" + trade_status
	resultdata += "&uid=" + uid
	resultdata += "&user_info" + user_info
	resultdata += ":" + appSecret

	glog.Info("resultdata=", resultdata)
	result := meizudoCheck(resultdata, sign)

	res := &meizuRes{}

	if result == false {
		glog.Info("===>魅族充值签名验证失败", " sign:", sign, " addr:", r.RemoteAddr)
		res.Code = "120014"
		res.Message = ""
		resbyte, _ := json.Marshal(res)
		w.Write(resbyte)
		glog.Info("发货失败")
		return
	} else {
		res.Code = "200"
		res.Message = ""
		resbyte, _ := json.Marshal(res)
		w.Write(resbyte)
		glog.Info("成功发货")

	}

	//支付失败
	if trade_status != "3" {
		glog.Error("魅族服务器返回支付失败")
		return
	}

	priceReal, err := strconv.ParseFloat(total_price, 64)
	glog.Info("===>魅族充值金额priceReal:", priceReal)
	if err != nil {
		fmt.Println(err)
		return
	}
	if priceReal == 0 {
		return
	}

	payType := fmt.Sprintf("%v", pay_type)
	money := fmt.Sprintf("%v", total_price)
	commonPay(cp_order_id, "meizu", money, payType, "success", order_id)
}
