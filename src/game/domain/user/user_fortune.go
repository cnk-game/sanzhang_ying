package user

import (
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"util"
)

type UserFish struct {
	RecordId string    `bson:"recordId"`
	FromUid  string    `bson:"fromId"`
	GetTime  time.Time `bson:"getTime"`
	FishType int       `bson:"fishType"`
	Count    int       `bson:"count"`
}

type FToken struct {
	Token     string `bson:"token"`
	TimeStamp int    `bson:"timeStamp"`
}

type UserVipState struct {
	VipTaskId    int   `bson:"vipTaskId"`
	StartTime    int64 `bson:"startTime"`
	LastGainTime int64 `bson:"lastGainTime"`
}

type BoxLog struct {
	Reason   string    `bson:"reason"`
	Gold     int64     `bson:"gold"`
	Savings  int64     `bson:"savings"`
	DateTime time.Time `bson:"DateTime"`
	ToUid    string    `bson:"toUid"`
	LogId    string    `bson:"logId"`
}

type UserBox struct {
	IsOpen   bool              `bson:"isOpen"`
	Pwd      string            `bson:"pwd"`
	Savings  int64             `bson:"savings"`
	IsSetPwd bool              `bson:"isSetPwd"`
	BoxLogs  map[string]BoxLog `bson:"logs"`
}

type UserFortune struct {
	UserId                   string                  `bson:"userId"`
	IsRobot                  bool                    `bson:"isRobot"`
	Gold                     int64                   `bson:"gold"`
	Diamond                  int                     `bson:"diamond"`
	Score                    int                     `bson:"score"`
	Exp                      int                     `bson:"exp"`
	VipLevel                 int                     `bson:"vipLevel"`
	VipStartTime             time.Time               `bson:"vipStartTime"`
	VipValidDays             int                     `bson:"vipValidDays"`
	VipTaskStates            map[string]UserVipState `bson:"vipTaskStates"`
	TodayRecharge            int                     `bson:"todayRecharge"`
	CurWeekRecharge          int                     `bson:"curWeekRecharge"`
	TodayEarnings            int                     `bson:"todayEarnings"`
	CurWeekEarnings          int                     `bson:"curWeekEarnings"`
	LastRechargeTime         time.Time               `bson:"lastRechargeTime"`
	LastEarnTime             time.Time               `bson:"lastEarnTime"`
	DoubleCardCount          int                     `bson:"doubleCardCount"`
	ForbidCardCount          int                     `bson:"forbidCardCount"`
	ChangeCardCount          int                     `bson:"changeCardCount"`
	FirstRecharge10          bool                    `bson:"firstRecharge10"`
	FirstRecharge20          bool                    `bson:"firstRecharge20"`
	FirstRecharge30          bool                    `bson:"firstRecharge30"`
	FirstRecharge50          bool                    `bson:"firstRecharge50"`
	FirstRecharge100         bool                    `bson:"firstRecharge100"`
	FirstRecharge500         bool                    `bson:"firstRecharge500"`
	FirstRecharge            bool                    `bson:"firstRecharge"`
	DailyGiftBagUseTime      time.Time               `bson:"dailyGiftBagUseTime"`
	BuyDailyGiftBag          bool                    `bson:"buyDailyGiftBag"`
	LastPayTime              time.Time               `bson:"lastPayTime"`
	GainedFirstRechargeBonus bool                    `bson:"gainedFirstRechargeBonus"`
	// add by wangsq start
	Charm             int                 `bson:"charm"`
	FishToken         FToken              `bson:"fishToken"`
	FishInfo          map[string]UserFish `bson:"fishInfo"`
	Horn              int                 `bson:"horn"`
	CharmExchangeInfo map[string]int      `bson:"charmExchangeInfo"`
	SafeBox           UserBox             `bson:"safeBox"`
	// add by wangsq end

	// add by yelong
	ChangeGameTypeNotifyDay int `bson:"changeGameTypeNotifyDay"`

	isDailyGiftBagPay bool
	CurMonRec         int `bson:"curMonRec"`
	CurMonEarn        int `bson:"curMonEarn"`
	CurMonWin         int `bson:"curMonWin"`
}

