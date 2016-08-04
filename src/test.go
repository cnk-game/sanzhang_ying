package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

type GameType int

const (
	MaxPlayerCount                   = 5
	GameType_Common_Level_1 GameType = 1 // 常规赛初级场
	GameType_Common_Level_2 GameType = 2 // 常规赛中级场
	GameType_Common_Level_3 GameType = 3 // 常规赛高级场

	GameType_Props_Level_1 GameType = 11 // 道具赛菜鸟场
	GameType_Props_Level_2 GameType = 12 // 道具赛精英场
	GameType_Props_Level_3 GameType = 13 // 道具赛大师场

	GameType_SNG_Level_1 GameType = 21 // SNG
	GameType_SNG_Level_2 GameType = 22 // SNG
	GameType_SNG_Level_3 GameType = 23 // SNG

	GameType_WAN_REN = 30 // 万人场
)

type GameDeskUser struct {
	UserId  string
	IsRobot bool
}

type GameDesk struct {
	gameType GameType
	gameId   int
	players  map[string]*GameDeskUser
	Locked   bool
}

type GameDeskManager struct {
	wanRenDesk *GameDesk
	desks      map[int]*GameDesk // gameId==>gameDesk
	players    map[string]int    // userId==>gameId
	gameId     int
}

func IsGameTypeWanRen(gameType GameType) bool {
	return gameType == GameType_WAN_REN
}

func (m *GameDeskManager) selectGameDesk(gameType GameType, ignoreGameId int, isRobot bool) *GameDesk {
	if IsGameTypeWanRen(gameType) {
		return m.wanRenDesk
	}

	if len(m.desks) <= 0 {
		return nil
	}

	count := 0
	gameId := 0
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

		itemCount := len(item.players)
		if count == 0 {
			count = itemCount
			gameId = item.gameId
		}

		if itemCount < 5 && itemCount > count {
			count = itemCount
			gameId = item.gameId
		}
	}

	// 当前所有游戏都已满员
	if count >= 5 {
		return nil
	}

	return m.desks[gameId]
}

func IsGameTypeSNG(gameType GameType) bool {
	if gameType == GameType_SNG_Level_1 {
		return true
	}
	if gameType == GameType_SNG_Level_2 {
		return true
	}
	if gameType == GameType_SNG_Level_3 {
		return true
	}
	return false
}

func (desk *GameDesk) allowRobotIn() bool {
	if IsGameTypeSNG(desk.gameType) {
		return len(desk.players) <= 3
	}

	playerCount := 0
	for _, p := range desk.players {
		if p.IsRobot {
			playerCount++
		}
	}

	return playerCount <= 1
}

var m *GameDeskManager

func init() {
	m = &GameDeskManager{}
	m.desks = make(map[int]*GameDesk)
	m.players = make(map[string]int)
}

func testGameDesk(gameType GameType, ignoreGameId int, isRobot bool) {
	desk := m.selectGameDesk(GameType_Common_Level_1, 0, false)
	if desk == nil {
		desk = &GameDesk{}
		desk.gameType = gameType
		desk.gameId++
		desk.players = make(map[string]*GameDeskUser)
		if IsGameTypeWanRen(gameType) {
			m.wanRenDesk = desk
		} else {
			m.desks[desk.gameId] = desk
		}
	}
}

func main() {
	go func() {
		fmt.Println("start")
		start := time.Now()
		for i := 0; i < 10000000; i++ {
			testGameDesk(GameType_Common_Level_1, 0, false)
			testGameDesk(GameType_Common_Level_2, 0, false)
			testGameDesk(GameType_Common_Level_3, 0, false)
		}
		fmt.Println(" elapsed:", time.Since(start))
	}()
	log.Fatal(http.ListenAndServe(":9000", nil))
}
