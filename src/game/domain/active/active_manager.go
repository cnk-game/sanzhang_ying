package active

import (
	mgo "gopkg.in/mgo.v2"
	"sync"
)

type ActiveInfoManager struct {
	sync.RWMutex
	items map[string]string
}

var activeInfoManager *ActiveInfoManager

func init() {
	activeInfoManager = &ActiveInfoManager{}
	activeInfoManager.items = make(map[string]string)
}

func GetActiveManager() *ActiveInfoManager {
	return activeInfoManager
}

func (m *ActiveInfoManager) GetActiveStatus(userId string) string {
	m.Lock()
	defer m.Unlock()

	id := ""
	id = m.items[userId]
	return id
}

func (m *ActiveInfoManager) AddItem(userId string, id string) bool {
	m.Lock()
	defer m.Unlock()

	m.items[userId] = id
	return true
}

func (m *ActiveInfoManager) GetStatus(userId string) string {
	id := m.GetActiveStatus(userId)
	if id == "" {
		info, err1 := FindUserActive(userId)
		if err1 == mgo.ErrNotFound {
			m.AddItem(info.UserId, "0")
			return "0"
		} else {
			m.AddItem(info.UserId, info.Id)
			return info.Id
		}
	} else {
		return id
	}

}

func (m *ActiveInfoManager) SaveStatus(userId string, id string) bool {
	SaveUserActive(userId, id)
	return true
}
