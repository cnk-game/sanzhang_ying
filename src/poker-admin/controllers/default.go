package controllers

import (
	"github.com/astaxie/beego"
    "fmt"
    "poker-admin/models"
)



////////////////////////////////////////////////////////////////////////////////
// 登陆
type LoginController struct {
    beego.Controller
}
func (c *LoginController) Get() {
    c.TplNames = "login.html"
}

func (c *LoginController) Post() {
    username := c.GetString("username")
    userpwd := c.GetString("userpwd")
    if len(username) != 0 || len(userpwd) != 0 {
        result, userinfo := models.AdminMgr.Login(username, userpwd)
	
        if 0 == result {
            c.SetSession(models.SessionKey_UserId, userinfo.UserId)
            c.SetSession(models.SessionKey_UserName, userinfo.UserName)
            c.SetSession(models.SessionKey_UserType, userinfo.UserType)
            c.SetSession(models.SessionKey_Channel, userinfo.Channel)
            fmt.Println("=========> user:", username, " login succeed")
        }
        c.Data["json"] = result
    } else {
        c.Data["json"] = 3
    }
    c.ServeJson()
}

////////////////////////////////////////////////////////////////////////////////
// 登出
type LogoutController struct {
    beego.Controller
}
func (c *LogoutController) Get() {
    username := c.GetSession(models.SessionKey_UserName).(string)
    if len(username) > 0 {
        models.AdminMgr.Logout(username)
    }
    c.DestroySession()
    c.Redirect("/login", 302)
}






type NotifyMgrController struct {
    beego.Controller
}
func (c *NotifyMgrController) Get() {
    c.TplNames = "notify_mgr.html"
}

type PlayerMgrController struct {
    beego.Controller
}
func (c *PlayerMgrController) Get() {
    c.TplNames = "player_mgr.html"
}


/**mim
 * 修改密碼 
 * PwdMgrController
 */
type PwdMgrController struct {
    beego.Controller
}
func (c *PwdMgrController) Get() {
    c.TplNames = "pwd_mgr.html"
}

/**
 * 保存密碼 
 * PwdMgrController
 */
type SavePwdController struct {
    beego.Controller
}
func (c *SavePwdController) Post() {
    password := c.GetString("password")
    if len(password) >= 6  {
        userid := c.GetSession(models.SessionKey_UserId).(string)

        fmt.Println("=========> UserId:", userid, " =========> password: ", password)
        models.SetUserPwd(userid, password)
        c.Data["json"] = 1
    } else {
        c.Data["json"] = 0
    }
    c.ServeJson()
}



