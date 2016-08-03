package user

import (
	"code.google.com/p/goprotobuf/proto"
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"pb"
	"time"
	"util"
)

type MatchRecord struct {
	UserId                string    `bson:"userId"`
	WinTimes              int       `bson:"winTimes"`
	LoseTimes             int       `bson:"loseTimes"`
	CurDayEarnGold        int       `bson:"curDayEarnGold"`
	CurWeekEarnGold       int       `bson:"curWeekEarnGold"`
	MaxCards              []int32   `bson:"maxCards"`
	DayEarnGoldResetTime  time.Time `bson:"dayEarnGoldResetTime"`
	WeekEarnGoldResetTime time.Time `bson:"weekEarnGoldResetTime"`
	CurMonWin             int       `bson:"curMonWin"`
	MonWinRetTime         time.Time `bson:"monWinRetTime"`
}

func (r *MatchRecord) BuildMessage() *pb.UserMatchRecordDef {
	r.resetEarnGold()

	msg := &pb.UserMatchRecordDef{}
	msg.PlayWinCount = proto.Int(r.WinTimes)
	msg.PlayTotalCount = proto.Int(r.WinTimes + r.LoseTimes)
	msg.TheDayWinGold = proto.Int(r.CurDayEarnGold)
	msg.TheWeekWinGold = proto.Int(r.CurWeekEarnGold)
	msg.MaxCards = r.MaxCards

	return msg
}

func (r *MatchRecord) resetEarnGold() {
	now := time.Now()
	if !util.CompareDate(now, r.DayEarnGoldResetTime) {
		r.CurDayEarnGold = 0
		r.DayEarnGoldResetTime = now
	}

	if !util.CompareDate(now, r.WeekEarnGoldResetTime) && now.Weekday() == time.Monday {
		r.CurWeekEarnGold = 0
		r.WeekEarnGoldResetTime = now
	}
}

func (r *MatchRecord) AddEarnGold(gold int) {
	r.resetEarnGold()
	r.CurDayEarnGold += gold
	r.CurWeekEarnGold += gold
	if gold > 0 {
		r.AddMonWindTimes()
	}
}

func (r *MatchRecord) AddMonWindTimes() {
	now := time.Now()
	if !util.CompareDate(now, r.MonWinRetTime) {
		if now.Day() == 1 {
			r.MonWinRetTime = now
			r.CurMonWin = 0
		}
	}
	r.CurMonWin += 1
	GetUserFortuneManager().UpdateCurMonWinTimesRankingList(r.UserId, r.CurMonWin)
}

const (
	matchRecordC = "match_record"
)

func FindMatchRecord(userId string) *MatchRecord {
	r := &MatchRecord{}
	r.UserId = userId
	matchRecordTableC, er := GetUserMatchRecordTable(userId)
	if er != nil {
		glog.Error("FindMatchRecord GetUserMatchRecordTable err")
	}
	err := util.WithUserCollection(matchRecordTableC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"userId": userId}).One(r)
	})
	if err != nil {
		glog.Error("加载玩家比赛记录失败err:", err)
	}
	return r
}

func SaveMatchRecord(r *MatchRecord) error {
	matchRecordTableC, er := GetUserMatchRecordTable(r.UserId)
	if er != nil {
		return er
	}

	return util.WithSafeUserCollection(matchRecordTableC, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"userId": r.UserId}, r)
		return err
	})
}
