package stats

import (
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"time"
	"util"
)

type OnlineLog struct {
	OnlinePlayerCount int       `bson:"onlinePlayerCount"`
	Datetime          time.Time `bson:"datetime"`
}

const (
	onlineLogC = "online_log"
)

func SaveOnlineLog(count int) error {
	l := &OnlineLog{}
	l.OnlinePlayerCount = count
	l.Datetime = time.Now()
	glog.V(2).Info("===>当前在线:", count)
	return util.WithLogCollection(onlineLogC, func(c *mgo.Collection) error {
		return c.Insert(l)
	})
}
