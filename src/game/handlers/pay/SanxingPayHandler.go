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
	SanxingPublicKey = "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDSkFLf1d47EEcR9Tbu5O5scMapVnJLP46ww97xHeHUg4Hv8rgpTs1aYPMuIzDbWHKcRxmKCameTzjrnB8AFa7ToDsW/KmqlS2dTIMOf0MJkn98XBjuRaGi038oVZiCTOwcUVQ9HNWSBNwLEKUlyJu1OXaiymeo/EIX8vX7sllmrwIDAQAB"
)

type SanxingRes struct {
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

func SanxingdoCheck(content string, sign string) bool {

	glog.Info("SanxingPublicKey=", SanxingPublicKey)

	//base64解码

	publickey, _ := base64.URLEncoding.DecodeString(LianxiangPublicKey)

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
		glog.Info("Sanxing verify erren=", erren)
		return false
	} else {
		glog.Info("Sanxing docheck success")
		return true
	}

}

func SanxingPayHandler(w http.ResponseWriter, r *http.Request) {

	glog.Info("SanxingPayHandler:")

	r.ParseForm()

	glog.Info("Sanxing:", r.Form)

	//post方式
	transdata := r.PostFormValue("transdata")
	sign := r.PostFormValue("sign")
	//signtype := r.PostFormValue("signtype")

	glog.Info("transdata = ", transdata)
	glog.Info("sign = ", sign)

	//result := SanxingdoCheck(transdata, sign)

	//	if result == false {
	//		glog.Info("===>三星充值签名验证失败", " sign:", sign, " addr:", r.RemoteAddr)
	//		w.Write([]byte(`FAILURE`))
	//		return
	//	}

	w.Write([]byte(`SUCCESS`))

	res := &SanxingRes{}
	err := json.Unmarshal([]byte(transdata), res)
	if err != nil {
		glog.Error("三星解析数据失败err:", err)
		glog.Error("三星解析数据失败数据:", transdata)
		return
	}

	//交易失败
	if res.Result != 0 {
		glog.Error("三星服务器返回交易失败")
		return
	}

	if res.Appid != "3003773855" && res.Money < 0.01 {
		glog.Error("三星服务器返回appid错误", res.Appid)
		glog.Error("三星服务器返回Money错误", res.Money)
		return
	}

	payType := fmt.Sprintf("%v", res.Paytype)
	money := fmt.Sprintf("%v", res.Money)
	commonPay(res.Cporderid, "Sanxing", money, payType, "success", res.Transid)
}
