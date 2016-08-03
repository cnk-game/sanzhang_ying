package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	mgo "gopkg.in/mgo.v2"
	"log"
	"math/rand"
	"os"
	"time"
	"util"
)

const (
	CD_KEY   = "cd_key"
	cdKeyC   = "cd_key"
	maxValue = 99999
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randChar() int {
	return 'a' + rand.Int()%25 + 1
}

func genCDKey() string {
	v := util.NextSequenceValue(CD_KEY)
	if v == -1 {
		panic(errors.New("gen cdkey failed"))
	}
	if v > maxValue {
		panic(errors.New(fmt.Sprintf("gen cdkey failed v exceed maxValue v:", v)))
	}
	key := fmt.Sprintf("%05v", v)
	return fmt.Sprintf("%c%c%c%c%c%c%c%c", randChar(), key[0], key[1], randChar(), key[2], key[3], randChar(), key[4])
}

type CDKey struct {
	Type      int    `bson:"type"`
	Key       string `bson:"key"`
	Gold      int    `bson:"gold"`
	Diamond   int    `bson:"diamond"`
	Score     int    `bson:"score"`
	ItemType  int    `bson:"itemType"`
	ItemCount int    `bson:"itemCount"`
}

func SaveCDKey(cdkey *CDKey) error {
	return util.WithSafeUserCollection(cdKeyC, func(c *mgo.Collection) error {
		return c.Insert(cdkey)
	})
}

func DumpToFile(path string, b []byte) {
	log.Println("dump file:", path)
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f.Write(b)
	f.Sync()
}

func Marshal(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "")
}

func GenCdKeys(keyType, count, gold, diamond, itemType, itemCount int, filename string) {
	result := []string{}

	for i := 0; i < count; i++ {
		key := &CDKey{}
		key.Type = keyType
		key.Key = genCDKey()
		key.Gold = gold
		key.Diamond = diamond
		key.Score = 0
		key.ItemType = itemType
		key.ItemCount = itemCount

		err := SaveCDKey(key)
		if err == nil {
			result = append(result, key.Key)
		}
	}

	b, err := Marshal(result)
	if err != nil {
		fmt.Println("marshal failed err:", err)
		return
	}
	DumpToFile(filename, b)
}

func main() {
	flag.Parse()

	GenCdKeys(1, 2000, 2000, 0, 0, 0, "qipai_key_15_11_12_v1.json")
}
