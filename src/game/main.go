package main

import (
	"config"
	"flag"
	activeUser "game/domain/activedata"
	domainconfig "game/domain/config"
	"game/domain/exchange"
	forbidWords "game/domain/forbidWords"
	"game/domain/game"
	iosActive "game/domain/iosActive"
	"game/domain/littlegame"
	newUserTask "game/domain/newusertask"
	domainPay "game/domain/pay"
	domainPrize "game/domain/prize"
	"game/domain/randBugle"
	"game/domain/rankingList"
	"game/domain/slots"
	domainSlot "game/domain/slots"
	"game/domain/stats"
	"game/domain/user"
	domainUser "game/domain/user"
	"game/handlers"
	handlerPay "game/handlers/pay"
	"game/server"
	"github.com/golang/glog"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	flag.Parse()
	glog.Info("===>启动游戏服务器")

	initialize()

	runtime.GOMAXPROCS(runtime.NumCPU())

	go stopHandler(server.GetServerInstance().GetSigChan())

	server.GetServerInstance().StartServer(handlers.GetMsgRegistry())
}

func stopHandler(c chan os.Signal) {
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	<-c

	rankingList.GetRankingList().Save()
	domainUser.GetRankingListUpdater().DoUpdate()
	stats.GetMatchLogManager().SaveAllMatchLogs()
	stats.GetAiFortuneLogManager().SaveLog()

	server.GetServerInstance().SetRefuseService()
	server.GetServerInstance().WaitStopServer()
	glog.Flush()
	os.Exit(0)
}

func initialize() {
	config.GetConfigManager().Init()
	config.GetCardConfigManager().Init()
	config.GetVipPriceConfigManager().Init()
	domainPrize.GetOnlinePrizeManager().Init()
	domainPrize.GetVipPrizeManager().Init()
	domainPrize.GetTaskPrizeManager().Init()
	domainPrize.GetExchangeGoodsManager().Init()
	slots.GetWheelConfigManager().Init()
	stats.GetMatchLogManager().Init()
	stats.GetAiFortuneLogManager().Init()
	domainSlot.Init()
	forbidWords.InitForbidWords()
	domainPay.Init()
	iosActive.GetIosActiveManager().Init()
	iosActive.GetUserIosActiveManager().Init()
	config.GetIosOnlineConfigManager().Init()
	activeUser.GetActiveManager().Init()        //wjs newActive
	exchange.GetGoodManager().Init()            //wjs exchange Goods
	domainconfig.GetMatchConfigManager().Init() //wjs match config
	littlegame.GetCardConfigManager().Init()    //wjs litgame config
	//littlegame.GetLigameLogicManager().Test()   //wjs gailvTest

	user.GetUserFortuneManager().UpdateGoldInGameFunc = game.GetDeskManager().UpdateGoldInGame

	rankingList.GetRankingList().Init()
	randBugle.GetRandBugleManager().Init()
	handlerPay.InitPayCode()
	newUserTask.GetNewUserTaskManager().Init()
	saveOnlineLog()
	domainPay.GetPayActiveConfig()
}

func saveOnlineLog() {
	go func() {
		for {
			stats.SaveOnlineLog(domainUser.GetPlayerManager().GetOnlineCount())
			time.Sleep(5 * time.Minute)
		}
	}()
}
