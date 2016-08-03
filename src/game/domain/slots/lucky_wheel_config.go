package slots

import (
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"math/rand"
	"sync"
	"time"
	"util"
)

type WheelPrize struct {
	WheelSerialId  int     `bson:"wheelSerialId"`
	Probability    float64 `bson:"probability"`
	PrizeName      string  `bson:"prizeName"`
	PrizeGold      int     `bson:"prizeGold"`
	PrizeDiamond   int     `bson:"prizeDiamond"`
	PrizeExp       int     `bson:"prizeExp"`
	PrizeScore     int     `bson:"prizeScore"`
	PrizeItemType  int     `bson:"prizeItemType"`
	PrizeItemCount int     `bson:"prizeItemCount"`
}

const (
	wheelPrizeC = "wheel_prize"
)

func FindWheelPrizes() ([]*WheelPrize, error) {
	prizes := []*WheelPrize{}

	err := util.WithGameCollection(wheelPrizeC, func(c *mgo.Collection) error {
		return c.Find(nil).All(&prizes)
	})
	return prizes, err
}

type LuckyWheelConfigManager struct {
	sync.RWMutex
	configs map[int]WheelPrize
}

var wheelConfigManager *LuckyWheelConfigManager

func init() {
	wheelConfigManager = &LuckyWheelConfigManager{}
	wheelConfigManager.configs = make(map[int]WheelPrize)
}

func GetWheelConfigManager() *LuckyWheelConfigManager {
	return wheelConfigManager
}

func (m *LuckyWheelConfigManager) Init() {
	rand.Seed(time.Now().UnixNano())

	prizes, err := FindWheelPrizes()
	if err != nil {
		glog.Fatal("加载轮盘配置失败err:", err)
		return
	}

	m.Lock()
	defer m.Unlock()

	for _, prize := range prizes {
		m.configs[prize.WheelSerialId] = *prize
	}

	for i := 1; i <= 8; i++ {
		if _, ok := m.configs[i]; !ok {
			glog.Fatal("轮盘配置错误:", i)
		}
	}
}

func (m *LuckyWheelConfigManager) GetWheelConfig(id int) (WheelPrize, bool) {
	m.Lock()
	defer m.Unlock()

	c, ok := m.configs[id]
	return c, ok
}

func (m *LuckyWheelConfigManager) RandomPrize() (WheelPrize, bool) {
	m.Lock()
	defer m.Unlock()

	r := rand.Float64()

	var offset float64 = 0

	for i := 1; i <= 8; i++ {
		if r >= offset && r < offset+m.configs[i].Probability {
			return m.configs[i], true
		}
		offset += m.configs[i].Probability
	}

	return WheelPrize{}, false
}
