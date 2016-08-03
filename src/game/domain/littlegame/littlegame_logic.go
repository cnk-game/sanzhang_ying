package littlegame

import (
	"fmt"
	domainUser "game/domain/user"
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"math/rand"
	"pb"
	"sort"
	"time"
	"util"
)

const (
	MASK_COLOR = 0xF0 // 花色掩码
	MASK_VALUE = 0x0F // 数值掩码
	LittleLogC = "littlegame_log"
)

func PrintCardType(cardType int) string {
	switch cardType {
	case CARD_TYPE_SINGLE:
		return fmt.Sprintf("单牌类型")
	case CARD_TYPE_DOUBLE:
		return fmt.Sprintf("对子类型")
	case CARD_TYPE_SHUN_ZI:
		return fmt.Sprintf("顺子类型")
	case CARD_TYPE_JIN_HUA:
		return fmt.Sprintf("金花类型")
	case CARD_TYPE_SHUN_JIN:
		return fmt.Sprintf("顺金类型")
	case CARD_TYPE_BAO_ZI:
		return fmt.Sprintf("豹子类型")
	case CARD_TYPE_BAO_A:
		return fmt.Sprintf("豹子A类型")
	case CARD_TYPE_SEPCIAL:
		return fmt.Sprintf("地龙类型")
	}
	return "未知"
}

type LittleGameLogic struct {
	cardList []byte
	gameType int
}

type LittleGameLog struct {
	Gametype   int       `bson:"gametype"`
	Chip       int       `bson:"chip"`
	Money      int       `bson:"money"`
	Userid     string    `bson:"userid"`
	Paitype    int       `bson:"paitype"`
	Createtime time.Time `bson:"createtime"`
}

var gameLogicM *LittleGameLogic

func init() {
	rand.Seed(time.Now().UnixNano())

	gameLogicM = NewLittleGameLogic(1)

}

func GetLigameLogicManager() *LittleGameLogic {
	return gameLogicM
}

func NewLittleGameLogic(gameType int) *LittleGameLogic {
	logic := &LittleGameLogic{}
	logic.gameType = gameType
	logic.cardList = []byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, //方块 A - K
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, //梅花 A - K
		0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D, //红桃 A - K
		0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D, //黑桃 A - K
	}

	return logic
}

type CardComp []byte

func (p CardComp) Len() int { return len(p) }

func (p CardComp) Less(i, j int) bool {
	return getCardLogicValue(p[i]) < getCardLogicValue(p[j])
}

func (p CardComp) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func GetCardType(cards []byte) int {
	return getCardType(cards)
}

func getCurrentDate() string {
	now := time.Now()

	year, mon, day := now.Date()

	//zone, _ := now.Zone()
	return fmt.Sprintf("%d%02d%02d", year, mon, day)

}

func SaveLog(chip int, Gametype int, Money int, Paitype int, Userid string) error {

	log := &LittleGameLog{}
	log.Chip = chip
	log.Gametype = Gametype
	log.Money = Money
	log.Paitype = Paitype
	log.Userid = Userid
	now := time.Now()
	log.Createtime = now

	cur_C := LittleLogC + "_" + getCurrentDate() //strconv.Itoa(int(now.Year())) + strconv.Itoa(int(now.Month())) + strconv.Itoa(int(now.Day()))

	return util.WithLogCollection(cur_C, func(c *mgo.Collection) error {
		return c.Insert(log)
	})
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

	//豹子A
	isBaoZiA := false
	if getCardLogicValue(cards[0]) == 14 && getCardLogicValue(cards[1]) == 14 && getCardLogicValue(cards[2]) == 14 {
		isBaoZiA = true
	}

	if isDouble {
		if isBaoZi {
			if isBaoZiA {
				return CARD_TYPE_BAO_A
			}
			return CARD_TYPE_BAO_ZI
		}
		return CARD_TYPE_DOUBLE
	}

	// 地龙特殊235
	if getCardValue(cards[0]) == 2 && getCardValue(cards[1]) == 3 && getCardValue(cards[2]) == 5 {
		return CARD_TYPE_SEPCIAL
	}

	return CARD_TYPE_SINGLE
}

func getCardValue(card byte) byte {
	return card & MASK_VALUE
}

func getCardColor(card byte) byte {
	return card & MASK_COLOR
}

func getCardColorString(card int32) string {

	if card&MASK_COLOR == 0x00 {
		return "方块"
	} else if card&MASK_COLOR == 0x10 {
		return "梅花"
	} else if card&MASK_COLOR == 0x20 {
		return "红桃"
	} else if card&MASK_COLOR == 0x30 {
		return "黑桃"
	}
	return "方块"
}

func getCardLogicValue(card byte) byte {
	v := getCardValue(card)
	if v == 1 {
		return v + 13
	}
	return v
}

