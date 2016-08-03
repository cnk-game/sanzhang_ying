package game

import (
	domainUser "game/domain/user"
	"github.com/golang/glog"
	"pb"
	"sync"
	"time"
	"util"
)

type GameManager struct {
	sync.RWMutex
	items map[int]*GameItem // gameId==>Item
	chs   map[int]chan int
}

var (
	ShowUser = "QF10000228"
)

var gameManager *GameManager

func init() {
	gameManager = &GameManager{}
	gameManager.items = make(map[int]*GameItem)
	gameManager.chs = make(map[int]chan int)

	go gameManager.playerCountByGameType()
}

func GetGameManager() *GameManager {
	return gameManager
}

func (m *GameManager) playerCountByGameType() {
	for {
		time.Sleep(120 * time.Second)
		if len(gameManager.chs) <= 0 {
			continue
		}

		err := util.MongoLog_ClearPlayerCountByGType()
		if err != nil {
			glog.Info("MongoLog_ClearPlayerCountByGType error, err=", err)
		}
		for gid, ch := range m.chs {
			item := m.getGame(gid)
			if item == nil {
				continue
			}
			pcount := <-ch
			util.MongoLog_SetPlayerCountByGType(int(item.GameType), pcount)
		}
	}
}

func (m *GameManager) getGame(gameId int) *GameItem {
	m.RLock()
	defer m.RUnlock()

	return m.items[gameId]
}

func (m *GameManager) FindGame(gameId int) *GameItem {
	return m.getGame(gameId)
}

func (m *GameManager) EnterGame(gameType util.GameType, gameId int, user *pb.UserDef) bool {
	m.Lock()
	item := m.items[gameId]
	if item == nil {
		m.chs[gameId] = make(chan int)
		item = NewGameItem(gameType, gameId, m.chs[gameId])
		m.items[gameId] = item
	}
	m.Unlock()

	return item.enterGame(user)
}

func (m *GameManager) LeaveGame(gameId int, userId string, changeDesk bool) {
	item := m.getGame(gameId)
	if item == nil {
		return
	}

	item.LeaveGame(userId, changeDesk)	
}

//wjs 修复被踢人不离桌问题
func (m *GameManager) KickedLeaveGame(gameId int, userId string, changeDesk bool) {
	item := m.getGame(gameId)
	if item == nil {
		return
	}

	item.KickedLeaveGame(userId, changeDesk)	
}


// add by wangsq start --- 踢人
func (m *GameManager) KickPlayer(gameId int, from_userId string, target_userId string) bool {
	item := m.getGame(gameId)
	if item == nil {
		return false
	}

	return item.KickPlayer(from_userId, target_userId)
}

// add by wangsq end

func (m *GameManager) onOpCards(gameId int, userId string, msg *pb.MsgOpCardReq) {
	item := m.getGame(gameId)
	if item == nil {
		return
	}

	item.onOpCards(userId, msg)
}

func (m *GameManager) onChatMsg(gameId int, userId string, msg *pb.MsgChat) {
	item := m.getGame(gameId)
	if item == nil {
		return
	}

	item.onChatMsg(userId, msg)
}

func (m *GameManager) onConsumeProps(gameId int, userId string, itemType pb.MagicItemType, replaceCard int) {
	item := m.getGame(gameId)
	if item == nil {
		return
	}

	item.onConsumeProps(userId, itemType, replaceCard)
}

func (m *GameManager) onJoinWaitQueue(gameId int, userId string) {
	item := m.getGame(gameId)
	if item == nil {
		return
	}

	item.onJoinWaitQueue(userId)
}

func (m *GameManager) onLeaveWaitQueue(gameId int, userId string) {
	item := m.getGame(gameId)
	if item == nil {
		return
	}

	item.onLeaveWaitQueue(userId)
}

func (m *GameManager) onLookupBetGold(gameId int, userId string, betUserId string) {
	item := m.getGame(gameId)
	if item == nil {
		return
	}

	item.onLookupBetGold(userId, betUserId)
}

func (m *GameManager) onUpdateGoldInGame(gameId int, userId string) {
	item := m.getGame(gameId)
	if item == nil {
		return
	}

	item.onUpdateGoldInGame(userId)
}

func (m *GameManager) onRewardInGame(gameId int, userId string) bool {
	item := m.getGame(gameId)
	if item == nil {
		return false
	}

	return item.onRewardInGame(userId)
}

func (m *GameManager) onEnterBackground(gameId int, userId string) {
	item := m.getGame(gameId)
	glog.V(2).Info("====>onEnterBackground userId:", userId, " item:", item)
	if item != nil {
		item.onEnterBackground(userId)
	}
}

func (m *GameManager) onEnterForeground(gameId int, userId string) {
	item := m.getGame(gameId)
	glog.V(2).Info("====>onEnterForeground userId:", userId, " item:", item)
	if item != nil {
		item.onEnterForeground(userId)
	} else {
		// 不在游戏中
		domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_APP_ENTER_FOREGROUND), nil)
	}
}

func (m *GameManager) UpdateGoldInGame(gameId int, userId string) {
	item := m.getGame(gameId)
	if item == nil {
		return
	}
	item.onUpdateGoldInGame(userId)
}
