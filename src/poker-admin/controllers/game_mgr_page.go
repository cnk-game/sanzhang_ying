package controllers

import (
    "github.com/astaxie/beego"
    "fmt"
    "poker-admin/models/http"
    "poker-admin/models"
)

///////////////////////////////////////////////////////////////////////////////////////////////
type GetPrizeVersionController struct {
    beego.Controller
}
func (c *GetPrizeVersionController) Get() {
    usertype := c.GetSession(models.SessionKey_UserType)
    if usertype == models.UserType_Admin {
        version := models.GetPrizeVersion()
        fmt.Println("=======================> prize version : ", version)
        c.Data["json"] = version
    } else {
        c.Data["json"] = nil
    }
    c.ServeJson()
}
///////////////////////////////////////////////////////////////////////////////////////////////
type SavePrizeVersionController struct {
    beego.Controller
}
func (c *SavePrizeVersionController) Get() {
    usertype := c.GetSession(models.SessionKey_UserType)
    if usertype == models.UserType_Admin {
        version := c.GetString("version")
        result := http.SetPrizeVersion(version)
        if result == http.SucceedResult {
            models.SavePrizeVersion(version)
        }
        c.Data["json"] = result
    } else {
        c.Data["json"] = "failed"
    }
    c.ServeJson()
}

///////////////////////////////////////////////////////////////////////////////////////////////
type GetMatchConfigController struct {
    beego.Controller
}
func (c *GetMatchConfigController) Get() {
    usertype := c.GetSession(models.SessionKey_UserType)
    if usertype == models.UserType_Admin {
        config := models.GetMatchConfig()
        if (nil == config) {
            config = &models.MatchConfig{}
        }
        c.Data["json"] = config
    } else {
        c.Data["json"] = nil
    }
    c.ServeJson()
}

///////////////////////////////////////////////////////////////////////////////////////////////
type SaveMatchConfigController struct {
    beego.Controller
}
func (c *SaveMatchConfigController) Get() {
    usertype := c.GetSession(models.SessionKey_UserType)
    if usertype == models.UserType_Admin {
        GameType, _ := c.GetInt("GameType")
        Single, _ := c.GetInt("Single")
        Double, _ := c.GetInt("Double")
        ShunZi, _ := c.GetInt("ShunZi")
        JinHua, _ := c.GetInt("JinHua")
        ShunJin, _ := c.GetInt("ShunJin")
        BaoZi, _ := c.GetInt("BaoZi")
        WinGold, _ := c.GetInt("WinGold")
        WinRateHigh, _ := c.GetInt("WinRateHigh")
        LoseGold, _ := c.GetInt("LoseGold")
        WinRateLow, _ := c.GetInt("WinRateLow")

        config := models.GetMatchConfig()
        if (nil == config) {
            config = &models.MatchConfig{}
        }
        result := http.SetMatchConfig(GameType, Single, Double, ShunZi, JinHua, ShunJin, BaoZi, WinGold, WinRateHigh, LoseGold, WinRateLow)
        if result == http.SucceedResult {
            if GameType == 1 {
                config.CommonLevel1Single = Single
                config.CommonLevel1Double = Double
                config.CommonLevel1ShunZi = ShunZi
                config.CommonLevel1JinHua = JinHua
                config.CommonLevel1ShunJin = ShunJin
                config.CommonLevel1BaoZi = BaoZi
                config.CommonLevel1WinGold = WinGold
                config.CommonLevel1WinRateHigh = WinRateHigh
                config.CommonLevel1LoseGold = LoseGold
                config.CommonLevel1WinRateLow = WinRateLow
                models.SaveMatchConfig(config)
            }
            if GameType == 2 {
                config.CommonLevel2Single = Single
                config.CommonLevel2Double = Double
                config.CommonLevel2ShunZi = ShunZi
                config.CommonLevel2JinHua = JinHua
                config.CommonLevel2ShunJin = ShunJin
                config.CommonLevel2BaoZi = BaoZi
                config.CommonLevel2WinGold = WinGold
                config.CommonLevel2WinRateHigh = WinRateHigh
                config.CommonLevel2LoseGold = LoseGold
                config.CommonLevel2WinRateLow = WinRateLow
                models.SaveMatchConfig(config)
            }
            if GameType == 3 {
                config.CommonLevel3Single = Single
                config.CommonLevel3Double = Double
                config.CommonLevel3ShunZi = ShunZi
                config.CommonLevel3JinHua = JinHua
                config.CommonLevel3ShunJin = ShunJin
                config.CommonLevel3BaoZi = BaoZi
                config.CommonLevel3WinGold = WinGold
                config.CommonLevel3WinRateHigh = WinRateHigh
                config.CommonLevel3LoseGold = LoseGold
                config.CommonLevel3WinRateLow = WinRateLow
                models.SaveMatchConfig(config)
            }
            if GameType == 11 {
                config.ItemLevel1Single = Single
                config.ItemLevel1Double = Double
                config.ItemLevel1ShunZi = ShunZi
                config.ItemLevel1JinHua = JinHua
                config.ItemLevel1ShunJin = ShunJin
                config.ItemLevel1BaoZi = BaoZi
                config.ItemLevel1WinGold = WinGold
                config.ItemLevel1WinRateHigh = WinRateHigh
                config.ItemLevel1LoseGold = LoseGold
                config.ItemLevel1WinRateLow = WinRateLow
                models.SaveMatchConfig(config)
            }
            if GameType == 12 {
                config.ItemLevel2Single = Single
                config.ItemLevel2Double = Double
                config.ItemLevel2ShunZi = ShunZi
                config.ItemLevel2JinHua = JinHua
                config.ItemLevel2ShunJin = ShunJin
                config.ItemLevel2BaoZi = BaoZi
                config.ItemLevel2WinGold = WinGold
                config.ItemLevel2WinRateHigh = WinRateHigh
                config.ItemLevel2LoseGold = LoseGold
                config.ItemLevel2WinRateLow = WinRateLow
                models.SaveMatchConfig(config)
            }
            if GameType == 13 {
                config.ItemLevel3Single = Single
                config.ItemLevel3Double = Double
                config.ItemLevel3ShunZi = ShunZi
                config.ItemLevel3JinHua = JinHua
                config.ItemLevel3ShunJin = ShunJin
                config.ItemLevel3BaoZi = BaoZi
                config.ItemLevel3WinGold = WinGold
                config.ItemLevel3WinRateHigh = WinRateHigh
                config.ItemLevel3LoseGold = LoseGold
                config.ItemLevel3WinRateLow = WinRateLow
                models.SaveMatchConfig(config)
            }
            if GameType == 21 {
                config.SngLevel1Single = Single
                config.SngLevel1Double = Double
                config.SngLevel1ShunZi = ShunZi
                config.SngLevel1JinHua = JinHua
                config.SngLevel1ShunJin = ShunJin
                config.SngLevel1BaoZi = BaoZi
                config.SngLevel1WinGold = WinGold
                config.SngLevel1WinRateHigh = WinRateHigh
                config.SngLevel1LoseGold = LoseGold
                config.SngLevel1WinRateLow = WinRateLow
                models.SaveMatchConfig(config)
            }
            if GameType == 22 {
                config.SngLevel2Single = Single
                config.SngLevel2Double = Double
                config.SngLevel2ShunZi = ShunZi
                config.SngLevel2JinHua = JinHua
                config.SngLevel2ShunJin = ShunJin
                config.SngLevel2BaoZi = BaoZi
                config.SngLevel2WinGold = WinGold
                config.SngLevel2WinRateHigh = WinRateHigh
                config.SngLevel2LoseGold = LoseGold
                config.SngLevel2WinRateLow = WinRateLow
                models.SaveMatchConfig(config)
            }
            if GameType == 23 {
                config.SngLevel3Single = Single
                config.SngLevel3Double = Double
                config.SngLevel3ShunZi = ShunZi
                config.SngLevel3JinHua = JinHua
                config.SngLevel3ShunJin = ShunJin
                config.SngLevel3BaoZi = BaoZi
                config.SngLevel3WinGold = WinGold
                config.SngLevel3WinRateHigh = WinRateHigh
                config.SngLevel3LoseGold = LoseGold
                config.SngLevel3WinRateLow = WinRateLow
                models.SaveMatchConfig(config)
            }
            if GameType == 30 {
                config.WanSingle = Single
                config.WanDouble = Double
                config.WanShunZi = ShunZi
                config.WanJinHua = JinHua
                config.WanShunJin = ShunJin
                config.WanBaoZi = BaoZi
                config.WanWinGold = WinGold
                config.WanWinRateHigh = WinRateHigh
                config.WanLoseGold = LoseGold
                config.WanWinRateLow = WinRateLow
                models.SaveMatchConfig(config)
            }
        }
        c.Data["json"] = result
    } else {
        c.Data["json"] = "failed"
    }
    c.ServeJson()
}