func getCardLogicValueTest(card int32) byte {
	v := getCardValue(byte(card))

	return v
}

func getTypeName(gameType int) string {

	if gameType == 1 {
		return "菜鸟场"
	} else if gameType == 4 {
		return "中级场"
	} else if gameType == 2 {
		return "高手场"
	} else if gameType == 2 {
		return "精英场"
	}
	return "菜鸟场"
}

func (logic *LittleGameLogic) ShuffleCards(gameType int, chip int, userid string, nickname string) (int, int, int, []int32) {

	glog.Info("littlegame_ShuffleCards in.")

	f, ok := domainUser.GetUserFortuneManager().GetUserFortune(userid)

	if !ok {
		glog.Info("===>小游戏开始检测，查询用户财富信息失败userId:", userid)
		return 1, 0, 0, []int32{}
	}

	if f.Gold < int64(chip) {
		glog.Info("littlegame_ShuffleCards金币不足，无法下注 f.Gold=,chip=", f.Gold, chip)
		return 1, 0, 0, []int32{}
	}

	// 按概率抽牌
	t := GetCardConfigManager().GetRandCardType(gameType)
	glog.Info("littlegame_ShuffleCards in. 牌型=", t)

	sitepos := 0
	if t == CARD_TYPE_SINGLE {
		sitepos = logic.changeSingle()
		//单牌的话用户输掉下注的金额
		// 扣费
		curGold, _, ok := domainUser.GetUserFortuneManager().ConsumeGold(userid, int64(chip), true, "小游戏单牌扣金币")
		if !ok {
			glog.Info("==>扣款失败userId:", userid, " betGold:", chip)
			return 2, 0, 0, []int32{}
		}
		glog.Info("==>扣款成功userId:", userid, " betGold:", chip, " curGold:", curGold)

		glog.Info("单牌", getCardColorString(logic.GetCardsInt32(sitepos)[0]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[0]), " ", getCardColorString(logic.GetCardsInt32(sitepos)[1]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[1]), " ", getCardColorString(logic.GetCardsInt32(sitepos)[2]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[2]))
		SaveLog(chip, gameType, chip, t, userid)

		return 0, CARD_TYPE_SINGLE, 0, logic.GetCardsInt32(sitepos) //单牌赢的金币返回0

	} else if t == CARD_TYPE_DOUBLE {
		sitepos = logic.changeDouble()
		//对子加相应倍数的金币
		domainUser.GetUserFortuneManager().EarnGold(userid, int64(chip*(GetCardConfigManager().multiple.Double-1)), "小游戏对子加金币")

		glog.Info("对子", getCardColorString(logic.GetCardsInt32(sitepos)[0]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[0]), " ", getCardColorString(logic.GetCardsInt32(sitepos)[1]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[1]), " ", getCardColorString(logic.GetCardsInt32(sitepos)[2]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[2]))

		SaveLog(chip, gameType, 0-chip*(GetCardConfigManager().multiple.Double-1), t, userid)
		return 0, CARD_TYPE_DOUBLE, chip * GetCardConfigManager().multiple.Double, logic.GetCardsInt32(sitepos)
	} else if t == CARD_TYPE_SHUN_ZI {
		sitepos = logic.changeShunZi()
		//顺子加相应倍数的金币
		domainUser.GetUserFortuneManager().EarnGold(userid, int64(chip*(GetCardConfigManager().multiple.ShunZi-1)), "小游戏顺子加金币")

		glog.Info("顺子", getCardColorString(logic.GetCardsInt32(sitepos)[0]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[0]), " ", getCardColorString(logic.GetCardsInt32(sitepos)[1]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[1]), " ", getCardColorString(logic.GetCardsInt32(sitepos)[2]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[2]))

		SaveLog(chip, gameType, 0-chip*(GetCardConfigManager().multiple.ShunZi-1), t, userid)
		return 0, CARD_TYPE_SHUN_ZI, chip * GetCardConfigManager().multiple.ShunZi, logic.GetCardsInt32(sitepos)
	} else if t == CARD_TYPE_JIN_HUA {
		sitepos = logic.changeJinHua()
		//金花加相应倍数的金币
		domainUser.GetUserFortuneManager().EarnGold(userid, int64(chip*(GetCardConfigManager().multiple.ShunZi-1)), "小游戏金花加金币")

		glog.Info("金花", getCardColorString(logic.GetCardsInt32(sitepos)[0]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[0]), " ", getCardColorString(logic.GetCardsInt32(sitepos)[1]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[1]), " ", getCardColorString(logic.GetCardsInt32(sitepos)[2]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[2]))

		SaveLog(chip, gameType, 0-chip*(GetCardConfigManager().multiple.ShunZi-1), t, userid)
		return 0, CARD_TYPE_JIN_HUA, chip * GetCardConfigManager().multiple.JinHua, logic.GetCardsInt32(sitepos)
	} else if t == CARD_TYPE_SHUN_JIN {
		sitepos = logic.changeShunJin()
		//顺金加相应倍数的金币
		domainUser.GetUserFortuneManager().EarnGold(userid, int64(chip*(GetCardConfigManager().multiple.ShunJin-1)), "小游戏顺金加金币")

		glog.Info("顺金", getCardColorString(logic.GetCardsInt32(sitepos)[0]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[0]), " ", getCardColorString(logic.GetCardsInt32(sitepos)[1]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[1]), " ", getCardColorString(logic.GetCardsInt32(sitepos)[2]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[2]))

		//发通告   XXX在XXX场玩翻翻乐中了XXX获得了XXX金币
		domainUser.GetPlayerManager().BroadcastClientMsg(int32(pb.MessageId_CHAT), util.BuildSysBugle(fmt.Sprintf("恭喜%v在%v玩翻翻乐中了同花顺获得了%v金币！", nickname, getTypeName(gameType), chip*GetCardConfigManager().multiple.ShunJin)))

		SaveLog(chip, gameType, 0-chip*(GetCardConfigManager().multiple.ShunJin-1), t, userid)
		return 0, CARD_TYPE_SHUN_JIN, chip * GetCardConfigManager().multiple.ShunJin, logic.GetCardsInt32(sitepos)
	} else if t == CARD_TYPE_BAO_ZI {
		sitepos = logic.changeBaoZhi()
		//豹子加相应倍数的金币
		domainUser.GetUserFortuneManager().EarnGold(userid, int64(chip*(GetCardConfigManager().multiple.BaoZi-1)), "小游戏豹子加金币")

		glog.Info("豹子", getCardColorString(logic.GetCardsInt32(sitepos)[0]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[0]), " ", getCardColorString(logic.GetCardsInt32(sitepos)[1]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[1]), " ", getCardColorString(logic.GetCardsInt32(sitepos)[2]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[2]))

		//发通告   XXX在XXX场玩翻翻乐中了XXX获得了XXX金币
		domainUser.GetPlayerManager().BroadcastClientMsg(int32(pb.MessageId_CHAT), util.BuildSysBugle(fmt.Sprintf("恭喜%v在%v玩翻翻乐中了豹子获得了%v金币！", nickname, getTypeName(gameType), chip*GetCardConfigManager().multiple.BaoZi)))

		SaveLog(chip, gameType, 0-chip*(GetCardConfigManager().multiple.BaoZi-1), t, userid)
		return 0, CARD_TYPE_BAO_ZI, chip * GetCardConfigManager().multiple.BaoZi, logic.GetCardsInt32(sitepos)
	} else if t == CARD_TYPE_BAO_A {
		sitepos = logic.changeBaoZhiA()
		//豹子A加相应倍数的金币
		domainUser.GetUserFortuneManager().EarnGold(userid, int64(chip*(GetCardConfigManager().multiple.BaoA-1)), "小游戏豹子A加金币")

		glog.Info("豹子A", getCardColorString(logic.GetCardsInt32(sitepos)[0]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[0]), " ", getCardColorString(logic.GetCardsInt32(sitepos)[1]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[1]), " ", getCardColorString(logic.GetCardsInt32(sitepos)[2]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[2]))

		//发通告   XXX在XXX场玩翻翻乐中了XXX获得了XXX金币
		domainUser.GetPlayerManager().BroadcastClientMsg(int32(pb.MessageId_CHAT), util.BuildSysBugle(fmt.Sprintf("恭喜%v在%v玩翻翻乐中了豹子A获得了%v金币！", nickname, getTypeName(gameType), chip*GetCardConfigManager().multiple.BaoA)))

		SaveLog(chip, gameType, 0-chip*(GetCardConfigManager().multiple.BaoA-1), t, userid)
		return 0, CARD_TYPE_BAO_A, chip * GetCardConfigManager().multiple.BaoA, logic.GetCardsInt32(sitepos)
	} else if t == CARD_TYPE_SEPCIAL {
		sitepos = logic.changeSpecial()
		//地龙加相应倍数的金币
		domainUser.GetUserFortuneManager().EarnGold(userid, int64(chip*(GetCardConfigManager().multiple.Special-1)), "小游戏地龙加金币")

		glog.Info("地龙", getCardColorString(logic.GetCardsInt32(sitepos)[0]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[0]), " ", getCardColorString(logic.GetCardsInt32(sitepos)[1]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[1]), " ", getCardColorString(logic.GetCardsInt32(sitepos)[2]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[2]))

		//发通告   XXX在XXX场玩翻翻乐中了XXX获得了XXX金币
		domainUser.GetPlayerManager().BroadcastClientMsg(int32(pb.MessageId_CHAT), util.BuildSysBugle(fmt.Sprintf("恭喜%v在%v玩翻翻乐中了地龙获得了%v金币！", nickname, getTypeName(gameType), chip*GetCardConfigManager().multiple.Special)))

		SaveLog(chip, gameType, 0-chip*(GetCardConfigManager().multiple.Special-1), t, userid)
		return 0, CARD_TYPE_SEPCIAL, chip * GetCardConfigManager().multiple.Special, logic.GetCardsInt32(sitepos)
	}
	return 0, CARD_TYPE_SINGLE, 0, logic.GetCardsInt32(sitepos) //单牌赢的金币返回0
}

