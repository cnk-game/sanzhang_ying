package util

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"fmt"
)

const (
    FINE = 1
    BAD  = 2
)

const (
	systemTipC = "system_tip"
    slotPoolLogC = "slot_pool_log"
    gameFeeLogC = "game_fee_log"
    gameGiftLogC = "game_gift_log"
    gameCharmPoolC = "game_charmpool_log"
    charmExchangeLogC = "charm_exchange_log"
    pCountGTypeLogC = "player_count_gtype_log"
)

func MongoLog_SetPlayerCountByGType(gameType int, pCount int) error {
    err := WithLogCollection(pCountGTypeLogC, func(c *mgo.Collection) error {
        _, err := c.Upsert(bson.M{"gameType": gameType}, bson.M{"$inc": bson.M{"playerCount": pCount}})
        return err
    })

    return err
}

func MongoLog_SystemTip(gold int) {
    tm := time.Now()
    date := fmt.Sprintf("%04d-%02d-%02d", tm.Year(), tm.Month(), tm.Day())
    go WithGameCollection(systemTipC, func(c *mgo.Collection) error {
        _, err := c.Upsert(bson.M{"date": date}, bson.M{"$inc": bson.M{"tipGold": gold}})
        return err
    })
}


func MongoLog_SlotPool(gold int, date string) {
    go WithGameCollection(slotPoolLogC, func(c *mgo.Collection) error {
        _, err := c.Upsert(bson.M{"date": date}, bson.M{"$inc": bson.M{"poolGold": gold}})
        return err
    })
}


func MongoLog_GameFee(gold int) {
    tm := time.Now()
    date := fmt.Sprintf("%04d-%02d-%02d", tm.Year(), tm.Month(), tm.Day())
    go WithGameCollection(gameFeeLogC, func(c *mgo.Collection) error {
        _, err := c.Upsert(bson.M{"date": date}, bson.M{"$inc": bson.M{"feeGold": gold}})
        return err
    })
}

func MongoLog_GameGiftFine(gold int) {
    tm := time.Now()
    date := fmt.Sprintf("%04d-%02d-%02d", tm.Year(), tm.Month(), tm.Day())
    go WithGameCollection(gameGiftLogC, func(c *mgo.Collection) error {
        _, err := c.Upsert(bson.M{"date": date, "type":"fine"}, bson.M{"$inc": bson.M{"giftGold": gold}})
        return err
    })
}

func MongoLog_GameGiftBad(gold int) {
    tm := time.Now()
    date := fmt.Sprintf("%04d-%02d-%02d", tm.Year(), tm.Month(), tm.Day())
    go WithGameCollection(gameGiftLogC, func(c *mgo.Collection) error {
        _, err := c.Upsert(bson.M{"date": date, "type":"bad"}, bson.M{"$inc": bson.M{"giftGold": gold}})
        return err
    })
}

func MongoLog_CharmPool(charm int) {
    go WithGameCollection(gameCharmPoolC, func(c *mgo.Collection) error {
        _, err := c.Upsert(bson.M{"name": "charmPool"}, bson.M{"$inc": bson.M{"charm": charm}})
        return err
    })
}

func MongoLog_CharmExchange(charm int) {
    tm := time.Now()
    date := fmt.Sprintf("%04d-%02d-%02d", tm.Year(), tm.Month(), tm.Day())
    go WithGameCollection(charmExchangeLogC, func(c *mgo.Collection) error {
        _, err := c.Upsert(bson.M{"date": date}, bson.M{"$inc": bson.M{"expendCharm": charm}})
        return err
    })

    go WithGameCollection(gameCharmPoolC, func(c *mgo.Collection) error {
        _, err := c.Upsert(bson.M{"name": "charmPool"}, bson.M{"$inc": bson.M{"expendCharm": charm}})
        return err
    })
}

//db.userInfo.find({age: {$gte: 23, $lte: 26}})
type SystemTip struct {
	Date      string   `bson:"date"`
	TipGold   int      `bson:"tipGold"`
}

