package main

import "sort"

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

func getCardValue(card byte) byte {
	return card & MASK_VALUE
}

func getCardColor(card byte) byte {
	return card & MASK_COLOR
}

func GetCardType(cards []byte) int {
	return getCardType(cards)
}

func getCardLogicValue(card byte) byte {
	v := getCardValue(card)
	if v == 1 {
		return v + 13
	}
	return v
}

type CardComp []byte

func (p CardComp) Len() int { return len(p) }

func (p CardComp) Less(i, j int) bool {
	return getCardLogicValue(p[i]) < getCardLogicValue(p[j])
}

func (p CardComp) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func getCardType(cards []byte) int {
	sort.Sort(CardComp(cards))

	// 同花判断
	sameColor := false
	if getCardColor(cards[0]) == getCardColor(cards[1]) && getCardColor(cards[1]) == getCardColor(cards[2]) {
		sameColor = true
	}

	// 顺子判断
	lineCards := false
	if getCardLogicValue(cards[0])+1 == getCardLogicValue(cards[1]) && getCardLogicValue(cards[1])+1 == getCardLogicValue(cards[2]) {
		lineCards = true
	}

	// 顺子A23
	if !lineCards {
		if getCardValue(cards[0]) == 2 && getCardValue(cards[1]) == 3 && getCardValue(cards[2]) == 1 {
			lineCards = true
		}
	}

	// 同花顺类型
	if sameColor && lineCards {
		return CARD_TYPE_SHUN_JIN
	}

	// 顺子类型
	if !sameColor && lineCards {
		return CARD_TYPE_SHUN_ZI
	}

	// 金花类型
	if sameColor && !lineCards {
		return CARD_TYPE_JIN_HUA
	}

	// 对子
	isDouble := false
	if getCardLogicValue(cards[0]) == getCardLogicValue(cards[1]) ||
		getCardLogicValue(cards[0]) == getCardLogicValue(cards[2]) ||
		getCardLogicValue(cards[1]) == getCardLogicValue(cards[2]) {
		isDouble = true
	}

	// 豹子
	isBaoZi := false
	if getCardLogicValue(cards[0]) == getCardLogicValue(cards[1]) && getCardLogicValue(cards[1]) == getCardLogicValue(cards[2]) {
		isBaoZi = true
	}

	if isDouble {
		if isBaoZi {
			return CARD_TYPE_BAO_ZI
		}
		return CARD_TYPE_DOUBLE
	}

	// 特殊235
	if getCardValue(cards[0]) == 2 && getCardValue(cards[1]) == 3 && getCardValue(cards[2]) == 5 {
		return CARD_TYPE_SPECIAL
	}

	return CARD_TYPE_SINGLE
}