func (logic *LittleGameLogic) GetCards(pos int) []byte {

	return []byte{logic.cardList[pos], logic.cardList[pos+1], logic.cardList[pos+2]}
}

func (logic *LittleGameLogic) GetCardsInt32(pos int) []int32 {
	cards := logic.GetCards(pos)

	return []int32{int32(cards[0]), int32(cards[1]), int32(cards[2])}
}

func (logic *LittleGameLogic) getCards() []byte {
	return []byte{logic.cardList[0], logic.cardList[1], logic.cardList[2]}
}

func PrintCards(s string, cards []byte) {
	fmt.Print("==>s:", s)
	for _, v := range cards {
		fmt.Printf("%x ", v)
	}
	PrintCardType(getCardType(cards))
}

// 替换地龙
func (logic *LittleGameLogic) changeSpecial() int {
begin:
	changePositions := []int{}
	pos1 := rand.Int() % 10
	for i := pos1; i < 52; i++ {
		for j := i + 1; j < 52; j++ {
			for k := j + 1; k < 52; k++ {
				if getCardType([]byte{logic.cardList[i], logic.cardList[j], logic.cardList[k]}) == CARD_TYPE_SEPCIAL {
					changePositions = []int{i, j, k}
					goto end
				}
			}
		}
	}

end:
	sort.Ints(changePositions)

	if len(changePositions) >= 3 {
		logic.cardList[pos1*3], logic.cardList[changePositions[0]] = logic.cardList[changePositions[0]], logic.cardList[pos1*3]
		logic.cardList[pos1*3+1], logic.cardList[changePositions[1]] = logic.cardList[changePositions[1]], logic.cardList[pos1*3+1]
		logic.cardList[pos1*3+2], logic.cardList[changePositions[2]] = logic.cardList[changePositions[2]], logic.cardList[pos1*3+2]
	} else {
		goto begin
	}
	return pos1
}

