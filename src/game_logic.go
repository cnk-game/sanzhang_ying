package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"time"
	"util"
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

func PrintCardType(cardType int) {
	switch cardType {
	case CARD_TYPE_SINGLE:
		fmt.Println("单牌类型")
	case CARD_TYPE_DOUBLE:
		fmt.Println("对子类型")
	case CARD_TYPE_SHUN_ZI:
		fmt.Println("顺子类型")
	case CARD_TYPE_JIN_HUA:
		fmt.Println("金花类型")
	case CARD_TYPE_SHUN_JIN:
		fmt.Println("顺金类型")
	case CARD_TYPE_BAO_ZI:
		fmt.Println("豹子类型")
	case CARD_TYPE_SPECIAL:
		fmt.Println("特殊类型")
	}
}

type GameLogic struct {
	cardList  []byte
	pos       []int
	sortedPos []int
	sorted    bool
	gameType  int
}

type ByteSlice []byte

func (p ByteSlice) Len() int           { return len(p) }
func (p ByteSlice) Less(i, j int) bool { return rand.Float32() < 0.5 }
func (p ByteSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type CardComp []byte

func (p CardComp) Len() int { return len(p) }

func (p CardComp) Less(i, j int) bool {
	return getCardLogicValue(p[i]) < getCardLogicValue(p[j])
}

func (p CardComp) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// 牌型排序
type CardTypeComp struct {
	pos   []int
	logic *GameLogic
}

func (p CardTypeComp) Len() int { return len(p.pos) }

func (p CardTypeComp) Less(i, j int) bool {
	return p.logic.CompareCards(p.logic.GetCards(p.pos[i]), p.logic.GetCards(p.pos[j]))
}

func (p CardTypeComp) Swap(i, j int) { p.pos[i], p.pos[j] = p.pos[j], p.pos[i] }

func NewCardTypeComp(logic *GameLogic) *CardTypeComp {
	t := &CardTypeComp{}
	t.logic = logic

	return t
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewGameLogic(gameType int) *GameLogic {
	logic := &GameLogic{}
	logic.gameType = gameType
	logic.cardList = []byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, //方块 A - K
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, //梅花 A - K
		0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D, //红桃 A - K
		0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D, //黑桃 A - K
	}

	return logic
}

func GetCardType(cards []byte) int {
	return getCardType(cards)
}

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
	if (getCardValue(cards[0]) == 2 && getCardValue(cards[1]) == 3 && getCardValue(cards[2]) == 5) && (getCardColor(cards[0]) != getCardColor(cards[1]) && getCardColor(cards[0]) != getCardColor(cards[2]) && getCardColor(cards[1]) != getCardColor(cards[2])) {
		return CARD_TYPE_SPECIAL
	}

	return CARD_TYPE_SINGLE
}

func getCardValue(card byte) byte {
	return card & MASK_VALUE
}

func getCardColor(card byte) byte {
	return card & MASK_COLOR
}

func getCardLogicValue(card byte) byte {
	v := getCardValue(card)
	if v == 1 {
		return v + 13
	}
	return v
}

func (logic *GameLogic) ShuffleCards2() {
	sort.Sort(ByteSlice(logic.cardList))
}

func (logic *GameLogic) ShuffleCards(pos []int, robotPositions []int) {
	sort.Sort(ByteSlice(logic.cardList))

	logic.pos = pos
	logic.sortedPos = []int{}
	logic.sorted = false

	// 按概率抽牌
	for i := 0; i < len(pos); i++ {
		//		t := config.GetCardConfigManager().GetRandCardType(logic.gameType)

		t := CARD_TYPE_JIN_HUA

		if t == CARD_TYPE_SINGLE {
			logic.changeSingle(i)
		} else if t == CARD_TYPE_DOUBLE {
			logic.changeDouble(i)
		} else if t == CARD_TYPE_SHUN_ZI {
			logic.changeShunZi(i)
		} else if t == CARD_TYPE_JIN_HUA {
			logic.changeJinHua(i)
		} else if t == CARD_TYPE_SHUN_JIN {
			logic.changeShunJin(i)
		} else if t == CARD_TYPE_BAO_ZI {
			logic.changeBaoZhi(i)
		}
		logic.sortedPos = append(logic.sortedPos, i)
	}

	cardTypeComp := NewCardTypeComp(logic)
	cardTypeComp.pos = logic.sortedPos
	sort.Sort(cardTypeComp)
	logic.sortedPos = cardTypeComp.pos
	logic.sorted = true

	// 没有机器人
	if len(robotPositions) <= 0 {
		return
	}

	// 机器人拿牌最大
	for _, v := range robotPositions {
		if logic.sortedPos[0] == v {
			return
		}
	}

	maxCards := []byte{logic.cardList[logic.sortedPos[0]*3], logic.cardList[logic.sortedPos[0]*3+1], logic.cardList[logic.sortedPos[0]*3+2]}
	cardType := GetCardType(maxCards)
	robotPos := robotPositions[rand.Int()%len(robotPositions)]

	if cardType == CARD_TYPE_BAO_ZI {
		if rand.Int()%100 < 10 {
			// 替换更大豹子
			logic.changeGreaterBaoZhi(maxCards, robotPos)
		}
	} else if cardType == CARD_TYPE_SHUN_JIN {
		if rand.Int()%100 < 10 {
			// 替换豹子
			logic.changeBaoZhi(robotPos)
		}
	} else if cardType == CARD_TYPE_JIN_HUA {
		if rand.Int()%100 < 10 {
			// 优先替换同花顺,豹子次之
			if !logic.changeShunJin(robotPos) {
				logic.changeBaoZhi(robotPos)
			}
		}
	}
}

