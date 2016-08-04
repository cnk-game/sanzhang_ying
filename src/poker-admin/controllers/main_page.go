package controllers

import (
    "github.com/astaxie/beego"
    "poker-admin/models/http"
    "poker-admin/models"
    "fmt"
)




///////////////////////////////////////////////////////////////
type MainController struct {
    beego.Controller
}
func (c *MainController) Get() {
    c.TplNames = "main.html"
}

///////////////////////////////////////////////////////////////
type MainTodayQueryController struct {
    beego.Controller
}
func (c *MainTodayQueryController) Get() {
    channel := c.GetSession(models.SessionKey_Channel).(string)
    c.Data["json"] = models.StatisticsMgr.GetMainGameInfo(channel)
    c.ServeJson()
}

///////////////////////////////////////////////////////////////
type MainNowQueryController struct {
    beego.Controller
}
func (c *MainNowQueryController) Get() {
    c.Data["json"] = http.GetMainNowQuery()
    c.ServeJson()
}

///////////////////////////////////////////////////////////////
type OnlineStatus struct {
    Period string
    OnlinePlayerCount int
}
type TodayOnlineStatusController struct {
    beego.Controller
}
func (c *TodayOnlineStatusController) Get() {
    channel := c.GetSession(models.SessionKey_Channel).(string)
    if !models.IsAdminChannel(channel) {
        c.Data["json"] = nil
        c.ServeJson()
    } else {
        list := []*OnlineStatus{}
        logs := models.StatisticsMgr.GetOnlineInfo()
        for _, log := range logs {
            item := &OnlineStatus{}
            item.Period = fmt.Sprintf("%.4d-%.2d-%.2d %.2d:%.2d", log.DateTime.Year(), int(log.DateTime.Month()), log.DateTime.Day(), log.DateTime.Hour(), log.DateTime.Minute())
            item.OnlinePlayerCount = log.OnlinePlayerCount
            list = append(list, item)
        }
        c.Data["json"] = list
        c.ServeJson()
    }
}