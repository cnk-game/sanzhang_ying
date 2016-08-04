package util

import (
	"bytes"
	"encoding/json"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type GetCDKeyReqMsg struct {
	AppId       int    `json:"appId"`
	Code        string `json:"code"`
	UserChannel int    `json:"userChannel"`
	Imei        string `json:"imei"`
	UserId      int    `json:"userId"`
}

type GetCDKeyAckMsg struct {
	Result int    `json:"result"`
	Count  int    `json:"count"`
	Desc   string `json:"desc"`
}

func CheckCDKey(appId int, code string, userChannel int, imei string, userIdStr string) (int, int, int) {
	req := GetCDKeyReqMsg{}
	req.AppId = appId
	req.Code = code
	req.UserChannel = userChannel
	req.Imei = imei

	glog.Info("CheckCDKey appId:", appId)
	glog.Info("CheckCDKey code:", code)
	glog.Info("CheckCDKey userChannel:", userChannel)
	glog.Info("CheckCDKey imei:", imei)
	glog.Info("CheckCDKey userIdStr:", userIdStr)

	tem := strings.Trim(userIdStr, "QF1")
	userId, error := strconv.Atoi(tem)
	if error != nil {
		glog.Info("CheckCDKey err:", error)
		return -1, 0, 0
	}
	glog.Info("CheckCDKey userId:", userId)
	req.UserId = userId
	b, err := json.Marshal(req)
	if err != nil {
		glog.Info("CheckCDKey err:", err)
		return -1, 0, 0
	}

	body := bytes.NewBuffer([]byte(b))
	res, err := http.Post("http://103.26.1.220:7001/checkCode", "application/json;charset=utf-8", body)
	if err != nil {
		glog.Info("CheckCDKey err:", err)
		return -1, 0, 0
	}

	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		glog.Info("CheckCDKey err:", err)
		return -1, 0, 0
	}

	resmsg := GetCDKeyAckMsg{}
	err = json.Unmarshal(result, &resmsg)
	if err != nil {
		glog.Info("CheckCDKey err:", err)
		return -1, 0, 0
	}

	resultBack := resmsg.Result

	if !(resultBack == 1) {
		glog.Info("CheckCDKey err, Result = ", result)
		return resultBack, 0, 0
	}

	resultCount := resmsg.Count
	resultDesc := resmsg.Desc
	glog.Info("CheckCDKey count ", resultCount, " desc ", resultDesc)
	intDesc, err := strconv.Atoi(resultDesc)
	if err != nil {
		return -1, 0, 0
	}

	return 1, resultCount, intDesc
}
