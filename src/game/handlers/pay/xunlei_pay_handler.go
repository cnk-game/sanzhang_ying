package pay

import (
	"crypto/md5"
	"fmt"
	"github.com/golang/glog"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

const (
	XunleiKey = "XSaSS5sGLyJmNBlqPKcDDHVs"
)

func XunLeshiPayHandler(w http.ResponseWriter, r *http.Request) {
	glog.Info("XunLeshiPayHandler in")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		glog.Error(err)
		return
	}

	l, _ := url.ParseQuery(string(b))

	uId := l["user"][0]

	orderId := l["orderid"][0]
	gold := l["gold"][0]
	money := l["money"][0]
	time := l["time"][0]
	ext := l["ext"][0]
	sign := l["sign"][0]

	signStr := ""
	signStr += orderId
	signStr += uId
	signStr += gold
	signStr += money
	signStr += time
	signStr += XunleiKey

	glog.Info("XunLeshiPayHandler data:", signStr)

	h := md5.New()
	io.WriteString(h, signStr)
	localSign := fmt.Sprintf("%x", h.Sum(nil))

	if localSign == sign {
		w.Write([]byte(`1`))
	} else {
		glog.Info("===>迅雷充值签名验证失败localSign:", localSign, " sign:", sign, " addr:", r.RemoteAddr)
		w.Write([]byte(`-2`))
		return
	}

	priceReal, err := strconv.ParseFloat(money, 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	if priceReal == 0 {
		return
	}

	commonPay(ext, "xunlei", money, "", "success", orderId)
}
