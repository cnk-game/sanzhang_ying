package slots

import (
    "math/rand"
    "time"
    "sort"
	"github.com/golang/glog"
	"config"
)

type SlotMachine struct {
	Coin             int
	cards            []byte
	ReplaceCardTimes int
}

const (
	MASK_COLOR         = 0xF0
	MASK_VALUE         = 0x0F
)

const (
    TYPE_SINGLE   = 1
    TYPE_PAIR     = 2
    TYPE_STR      = 3
    TYPE_FLUSH    = 4
    TYPE_SPECIAL  = 7
    TYPE_FLUSHSTR = 5
    TYPE_SET      = 6
)


type Colors []int
func (arr Colors) Len() int { return len(arr) }
func (arr Colors) Less(i, j int) bool { return rand.Float32() < 0.5 }
func (arr Colors) Swap(i, j int) { arr[i], arr[j] = arr[j], arr[i] }


type Cards []int
func (arr Cards) Len() int { return len(arr) }
func (arr Cards) Less(i, j int) bool { return arr[i] < arr[j] }
func (arr Cards) Swap(i, j int) { arr[i], arr[j] = arr[j], arr[i] }


type RandCards []int
func (arr RandCards) Len() int { return len(arr) }
func (arr RandCards) Less(i, j int) bool { return rand.Float32() < 0.5 }
func (arr RandCards) Swap(i, j int) { arr[i], arr[j] = arr[j], arr[i] }


func (m *SlotMachine) GetCardsLen() int {
	return len(m.cards)
}

func (m *SlotMachine) GetCards() []byte {
	b := make([]byte, 3)
	copy(b, m.cards)

	return b
}

func (m *SlotMachine) GetCard(pos int) int {
    rand.Seed(time.Now().UnixNano())
	return int(m.cards[pos])
}

func (m *SlotMachine) ReplaceCard(pos int, card byte) {
	m.cards[pos] = card
}

func (m *SlotMachine) SetCards(cards []byte) {
	m.cards = cards
}

func (m *SlotMachine) Reset() {
	m.Coin = 0
	m.cards = []byte{}
	m.ReplaceCardTimes = 0
}

func (m *SlotMachine) RandomCardType(change_config map[int]int) int {
    if change_config == nil || len(change_config) == 0 {
        change_config = m.GetSlotConfig()
    }
	ranInt := rand.Int() % 1000
	offset := 0

	for i := TYPE_SINGLE; i <= TYPE_SPECIAL; i++ {
	    offset += int(change_config[i])
	    if ranInt <= offset {
	        glog.Info("RandomCardType-------->",i)
	        glog.Info("RandomCardType,",ranInt)
	        return i
	    }
	}

	return 0
}

func (m *SlotMachine) GetSlotCardList() []byte {
    Slot_CardList := []byte{
        0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D,
        0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D,
        0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D,
        0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3A, 0x3B, 0x3C, 0x3D,
    }

    return Slot_CardList
}

func (m *SlotMachine) GetSlotConfig() map[int]int {
    Slot_Config := make(map[int] int)
    Slot_Config[TYPE_SINGLE] = 585
    Slot_Config[TYPE_PAIR] = 300
    Slot_Config[TYPE_FLUSH] = 50
    Slot_Config[TYPE_STR] = 30
    Slot_Config[TYPE_SPECIAL] = 20
    Slot_Config[TYPE_FLUSHSTR] = 10
    Slot_Config[TYPE_SET] = 5

    return Slot_Config
}

func (m *SlotMachine) GetSlotCards(card_type int) []byte {
    switch card_type {
    case TYPE_SINGLE:
        return m.GetSingle()
    case TYPE_PAIR:
        return m.GetPair()
    case TYPE_FLUSH:
        return m.GetFlush()
    case TYPE_STR:
        return m.GetStr()
    case TYPE_SPECIAL:
        return m.GetSpecial()
    case TYPE_FLUSHSTR:
        return m.GetFlushStr()
    case TYPE_SET:
        return m.GetSet()
    }

    return nil
}

func (m *SlotMachine) GetSet() []byte {
    num := rand.Int() % 13
    colors := Colors{0, 1, 2, 3}
    sort.Sort(colors)
    Slot_CardList := m.GetSlotCardList()

    return []byte{Slot_CardList[num+colors[0]*13], Slot_CardList[num+colors[1]*13], Slot_CardList[num+colors[2]*13]}
}

