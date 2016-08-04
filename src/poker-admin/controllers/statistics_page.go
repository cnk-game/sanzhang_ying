package controllers

import (
    "github.com/astaxie/beego"
    "fmt"
    "poker-admin/models"
)


///////////////////////////////////////////////////////////////
type SynthesizeInfoController struct {
    beego.Controller
}
func (c *SynthesizeInfoController) Get() {
    c.TplNames = "synthesize_info.html"
}


type SynthesizeQueryController struct {
    beego.Controller
}
func (c *SynthesizeQueryController) Get() {
    type SynthesizeLogInfo struct {
        LogList []*models.SynthesizeInfoLog
        TotalPage int
        CurPage int
        TotalCount int
    }
    channle := c.GetSession(models.SessionKey_Channel).(string)

    logInfo := SynthesizeLogInfo{}

    bYear, _ := c.GetInt("bYear")
    bMonth, _ := c.GetInt("bMonth")
    bDay, _ := c.GetInt("bDay")
    eYear, _ := c.GetInt("eYear")
    eMonth, _ := c.GetInt("eMonth")
    eDay, _ := c.GetInt("eDay")
    logInfo.CurPage, _ = c.GetInt("pageIdx")

    logInfo.LogList = []*models.SynthesizeInfoLog{}
    logInfo.TotalCount = models.LoadSynthesizeLogList(channle, bYear, bMonth, bDay, eYear, eMonth, eDay, logInfo.CurPage, &logInfo.LogList)
    fmt.Println("synthesize channel ========================>")

    logInfo.TotalPage = logInfo.TotalCount / models.PageItemCount
    if (logInfo.TotalCount % models.PageItemCount) > 0 {
        logInfo.TotalPage++
    }
    if logInfo.CurPage > logInfo.TotalPage {
        logInfo.CurPage = logInfo.TotalPage
    }
    c.Data["json"] = logInfo
    c.ServeJson()
}



///////////////////////////////////////////////////////////////
type NewPlayerInfoController struct {
    beego.Controller
}
func (c *NewPlayerInfoController) Get() {
    c.TplNames = "new_player_info.html"
}


type NewPlayerQueryController struct {
    beego.Controller
}
func (c *NewPlayerQueryController) Get() {
    type NewPlayerPage struct {
        UserList []*models.NewUserLog
        TotalPage int
        CurPage int
        TotalCount int
    }
    channle := c.GetSession(models.SessionKey_Channel).(string)
    pageInfo := NewPlayerPage{}
    bYear, _ := c.GetInt("bYear")
    bMonth, _ := c.GetInt("bMonth")
    bDay, _ := c.GetInt("bDay")
    eYear, _ := c.GetInt("eYear")
    eMonth, _ := c.GetInt("eMonth")
    eDay, _ := c.GetInt("eDay")
    pageInfo.CurPage, _ = c.GetInt("pageIdx")

    pageInfo.UserList = []*models.NewUserLog{}
    pageInfo.TotalCount = models.GetNewPlayers(bYear, bMonth, bDay, 0, 0, eYear, eMonth, eDay, 0, 0, pageInfo.CurPage, &pageInfo.UserList, channle)

    pageInfo.TotalPage = pageInfo.TotalCount / models.PageItemCount
    if (pageInfo.TotalCount % models.PageItemCount) > 0 {
        pageInfo.TotalPage++
    }
    if pageInfo.CurPage > pageInfo.TotalPage {
        pageInfo.CurPage = pageInfo.TotalPage
    }
    c.Data["json"] = pageInfo
    c.ServeJson()
}



///////////////////////////////////////////////////////////////
type PayInfoController struct {
    beego.Controller
}
func (c *PayInfoController) Get() {
    c.TplNames = "pay_info.html"
}


type PayQueryController struct {
    beego.Controller
}
func (c *PayQueryController) Get() {
    type PayPage struct {
        PayLogList []*models.PayLog
        TotalPage int
        CurPage int
        TotalCount int
    }
    channle := c.GetSession(models.SessionKey_Channel).(string)
    pageInfo := PayPage{}
    for ;; {
        orderId := c.GetString("orderId")
        if (len(orderId) > 0) {
            // 按订单查询
            fmt.Println("orderId : ", orderId)
            log := models.GetPayLogByOrderId(orderId, channle)
            c.Data["json"] = log
            break
        }
        userId := c.GetString("userId")
        if (len(userId) > 0) {
            // 按用户查询
            fmt.Println("userId : ", userId)
            pageInfo.CurPage, _ = c.GetInt("pageIdx")

            pageInfo.PayLogList = []*models.PayLog{}
            pageInfo.TotalCount = models.GetPayLogByUserId(userId, pageInfo.CurPage, &pageInfo.PayLogList, channle)
        } else {
            fmt.Println("query pay by time")
            bYear, _ := c.GetInt("bYear")
            bMonth, _ := c.GetInt("bMonth")
            bDay, _ := c.GetInt("bDay")
            eYear, _ := c.GetInt("eYear")
            eMonth, _ := c.GetInt("eMonth")
            eDay, _ := c.GetInt("eDay")
            pageInfo.CurPage, _ = c.GetInt("pageIdx")
            pageInfo.PayLogList = []*models.PayLog{}
            pageInfo.TotalCount = models.GetPayLogList(bYear, bMonth, bDay, 0, 0, 0, eYear, eMonth, eDay, 0, 0, 0, pageInfo.CurPage, &pageInfo.PayLogList, channle)
        }
        pageInfo.TotalPage = pageInfo.TotalCount / models.PageItemCount
        if (pageInfo.TotalCount % models.PageItemCount) > 0 {
            pageInfo.TotalPage++
        }
        if pageInfo.CurPage > pageInfo.TotalPage {
            pageInfo.CurPage = pageInfo.TotalPage
        }
        c.Data["json"] = pageInfo
        break
    }
    c.ServeJson()
}