package stats

import (
	"config"
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"sync"
	"time"
	"util"
)

const (
	WIN_TIME  = 0
	LOSE_TIME = 1
)

type AiFortuneLog struct {
	MatchType    int       `bson:"matchType"`
	EarnGold     int64     `bson:"earnGold"`
	SwitchStatus int       `bson:"switchStatus"`
	Time         time.Time `bson:"time"`
}

type SwitchTimeLog struct {
	MatchType    int       `bson:"matchType"`
	SwitchStatus string    `bson:"switchStatus"`
	SwitchGold   int64     `bson:"switchGold"`
	Time         time.Time `bson:"time"`
}

const (
	aiFortuneLogC  = "ai_fortune_log"
	switchWinTimeC = "switch_time_log"
)

func FindAiFortuneLogs() ([]*AiFortuneLog, error) {
	l := []*AiFortuneLog{}
	err := util.WithLogCollection(aiFortuneLogC, func(c *mgo.Collection) error {
		return c.Find(nil).All(&l)
	})
	return l, err
}

func SaveAiFortuneLog(l *AiFortuneLog) error {
	l.Time = time.Now()
	return util.WithLogCollection(aiFortuneLogC, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"matchType": l.MatchType}, l)
		return err
	})
}

func SaveSwitchTimeLog(l *SwitchTimeLog) error {
	l.Time = time.Now()
	return util.WithLogCollection(switchWinTimeC, func(c *mgo.Collection) error {
		return c.Insert(l)
	})
}

type AiFortuneLogManager struct {
	sync.RWMutex
	logs map[int]*AiFortuneLog
}

var aiFortuneLogManager *AiFortuneLogManager

func init() {
	aiFortuneLogManager = &AiFortuneLogManager{}
	aiFortuneLogManager.logs = make(map[int]*AiFortuneLog)
}

func GetAiFortuneLogManager() *AiFortuneLogManager {
	return aiFortuneLogManager
}

func (m *AiFortuneLogManager) Init() {
	logs, _ := FindAiFortuneLogs()
	for _, item := range logs {
		m.logs[item.MatchType] = item
	}
}

func (m *AiFortuneLogManager) AddEarnGold(matchType int, gold int64) {
	winGold, loseGold := config.GetCardConfigManager().GetTimeConfig(matchType)

	glog.Info("===>切换所需matchType:", matchType, " winGold:", winGold, " loseGold:", loseGold)

	m.Lock()
	defer m.Unlock()

	item := m.logs[matchType]
	if item == nil {
		item = &AiFortuneLog{}
		item.MatchType = matchType
		item.SwitchStatus = WIN_TIME
		m.logs[item.MatchType] = item
	}

	item.EarnGold += gold

	switchGold := item.EarnGold

	if item.SwitchStatus == WIN_TIME {
		// 检测是否切换到LOSE_TIME
		glog.Info("当前为吸金期gameType:", matchType, ",吸入:", item.EarnGold, " 切换状态所需:", winGold)
		if item.EarnGold >= int64(winGold) {
			item.EarnGold = 0
			item.SwitchStatus = LOSE_TIME

			switchLog := &SwitchTimeLog{}
			switchLog.MatchType = matchType
			switchLog.SwitchGold = switchGold
			switchLog.SwitchStatus = "吐金"
			glog.Info("切换财富周期状态吐金期", switchLog)
			SaveSwitchTimeLog(switchLog)
		}
	} else {
		// 检测是否切换到WIN_TIME
		glog.Info("当前为吐金期gameType:", matchType, ",吐出:", -item.EarnGold, " 切换状态所需:", loseGold)
		if -item.EarnGold >= int64(loseGold) {
			item.EarnGold = 0
			item.SwitchStatus = WIN_TIME

			switchLog := &SwitchTimeLog{}
			switchLog.MatchType = matchType
			switchLog.SwitchGold = switchGold
			switchLog.SwitchStatus = "吸金"
			glog.Info("切换财富周期状态吸金期", switchLog)
			SaveSwitchTimeLog(switchLog)
		}
	}
}

func (m *AiFortuneLogManager) GetWinRate(matchType int) int {
	winRateHigh, winRateLow := config.GetCardConfigManager().GetWinRate(matchType)

	m.RLock()
	defer m.RUnlock()

	item := m.logs[matchType]
	if item == nil {
		glog.Info("===>找不到财富周期配置，默认50")
		return 50
	}

	if item.SwitchStatus == WIN_TIME {
		glog.Info("====>财富周期胜率matchType:", matchType, " 吸金期:", winRateHigh)
		return winRateHigh
	}

	glog.Info("====>财富周期胜率matchType:", matchType, " 吐金期:", winRateLow)

	return winRateLow
}

func (m *AiFortuneLogManager) SaveLog() {
	m.Lock()
	defer m.Unlock()

	for _, item := range m.logs {
		SaveAiFortuneLog(item)
	}
}
