package models

import (
	"fmt"
	"labix.org/v2/mgo/bson"
	"poker-admin/util"
	"time"
)

const (
	online_logC = "online_log"
)

type OnlineLog struct {
	OnlinePlayerCount int       `bson:"onlinePlayerCount"`
	DateTime          time.Time `bson:"datetime"`
}

func LoadTodayOnlineLog(isGetAll bool, lastTime time.Time, result interface{}) {
	session := util.GetLogSession()
	c := session.DB(util.LogDbName).C(online_logC)
	defer session.Close()

	nextOneDay := time.Date(lastTime.Year(), lastTime.Month(), lastTime.Day(), 0, 0, 0, 0, time.Local)
	nextOneDay = nextOneDay.AddDate(0, 0, 1)
	end := util.ParseTime(fmt.Sprintf("%v-%02v-%02v 00:00:00", nextOneDay.Year(), int(nextOneDay.Month()), nextOneDay.Day()))
	var begin time.Time
	if isGetAll {
		begin = util.ParseTime(fmt.Sprintf("%v-%02v-%02v 00:00:00", lastTime.Year(), int(lastTime.Month()), lastTime.Day()))
	} else {
		begin = util.ParseTime(fmt.Sprintf("%v-%02v-%02v %02v:%02v:%02v", lastTime.Year(), int(lastTime.Month()), lastTime.Day(), lastTime.Hour(), lastTime.Minute(), lastTime.Second()))
	}
	fmt.Println("online log => ", begin, end)
	err := c.Find(bson.M{"datetime": bson.M{"$gte": begin, "$lt": end}}).All(result)
	if err != nil {
		fmt.Println("LoadTodayOnlineLog => error")
	}
}