func MongoLog_GetSystemTipLog(start, end string) []SystemTip {
    results := []SystemTip{}
    WithGameCollection(systemTipC, func(c *mgo.Collection) error{
        iter := c.Find(bson.M{"date": bson.M{"$gte": start, "$lte": end}}).Iter()
        res := SystemTip{}
        for iter.Next(&res) {
            results = append(results, res)
        }
        return nil
    })

    return results
}

type SlotPool struct {
	Date      string   `bson:"date"`
	PoolGold  int      `bson:"poolGold"`
}

func MongoLog_GetSlotPoolLog(start, end string) []SlotPool {
    results := []SlotPool{}
    WithGameCollection(slotPoolLogC, func(c *mgo.Collection) error{
        iter := c.Find(bson.M{"date": bson.M{"$gte": start, "$lte": end}}).Iter()
        res := SlotPool{}
        for iter.Next(&res) {
            results = append(results, res)
        }
        return nil
    })

    return results
}

type GameFee struct {
	Date      string   `bson:"date"`
	FeeGold   int      `bson:"feeGold"`
}

func MongoLog_GetGameFeeLog(start, end string) []GameFee {
    results := []GameFee{}
    WithGameCollection(gameFeeLogC, func(c *mgo.Collection) error{
        iter := c.Find(bson.M{"date": bson.M{"$gte": start, "$lte": end}}).Iter()
        res := GameFee{}
        for iter.Next(&res) {
            results = append(results, res)
        }
        return nil
    })

    return results
}

type GameGift struct {
	Date       string   `bson:"date"`
	GiftGold   int      `bson:"giftGold"`
}

func MongoLog_GetGameGiftFineLog(start, end string) []GameGift {
    results := []GameGift{}
    WithGameCollection(gameGiftLogC, func(c *mgo.Collection) error{
        iter := c.Find(bson.M{"type":"fine", "date": bson.M{"$gte": start, "$lte": end}}).Iter()
        res := GameGift{}
        for iter.Next(&res) {
            results = append(results, res)
        }
        return nil
    })

    return results
}

func MongoLog_GetGameGiftBadLog(start, end string) []GameGift {
    results := []GameGift{}
    WithGameCollection(gameGiftLogC, func(c *mgo.Collection) error{
        iter := c.Find(bson.M{"type":"bad", "date": bson.M{"$gte": start, "$lte": end}}).Iter()
        res := GameGift{}
        for iter.Next(&res) {
            results = append(results, res)
        }
        return nil
    })

    return results
}


type CharmPool struct {
	Name         string   `bson:"name"`
	Charm        int      `bson:"charm"`
	ExpendCharm  int      `bson:"expendCharm"`
}

func MongoLog_GetCharmPool() (*CharmPool, error) {
    ret := &CharmPool{}
    err := WithGameCollection(gameCharmPoolC, func(c *mgo.Collection) error{
        return c.Find(bson.M{"name":"charmPool"}).One(ret)
    })

    return ret, err
}

type CharmExchange struct {
	Date          string   `bson:"date"`
	ExpendCharm   int      `bson:"expendCharm"`
}

func MongoLog_GetCharmExchangeLog(start, end string) []CharmExchange {
    results := []CharmExchange{}
    WithGameCollection(charmExchangeLogC, func(c *mgo.Collection) error{
        iter := c.Find(bson.M{"date": bson.M{"$gte": start, "$lte": end}}).Iter()
        res := CharmExchange{}
        for iter.Next(&res) {
            results = append(results, res)
        }
        return nil
    })
    return results
}

type PlayerCountGTypeInfo struct {
	GameType      int      `bson:"gameType"`
	PlayerCount   int      `bson:"playerCount"`
}

func MongoLog_GetPlayerCountByGType() ([]*PlayerCountGTypeInfo, error) {
    ret := []*PlayerCountGTypeInfo{}
    err := WithLogCollection(pCountGTypeLogC, func(c *mgo.Collection) error{
        return c.Find(nil).All(&ret)
    })

    return ret, err
}

func MongoLog_ClearPlayerCountByGType() error {
	err := WithLogCollection(pCountGTypeLogC, func(c *mgo.Collection) error {
		_, e := c.RemoveAll(nil)
		return e
	})
	return err
}
