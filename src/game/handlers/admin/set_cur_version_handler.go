package admin

import (
	"config"
	"github.com/golang/glog"
	"net/http"
)

func SetCurVersionHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.FormValue("key") != config.ControlKey {
		glog.Error("key不符Addr:", r.RemoteAddr)
		return
	}

	version := r.FormValue("version")
	if version == "" {
		w.Write([]byte(`0`))
		return
	}

	config.GetConfigManager().SetCurVersion(version)

	glog.Info("===>设置当前版本号:", version)

	w.Write([]byte(`1`))
}
