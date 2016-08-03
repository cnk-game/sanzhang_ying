package pay

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	"net/http"
)

const (
	KupaiPublicKey = "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCbAe6/sQHhLIBMQJXKHm12noqAMb+LXzW2rBVIfLJZj0tVRAA5lEeACdB5AJpCLE3klnO+m7gxAuZn6PRdON7MM4oF4qxtJdhKmjAyscjGm+cQEz0VnrJJbwXcANorogX6hwjpp+JwuOqFGtcbkUUGTo5b+FcuFGd6ZGhdyA/K2wIDAQAB"
)

type KupaiRes struct {
	Transtype int     `json:"transtype"`
	Cporderid string  `json:"cporderid"`
	Transid   string  `json:"transid"`
	Appuserid string  `json:"appuserid"`
	Appid     string  `json:"appid"`
	Waresid   int     `json:"waresid"`
	Feetype   int     `json:"feetype"`
	Money     float32 `json:"money"`
	Currency  string  `json:"currency"`
	Result    int     `json:"result"`
	Transtime string  `json:"transtime"`
	Cpprivate string  `json:"cpprivate"`
	Paytype   int     `json:"paytype"`
}

func kupaidoCheck(content string, sign string) bool {

	glog.Info("KupaiPublicKey=", KupaiPublicKey)

	//base64解码

	publickey, _ := base64.URLEncoding.DecodeString(KupaiPublicKey)

	glog.Info("base64解码后publickey=", publickey)

	//	block, _ := pem.Decode([]byte(KupaiPublicKey))
	//	if block == nil { // 失败情况
	//		glog.Info("pem fiail")
	//		return false
	//	}

	//pub, err := x509.ParsePKIXPublicKey(block.Bytes)

	pub, err := x509.ParsePKIXPublicKey(publickey)

	if err != nil {
		glog.Info("x509 err=", err)
		return false
	}

	baseSing, _ := base64.URLEncoding.DecodeString(sign)

	erren := rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.MD5, []byte(content), baseSing)

	if erren != nil {
		glog.Info("kupai verify erren=", erren)
		return false
	} else {
		glog.Info("kupai docheck success")
		return true
	}

}

func KupaiPayHandler(w http.ResponseWriter, r *http.Request) {

	glog.Info("KupaiPayHandler:")

	r.ParseForm()

	glog.Info("kupai:", r.Form)

	//post方式
	transdata := r.PostFormValue("transdata")
	sign := r.PostFormValue("sign")
	//signtype := r.PostFormValue("signtype")

	glog.Info("transdata = ", transdata)
	glog.Info("sign = ", sign)

	//result := kupaidoCheck(transdata, sign)

	//	if result == false {
	//		glog.Info("===>酷派充值签名验证失败", " sign:", sign, " addr:", r.RemoteAddr)
	//		w.Write([]byte(`FAILURE`))
	//		return
	//	}

	w.Write([]byte(`SUCCESS`))

	res := &KupaiRes{}
	err := json.Unmarshal([]byte(transdata), res)
	if err != nil {
		glog.Error("酷派解析数据失败err:", err)
		glog.Error("酷派解析数据失败数据:", transdata)
		return
	}

	//交易失败
	if res.Result != 0 {
		glog.Error("酷派服务器返回交易失败")
		return
	}

	if res.Appid != "5000002955" && res.Money < 0.01 {
		glog.Error("酷派服务器返回appid错误", res.Appid)
		glog.Error("酷派服务器返回Money错误", res.Money)
		return
	}

	payType := fmt.Sprintf("%v", res.Paytype)
	money := fmt.Sprintf("%v", res.Money)
	commonPay(res.Cporderid, "kupai", money, payType, "success", res.Transid)
}
