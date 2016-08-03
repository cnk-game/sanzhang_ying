package config

import (
	mgo "gopkg.in/mgo.v2"
	"sync"
	"github.com/golang/glog"
	"util"
)

type MatchConfig_Data struct {
	MatchType int `bson:"matchType"`	
	Name    string   `bson:"name"`
	ChipList string	`bson:"chipList"`
	DefaultBet	int	`bson:"defaultBet"`
	AutoQuitLimit int `bson:"autoQuitLimit"`
	EnterLimit int `bson:"enterLimit"`
	Reward int `bson:"reward"`	
	MaxLimit int `bson:"maxLimit"`	
	QuickLimit int `bson:"quickLimit"`	
		
}

const (
	Match_ConfigC = "match_config"
)


type MatchConfigManager struct {
	sync.RWMutex
	items map[string]*MatchConfig_Data
}

var matchConfigManager *MatchConfigManager

func init() {
	matchConfigManager = &MatchConfigManager{}
	matchConfigManager.items = make(map[string]*MatchConfig_Data)
	
}

func GetMatchConfigManager() *MatchConfigManager {
	return matchConfigManager
}

func FindMatchDatas() ([]*MatchConfig_Data, error) {
	datas := []*MatchConfig_Data{}

	err := util.WithGameCollection(Match_ConfigC, func(c *mgo.Collection) error {
		return c.Find(nil).All(&datas)
	})
	return datas, err
}


func (m *MatchConfigManager) Init() {
	m.Lock()
	defer m.Unlock()

	datas, err := FindMatchDatas()
	if err != nil {
		glog.Fatal(err)
	}

	for _, data := range datas {
		m.items[data.Name] = data		
	}
	
	glog.Info("match config data request")
}


func (m *MatchConfigManager) GetData() map[string]*MatchConfig_Data {
	m.Lock()
	defer m.Unlock()
	
	item := m.items
	return item
	

}

