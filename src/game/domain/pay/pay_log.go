package pay

import (
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"util"
)

type PayLog struct {
	OrderId    string    `bson:"orderId"`
	UserId     string    `bson:"userId"`
	Amount     int       `bson:"amount"`
	PayChannel string    `bson:"payChannel"`
	PayType    string    `bson:"payType"`
	Channel    string    `bson:"channel"`
	Time       time.Time `bson:"time"`
	PayCode    string    `bson:"payCode"`
}

const (
	payLogC          = "pay_log"
	payLogActiveC    = "pay_log_active"
	payActiveAddLogC = "pay_active_add_log"
	payActiveConfigC = "pay_active_config"
)

//0206 1454659200
//0221 1455955200
var (
	Pay_Atvice_begin     = time.Unix(1454688000, 0)
	Pay_Atvice_end       = time.Unix(1455984000, 0)
	Pay_Atvice_Begin_Int = int64(1454688000)
	Pay_Atvice_End_Int   = int64(1455984000)
)

func SavePayLog(log *PayLog) error {
	log.Time = time.Now()
	return util.WithLogCollection(payLogC, func(c *mgo.Collection) error {
		SavePayActiveLog(log)
		return c.Insert(log)
	})
}

func SavePayActiveLog(log *PayLog) error {
	tm := int64(time.Now().Unix())
	if tm < Pay_Atvice_Begin_Int || tm > Pay_Atvice_End_Int {
		return nil
	}
	log.Time = time.Now()
	return util.WithLogCollection(payLogActiveC, func(c *mgo.Collection) error {
		return c.Insert(log)
	})
}

type PayActiveAddCoinsLog struct {
	UserId string    `bson:"userId"`
	Money  int       `bson:"money"`
	Coins  int       `bson:"coins"`
	Time   time.Time `bson:"time"`
}

func SavePayAcitveAddCoinsLog(userId string, money int, coins int) error {
	tm := &PayActiveAddCoinsLog{}
	tm.UserId = userId
	tm.Money = money
	tm.Coins = coins
	tm.Time = time.Now()
	return util.WithLogCollection(payActiveAddLogC, func(c *mgo.Collection) error {
		return c.Insert(tm)
	})
}

func GetPayAcitveAddCoinsLog(userId string, money int) error {
	tm := &PayActiveAddCoinsLog{}
	err := util.WithLogCollection(payActiveAddLogC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"userId": userId, "money": money}).One(tm)
	})
	return err
}

type PayActiveConfig struct {
	BeginTime int64 `bson:"beginTime"`
	EndTime   int64 `bson:"endTime"`
}

func GetPayActiveConfig() error {
	tm := &PayActiveConfig{}
	err := util.WithGameCollection(payActiveConfigC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"id": 1}).One(tm)
	})

	if err == mgo.ErrNotFound {
		glog.Error("GetPayActiveConfig ErrNotFound")
		return nil
	} else {
		glog.Error("GetPayActiveConfig ", tm)
	}
	Pay_Atvice_begin = time.Unix(tm.BeginTime, 0)
	Pay_Atvice_end = time.Unix(tm.EndTime, 0)
	Pay_Atvice_Begin_Int = int64(tm.BeginTime)
	Pay_Atvice_End_Int = int64(tm.EndTime)
	return err
}
