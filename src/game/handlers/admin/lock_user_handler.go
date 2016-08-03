package admin

import (
	"code.google.com/p/goprotobuf/proto"
	"config"
	"game/domain/user"
	"github.com/golang/glog"
	"net/http"
	"pb"
)

func LockUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("key") != config.ControlKey {
		glog.Error("key不符Addr:", r.RemoteAddr)
		return
	}

	userId := r.FormValue("userId")
	if userId == "" {
		w.Write([]byte(`0`))
		return
	}

	isLocked := r.FormValue("isLocked")

	_, err := user.FindByUserId(userId)
	if err != nil {
		glog.Error(err)
		w.Write([]byte(`0`))
		return
	}

	if user.GetPlayerManager().IsOnline(userId) {
		user.GetPlayerManager().Kickout(userId)
		glog.Info("===>管理员锁定玩家userId:", userId, " addr:", r.RemoteAddr)
		lockMsg := &pb.MsgLockUser{}
		if isLocked == "true" {
			lockMsg.Lock = proto.Bool(true)
		} else {
			lockMsg.Lock = proto.Bool(false)
		}
		user.GetPlayerManager().SendServerMsg("", []string{userId}, int32(pb.ServerMsgId_MQ_LOCK_USER), lockMsg)
	}

	if isLocked == "true" {
		user.SetLocked(userId, true)
	} else {
		user.SetLocked(userId, false)
	}

	w.Write([]byte(`1`))
}
