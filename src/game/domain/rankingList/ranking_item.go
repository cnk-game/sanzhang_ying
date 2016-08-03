package rankingList

import (
	"code.google.com/p/goprotobuf/proto"
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"pb"
	"sort"
	"sync"
	"time"
	"util"
)

type rankingItem struct {
	UserId   string `bson:"userId"`
	Value    int64  `bson:"value"`
	UserInfo []byte `bson:"userInfo"`
	userInfo *pb.UserDef
}

func (item *rankingItem) buildMessage(order int) *pb.RankingItem {
	msg := &pb.RankingItem{}
	msg.User = item.userInfo
	msg.Order = proto.Int(order)
	msg.RankingValue = proto.Int64(item.Value)

	return msg
}

type RankingItem struct {
	sync.RWMutex `bson:",omitempty"`
	RankingType  int            `bson:"rankingType"`
	MaxCount     int            `bson:"maxCount"`
	Items        []*rankingItem `bson:"items"`
	ResetTime    time.Time      `bson:"resetTime"`
}

func (ranking *RankingItem) BuildMessage(rankingType pb.RankingType) *pb.RankingList {
	msg := &pb.RankingList{}
	msg.Type = rankingType.Enum()

	switch rankingType {
	case pb.RankingType_COMPETITION:
		msg.ItemType = proto.Int32(1)
	}

	for i, item := range ranking.Items {
		msg.Items = append(msg.Items, item.buildMessage(i+1))
	}

	return msg
}

func (ranking *RankingItem) GetOrder1UserNickname() string {
	if len(ranking.Items) <= 0 {
		return ""
	}
	if ranking.Items[0].userInfo == nil {
		return ""
	}
	return ranking.Items[0].userInfo.GetNickName()
}

func NewRankingItem(rankingType, maxCount int) *RankingItem {
	item, err := findRankingItem(rankingType)
	if err != nil {
		item = &RankingItem{}
		item.RankingType = rankingType
		item.MaxCount = maxCount
	}

	return item
}

const (
	rankingListC = "ranking_list"
)

func findRankingItem(rankingType int) (*RankingItem, error) {
	item := &RankingItem{}
	err := util.WithUserCollection(rankingListC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"rankingType": rankingType}).One(item)
	})

	glog.Info("====>加载排行榜:", item, " err:", err)

	for _, item := range item.Items {
		item.userInfo = &pb.UserDef{}
		err = proto.Unmarshal(item.UserInfo, item.userInfo)
		if err != nil {
			glog.Error(err)
		}
	}

	return item, err
}

func saveRankingItem(item *RankingItem) error {
	return util.WithSafeUserCollection(rankingListC, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"rankingType": item.RankingType}, item)
		return err
	})
}

func SaveAsRankingItem(item *RankingItem, rankingListSaveC string) error {
	return util.WithLogCollection(rankingListSaveC, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"rankingType": item.RankingType}, item)
		return err
	})
}

func (ranking *RankingItem) SaveRankingItem() {
	for _, item := range ranking.Items {
		if item.userInfo != nil {
			item.UserInfo, _ = proto.Marshal(item.userInfo)
		}
	}
	err := saveRankingItem(ranking)
	if err != nil {
		glog.Error(err)
	}
}

func (ranking *RankingItem) SaveAsRankingItem(rankingListSaveC string) {
	for _, item := range ranking.Items {
		if item.userInfo != nil {
			item.UserInfo, _ = proto.Marshal(item.userInfo)
		}
	}
	err := SaveAsRankingItem(ranking, rankingListSaveC)
	if err != nil {
		glog.Error(err)
	}
}

type RankingItemSlice []*rankingItem

func (p RankingItemSlice) Len() int           { return len(p) }
func (p RankingItemSlice) Less(i, j int) bool { return p[i].Value > p[j].Value }
func (p RankingItemSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (ranking *RankingItem) UpdateRankingItem(user *pb.UserDef, value int64) {
	glog.V(2).Info("====>UpdateRankingItem rankingType:", ranking.RankingType, " userId:", user.GetUserId(), " value:", value)
	ranking.Lock()
	defer ranking.Unlock()

	if ranking.RankingType == int(RankingType_Charm) {
		if value <= 0 {
			index := -1
			for i, item := range ranking.Items {
				if item.UserId == user.GetUserId() {
					index = i
					break
				}
			}
			if index != -1 {
				ranking.Items = append(ranking.Items[:index], ranking.Items[index+1:]...)
				return
			}
		}
	}

	exist := false
	for _, item := range ranking.Items {
		if item.UserId == user.GetUserId() {
			item.Value = value
			item.userInfo = user
			exist = true
			break
		}
	}
	if !exist {
		item := &rankingItem{}
		item.UserId = user.GetUserId()
		item.Value = value
		item.userInfo = user
		ranking.Items = append(ranking.Items, item)
	}
	sort.Sort(RankingItemSlice(ranking.Items))

	if len(ranking.Items) > ranking.MaxCount {
		ranking.Items = ranking.Items[:ranking.MaxCount]
	}
}

func (ranking *RankingItem) ResetRankingItem() {
	ranking.Lock()
	defer ranking.Unlock()

	ranking.Items = []*rankingItem{}
	ranking.ResetTime = time.Now()
}
