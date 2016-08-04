package util

import (
    "net/http"
    "io/ioutil"
    "bytes"
    "encoding/json"
	"github.com/golang/glog"
)

type GetVerifyReqMsg struct {
    MsgId    int       `json:"msgId"`
    Phone    string    `json:"phone"`
}

type GetVerifyResMsg struct {
    Result   int     `json:"result"`
    Phone    string  `json:"phone"`
}

func GetVerify(phone string) bool {
    req := GetVerifyReqMsg{}
    req.MsgId = 1
    req.Phone = phone
    b, err := json.Marshal(req)
    if err != nil {
        glog.Info("GetVerify err:", err)
        return false
    }

    body := bytes.NewBuffer([]byte(b))
    res, err := http.Post("http://sms.dapai2.com:4001", "application/json;charset=utf-8", body)
    if err != nil {
        glog.Info("GetVerify err:", err)
        return false
    }

    result, err := ioutil.ReadAll(res.Body)
    res.Body.Close()
    if err != nil {
        glog.Info("GetVerify err:", err)
        return false
    }

    resmsg := GetVerifyResMsg{}
    err = json.Unmarshal(result, &resmsg)
    if err != nil {
        glog.Info("GetVerify err:", err)
        return false
    }

    if !(resmsg.Result == 1 && resmsg.Phone == phone) {
        glog.Info("GetVerify err, Result=", resmsg.Result, "|Phone=", resmsg.Phone)
        return false
    }
    return true
}

type CheckVerifyReqMsg struct {
    MsgId    int     `json:"msgId"`
    Phone    string  `json:"phone"`
    Code     int     `json:"code"`
}

type CheckVerifyResMsg struct {
    Result   int     `json:"result"`
    Phone    string  `json:"phone"`
}

func CheckVerify(phone string, code int) bool {
    req := CheckVerifyReqMsg{}
    req.MsgId = 2
    req.Phone = phone
    req.Code = code
    b, err := json.Marshal(req)
    if err != nil {
        glog.Info("CheckVerify err:", err)
        return false
    }

    body := bytes.NewBuffer([]byte(b))
    res, err := http.Post("http://sms.dapai2.com:4001", "application/json;charset=utf-8", body)
    if err != nil {
        glog.Info("CheckVerify err:", err)
        return false
    }

    result, err := ioutil.ReadAll(res.Body)
    res.Body.Close()
    if err != nil {
        glog.Info("CheckVerify err:", err)
        return false
    }

    resmsg := CheckVerifyResMsg{}
    err = json.Unmarshal(result, &resmsg)
    if err != nil {
        glog.Info("CheckVerify err:", err)
        return false
    }

    if !(resmsg.Result == 1 && resmsg.Phone == phone) {
        glog.Info("CheckVerify err, Result=", resmsg.Result, "|Phone=", resmsg.Phone)
        return false
    }
    return true
}