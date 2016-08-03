package prize

import (
	"code.google.com/p/goprotobuf/proto"
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"pb"
	"time"
	"util"
)

type SignInRecord struct {
	UserId       string    `bson:"userId"`
	Day1         bool      `bson:"day1"`
	Day2         bool      `bson:"day2"`
	Day3         bool      `bson:"day3"`
	Day4         bool      `bson:"day4"`
	Day5         bool      `bson:"day5"`
	LastSignTime time.Time `bson:"lastSignTime"`
}

const (
	signInRecordC = "sign_record"
)

func FindSignInRecord(userId string) (*SignInRecord, error) {
	r := &SignInRecord{}
	r.UserId = userId
	err := util.WithUserCollection(signInRecordC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"userId": userId}).One(r)
	})
	glog.V(2).Info("===>读取签到记录userId:", userId, " r:", r, " err:", err)
	return r, err
}

func SaveSignInRecord(r *SignInRecord) error {
	return util.WithSafeUserCollection(signInRecordC, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"userId": r.UserId}, r)
		return err
	})
}

func (r *SignInRecord) IsSignIn() bool {
	return util.CompareDate(r.LastSignTime, time.Now())
}

func (r *SignInRecord) ResetSignRecord() {
	t := time.Now()
	if util.CompareDate(r.LastSignTime, t) {
		return
	}

	dur, _ := time.ParseDuration("-24h")
	if !util.CompareDate(r.LastSignTime, t.Add(dur)) {
		glog.V(2).Info("====>重置签到记录lastSignTime:", r.LastSignTime, " 昨天:", t.Add(dur))

		r.Day1 = false
		r.Day2 = false
		r.Day3 = false
		r.Day4 = false
		r.Day5 = false
	}
}

func (r *SignInRecord) SetSignIn() int {
	r.LastSignTime = time.Now()

	if !r.Day1 {
		r.Day1 = true
		return 1
	}

	if !r.Day2 {
		r.Day2 = true
		return 2
	}

	if !r.Day3 {
		r.Day3 = true
		return 3
	}

	if !r.Day4 {
		r.Day4 = true
		return 4
	}

	if !r.Day5 {
		r.Day5 = true
		return 5
	}
	return 5
}

func (r *SignInRecord) BuildMessage() *pb.MsgSignInRecord {
	msg := &pb.MsgSignInRecord{}

	msg.Day1 = proto.Bool(r.Day1)
	msg.Day2 = proto.Bool(r.Day2)
	msg.Day3 = proto.Bool(r.Day3)
	msg.Day4 = proto.Bool(r.Day4)
	msg.Day5 = proto.Bool(r.Day5)
	msg.LastSignInTime = proto.Int64(r.LastSignTime.Unix())

	glog.V(2).Info("===>签到记录msg:", msg)

	return msg
}
