package http

import (
    "fmt"
    "github.com/astaxie/beego/httplib"
    "encoding/json"
    "github.com/astaxie/beego"
)

const (
    SucceedResult = "succeed"
    FailedResult = "failed"
    TimeoutResult = "timeout"
)

var gameServerAddr, gameServerKey string
var gameServerPort int

type PrizeMail struct {
    UserId string `json:"userId"`
    Content string `json:"content"`
    Gold int `json:"gold"`
    Diamond int `json:"diamond"`
    Exp int `json:"exp"`
    Score int `json:"score"`
    ItemType int `json:"itemType"`
    ItemCount int `json:"itemCount"`
}


type UserInfoResp struct {
    Ok    bool       `json:"ok"`
    Infos []UserInfo `json:"infos"`
}

type OnlineTypeInfo struct {
    GameType    int    `json:"gameType"`
    Count       int    `json:"count"`
}

type OnlineTypeResp struct {
    Infos    []OnlineTypeInfo    `json:"infos"`
}

type GoldLimitUserCount struct {

}

type UserInfo struct {
    UserId     string `json:"userId"`
    Username   string `json:"username"`
    Nickname   string `json:"nickname"`
    Level      int    `json:"level"`
    VipLevel   int    `json:"vipLevel"`
    Gold       int    `json:"gold"`
    Diamond    int    `json:"diamond"`
    CreateTime int64  `json:"createTime"`
}

type CardConfigData struct {
    GameType int `json:"gameType"`
    Single   int `json:"single"`
    Double   int `json:"double"`
    ShunZi   int `json:"shuZi"`
    JinHua   int `json:"jinHua"`
    ShunJin  int `json:"shunJin"`
    BaoZi    int `json:"baoZi"`
    WinGold  int `json:"winGold"`
    WinRateHigh  int `json:"winRateHigh"`
    LoseGold  int `json:"loseGold"`
    WinRateLow  int `json:"winRateLow"`
}


func init() {
    gameServerAddr = beego.AppConfig.String("game_server_addr")
    gameServerPort, _ = beego.AppConfig.Int("game_server_port")
    gameServerKey = beego.AppConfig.String("game_server_key")
    fmt.Println("gameServerAddr=>", gameServerAddr)
    fmt.Println("gameServerPort=>", gameServerPort)
    fmt.Println("gameServerKey=>", gameServerKey)
}

func getUrl() (string) {
    fmt.Println("getUrl => gameServerPort => ", gameServerPort)
    if 0 != gameServerPort {
        return "http://" + gameServerAddr + ":" + fmt.Sprintf("%d", gameServerPort) + "/"
    } else {
        return "http://" + gameServerAddr + "/"
    }
}

func getKeyUrl(apiName string) (string) {
    return getUrl() + apiName + "?key=" + gameServerKey
}

// 设置奖励版本
func SetPrizeVersion(version string) (string) {
    req := httplib.Get(getKeyUrl("setCurVersion"))
    req.Param("version", version)
    req.Debug(true).Response()
    resultStr, err := req.String()
    if nil == err && resultStr == "1" {
        return SucceedResult
    }
    return FailedResult
}

// 设置赛事配置
func SetMatchConfig(gameType, single, double, shunzi, jinHua, shunJin, baoZi, winGold, winRateHigh, loseGold, winRateLow int) (string) {
    req := httplib.Post(getKeyUrl("setCardConfig"))
    matchConfig := CardConfigData{}
    matchConfig.GameType = gameType
    matchConfig.Single = single
    matchConfig.Double = double
    matchConfig.ShunZi = shunzi
    matchConfig.JinHua = jinHua
    matchConfig.ShunJin = shunJin
    matchConfig.BaoZi = baoZi
    matchConfig.WinGold = winGold
    matchConfig.WinRateHigh = winRateHigh
    matchConfig.LoseGold = loseGold
    matchConfig.WinRateLow = winRateLow
    jsonContent, _ := json.Marshal(matchConfig)
    req.Body(jsonContent)
    req.Debug(true).Response()
    resultStr, err := req.String()
    fmt.Println("resultStr => ", resultStr)
    if nil == err && resultStr == "1" {
        return SucceedResult
    }
    return FailedResult
}

func SetUserFortune(userId, gold, diamond string) (string) {
    req := httplib.Get(getKeyUrl("setUserFortune"))
    req.Param("userId", userId)
    req.Param("gold", gold)
    req.Param("diamond", diamond)
    req.Response()
    resultStr, err := req.String()
    if nil == err && resultStr == "1" {
        return SucceedResult
    }
    return FailedResult
}

