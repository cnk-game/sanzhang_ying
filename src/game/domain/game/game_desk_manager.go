package game

import (
	domainUser "game/domain/user"
	"github.com/golang/glog"
	"math/rand"
	"pb"
	"sync"
	"time"
	"util"
)

type GameDeskUser struct {
	UserId  string
	IsRobot bool
}

type GameDesk struct {
	gameType util.GameType
	gameId   int
	players  map[string]*GameDeskUser
	Locked   bool
}

func (desk *GameDesk) hasPlayer(userId string) bool {
	_, ok := desk.players[userId]
	return ok
}

func (desk *GameDesk) allowRobotIn() bool {
	robotCount := 0
	playerCount := 0
	for _, p := range desk.players {
		if p.IsRobot {
			robotCount++
		} else {
			playerCount++
		}
	}

	// 桌内没有真实玩家，不允许机器人进入
	if playerCount <= 0 {
		return false
	}

	if util.IsGameTypeSNG(desk.gameType) {
		return true
	}

	// 桌内有一名以下机器人时，允许机器人进入（即一桌最多允许进入2个机器人）
	return robotCount <= 2
}

type GameDeskManager struct {
	sync.RWMutex
	wanRenDesk *GameDesk
	desks      map[int]*GameDesk // gameId==>gameDesk
	players    map[string]int    // userId==>gameId
	gameId     int
}

var deskManager *GameDeskManager

func init() {
	rand.Seed(time.Now().Unix())
	deskManager = &GameDeskManager{}
	deskManager.desks = make(map[int]*GameDesk)
	deskManager.players = make(map[string]int)
}

func GetDeskManager() *GameDeskManager {
	return deskManager
}

func (m *GameDeskManager) genGameId() int {
	m.gameId++
	return m.gameId
}

func (m *GameDeskManager) selectGameDesk(gameType util.GameType, ignoreGameId int, isRobot bool,userid string) *GameDesk {
	if util.IsGameTypeWanRen(gameType) {
		return m.wanRenDesk
	}

	if len(m.desks) <= 0 {
		return nil
	}

	gameIds := []int{}

	for _, item := range m.desks {
		if item.gameType != gameType {
			continue
		}

		if item.Locked {
			continue
		}

		if item.gameId == ignoreGameId {
			continue
		}

		if isRobot && !item.allowRobotIn() {
			continue
		}
		
		//wjs 游戏内被踢2分钟内不能进入此桌
		if(GetKickedUserManager().KickedLog[userid]!=nil && item.gameId == GetKickedUserManager().KickedLog[userid].GameId && (time.Now().Unix()-GetKickedUserManager().KickedLog[userid].KickTime) < 120){
		
			glog.Info("select target_userId=",userid)
			glog.Info("select gameId=",GetKickedUserManager().KickedLog[userid].GameId)	
			glog.Info("select time.Now().Unix()=",time.Now().Unix())	
			continue
		}		
		//wjs --end

		count := len(item.players)
		if count < util.MaxPlayerCount && count >= 1 {
			gameIds = append(gameIds, item.gameId)
		}
		if len(gameIds) >= 20 {
			break
		}
	}

	// 当前所有游戏都已满员
	count := len(gameIds)

	if count <= 0 {
		return nil
	}

	return m.desks[gameIds[rand.Int()%count]]
}

// add by wangsq start --- 踢人
func (m *GameDeskManager) KickPlayer(from_userId string, target_userId string) {
    m.Lock()
    defer m.Unlock()

    // 参数判断
    if m.players[from_userId] != m.players[target_userId] {
        glog.Infof("KickPlayer error. players not in same game desk.")
        return
    }

    gameId := m.players[from_userId]
    item := m.desks[gameId]
    ok := GetGameManager().KickPlayer(gameId, from_userId, target_userId)

    if ok && item != nil {
        delete(item.players, target_userId)
        delete(m.players, target_userId)
       // GetGameManager().LeaveGame(gameId, target_userId, false)
		GetGameManager().KickedLeaveGame(gameId, target_userId, false)
    }

//	//wjs 记录被踢时的数据
	GetKickedUserManager().KickedLog[target_userId] = &KickedUser{gameId, time.Now().Unix()}
	glog.Info("kicked target_userId=",target_userId)
	glog.Info("kicked gameId=%d",gameId)	
	glog.Info("kicked time.Now().Unix()=%d",time.Now().Unix())	
		

    return
}
// add by wangsq end


