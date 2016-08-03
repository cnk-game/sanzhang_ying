package slots

import (
	"fmt"
	domainGame "game/domain/game"
	"testing"
)

const (
	MASK_COLOR         = 0xF0 // 花色掩码
	MASK_VALUE         = 0x0F // 数值掩码
	CARD_TYPE_SINGLE   = 1    // 单牌类型
	CARD_TYPE_DOUBLE   = 2    // 对子类型
	CARD_TYPE_SHUN_ZI  = 3    // 顺子类型
	CARD_TYPE_JIN_HUA  = 4    // 金花类型
	CARD_TYPE_SHUN_JIN = 5    // 同花顺类型(顺金)
	CARD_TYPE_BAO_ZI   = 6    // 豹子类型
	CARD_TYPE_SPECIAL  = 7    // 特殊类型(不同花色235)
)

func PrintCardType(cardType int) string {
	switch cardType {
	case CARD_TYPE_SINGLE:
		return "单牌类型"
	case CARD_TYPE_DOUBLE:
		return "对子类型"
	case CARD_TYPE_SHUN_ZI:
		return "顺子类型"
	case CARD_TYPE_JIN_HUA:
		return "金花类型"
	case CARD_TYPE_SHUN_JIN:
		return "顺金类型"
	case CARD_TYPE_BAO_ZI:
		return "豹子类型"
	case CARD_TYPE_SPECIAL:
		return "特殊类型"
	}
	return ""
}

func TestRandomCards(t *testing.T) {

	cardTypes := make(map[int]int)

	for i := 0; i < 100000; i++ {
		cards := GetRandomCards()
		cardTypes[domainGame.GetCardType(cards)]++
	}

	total := 0
	for _, v := range cardTypes {
		total += v
	}

	for k, v := range cardTypes {
		fmt.Println("牌型:", PrintCardType(k), " 数目:", v, " 百分比:", float64(v)/float64(total)*100, "%")
	}
}
