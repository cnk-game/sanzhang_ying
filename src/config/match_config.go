package config

import (
	"sync"
	"util"
)

type MatchConfig struct {
	MatchType    int
	SingleBet    int
	ChipList     []int
	KickoutLimit int
	EnterLimit   int
	TipChip      int
}

type SNGMatchConfig struct {
	MatchType     int
	RewardGold    int
	RewardScore   int
	EnterCostGold int
	SingleBet     int // 底注
	AnteIncGold   int // 每场底注增加
	ChipList      []int
}

type WanRenMatchConfig struct {
	MatchType          int
	SingleBet          int
	JoinWaitQueueLimit int
	KickoutLimit       int
	ChipList           []int
}

type MatchConfigManager struct {
	sync.RWMutex
	matchConfigs map[int]MatchConfig
	sngConfigs   map[int]SNGMatchConfig
	wanRenConfig WanRenMatchConfig
}

var matchConfigManager *MatchConfigManager

func init() {
	matchConfigManager = &MatchConfigManager{}
	matchConfigManager.matchConfigs = make(map[int]MatchConfig)
	matchConfigManager.sngConfigs = make(map[int]SNGMatchConfig)

	sng1 := SNGMatchConfig{}
	sng1.MatchType = int(util.GameType_SNG_Level_1)
	sng1.RewardGold = 100000
	sng1.RewardScore = 1
	sng1.EnterCostGold = 20000
	sng1.SingleBet = 100
	sng1.AnteIncGold = 1000
	sng1.ChipList = []int{100, 300, 500, 800, 1000}
	matchConfigManager.sngConfigs[sng1.MatchType] = sng1

	sng2 := SNGMatchConfig{}
	sng2.MatchType = int(util.GameType_SNG_Level_2)
	sng2.RewardGold = 1000000
	sng2.RewardScore = 30
	sng2.EnterCostGold = 500000
	sng2.SingleBet = 5000
	sng2.AnteIncGold = 2000
	sng2.ChipList = []int{5000, 10000, 15000, 20000, 30000}
	matchConfigManager.sngConfigs[sng2.MatchType] = sng2

	sng3 := SNGMatchConfig{}
	sng3.MatchType = int(util.GameType_SNG_Level_3)
	sng3.RewardGold = 0
	sng3.RewardScore = 100
	sng3.EnterCostGold = 250000
	sng3.SingleBet = 5000
	sng3.AnteIncGold = 2000
	sng3.ChipList = []int{5000, 10000, 15000, 20000, 30000}
	matchConfigManager.sngConfigs[sng3.MatchType] = sng3

	//modify by yelong
	c1 := MatchConfig{}
	c1.MatchType = int(util.GameType_Common_Level_1)
	c1.SingleBet = 100
	c1.ChipList = []int{100, 300, 500, 800, 1000}
	c1.KickoutLimit = 100
	c1.EnterLimit = 1000
	c1.TipChip = 100
	matchConfigManager.matchConfigs[c1.MatchType] = c1
	//
	c2 := MatchConfig{}
	c2.MatchType = int(util.GameType_Common_Level_2)
	c2.SingleBet = 2000
	c2.ChipList = []int{2000, 5000, 8000, 10000, 15000}
	c2.KickoutLimit = 2000
	c2.EnterLimit = 100000
	c2.TipChip = 100
	matchConfigManager.matchConfigs[c2.MatchType] = c2

	c3 := MatchConfig{}
	c3.MatchType = int(util.GameType_Common_Level_3)
	c3.SingleBet = 20000
	c3.ChipList = []int{20000, 50000, 80000, 100000, 150000}
	c3.KickoutLimit = 20000
	c3.EnterLimit = 1000000
	c3.TipChip = 100
	matchConfigManager.matchConfigs[c3.MatchType] = c3
	//中级场添加 add by yelong
	c7 := MatchConfig{}
	c7.MatchType = int(util.GameType_Common_Level_4)
	c7.SingleBet = 1000
	c7.ChipList = []int{1000, 2000, 3000, 5000, 8000}
	c7.KickoutLimit = 1000
	c7.EnterLimit = 50000
	c7.TipChip = 100
	matchConfigManager.matchConfigs[c7.MatchType] = c7
	//
	c4 := MatchConfig{}
	c4.MatchType = int(util.GameType_Props_Level_1)
	c4.SingleBet = 100
	c4.ChipList = []int{100, 300, 500, 800, 1000}
	c4.KickoutLimit = 100
	c4.EnterLimit = 1000
	c4.TipChip = 100
	matchConfigManager.matchConfigs[c4.MatchType] = c4

	c5 := MatchConfig{}
	c5.MatchType = int(util.GameType_Props_Level_2)
	c5.SingleBet = 5000
	c5.ChipList = []int{5000, 10000, 15000, 20000, 30000}
	c5.KickoutLimit = 50000
	c5.EnterLimit = 300000
	c5.TipChip = 100
	matchConfigManager.matchConfigs[c5.MatchType] = c5

	c6 := MatchConfig{}
	c6.MatchType = int(util.GameType_Props_Level_3)
	c6.SingleBet = 50000
	c6.ChipList = []int{50000, 100000, 150000, 200000, 300000}
	c6.KickoutLimit = 500000
	c6.EnterLimit = 2000000
	c6.TipChip = 100
	matchConfigManager.matchConfigs[c6.MatchType] = c6

	matchConfigManager.wanRenConfig = WanRenMatchConfig{}
	matchConfigManager.wanRenConfig.SingleBet = 50000
	matchConfigManager.wanRenConfig.JoinWaitQueueLimit = 5000000
	matchConfigManager.wanRenConfig.KickoutLimit = 2000000
	matchConfigManager.wanRenConfig.ChipList = []int{50000, 150000, 250000, 400000, 500000}
}

