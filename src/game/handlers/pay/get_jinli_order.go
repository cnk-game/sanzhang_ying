package pay

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"

	"crypto/x509"
	"encoding/base64"
	"encoding/json"

	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	//PRODUCT_VIP1                 = "100559"
	//PRODUCT_VIP2                 = "100560"
	//PRODUCT_VIP3                 = "100561"
	//PRODUCT_VIP4                 = "100562"
	//QUICK_PAY10                  = "100563"
	//QUICK_PAY20 = "100577"
	//QUICK_PAY50                  = "100564"
	//QUICK_PAY500                 = "100565"
	//FIRST_PAY30                  = "100576"
	//ACTIVE_10                    = "100612"
	//ACTIVE_30                    = "100613"
	DIAMON_10                    = "100519"
	DIAMON_50                    = "100520"
	DIAMON_100                   = "100521"
	DIAMON_300                   = "100522"
	DIAMON_500                   = "100523"
	DIAMON_1000                  = "100524"
	jinliAPIkey                  = "918AD85BD7214741B04B0E2ABA1AD3AD"
	deliver_type                 = "1"
	GIONEE_PAY_INIT              = "https://pay.gionee.com/order/create"
	CREATE_SUCCESS_RESPONSE_CODE = "200010000"
	PrivateKey                   = "MIICdwIBADANBgkqhkiG9w0BAQEFAASCAmEwggJdAgEAAoGBAI1zpRXvhN403LHjQy2AYL4ocgIREo0DN70TFwo4qBYWa1B7mOiUnDTYEiHBHC7GSV68IkTTFLHRORInKfgD/iggVcttLhjO7HwdBc6OQk2C2AzluKdZo5zR7VKxN4KZ60ZME7IdKUUSC5OxaTeT0JZE/aNmQtG+0OtYykT/lR/RAgMBAAECgYB2TSDlnqO2H+nwEM0PUg4XG8Z1+gdUzjsgs8WSM95IUsm/zw8MfbXC8G8Bsfs7F3UneRsECrGFIiGkEAMmvVgiwqk8mGS8oIgDtdZdbEsRf2rPDGx1s1MbrTvs95wPacVSXQkD0Ur4n9luRyAHw4EEzwgHbSVU4cdYrVSEV5G8TQJBAMRz27xYEuwbAQGGu9hfCDwBhZnV1F6krtOMdS8/qtAQE8JnnCux9Fql0QY+p58Sc8DqcOtD88vHZHbufon3Pd8CQQC4U+I0LSEm1Up8uUumvR4ULLQXTVYkID/rQkIGlu2PFl5Z2STFW6NCuv5Xm51W5qVnidU1NWteq4Y8C/owx/hPAkEAqo2vam/IVcUH9YxEjw/KNVZY5/qVemldAnqBzjhnEmWy0edj1SeU7hHhS5ufqOG7LvQafpYrFXKRTRO3Ng45XwJAAM6lO/NCpOfkNp2dHjLP0ejMNRnqmafmf8I/hcXdbnX7ncscpRycn2swN/P/gWTrLoPlAiGkwbpgkRzAULxfcwJBALAF3HtSpYSaHYx+42t1cKxnSS9iR6QeI7z4bINXvaS1ws4u30o4JJG614/sM3dM/sEiynWqGjlQw61GkIW3Cz8="
)

var JinliOrderMu sync.RWMutex

type JinliOrderRes struct {
	Status      string `json:"status"`
	Apikey      string `json:"api_key"`
	Description string `json:"description"`
	OutOrderNo  string `json:"out_order_no"`
	SubmitTime  string `json:"submit_time"`
	OrderNo     string `json:"order_no"`
}

type JinliOrderReq struct {
	ApiKey      string `json:"api_key"`
	DealPrice   string `json:"deal_price"`
	DeliverType string `json:"deliver_type"`
	OutOrderNo  string `json:"out_order_no"`
	Subject     string `json:"subject"`
	SubmitTime  string `json:"submit_time"`
	TotalFee    string `json:"total_fee"`
	Sign        string `json:"sign"`
	PlayerId    string `json:"player_id"`
}

func getCurrentTime() string {
	now := time.Now()

	year, mon, day := now.Date()
	hour, min, sec := now.Clock()
	//zone, _ := now.Zone()
	return fmt.Sprintf("%d%02d%02d%02d%02d%02d", year, mon, day, hour, min, sec)

}

