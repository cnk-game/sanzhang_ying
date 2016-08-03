package game

import (
	mgo "gopkg.in/mgo.v2"
	"strconv"
	"time"
	"util"
)

type GameLog struct {
	Username   string    `bson:"username"`
	GameId     int       `bson:"gameId"`
	GameType   int       `bson:"gameType"`
	SingleBet  int       `bson:"singleBet"`
	TotalBet   int       `bson:"totalBet"`
	CurRound   int       `bson:"curRound"`
	MatchTimes int       `bson:"matchTimes"`
	SeenCard   bool      `bson:"seenCard"`
	BetGold    int       `bson:"betGold"`
	CurGold    int       `bson:"curGold"`
	Winner     string    `bson:"winner"`
	EarnGold   int       `bson:"earnGold"`
	Event      string    `bson:"event"`
	Time       time.Time `bson:"time"`
}

type CharmLog struct {
	UserId  string    `bson:"userId"`
	Count   int       `bson:"count"`
	Channel int       `bson:"channel"`
	Time    time.Time `bson:"time"`
}

const (
	gameLogC  = "game_log"
	charmLogC = "charm_log"
)

func SaveGameLog(log *GameLog) error {
	now := time.Now()
	log.Time = now
	cur_C := gameLogC + "_" + strconv.Itoa(int(now.Year())) + strconv.Itoa(int(now.Month())) + strconv.Itoa(int(now.Day()))

	return util.WithLogCollection(cur_C, func(c *mgo.Collection) error {
		return c.Insert(log)
	})
}

func SaveCharmLog(userId string, count int, channel int) error {
	now := time.Now()
	cur_C := charmLogC + "_" + strconv.Itoa(int(now.Year())) + strconv.Itoa(int(now.Month())) + strconv.Itoa(int(now.Day()))
	log := CharmLog{userId, count, channel, now}

	return util.WithLogCollection(cur_C, func(c *mgo.Collection) error {
		return c.Insert(&log)
	})
}
