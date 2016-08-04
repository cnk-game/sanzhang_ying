package models

import (
    "time"
    "poker-admin/util"
    "fmt"
    "labix.org/v2/mgo/bson"
)

type NewUserLog struct {
    UserId string `bson:"userId"`
    UserName string `bson:"userName"`
    Channel string `bson:"channel"`
    Model string `bson:"model"`
    TotalOnlineSeconds int `bson:"totalOnlineSeconds"`
    MatchTimes int `bson:"matchTimes"`
    CreateTime time.Time `bson:"createTime"`
}

type LoginRecord struct {
    UserId string `bson:"userId"`
    UserName string `bson:"userName"`
    Channel string `bson:"channel"`
    LoginTime time.Time `bson:"loginTime"`
    LogoutTime time.Time `bson:"logoutTime"`
}


// 总用户数
func GetTotalPlayerCount(channel string) (int) {
    session := util.GetLogSession()
    c := session.DB("poker_log").C("user_log")
    defer session.Close()
    var count int = 0
    var err error
    if len(channel) > 0 {
        count, err = c.Find(bson.M{"channel":channel}).Count()
    } else {
        count, err = c.Find(bson.M{}).Count()
    }
    if err != nil {
        fmt.Println("GetTotalPlayerCount => error")
    }
    return count
}

// 指定日期内的所有新玩家
func GetAllNewPlayers(year, month, day int, result interface{}, channel string) {
    session := util.GetLogSession()
    c := session.DB("poker_log").C("user_log")
    defer session.Close()
    
    nextOneDay := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
    nextOneDay = nextOneDay.AddDate(0, 0, 1)
    begin := util.ParseTime(fmt.Sprintf("%v-%02v-%02v 00:00:00", year, month, day))
    end := util.ParseTime(fmt.Sprintf("%v-%02v-%02v 00:00:00", nextOneDay.Year(), int(nextOneDay.Month()), nextOneDay.Day()))
    var err error
    if IsAdminChannel(channel) {
        err = c.Find(bson.M{"createTime":bson.M{"$gte": begin, "$lt": end}}).All(result)
    } else {
        err = c.Find(bson.M{"channel":channel, "createTime":bson.M{"$gte": begin, "$lt": end}}).All(result)
    }
    if err != nil {
        fmt.Println("GetAllNewPlayers => error")
    }
}

func GetNewPlayers(b_year, b_month, b_day, b_hour, b_minute, e_year, e_month, e_day, e_hour, e_minute, pageIdx int, result interface{}, channel string) (int) {
    session := util.GetLogSession()
    c := session.DB("poker_log").C("user_log")
    defer session.Close()
    
    begin := util.ParseTime(fmt.Sprintf("%v-%02v-%02v %02v:%02v:00", b_year, b_month, b_day, b_hour, b_minute))
    end := util.ParseTime(fmt.Sprintf("%v-%02v-%02v %02v:%02v:00", e_year, e_month, e_day, e_hour, e_minute))

    var err error
    if IsAdminChannel(channel) {
        err = c.Find(bson.M{"createTime":bson.M{"$gte": begin, "$lt": end}}).Sort("-createTime").Skip(page_count * pageIdx).Limit(page_count).All(result)
    } else {
        err = c.Find(bson.M{"channel":channel, "createTime":bson.M{"$gte": begin, "$lt": end}}).Sort("-createTime").Skip(page_count * pageIdx).Limit(page_count).All(result)
    }
    if err != nil {
        fmt.Println("GetNewPlayers => error1")
    }

    count := 0
    if IsAdminChannel(channel) {
        count, err = c.Find(bson.M{"createTime":bson.M{"$gte": begin, "$lt": end}}).Count()
    } else {
        count, err = c.Find(bson.M{"channel":channel, "createTime":bson.M{"$gte": begin, "$lt": end}}).Count()
    }
    if err != nil {
        fmt.Println("GetNewPlayers => error2")
    }

    return count
}

// 指定日期内的登陆记录
func GetLoginRecords(year, month, day int, result interface{}) {
    session := util.GetLogSession()
    c := session.DB("poker_log").C("login_record")
    defer session.Close()

    nextOneDay := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
    nextOneDay = nextOneDay.AddDate(0, 0, 1)
    begin := util.ParseTime(fmt.Sprintf("%v-%02v-%02v 00:00:00", year, month, day))
    end := util.ParseTime(fmt.Sprintf("%v-%02v-%02v 00:00:00", nextOneDay.Year(), int(nextOneDay.Month()), nextOneDay.Day()))
    err := c.Find(bson.M{"loginTime":bson.M{"$gte": begin, "$lt": end}}).All(result)
    if err != nil {
        fmt.Println("getLoginRecords ======> error")
    }
}