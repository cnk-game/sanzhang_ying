package user

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
	"util"
)

type FortuneLog struct {
	UserId     string    `bson:"userId"`
	Gold       int       `bson:"gold"`
	CurGold    int64     `bson:"curGold"`
	Diamond    int       `bson:"diamond"`
	CurDiamond int       `bson:"curDiamond"`
	Score      int       `bson:"score"`
	CurScore   int       `bson:"curScore"`
	Reason     string    `bson:"reason"`
	Time       time.Time `bson:"time"`
}

type GiftFishLog struct {
	FromId   string    `bson:"fromId"`
	ToId     string    `bson:"toId"`
	FishType int       `bson:"fishType"`
	Count    int       `bson:"count"`
	Time     time.Time `bson:"time"`
}

const (
	earnFortuneLogC    = "earn_fortune_log"
	consumeFortuneLogC = "consume_fortune_log"
	giftFishLogC       = "gift_fish"
)

func FindEarnFortuneLogs() ([]*FortuneLog, error) {
	logs := []*FortuneLog{}
	err := util.WithLogCollection(earnFortuneLogC, func(c *mgo.Collection) error {
		return c.Find(nil).All(&logs)
	})

	return logs, err
}

func SaveEarnFortuneLog(log *FortuneLog) error {
	now := time.Now()
	cur_C := earnFortuneLogC + "_" + strconv.Itoa(int(now.Year())) + strconv.Itoa(int(now.Month())) + strconv.Itoa(int(now.Day()))
	log.Time = now
	return util.WithLogCollection(cur_C, func(c *mgo.Collection) error {
		return c.Insert(log)
	})
}

func SaveGiftFishLog(fromId, toId string, fishType, count int) error {
	log := GiftFishLog{fromId, toId, fishType, count, time.Now()}
	return util.WithLogCollection(giftFishLogC, func(c *mgo.Collection) error {
		return c.Insert(log)
	})
}

func LoadGiftFishLog(userId string) ([]*GiftFishLog, error) {
	logs := []*GiftFishLog{}
	err := util.WithLogCollection(giftFishLogC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"fromId": userId}).All(&logs)
	})
	return logs, err
}

func FindConsumeFortuneLogs() ([]*FortuneLog, error) {
	logs := []*FortuneLog{}
	err := util.WithLogCollection(consumeFortuneLogC, func(c *mgo.Collection) error {
		return c.Find(nil).All(&logs)
	})

	return logs, err
}

func SaveConsumeFortuneLog(log *FortuneLog) error {
	now := time.Now()
	cur_C := consumeFortuneLogC + "_" + strconv.Itoa(int(now.Year())) + strconv.Itoa(int(now.Month())) + strconv.Itoa(int(now.Day()))
	log.Time = now
	return util.WithLogCollection(cur_C, func(c *mgo.Collection) error {
		return c.Insert(log)
	})
}
