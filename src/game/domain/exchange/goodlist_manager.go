package exchange

import (
	mgo "gopkg.in/mgo.v2"
	"sync"
	"github.com/golang/glog"
	"util"
)

type Good_Data struct {
	GoodsId int `bson:"goodsId"`	
	Name    string   `bson:"name"`
	Desc string	`bson:"desc"`
	IconRes	string	`bson:"icoRes"`
	NeedScore int `bson:"needScore"`
	TotalCount int `bson:"totalCount"`
	RemainderCount int `bson:"remainderCount"`
	
}

const (
	Good_DataC = "exchange_goods"
)



type GoodDataManager struct {
	sync.RWMutex
	items map[string]*Good_Data
}

var goodDataManager *GoodDataManager

func init() {
	goodDataManager = &GoodDataManager{}
	goodDataManager.items = make(map[string]*Good_Data)
	
}

func GetGoodManager() *GoodDataManager {
	return goodDataManager
}

func FindGoodDatas() ([]*Good_Data, error) {
	datas := []*Good_Data{}

	err := util.WithGameCollection(Good_DataC, func(c *mgo.Collection) error {
		return c.Find(nil).All(&datas)
	})
	return datas, err
}


func (m *GoodDataManager) Init() {
	m.Lock()
	defer m.Unlock()

	datas, err := FindGoodDatas()
	if err != nil {
		glog.Fatal(err)
	}

	for _, data := range datas {
		m.items[data.Name] = data		
	}
	
	glog.Info("exchange goods data request")
}


func (m *GoodDataManager) GetData() map[string]*Good_Data {
	m.Lock()
	defer m.Unlock()
	
	item := m.items
	return item
	

}

