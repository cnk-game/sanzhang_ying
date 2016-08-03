package cdkey

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"util"
)

type CDKeyGainRecord struct {
	UserId  string `bson:"userId"`
	Records []int  `bson:"records"`
}

func (r *CDKeyGainRecord) IsGained(keyType int) bool {
	for _, t := range r.Records {
		if t == keyType {
			return true
		}
	}
	return false
}

func (r *CDKeyGainRecord) SetGained(keyType int) {
	if r.IsGained(keyType) {
		return
	}
	r.Records = append(r.Records, keyType)
}

const (
	cdKeyGainRecordC = "cd_key_gain_record"
)

func FindCdKeyGainRecord(userId string) (*CDKeyGainRecord, error) {
	r := &CDKeyGainRecord{}
	r.UserId = userId
	err := util.WithUserCollection(cdKeyGainRecordC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"userId": userId}).One(r)
	})
	return r, err
}

func SaveCdKeyGainRecord(r *CDKeyGainRecord) error {
	return util.WithSafeUserCollection(cdKeyGainRecordC, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"userId": r.UserId}, r)
		return err
	})
}