// 替换豹子A
func (logic *LittleGameLogic) changeBaoZhiA() int {
begin:
	changePositions := []int{}
	pos1 := 0 //rand.Int() % 5
	for i := 1; i < 52; i++ {
		for j := i + 1; j < 52; j++ {
			for k := j + 1; k < 52; k++ {
				if getCardType([]byte{logic.cardList[i], logic.cardList[j], logic.cardList[k]}) == CARD_TYPE_BAO_A {
					changePositions = []int{i, j, k}
					goto end
				}
			}
		}
	}

end:
	sort.Ints(changePositions)

	if len(changePositions) >= 3 {
		logic.cardList[pos1*3], logic.cardList[changePositions[0]] = logic.cardList[changePositions[0]], logic.cardList[pos1*3]
		logic.cardList[pos1*3+1], logic.cardList[changePositions[1]] = logic.cardList[changePositions[1]], logic.cardList[pos1*3+1]
		logic.cardList[pos1*3+2], logic.cardList[changePositions[2]] = logic.cardList[changePositions[2]], logic.cardList[pos1*3+2]
	} else {
		goto begin
	}
	return pos1
}

// 替换豹子
func (logic *LittleGameLogic) changeBaoZhi() int {
begin:
	changePositions := []int{}
	pos1 := rand.Int() % 5
	for i := pos1; i < 52; i++ {
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
		logic.cardList[pos1], logic.cardList[changePositions[0]] = logic.cardList[changePositions[0]], logic.cardList[pos1]
		logic.cardList[pos1+1], logic.cardList[changePositions[1]] = logic.cardList[changePositions[1]], logic.cardList[pos1+1]
		logic.cardList[pos1+2], logic.cardList[changePositions[2]] = logic.cardList[changePositions[2]], logic.cardList[pos1+2]
	} else {
		goto begin
	}
	return pos1
}

// 替换顺金
func (logic *LittleGameLogic) changeShunJin() int {
begin:
	changePositions := []int{}
	pos1 := rand.Int() % 3
	for i := pos1; i < 52; i++ {
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
		logic.cardList[pos1], logic.cardList[changePositions[0]] = logic.cardList[changePositions[0]], logic.cardList[pos1]
		logic.cardList[pos1+1], logic.cardList[changePositions[1]] = logic.cardList[changePositions[1]], logic.cardList[pos1+1]
		logic.cardList[pos1+2], logic.cardList[changePositions[2]] = logic.cardList[changePositions[2]], logic.cardList[pos1+2]

	} else {
		goto begin
	}
	return pos1
}

