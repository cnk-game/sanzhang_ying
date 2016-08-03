package slots

import (
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"fmt"
	"util"
	"sync"
)

const (
    slotPoolValueC = "slot_pool_value"
    slotName = "slotPool"
)

type PoolInfo struct {
	PoolValue   int     `bson:"poolValue"`
	PoolName    string  `bson:"slotPool"`
}

type SlotGlobal struct {
	sync.RWMutex
	CurrentDateKey     string
	CurrentPoolValue   int
}

var slotGlobal *SlotGlobal

func Init() {
	slotGlobal = &SlotGlobal{}
	tm := time.Now()
    slotGlobal.CurrentDateKey = fmt.Sprintf("%04d-%02d-%02d", tm.Year(), tm.Month(), tm.Day())
    value, err := slotGlobal.Mongo_GetPoolValue()
    if err != nil {
        slotGlobal.CurrentPoolValue = 0
    } else {
        slotGlobal.CurrentPoolValue = value
    }
    glog.Info("SlotGlobal init ok >>>>>>")
    glog.Info(slotGlobal.CurrentDateKey)
    glog.Info(slotGlobal.CurrentPoolValue)
    glog.Info("SlotGlobal init end <<<<<<")
}

func GetSlotGlobal() *SlotGlobal {
	return slotGlobal
}

func (s *SlotGlobal) InputPool(gold int) {
    s.Lock()
    defer s.Unlock()

    s.CurrentPoolValue += gold
    if s.CurrentPoolValue < 0 {
        s.CurrentPoolValue = 0
        s.Mongo_SetPoolValue(0)
    } else {
        s.Mongo_InputPoolValue(gold)
    }
}

func (s *SlotGlobal)GetPoolValue() int {
    s.Lock()
    defer s.Unlock()

    return s.CurrentPoolValue
}

func (s *SlotGlobal)Mongo_GetPoolValue() (int, error){
    s.Lock()
    defer s.Unlock()

    info := &PoolInfo{}
    err := util.WithGameCollection(slotPoolValueC, func(c *mgo.Collection) error {
        err := c.Find(bson.M{"name": slotName}).One(info)
        return err
    })

    return info.PoolValue, err
}

func (s *SlotGlobal)Mongo_SetPoolValue(gold int) error {
    err := util.WithGameCollection(slotPoolValueC, func(c *mgo.Collection) error {
        err := c.Update(bson.M{"name": slotName}, bson.M{"$set": bson.M{"poolValue": gold}})
        return err
    })

    return err
}

func (s *SlotGlobal)Mongo_InputPoolValue(gold int) error {
    go s.Mongo_CheckPrintLog()
    err := util.WithGameCollection(slotPoolValueC, func(c *mgo.Collection) error {
        _, err := c.Upsert(bson.M{"name": slotName}, bson.M{"$inc": bson.M{"poolValue": gold}})
        return err
    })

    return err
}

func (s *SlotGlobal)Mongo_CheckPrintLog() {
    tm := time.Now()
    date := fmt.Sprintf("%04d-%02d-%02d", tm.Year(), tm.Month(), tm.Day())
    if date != s.CurrentDateKey {
        util.MongoLog_SlotPool(s.CurrentPoolValue, s.CurrentDateKey)
        s.CurrentDateKey = date
    }
}
