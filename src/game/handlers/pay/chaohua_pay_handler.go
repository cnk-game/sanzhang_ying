package pay

import (
	"crypto/md5"
	"fmt"
	"github.com/golang/glog"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	CaoHuaKey = "B74D8B70788276F7FAEBB89A232CADAF"
)

func ChaohuaPayHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	glog.Info("chaohua:", r.Form)

	orderNo := r.FormValue("OrderNo")
	outPayNo := r.FormValue("OutPayNo")
	userID := r.FormValue("UserID")
	serverNo := r.FormValue("ServerNo")
	payType := r.FormValue("PayType")
	money := r.FormValue("Money")
	pMoney := r.FormValue("PMoney")
	payTime := r.FormValue("PayTime")
	sign := r.FormValue("Sign")

	signStr := ""
	signStr += fmt.Sprintf("%v", orderNo)
	signStr += fmt.Sprintf("%v", outPayNo)
	signStr += fmt.Sprintf("%v", userID)
	signStr += fmt.Sprintf("%v", serverNo)
	signStr += fmt.Sprintf("%v", payType)
	signStr += fmt.Sprintf("%v", money)
	signStr += fmt.Sprintf("%v", pMoney)
	signStr += fmt.Sprintf("%v", payTime)
	signStr += fmt.Sprintf("%v", CaoHuaKey)

	h := md5.New()
	io.WriteString(h, signStr)
	localSign := fmt.Sprintf("%x", h.Sum(nil))
	localSign = strings.ToUpper(localSign)

	w.Write([]byte(`1`))

	if localSign != sign {
		glog.Info("===>草花充值签名验证失败localSign:", localSign, " sign:", sign, " addr:", r.RemoteAddr)
		return
	}

	payType = fmt.Sprintf("%v", payType)
	money = fmt.Sprintf("%v", money)
	glog.Info("===>草花支付类型:", payType)
	glog.Info("===>草花充值金额1:", money)

	priceReal, err := strconv.ParseFloat(money, 64)
	glog.Info("===>草花充值金额priceReal:", priceReal)
	if err != nil {
		fmt.Println(err)
		return
	}
	if priceReal == 0 {
		return
	}

	commonPay(outPayNo, "caohua", money, payType, "success", orderNo)
}
