package user

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
	"util"
)

type UserLog struct {
	UserId             string    `bson:"userId"`
	UserName           string    `bson:"userName"`
	TotalOnlineSeconds int       `bson:"totalOnlineSeconds"`
	MatchTimes         int       `bson:"matchTimes"`
	CreateTime         time.Time `bson:"createTime"`
	Channel            string    `bson:"channel"`
	Model              string    `bson:"model"`
}

type LoginRecord struct {
	UserId     string    `bson:"userId"`
	UserName   string    `bson:"userName"`
	Channel    string    `bson:"channel"`
	LoginTime  time.Time `bson:"loginTime"`
	LogoutTime time.Time `bson:"logoutTime"`
	LoginIP    string    `bson:"loginIP"`
	DeviceId   string    `bson:"deviceId"`
}

type SlowMsgRecord struct {
	UserId     string    `bson:"userId"`
	MsgId      string    `bson:"msgId"`
	StartTime  time.Time `bson:"startTime"`
	ElapseTime string    `bson:"elapseTime"`
}

const (
	userLogC       = "user_log"
	loginRecordC   = "login_record"
	slowMsgRecordC = "slow_msg"
)

func FindUserLog(userId string) (*UserLog, error) {
	l := &UserLog{}
	l.UserId = userId
	err := util.WithLogCollection(userLogC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"userId": userId}).One(l)
	})
	return l, err
}

func SaveUserLog(l *UserLog) error {
	return util.WithLogCollection(userLogC, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"userId": l.UserId}, l)
		return err
	})
}

func InsertLoginRecord(r *LoginRecord) error {
	now := time.Now()
	cur_C := loginRecordC + "_" + strconv.Itoa(int(now.Year())) + strconv.Itoa(int(now.Month())) + strconv.Itoa(int(now.Day()))
	return util.WithLogCollection(cur_C, func(c *mgo.Collection) error {
		return c.Insert(r)
	})
}

func SaveSlowMsg(userId string, msgId string, startTime time.Time, elapseTime string) error {
	msg := &SlowMsgRecord{}
	msg.UserId = userId
	msg.MsgId = msgId
	msg.StartTime = startTime
	msg.ElapseTime = elapseTime

	return util.WithLogCollection(slowMsgRecordC, func(c *mgo.Collection) error {
		return c.Insert(msg)
	})
}
