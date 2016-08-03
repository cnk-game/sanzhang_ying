package user

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"util"
)

const (
	DeviceUserC = "device_user_count"
)

type Imei_User_Count struct {
	Imei  string `bson:"imei"`
	Count int    `bson:"count"`
}

func GetDeviceUserCount(imei string) (int, error) {
	imeiUser := Imei_User_Count{"", 0}
	err := util.WithUserCollection(DeviceUserC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"imei": imei}).One(&imeiUser)
	})
	return imeiUser.Count, err
}

func SaveDeivceUserCount(imei string, count int) error {
	imeiUser := Imei_User_Count{imei, count}

	err := util.WithUserCollection(DeviceUserC, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"imei": imei}, &imeiUser)
		return err
	})

	return err
}
