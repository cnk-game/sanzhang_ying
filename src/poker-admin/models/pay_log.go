package models

import (
	"fmt"
	"labix.org/v2/mgo/bson"
	"poker-admin/util"
	"time"
)

type PayLog struct {
	OrderId    string    `bson:"orderId"`
	UserId     string    `bson:"userId"`
	Amount     int       `bson:"amount"`
	PayChannel string    `bson:"payChannel"`
	PayType    string    `bson:"payType"`
	Channel    string    `bson:"channel"`
	Time       time.Time `bson:"time"`
}

func GetPayLogList(b_year, b_month, b_day, b_hour, b_minute, b_second, e_year, e_month, e_day, e_hour, e_minute, e_second, pageIdx int, result interface{}, channel string) int {
	session := util.GetLogSession()
	c := session.DB(util.LogDbName).C(pay_logC)
	defer session.Close()

	begin := util.ParseTime(fmt.Sprintf("%v-%02v-%02v %02v:%02v:%02v", b_year, b_month, b_day, b_hour, b_minute, b_second))
	end := util.ParseTime(fmt.Sprintf("%v-%02v-%02v %02v:%02v:%02v", e_year, e_month, e_day, e_hour, e_minute, e_second))

	var err error
	if IsAdminChannel(channel) {
		err = c.Find(bson.M{"time": bson.M{"$gte": begin, "$lt": end}}).Sort("time").Skip(page_count * pageIdx).Limit(page_count).All(result)
	} else {
		err = c.Find(bson.M{"payChannel": channel, "time": bson.M{"$gte": begin, "$lt": end}}).Sort("time").Skip(page_count * pageIdx).Limit(page_count).All(result)
	}
	if err != nil {
		fmt.Println("GetPayLogList => error1")
	}

	var count int = 0
	if IsAdminChannel(channel) {
		count, err = c.Find(bson.M{"time": bson.M{"$gte": begin, "$lt": end}}).Count()
	} else {
		count, err = c.Find(bson.M{"payChannel": channel, "time": bson.M{"$gte": begin, "$lt": end}}).Count()
	}
	if err != nil {
		fmt.Println("GetPayLogList => error2")
	}
	//fmt.Println("query count => ", count, page_count, pageIdx)

	return count
}

func GetAllPayLogList(year, month, day int, result interface{}) {
	session := util.GetLogSession()
	c := session.DB(util.LogDbName).C(pay_logC)
	defer session.Close()
	nextOneDay := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	nextOneDay = nextOneDay.AddDate(0, 0, 1)
	begin := util.ParseTime(fmt.Sprintf("%v-%02v-%02v 00:00:00", year, month, day))
	end := util.ParseTime(fmt.Sprintf("%v-%02v-%02v 00:00:00", nextOneDay.Year(), int(nextOneDay.Month()), nextOneDay.Day()))
	err := c.Find(bson.M{"time": bson.M{"$gte": begin, "$lt": end}}).Sort("loginTime").All(result)
	if err != nil {
		fmt.Println("GetAllPayLogList => error")
	}
}

func GetPayLogByUserId(userId string, pageIdx int, result interface{}, channel string) int {
	session := util.GetLogSession()
	c := session.DB(util.LogDbName).C(pay_logC)
	defer session.Close()

	var err error
	if IsAdminChannel(channel) {
		err = c.Find(bson.M{"userId": userId}).Sort("loginTime").Skip(page_count * pageIdx).Limit(page_count).All(result)
	} else {
		err = c.Find(bson.M{"userId": userId, "payChannel": channel}).Sort("loginTime").Skip(page_count * pageIdx).Limit(page_count).All(result)
	}
	if err != nil {
		fmt.Println("GetPayLogByUserId => error1")
	}

	var count int = 0
	if IsAdminChannel(channel) {
		count, err = c.Find(bson.M{"userId": userId}).Count()
	} else {
		count, err = c.Find(bson.M{"userId": userId, "payChannel": channel}).Count()
	}

	if err != nil {
		fmt.Println("GetPayLogByUserId => error2")
	}
	return count
}

func GetPayLogByOrderId(orderId, channel string) *PayLog {
	session := util.GetLogSession()
	c := session.DB(util.LogDbName).C(pay_logC)
	defer session.Close()
	payLog := &PayLog{}
	var err error
	if IsAdminChannel(channel) {
		err = c.Find(bson.M{"orderId": orderId}).One(payLog)
	} else {
		err = c.Find(bson.M{"orderId": orderId, "payChannel": channel}).One(payLog)
	}
	if err != nil {
		fmt.Println("GetPayLogByOrderId error")
		return nil
	}
	return payLog
}

/*
计算数值
返回值：付费人数、付费次数、付费金额
*/
func CalcPay(b_year, b_monty, b_day, b_hour, b_minute, b_second, e_year, e_month, e_day, e_hour, e_minute, e_second int) (int, int, int) {
	session := util.GetLogSession()
	c := session.DB(util.LogDbName).C(pay_logC)
	defer session.Close()
	beginTime := util.ParseTime(fmt.Sprintf("%v-%02v-%02v %02v:%02v:%02v", b_year, b_monty, b_day, b_hour, b_minute, b_second))
	endTime := util.ParseTime(fmt.Sprintf("%v-%02v-%02v %02v:%02v:%02v", e_year, e_month, e_day, e_hour, e_minute, e_second))

	count, err := c.Find(bson.M{"time": bson.M{"$gte": beginTime, "$lt": endTime}}).Count()
	if err != nil {
		fmt.Println("calcPay error => 1")
		return 0, 0, 0
	}
	if count > 2000 {
		fmt.Println("calcPay error => 2")
		return 0, 0, 0
	}
	logs := []*PayLog{}
	err = c.Find(bson.M{"time": bson.M{"$gte": beginTime, "$lt": endTime}}).Sort("loginTime").All(&logs)
	if err != nil {
		fmt.Println("calcPay error => 3")
		return 0, 0, 0
	}
	totalPayCount := len(logs)
	totalPayAmount := 0
	playerPayMap := make(map[string]int)
	for _, item := range logs {
		totalPayAmount += item.Amount
		payCount, exists := playerPayMap[item.UserId]
		if exists {
			playerPayMap[item.UserId] = payCount + 1
		} else {
			playerPayMap[item.UserId] = 1
		}
	}
	return len(playerPayMap), totalPayCount, totalPayAmount
}
