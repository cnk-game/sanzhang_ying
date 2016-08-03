package prize

import (
	mgo "gopkg.in/mgo.v2"
	"time"
	"util"
)

type ExchangeGoodsLog struct {
	UserId                 string    `bson:"userId"`
	UserName               string    `bson:"userName"`
	ItemId                 int       `bson:"itemId"`
	IsShipped              bool      `bson:"isShipped"`
	ShippingAddressName    string    `bson:"shippingAddressName"`
	ShippingAddressPhone   string    `bson:"shippingAddressPhone"`
	ShippingAddressAddress string    `bson:"shippingAddressAddress"`
	ShippingAddressZipCode string    `bson:"shippingAddressZipCode"`
	Time                   time.Time `bson:"time"`
	ShippingAddressQQ      string    `bson:"shippingAddressZipQQ"`
}

const (
	exchangeGoodsLogC = "exchange_goods_log"
)

func SaveExchangeGoodsLog(l *ExchangeGoodsLog) error {
	return util.WithLogCollection(exchangeGoodsLogC, func(c *mgo.Collection) error {
		return c.Insert(l)
	})
}