const (
	userFortuneC     = "user_fortune"
	exchangeGoldRate = 10000 // 钻石兑换金币比率
)

func FindUserFortune(userId string) (*UserFortune, error) {
	fortune := &UserFortune{}
	fortune.SafeBox = UserBox{false, "", 0, false, map[string]BoxLog{}}
	fortune.UserId = userId

	userFortuneTableC, er := GetUserFortuneTable(userId)
	if er != nil {
		return fortune, er
	}

	err := util.WithUserCollection(userFortuneTableC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"userId": userId}).One(fortune)
	})

	glog.V(2).Info("===>用户财富userId:", userId, " vipLevel:", fortune.VipLevel)

	/*if err == mgo.ErrNotFound {
		return fortune, nil
	}*/

	return fortune, err
}

func SaveFortune(fortune *UserFortune) error {
	userId := fortune.UserId
	userFortuneTableC, er := GetUserFortuneTable(userId)
	if er != nil {
		return er
	}

	return util.WithSafeUserCollection(userFortuneTableC, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"userId": fortune.UserId}, fortune)
		return err
	})
}

func SaveDiamond(userId string, diamond int, isRobot bool) error {
	if isRobot {
		return nil
	}

	userFortuneTableC, er := GetUserFortuneTable(userId)
	if er != nil {
		return er
	}

	return util.WithSafeUserCollection(userFortuneTableC, func(c *mgo.Collection) error {
		return c.Update(bson.M{"userId": userId}, bson.M{"$set": bson.M{"diamond": diamond}})
	})
}

func SaveGold(userId string, gold int64, isRobot bool) error {
	if isRobot {
		return nil
	}

	userFortuneTableC, er := GetUserFortuneTable(userId)
	if er != nil {
		return er
	}

	return util.WithSafeUserCollection(userFortuneTableC, func(c *mgo.Collection) error {
		return c.Update(bson.M{"userId": userId}, bson.M{"$set": bson.M{"gold": gold}})
	})
}

// add by wangsq start
func SaveCharm(userId string, charm int, isRobot bool) error {
	if isRobot {
		return nil
	}

	userFortuneTableC, er := GetUserFortuneTable(userId)
	if er != nil {
		return er
	}

	return util.WithSafeUserCollection(userFortuneTableC, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"userId": userId}, bson.M{"$set": bson.M{"charm": charm}})
		return err
	})
}

func SaveHorn(userId string, horn int, isRobot bool) error {
	if isRobot {
		return nil
	}

	userFortuneTableC, er := GetUserFortuneTable(userId)
	if er != nil {
		return er
	}

	return util.WithSafeUserCollection(userFortuneTableC, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"userId": userId}, bson.M{"$set": bson.M{"horn": horn}})
		return err
	})
}

func FindGoldLimitUserCount(limitLow, limitHigh int) (int, error) {
	count := 0
	err := util.WithSafeUserCollection(userFortuneC, func(c *mgo.Collection) error {
		n, err := c.Find(bson.M{"gold": bson.M{"$gte": limitLow, "$lte": limitHigh}}).Count()
		count = n
		return err
	})
	return count, err
}

func SumAllGolds() (int64, error) {
	var result []struct {
		Id    string "_id"
		Value int64  "value"
	}
	err := util.WithSafeUserCollection(userFortuneC, func(c *mgo.Collection) error {
		job := &mgo.MapReduce{
			Map:    "function() { emit('all', this.gold); }",
			Reduce: "function(key, values) { var sum = 0; for (var i in values) { sum += values[i]; }; return sum; }",
		}

		_, err := c.Find(nil).MapReduce(job, &result)
		return err
	})
	return result[0].Value, err
}

// add by wangsq end
