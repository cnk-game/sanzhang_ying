package config

import (
	"sync"
	"util"
)

const (
	FishMD5Key = "QF_Fish"
	FIshFuncClose = false
)

type FishConfig struct {
	FishType    int
	FishName    string
	Price       int
	Harvest     int
	GrowUp      int
	Profit      int
	Url         string
}


type FishConfigManager struct {
	sync.RWMutex
	fishConfigs map[int]FishConfig
}

var fishConfigManager *FishConfigManager


func init() {
	fishConfigManager = &FishConfigManager{}
	fishConfigManager.fishConfigs = make(map[int]FishConfig)

	fish1 := FishConfig{}
	fish1.FishType = int(util.FishType_Fish_ID_1)
	fish1.FishName = string(util.FishType_Fish_Name_1)
	fish1.Price = int(util.FishType_Fish_Price_1)
	fish1.Harvest = int(util.FishType_Fish_Price_1)
	fish1.GrowUp = 30
	fish1.Profit = 140
	fish1.Url = "fish_1"
	fishConfigManager.fishConfigs[fish1.FishType] = fish1

	fish2 := FishConfig{}
    fish2.FishType = int(util.FishType_Fish_ID_2)
    fish2.FishName = string(util.FishType_Fish_Name_2)
    fish2.Price = int(util.FishType_Fish_Price_2)
    fish2.Harvest = int(util.FishType_Fish_Price_2)
    fish2.GrowUp = 30
    fish2.Profit = 130
    fish2.Url = "fish_2"
    fishConfigManager.fishConfigs[fish2.FishType] = fish2

    fish3 := FishConfig{}
    fish3.FishType = int(util.FishType_Fish_ID_3)
    fish3.FishName = string(util.FishType_Fish_Name_3)
    fish3.Price = int(util.FishType_Fish_Price_3)
    fish3.Harvest = int(util.FishType_Fish_Price_3)
    fish3.GrowUp = 30
    fish3.Profit = 65
    fish3.Url = "fish_3"
    fishConfigManager.fishConfigs[fish3.FishType] = fish3

    fish4 := FishConfig{}
    fish4.FishType = int(util.FishType_Fish_ID_4)
    fish4.FishName = string(util.FishType_Fish_Name_4)
    fish4.Price = int(util.FishType_Fish_Price_4)
    fish4.Harvest = int(util.FishType_Fish_Price_4)
    fish4.GrowUp = 30
    fish4.Profit = 13
    fish4.Url = "fish_4"
    fishConfigManager.fishConfigs[fish4.FishType] = fish4

}

func GetFishConfigManager() *FishConfigManager {
	return fishConfigManager
}

func (m *FishConfigManager) GetFishConfig(fishType int) (FishConfig, bool) {
	m.RLock()
	defer m.RUnlock()

	c, ok := m.fishConfigs[fishType]
	return c, ok
}

func (m *FishConfigManager) GetFishAll() ([]FishConfig) {
    result := []FishConfig{}
    for _, value := range m.fishConfigs {
        result = append(result, value)
    }

    return result
}