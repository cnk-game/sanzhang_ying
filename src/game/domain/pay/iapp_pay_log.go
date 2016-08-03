package pay

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"util"
)

type IAppPayLog struct {
	Transtype int       `bson:"transtype"`
	Cporderid string    `bson:"cporderid"`
	Transid   string    `bson:"transid"`
	Appuserid string    `bson:"appuserid"`
	Appid     string    `bson:"appid"`
	Waresid   int       `bson:"waresid"`
	Feetype   int       `bson:"feetype"`
	Money     float32   `bson:"money"`
	Currency  string    `bson:"currency"`
	Result    int       `bson:"result"`
	Transtime string    `bson:"transtime"`
	Cpprivate string    `bson:"cpprivate"`
	Paytype   int       `bson:"paytype"`
	Time      time.Time `bson:"time"`
}

const (
	iappPayLogC = "iapp_pay_log"
)

func FindIAppPayLog(transid string) (*IAppPayLog, error) {
	l := &IAppPayLog{}
	err := util.WithUserCollection(iappPayLogC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"transid": transid}).One(l)
	})
	return l, err
}

func SaveIAppPayLog(log *IAppPayLog) error {
	log.Time = time.Now()
	return util.WithUserCollection(iappPayLogC, func(c *mgo.Collection) error {
		return c.Insert(log)
	})
}