func (m *GameDeskManager) EnterGame(gameType util.GameType, user *pb.UserDef, ignoreGameId int) (int, bool) {
	m.Lock()

	desk := m.selectGameDesk(gameType, ignoreGameId, user.GetIsRobot(),*user.UserId)
	if desk == nil {
		if user.GetIsRobot() {
			m.Unlock()
			return 0, false
		}
		desk = &GameDesk{}
		desk.gameType = gameType
		desk.gameId = m.genGameId()
		desk.players = make(map[string]*GameDeskUser)
		if util.IsGameTypeWanRen(gameType) {
			m.wanRenDesk = desk
		} else {
			m.desks[desk.gameId] = desk
		}
	}

	m.Unlock()

	if !GetGameManager().EnterGame(gameType, desk.gameId, user) {
		return 0, false
	}

	m.Lock()
	defer m.Unlock()

	p := &GameDeskUser{}
	p.UserId = user.GetUserId()
	p.IsRobot = user.GetIsRobot()
	desk.players[user.GetUserId()] = p

	if !util.IsGameTypeWanRen(gameType) {
		m.players[user.GetUserId()] = desk.gameId
	}

	return desk.gameId, true
}

func (m *GameDeskManager) LeaveGame(userId string, gameType util.GameType, changeDesk bool) {
	glog.V(2).Info("===>LeaveGame userId:", userId, " type:", gameType, " changeDesk:", changeDesk)

	if changeDesk {
		// 换桌，非万人场
		m.Lock()
		gameId := m.players[userId]
		item := m.desks[gameId]
		if item != nil {
			delete(item.players, userId)
			delete(m.players, userId)
			m.Unlock()
			GetGameManager().LeaveGame(gameId, userId, changeDesk)
		} else {
			m.Unlock()
		}
		return
	}

	if util.IsGameTypeWanRen(gameType) {
		if m.wanRenDesk != nil {
			m.Lock()
			delete(m.wanRenDesk.players, userId)
			m.Unlock()
			GetGameManager().LeaveGame(m.wanRenDesk.gameId, userId, false)
		}
		return
	}

	m.Lock()
	gameId := m.players[userId]
	item := m.desks[gameId]

	glog.V(2).Info("===>userId:", userId, " gameId:", gameId, " item:", item)
	if item != nil {
		delete(item.players, userId)
		delete(m.players, userId)
		m.Unlock()
		GetGameManager().LeaveGame(gameId, userId, changeDesk)
	} else {
		m.Unlock()
	}
}

func (m *GameDeskManager) LockGame(gameId int, lock bool) {
	m.Lock()
	defer m.Unlock()

	item := m.desks[gameId]
	if item == nil {
		return
	}
	item.Locked = lock
}

func (m *GameDeskManager) getWanRenInfo(userId string) (gameId int, ok bool) {
	m.RLock()
	defer m.RUnlock()

	if m.wanRenDesk != nil && m.wanRenDesk.hasPlayer(userId) {
		return m.wanRenDesk.gameId, true
	}

	return 0, false
}

func (m *GameDeskManager) getGameInfo(userId string) (gameId int, ok bool) {
	m.RLock()
	defer m.RUnlock()

	gameId, ok = m.players[userId]
	return
}

func (m *GameDeskManager) AppEnterBackground(userId string) {
	wanRenGameId, wanRenOk := m.getWanRenInfo(userId)
	gameId, ok := m.getGameInfo(userId)

	if wanRenOk && ok {
		// 即在普通游戏，又在万人场直播，直接踢出直播
		m.LeaveGame(userId, util.GameType_WAN_REN, false)
	}

	if wanRenOk {
		if ok {
			// 即在普通游戏，又在万人场直播，直接踢出直播
			m.LeaveGame(userId, util.GameType_WAN_REN, false)
		} else {
			GetGameManager().onEnterBackground(wanRenGameId, userId)
			return
		}
	}

	if ok {
		GetGameManager().onEnterBackground(gameId, userId)
	}
}

func (m *GameDeskManager) AppEnterForeground(userId string) {
	wanRenGameId, wanRenOk := m.getWanRenInfo(userId)
	gameId, ok := m.getGameInfo(userId)

	if wanRenOk {
		GetGameManager().onEnterForeground(wanRenGameId, userId)
		return
	}

	if ok {
		GetGameManager().onEnterForeground(gameId, userId)
		return
	}
	domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_APP_ENTER_FOREGROUND), nil)
}

