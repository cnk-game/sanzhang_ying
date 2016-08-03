package admin

import (
	domainPay "game/domain/pay"
	"net/http"
)

func ReloadPayActiveConfigHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userId := r.FormValue("userId")
	if userId != "admin" {
		w.Write([]byte(`userId error`))
	} else {
		domainPay.GetPayActiveConfig()
		w.Write([]byte(`success`))
	}
}
