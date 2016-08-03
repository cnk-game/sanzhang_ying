package active

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"util"
)

type UserActive struct {
	UserId string    `bson:"userId"`
	Id     string    `bson:"id"`
	Time   time.Time `bson:"createTime"`
}

const (
	active_christmasC = "active_newyear"
)

func FindUserActive(userId string) (*UserActive, error) {
	active := &UserActive{}
	active.UserId = userId
	err := util.WithUserCollection(active_christmasC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"userId": userId}).One(active)
	})

	return active, err
}

func SaveUserActive(userId string, id string) error {
	active := &UserActive{}
	active.UserId = userId
	active.Id = id
	active.Time = time.Now()

	return util.WithUserCollection(active_christmasC, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"userId": active.UserId}, active)
		return err
	})
}
