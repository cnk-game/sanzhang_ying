package user

import (
	"fmt"
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"strings"
	"util"
)

const (
	userIdC = "user_id_index"
)

type sequence struct {
	Name string `bson:"name"`
	Val  int64  `bson:"seq"`
}

func GetNewUserId() string {
	ch := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"seq": 1}},
		ReturnNew: true,
	}

	next := &sequence{}
	err := util.WithUserCollection(userIdC, func(c *mgo.Collection) error {
		_, err := c.Find(bson.M{"name": "index"}).Apply(ch, &next)
		return err
	})

	if err != nil {
		return "error"
	}

	glog.Info("GetNewUserId in. next=", next)
	newId := "QF" + strconv.Itoa(int(next.Val))

	return newId
}

func GetUserIdTable(userIdStr string) (string, error) {
	glog.Info("GetUserIdTable userIdInt:", userIdStr)
	tem := strings.Trim(userIdStr, "QF1")
	glog.Info("GetUserIdTable tem:", tem)
	userIdInt, err := strconv.Atoi(tem)
	if err != nil {
		glog.Info("GetUserIdTable err:", err)
		return "", err
	}
	glog.Info("GetUserIdTable userIdInt:", userIdInt)
	index := userIdInt/50000 + 1
	glog.Info("GetUserIdTable index:", index)
	tableName := "user_" + fmt.Sprintf("%v", index)
	return tableName, nil
}

func GetUserFortuneTable(userIdStr string) (string, error) {
	tem := strings.Trim(userIdStr, "QF1")
	userIdInt, err := strconv.Atoi(tem)
	if err != nil {
		glog.Info("GetUserFortuneTable err:", err)
		return "", err
	}

	index := userIdInt/30000 + 1
	tableName := "user_fortune_" + fmt.Sprintf("%v", index)
	return tableName, nil
}

func GetUserMatchRecordTable(userIdStr string) (string, error) {
	tem := strings.Trim(userIdStr, "QF1")
	userIdInt, err := strconv.Atoi(tem)
	if err != nil {
		glog.Info("GetUserMatchRecordTable err:", err)
		return "", err
	}

	index := userIdInt/100000 + 1
	tableName := "match_record_" + fmt.Sprintf("%v", index)
	return tableName, nil
}
