package prize

import (
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"sync"
	"util"
)

type OnlinePrize struct {
	PrizeID        int    `bson:"prizeID"`
	PrizeTitle     string `bson:"prizeTitle"`
	IcoRes         string `bson:"icoRes"`
	BeginTime      int    `bson:"beginTime"`
	EndTime        int    `bson:"endTime"`
	PrizeGold      int    `bson:"prizeGold"`
	PrizeDiamond   int    `bson:"prizeDiamond"`
	PrizeExp       int    `bson:"prizeExp"`
	PrizeScore     int    `bson:"prizeScore"`
	PrizeItemType  int    `bson:"prizeItemType"`
	PrizeItemCount int    `bson:"prizeItemCount"`
}

const (
	onlinePrizeC = "online_prize"
)

func FindOnlinePrizes() ([]*OnlinePrize, error) {
	hrs := []*OnlinePrize{}
	err := util.WithGameCollection(onlinePrizeC, func(c *mgo.Collection) error {
		return c.Find(nil).All(&hrs)
	})
	return hrs, err
}

type OnlinePrizeManager struct {
	sync.RWMutex
	prizes map[int]*OnlinePrize
}

var onlinePrizeManager *OnlinePrizeManager

func init() {
	onlinePrizeManager = &OnlinePrizeManager{}
	onlinePrizeManager.prizes = make(map[int]*OnlinePrize)
}

func GetOnlinePrizeManager() *OnlinePrizeManager {
	return onlinePrizeManager
}

func (m *OnlinePrizeManager) Init() {
	prizes, err := FindOnlinePrizes()
	if err != nil {
		glog.Fatal("加载在线奖励失败err:", err)
	}
	for _, p := range prizes {
		m.prizes[p.PrizeID] = p
	}
}

func (m *OnlinePrizeManager) GetOnlinePrize(id int) (OnlinePrize, bool) {
	m.Lock()
	defer m.Unlock()

	p := m.prizes[id]
	if p != nil {
		return *p, true
	}

	return OnlinePrize{}, false
}

func (m *OnlinePrizeManager) GetOnlineAllPrize() map[int]*OnlinePrize {
	return m.prizes
}