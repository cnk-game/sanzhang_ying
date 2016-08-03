package randBugle

import (
	domainUser "game/domain/user"
	"github.com/golang/glog"
	"math/rand"
	"pb"
	"time"
    "io/ioutil"
    "strings"
    "fmt"
    "util"
)

const (
    format1 = "恭喜%v玩家在菜鸟场单场赢得%v金币"
    format2 = "恭喜%v玩家在高手场单场赢得%v金币"
    format3 = "恭喜%v玩家在精英场单场赢得%v金币"
    format4 = "恭喜%v在一夜暴富中奖！豹子！豹子！豹子！"
    format5 = "恭喜%v在一夜暴富中奖！顺金！顺金！顺金！"
    format6 = "恭喜%v获得了特权2权限！领取魅力值到商场兑换心仪的礼品！"
    format7 = "恭喜%v获得了特权3权限！领取魅力值到商场兑换心仪的礼品！"
    format8 = "恭喜%v获得了特权4权限！领取魅力值到商场兑换心仪的礼品！"
)

type RandBugleManager struct {
	names []string
}

var randBugleManager *RandBugleManager

func init() {
	randBugleManager = &RandBugleManager{}
}

func GetRandBugleManager() *RandBugleManager {
	return randBugleManager
}

func (m *RandBugleManager) Init() {
	//m.LoadNames()
	go func() {
		for {
		    if len(m.names) == 0 {
		        break
            }
		    name := m.names[rand.Int() % len(m.names)]
		    if name == "" {
		        continue
            }
		    time.Sleep(time.Duration(60+rand.Int()%180) * time.Second)

		    bugleType := rand.Int() % 8
		    msg := m.buildBugleMsg(name, bugleType)

			if msg == "" {
				continue
			}
			domainUser.GetPlayerManager().BroadcastClientMsg(int32(pb.MessageId_CHAT), util.BuildSysBugle(msg))
		}
	}()
}

func (m *RandBugleManager)buildBugleMsg(name string, bugleType int) string {
    msg := ""
    switch bugleType {
    case 0:
        randGold := rand.Int() % 80000 + 100000
        msg = fmt.Sprintf(format1, name, randGold)
    case 1:
        randGold := rand.Int() % 500000 + 800000
        msg = fmt.Sprintf(format2, name, randGold)
    case 2:
        randGold := rand.Int() % 5000000 + 8000000
        msg = fmt.Sprintf(format3, name, randGold)
    case 3:
        msg = fmt.Sprintf(format4, name)
    case 4:
        msg = fmt.Sprintf(format5, name)
    case 5:
        msg = fmt.Sprintf(format6, name)
    case 6:
        msg = fmt.Sprintf(format7, name)
    case 7:
        msg = fmt.Sprintf(format8, name)
    }
    glog.Info("buildBugleMsg ok. msg=",msg)
    return msg
}

func (m *RandBugleManager)LoadNames() {
	buff, err := ioutil.ReadFile("./names.conf")
	if err != nil {
        panic("open file failed!")
    }

    str := string(buff)
    names := strings.Split(str, "\r\n")
    for _, v := range names {
        m.names = append(m.names, v)
    }
    //glog.Info("LoadNames ok. names=",m.names)
}