package controllers

import (
    "github.com/astaxie/beego"
    "fmt"
    "poker-admin/models"
)


////////////////////////////////////////////////////////////////////////////////
type UserMgrController struct {
    beego.Controller
}
func (c *UserMgrController) Get() {
    c.TplNames = "user_mgr.html"
}



type AddUserResult struct {
    Result int
    Channel string
    UserName string
    Mark string
}
type AddUserController struct {
    beego.Controller
}
func (c *AddUserController) Post() {
    channel := c.GetString("channel")
    username := c.GetString("username")
    userpwd := c.GetString("userpwd")
    mark := c.GetString("mark")

    fmt.Println(channel, username, userpwd, mark)

    result := AddUserResult{}
    if 0 == len(channel) || 0 == len(username) || 0 == len(userpwd) || 0 == len(mark)  {
        result.Result = 1
    } else {
        usertype := c.GetSession(models.SessionKey_UserType)
        if usertype != models.UserType_Admin {
            result.Result = 2
        } else {
            userinfo := models.AdminMgr.GetInfo(username)
            if nil == userinfo {
                userinfo := models.AdminMgr.AddUser(username, userpwd, channel, mark)
                result.Result = 0
                result.Channel = userinfo.Channel
                result.UserName = userinfo.UserName
                result.Mark = userinfo.Mark
            } else {
                result.Result = 3
            }
        }
    }
    c.Data["json"] = result
    c.ServeJson()
}



type UserInfoDef struct {
    Channel string `bson:"Channel"`
    UserName string `bson:"UserName"`
    Mark string `bson:"Mark"`
}
type QueryUserListController struct {
    beego.Controller
}
func (c *QueryUserListController) Get() {
    usertype := c.GetSession(models.SessionKey_UserType)
    if usertype == models.UserType_Admin {
        userlist := []*UserInfoDef{}
        models.AdminMgr.GetAllUser(&userlist)
        c.Data["json"] = userlist
    } else {
        c.Data["json"] = nil
    }
    c.ServeJson()
}



type RemoveUserController struct {
    beego.Controller
}
func (c *RemoveUserController) Post() {
    username := c.GetString("username")
    if 0 == len(username) {
        c.Data["json"] = 1
    } else {
        session_username := c.GetSession(models.SessionKey_UserName)
        if session_username == username {
            c.Data["json"] = 2
        } else {
            usertype := c.GetSession(models.SessionKey_UserType)
            if usertype == models.UserType_Admin {
                models.AdminMgr.RemoveUser(username)
                c.Data["json"] = 0
            } else {
                c.Data["json"] = 3
            }
        }
    }
    c.ServeJson()
}



////////////////////////////////////////////////////////////////////////////////
type ConfigMgrController struct {
    beego.Controller
}
func (c *ConfigMgrController) Get() {
    c.TplNames = "config_mgr.html"
}

type ConfigInfo struct {
    VaildNewPlayerOnlineSecond int
    VaildNewPlayerMatchCount int
}
type QueryConfigController struct {
    beego.Controller
}
func (c *QueryConfigController) Get() {
    fmt.Println("==============QueryConfigController")
    usertype := c.GetSession(models.SessionKey_UserType)
    if usertype == models.UserType_Admin {
        info := &ConfigInfo{}
        info.VaildNewPlayerOnlineSecond = models.ConfigMgr.VaildNewPlayerOnlineSecond
        info.VaildNewPlayerMatchCount = models.ConfigMgr.VaildNewPlayerMatchCount
        c.Data["json"] = info
    } else {
        c.Data["json"] = nil
    }
    c.ServeJson()
}

type UpdateConfigController struct {
    beego.Controller
}
func (c *UpdateConfigController) Get() {
    usertype := c.GetSession(models.SessionKey_UserType)
    if usertype == models.UserType_Admin {
        vaild_new_player_online_second, _ := c.GetInt("vaild_new_player_online_second")
        vaild_new_player_match_count, _ := c.GetInt("vaild_new_player_match_count")
        models.ConfigMgr.VaildNewPlayerOnlineSecond = vaild_new_player_online_second
        models.ConfigMgr.VaildNewPlayerMatchCount = vaild_new_player_match_count
        models.UpdateChannelCheckConfig(vaild_new_player_online_second, vaild_new_player_match_count)
        c.Data["json"] = 0
    } else {
        c.Data["json"] = 1
    }
    c.ServeJson()
}


type QueryAdminController struct {
    beego.Controller
}
func (c *QueryAdminController) Get() {
    channel := c.GetSession(models.SessionKey_Channel).(string)
    if len(channel) > 0 {
        c.Data["json"] = 0
    } else {
        c.Data["json"] = 1
    }
    c.ServeJson()
}