func (m *SlotMachine) GetFlushStr() []byte {
    start := rand.Int() % 12
    colors := Colors{0, 1, 2, 3}
    sort.Sort(colors)
    color := colors[0]
    Slot_CardList := m.GetSlotCardList()

    ret := []byte{}
    if start == 11 {
        ret = []byte{Slot_CardList[start+color*13], Slot_CardList[start+1+color*13], Slot_CardList[color*13]}
    } else {
        ret = []byte{Slot_CardList[start+color*13], Slot_CardList[start+1+color*13], Slot_CardList[start+2+color*13]}
    }

    return ret
}

func (m *SlotMachine) GetSpecial() []byte {
    card1, card2, card3 := 2, 3, 5
    colors := Colors{0, 1, 2, 3}
    sort.Sort(colors)
    Slot_CardList := m.GetSlotCardList()

    return []byte{Slot_CardList[card1+colors[0]*13 - 1], Slot_CardList[card2+colors[1]*13 - 1], Slot_CardList[card3+colors[2]*13] - 1}
}

func (m *SlotMachine) GetStr() []byte {
    start := rand.Int() % 12
    color1, color2, color3 := 0, 0, 0

    for {
        if !(color1 == color2 && color1 == color3 && color2 == color3) {
            break
        }
        color1 = rand.Int() % 4
        color2 = rand.Int() % 4
        color3 = rand.Int() % 4
    }

    Slot_CardList := m.GetSlotCardList()

    ret := []byte{}
    if start == 11 {
        ret = []byte{Slot_CardList[start+color1*13], Slot_CardList[start+1+color2*13], Slot_CardList[color3*13]}
    } else {
        ret = []byte{Slot_CardList[start+color1*13], Slot_CardList[start+1+color2*13], Slot_CardList[start+2+color3*13]}
    }

    return ret
}

func (m *SlotMachine) GetFlush() []byte {
    cards := Cards{0, 0, 0}
    colors := Colors{0, 1, 2, 3}
    sort.Sort(colors)
    color := colors[0]
    Slot_CardList := m.GetSlotCardList()

    for {
        if !(cards[0] == cards[1] || cards[0] == cards[2] || cards[1] == cards[2]) && !(cards[0]+1 == cards[1] && cards[1]+1 == cards[2]) && !(cards[0] == 0 && cards[1] == 11 && cards[2] == 12) {
            break
        }
        cards[0] = rand.Int() % 13
        cards[1] = rand.Int() % 13
        cards[2] = rand.Int() % 13

        sort.Sort(cards)
    }

    ret := []byte{Slot_CardList[cards[0]+color*13], Slot_CardList[cards[1]+color*13], Slot_CardList[cards[2]+color*13]}

    return ret
}

func (m *SlotMachine) GetSingle() []byte {
    cards := Cards{0, 0, 0}
    ret := []byte{}
    Slot_CardList := m.GetSlotCardList()
    for {
        cards_ok := false

        cards[0] = rand.Int() % 13
        cards[1] = rand.Int() % 13
        cards[2] = rand.Int() % 13
        sort.Sort(cards)
        if !(cards[0] == cards[1] || cards[0] == cards[2] || cards[1] == cards[2]) && !(cards[0]+1 == cards[1] && cards[1]+1 == cards[2]) && !(cards[0] == 0 && cards[1] == 11 && cards[2] == 12) && !(cards[0] == 2-1 || cards[1] == 3-1 || cards[2] == 5-1) {
            cards_ok = true
        }

        colors_ok := false

        color1 := rand.Int() % 4
        color2 := rand.Int() % 4
        color3 := rand.Int() % 4

         if !(color1 == color2 && color1 == color3 && color2 == color3) {
            colors_ok = true
         }

         if cards_ok && colors_ok {
            ret = append(ret, Slot_CardList[cards[0]+color1*13])
            ret = append(ret, Slot_CardList[cards[1]+color2*13])
            ret = append(ret, Slot_CardList[cards[2]+color3*13])
            break
         }
    }

    return ret
}

func (m *SlotMachine) GetPair() []byte {
    num := rand.Int() % 13
    single := 0
    for {
        single = rand.Int() % 13
        if single != num { break }
    }
    colors := Colors{0, 1, 2, 3}
    sort.Sort(colors)
    Slot_CardList := m.GetSlotCardList()

    ret := []byte{Slot_CardList[num+colors[0]*13], Slot_CardList[num+colors[1]*13]}
    sort.Sort(colors)
    ret = append(ret, Slot_CardList[single+colors[0]*13])

    return ret
}

func getCardValue(card byte) int {
	return int(card & MASK_VALUE) - 1
}