///////////////////////////////////////////////////////////////////////////////////////////////
type SetUserFortuneController struct {
    beego.Controller
}
func (c *SetUserFortuneController) Get() {
    userId := c.GetString("userId")
    gold := c.GetString("gold")
    diamond := c.GetString("diamond")
    c.Data["json"] = http.SetUserFortune(userId, gold, diamond)
    c.ServeJson()
}



///////////////////////////////////////////////////////////////////////////////////////////////
type LockUserController struct {
    beego.Controller
}
func (c *LockUserController) Get() {
    userId := c.GetString("userId")
    isLock := c.GetString("isLock")
    c.Data["json"] = http.LockUser(userId, isLock)
    c.ServeJson()
}



///////////////////////////////////////////////////////////////////////////////////////////////
type SendSystemMessageController struct {
    beego.Controller
}
func (c *SendSystemMessageController) Get() {
    content := c.GetString("content")
    c.Data["json"] = http.SendSystemMessage(content)
    c.ServeJson()
}



///////////////////////////////////////////////////////////////////////////////////////////////
type SendUserPrizeMailController struct {
    beego.Controller
}
func (c *SendUserPrizeMailController) Get() {
    userId := c.GetString("userId")
    content := c.GetString("content")
    gold, _ := c.GetInt("gold")
    diamond, _ := c.GetInt("diamond")
    //exp, _ := c.GetInt("exp")
    //score, _ := c.GetInt("score")
    itemType, _ := c.GetInt("itemType")
    itemCount, _ := c.GetInt("itemCount")
    c.Data["json"] = http.SendUserPrizeMail(userId, content, gold, diamond, 0, 0, itemType, itemCount)
    c.ServeJson()
}



///////////////////////////////////////////////////////////////////////////////////////////////
type QueryUserController struct {
    beego.Controller
}
func (c *QueryUserController) Get() {
    userId := c.GetString("userId")
    nickname := c.GetString("nickname")
    fmt.Println("===================> ", "userId:", userId, "nickname:", nickname)
    resp := http.QueryUser(userId, nickname)
    if nil != resp {
        fmt.Println("resp => ", resp)
        c.Data["json"] = resp.Infos
    } else {
        c.Data["json"] = nil
    }
    c.ServeJson()
}