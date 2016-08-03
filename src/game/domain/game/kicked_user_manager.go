package game

import (
	"sync"
)

type KickedUser struct {
	GameId   int
	KickTime int64
}

type KickedUserManager struct {
	sync.RWMutex
	//记录用户被踢的数据
	KickedLog map[string]*KickedUser //userId ==> KickedUser
}

var kicdedManager *KickedUserManager

func init() {

	kicdedManager = &KickedUserManager{}
	kicdedManager.KickedLog = make(map[string]*KickedUser)

}

func GetKickedUserManager() *KickedUserManager {
	return kicdedManager
}