func getCardColor(card byte) int {
	return int(card & MASK_COLOR) / 16
}

func (m *SlotMachine)UpdatePerByPoolValue(poolValue, betValue int) map[int]int {
    configs := m.GetSlotConfig()
    if config.Set_times * betValue > poolValue {
        configs[TYPE_SINGLE] += configs[TYPE_SET]
        configs[TYPE_SET] = 0
    }

    if config.FlushStr_times * betValue > poolValue {
        configs[TYPE_SINGLE] += configs[TYPE_FLUSHSTR]
        configs[TYPE_FLUSHSTR] = 0
    }

    if config.Special_times * betValue > poolValue {
        configs[TYPE_SINGLE] += configs[TYPE_SPECIAL]
        configs[TYPE_SPECIAL] = 0
    }

    if config.Str_times * betValue > poolValue {
        configs[TYPE_SINGLE] += configs[TYPE_STR]
        configs[TYPE_STR] = 0
    }

    if config.Flush_times * betValue > poolValue {
        configs[TYPE_SINGLE] += configs[TYPE_FLUSH]
        configs[TYPE_FLUSH] = 0
    }

    if config.Pair_times * betValue > poolValue {
        configs[TYPE_SINGLE] += configs[TYPE_PAIR]
        configs[TYPE_PAIR] = 0
    }

    return configs
}

func (m *SlotMachine)MergeCardPer(conf_pool, conf_card map[int]int) map[int]int {
    glog.Info("MergeCardPer in. >>>>>>")
    glog.Info("conf_pool=", conf_pool)
    glog.Info("conf_card=", conf_card)
    for k, v := range conf_pool {
        if v == 0 && conf_card[k] != 0 {
            conf_card[TYPE_SINGLE] += conf_card[k]
            conf_card[k] = 0
        }
    }
    glog.Info("ret=", conf_card)
    glog.Info("MergeCardPer out. <<<<<<")
    return conf_card
}

func (m *SlotMachine)CheckUpdate(change_pos int) (map[int]int, bool){
    others := []byte{}
    for i := 0; i < len(m.cards); i ++ {
        if i != change_pos {
            others = append(others, m.cards[i])
        }
    }

    configs := m.GetSlotConfig()
    if len(others) != 2 {
        glog.Info("CheckUpdate error, len(others) != 2.", len(others))
        return configs, false
    }

    if m.CheckUpdateSingle(others) {
        if !m.CheckUpdateFlush(others) {
            configs[TYPE_SINGLE] += configs[TYPE_FLUSH]
            configs[TYPE_FLUSH] = 0
        }

        if !m.CheckUpdateStr(others) {
            configs[TYPE_SINGLE] += configs[TYPE_STR]
            configs[TYPE_STR] = 0
        }

        if !m.CheckUpdateSpecial(others) {
            configs[TYPE_SINGLE] += configs[TYPE_SPECIAL]
            configs[TYPE_SPECIAL] = 0
        }

        if !m.CheckUpdateStr(others) || !m.CheckUpdateFlush(others) {
            configs[TYPE_SINGLE] += configs[TYPE_FLUSHSTR]
            configs[TYPE_FLUSHSTR] = 0
        }
        configs[TYPE_SINGLE] += configs[TYPE_SET]
        configs[TYPE_SET] = 0
    } else {
        configs[TYPE_PAIR] += configs[TYPE_SINGLE]
        configs[TYPE_SINGLE] = 0
        configs[TYPE_PAIR] += configs[TYPE_FLUSH]
        configs[TYPE_FLUSH] = 0
        configs[TYPE_PAIR] += configs[TYPE_STR]
        configs[TYPE_STR] = 0
        configs[TYPE_PAIR] += configs[TYPE_SPECIAL]
        configs[TYPE_SPECIAL] = 0
        configs[TYPE_PAIR] += configs[TYPE_FLUSHSTR]
        configs[TYPE_FLUSHSTR] = 0
    }

    glog.Info("CheckUpdate out, configs=", configs)
    return configs, true
}

func (m *SlotMachine) CheckUpdateSingle(others []byte) bool {
    card1 := getCardValue(others[0])
    card2 := getCardValue(others[1])

    if card1 == card2 {
        return false
    }

    return true
}

func (m *SlotMachine) CheckUpdateFlush(others []byte) bool {
    color1 := getCardColor(others[0])
    color2 := getCardColor(others[1])

    if color1 != color2 {
        return false
    }

    return true
}

