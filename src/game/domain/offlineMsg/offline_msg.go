package offlineMsg

import (
	"code.google.com/p/goprotobuf/proto"
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"util"
)

type OfflineMsg struct {
	UserId  string    `bson:"userId"`
	MsgId   int32     `bson:"msgId"`
	MsgBody []byte    `bson:"msg"`
	Time    time.Time `bson:"time"`
}

const (
	offlineMsgC    = "offline_msg"
	offlineMsgLogC = "offline_msg_log"
)

func FindOfflineMsg(userId string) ([]*OfflineMsg, error) {
	defer RemoveOfflineMsg(userId)

	msg := []*OfflineMsg{}
	err := util.WithUserCollection(offlineMsgC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"userId": userId}).All(&msg)
	})
	return msg, err
}

func SaveOfflineMsg(msg *OfflineMsg) error {
	msg.Time = time.Now()
	return util.WithUserCollection(offlineMsgC, func(c *mgo.Collection) error {
		return c.Insert(msg)
	})
}

func RemoveOfflineMsg(userId string) error {
	return util.WithUserCollection(offlineMsgC, func(c *mgo.Collection) error {
		_, err := c.RemoveAll(bson.M{"userId": userId})
		return err
	})
}

func PutOfflineMsg(userId string, msgId int32, msgBody proto.Message) {
	msg := &OfflineMsg{}
	msg.UserId = userId
	msg.MsgId = msgId

	if msgBody != nil {
		b, err := proto.Marshal(msgBody)
		if err != nil {
			glog.Error(err)
			return
		}
		msg.MsgBody = b
	}

	SaveOfflineMsg(msg)
}

type PrizeMailLog struct {
	UserId  string    `json:"userId"`
	Gold    int       `json:"gold"`
	Diamond int       `json:"diamond"`
	LogTime time.Time `json:logTime`
}

func SaveOfflineMsgLog(userId string, gold int, diamond int) error {
	msg := PrizeMailLog{userId, gold, diamond, time.Now()}
	err := util.WithUserCollection(offlineMsgLogC, func(c *mgo.Collection) error {
		return c.Insert(&msg)
	})
	return err
}
