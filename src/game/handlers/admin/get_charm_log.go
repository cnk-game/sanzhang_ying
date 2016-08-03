package admin

import (
	"github.com/golang/glog"
	"net/http"
	"util"
    "encoding/json"
)

type CharmExchange struct {
	Date          string
	ExpendCharm   int
}

type GameGift struct {
	Date       string
	GiftGold   int
}

type CharmInfo struct {
    AllCharm         int
    AllExpendCharm   int
    DayExchanged     []CharmExchange
    FinePool         []GameGift
    BadPool          []GameGift
}

func GetCharmLogHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
	start := r.FormValue("start")
	end := r.FormValue("end")
	glog.Info("GetCharmLogHandler in,", start, end)

    ret := CharmInfo{}

    charmPool, _ := util.MongoLog_GetCharmPool()
    ret.AllCharm = charmPool.Charm
    ret.AllExpendCharm = charmPool.ExpendCharm

    charmExchange := util.MongoLog_GetCharmExchangeLog(start, end)
    for _, v := range charmExchange {
        exchange := CharmExchange{v.Date, v.ExpendCharm}
        ret.DayExchanged = append(ret.DayExchanged, exchange)
    }

    giftFine := util.MongoLog_GetGameGiftFineLog(start, end)
    for _, v := range giftFine {
        gameGift := GameGift{v.Date, v.GiftGold}
        ret.FinePool = append(ret.FinePool, gameGift)
    }

    giftBad := util.MongoLog_GetGameGiftBadLog(start, end)
    for _, v := range giftBad {
        gameGift := GameGift{v.Date, v.GiftGold}
        ret.BadPool = append(ret.BadPool, gameGift)
    }

    ret_str, _ := json.Marshal(ret)

    w.Write([]byte(ret_str))
}
