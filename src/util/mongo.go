package util

import (
	"flag"
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var (
	mgoSession *mgo.Session
	logSession *mgo.Session
	sequenceC  = "sequence"
)

var dbUrl string
var gameDb string
var userDb string
var logDb string
var logDbUrl string
var dbPoolLimit int

func init() {
	flag.StringVar(&dbUrl, "db", "127.0.0.1", "db url.")
	flag.StringVar(&logDbUrl, "logDBUrl", "", "logdb url.")
	flag.StringVar(&gameDb, "gameDB", "poker_game", "gameDB名称")
	flag.StringVar(&userDb, "userDB", "poker_user", "userDB名称")
	flag.StringVar(&logDb, "logDB", "poker_log", "logDB名称")
	flag.IntVar(&dbPoolLimit, "dbPoolLimit", 200, "数据库连接数限制")
}

func DumpDBInfo() string {
	return "dbUrl:" + dbUrl + " gameDb:" + gameDb + " userDb:" + userDb + " logDb:" + logDb
}

func getSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(dbUrl)
		if err != nil {
			panic(err)
		}
		mgoSession.SetPoolLimit(dbPoolLimit)
	}
	return mgoSession.Clone()
}

func getLogSession() *mgo.Session {
	if logDbUrl == "" {
		return getSession()
	}

	if logSession == nil {
		var err error
		logSession, err = mgo.Dial(logDbUrl)
		if err != nil {
			panic(err)
		}
		logSession.SetPoolLimit(dbPoolLimit)
	}
	return logSession.Clone()
}

func WithGameCollection(collection string, s func(*mgo.Collection) error) error {
	session := getSession()
	defer session.Close()
	c := session.DB(gameDb).C(collection)
	return s(c)
}

func WithUserCollection(collection string, s func(*mgo.Collection) error) error {
	t := time.Now()
	var t1 time.Time
	defer func() {
		if time.Since(t).Seconds() > 0.05 {
			glog.V(2).Info("==>collection:", collection, " total:", time.Since(t), " query:", time.Since(t1), " caller:", GetFunCaller(3))
		}
	}()
	session := getSession()
	defer session.Close()
	t1 = time.Now()
	c := session.DB(userDb).C(collection)
	return s(c)
}

func WithSafeUserCollection(collection string, s func(*mgo.Collection) error) error {
	t := time.Now()
	var t1 time.Time
	defer func() {
		if time.Since(t).Seconds() > 0.005 {
			glog.V(2).Info("==>collection:", collection, " total:", time.Since(t), " query:", time.Since(t1), " caller:", GetFunCaller(3))
		}
	}()
	session := getSession()
	defer session.Close()
	t1 = time.Now()
	session.SetSafe(&mgo.Safe{})
	c := session.DB(userDb).C(collection)
	return s(c)
}

func WithLogCollection(collection string, s func(*mgo.Collection) error) error {
	session := getLogSession()
	defer session.Close()
	c := session.DB(logDb).C(collection)
	return s(c)
}

func InitSequenceValue(name string, v int32) error {
	s := &sequence{}
	s.Name = name
	s.Val = int64(v)
	return WithSafeUserCollection(sequenceC, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"name": name}, s)
		return err
	})
}

type sequence struct {
	Name string
	Val  int64
}

func NextSequenceValue(name string) int32 {
	ch := mgo.Change{
		Update:    bson.M{"$inc": bson.M{"val": 1}},
		Upsert:    true,
		ReturnNew: true,
	}

	next := &sequence{}
	err := WithSafeUserCollection(sequenceC, func(c *mgo.Collection) error {
		_, err := c.Find(bson.M{"name": name}).Apply(ch, &next)
		return err
	})

	if err != nil {
		return -1
	}

	return int32(next.Val % 2147483648)
}