func (logic *GameLogic) getSortedPos(i int) int {
	if !logic.sorted {
		return i
	}

	if len(logic.sortedPos) > 0 && len(logic.pos) > 0 && len(logic.sortedPos) == len(logic.pos) {
		for idx, v := range logic.pos {
			if v == i {
				return logic.sortedPos[idx]
			}
		}
	}

	return i
}

func (logic *GameLogic) GetCards(i int) []byte {
	pos := logic.getSortedPos(i)

	return []byte{logic.cardList[pos*3], logic.cardList[pos*3+1], logic.cardList[pos*3+2]}
}

func (logic *GameLogic) GetCardsInt32(i int) []int32 {
	cards := logic.GetCards(i)
	sort.Sort(CardComp(cards))
	return []int32{int32(cards[0]), int32(cards[1]), int32(cards[2])}
}

func (logic *GameLogic) getCards(i int) []byte {
	return []byte{logic.cardList[i*3], logic.cardList[i*3+1], logic.cardList[i*3+2]}
}

func (logic *GameLogic) CompareCardsByPos(pos1, pos2 int) bool {
	return logic.CompareCards(logic.GetCards(pos1), logic.GetCards(pos2))
}

func (logic *GameLogic) CompareCards2(firstCards []int32, nextCards []int32) bool {
	if len(firstCards) != 3 {
		return true
	}
	if len(nextCards) != 3 {
		return true
	}
	fCards := []byte{byte(firstCards[0]), byte(firstCards[1]), byte(firstCards[2])}
	nCards := []byte{byte(nextCards[0]), byte(nextCards[1]), byte(nextCards[2])}

	return logic.CompareCards(fCards, nCards)
}

// 相同牌，先开者输(即firstCards)
func (logic *GameLogic) CompareCards(firstCards []byte, nextCards []byte) bool {
	firstType := getCardType(firstCards)
	nextType := getCardType(nextCards)

	// 特殊情况
	if firstType == CARD_TYPE_SPECIAL && nextType == CARD_TYPE_BAO_ZI {
		return true
	}

	if firstType == CARD_TYPE_BAO_ZI && nextType == CARD_TYPE_SPECIAL {
		return false
	}

	if firstType == CARD_TYPE_SPECIAL {
		firstType = CARD_TYPE_SINGLE
	}

	if nextType == CARD_TYPE_SPECIAL {
		nextType = CARD_TYPE_SINGLE
	}

	// 不同类型
	if firstType != nextType {
		return firstType > nextType
	}

	// 豹子，单牌，金花
	if firstType == CARD_TYPE_BAO_ZI || firstType == CARD_TYPE_SINGLE || firstType == CARD_TYPE_JIN_HUA {
		for i := 2; i >= 0; i-- {
			if getCardLogicValue(firstCards[i]) != getCardLogicValue(nextCards[i]) {
				return getCardLogicValue(firstCards[i]) > getCardLogicValue(nextCards[i])
			}
		}
		// 先开者输
		return false
	}

	// 顺子，顺金
	if firstType == CARD_TYPE_SHUN_ZI || firstType == CARD_TYPE_SHUN_JIN {
		firstV := getCardLogicValue(firstCards[2])
		nextV := getCardLogicValue(nextCards[2])

		if firstV == 14 && getCardLogicValue(firstCards[1]) == 3 {
			firstV = 3
		}
		if nextV == 14 && getCardLogicValue(nextCards[1]) == 3 {
			nextV = 3
		}

		if firstV != nextV {
			return firstV > nextV
		}

		// 先开者输
		return false
	}

	// 对子
	if firstType == CARD_TYPE_DOUBLE {
		var firstDouble byte = 0
		var firstSingle byte = 0
		var nextDouble byte = 0
		var nextSingle byte = 0

		if getCardLogicValue(firstCards[0]) == getCardLogicValue(firstCards[1]) {
			firstDouble = getCardLogicValue(firstCards[0])
			firstSingle = getCardLogicValue(firstCards[2])
		} else {
			firstDouble = getCardLogicValue(firstCards[2])
			firstSingle = getCardLogicValue(firstCards[0])
		}

		if getCardLogicValue(nextCards[0]) == getCardLogicValue(nextCards[1]) {
			nextDouble = getCardLogicValue(nextCards[0])
			nextSingle = getCardLogicValue(nextCards[2])
		} else {
			nextDouble = getCardLogicValue(nextCards[2])
			nextSingle = getCardLogicValue(nextCards[0])
		}

		if firstDouble != nextDouble {
			return firstDouble > nextDouble
		}
		if firstSingle != nextSingle {
			return firstSingle > nextSingle
		}

		// 先开者输
		return false
	}

	return true
}

