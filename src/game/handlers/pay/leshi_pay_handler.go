package pay

import (
	"crypto/md5"
	"fmt"
	"github.com/golang/glog"
	"io"
	"net/http"
	"sort"
	"strconv"
)

const (
	LeshiKey = "86491a07bcff33a94dda22d3aad1aef7"
)

func LeshiPayHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	glog.Info("leshi:", r.Form)

	params := []string{}
	for key := range r.Form {
		if r.FormValue(key) == "null" {
			continue
		}

		if key != "sign" {
			params = append(params, key)
		}
	}

	sort.Strings(params)

	signStr := ""
	for _, key := range params {
		temp := r.FormValue(key)
		signStr += fmt.Sprintf("%v=%v&", key, temp)

	}

	signStr += fmt.Sprintf("key=%v", LeshiKey)

	h := md5.New()
	io.WriteString(h, signStr)
	localSign := fmt.Sprintf("%x", h.Sum(nil))

	w.Write([]byte(`success`))

	if localSign != r.FormValue("sign") {
		glog.Info("===>乐视充值签名验证失败localSign:", localSign, " sign:", r.FormValue("sign"), " addr:", r.RemoteAddr)
		return
	}

	marchantNo := r.FormValue("merchant_no")
	lePayOrderNo := r.FormValue("lepay_order_no")
	price := r.FormValue("price")

	priceReal, err := strconv.ParseFloat(price, 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	if priceReal == 0 {
		return
	}

	commonPay(marchantNo, "leshi", price, "", "success", lePayOrderNo)
}
