package config

import (
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"sync"
	"util"
)

const (
	vipPriceConfigC = "vip_price"
)

type VipPriceConfig struct {
	Name           string  `bson:"name"`
	PrizeGold      int     `bson:"prizeGold"`
	PrizeDays      int     `bson:"prizeDays"`
	Level          int     `bson:"level"`
	PrizeCharm     int     `bson:"prizeCharm"`
	PrizeGoldNow   int     `bson:"prizeGoldNow"`
}

type VipPriceConfigManager struct {
	sync.RWMutex
	configs       map[int]*VipPriceConfig
}

var vipPriceConfigM *VipPriceConfigManager

func GetVipPriceConfigManager() *VipPriceConfigManager {
	return vipPriceConfigM
}

func FindVipPriceConfigs() ([]*VipPriceConfig, error) {
	configs := []*VipPriceConfig{}
	err := util.WithGameCollection(vipPriceConfigC, func(c *mgo.Collection) error {
		return c.Find(nil).All(&configs)
	})

	return configs, err
}

func (m *VipPriceConfigManager) GetVipConfig() map[int]*VipPriceConfig {
    return m.configs
}

func (m *VipPriceConfigManager) Init() bool {
    glog.Info("==>VipPriceConfigManager.Init in.")
	m.Lock()
	defer m.Unlock()

	configs, err := FindVipPriceConfigs()
	if err != nil && err != mgo.ErrNotFound {
		glog.Error(err)
		return false
	}

	for _, config := range configs {
		glog.Info("==>vipPrice_config:", config)
		m.configs[config.Level] = config
	}

	return true
}

func init() {
    vipPriceConfigM = &VipPriceConfigManager{}
    vipPriceConfigM.configs = make(map[int]*VipPriceConfig)
}
