package stats

import (
	"fmt"
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"sync"
	"time"
	"util"
)

type MatchLog struct {
	Year             int   `bson:"year"`
	Month            int   `bson:"month"`
	Day              int   `bson:"day"`
	MatchType        int   `bson:"matchType"`
	WinGold          int64 `bson:"winGold"`
	WinGoldWithRobot int64 `bson:"winGoldWithRobot"`
	RoundCount       int   `bson:"roundCount"`
}

const (
	matchLogC = "match_log"
)

func FindMatchLogs() ([]*MatchLog, error) {
	logs := []*MatchLog{}
	now := time.Now()
	err := util.WithLogCollection(matchLogC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"year": now.Year(), "month": int(now.Month()), "day": now.Day()}).All(&logs)
	})
	return logs, err
}

func SaveMatchLog(l *MatchLog) error {
	return util.WithLogCollection(matchLogC, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"year": l.Year, "month": l.Month, "day": l.Day, "matchType": l.MatchType}, l)
		return err
	})
}

type MatchLogManager struct {
	sync.RWMutex
	logs           map[string]*MatchLog
	lastSaveDbTime time.Time
}

var matchLogManager *MatchLogManager

func init() {
	matchLogManager = &MatchLogManager{}
	matchLogManager.logs = make(map[string]*MatchLog)
}

func GetMatchLogManager() *MatchLogManager {
	return matchLogManager
}

func (m *MatchLogManager) Init() bool {
	logs, err := FindMatchLogs()
	if err != nil {
		glog.Error("加载比赛记录失败!")
		return false
	}

	m.Lock()

	for _, l := range logs {
		m.logs[fmt.Sprintf("%v-%v-%v-%v", l.Year, l.Month, l.Day, l.MatchType)] = l
	}
	m.lastSaveDbTime = time.Now()

	m.Unlock()

	go m.doSave()

	return true
}

func (m *MatchLogManager) AddMatchLog(matchType int, winGold int, winGoldWithRobot int) {
	m.Lock()
	defer m.Unlock()

	now := time.Now()

	k := fmt.Sprintf("%v-%v-%v-%v", now.Year(), int(now.Month()), now.Day(), matchType)
	l := m.logs[k]
	if l == nil {
		l = &MatchLog{}
		l.Year = now.Year()
		l.Month = int(now.Month())
		l.Day = now.Day()
		l.MatchType = matchType
		m.logs[k] = l
	}
	l.WinGold += int64(winGold)
	l.WinGoldWithRobot += int64(winGoldWithRobot)
	l.RoundCount++
}

func (m *MatchLogManager) doSave() {
	for {
		if !util.CompareDate(time.Now(), m.lastSaveDbTime) {
			// 保存
			m.SaveMatchLogs()
			m.lastSaveDbTime = time.Now()
		}
		time.Sleep(time.Minute)
	}
}

func (m *MatchLogManager) SaveMatchLogs() {
	m.Lock()
	defer m.Unlock()

	t := time.Now().AddDate(0, 0, -1)
	date := fmt.Sprintf("%v-%v-%v", t.Year(), int(t.Month()), t.Day())
	for k, l := range m.logs {
		if strings.HasPrefix(k, date) {
			SaveMatchLog(l)
		}
	}
}

func (m *MatchLogManager) SaveAllMatchLogs() {
	m.Lock()
	defer m.Unlock()

	for _, l := range m.logs {
		SaveMatchLog(l)
	}
}
