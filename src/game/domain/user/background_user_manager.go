package user

import (
	"pb"
	"sync"
)

type BackgroundUserManager struct {
	sync.RWMutex
	userIds map[string]bool
}

var backgroundUserManager *BackgroundUserManager

func init() {
	backgroundUserManager = &BackgroundUserManager{}
	backgroundUserManager.userIds = make(map[string]bool)
}

func GetBackgroundUserManager() *BackgroundUserManager {
	return backgroundUserManager
}

func (m *BackgroundUserManager) SetUser(userId string, isRobot bool) {
	m.Lock()
	defer m.Unlock()

	m.userIds[userId] = isRobot
}

func (m *BackgroundUserManager) DelUser(userId string) {
	m.Lock()
	defer m.Unlock()

	delete(m.userIds, userId)
}

func (m *BackgroundUserManager) IsHaveUser(userId string) bool {
	m.RLock()
	defer m.RUnlock()

	if _, ok := m.userIds[userId]; ok {
		return true
	}
	return false
}

func (m *BackgroundUserManager) Filter(userId string, msgId int32) bool {
	m.RLock()
	defer m.RUnlock()

	// 玩家不在后台，不过滤，继续发送消息
	if _, ok := m.userIds[userId]; !ok {
		return false
	}

	// 玩家在后台且为非机器人,阻止消息继续发送
	if !m.userIds[userId] {
		if msgId == int32(pb.MessageId_BUY_DAILY_GIFT_BAG_OK) {
			return false
		} else if msgId == int32(pb.MessageId_GET_PRIZE_MAIL_LIST) {
			return false
		} else if msgId == int32(pb.MessageId_UPDATE_VIP_TASK) {
		    return false
		}
		return true
	}

	// 为机器人，过滤聊天消息
	return msgId == int32(pb.MessageId_CHAT)
}
