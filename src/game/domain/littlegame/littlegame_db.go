package littlegame

import (
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"math/rand"
	"sync"
	"time"
	"util"
)

type LittleCardConfig struct {
	Single   int `bson:"single"`
	Double   int `bson:"double"`
	ShunZi   int `bson:"shuZi"`
	JinHua   int `bson:"jinHua"`
	ShunJin  int `bson:"shunJin"`
	BaoZi    int `bson:"baoZi"`
	BaoA     int `bson:"baoA"`
	Special  int `bson:"special"`
	Gametype int `bson:"gametype"`
	total    int
}

type LittleCardMultiple struct {
	Single  int `bson:"single"`
	Double  int `bson:"double"`
	ShunZi  int `bson:"shuZi"`
	JinHua  int `bson:"jinHua"`
	ShunJin int `bson:"shunJin"`
	BaoZi   int `bson:"baoZi"`
	BaoA    int `bson:"baoA"`
	Special int `bson:"special"`
}

type LittleChip struct {
	LevelId int `bson:"levelId"`
	Chip1   int `bson:"chip1"`
	Chip2   int `bson:"chip2"`
	Chip3   int `bson:"chip3"`
}

func (c *LittleCardMultiple) calcSort() bool {
	return (c.BaoA > c.BaoZi) && (c.BaoZi > c.ShunJin) && (c.ShunJin > c.JinHua) && (c.JinHua > c.ShunZi) && (c.ShunZi > c.Double) && (c.Double > c.Single)
}

func (c *LittleCardConfig) calcTotal() {
	c.total = c.Single + c.Double + c.ShunZi + c.JinHua + c.ShunJin + c.BaoZi + c.BaoA + c.Special
}

const (
	littleConfigC   = "littlegame_config"
	littleMultipleC = "littlegame_multiple"
	LittleChipC     = "ligame_chip_config"
)

const (
	CARD_TYPE_SINGLE   = 0 // 单牌类型
	CARD_TYPE_DOUBLE   = 7 // 对子类型
	CARD_TYPE_SHUN_ZI  = 6 // 顺子类型
	CARD_TYPE_JIN_HUA  = 5 // 金花类型
	CARD_TYPE_SHUN_JIN = 4 // 同花顺类型(顺金)
	CARD_TYPE_BAO_ZI   = 3 // 豹子类型
	CARD_TYPE_BAO_A    = 2 // 豹子A类型
	CARD_TYPE_SEPCIAL  = 1 //地龙类型
)

func FindLittleCardConfigs() ([]*LittleCardConfig, error) {
	configs := []*LittleCardConfig{}
	err := util.WithGameCollection(littleConfigC, func(c *mgo.Collection) error {
		return c.Find(nil).All(&configs)
	})

	return configs, err
}

func FindLittleCardMultiples() ([]*LittleCardMultiple, error) {
	multiples := []*LittleCardMultiple{}
	err := util.WithGameCollection(littleMultipleC, func(c *mgo.Collection) error {
		return c.Find(nil).All(&multiples)
	})

	return multiples, err
}

func FindLittleChips() ([]*LittleChip, error) {
	chips := []*LittleChip{}
	err := util.WithGameCollection(LittleChipC, func(c *mgo.Collection) error {
		return c.Find(nil).All(&chips)
	})

	return chips, err
}

type LittleCardConfigManager struct {
	sync.RWMutex
	configs         map[int]*LittleCardConfig
	defaultConfig   *LittleCardConfig
	multiple        *LittleCardMultiple
	defaultMultiple *LittleCardMultiple
	chips           map[int]*LittleChip
}

var cardConfigM *LittleCardConfigManager

