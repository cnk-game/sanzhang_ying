package config

import (
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"sync"
	"util"
)

type Ios_Online_Config struct {
	VersionId string `bson:"version"`
	IsOpen    bool   `bson:"isOpen"`
}

const (
	ios_online_configC = "ios_online_config"
)

type IosOnlineConfigManager struct {
	sync.RWMutex
	iosOnlineConfig map[string]*Ios_Online_Config
}

var iosOnlineConfigManager *IosOnlineConfigManager

func init() {
	iosOnlineConfigManager = &IosOnlineConfigManager{}
	iosOnlineConfigManager.iosOnlineConfig = make(map[string]*Ios_Online_Config)
}

func GetIosOnlineConfigManager() *IosOnlineConfigManager {
	return iosOnlineConfigManager
}

func FindOnlineConfigs() ([]*Ios_Online_Config, error) {
	configs := []*Ios_Online_Config{}

	err := util.WithGameCollection(ios_online_configC, func(c *mgo.Collection) error {
		return c.Find(nil).All(&configs)
	})
	return configs, err
}

func (m *IosOnlineConfigManager) Init() {
	m.Lock()
	defer m.Unlock()

	configs, err := FindOnlineConfigs()
	if err != nil {
		glog.Fatal(err)
	}

	for _, config := range configs {
		m.iosOnlineConfig[config.VersionId] = config
	}
}

func (m *IosOnlineConfigManager) GetConfig(version string) (Ios_Online_Config, bool) {
	m.Lock()
	defer m.Unlock()

	item, ok := m.iosOnlineConfig[version]
	return *item, ok
}
