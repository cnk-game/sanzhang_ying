package models

import (
	"fmt"
	"labix.org/v2/mgo/bson"
	"poker-admin/util"
	"time"
)

const (
	synthesize_logC = "synthesize_log"
)

type SynthesizeInfoLog struct {
	Channel                  string    `bson:"channel"`
	NewUserCount             int       `bson:"newUserCount"`
	VaildNewUserCount        int       `bson:"vaildNewUserCount"`
	LoginUserCount           int       `bson:"loginUserCount"`
	PayCount                 int       `bson:"payCount"`
	PayTotalAmount           int       `bson:"payTotalAmount"`
	PayPlayerCount           int       `bson:"payPlayerCount"`
	After1DayRemainUserCount int       `bson:"after1DayRemainUserCount"`
	After2DayRemainUserCount int       `bson:"after2DayRemainUserCount"`
	After3DayRemainUserCount int       `bson:"after3DayRemainUserCount"`
	After5DayRemainUserCount int       `bson:"after5DayRemainUserCount"`
	After7DayRemainUserCount int       `bson:"after7DayRemainUserCount"`
	DateTime                 time.Time `bson:"datetime"`
}

///////////////////////////////////////////////////////////////////////////////////////
//
func InsertSynthesizeInfoLog(log *SynthesizeInfoLog) {
	session := util.GetLogSession()
	c := session.DB(util.LogDbName).C(synthesize_logC)
	defer session.Close()

	count, err := c.Find(bson.M{"datetime": log.DateTime, "channel": log.Channel}).Count()
	if err != nil {
		fmt.Println("UpdateXDayRemainBySynthesizeInfoLog => error1")
	}
	if count <= 0 {
		err := c.Insert(log)
		if err != nil {
			fmt.Println("InsertSynthesizeInfoLog => error")
		}
	}
}

func UpdateXDayRemainBySynthesizeInfoLog(afterXDayRemain, remainCount int, channel string, queryDate time.Time) {
	var updateKey = "after1DayRemainUserCount"
	if afterXDayRemain == 2 {
		updateKey = "after2DayRemainUserCount"
	} else if afterXDayRemain == 3 {
		updateKey = "after3DayRemainUserCount"
	} else if afterXDayRemain == 5 {
		updateKey = "after5DayRemainUserCount"
	} else if afterXDayRemain == 7 {
		updateKey = "after7DayRemainUserCount"
	}
	session := util.GetLogSession()
	c := session.DB(util.LogDbName).C(synthesize_logC)
	defer session.Close()

	count, err := c.Find(bson.M{"datetime": queryDate, "channel": channel}).Count()
	if err != nil {
		fmt.Println("UpdateXDayRemainBySynthesizeInfoLog => error1")
	}
	if count > 0 {
		err := c.Update(bson.M{"datetime": queryDate, "channel": channel}, bson.M{"$set": bson.M{updateKey: remainCount}})
		if err != nil {
			fmt.Println("UpdateXDayRemainBySynthesizeInfoLog => error2", channel, updateKey, queryDate)
		}
	}
}

func LoadSynthesizeLogList(channel string, b_year, b_month, b_day, e_year, e_month, e_day, pageIdx int, result interface{}) int {
	session := util.GetLogSession()
	c := session.DB(util.LogDbName).C(synthesize_logC)
	defer session.Close()
	beginTime := util.ParseTime(fmt.Sprintf("%v-%02v-%02v 00:00:00", b_year, b_month, b_day))
	endTime := util.ParseTime(fmt.Sprintf("%v-%02v-%02v 00:00:00", e_year, e_month, e_day))

	var err error
	var count = 0

	if len(channel) == 0 {
		err = c.Find(bson.M{"datetime": bson.M{"$gte": beginTime, "$lt": endTime}}).Sort("-datetime").Skip(page_count * pageIdx).Limit(page_count).All(result)
	} else {
		err = c.Find(bson.M{"channel": channel, "datetime": bson.M{"$gte": beginTime, "$lt": endTime}}).Sort("-datetime").Skip(page_count * pageIdx).Limit(page_count).All(result)
	}
	if err != nil {
		fmt.Println("LoadSynthesizeLogList => error1")
	}
	if len(channel) == 0 {
		count, err = c.Find(bson.M{"datetime": bson.M{"$gte": beginTime, "$lt": endTime}}).Count()
	} else {
		count, err = c.Find(bson.M{"channel": channel, "datetime": bson.M{"$gte": beginTime, "$lt": endTime}}).Count()
	}
	if err != nil {
		fmt.Println("LoadSynthesizeLogList => error2")
	}
	return count
}

func GetSynthesizeLogCount(year, month, day int) int {
	session := util.GetLogSession()
	c := session.DB(util.LogDbName).C(synthesize_logC)
	defer session.Close()

	logTime := util.ParseTime(fmt.Sprintf("%v-%02v-%02v 00:00:00", year, month, day))

	count, err := c.Find(bson.M{"datetime": logTime}).Count()
	if err != nil {
		fmt.Println("LoadSynthesizeLogList => error2")
	}
	return count
}

///////////////////////////////////////////////////////////////////////////////////////
// 每日综合信息

