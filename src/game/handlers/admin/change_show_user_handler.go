package admin

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	domainGame "game/domain/game"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
)

type ChangeUserRes struct {
	UserId interface{} `json:"userId"`
	Key    interface{} `json:"key"`
}

func ChangeShowUserHandler(w http.ResponseWriter, r *http.Request) {
	glog.Info("ChangeShowUserHandler in")
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		glog.Error(err)
		return
	}

	jsonData, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		glog.Error(err)
		return
	}

	res := &ChangeUserRes{}
	err = json.Unmarshal(jsonData, res)
	glog.Info("ChangeShowUserHandler:", string(jsonData))
	if err != nil {
		glog.Error("解析ChangeShowUserHandler数据失败err:", err)
		return
	}

	userId := fmt.Sprintf("%v", res.UserId)
	key := fmt.Sprintf("%v", res.Key)

	if key == "qifanthreekey123" {
		w.Write([]byte(`success`))
		glog.Info("old show user:", domainGame.ShowUser)
		domainGame.ShowUser = userId
	} else {
		w.Write([]byte(`fail`))
	}
}