func GetMatchConfigManager() *MatchConfigManager {
	return matchConfigManager
}

func (m *MatchConfigManager) GetMatchConfig(gameType int) (MatchConfig, bool) {
	m.RLock()
	defer m.RUnlock()

	c, ok := m.matchConfigs[gameType]
	return c, ok
}

func (m *MatchConfigManager) GetSNGMatchConfig(gameType int) (SNGMatchConfig, bool) {
	m.RLock()
	defer m.RUnlock()

	c, ok := m.sngConfigs[gameType]
	return c, ok
}

func (m *MatchConfigManager) GetWanRenConfig() WanRenMatchConfig {
	return m.wanRenConfig
}

func (m *MatchConfigManager) IsRaiseBetRight(gameType int, raiseBet int) bool {
	m.RLock()
	defer m.RUnlock()

	if gameType == int(util.GameType_WAN_REN) {
		for _, c := range m.wanRenConfig.ChipList {
			if c == raiseBet {
				return true
			}
		}
		return false
	}

	c1, ok := m.matchConfigs[gameType]
	if ok {
		for _, c := range c1.ChipList {
			if c == raiseBet {
				return true
			}
		}
		return false
	}

	c2, ok := m.sngConfigs[gameType]
	if ok {
		for _, c := range c2.ChipList {
			if c == raiseBet {
				return true
			}
		}
		return false
	}

	return false
}

func (m *MatchConfigManager) GetMaxChip(gameType int) (int, bool) {
	m.RLock()
	defer m.RUnlock()

	if gameType == int(util.GameType_WAN_REN) {
		max := 0
		for _, c := range m.wanRenConfig.ChipList {
			if c > max {
				max = c
			}
		}
		return max, true
	}

	c1, ok := m.matchConfigs[gameType]
	if ok {
		max := 0
		for _, c := range c1.ChipList {
			if c > max {
				max = c
			}
		}
		return max, true
	}

	c2, ok := m.sngConfigs[gameType]
	if ok {
		max := 0
		for _, c := range c2.ChipList {
			if c > max {
				max = c
			}
		}
		return max, true
	}

	return 0, false
}

func (m *MatchConfigManager) GetTipChip(gameType int) int {
	m.RLock()
	defer m.RUnlock()

	c1, ok := m.matchConfigs[gameType]
	if ok {
		return c1.TipChip
	}

	return 100
}

func (m *MatchConfigManager) GetMinChip(gameType int) (int, bool) {
	m.RLock()
	defer m.RUnlock()

	if gameType == int(util.GameType_WAN_REN) {
		return m.wanRenConfig.ChipList[0], true
	}

	c1, ok := m.matchConfigs[gameType]
	if ok {
		return c1.ChipList[0], true
	}

	c2, ok := m.sngConfigs[gameType]
	if ok {
		return c2.ChipList[0], true
	}

	return 0, false
}

func (m *MatchConfigManager) GetEnterLimit(gameType int) (int, bool) {
	m.RLock()
	defer m.RUnlock()

	if gameType == int(util.GameType_WAN_REN) {
		return 0, true
	}

	c1, ok := m.matchConfigs[gameType]
	if ok {
		return c1.EnterLimit, true
	}

	c2, ok := m.sngConfigs[gameType]
	if ok {
		return c2.EnterCostGold, true
	}

	return 0, false
}

func (m *MatchConfigManager) GetNextLevelEnterLimit(gameType int) (int, bool) {
	if gameType == int(util.GameType_Common_Level_1) {
		return m.GetEnterLimit(int(util.GameType_Common_Level_4))
	} else if gameType == int(util.GameType_Common_Level_4) {
		return m.GetEnterLimit(int(util.GameType_Common_Level_2))
	} else if gameType == int(util.GameType_Common_Level_2) {
		return m.GetEnterLimit(int(util.GameType_Common_Level_3))
	} else if gameType == int(util.GameType_Common_Level_3) {
		return 0, false
	} else if gameType == int(util.GameType_Props_Level_1) {
		return m.GetEnterLimit(int(util.GameType_Props_Level_2))
	} else if gameType == int(util.GameType_Props_Level_2) {
		return m.GetEnterLimit(int(util.GameType_Props_Level_3))
	} else if gameType == int(util.GameType_Props_Level_3) {
		return 0, false
	} else if gameType == int(util.GameType_SNG_Level_1) {
		return m.GetEnterLimit(int(util.GameType_SNG_Level_2))
	} else if gameType == int(util.GameType_SNG_Level_2) {
		return m.GetEnterLimit(int(util.GameType_SNG_Level_3))
	} else if gameType == int(util.GameType_SNG_Level_3) {
		return 0, false
	}

	return 0, false
}
