package activedata

import (
	"fmt"
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"sync"
	"time"
	"util"
)

type Active_Data struct {
	ActivityId           string `bson:"activityId"`
	Name                 string `bson:"name"`
	IconRes              string `bson:"iconRes"`
	Desc                 string `bson:"desc"`
	ButtonTitle          string `bson:"buttonTitle"`
	ShowViewId           int    `bson:"showViewId"`
	DailyBeginShowSecond int    `bson:"dailyBeginShowSecond"`
	DailyEndShowSecond   int    `bson:"dailyEndShowSecond"`
	IsCompleteClose      int    `bson:"isCompleteClose"`
	StartDate            string `bson:"startDate"`
	ShowLeftButton       bool   `bson:"showLeftButton"`
	LeftButtonTitle      string `bson:"leftButtonTitle"`
	OpenUrl              string `bson:"openUrl"`
	OnlineTime           string `bson:"onlineTime"`
	OfflineTime          string `bson:"offlineTime"`
}

const (
	activity_dataC     = "activity_data"
	activity_dataxy    = "activity_data_xy"
	activity_datajinli = "activity_data_jinli"
	activity_datakupai = "activity_data_kupai"
	chnnel_jinli       = "212"
	chnnel_kupai       = "215"
	chnnel_xy          = "178"
)

type ActiveDataManager struct {
	sync.RWMutex
	items       map[string]*Active_Data
	items_xy    map[string]*Active_Data
	items_jinli map[string]*Active_Data
	items_kupai map[string]*Active_Data
}

var activeDataManager *ActiveDataManager

func init() {
	activeDataManager = &ActiveDataManager{}
	activeDataManager.items = make(map[string]*Active_Data)
	activeDataManager.items_xy = make(map[string]*Active_Data)
	activeDataManager.items_jinli = make(map[string]*Active_Data)
	activeDataManager.items_kupai = make(map[string]*Active_Data)

}

func GetActiveManager() *ActiveDataManager {
	return activeDataManager
}

func FindActiveDatas() ([]*Active_Data, error) {
	datas := []*Active_Data{}

	err := util.WithGameCollection(activity_dataC, func(c *mgo.Collection) error {
		return c.Find(nil).All(&datas)
	})
	return datas, err
}

func FindActiveDatas_xy() ([]*Active_Data, error) {
	datas := []*Active_Data{}

	err := util.WithGameCollection(activity_dataxy, func(c *mgo.Collection) error {
		return c.Find(nil).All(&datas)
	})
	return datas, err
}

func FindActiveDatas_jinli() ([]*Active_Data, error) {
	datas := []*Active_Data{}

	err := util.WithGameCollection(activity_datajinli, func(c *mgo.Collection) error {
		return c.Find(nil).All(&datas)
	})
	return datas, err
}

func FindActiveDatas_kupai() ([]*Active_Data, error) {
	datas := []*Active_Data{}

	err := util.WithGameCollection(activity_datakupai, func(c *mgo.Collection) error {
		return c.Find(nil).All(&datas)
	})
	return datas, err
}

func (m *ActiveDataManager) Init() {
	glog.Info("activ data request")
	m.Lock()
	defer m.Unlock()

	datas, err := FindActiveDatas()
	if err != nil {
		glog.Fatal(err)
	}

	for _, data := range datas {
		m.items[data.ActivityId] = data
	}

	datas_xy, err_xy := FindActiveDatas_xy()
	if err_xy != nil {
		glog.Fatal(err_xy)
	}

	for _, data_xy := range datas_xy {
		m.items_xy[data_xy.ActivityId] = data_xy
	}

	datas_jinli, err_jinli := FindActiveDatas_jinli()
	if err_jinli != nil {
		glog.Fatal(err_jinli)
	}

	for _, data_jinli := range datas_jinli {
		m.items_jinli[data_jinli.ActivityId] = data_jinli
	}

	datas_kupai, err_kupai := FindActiveDatas_kupai()
	if err_kupai != nil {
		glog.Fatal(err_kupai)
	}

	for _, data_kupai := range datas_kupai {
		m.items_kupai[data_kupai.ActivityId] = data_kupai
	}

}

func getCurrentTime() string {
	now := time.Now()

	year, mon, day := now.Date()
	hour, min, sec := now.Clock()
	//zone, _ := now.Zone()
	return fmt.Sprintf("%d%02d%02d%02d%02d%02d", year, mon, day, hour, min, sec)

}

func (m *ActiveDataManager) GetData(chnnelid string) map[string]*Active_Data {
	m.Lock()
	defer m.Unlock()
	curtime := getCurrentTime()
	glog.Info("curtime=", curtime)

	if chnnelid == chnnel_xy {

		items_xy := map[string]*Active_Data{}
		for _, data_xy := range m.items_xy {

			if curtime < data_xy.OnlineTime || curtime > data_xy.OfflineTime {

			} else {
				items_xy[data_xy.ActivityId] = data_xy

			}
		}

		return items_xy
	} else if chnnelid == chnnel_jinli {

		items_jinli := map[string]*Active_Data{}
		for _, data_jinli := range m.items_jinli {
			if curtime < data_jinli.OnlineTime || curtime > data_jinli.OfflineTime {

			} else {
				items_jinli[data_jinli.ActivityId] = data_jinli
			}
		}

		return items_jinli
	} else if chnnelid == chnnel_kupai {

		items_kupai := map[string]*Active_Data{}
		for _, data_kupai := range m.items_kupai {
			if curtime < data_kupai.OnlineTime || curtime > data_kupai.OfflineTime {

			} else {
				items_kupai[data_kupai.ActivityId] = data_kupai
			}
		}

		return items_kupai
	} else {

		glog.Info("m.items=", m.items)
		items := map[string]*Active_Data{}
		for _, data := range m.items {
			glog.Info("data=", data)
			if curtime < data.OnlineTime || curtime > data.OfflineTime {

			} else {
				items[data.ActivityId] = data
			}
		}
		return items
	}

}