func init() {
	rand.Seed(time.Now().UnixNano())
	cardConfigM = &LittleCardConfigManager{}
	cardConfigM.configs = make(map[int]*LittleCardConfig)
	cardConfigM.chips = make(map[int]*LittleChip)

	config := &LittleCardConfig{}
	config.Single = 57000
	config.Double = 37000
	config.ShunZi = 5000
	config.JinHua = 1000
	config.ShunJin = 80
	config.BaoZi = 5
	config.BaoA = 2
	config.Special = 0
	config.calcTotal()
	cardConfigM.defaultConfig = config

	cardConfigM.multiple = &LittleCardMultiple{}

	multiple := &LittleCardMultiple{}
	multiple.Single = 0
	multiple.Double = 2
	multiple.ShunZi = 3
	multiple.JinHua = 5
	multiple.ShunJin = 30
	multiple.BaoZi = 80
	multiple.BaoA = 800
	multiple.Special = 1000
	cardConfigM.defaultMultiple = multiple

	cardConfigM.chips = make(map[int]*LittleChip)
	littlechip1 := &LittleChip{}
	littlechip1.LevelId = 1
	littlechip1.Chip1 = 100
	littlechip1.Chip2 = 200
	littlechip1.Chip3 = 500
	cardConfigM.chips[littlechip1.LevelId] = littlechip1

	littlechip2 := &LittleChip{}
	littlechip2.LevelId = 2
	littlechip2.Chip1 = 1000
	littlechip2.Chip2 = 5000
	littlechip2.Chip3 = 10000
	cardConfigM.chips[littlechip2.LevelId] = littlechip2

	littlechip3 := &LittleChip{}
	littlechip3.LevelId = 3
	littlechip3.Chip1 = 10000
	littlechip3.Chip2 = 20000
	littlechip3.Chip3 = 50000
	cardConfigM.chips[littlechip3.LevelId] = littlechip3

	littlechip4 := &LittleChip{}
	littlechip4.LevelId = 4
	littlechip4.Chip1 = 500
	littlechip4.Chip2 = 2000
	littlechip4.Chip3 = 5000
	cardConfigM.chips[littlechip4.LevelId] = littlechip4

}

func GetCardConfigManager() *LittleCardConfigManager {
	return cardConfigM
}

func (m *LittleCardConfigManager) Init() bool {
	m.Lock()
	defer m.Unlock()

	configs, err := FindLittleCardConfigs()
	if err != nil && err != mgo.ErrNotFound {
		glog.Error(err)
		return false
	}

	for _, config := range configs {
		glog.Info("==>littlegamecard_config:", config)
		config.calcTotal()
		if config.total == 0 {
			glog.Error("牌型配置错误,概率总和为0 config:", config)
			continue
		}
		m.configs[config.Gametype] = config
	}

	multiples, err2 := FindLittleCardMultiples()
	if err2 != nil && err2 != mgo.ErrNotFound {
		glog.Error(err2)
		return false
	}

	for _, multiple := range multiples {
		glog.Info("==>littlegamecard_mulitple:", multiple)
		multiple.calcSort()
		if !multiple.calcSort() {
			glog.Error("倍数配置错误,不是呈递增关系 multiple:", multiple)
			continue
		}
		m.multiple = multiple
	}

	chips, err3 := FindLittleChips()
	if err3 != nil && err3 != mgo.ErrNotFound {
		glog.Error(err3)
		return false
	}

	for _, chip := range chips {
		glog.Info("==>littlegamecard_chip:", chip)

		//m.chips = append(m.chips, chip)
		m.chips[chip.LevelId] = chip
	}

	return true
}

func (m *LittleCardConfigManager) GetRandCardType(gameType int) int {
	m.RLock()
	defer m.RUnlock()

	config := m.configs[gameType]
	if config == nil {
		glog.Info(" 找不到牌型配置，使用默认配置")
		config = m.defaultConfig
	}

	offset := 0
	r := rand.Int() % config.total
	if r >= offset && r < offset+config.Single {
		return CARD_TYPE_SINGLE
	}

	offset += config.Single
	if r >= offset && r < offset+config.Double {
		return CARD_TYPE_DOUBLE
	}

	offset += config.Double
	if r >= offset && r < offset+config.ShunZi {
		return CARD_TYPE_SHUN_ZI
	}

	offset += config.ShunZi
	if r >= offset && r < offset+config.JinHua {
		return CARD_TYPE_JIN_HUA
	}

	offset += config.JinHua
	if r >= offset && r < offset+config.ShunJin {
		return CARD_TYPE_SHUN_JIN
	}

	offset += config.ShunJin
	if r >= offset && r < offset+config.BaoZi {
		return CARD_TYPE_BAO_ZI
	}

	offset += config.BaoZi
	if r >= offset && r < offset+config.BaoA {
		return CARD_TYPE_BAO_A
	}

	offset += config.Special
	if r >= offset && r < offset+config.Special {
		return CARD_TYPE_SEPCIAL
	}

	return CARD_TYPE_SINGLE
}

func (m *LittleCardConfigManager) GetChipData() map[int]*LittleChip {
	m.Lock()
	defer m.Unlock()

	chips := m.chips
	return chips

}

func (m *LittleCardConfigManager) GetMultipeData() *LittleCardMultiple {
	m.Lock()
	defer m.Unlock()

	multiple := m.multiple
	return multiple

}
