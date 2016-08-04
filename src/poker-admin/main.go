package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"poker-admin/models"
	_ "poker-admin/routers"
	_ "poker-admin/util"
)

var FilterLoginCheck = func(ctx *context.Context) {
	fmt.Println("FilterLoginCheck Check=======================", ctx.Request.RequestURI)
	if  ctx.Request.RequestURI != "/login" && ctx.Request.RequestURI != "/" && ctx.Request.RequestURI != "/logout" {
		// 检查登陆状态
		userId, ok := ctx.Input.Session(models.SessionKey_UserId).(string)
		if !ok || 0 == len(userId) {
			ctx.Redirect(302, "/login")
		}
		fmt.Println("FilterLoginCheck Session => ", ok, " value => ", userId)
	}
}

func main() {

	models.CheckOldLogs()

	go models.DailyCalc()

	beego.SessionOn = true

	beego.InsertFilter("/*", beego.BeforeRouter, FilterLoginCheck)
	beego.InsertFilter("/", beego.BeforeRouter, FilterLoginCheck)
	beego.Run()
}