func (logic *GameLogic) test() {
	sort.Sort(ByteSlice(logic.cardList))
	fmt.Println("cardlist:", logic.cardList)
}

func PrintCards(s string, cards []byte) {
	fmt.Print("==>s:", s)
	for _, v := range cards {
		fmt.Printf("%x ", v)
	}
	PrintCardType(getCardType(cards))
}

func (logic *GameLogic) getRandPositions() []int {
	randPositions := []int{}
	for i := 0; i < 17; i++ {
		exist := false
		for j := 0; j < len(logic.sortedPos); j++ {
			if logic.sortedPos[j] == i {
				exist = true
				break
			}
		}
		if !exist {
			randPositions = append(randPositions, i)
		}
	}

	return randPositions
}

func (logic *GameLogic) ReplaceCard(i, card int) int {
	startPos := logic.getSortedPos(i)

	randPositions := logic.getRandPositions()
	if len(randPositions) <= 0 {
		return 0
	}

	r := randPositions[rand.Int()%len(randPositions)]
	r = r*3 + rand.Int()%3

	for pos := startPos * 3; pos <= startPos*3+2; pos++ {
		if int(logic.cardList[pos]) == card {
			logic.cardList[pos], logic.cardList[r] = logic.cardList[r], logic.cardList[pos]
			return int(logic.cardList[pos])
		}
	}

	return 0
}

// 替换更大的豹子
func (logic *GameLogic) changeGreaterBaoZhi(cards []byte, pos int) {
	pos = logic.getSortedPos(pos)

	if getCardValue(cards[0]) == 1 {
		// 已经是最大的豹子
		return
	}

	randPositions := logic.getRandPositions()
	if len(randPositions) <= 0 {
		return
	}

	cardCounts := make(map[int]int)
	for i := 0; i < len(logic.pos); i++ {
		cs := logic.GetCards(i)
		for _, v := range cs {
			cardCounts[int(getCardLogicValue(v))]++
		}
	}

	values := []int{}
	for i := 2; i <= 14; i++ {
		if cardCounts[i] <= 1 && i > int(getCardLogicValue(cards[0])) {
			values = append(values, i)
		}
	}

	if len(values) <= 0 {
		return
	}

	randCard := values[rand.Int()%len(values)]

	changePositions := []int{}

	for _, p := range randPositions {
		cs := logic.getCards(p)
		for i, v := range cs {
			if int(getCardLogicValue(v)) == randCard {
				changePositions = append(changePositions, p*3+i)
			}
		}
	}

	if int(getCardLogicValue(logic.cardList[51])) == randCard {
		changePositions = append(changePositions, 51)
	}

	if len(changePositions) >= 3 {
		logic.cardList[pos*3], logic.cardList[changePositions[0]] = logic.cardList[changePositions[0]], logic.cardList[pos*3]
		logic.cardList[pos*3+1], logic.cardList[changePositions[1]] = logic.cardList[changePositions[1]], logic.cardList[pos*3+1]
		logic.cardList[pos*3+2], logic.cardList[changePositions[2]] = logic.cardList[changePositions[2]], logic.cardList[pos*3+2]
	}
}

