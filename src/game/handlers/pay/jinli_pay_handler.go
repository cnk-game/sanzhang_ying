package pay

import (
	"crypto"

	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"

	"fmt"
	"github.com/golang/glog"
	"net/http"
)

const (
	JinliPublicKey = "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCAqs7dUT6UoTruQGrKFiRBlj0S/yd7KB7Nt3OOCJ6uAbnJvMSV1MrKwUrmE/4CUZkXh9u3yyDoKRAZCajsDFI4sCOgR+1iG4kNVOa6JWuXZBF7zhd+MfEaEtGQ6fPalyU28JijEUwm7jL0HpbrXTYeDE7upFKLJv6PCaIBBCHgsQIDAQAB"
)

//var JinliPayMu sync.RWMutex

func doCheck(content string, sign string) bool {

	//base64解码
	//publickey, _ := base64.StdEncoding.DecodeString(JinliPublicKey)

	//glog.Info("publickey=", publickey)
	//	var publickey = []byte(`
	//	-----BEGIN PUBLIC KEY-----
	//	MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCAqs7dUT6UoTruQGrKFiRBlj0S
	//	/yd7KB7Nt3OOCJ6uAbnJvMSV1MrKwUrmE/4CUZkXh9u3yyDoKRAZCajsDFI4sCOg
	//	R+1iG4kNVOa6JWuXZBF7zhd+MfEaEtGQ6fPalyU28JijEUwm7jL0HpbrXTYeDE7u
	//	pFKLJv6PCaIBBCHgsQIDAQAB
	//	-----END PUBLIC KEY-----
	//	`)

	block, _ := pem.Decode([]byte(JinliPublicKey))
	if block == nil { // 失败情况
		glog.Info("pem fiail")
		return false
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)

	if err != nil {
		glog.Info("err=", err)
		return false
	}

	//base64解密
	signen, _ := base64.StdEncoding.DecodeString(sign)
	glog.Info("sign=", sign)

	erren := rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA1, []byte(content), signen)

	if erren != nil {
		glog.Info("erren=", erren)
		return false
	} else {
		glog.Info("yangzhengchenggong")
		return true
	}

}

func JinliPayHandler(w http.ResponseWriter, r *http.Request) {
	//	JinliPayMu.Lock()
	//	defer JinliPayMu.Unlock()
	glog.Info("JinliPayHandler:")

	r.ParseForm()

	glog.Info("jinli:", r.Form)

	//post方式
	apiKey := r.PostFormValue("api_key")
	closeTime := r.PostFormValue("close_time")
	createTime := r.PostFormValue("create_time")
	dealPrice := r.PostFormValue("deal_price")
	outOrderNo := r.PostFormValue("out_order_no")
	payChannel := r.PostFormValue("pay_channel")
	submitTime := r.PostFormValue("submit_time")
	//sign := r.FormValue("sign")

	signStr := ""
	signStr += fmt.Sprintf("api_key=%v", apiKey)
	signStr += fmt.Sprintf("&close_time=%v", closeTime)
	signStr += fmt.Sprintf("&create_time=%v", createTime)
	signStr += fmt.Sprintf("&deal_price=%v", dealPrice)
	signStr += fmt.Sprintf("&out_order_no=%v", outOrderNo)
	signStr += fmt.Sprintf("&pay_channel=%v", payChannel)
	signStr += fmt.Sprintf("&submit_time=%v", submitTime)
	signStr += fmt.Sprintf("&user_id=null")

	glog.Info("signStr = ", signStr)

	//	h := md5.New()
	//	io.WriteString(h, signStr)
	//	localSign := fmt.Sprintf("%x", h.Sum(nil))
	//	localSign = strings.ToUpper(localSign)

	//	result := doCheck(signStr, sign)

	//	if result == false {
	//		glog.Info("===>金立充值签名验证失败", " sign:", sign, " addr:", r.RemoteAddr)
	//		w.Write([]byte(`fail`))
	//		return
	//	}

	w.Write([]byte(`success`))

	payType := fmt.Sprintf("%v", payChannel)
	money := fmt.Sprintf("%v", dealPrice)
	commonPay(outOrderNo, "jinli", money, payType, "success", submitTime)
}