func (m *SlotMachine) CheckUpdateStr(others []byte) bool {
    card1 := getCardValue(others[0])
    card2 := getCardValue(others[1])
    if card1 > card2 {
        card1, card2 = card2, card1
    }

    if card1+1 == card2 || card1+2 == card2 {
        return true
    }

    if (card1 == 0 && card2 == 11) || (card1 == 0 && card2 == 12) {
        return true
    }

    return false
}

func (m *SlotMachine) CheckUpdateSpecial(others []byte) bool {
    card1, color1 := getCardValue(others[0]), getCardColor(others[0])
    card2, color2 := getCardValue(others[1]), getCardColor(others[1])
    if card1 > card2 {
        card1, card2 = card2, card1
    }

    num_ok := false
    if card1 == 2-1 && card2 == 3-1 || card1 == 2-1 && card2 == 5-1 || card1 == 3-1 && card2 == 5-1 {
        num_ok = true
    }

    color_ok := false
    if color1 != color2 {
        color_ok = true
    }

    if num_ok && color_ok {
        return true
    }

    return false
}

func (m *SlotMachine) UpdateCardByType(card_type int, change_pos int) (byte, bool) {
    must := []byte{}
    for i := 0; i < len(m.cards); i ++ {
        if i != change_pos {
            must = append(must, m.cards[i])
        }
    }

    var card byte = 0
    switch card_type {
    case TYPE_SINGLE:
        card = m.Update2Single(must)
    case TYPE_PAIR:
        card = m.Update2Pair(must)
    case TYPE_FLUSH:
        card = m.Update2Flush(must)
    case TYPE_STR:
        card = m.Update2Str(must)
    case TYPE_SPECIAL:
        card = m.Update2Special(must)
    case TYPE_FLUSHSTR:
        card = m.Update2FlushStr(must)
    case TYPE_SET:
        card = m.Update2Set(must)
    }

    glog.Info("UpdateCardByType--------",card_type)
    glog.Info(must)
    glog.Info(card)
    glog.Info("UpdateCardByType--------")

    return card, true
}

func (m *SlotMachine) Update2Set(must []byte) byte {
    card1, color1 := getCardValue(must[0]), getCardColor(must[0])
    color2 := getCardColor(must[1])
    glog.Info("Update2FlushStr in---card1, color1, color2", card1, color1, color2)
    ret_card := card1
    ret_color := 0
    for {
        ret_color = rand.Int() % 4
        if ret_color != color1 && ret_color != color2 { break }
    }
    Slot_CardList := m.GetSlotCardList()


    return Slot_CardList[ret_card+ret_color*13]
}

func (m *SlotMachine) Update2FlushStr(must []byte) byte {
    card1, color1 := getCardValue(must[0]), getCardColor(must[0])
    card2 := getCardValue(must[1])
    glog.Info("Update2FlushStr in---card1, color1, card2", card1, color1, card2)
    ret_color := color1
    can_cards := RandCards{}

    if card1 > card2 {
        card1, card2 = card2, card1
    }

    if card1 != 0 {
        if card1+1 == card2 {
            can_cards = append(can_cards, card1-1)
            can_cards = append(can_cards, card2+1)
        } else {
            can_cards = append(can_cards, card1+1)
        }
    } else {
        if card1+1 == card2 {
            can_cards = append(can_cards, card2+1)
        } else if card1+2 == card2 {
            can_cards = append(can_cards, card1+1)
        } else if card2 == 11 {
            can_cards = append(can_cards, 12)
        } else {
            can_cards = append(can_cards, 11)
        }
    }
    sort.Sort(can_cards)
    Slot_CardList := m.GetSlotCardList()

    return Slot_CardList[can_cards[0]+ret_color*13]
}

func (m *SlotMachine) Update2Special(must []byte) byte {
    card1, color1 := getCardValue(must[0]), getCardColor(must[0])
    card2, color2 := getCardValue(must[1]), getCardColor(must[1])
    glog.Info("Update2Special in---card1, color1, card2, color2", card1, color1, card2, color2)
    ret_card, ret_color := 0, 0

    if card1 > card2 {
        card1, card2 = card2, card1
    }

    if card1 == 2-1 && card2 == 3-1 {
        ret_card = 5-1
    } else if card1 == 2-1 && card2 == 5-1 {
        ret_card = 3-1
    } else {
        ret_card = 2-1
    }

    for {
        ret_color = rand.Int() % 4
        if ret_color != color1 && ret_color != color2 { break }
    }
    Slot_CardList := m.GetSlotCardList()


    return Slot_CardList[ret_card+ret_color*13]
}

