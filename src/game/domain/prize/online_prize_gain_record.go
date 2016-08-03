package prize

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"util"
)

type OnlinePrizeGainRecord struct {
	UserId   string    `bson:"userId"`
	PrizeId  int       `bson:"prizeId"`
	GainTime time.Time `bson:"gainTime"`
	hashCode *util.HashCode
}

func (r *OnlinePrizeGainRecord) HashCode() *util.HashCode {
	return r.hashCode
}

func (r *OnlinePrizeGainRecord) SetHashCode(hashCode *util.HashCode) {
	r.hashCode = hashCode
}

const (
	onlinePrizeGainRecordC = "online_prize_gain_record"
)

func FindOnlinePrizeGainRecords(userId string) ([]*OnlinePrizeGainRecord, error) {
	rs := []*OnlinePrizeGainRecord{}
	err := util.WithUserCollection(onlinePrizeGainRecordC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"userId": userId}).All(&rs)
	})

	if err == nil {
		for _, r := range rs {
			r.SetHashCode(util.NewHashCode(r))
		}
	}
	return rs, err
}

func SaveOnlinePrizeGainRecord(r *OnlinePrizeGainRecord) error {
	hashCode := util.NewHashCode(r)
	if r.HashCode() != nil && r.HashCode().Compare(hashCode) {
		return nil
	}

	return util.WithSafeUserCollection(onlinePrizeGainRecordC, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"userId": r.UserId, "prizeId": r.PrizeId}, r)
		if err == nil {
			// 保存成功
			r.SetHashCode(hashCode)
		}
		return err
	})
}
