package models

import (
	"fmt"
	"labix.org/v2/mgo/bson"
	"poker-admin/util"
)

const (
	pay_logC      = "pay_log"
	page_count    = 20
	PageItemCount = page_count
)

var ConfigMgr *ConfigManager

type ConfigManager struct {
	VaildNewPlayerOnlineSecond int
	VaildNewPlayerMatchCount   int
}

func init() {
	value1, value2 := GetChannelCheckConfig()
	if 0 == value1 {
		value1 = 5 * 60
	}
	if 0 == value2 {
		value2 = 1
	}
	ConfigMgr = &ConfigManager{value1, value2}
}

type ChannelCheckConfig struct {
	OnlineSecond int `bson:"OnlineSecond"`
	MatchCount   int `bson:"MatchCount"`
}

func GetChannelCheckConfig() (onlineSecond int, matchCount int) {
	session := util.GetSession()
	c := session.DB(util.WebsiteDBName).C("channel_check_config")
	defer session.Close()

	config := ChannelCheckConfig{}
	err := c.Find(nil).One(&config)
	if err != nil {
		fmt.Println("GetChannelCheckConfig error => 1")
	}

	return config.OnlineSecond, config.MatchCount
}

func UpdateChannelCheckConfig(onlineSecond, matchCount int) {
	session := util.GetSession()
	c := session.DB(util.WebsiteDBName).C("channel_check_config")
	defer session.Close()

	_, err := c.RemoveAll(nil)
	if err != nil {
		fmt.Println("UpdateChannelCheckConfig error => 1")
	}

	err = c.Insert(ChannelCheckConfig{onlineSecond, matchCount})
	if err != nil {
		fmt.Println("UpdateChannelCheckConfig error => 2")
	}
}