func sign(content string, privateKey string) (signdata string, result bool) {

	//base64解密
	key, _ := base64.StdEncoding.DecodeString(privateKey)

	glog.Info("key=", key)

	//	block, _ := pem.Decode([]byte(key))
	//	glog.Info("block=", block)

	//	if block == nil { // 失败情况
	//		return "", false
	//	}

	//private, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	private, err := x509.ParsePKCS8PrivateKey(key)
	glog.Info("private=", private)

	if err != nil {
		return "", false
	}

	h := crypto.Hash.New(crypto.SHA1)

	h.Write([]byte(content))

	hashed := h.Sum(nil)

	glog.Info("hashed=", hashed)

	// 进行rsa加密签名
	signedData, err := rsa.SignPKCS1v15(rand.Reader, private.(*rsa.PrivateKey), crypto.SHA1, hashed)
	if err != nil {
		return "", false
	}
	glog.Info("signedData=", signedData)

	nonce := make([]byte, 200)
	base64.StdEncoding.Encode(nonce, signedData)

	glog.Info("nonce=", nonce)

	if nonce != nil {
		return string(nonce), true
	} else {
		return "", false
	}

}

func GetJinliOrder(userId string, productid string, oerderid string) (subTime string) {
	glog.Info("GetJinliOrder in")
	glog.Info("userId=", userId)
	glog.Info("productid=", productid)
	glog.Info("oerderid=", oerderid)

	JinliOrderMu.Lock()
	defer JinliOrderMu.Unlock()

	productPrice := ""
	productName := ""

	if productid == PRODUCT_VIP1 {
		productPrice = "28"
		productName = "特权1"
	} else if productid == PRODUCT_VIP2 {
		productPrice = "188"
		productName = "特权2"
	} else if productid == PRODUCT_VIP3 {
		productPrice = "588"
		productName = "特权3"
	} else if productid == PRODUCT_VIP4 {
		productPrice = "888"
		productName = "特权4"
	} else if productid == QUICK_PAY10 {
		productPrice = "10"
		productName = "快充10万"
	} else if productid == QUICK_PAY10 {
		productPrice = "10"
		productName = "快充10万"
	} else if productid == QUICK_PAY20 {
		productPrice = "20"
		productName = "快充20万"
	} else if productid == QUICK_PAY50 {
		productPrice = "50"
		productName = "快充50万"
	} else if productid == QUICK_PAY500 {
		productPrice = "500"
		productName = "快充500万"
	} else if productid == FIRST_PAY30 {
		productPrice = "30"
		productName = "首充礼包"
	} else if productid == ACTIVE_10 {
		productPrice = "10"
		productName = "特权限时特惠"
	} else if productid == ACTIVE_30 {
		productPrice = "30"
		productName = "双旦充值送双倍"
	} else if productid == DIAMON_10 {
		productPrice = "10"
		productName = "10钻石"
	} else if productid == DIAMON_50 {
		productPrice = "50"
		productName = "50钻石"
	} else if productid == DIAMON_100 {
		productPrice = "100"
		productName = "100钻石"
	} else if productid == DIAMON_300 {
		productPrice = "300"
		productName = "300钻石"
	} else if productid == DIAMON_500 {
		productPrice = "500"
		productName = "500钻石"
	} else if productid == DIAMON_1000 {
		productPrice = "1000"
		productName = "1000钻石"
	}

	glog.Info("productPrice=", productPrice)
	glog.Info("productName=", productName)

	submitTime := getCurrentTime()

	glog.Info("submitTime=", submitTime)

	signContent := jinliAPIkey + productPrice + deliver_type + oerderid + productName + submitTime + productPrice

	glog.Info("signContent=", signContent)

	sindata, result := sign(signContent, PrivateKey)
	if result == false {
		return ""
	}

	reqstr := JinliOrderReq{}
	reqstr.ApiKey = jinliAPIkey
	reqstr.DealPrice = productPrice
	reqstr.DeliverType = deliver_type
	reqstr.OutOrderNo = oerderid
	reqstr.Subject = productName
	reqstr.SubmitTime = submitTime
	reqstr.TotalFee = productPrice
	reqstr.Sign = sindata
	reqstr.PlayerId = userId

	glog.Info("reqstr=", reqstr)

	sendStr, _ := json.Marshal(reqstr)

	glog.Info("sendStr=", string(sendStr))

	//	pool := x509.NewCertPool()
	//	tr := &http.Transport{
	//		TLSClientConfig: &tls.Config{
	//			RootCAs:      pool,
	//			Certificates: []tls.Certificate{cliCrt},
	//		},
	//	}

	//client := &http.Client{Transport: tr}
	client := &http.Client{}
	req, err := http.NewRequest("POST", GIONEE_PAY_INIT, strings.NewReader(string(sendStr)))
	if err != nil {
		fmt.Println("Get error:", err)
		return ""
	}

	//	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//	req.Header.Set("Authorization", builderAuthorization())

	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))

	resultSrt := fmt.Sprintf("%v", body)
	glog.Info("JInLiorder result:", resultSrt)

	res := &JinliOrderRes{}
	err = json.Unmarshal(body, res)
	if err != nil {
		glog.Error("解析数据失败err:", err)
		return ""
	}

	glog.Info("SubmitTime = ", res.SubmitTime)

	if res.Status == CREATE_SUCCESS_RESPONSE_CODE {

		return res.SubmitTime
	} else {
		glog.Info("JInLiorder result=", res.Status)
		return ""
	}
}
