package littlegame

import (
	"game/domain/littlegame"
	"github.com/golang/glog"
	"net/http"
)

func ChangeProbablyHandler(w http.ResponseWriter, r *http.Request) {
	glog.Info("ChangeProbablyHandler in.")

	result := littlegame.GetCardConfigManager().Init()
	glog.Infof("%v", result)

	if result == true {
		w.Write([]byte(`ChangeProbablyHandler success`))
	} else {
		w.Write([]byte(`ChangeProbablyHandler fail`))
	}

}