////////////////////////////////////////////////////////////////////////
type MatchConfig struct {
	CommonLevel1Single      int `bson:"CommonLevel1Single"`
	CommonLevel1Double      int `bson:"CommonLevel1Double"`
	CommonLevel1ShunZi      int `bson:"CommonLevel1ShunZi"`
	CommonLevel1JinHua      int `bson:"CommonLevel1JinHua"`
	CommonLevel1ShunJin     int `bson:"CommonLevel1ShunJin"`
	CommonLevel1BaoZi       int `bson:"CommonLevel1BaoZi"`
	CommonLevel1WinGold     int `bson:"CommonLevel1WinGold"`
	CommonLevel1WinRateHigh int `bson:"CommonLevel1WinRateHigh"`
	CommonLevel1LoseGold    int `bson:"CommonLevel1LoseGold"`
	CommonLevel1WinRateLow  int `bson:"CommonLevel1WinRateLow"`

	CommonLevel2Single      int `bson:"CommonLevel2Single"`
	CommonLevel2Double      int `bson:"CommonLevel2Double"`
	CommonLevel2ShunZi      int `bson:"CommonLevel2ShunZi"`
	CommonLevel2JinHua      int `bson:"CommonLevel2JinHua"`
	CommonLevel2ShunJin     int `bson:"CommonLevel2ShunJin"`
	CommonLevel2BaoZi       int `bson:"CommonLevel2BaoZi"`
	CommonLevel2WinGold     int `bson:"CommonLevel2WinGold"`
	CommonLevel2WinRateHigh int `bson:"CommonLevel2WinRateHigh"`
	CommonLevel2LoseGold    int `bson:"CommonLevel2LoseGold"`
	CommonLevel2WinRateLow  int `bson:"CommonLevel2WinRateLow"`

	CommonLevel3Single      int `bson:"CommonLevel3Single"`
	CommonLevel3Double      int `bson:"CommonLevel3Double"`
	CommonLevel3ShunZi      int `bson:"CommonLevel3ShunZi"`
	CommonLevel3JinHua      int `bson:"CommonLevel3JinHua"`
	CommonLevel3ShunJin     int `bson:"CommonLevel3ShunJin"`
	CommonLevel3BaoZi       int `bson:"CommonLevel3BaoZi"`
	CommonLevel3WinGold     int `bson:"CommonLevel3WinGold"`
	CommonLevel3WinRateHigh int `bson:"CommonLevel3WinRateHigh"`
	CommonLevel3LoseGold    int `bson:"CommonLevel3LoseGold"`
	CommonLevel3WinRateLow  int `bson:"CommonLevel3WinRateLow"`

	ItemLevel1Single      int `bson:"ItemLevel1Single"`
	ItemLevel1Double      int `bson:"ItemLevel1Double"`
	ItemLevel1ShunZi      int `bson:"ItemLevel1ShunZi"`
	ItemLevel1JinHua      int `bson:"ItemLevel1JinHua"`
	ItemLevel1ShunJin     int `bson:"ItemLevel1ShunJin"`
	ItemLevel1BaoZi       int `bson:"ItemLevel1BaoZi"`
	ItemLevel1WinGold     int `bson:"ItemLevel1WinGold"`
	ItemLevel1WinRateHigh int `bson:"ItemLevel1WinRateHigh"`
	ItemLevel1LoseGold    int `bson:"ItemLevel1LoseGold"`
	ItemLevel1WinRateLow  int `bson:"ItemLevel1WinRateLow"`

	ItemLevel2Single      int `bson:"ItemLevel2Single"`
	ItemLevel2Double      int `bson:"ItemLevel2Double"`
	ItemLevel2ShunZi      int `bson:"ItemLevel2ShunZi"`
	ItemLevel2JinHua      int `bson:"ItemLevel2JinHua"`
	ItemLevel2ShunJin     int `bson:"ItemLevel2ShunJin"`
	ItemLevel2BaoZi       int `bson:"ItemLevel2BaoZi"`
	ItemLevel2WinGold     int `bson:"ItemLevel2WinGold"`
	ItemLevel2WinRateHigh int `bson:"ItemLevel2WinRateHigh"`
	ItemLevel2LoseGold    int `bson:"ItemLevel2LoseGold"`
	ItemLevel2WinRateLow  int `bson:"ItemLevel2WinRateLow"`

	ItemLevel3Single      int `bson:"ItemLevel3Single"`
	ItemLevel3Double      int `bson:"ItemLevel3Double"`
	ItemLevel3ShunZi      int `bson:"ItemLevel3ShunZi"`
	ItemLevel3JinHua      int `bson:"ItemLevel3JinHua"`
	ItemLevel3ShunJin     int `bson:"ItemLevel3ShunJin"`
	ItemLevel3BaoZi       int `bson:"ItemLevel3BaoZi"`
	ItemLevel3WinGold     int `bson:"ItemLevel3WinGold"`
	ItemLevel3WinRateHigh int `bson:"ItemLevel3WinRateHigh"`
	ItemLevel3LoseGold    int `bson:"ItemLevel3LoseGold"`
	ItemLevel3WinRateLow  int `bson:"ItemLevel3WinRateLow"`

	SngLevel1Single      int `bson:"SngLevel1Single"`
	SngLevel1Double      int `bson:"SngLevel1Double"`
	SngLevel1ShunZi      int `bson:"SngLevel1ShunZi"`
	SngLevel1JinHua      int `bson:"SngLevel1JinHua"`
	SngLevel1ShunJin     int `bson:"SngLevel1ShunJin"`
	SngLevel1BaoZi       int `bson:"SngLevel1BaoZi"`
	SngLevel1WinGold     int `bson:"SngLevel1WinGold"`
	SngLevel1WinRateHigh int `bson:"SngLevel1WinRateHigh"`
	SngLevel1LoseGold    int `bson:"SngLevel1LoseGold"`
	SngLevel1WinRateLow  int `bson:"SngLevel1WinRateLow"`

	SngLevel2Single      int `bson:"SngLevel2Single"`
	SngLevel2Double      int `bson:"SngLevel2Double"`
	SngLevel2ShunZi      int `bson:"SngLevel2ShunZi"`
	SngLevel2JinHua      int `bson:"SngLevel2JinHua"`
	SngLevel2ShunJin     int `bson:"SngLevel2ShunJin"`
	SngLevel2BaoZi       int `bson:"SngLevel2BaoZi"`
	SngLevel2WinGold     int `bson:"SngLevel2WinGold"`
	SngLevel2WinRateHigh int `bson:"SngLevel2WinRateHigh"`
	SngLevel2LoseGold    int `bson:"SngLevel2LoseGold"`
	SngLevel2WinRateLow  int `bson:"SngLevel2WinRateLow"`

	SngLevel3Single      int `bson:"SngLevel3Single"`
	SngLevel3Double      int `bson:"SngLevel3Double"`
	SngLevel3ShunZi      int `bson:"SngLevel3ShunZi"`
	SngLevel3JinHua      int `bson:"SngLevel3JinHua"`
	SngLevel3ShunJin     int `bson:"SngLevel3ShunJin"`
	SngLevel3BaoZi       int `bson:"SngLevel3BaoZi"`
	SngLevel3WinGold     int `bson:"SngLevel3WinGold"`
	SngLevel3WinRateHigh int `bson:"SngLevel3WinRateHigh"`
	SngLevel3LoseGold    int `bson:"SngLevel3LoseGold"`
	SngLevel3WinRateLow  int `bson:"SngLevel3WinRateLow"`

	WanSingle      int `bson:"WanSingle"`
	WanDouble      int `bson:"WanDouble"`
	WanShunZi      int `bson:"WanShunZi"`
	WanJinHua      int `bson:"WanJinHua"`
	WanShunJin     int `bson:"WanShunJin"`
	WanBaoZi       int `bson:"WanBaoZi"`
	WanWinGold     int `bson:"WanWinGold"`
	WanWinRateHigh int `bson:"WanWinRateHigh"`
	WanLoseGold    int `bson:"WanLoseGold"`
	WanWinRateLow  int `bson:"WanWinRateLow"`
}

