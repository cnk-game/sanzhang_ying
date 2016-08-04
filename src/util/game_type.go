package util

import "pb"

type GameType int

const (
	MaxPlayerCount                   = 5
	GameType_Common_Level_1 GameType = 1 // 菜鸟
	GameType_Common_Level_2 GameType = 2 // 高级
	GameType_Common_Level_3 GameType = 3 // 精英
	GameType_Common_Level_4 GameType = 4 // 中级

	GameType_Props_Level_1 GameType = 11 // 道具赛菜鸟场
	GameType_Props_Level_2 GameType = 12 // 道具赛精英场
	GameType_Props_Level_3 GameType = 13 // 道具赛大师场

	GameType_SNG_Level_1 GameType = 21 // SNG
	GameType_SNG_Level_2 GameType = 22 // SNG
	GameType_SNG_Level_3 GameType = 23 // SNG

	GameType_WAN_REN = 30 // 万人场
)

// add by wangsq start -- 养鱼配置
type FishType int

const (
	Function_Switch    bool     = true // 功能开关 true开 - false关
	FishType_Fish_ID_1 FishType = 1    // 金龙鱼
	FishType_Fish_ID_2 FishType = 2    // 银龙鱼
	FishType_Fish_ID_3 FishType = 3    // 地图鱼
	FishType_Fish_ID_4 FishType = 4    // 凤尾鱼
)

const (
	FishType_Fish_Name_1 string = "河豚鱼"
	FishType_Fish_Name_2 string = "剑鱼"
	FishType_Fish_Name_3 string = "灯笼鱼"
	FishType_Fish_Name_4 string = "小丑鱼"
)

const (
	FishType_Fish_Price_1 int = 1000000
	FishType_Fish_Price_2 int = 500000
	FishType_Fish_Price_3 int = 100000
	FishType_Fish_Price_4 int = 10000
)

// add by wangsq end

func ToMatchType(gameType GameType) *pb.MatchType {
	switch gameType {
	case GameType_Common_Level_1:
		return pb.MatchType_COMMON_LEVEL1.Enum()
	case GameType_Common_Level_2:
		return pb.MatchType_COMMON_LEVEL2.Enum()
	case GameType_Common_Level_3:
		return pb.MatchType_COMMON_LEVEL3.Enum()
	case GameType_Common_Level_4:
		return pb.MatchType_COMMON_LEVEL4.Enum()
	case GameType_Props_Level_1:
		return pb.MatchType_MAGIC_ITEM_LEVEL1.Enum()
	case GameType_Props_Level_2:
		return pb.MatchType_MAGIC_ITEM_LEVEL2.Enum()
	case GameType_Props_Level_3:
		return pb.MatchType_MAGIC_ITEM_LEVEL3.Enum()
	case GameType_SNG_Level_1:
		return pb.MatchType_SNG_LEVEL1.Enum()
	case GameType_SNG_Level_2:
		return pb.MatchType_SNG_LEVEL2.Enum()
	case GameType_SNG_Level_3:
		return pb.MatchType_SNG_LEVEL3.Enum()
	case GameType_WAN_REN:
		return pb.MatchType_WAN_REN_GAME.Enum()
	}
	return nil
}

func IsGameTypeCommon(gameType GameType) bool {
	if gameType == GameType_Common_Level_1 {
		return true
	}
	if gameType == GameType_Common_Level_2 {
		return true
	}
	if gameType == GameType_Common_Level_3 {
		return true
	}
	if gameType == GameType_Common_Level_4 {
    	return true
    }
	return false
}

func IsGameTypeProps(gameType GameType) bool {
	if gameType == GameType_Props_Level_1 {
		return true
	}
	if gameType == GameType_Props_Level_2 {
		return true
	}
	if gameType == GameType_Props_Level_3 {
		return true
	}
	return false
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

func IsGameTypeWanRen(gameType GameType) bool {
	return gameType == GameType_WAN_REN
}

func GetGameSNGGoldPrice(gameType GameType) int {
	if gameType == GameType_SNG_Level_1 {
		return 1
	}
	if gameType == GameType_SNG_Level_2 {
		return 2
	}
	if gameType == GameType_SNG_Level_3 {
		return 3
	}
	return 0
}
