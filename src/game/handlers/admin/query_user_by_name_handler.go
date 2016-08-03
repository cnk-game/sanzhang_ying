package admin

import (
	"config"
	"encoding/json"
	domainUser "game/domain/user"
	"github.com/golang/glog"
	"net/http"
)

func QueryUserByNameHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.FormValue("key") != config.ControlKey {
		glog.Error("key不符Addr:", r.RemoteAddr)
		return
	}

	nickname := r.FormValue("nickname")
	if nickname == "" {
		w.Write([]byte(`{ok:false}`))
		return
	}

	users, err := domainUser.FindByNickname(nickname)
	if err != nil {
		glog.Error(err)
		w.Write([]byte(`{ok:false}`))
		return
	}

	if len(users) <= 0 {
		w.Write([]byte(`{ok:false}`))
		return
	}

	resp := &UserInfoResp{}
	resp.Ok = true

	for _, u := range users {
		user := UserInfo{}
		user.UserId = u.UserId
		user.Username = u.UserName
		user.Nickname = u.Nickname
		user.IsLocked = u.IsLocked
		user.CreateTime = u.CreateTime.Unix()

		f, ok := domainUser.GetUserFortuneManager().GetUserFortune(u.UserId)
		if ok {
			user.Level = f.Exp
			user.VipLevel = f.VipLevel
			user.Diamond = f.Diamond
		}
		resp.Infos = append(resp.Infos, user)
	}

	b, err := json.Marshal(resp)
	if err != nil {
		w.Write([]byte(`{ok:false}`))
		return
	}

	w.Write(b)

	return
}
