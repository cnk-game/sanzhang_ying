package user

import (
	"github.com/golang/glog"
	"math/rand"
	"sync"
	"time"
)

type FakeRankingList struct {
	sync.RWMutex
	l map[string]int
}

var fakeRankingList *FakeRankingList

func init() {
	fakeRankingList = &FakeRankingList{}
	fakeRankingList.l = make(map[string]int)
	go func() {
		for {
			fakeRankingList.FakeRankingList()
			time.Sleep(time.Duration(600+rand.Int()%1200) * time.Second)
		}
	}()
}

func GetFakeRankingList() *FakeRankingList {
	return fakeRankingList
}

func (m *FakeRankingList) SetRobot(userId string) {
	m.Lock()
	defer m.Unlock()

	m.l[userId] = 0
	glog.V(2).Info("===>充值排行榜添加机器人userId:", userId)
}

func (m *FakeRankingList) RemoveRobot(userId string) {
	m.Lock()
	defer m.Unlock()

	delete(m.l, userId)
	glog.V(2).Info("===>充值排行榜移除机器人userId:", userId)
}

func (m *FakeRankingList) FakeRankingList() {
	m.Lock()
	defer m.Unlock()

	glog.V(2).Info("====>添加充值排行榜len:", len(m.l))
	if len(m.l) <= 0 {
		return
	}

	r := rand.Int() % len(m.l)

	index := 0
	for userId := range m.l {
		if index == r {
			amount := m.getRandDiamond()
			GetUserFortuneManager().UpdateRechargeRankingList(userId, amount)
			glog.V(2).Info("====>更新充值排行榜userId:", userId, " 10")
			return
		}
		index++
	}
}

func (m *FakeRankingList) getRandDiamond() int {
	r := rand.Int() % 100
	if r >= 0 && r < 70 {
		return []int{2, 10, 30, 50}[rand.Int()%4]
	} else if r >= 70 && r < 90 {
		return []int{98, 100}[rand.Int()%2]
	} else if r >= 90 && r < 95 {
		return 298
	} else if r >= 95 && r < 99 {
		return 500
	} else {
		return 698
	}
}
