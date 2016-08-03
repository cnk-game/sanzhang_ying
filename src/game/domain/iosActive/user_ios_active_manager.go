package iosActive

import (
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"sync"
	"time"
	"util"
)

type UserIosActiveManager struct {
	sync.RWMutex
	items         map[string]*UserIosActive
	activeContent *IosActiveContent
	isContinue    bool
	beginInt      int64
	endInt        int64
}

var userIosActiveManager *UserIosActiveManager

func init() {
	userIosActiveManager = &UserIosActiveManager{}
	userIosActiveManager.items = make(map[string]*UserIosActive)
}

func (m *UserIosActiveManager) Init() {
	m.Lock()
	defer m.Unlock()
	m.activeContent, _ = GetIosActiveContent()
	m.beginInt = int64(util.ParseTime(m.activeContent.BeginTime).Unix())
	m.endInt = int64(util.ParseTime(m.activeContent.EndTime).Unix())
	glog.Info("Init++GetIosActiveContent content ", m.beginInt)
	glog.Info("Init++GetIosActiveContent content ", m.endInt)
	cur := time.Now()
	if m.IsInActive(cur) {
		m.isContinue = true
	}
}

func (m *UserIosActiveManager) GetIosActiveContent() (string, string, string) {
	m.Lock()
	defer m.Unlock()

	if m.activeContent == nil {
		return "", "", ""
	} else {
		return m.activeContent.Content, m.activeContent.BeginTime, m.activeContent.EndTime
	}
}

func (m *UserIosActiveManager) IsInActive(timeIn time.Time) bool {
	curInt := int64(timeIn.Unix())
	if curInt > m.beginInt && curInt < m.endInt {
		return true
	} else {
		return false
	}
}

func (m *UserIosActiveManager) IsActiveContinue() bool {
	return m.IsInActive(time.Now())
	//return m.isContinue
}

func GetUserIosActiveManager() *UserIosActiveManager {
	return userIosActiveManager
}

func (m *UserIosActiveManager) GetActiveStatus(userId string) UserIosActive {
	m.Lock()
	defer m.Unlock()

	item := m.items[userId]
	if item != nil {
		return *item
	} else {
		info, err1 := FindUserIosActive(userId)
		if err1 == mgo.ErrNotFound {
			active := &UserIosActive{}
			active.UserId = userId
			active.Gold = int64(0)
			active.Time = time.Now()
			m.items[userId] = active
			return *active
		} else {
			m.items[userId] = info
			return *info
		}
	}

}

func (m *UserIosActiveManager) AddGold(userId string, gold int64) int64 {
	m.Lock()
	defer m.Unlock()

	item := m.items[userId]
	now := time.Now()
	if item == nil {
		info, err1 := FindUserIosActive(userId)
		if err1 == mgo.ErrNotFound {
			active := &UserIosActive{}
			active.UserId = userId
			active.Gold = gold
			active.Time = now
			m.items[userId] = active
			return active.Gold
		} else {
			m.items[userId] = info
			if m.IsInActive(info.Time) {
				info.Gold += gold

			} else {
				info.Gold = gold
			}

			info.Time = now
			return info.Gold
		}

	} else {
		if m.IsInActive(item.Time) {
			item.Gold += gold

		} else {
			item.Gold = gold
		}
		item.Time = now
		return item.Gold
	}
}

func (m *UserIosActiveManager) SaveStatus(userId string) bool {
	m.Lock()
	defer m.Unlock()

	item := m.items[userId]
	if item == nil {
		return false
	} else {
		gold := item.Gold
		SaveUserIosActive(userId, int64(gold))
		delete(m.items, userId)
		return true
	}
}