type ChannelGameInfo struct {
	Channel                  string
	TotalPlayerCount         int
	TodayNewPlayerCount      int
	TodayTotalNewPlayerCount int
	TodayLoginPlayerCount    int
	TodayPayPlayerCount      int
	TodayPayCount            int
	TodayTotalPay            int
}
type GameInfo struct {
	TotalPlayerCount         int
	TodayNewPlayerCount      int
	TodayTotalNewPlayerCount int
	TodayLoginPlayerCount    int
	TodayPayPlayerCount      int
	TodayPayCount            int
	TodayTotalPay            int
	ChannelGameInfoList      []*ChannelGameInfo
}

func GetNowGameInfo(gameInfo *GameInfo) {
	now := time.Now()
	//gameInfo := &GameInfo{}
	// 新用户列表
	newUsers := []*NewUserLog{}
	GetAllNewPlayers(now.Year(), int(now.Month()), now.Day(), &newUsers, "")
	// 有效新用户
	channelMap := make(map[string]*ChannelGameInfo)
	invaildUserIds := make(map[string]int)
	for _, item := range newUsers {
		isVaild := false
		if item.TotalOnlineSeconds >= ConfigMgr.VaildNewPlayerOnlineSecond && item.MatchTimes >= ConfigMgr.VaildNewPlayerMatchCount {
			isVaild = true
		}
		channeInfo, exist := channelMap[item.Channel]
		if exist {
			channeInfo.TodayTotalNewPlayerCount++
			if isVaild {
				channeInfo.TodayNewPlayerCount++
			}
		} else {
			channeInfo = &ChannelGameInfo{}
			channeInfo.Channel = item.Channel
			channeInfo.TodayTotalNewPlayerCount = 1
			if isVaild {
				channeInfo.TodayNewPlayerCount = 1
			} else {
				invaildUserIds[item.UserId] = 1
			}
			channelMap[item.Channel] = channeInfo
		}
	}
	// 登陆记录
	records := []*LoginRecord{}
	GetLoginRecords(now.Year(), int(now.Month()), now.Day(), &records)
	// 各渠道去除重复汇总
	recordCountMap := make(map[string]int)
	userIdMap := make(map[string]int)
	for _, item := range records {
		// 去除重复用户
		_, isExistUser := userIdMap[item.UserId]
		if !isExistUser {
			userIdMap[item.UserId] = 1
		} else {
			continue
		}
		// 是否有效
		_, isExistUser = invaildUserIds[item.UserId]
		if isExistUser {
			continue
		}
		count, exist := recordCountMap[item.Channel]
		if exist {
			recordCountMap[item.Channel] = count + 1
		} else {
			recordCountMap[item.Channel] = 1
		}
	}
	for channel, count := range recordCountMap {
		channeInfo, exist := channelMap[channel]
		if exist {
			channeInfo.TodayLoginPlayerCount = count
		} else {
			channeInfo = &ChannelGameInfo{}
			channeInfo.Channel = channel
			channeInfo.TodayLoginPlayerCount = count
			channelMap[channel] = channeInfo
		}
	}
	// 支付记录
	payRecords := []*PayLog{}
	GetAllPayLogList(now.Year(), int(now.Month()), now.Day(), &payRecords)
	payUserCountRecordMap := make(map[string]int)
	payUserIdMap := make(map[string]int)
	for _, item := range payRecords {
		channeInfo, exist := channelMap[item.Channel]
		if exist {
			// 付费次数
			channeInfo.TodayPayCount++
			// 付费金额
			channeInfo.TodayTotalPay += item.Amount
		} else {
			channeInfo = &ChannelGameInfo{}
			channeInfo.Channel = item.Channel
			channeInfo.TodayPayCount = 1
			channeInfo.TodayTotalPay = item.Amount
			channelMap[item.Channel] = channeInfo
		}
		// 付费人数汇总
		_, exist = payUserIdMap[item.UserId]
		if !exist {
			count, exist := payUserCountRecordMap[item.Channel]
			if exist {
				payUserCountRecordMap[item.Channel] = count + 1
			} else {
				payUserCountRecordMap[item.Channel] = 1
			}
			payUserIdMap[item.UserId] = 1
		}
	}
	// 付费人数
	for channel, count := range payUserCountRecordMap {
		channeInfo, exist := channelMap[channel]
		if exist {
			channeInfo.TodayPayPlayerCount = count
		} else {
			channeInfo = &ChannelGameInfo{}
			channeInfo.Channel = channel
			channeInfo.TodayPayPlayerCount = count
			channelMap[channel] = channeInfo
		}
	}
	// 汇总
	for _, channelInfo := range channelMap {
		channelInfo.TotalPlayerCount = GetTotalPlayerCount(channelInfo.Channel)
		gameInfo.TodayNewPlayerCount += channelInfo.TodayNewPlayerCount
		gameInfo.TodayTotalNewPlayerCount += channelInfo.TodayTotalNewPlayerCount
		gameInfo.TodayLoginPlayerCount += channelInfo.TodayLoginPlayerCount
		gameInfo.TodayPayPlayerCount += channelInfo.TodayPayPlayerCount
		gameInfo.TodayPayCount += channelInfo.TodayPayCount
		gameInfo.TodayTotalPay += channelInfo.TodayTotalPay
		gameInfo.ChannelGameInfoList = append(gameInfo.ChannelGameInfoList, channelInfo)
	}
	// 总人数
	gameInfo.TotalPlayerCount = GetTotalPlayerCount("")
	//fmt.Println(gameInfo)
}