// 替换顺子
func (logic *LittleGameLogic) changeShunZi() int {
begin:
	changePositions := []int{}
	pos1 := rand.Int() % 3
	for i := pos1; i < 52; i++ {
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
		logic.cardList[pos1], logic.cardList[changePositions[0]] = logic.cardList[changePositions[0]], logic.cardList[pos1]
		logic.cardList[pos1+1], logic.cardList[changePositions[1]] = logic.cardList[changePositions[1]], logic.cardList[pos1+1]
		logic.cardList[pos1+2], logic.cardList[changePositions[2]] = logic.cardList[changePositions[2]], logic.cardList[pos1+2]
	} else {
		goto begin
	}
	return pos1
}

// 替换金花
func (logic *LittleGameLogic) changeJinHua() int {
begin:
	changePositions := []int{}
	pos1 := rand.Int() % 3
	for i := pos1; i < 52; i++ {
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
		logic.cardList[pos1], logic.cardList[changePositions[0]] = logic.cardList[changePositions[0]], logic.cardList[pos1]
		logic.cardList[pos1+1], logic.cardList[changePositions[1]] = logic.cardList[changePositions[1]], logic.cardList[pos1+1]
		logic.cardList[pos1+2], logic.cardList[changePositions[2]] = logic.cardList[changePositions[2]], logic.cardList[pos1+2]
	} else {
		goto begin
	}
	return pos1
}

// 替换对子
func (logic *LittleGameLogic) changeDouble() int {
begin:
	changePositions := []int{}
	pos1 := rand.Int() % 15
	for i := pos1; i < 52; i++ {
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

		glog.Info("logic.cardList[i]", logic.cardList[changePositions[0]], "j=", logic.cardList[changePositions[1]], "k=", logic.cardList[changePositions[2]])

		logic.cardList[pos1], logic.cardList[changePositions[0]] = logic.cardList[changePositions[0]], logic.cardList[pos1]
		logic.cardList[pos1+1], logic.cardList[changePositions[1]] = logic.cardList[changePositions[1]], logic.cardList[pos1+1]
		logic.cardList[pos1+2], logic.cardList[changePositions[2]] = logic.cardList[changePositions[2]], logic.cardList[pos1+2]
	} else {
		goto begin
	}

	return pos1
}

// 替换单牌
func (logic *LittleGameLogic) changeSingle() int {
begin:
	changePositions := []int{}
	pos1 := rand.Int() % 15
	for i := pos1; i < 52; i++ {
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
		logic.cardList[pos1], logic.cardList[changePositions[0]] = logic.cardList[changePositions[0]], logic.cardList[pos1]
		logic.cardList[pos1+1], logic.cardList[changePositions[1]] = logic.cardList[changePositions[1]], logic.cardList[pos1+1]
		logic.cardList[pos1+2], logic.cardList[changePositions[2]] = logic.cardList[changePositions[2]], logic.cardList[pos1+2]
	} else {
		goto begin
	}

	return pos1
}

