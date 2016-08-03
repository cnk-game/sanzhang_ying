package pay

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"util"
)

type HmPayLog struct {
	NotifyTime  string    `bson:"notify_time"`
	AppId       string    `bson:"appId"`
	UserId      string    `bson:"userId"`
	OutTradeNo  string    `bson:"out_trade_no"`
	TotalFee    string    `bson:"total_fee"`
	Subject     string    `bson:"subject"`
	Body        string    `bson:"body"`
	TradeStatus string    `bson:"trade_status"`
	CreateTime  time.Time `bson:"create_time"`
}

const (
	hmPayLogC = "hm_pay_log"
)

func FindHmPayLog(out_trade_no string) (*HmPayLog, error) {
	l := &HmPayLog{}
	err := util.WithUserCollection(hmPayLogC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"out_trade_no": out_trade_no}).One(l)
	})
	return l, err
}

func SaveHmPayLog(log *HmPayLog) error {
	log.CreateTime = time.Now()
	return util.WithUserCollection(hmPayLogC, func(c *mgo.Collection) error {
		return c.Insert(log)
	})
}
