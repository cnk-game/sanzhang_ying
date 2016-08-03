package pay

import (
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"util"
)

type Pay_log_result struct {
	Id    int `bson:"_id"`
	Count int `bson:"pay_count"`
}

func GetPayCountBetweenTime(userId string) (error, int) {
	money := 0
	glog.Info("GetPayCountBetweenTime userId=", userId)
	tmp := []bson.M{
		{"$match": bson.M{"userId": userId,
			"payCode": bson.M{"$ne": "100576"},
			"time":    bson.M{"$gte": Pay_Atvice_begin, "$lte": Pay_Atvice_end}}},
		{"$group": bson.M{"_id": "$userId",
			"pay_count": bson.M{"$sum": "$amount"}}}}

	errr := util.WithLogCollection(payLogActiveC, func(c *mgo.Collection) error {
		pipe := c.Pipe(tmp)
		result := []Pay_log_result{{0, 0}}
		err := pipe.All(&result)
		if err != nil {
			return err
		}

		glog.Info("GetPayCountBetweenTime result=", result)
		for _, value := range result {
			money = value.Count
			glog.Info("GetPayCountBetweenTime money=", money)
		}
		return nil
	})

	return errr, money
}