func LockUser(userId, isLock string) (string) {
    req := httplib.Get(getKeyUrl("lockUser"))
    req.Param("userId", userId)
    // true | false
    req.Param("isLocked", isLock)
    req.Response()
    resultStr, err := req.String()
    if nil == err && resultStr == "1" {
        return SucceedResult
    }
    return FailedResult
}

func SendSystemMessage(content string) (string) {
    req := httplib.Get(getKeyUrl("sendSysMsg"))
    req.Param("msg", content)
    req.Response()
    resultStr, err := req.String()
    if nil == err && resultStr == "1" {
        return SucceedResult
    }
    return FailedResult
}

func SendUserPrizeMail(userId, content string, gold, diamond, exp, score, itemType, itemCount int) (string) {
    req := httplib.Post(getKeyUrl("sendPrizeMail"))
    mailInfo := PrizeMail{}
    mailInfo.UserId = userId
    mailInfo.Content = content
    mailInfo.Gold = gold
    mailInfo.Diamond = diamond
    mailInfo.Exp = exp
    mailInfo.Score = score
    mailInfo.ItemType = itemType
    mailInfo.ItemCount = itemCount
    jsonContent, _ := json.Marshal(mailInfo)
    req.Body(jsonContent)
    req.Debug(true).Response()
    resultStr, err := req.String()
    if nil == err && resultStr == "1" {
        return SucceedResult
    }
    return FailedResult
}

func QueryUser(userId, nickname string) (*UserInfoResp){
    var req *httplib.BeegoHttpRequest
    if len(userId) > 0 {
        req = httplib.Get(getKeyUrl("queryUserById"))
        req.Param("userId", userId)
    } else {
        req = httplib.Get(getKeyUrl("queryUserByName"))
        req.Param("nickname", nickname)
    }
    req.Debug(true).Response()
    resp := &UserInfoResp{}
    req.ToJson(&resp)
    fmt.Println("resp => ", resp)
    if resp.Ok {
        return resp
    }
    return nil
}

func GetMainNowQuery() (OnlineTypeResp){
    var req *httplib.BeegoHttpRequest
    req = httplib.Get(getKeyUrl("onlineByType"))
    req.Debug(true).Response()
    resp := OnlineTypeResp{}
    req.ToJson(&resp)
    fmt.Println("onlineByType resp => ", resp)
    return resp
}

type HostInfoResp struct {
    Date string `json:"Date"`
    TipGold int `json:"TipGold"`
}

func GetHostQuery(b,e string) ([]*HostInfoResp){
    var req *httplib.BeegoHttpRequest
    req = httplib.Get(getKeyUrl("log/getSysTipLog"))
    req.Param("start", b)
    req.Param("end",e)
    req.Debug(true).Response()

    resp := []*HostInfoResp{}
    req.ToJson(&resp)
    fmt.Println("getSysTipLog resp => ", resp)
    return resp
}

type SlotInfoResp struct {
    Date string `json:"Date"`
    PoolGold int `json:"PoolGold"`
}

func GetSlotQuery(b,e string) ([]*SlotInfoResp){
    var req *httplib.BeegoHttpRequest
    req = httplib.Get(getKeyUrl("log/getSlotPoolLog"))
    req.Param("start", b)
    req.Param("end",e)
    req.Debug(true).Response()

    resp := []*SlotInfoResp{}
    req.ToJson(&resp)
    fmt.Println("getSlotPoolLog resp => ", resp)
    return resp
}

type FeeInfoResp struct {
    Date string `json:"Date"`
    FeeGold int `json:"FeeGold"`
}

func GetFeeQuery(b,e string) ([]*FeeInfoResp){
    var req *httplib.BeegoHttpRequest
    req = httplib.Get(getKeyUrl("log/getGameFeeLog"))
    req.Param("start", b)
    req.Param("end",e)
    req.Debug(true).Response()

    resp := []*FeeInfoResp{}
    req.ToJson(&resp)
    fmt.Println("getGameFeeLog resp => ", resp)
    return resp
}

type CharmInfoResp struct {
    AllCharm int `json:"AllCharm"`
    AllExpendCharm int `json:"AllExpendCharm"`
    DayExchanged []*ExpendCharm
    FinePool []*GiftGold
    BadPool []*GiftGold
}
type ExpendCharm struct {
    Date string `json:"Date"`
    ExpendCharm int `json:"ExpendCharm"`
}
type GiftGold struct {
    Date string `json:"Date"`
    GiftGold int `json:"GiftGold"`
}
func GetCharmQuery(b,e string) (CharmInfoResp){
    var req *httplib.BeegoHttpRequest
    req = httplib.Get(getKeyUrl("log/getCharmLog"))
    req.Param("start", b)
    req.Param("end",e)
    req.Debug(true).Response()

    resp := CharmInfoResp{}
    req.ToJson(&resp)
    fmt.Println("getCharmLog resp => ", resp)
    return resp
}