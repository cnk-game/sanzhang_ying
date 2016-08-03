package prize

import (
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"sync"
	"util"
	"gopkg.in/mgo.v2/bson"
)

type ExchangeGoods struct {
	GoodsId   int    `bson:"goodsId"`
	Name      string `bson:"name"`
	Desc      string `bson:"desc"`
	IcoRes    string `bson:"icoRes"`
	NeedScore int    `bson:"needScore"`
	MaxCount  int    `bson:"maxCount"`
	TotalCount int	`bson:"totalCount"`
	RemainderCount int   `bson:"remainderCount"`
}

const (
	exchangeGoodsC = "exchange_goods"
)

func FindExchangeGoods() ([]*ExchangeGoods, error) {
	goods := []*ExchangeGoods{}
	err := util.WithGameCollection(exchangeGoodsC, func(c *mgo.Collection) error {
		return c.Find(nil).All(&goods)
	})
	return goods, err
}

type ExchangeGoodsManager struct {
	goods map[int]ExchangeGoods
	sync.RWMutex
}

var exchangeGoodsManager *ExchangeGoodsManager

func init() {
	exchangeGoodsManager = &ExchangeGoodsManager{}
	exchangeGoodsManager.goods = make(map[int]ExchangeGoods)
}

func GetExchangeGoodsManager() *ExchangeGoodsManager {
	return exchangeGoodsManager
}

func (m *ExchangeGoodsManager) Init() {
	m.Lock()
	defer m.Unlock()

	goods, err := FindExchangeGoods()
	if err != nil {
		glog.Error("加载兑换物品配置失败err:", err)
	}

	for _, item := range goods {
		m.goods[item.GoodsId] = *item
	}
}

func (m *ExchangeGoodsManager) GetExchangeGoods(itemId int) (ExchangeGoods, bool) {
	m.Lock()
	defer m.Unlock()
	
	item, ok := m.goods[itemId]
	return item, ok
}

//wjs 更新物品剩余数量到数据库
func (m *ExchangeGoodsManager) SetRemainderToDB(goodsId int, remainderCount int) error {

		//glog.Infof("goods.goodsid=",goodsId);
		//glog.Infof("goods.RemainderCount=",remainderCount);
	return util.WithGameCollection(exchangeGoodsC, func(c *mgo.Collection) error {
		return c.Update(bson.M{"goodsId": goodsId}, bson.M{"$set": bson.M{"remainderCount": remainderCount}})
	})
}

	//wjs 更新物品剩余数量到内存
func (m *ExchangeGoodsManager) SetRemainderToMem(goodsId int, good ExchangeGoods) {

	m.Lock()
	defer m.Unlock()
		if(goodsId!=0){
			delete(m.goods,goodsId);
			m.goods[goodsId] = good;			
		}				
	
}
