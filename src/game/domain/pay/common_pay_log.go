package pay

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"util"
)

type CommonPayLog struct {
	UserId     string    `bson:"userId"`
	Order      string    `bson:"order"`
	Price      string    `bson:"price"`
	PayType    string    `bson:"payType"`
	PayCode    string    `bson:"payCode"`
	State      string    `bson:"state"`
	Time       time.Time `bson:"time"`
	ThirdOrder string    `bson:"thirdOrder"`
	Channel    string    `bson:"channel"`
}

const (
	commonPayLogC = "common_pay_log"
)

func FindCommonPayLog(orderId string) (*CommonPayLog, error) {
	l := &CommonPayLog{}
	err := util.WithUserCollection(commonPayLogC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"order": orderId}).One(l)
	})
	return l, err
}

func SaveCommonPayLog(log *CommonPayLog) error {
	log.Time = time.Now()
	return util.WithUserCollection(commonPayLogC, func(c *mgo.Collection) error {
		return c.Insert(log)
	})
}
