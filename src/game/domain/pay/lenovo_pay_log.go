package pay

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"util"
)

type LenovoPayLog struct {
	OrderId         string    `bson:"orderId"`
	MerchantOrderId string    `bson:"merchantOrderId"`
	Amount          int       `bson:"amount"`
	AppId           string    `bson:"appId"`
	PayTime         string    `bson:"payTime"`
	Attach          string    `bson:"attach"`
	Status          string    `bson:"status"`
	Sign            string    `bson:"sign"`
	Time            time.Time `bson:"time"`
}

const (
	lenovoPayLogC = "lenovo_pay_log"
)

func FindLenovoPayLog(orderId string) (*LenovoPayLog, error) {
	l := &LenovoPayLog{}
	err := util.WithUserCollection(lenovoPayLogC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"orderId": orderId}).One(l)
	})
	return l, err
}

func SaveLenovoPayLog(log *LenovoPayLog) error {
	log.Time = time.Now()
	return util.WithUserCollection(lenovoPayLogC, func(c *mgo.Collection) error {
		return c.Insert(log)
	})
}