// 替换豹子
func (logic *GameLogic) changeBaoZhi(pos int) {
	pos = logic.getSortedPos(pos)

	changePositions := []int{}
	for i := len(logic.sortedPos) * 3; i < 52; i++ {
		for j := i + 1; j < 52; j++ {
			for k := j + 1; k < 52; k++ {
				if getCardType([]byte{logic.cardList[i], logic.cardList[j], logic.cardList[k]}) == CARD_TYPE_BAO_ZI {
					changePositions = []int{i, j, k}
					goto end
				}
			}
		}
	}

end:
	sort.Ints(changePositions)

	if len(changePositions) >= 3 {
		logic.cardList[pos*3], logic.cardList[changePositions[0]] = logic.cardList[changePositions[0]], logic.cardList[pos*3]
		logic.cardList[pos*3+1], logic.cardList[changePositions[1]] = logic.cardList[changePositions[1]], logic.cardList[pos*3+1]
		logic.cardList[pos*3+2], logic.cardList[changePositions[2]] = logic.cardList[changePositions[2]], logic.cardList[pos*3+2]
	}
}

// 替换顺金
func (logic *GameLogic) changeShunJin(pos int) bool {
	pos = logic.getSortedPos(pos)

	changePositions := []int{}
	for i := len(logic.sortedPos) * 3; i < 52; i++ {
		for j := i + 1; j < 52; j++ {
			for k := j + 1; k < 52; k++ {
				if getCardType([]byte{logic.cardList[i], logic.cardList[j], logic.cardList[k]}) == CARD_TYPE_SHUN_JIN {
					changePositions = []int{i, j, k}
					goto end
				}
			}
		}
	}

end:
	sort.Ints(changePositions)

	if len(changePositions) >= 3 {
		logic.cardList[pos*3], logic.cardList[changePositions[0]] = logic.cardList[changePositions[0]], logic.cardList[pos*3]
		logic.cardList[pos*3+1], logic.cardList[changePositions[1]] = logic.cardList[changePositions[1]], logic.cardList[pos*3+1]
		logic.cardList[pos*3+2], logic.cardList[changePositions[2]] = logic.cardList[changePositions[2]], logic.cardList[pos*3+2]
		return true
	}
	return false
}

// 替换顺子
func (logic *GameLogic) changeShunZi(pos int) {
	pos = logic.getSortedPos(pos)

	changePositions := []int{}
	for i := len(logic.sortedPos) * 3; i < 52; i++ {
		for j := i + 1; j < 52; j++ {
			for k := j + 1; k < 52; k++ {
				if getCardType([]byte{logic.cardList[i], logic.cardList[j], logic.cardList[k]}) == CARD_TYPE_SHUN_ZI {
					changePositions = []int{i, j, k}
					goto end
				}
			}
		}
	}

end:
	sort.Ints(changePositions)

	if len(changePositions) >= 3 {
		logic.cardList[pos*3], logic.cardList[changePositions[0]] = logic.cardList[changePositions[0]], logic.cardList[pos*3]
		logic.cardList[pos*3+1], logic.cardList[changePositions[1]] = logic.cardList[changePositions[1]], logic.cardList[pos*3+1]
		logic.cardList[pos*3+2], logic.cardList[changePositions[2]] = logic.cardList[changePositions[2]], logic.cardList[pos*3+2]
	}
}

// 替换金花
func (logic *GameLogic) changeJinHua(pos int) {
	pos = logic.getSortedPos(pos)

	changePositions := []int{}
	for i := len(logic.sortedPos) * 3; i < 52; i++ {
		for j := i + 1; j < 52; j++ {
			for k := j + 1; k < 52; k++ {
				if getCardColor(logic.cardList[i]) == getCardColor(logic.cardList[j]) && getCardColor(logic.cardList[j]) == getCardColor(logic.cardList[k]) {
					if getCardType([]byte{logic.cardList[i], logic.cardList[j], logic.cardList[k]}) == CARD_TYPE_SHUN_JIN {
						continue
					}
					changePositions = []int{i, j, k}
					goto end
				}
			}
		}
	}

end:
	sort.Ints(changePositions)

	if len(changePositions) >= 3 {
		logic.cardList[pos*3], logic.cardList[changePositions[0]] = logic.cardList[changePositions[0]], logic.cardList[pos*3]
		logic.cardList[pos*3+1], logic.cardList[changePositions[1]] = logic.cardList[changePositions[1]], logic.cardList[pos*3+1]
		logic.cardList[pos*3+2], logic.cardList[changePositions[2]] = logic.cardList[changePositions[2]], logic.cardList[pos*3+2]
	}
}

