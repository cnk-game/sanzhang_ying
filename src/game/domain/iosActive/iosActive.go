package iosActive

import (
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"sync"
	"util"
)

type IosActive struct {
	Id      int    `bson:"id"`
	Remains int    `bson:"remains"`
	Total   int    `bson:"total"`
	Name    string `bson:"name"`
	Price   int    `bson:"price"`
	Desc    string `bson:"desc"`
	IconUrl string `bson:"iconUrl"`
}

const (
	ios_active_contentC = "active_ios_prize"
)

type IosActiveManager struct {
	sync.RWMutex
	actives map[int]*IosActive
}

var iosActiveManager *IosActiveManager

func init() {

	iosActiveManager = &IosActiveManager{}
	iosActiveManager.actives = make(map[int]*IosActive)
}

func GetIosActiveManager() *IosActiveManager {
	return iosActiveManager
}

func FindIosAcives() ([]*IosActive, error) {
	actives := []*IosActive{}

	err := util.WithGameCollection(ios_active_contentC, func(c *mgo.Collection) error {
		return c.Find(nil).All(&actives)
	})
	return actives, err
}

func (m *IosActiveManager) Init() {
	m.Lock()
	defer m.Unlock()

	actives, err := FindIosAcives()
	if err != nil {
		glog.Fatal(err)
	}

	for _, active := range actives {
		m.actives[active.Id] = active
	}
}

func (m *IosActiveManager) GetIosActives() map[int]*IosActive {
	return m.actives
}

func (m *IosActiveManager) GetItemInfo(id int) (IosActive, bool) {
	m.Lock()
	defer m.Unlock()
	item, ok := m.actives[id]
	return *item, ok
}

func (m *IosActiveManager) UpdateItemCount(id int) error {
	m.Lock()
	defer m.Unlock()
	item, ok := m.actives[id]
	if !ok {
		return nil
	} else {
		item.Remains -= 1
		return util.WithGameCollection(ios_active_contentC, func(c *mgo.Collection) error {
			_, err := c.Upsert(bson.M{"id": id}, &item)
			return err
		})
		return nil
	}
}