//测试概率
func (logic *LittleGameLogic) Test() bool {

	gailv_single := 0
	gailv_double := 0
	gailv_SHUN_ZI := 0
	gailv_JIN_HUA := 0
	gailv_SHUN_JIN := 0
	gailv_BAO_ZI := 0
	gailv_BAO_A := 0
	gailv_SEPCIAL := 0

	errCount_single := 0
	errCount_double := 0
	errCount_SHUN_ZI := 0
	errCount_JIN_HUA := 0
	errCount_SHUN_JIN := 0
	errCount_BAO_ZI := 0
	errCount_BAO_A := 0
	errCount_SEPCIAL := 0

	totalmoney := 0
	chip := 0
	chip0count := 0
	chip1count := 0
	chip2count := 0
	gameType := 4
	userid := "111"

	for i := 0; i < 100000; i++ {
		// 按概率抽牌
		t := GetCardConfigManager().GetRandCardType(gameType)
		glog.Info("littlegame_ShuffleCards test 牌型=", t)
		PrintCardType(t)

		sitepos := 0
		chip = rand.Int() % 3

		if t == CARD_TYPE_SINGLE {
			gailv_single++
			sitepos = logic.changeSingle()
			glog.Info("单牌 ", getCardColorString(logic.GetCardsInt32(sitepos)[0]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[0]), " ", getCardColorString(logic.GetCardsInt32(sitepos)[1]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[1]), " ", getCardColorString(logic.GetCardsInt32(sitepos)[2]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[2]))

			if getCardType([]byte{byte(logic.GetCardsInt32(sitepos)[0]), byte(logic.GetCardsInt32(sitepos)[1]), byte(logic.GetCardsInt32(sitepos)[2])}) != CARD_TYPE_SINGLE {
				errCount_single++
			}

			//算钱
			if chip == 0 {
				totalmoney += GetCardConfigManager().chips[gameType].Chip1
				chip0count++

				SaveLog(chip, gameType, GetCardConfigManager().chips[gameType].Chip1, t, userid)
			} else if chip == 1 {
				totalmoney += GetCardConfigManager().chips[gameType].Chip2
				chip1count++

				SaveLog(chip, gameType, GetCardConfigManager().chips[gameType].Chip2, t, userid)
			} else {
				totalmoney += GetCardConfigManager().chips[gameType].Chip3
				chip2count++

				SaveLog(chip, gameType, GetCardConfigManager().chips[gameType].Chip3, t, userid)
			}

		} else if t == CARD_TYPE_DOUBLE {
			gailv_double++
			sitepos = logic.changeDouble()
			glog.Info("对子 ", getCardColorString(logic.GetCardsInt32(sitepos)[0]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[0]), " ", getCardColorString(logic.GetCardsInt32(sitepos)[1]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[1]), "", getCardColorString(logic.GetCardsInt32(sitepos)[2]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[2]))

			if getCardType([]byte{byte(logic.GetCardsInt32(sitepos)[0]), byte(logic.GetCardsInt32(sitepos)[1]), byte(logic.GetCardsInt32(sitepos)[2])}) != CARD_TYPE_DOUBLE {
				errCount_double++
			}

			//算钱
			if chip == 0 {
				totalmoney -= GetCardConfigManager().chips[gameType].Chip1 * (GetCardConfigManager().multiple.Double - 1)
				chip0count++

				SaveLog(chip, gameType, 0-GetCardConfigManager().chips[gameType].Chip1*(GetCardConfigManager().multiple.Double-1), t, userid)
			} else if chip == 1 {
				totalmoney -= GetCardConfigManager().chips[gameType].Chip2 * (GetCardConfigManager().multiple.Double - 1)
				chip1count++

				SaveLog(chip, gameType, 0-GetCardConfigManager().chips[gameType].Chip2*(GetCardConfigManager().multiple.Double-1), t, userid)
			} else {
				totalmoney -= GetCardConfigManager().chips[gameType].Chip3 * (GetCardConfigManager().multiple.Double - 1)
				chip2count++

				SaveLog(chip, gameType, 0-GetCardConfigManager().chips[gameType].Chip3*(GetCardConfigManager().multiple.Double-1), t, userid)
			}

		} else if t == CARD_TYPE_SHUN_ZI {
			gailv_SHUN_ZI++
			sitepos = logic.changeShunZi()
			glog.Info("顺子 ", getCardColorString(logic.GetCardsInt32(sitepos)[0]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[0]), " ", getCardColorString(logic.GetCardsInt32(sitepos)[1]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[1]), " ", getCardColorString(logic.GetCardsInt32(sitepos)[2]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[2]))

			if getCardType([]byte{byte(logic.GetCardsInt32(sitepos)[0]), byte(logic.GetCardsInt32(sitepos)[1]), byte(logic.GetCardsInt32(sitepos)[2])}) != CARD_TYPE_SHUN_ZI {
				errCount_SHUN_ZI++
			}

			//算钱
			if chip == 0 {
				totalmoney -= GetCardConfigManager().chips[gameType].Chip1 * (GetCardConfigManager().multiple.ShunZi - 1)
				chip0count++

				SaveLog(chip, gameType, 0-GetCardConfigManager().chips[gameType].Chip1*(GetCardConfigManager().multiple.ShunZi-1), t, userid)
			} else if chip == 1 {
				totalmoney -= GetCardConfigManager().chips[gameType].Chip2 * (GetCardConfigManager().multiple.ShunZi - 1)
				chip1count++

				SaveLog(chip, gameType, 0-GetCardConfigManager().chips[gameType].Chip2*(GetCardConfigManager().multiple.ShunZi-1), t, userid)
			} else {
				totalmoney -= GetCardConfigManager().chips[gameType].Chip3 * (GetCardConfigManager().multiple.ShunZi - 1)
				chip2count++

				SaveLog(chip, gameType, 0-GetCardConfigManager().chips[gameType].Chip3*(GetCardConfigManager().multiple.ShunZi-1), t, userid)
			}
		} else if t == CARD_TYPE_JIN_HUA {
			gailv_JIN_HUA++
			sitepos = logic.changeJinHua()
			glog.Info("金花 ", getCardColorString(logic.GetCardsInt32(sitepos)[0]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[0]), " ", getCardColorString(logic.GetCardsInt32(sitepos)[1]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[1]), " ", getCardColorString(logic.GetCardsInt32(sitepos)[2]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[2]))

			if getCardType([]byte{byte(logic.GetCardsInt32(sitepos)[0]), byte(logic.GetCardsInt32(sitepos)[1]), byte(logic.GetCardsInt32(sitepos)[2])}) != CARD_TYPE_JIN_HUA {
				errCount_JIN_HUA++
			}

			//算钱
			if chip == 0 {
				totalmoney -= GetCardConfigManager().chips[gameType].Chip1 * (GetCardConfigManager().multiple.JinHua - 1)
				chip0count++

				SaveLog(chip, gameType, 0-GetCardConfigManager().chips[gameType].Chip1*(GetCardConfigManager().multiple.JinHua-1), t, userid)
			} else if chip == 1 {
				totalmoney -= GetCardConfigManager().chips[gameType].Chip2 * (GetCardConfigManager().multiple.JinHua - 1)
				chip1count++

				SaveLog(chip, gameType, 0-GetCardConfigManager().chips[gameType].Chip2*(GetCardConfigManager().multiple.JinHua-1), t, userid)
			} else {
				totalmoney -= GetCardConfigManager().chips[gameType].Chip3 * (GetCardConfigManager().multiple.JinHua - 1)
				chip2count++

				SaveLog(chip, gameType, 0-GetCardConfigManager().chips[gameType].Chip3*(GetCardConfigManager().multiple.JinHua-1), t, userid)
			}

		} else if t == CARD_TYPE_SHUN_JIN {
			gailv_SHUN_JIN++
			sitepos = logic.changeShunJin()
			glog.Info("顺金 ", getCardColorString(logic.GetCardsInt32(sitepos)[0]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[0]), " ", getCardColorString(logic.GetCardsInt32(sitepos)[1]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[1]), " ", getCardColorString(logic.GetCardsInt32(sitepos)[2]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[2]))

			if getCardType([]byte{byte(logic.GetCardsInt32(sitepos)[0]), byte(logic.GetCardsInt32(sitepos)[1]), byte(logic.GetCardsInt32(sitepos)[2])}) != CARD_TYPE_SHUN_JIN {
				errCount_SHUN_JIN++
			}

			//算钱
			if chip == 0 {
				totalmoney -= GetCardConfigManager().chips[gameType].Chip1 * (GetCardConfigManager().multiple.ShunJin - 1)
				chip0count++

				SaveLog(chip, gameType, 0-GetCardConfigManager().chips[gameType].Chip1*(GetCardConfigManager().multiple.ShunJin-1), t, userid)
			} else if chip == 1 {
				totalmoney -= GetCardConfigManager().chips[gameType].Chip2 * (GetCardConfigManager().multiple.ShunJin - 1)
				chip1count++

				SaveLog(chip, gameType, 0-GetCardConfigManager().chips[gameType].Chip2*(GetCardConfigManager().multiple.ShunJin-1), t, userid)
			} else {
				totalmoney -= GetCardConfigManager().chips[gameType].Chip3 * (GetCardConfigManager().multiple.ShunJin - 1)
				chip2count++

				SaveLog(chip, gameType, 0-GetCardConfigManager().chips[gameType].Chip3*(GetCardConfigManager().multiple.ShunJin-1), t, userid)
			}

		} else if t == CARD_TYPE_BAO_ZI {
			gailv_BAO_ZI++
			sitepos = logic.changeBaoZhi()
			glog.Info("豹子 ", getCardColorString(logic.GetCardsInt32(sitepos)[0]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[0]), " ", getCardColorString(logic.GetCardsInt32(sitepos)[1]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[1]), " ", getCardColorString(logic.GetCardsInt32(sitepos)[2]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[2]))

			if getCardType([]byte{byte(logic.GetCardsInt32(sitepos)[0]), byte(logic.GetCardsInt32(sitepos)[1]), byte(logic.GetCardsInt32(sitepos)[2])}) != CARD_TYPE_BAO_ZI {
				errCount_BAO_ZI++
			}

			//算钱
			if chip == 0 {
				totalmoney -= GetCardConfigManager().chips[gameType].Chip1 * (GetCardConfigManager().multiple.BaoZi - 1)
				chip0count++

				SaveLog(chip, gameType, 0-GetCardConfigManager().chips[gameType].Chip1*(GetCardConfigManager().multiple.BaoZi-1), t, userid)
			} else if chip == 1 {
				totalmoney -= GetCardConfigManager().chips[gameType].Chip2 * (GetCardConfigManager().multiple.BaoZi - 1)
				chip1count++

				SaveLog(chip, gameType, 0-GetCardConfigManager().chips[gameType].Chip2*(GetCardConfigManager().multiple.BaoZi-1), t, userid)
			} else {
				totalmoney -= GetCardConfigManager().chips[gameType].Chip3 * (GetCardConfigManager().multiple.BaoZi - 1)
				chip2count++

				SaveLog(chip, gameType, 0-GetCardConfigManager().chips[gameType].Chip3*(GetCardConfigManager().multiple.BaoZi-1), t, userid)
			}

		} else if t == CARD_TYPE_BAO_A {
			gailv_BAO_A++
			sitepos = logic.changeBaoZhiA()
			glog.Info("豹子A ", getCardColorString(logic.GetCardsInt32(sitepos)[0]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[0]), " ", getCardColorString(logic.GetCardsInt32(sitepos)[1]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[1]), " ", getCardColorString(logic.GetCardsInt32(sitepos)[2]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[2]))

			if getCardType([]byte{byte(logic.GetCardsInt32(sitepos)[0]), byte(logic.GetCardsInt32(sitepos)[1]), byte(logic.GetCardsInt32(sitepos)[2])}) != CARD_TYPE_BAO_A {
				errCount_BAO_A++
			}

			//算钱
			if chip == 0 {
				totalmoney -= GetCardConfigManager().chips[gameType].Chip1 * (GetCardConfigManager().multiple.BaoA - 1)
				chip0count++

				SaveLog(chip, gameType, 0-GetCardConfigManager().chips[gameType].Chip1*(GetCardConfigManager().multiple.BaoA-1), t, userid)
			} else if chip == 1 {
				totalmoney -= GetCardConfigManager().chips[gameType].Chip2 * (GetCardConfigManager().multiple.BaoA - 1)
				chip1count++

				SaveLog(chip, gameType, 0-GetCardConfigManager().chips[gameType].Chip2*(GetCardConfigManager().multiple.BaoA-1), t, userid)
			} else {
				totalmoney -= GetCardConfigManager().chips[gameType].Chip3 * (GetCardConfigManager().multiple.BaoA - 1)
				chip2count++

				SaveLog(chip, gameType, 0-GetCardConfigManager().chips[gameType].Chip3*(GetCardConfigManager().multiple.BaoA-1), t, userid)
			}
		} else if t == CARD_TYPE_SEPCIAL {
			gailv_SEPCIAL++
			sitepos = logic.changeSpecial()
			glog.Info("地龙 ", getCardColorString(logic.GetCardsInt32(sitepos)[0]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[0]), " ", getCardColorString(logic.GetCardsInt32(sitepos)[1]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[1]), " ", getCardColorString(logic.GetCardsInt32(sitepos)[2]), ":", getCardLogicValueTest(logic.GetCardsInt32(sitepos)[2]))

			if getCardType([]byte{byte(logic.GetCardsInt32(sitepos)[0]), byte(logic.GetCardsInt32(sitepos)[1]), byte(logic.GetCardsInt32(sitepos)[2])}) != CARD_TYPE_SEPCIAL {
				errCount_SEPCIAL++
			}

			//算钱
			if chip == 0 {
				totalmoney -= GetCardConfigManager().chips[gameType].Chip1 * (GetCardConfigManager().multiple.Special - 1)
				chip0count++

				SaveLog(chip, gameType, 0-GetCardConfigManager().chips[gameType].Chip1*(GetCardConfigManager().multiple.Special-1), t, userid)
			} else if chip == 1 {
				totalmoney -= GetCardConfigManager().chips[gameType].Chip2 * (GetCardConfigManager().multiple.Special - 1)
				chip1count++

				SaveLog(chip, gameType, 0-GetCardConfigManager().chips[gameType].Chip2*(GetCardConfigManager().multiple.Special-1), t, userid)
			} else {
				totalmoney -= GetCardConfigManager().chips[gameType].Chip3 * (GetCardConfigManager().multiple.Special - 1)
				chip2count++

				SaveLog(chip, gameType, 0-GetCardConfigManager().chips[gameType].Chip3*(GetCardConfigManager().multiple.Special-1), t, userid)
			}
		}

	}

	glog.Info("单牌数量=", gailv_single)
	glog.Info("单牌错误数量=", errCount_single)
	glog.Info("对子数量=", gailv_double)
	glog.Info("对子错误数量=", errCount_double)
	glog.Info("顺子数量=", gailv_SHUN_ZI)
	glog.Info("顺子错误数量=", errCount_SHUN_ZI)
	glog.Info("金花数量=", gailv_JIN_HUA)
	glog.Info("金花错误数量=", errCount_JIN_HUA)
	glog.Info("顺金数量=", gailv_SHUN_JIN)
	glog.Info("顺金错误数量=", errCount_SHUN_JIN)
	glog.Info("豹子数量=", gailv_BAO_ZI)
	glog.Info("豹子错误数量=", errCount_BAO_ZI)
	glog.Info("豹子A数量=", gailv_BAO_A)
	glog.Info("豹子A错误数量=", errCount_BAO_A)
	glog.Info("地龙数量=", gailv_SEPCIAL)
	glog.Info("地龙错误数量=", errCount_SEPCIAL)
	glog.Info("totalmoney=", totalmoney)
	glog.Info("chip0count=", chip0count)
	glog.Info("chip1count=", chip1count)
	glog.Info("chip2count=", chip2count)

	return true
}