// 替换对子
func (logic *GameLogic) changeDouble(pos int) {
	pos = logic.getSortedPos(pos)

	changePositions := []int{}
	for i := len(logic.sortedPos) * 3; i < 52; i++ {
		for j := i + 1; j < 52; j++ {
			for k := j + 1; k < 52; k++ {
				if getCardLogicValue(logic.cardList[i]) == getCardLogicValue(logic.cardList[j]) && getCardLogicValue(logic.cardList[j]) != getCardLogicValue(logic.cardList[k]) {
					changePositions = []int{i, j, k}
					goto end
				}
			}
		}
	}
end:
	sort.Ints(changePositions)

	if len(changePositions) >= 3 {
		logic.cardList[pos*3], logic.cardList[changePositions[0]] = logic.cardList[changePositions[0]], logic.cardList[pos*3]
		logic.cardList[pos*3+1], logic.cardList[changePositions[1]] = logic.cardList[changePositions[1]], logic.cardList[pos*3+1]
		logic.cardList[pos*3+2], logic.cardList[changePositions[2]] = logic.cardList[changePositions[2]], logic.cardList[pos*3+2]
	}
}

// 替换单牌
func (logic *GameLogic) changeSingle(pos int) {
	pos = logic.getSortedPos(pos)

	changePositions := []int{}
	for i := len(logic.sortedPos) * 3; i < 52; i++ {
		for j := i + 1; j < 52; j++ {
			for k := j + 1; k < 52; k++ {
				if getCardType([]byte{logic.cardList[i], logic.cardList[j], logic.cardList[k]}) == CARD_TYPE_SINGLE {
					changePositions = []int{i, j, k}
					goto end
				}
			}
		}
	}

end:
	sort.Ints(changePositions)

	if len(changePositions) >= 3 {
		logic.cardList[pos*3], logic.cardList[changePositions[0]] = logic.cardList[changePositions[0]], logic.cardList[pos*3]
		logic.cardList[pos*3+1], logic.cardList[changePositions[1]] = logic.cardList[changePositions[1]], logic.cardList[pos*3+1]
		logic.cardList[pos*3+2], logic.cardList[changePositions[2]] = logic.cardList[changePositions[2]], logic.cardList[pos*3+2]
	}
}

func (logic *GameLogic) checkValid() {
	m := make(map[byte]int)
	for _, v := range logic.cardList {
		m[v]++
	}
	for k, v := range m {
		if v >= 2 {
			fmt.Println("k = ", k, " v = ", v, " cards:", logic.cardList)
			panic(errors.New("非法"))
		}
	}
}

func main() {
	//	flag.Parse()
	//	config.GetCardConfigManager().Init()
	t := time.Now()
	logic := NewGameLogic(int(util.GameType_Common_Level_3))

	if logic.CompareCards([]byte{0x23, 0x04, 0x12}, []byte{0x1a, 0x1b, 0x11}) {
		fmt.Println("大")
	} else {
		fmt.Println("小")
	}

	m := make(map[int]int)

	pos := []int{0, 3, 4}
	for j := 0; j < 10000; j++ {
		logic.ShuffleCards(pos, []int{})
		logic.checkValid()

		for _, i := range pos {
			m[getCardType(logic.GetCards(i))]++
			if getCardType(logic.GetCards(i)) != CARD_TYPE_JIN_HUA {
				panic(errors.New("不是对子"))
			}
			//			PrintCards(fmt.Sprintf("==>%v  ", i), logic.GetCards(i))
		}

		//		if GetCardType(logic.GetCards(0)) != CARD_TYPE_BAO_ZI {
		//			PrintCards("错误类型", logic.GetCards(0))
		//			return
		//		}
		//		if GetCardType(logic.GetCards(1)) != CARD_TYPE_SHUN_JIN {
		//			PrintCards("错误类型", logic.GetCards(1))
		//			return
		//		}
		//		if GetCardType(logic.GetCards(2)) != CARD_TYPE_SHUN_ZI {
		//			PrintCards("错误类型", logic.GetCards(2))
		//			return
		//		}
		//		if GetCardType(logic.GetCards(3)) != CARD_TYPE_DOUBLE {
		//			PrintCards("错误类型", logic.GetCards(3))
		//			return
		//		}
		//		if GetCardType(logic.GetCards(4)) != CARD_TYPE_JIN_HUA {
		//			PrintCards("错误类型", logic.GetCards(4))
		//			return
		//		}

		//		PrintCards("", logic.GetCards(i))
		//		PrintCardType(GetCardType(logic.GetCards(i)))
	}
	fmt.Println(m)
	fmt.Println(time.Since(t))

}
