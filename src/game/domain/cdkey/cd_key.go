package cdkey

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"util"
)

type CDKey struct {
	Type      int    `bson:"type"`
	Key       string `bson:"key"`
	Gold      int    `bson:"gold"`
	Diamond   int    `bson:"diamond"`
	Score     int    `bson:"score"`
	ItemType  int    `bson:"itemType"`
	ItemCount int    `bson:"itemCount"`
}

const (
	cdKeyC = "cd_key"
)

func FindCDKey(cdKey string) (*CDKey, error) {
	key := &CDKey{}
	err := util.WithUserCollection(cdKeyC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"key": cdKey}).One(key)
	})
	return key, err
}

func RemoveCDKey(cdKey string) error {
	return util.WithUserCollection(cdKeyC, func(c *mgo.Collection) error {
		return c.Remove(bson.M{"key": cdKey})
	})
}
