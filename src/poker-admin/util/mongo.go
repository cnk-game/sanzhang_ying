package util

import (
	"github.com/astaxie/beego"
	mgo "gopkg.in/mgo.v2"
)

var (
	mgoSession *mgo.Session
	logSession *mgo.Session
)

var WebsiteDBName string
var LogDbName string

var dbUrl string
var logDbUrl string
var dbPoolLimit int

func init() {
	dbUrl = beego.AppConfig.DefaultString("db_url", "192.168.1.128")
	WebsiteDBName = beego.AppConfig.DefaultString("db_name", "website_db")
	logDbUrl = beego.AppConfig.DefaultString("log_db_url", "192.168.1.128")
	LogDbName = beego.AppConfig.DefaultString("log_db_name", "poker_log")
	dbPoolLimit = beego.AppConfig.DefaultInt("db_pool_limit", 200)
}

func GetSession() *mgo.Session {
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

func GetLogSession() *mgo.Session {
	if logDbUrl == "" {
		return GetSession()
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