func (m *SlotMachine) Update2Str(must []byte) byte {
    card1, color1 := getCardValue(must[0]), getCardColor(must[0])
    card2, color2 := getCardValue(must[1]), getCardColor(must[1])
    glog.Info("Update2Str in---card1, color1, card2, color2", card1, color1, card2, color2)
    ret_color := 0
    can_cards := RandCards{}

    if card1 > card2 {
        card1, card2 = card2, card1
    }

    if card1 != 0 {
        if card1+1 == card2 {
            can_cards = append(can_cards, card1-1)
            can_cards = append(can_cards, card2+1)
        } else {
            can_cards = append(can_cards, card1+1)
        }
    } else {
        if card1+1 == card2 {
            can_cards = append(can_cards, card2+1)
        } else if card1+2 == card2 {
            can_cards = append(can_cards, card1+1)
        } else if card2 == 11 {
            can_cards = append(can_cards, 12)
        } else {
            can_cards = append(can_cards, 11)
        }
    }

    for {
        ret_color = rand.Int() % 4
        if ret_color != color1 || ret_color != color2 || color1 != color2 { break }
    }
    sort.Sort(can_cards)
    Slot_CardList := m.GetSlotCardList()

    return Slot_CardList[can_cards[0]+ret_color*13]
}

func (m *SlotMachine) Update2Flush(must []byte) byte {
    card1, color1 := getCardValue(must[0]), getCardColor(must[0])
    card2 := getCardValue(must[1])
    glog.Info("Update2Flush in---card1, color1, card2", card1, color1, card2)
    ret_color, ret_card := color1, 0

    for {
        ret_card = rand.Int() % 13
        if ret_card == card1 || ret_card == card2 {
            continue
        }
        s_cards := Cards{card1, card2}
        s_cards = append(s_cards, ret_card)
        sort.Sort(s_cards)

        if s_cards[0]+1 == s_cards[1] && s_cards[1]+1 == s_cards[2] {
            continue
        }

        if s_cards[0] == 0 && s_cards[1] == 11 && s_cards[2] == 12 {
            continue
        }
        break
    }
    Slot_CardList := m.GetSlotCardList()

    return Slot_CardList[ret_card+ret_color*13]
}

func (m *SlotMachine) Update2Pair(must []byte) byte {
    card1, color1 := getCardValue(must[0]), getCardColor(must[0])
    card2, color2 := getCardValue(must[1]), getCardColor(must[1])
    glog.Info("Update2Pair in---card1, color1, card2, color2", card1, color1, card2, color2)
    ret_card, ret_color := 0, 0

    if card1 == card2 {
        glog.Info("Update2Pair card1 == card2---->", card1, card2)
        for {
            ret_card = rand.Int() % 13
            if ret_card != card1 {
                break
            }
        }
        ret_color = rand.Int() % 4
    } else {
        glog.Info("Update2Pair card1 != card2---->", card1, card2)
        s_cards := RandCards{card1, card2}
        sort.Sort(s_cards)
        ret_card = s_cards[0]

        for {
            ret_color = rand.Int() % 4
            if ret_card == card1 {
                if ret_color != color1 { break }
            } else {
                if ret_color != color2 { break }
            }
        }
        glog.Info("Update2Pair ---->ret_card, ret_color", ret_card, ret_color)
    }
    Slot_CardList := m.GetSlotCardList()

    return Slot_CardList[ret_card+ret_color*13]
}

func (m *SlotMachine) Update2Single(must []byte) byte {
    card1, color1 := getCardValue(must[0]), getCardColor(must[0])
    card2, color2 := getCardValue(must[1]), getCardColor(must[1])
    glog.Info("Update2Single in---card1, color1, card2, color2", card1, color1, card2, color2)
    ret_card, ret_color := 0, 0

    for {
        ret_card = rand.Int() % 13
        if ret_card == card1 || ret_card == card2 {
            continue
        }
        s_cards := Cards{card1, card2}
        s_cards = append(s_cards, ret_card)
        sort.Sort(s_cards)

        if s_cards[0]+1 == s_cards[1] && s_cards[1]+1 == s_cards[2] {
            continue
        }

        if s_cards[0] == 0 && s_cards[1] == 11 && s_cards[2] == 12 {
            continue
        }

        if s_cards[0] == 2 && s_cards[1] == 3 && s_cards[2] == 5 {
            continue
        }
        break
    }

    for {
        ret_color = rand.Int() % 4
        if ret_color != color1 || ret_color != color2 || color1 != color2 { break }
    }

    Slot_CardList := m.GetSlotCardList()


    return Slot_CardList[ret_card+ret_color*13]
}
