package prize

import (
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"sync"
	"util"
)

type VipPrize struct {
	Level          int    `bson:"level"`
	IcoRes         string `bson:"icoRes"`
	PrizeGold      int    `bson:"prizeGold"`
	PrizeDiamond   int    `bson:"prizeDiamond"`
	PrizeExp       int    `bson:"prizeExp"`
	PrizeScore     int    `bson:"prizeScore"`
	PrizeItemType  int    `bson:"prizeItemType"`
	PrizeItemCount int    `bson:"prizeItemCount"`
}

const (
	vipPrizeC = "vip_prize"
)

func FindVipPrizes() ([]*VipPrize, error) {
	prizes := []*VipPrize{}
	err := util.WithGameCollection(vipPrizeC, func(c *mgo.Collection) error {
		return c.Find(nil).All(&prizes)
	})
	return prizes, err
}

type VipPrizeManager struct {
	sync.RWMutex
	prizes map[int]*VipPrize
}

var vipPrizeManager *VipPrizeManager

func init() {
	vipPrizeManager = &VipPrizeManager{}
	vipPrizeManager.prizes = make(map[int]*VipPrize)
}

func GetVipPrizeManager() *VipPrizeManager {
	return vipPrizeManager
}

func (m *VipPrizeManager) Init() {
	m.Lock()
	defer m.Unlock()

	prizes, err := FindVipPrizes()
	if err != nil {
		glog.Fatal(err)
	}

	for _, p := range prizes {
		m.prizes[p.Level] = p
	}
}

func (m *VipPrizeManager) GetVipPrize(vipLevel int) (VipPrize, bool) {
	m.Lock()
	defer m.Unlock()

	p := m.prizes[vipLevel]
	if p != nil {
		return *p, true
	}

	return VipPrize{}, false
}
