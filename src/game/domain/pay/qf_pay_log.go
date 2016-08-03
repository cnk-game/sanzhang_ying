package pay

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"util"
)

type QfPayLog struct {
	AppId      string    `bson:"appId"`
	UserId     string    `bson:"userId"`
	Order      string    `bson:"order"`
	Price      string    `bson:"price"`
	PayType    string    `bson:"payType"`
	PayCode    string    `bson:"payCode"`
	State      string    `bson:"state"`
	Time       string    `bson:"time"`
	GameOrder  string    `bson:"gameOrder"`
	CreateTime time.Time `bson:"createTime"`
}

const (
	qfPayLogC = "qf_pay_log"
)

func FindQfPayLog(orderId string) (*QfPayLog, error) {
	l := &QfPayLog{}
	err := util.WithUserCollection(qfPayLogC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"order": orderId}).One(l)
	})
	return l, err
}

func SaveQfPayLog(log *QfPayLog) error {
	log.CreateTime = time.Now()
	return util.WithUserCollection(qfPayLogC, func(c *mgo.Collection) error {
		return c.Insert(log)
	})
}
