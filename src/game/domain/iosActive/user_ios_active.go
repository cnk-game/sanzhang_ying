package iosActive

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"util"
)

type UserIosActive struct {
	UserId string    `bson:"userId"`
	Gold   int64     `bson:"gold"`
	Time   time.Time `bson:"updateTime"`
}

type IosActiveContent struct {
	Content   string `bson:"content"`
	BeginTime string `bson:"beginTime"`
	EndTime   string `bson:"endTime"`
}

const (
	active_iosC         = "active_ios"
	active_ios_contentC = "active_ios_content"
)

func FindUserIosActive(userId string) (*UserIosActive, error) {
	active := &UserIosActive{}
	active.UserId = userId
	err := util.WithGameCollection(active_iosC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"userId": userId}).One(active)
	})

	return active, err
}

func SaveUserIosActive(userId string, gold int64) error {
	active := &UserIosActive{}
	active.UserId = userId
	active.Gold = gold
	active.Time = time.Now()

	return util.WithGameCollection(active_iosC, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"userId": active.UserId}, active)
		return err
	})
}

func GetIosActiveContent() (*IosActiveContent, error) {
	active := &IosActiveContent{}
	err := util.WithGameCollection(active_ios_contentC, func(c *mgo.Collection) error {
		return c.Find(bson.M{}).One(active)
	})

	return active, err
}
