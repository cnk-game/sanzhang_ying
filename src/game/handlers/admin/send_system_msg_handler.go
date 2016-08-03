package admin

import (
	"config"
	"game/domain/user"
	"github.com/golang/glog"
	"net/http"
	"pb"
	"util"
)

func SendSystemMsgHandler(w http.ResponseWriter, r *http.Request) {
    glog.Info("SendSystemMsgHandler")
	if r.FormValue("key") != config.ControlKey {
		glog.Error("key不符Addr:", r.RemoteAddr)
		return
	}

	msg := r.FormValue("msg")
	if msg == "" {
		w.Write([]byte(`0`))
		return
	}

	user.GetPlayerManager().BroadcastClientMsg(int32(pb.MessageId_CHAT), util.BuildSysBugle(msg))

	w.Write([]byte(`1`))
}