func (m *GameDeskManager) JoinWaitQueue(userId string) {
	wanRenGameId, wanRenOk := m.getWanRenInfo(userId)
	if wanRenOk {
		GetGameManager().onJoinWaitQueue(wanRenGameId, userId)
	}
}

func (m *GameDeskManager) LeaveWaitQueue(userId string) {
	wanRenGameId, wanRenOk := m.getWanRenInfo(userId)
	if wanRenOk {
		GetGameManager().onLeaveWaitQueue(wanRenGameId, userId)
	}
}

func (m *GameDeskManager) LookupBetGold(userId string, betUserId string) {
	wanRenGameId, wanRenOk := m.getWanRenInfo(userId)
	if wanRenOk {
		GetGameManager().onLookupBetGold(wanRenGameId, userId, betUserId)
	}
}

func (m *GameDeskManager) OnOpCards(userId string, msg *pb.MsgOpCardReq) {
	gameId, ok := m.getGameInfo(userId)
	if ok {
		GetGameManager().onOpCards(gameId, userId, msg)
		return
	}

	wanRenGameId, wanRenOk := m.getWanRenInfo(userId)
	if wanRenOk {
		GetGameManager().onOpCards(wanRenGameId, userId, msg)
	}
}

func (m *GameDeskManager) OnRewardInGame(userId string) bool {
	gameId, ok := m.getGameInfo(userId)
	if ok {
		return GetGameManager().onRewardInGame(gameId, userId)
	}

	wanRenGameId, wanRenOk := m.getWanRenInfo(userId)
	if wanRenOk {
		return GetGameManager().onRewardInGame(wanRenGameId, userId)
	}

	return false
}

func (m *GameDeskManager) OnChatMsg(userId string, msg *pb.MsgChat) {
	gameId, ok := m.getGameInfo(userId)
	if ok {
		GetGameManager().onChatMsg(gameId, userId, msg)
		return
	}

	wanRenGameId, wanRenOk := m.getWanRenInfo(userId)
	if wanRenOk {
		GetGameManager().onChatMsg(wanRenGameId, userId, msg)
	}
}

func (m *GameDeskManager) IsPlayingWanRen(userId string) bool {
	_, wanRenOk := m.getWanRenInfo(userId)
	_, ok := m.getGameInfo(userId)

	return wanRenOk && !ok
}

func (m *GameDeskManager) IsPlayingNormal(userId string) bool {
	_, ok := m.getGameInfo(userId)
	return ok
}

func (m *GameDeskManager) IsPlayingProps(userId string) bool {
	m.RLock()
	defer m.RUnlock()

	gameId, ok := m.players[userId]
	if !ok {
		return false
	}

	item := m.desks[gameId]
	if item == nil {
		return false
	}

	return util.IsGameTypeProps(item.gameType)
}

func (m *GameDeskManager) OnConsumeProps(userId string, itemType pb.MagicItemType, replaceCard int) {
	gameId, ok := m.getGameInfo(userId)
	if ok {
		GetGameManager().onConsumeProps(gameId, userId, itemType, replaceCard)
	}
}

func (m *GameDeskManager) OnOffline(userId string) {
	m.Lock()

	gameId := m.players[userId]
	item := m.desks[gameId]
	if item != nil {
		delete(item.players, userId)
		delete(m.players, userId)
		m.Unlock()
		GetGameManager().LeaveGame(gameId, userId, false)
	} else {
		m.Unlock()
	}

	m.Lock()
	if m.wanRenDesk != nil && m.wanRenDesk.hasPlayer(userId) {
		delete(m.wanRenDesk.players, userId)
		m.Unlock()
		GetGameManager().LeaveGame(m.wanRenDesk.gameId, userId, false)
	} else {
		m.Unlock()
	}
}

func (m *GameDeskManager) UpdateGoldInGame(userId string) {
	gameId, ok := m.getGameInfo(userId)
	if ok {
		GetGameManager().UpdateGoldInGame(gameId, userId)
		return
	}

	wanRenGameId, wanRenOk := m.getWanRenInfo(userId)
	if wanRenOk {
		GetGameManager().UpdateGoldInGame(wanRenGameId, userId)
	}
}
