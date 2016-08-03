package admin

import (
	"config"
	"encoding/json"
	domainUser "game/domain/user"
	"github.com/golang/glog"
	"net/http"
)

type UserInfoResp struct {
	Ok    bool       `json:"ok"`
	Infos []UserInfo `json:"infos"`
}

type UserInfo struct {
	UserId     string  `json:"userId"`
	Username   string  `json:"username"`
	Nickname   string  `json:"nickname"`
	Level      int     `json:"level"`
	VipLevel   int     `json:"vipLevel"`
	Gold       int64   `json:"gold"`
	Charm      int     `json:"charm"`
	Diamond    int     `json:"diamond"`
	CreateTime int64   `json:"createTime"`
	IsLocked   bool    `bson:"isLocked"`
}

func QueryUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.FormValue("key") != config.ControlKey {
		glog.Error("key不符Addr:", r.RemoteAddr)
		return
	}

	resp := &UserInfoResp{}

	userId := r.FormValue("userId")
	u, err := domainUser.FindByUserId(userId)
	if err != nil {
		glog.Error(err)
		w.Write([]byte(`{ok:false}`))
		return
	}

	user := UserInfo{}
	user.UserId = u.UserId
	user.Username = u.UserName
	user.Nickname = u.Nickname
	user.IsLocked = u.IsLocked
	user.CreateTime = u.CreateTime.Unix()

	f, ok := domainUser.GetUserFortuneManager().GetUserFortune(userId)
	if ok {
		user.Level = f.Exp
		user.VipLevel = f.VipLevel
		user.Diamond = f.Diamond
		user.Gold = f.Gold
		user.Charm = f.Charm
	} else {
        f, err := domainUser.FindUserFortune(userId)
        if err != nil {
            glog.Error(err)
            w.Write([]byte(`{ok:false}`))
            return
        }
        user.Level = f.Exp
        user.VipLevel = f.VipLevel
        user.Diamond = f.Diamond
        user.Gold = f.Gold
        user.Charm = f.Charm
	}

	resp.Ok = true
	resp.Infos = append(resp.Infos, user)

	b, err := json.Marshal(resp)
	if err != nil {
		w.Write([]byte(`{ok:false}`))
		return
	}

	w.Write(b)

	return
}