func GetMatchConfig() *MatchConfig {
	session := util.GetLogSession()
	c := session.DB(util.WebsiteDBName).C("match_config")
	defer session.Close()
	config := &MatchConfig{}
	err := c.Find(bson.M{}).One(config)
	if err != nil {
		fmt.Println("GetMatchConfig error")
		return nil
	}
	return config
}

func SaveMatchConfig(config *MatchConfig) {
	session := util.GetLogSession()
	c := session.DB(util.WebsiteDBName).C("match_config")
	defer session.Close()

	_, err := c.RemoveAll(nil)
	if err != nil {
		fmt.Println("SaveMatchConfig error => 1")
	}

	err = c.Insert(config)
	if err != nil {
		fmt.Println("SaveMatchConfig error => 2")
	}
}

////////////////////////////////////////////////////////////////////////
type PrizeVersion struct {
	Version string
}

func GetPrizeVersion() string {
	session := util.GetLogSession()
	c := session.DB(util.WebsiteDBName).C("prize_version")
	defer session.Close()
	config := &PrizeVersion{}
	err := c.Find(bson.M{}).One(config)
	if err != nil {
		fmt.Println("GetMatchConfig error")
		return ""
	}
	return config.Version
}
func SavePrizeVersion(version string) {
	session := util.GetLogSession()
	c := session.DB(util.WebsiteDBName).C("prize_version")
	defer session.Close()

	_, err := c.RemoveAll(nil)
	if err != nil {
		fmt.Println("SavePrizeVersion error => 1")
	}
	config := &PrizeVersion{}
	config.Version = version
	err = c.Insert(config)
	if err != nil {
		fmt.Println("SavePrizeVersion error => 2")
	}
}
