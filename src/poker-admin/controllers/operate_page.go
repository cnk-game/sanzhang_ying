package controllers

import (
    "github.com/astaxie/beego"
    "fmt"
    "poker-admin/models/http"
)


///////////////////////////////////////////////////////////////
type HostMgrController struct {
    beego.Controller
}
func (c *HostMgrController) Get() {
    fmt.Println("=======================> prize version : ", "host")
    c.TplNames = "operate/host.html"
}

///////////////////////////////////////////////////////////////
type HostQueryController struct {
    beego.Controller
}
func (c *HostQueryController) Get() {
    bYear:= c.GetString("bYear")
    bMonth := c.GetString("bMonth")
    bDay := c.GetString("bDay")
    eYear := c.GetString("eYear")
    eMonth := c.GetString("eMonth")
    eDay := c.GetString("eDay")

    c.Data["json"] = http.GetHostQuery(bYear+"-"+bMonth+"-"+ bDay, eYear+"-"+eMonth+"-"+eDay);
    c.ServeJson()
}

///////////////////////////////////////////////////////////////
type SlotQueryController struct {
    beego.Controller
}
func (c *SlotQueryController) Get() {
    bYear:= c.GetString("bYear")
    bMonth := c.GetString("bMonth")
    bDay := c.GetString("bDay")
    eYear := c.GetString("eYear")
    eMonth := c.GetString("eMonth")
    eDay := c.GetString("eDay")

    c.Data["json"] = http.GetSlotQuery(bYear+"-"+bMonth+"-"+ bDay, eYear+"-"+eMonth+"-"+eDay);
    c.ServeJson()
}

///////////////////////////////////////////////////////////////
type FeeQueryController struct {
    beego.Controller
}
func (c *FeeQueryController) Get() {
    bYear:= c.GetString("bYear")
    bMonth := c.GetString("bMonth")
    bDay := c.GetString("bDay")
    eYear := c.GetString("eYear")
    eMonth := c.GetString("eMonth")
    eDay := c.GetString("eDay")

    c.Data["json"] = http.GetFeeQuery(bYear+"-"+bMonth+"-"+ bDay, eYear+"-"+eMonth+"-"+eDay);
    c.ServeJson()
}


///////////////////////////////////////////////////////////////
type CharmQueryController struct {
    beego.Controller
}
func (c *CharmQueryController) Get() {
    bYear:= c.GetString("bYear")
    bMonth := c.GetString("bMonth")
    bDay := c.GetString("bDay")
    eYear := c.GetString("eYear")
    eMonth := c.GetString("eMonth")
    eDay := c.GetString("eDay")

    c.Data["json"] = http.GetCharmQuery(bYear+"-"+bMonth+"-"+ bDay, eYear+"-"+eMonth+"-"+eDay);
    c.ServeJson()
}
///////////////////////////////////////////////////////////////
type SlotMgrController struct {
    beego.Controller
}
func (c *SlotMgrController) Get() {
    fmt.Println("=======================> prize version : ", "slot")
    c.TplNames = "operate/slot.html"
}

///////////////////////////////////////////////////////////////
type FeeMgrController struct {
    beego.Controller
}
func (c *FeeMgrController) Get() {
    fmt.Println("=======================> prize version : ", "fee")
    c.TplNames = "operate/fee.html"
}

///////////////////////////////////////////////////////////////
type CharmMgrController struct {
    beego.Controller
}
func (c *CharmMgrController) Get() {
    fmt.Println("=======================> prize version : ", "charm")
    c.TplNames = "operate/charm.html"
}

