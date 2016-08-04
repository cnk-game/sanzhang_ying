package main

import (
	"flag"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"pb"
	"sort"
	"syscall"
	"time"
	"util"
)

var url string
var origin string

func init() {
	flag.StringVar(&url, "url", "ws://192.168.1.133:8002/ws/", "url")
	flag.StringVar(&origin, "origin", "http://192.168.1.133/", "origin")
}

type RobotConfig struct {
	Username        bson.ObjectId `bson:"_id"`
	Nickname        string        `bson:"nickName"`
	Gender          int           `bson:"gender"`
	Sign            string        `bson:"sign"`
	Photo           string        `bson:"photo"`
	Gold            int           `bson:"gold"`
	Vip             int           `bson:"vip"`
	WinTimes        int           `bson:"winTimes"`
	LoseTimes       int           `bson:"loseTimes"`
	CurDayEarnGold  int           `bson:"curDayEarnGold"`
	CurWeekEarnGold int           `bson:"curWeekEarnGold"`
	MaxCards        []int         `bson:"maxCards"`
}

const (
	robotC = "robot"
)

func FindRobotConfigs() []*RobotConfig {
	configs := []*RobotConfig{}
	err := util.WithGameCollection(robotC, func(c *mgo.Collection) error {
		return c.Find(nil).All(&configs)
	})
	if err != nil {
		log.Fatal(err)
	}

	return configs
}

func getMatchType(i int) *pb.MatchType {
	switch i {
	case 0:
		return pb.MatchType_COMMON_LEVEL1.Enum()
	case 1:
		return pb.MatchType_COMMON_LEVEL2.Enum()
	case 2:
		return pb.MatchType_COMMON_LEVEL3.Enum()
	case 3:
		return pb.MatchType_COMMON_LEVEL1.Enum()
	//return pb.MatchType_MAGIC_ITEM_LEVEL1.Enum()
	case 4:
		return pb.MatchType_COMMON_LEVEL2.Enum()
	//return pb.MatchType_MAGIC_ITEM_LEVEL2.Enum()
	case 5:
		//return pb.MatchType_MAGIC_ITEM_LEVEL3.Enum()
		return pb.MatchType_COMMON_LEVEL3.Enum()
	case 6:
		return pb.MatchType_COMMON_LEVEL1.Enum()
	//		return pb.MatchType_SNG_LEVEL1.Enum()
	case 7:
		return pb.MatchType_COMMON_LEVEL2.Enum()
	//		return pb.MatchType_SNG_LEVEL2.Enum()
	case 8:
		return pb.MatchType_COMMON_LEVEL3.Enum()
	//		return pb.MatchType_SNG_LEVEL3.Enum()
	case 9:
		return pb.MatchType_WAN_REN_GAME.Enum()
	}
	return pb.MatchType_COMMON_LEVEL1.Enum()
}

type RobotRandSlice []*RobotConfig

func (p RobotRandSlice) Len() int { return len(p) }
func (p RobotRandSlice) Less(i, j int) bool { return rand.Float32() < 0.5 }
func (p RobotRandSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// Sort is a convenience method.
func (p RobotRandSlice) Sort() { sort.Sort(p) }

func main() {
	flag.Parse()

	robots := FindRobotConfigs()
	if len(robots) <= 0 {
		log.Fatal("获取机器人配置失败！")
	}

	sort.Sort(RobotRandSlice(robots))

	//		for i, robot := range robots {
	//			var r *Robot
	////			if i < 100 {
	////				r = NewRobot(robot, pb.MatchType_COMMON_LEVEL1.Enum())
	////			} else if i >= 100 && i < 200 {
	////				r = NewRobot(robot, pb.MatchType_COMMON_LEVEL2.Enum())
	////			} else if i >= 200 && i < 300 {
	////				r = NewRobot(robot, pb.MatchType_COMMON_LEVEL3.Enum())
	////			} else if i >= 300 && i < 380 {
	////				r = NewRobot(robot, pb.MatchType_WAN_REN_GAME.Enum())
	////			} else if i >= 380 && i < 490 {
	////				r = NewRobot(robot, pb.MatchType_SNG_LEVEL1.Enum())
	////			} else if i >= 490 && i < 500 {
	////				r = NewRobot(robot, pb.MatchType_SNG_LEVEL2.Enum())
	////			} else {
	////				break
	////			}
	//			r = NewRobot(robot, pb.MatchType_COMMON_LEVEL3.Enum())
	//			go r.Login(url, origin)
	//			if i >= 1 {
	//				break
	//			}
	//
	//			time.Sleep(20 * time.Millisecond)
	//		}

	count := len(robots) / 10
	//	count := 50
	print("robot count = ")
	print(count)
	for i, robot := range robots {
		r := NewRobot(robot, getMatchType(i/count))
		//r := NewRobot(robot, pb.MatchType_SNG_LEVEL1.Enum())
		print("start")
		print(url)
		go r.Login(url, origin)
		time.Sleep(20 * time.Millisecond)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	<-c
}
