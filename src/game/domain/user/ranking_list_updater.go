package user

import (
	"sync"
	"time"
)

type RankingListUpdater struct {
	sync.RWMutex
	userIds map[string]int
}

var rankingListUpdater *RankingListUpdater

func init() {
	rankingListUpdater = &RankingListUpdater{}
	rankingListUpdater.userIds = make(map[string]int)

	go rankingListUpdater.update()
}

func GetRankingListUpdater() *RankingListUpdater {
	return rankingListUpdater
}

func (m *RankingListUpdater) UpdateUser(userId string, gold int) {
	m.Lock()
	defer m.Unlock()

	m.userIds[userId] += gold
}

func (m *RankingListUpdater) update() {
	for {
		m.DoUpdate()
		time.Sleep(time.Minute)
	}
}

func (m *RankingListUpdater) DoUpdate() {
	userIds := make(map[string]int)

	m.Lock()
	for userId, gold := range m.userIds {
		userIds[userId] = gold
	}
	m.userIds = make(map[string]int)
	m.Unlock()

	for userId, gold := range userIds {
		GetUserFortuneManager().UpdateEarningsRankingList(userId, gold)
	}
}
