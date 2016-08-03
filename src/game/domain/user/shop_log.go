package user

import (
	"code.google.com/p/goprotobuf/proto"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"pb"
	"time"
	"util"
)

type UserShopLog struct {
	UserId          string    `bson:"userId"`
	RechargeDiamond int       `bson:"rechargeDiamond"`
	ExchangeGold    int       `bson:"exchangeGold"`
	ExchangeHorn    int       `bson:"exchangeHorn"`
	BuyGoodsId      int       `bson:"buyGoodsId"`
	BuyGoodsCount   int       `bson:"buyGoodsCount"`
	Time            time.Time `bson:"time"`
}

const (
	userShopLogC = "user_shop_log"
)

func (log *UserShopLog) BuildMessage() *pb.MsgGetShopLogRes_ShopLogDef {
	msg := &pb.MsgGetShopLogRes_ShopLogDef{}
	msg.BuyDiamond = proto.Int(log.RechargeDiamond)
	msg.BuyGold = proto.Int(log.ExchangeGold)
	msg.BuyGoodsId = proto.Int(log.BuyGoodsId)
	msg.BuyGoodsCount = proto.Int(log.BuyGoodsCount)
	msg.RecordTime = proto.Int64(log.Time.Unix())
	msg.BuyHorn = proto.Int(log.ExchangeHorn)

	return msg
}

func FindShopLogs(userId string) ([]*UserShopLog, error) {
	logs := []*UserShopLog{}
	err := util.WithUserCollection(userShopLogC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"userId": userId}).All(&logs)
	})
	return logs, err
}

func SaveShopLog(log *UserShopLog) error {
	return util.WithUserCollection(userShopLogC, func(c *mgo.Collection) error {
		return c.Insert(log)
	})
}
