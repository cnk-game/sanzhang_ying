package game

import (
	"code.google.com/p/goprotobuf/proto"
	"config"
	"fmt"
	newUserTask "game/domain/newusertask"
	"game/domain/stats"
	domainUser "game/domain/user"
	"github.com/golang/glog"
	"math/rand"
	"pb"
	"sort"
	"sync"
	"time"
	"util"
)

const (
	Flower_M_Id = 1
	Eggs_M_Id   = 2
	Cheers_M_Id = 3
	Shoe_M_Id   = 4
	Kiss_M_Id   = 5
	Bomb_M_Id   = 6
	Flower_F_Id = 7
	Eggs_F_Id   = 8
	Cheers_F_Id = 9
	Shoe_F_Id   = 10
	Kiss_F_Id   = 11
	Bomb_F_Id   = 12
)

var (
	TurnExpiredSeconds        = 15
	OfflineWaitSeconds        = 30
	MinBetGold                = 100
	MaxBetGold                = 1000
	MaxRound                  = 20
	MaxChangeCardTimes        = 3
	MaxForbidCompareCardTimes = 5
	SNGMinCount               = 5
	WanRenPondMinBet          = 20000
)

type GamePlayer struct {
	User               *pb.UserDef
	Pos                int
	IsPlaying          bool
	IsSeenCard         bool
	IsGiveUp           bool
	IsFailed           bool
	IsCompared         bool // 是否参与比牌，游戏结束时显示牌
	AllIn              bool
	BetGold            int       // 当前下注
	TimeoutTimes       int       // 超时次数
	IsOffline          bool      // 是否离线
	OfflineTime        time.Time // 离线开始时间
	IsUseDoubleCard    bool      // 是否使用翻倍卡（四倍）
	UseForbidCardRound int       // 使用禁止比牌卡回合数（持续5回合）
	ChangeCardTimes    int       // 换牌次数
	IsReady            bool      // sng比赛准备
	IsEnterBackground  bool      // 客户端是否进入后台，此时不向其发送游戏内消息
	IsAllInLook        bool
	AllInGold          int //最后全下的下注
}

func (p *GamePlayer) resetPlayerStatus() {
	p.IsPlaying = false
	p.IsSeenCard = false
	p.IsGiveUp = false
	p.IsFailed = false
	p.IsCompared = false
	p.AllIn = false
	p.BetGold = 0
	p.IsUseDoubleCard = false
	p.UseForbidCardRound = 0
	p.ChangeCardTimes = 0
	p.IsAllInLook = false
	p.AllInGold = 0
}

func (p *GamePlayer) BuildMessage(betPond int) *pb.UserPokerDeskDef {
	msg := &pb.UserPokerDeskDef{}

	msg.BaseInfo = p.User
	msg.PutIntoGold = proto.Int(p.BetGold)
	msg.IsPlaying = proto.Bool(p.IsPlaying)
	msg.IsSeenCard = proto.Bool(p.IsSeenCard)
	msg.IsGiveUp = proto.Bool(p.IsGiveUp)
	msg.IsFailed = proto.Bool(p.IsFailed)
	msg.Pos = proto.Int(p.Pos + 1)
	msg.IsUseDoubleCard = proto.Bool(p.IsUseDoubleCard)
	msg.ForbidCompareCardRound = proto.Int(p.UseForbidCardRound)
	msg.UseReplaceCardCount = proto.Int(p.ChangeCardTimes)
	msg.IsReady = proto.Bool(p.IsReady)
	msg.BetPond = proto.Int(betPond)

	return msg
}

func isBadLuckyPlayer(userId string) bool {
	switch userId {
	case "55053a421d4bd445b10008ce":
		return true
	case "55018dec1d4bd445b100067d":
		return true
	case "54cc5fc81d4bd45e75000084":
		return true
	case "551628401d4bd4280e00054a":
		return true
	case "55193ea21d4bd466850088ab":
		return true
	case "55109b891d4bd44515000ede":
		return true
	case "551e0bf21d4bd434e500af59":
		return true
	case "54e325d21d4bd41d720003ff":
		return true
	case "55224b491d4bd434e5016918":
		return true
	case "5523e7421d4bd434e501b212":
		return true
	case "5524a97b1d4bd434e501c5a8":
		return true
	case "55274c341d4bd40abd002de1":
		return true
	case "55274bf61d4bd40abd002dc5":
		return true
	}
	return false
}

func (p *GamePlayer) calcLuckyValue(gameType util.GameType) int {
	luckyValue := 50
	if p.User.GetIsRobot() {
		luckyValue = 0
	} else {
		luckyValue += int(p.User.GetVipLevel()*5) + int(p.User.GetLuckyValue())
		f, ok := domainUser.GetUserFortuneManager().GetUserFortune(p.User.GetUserId())
		if ok && f.Gold >= 5000000 {
			if f.VipLevel <= 0 || time.Since(f.VipStartTime).Hours() > float64(f.VipValidDays*24) {
				if isBadLuckyPlayer(p.User.GetUserId()) {
					luckyValue = 1 + rand.Int()%5
				}
			}
		}
	}

	if luckyValue <= 0 {
		luckyValue = 1
	}

	glog.V(2).Info("====>计算幸运值nickname:", p.User.GetNickName(), " robot:", p.User.GetIsRobot(), " pos:", p.Pos, " luckyValue:", luckyValue)

	return luckyValue
}

type GameItem struct {
	sync.RWMutex
	GameType             util.GameType // 游戏模式
	GameId               int
	Logic                *GameLogic
	Players              map[string]*GamePlayer
	Lookup               map[string]*pb.UserDef
	BetPond              map[string]int            // 万人场奖池
	TotalBetPond         int                       // 万人场本局奖池
	BetPondItems         map[string]map[string]int // 万人场旁观下注
	WaitQueue            []string
	IsStart              bool      // 游戏是否已经开始
	ZhuangJia            int       // 庄家
	CurTurn              int       // 当前回合
	CurTurnUserId        string    // 当前回合玩家
	CurTurnTime          time.Time // 当前回合开始时间
	CurRound             int       // 第几轮(第二轮开始可以比牌)
	WinnerUserId         string
	WinnerUser           *GamePlayer
	SingleBetGold        int // 本轮单注
	TotalBetGold         int // 总下注
	TotalUserBetGold     int // 玩家总下注（排除机器人）
	AllInGold            int
	Ante                 int // 底注
	MatchTimes           int // 比赛场次
	WaitGameEnd          bool
	SNGGameEnd           bool
	WanRenWaitingBet     bool
	WanRenWaitingBetTime time.Time
	WaitQueueOrder       int
	MaxCardUserId        string
	EndTime              time.Time
	BossRobotUserId      string // boss机器人
	CompareTime          time.Time
	IsEmpty              bool
	EmptyUserId          string
}

type ContinuousLoseTimesComp struct {
	pos      []int
	item     *GameItem
	gameType util.GameType
}

func (p ContinuousLoseTimesComp) Len() int { return len(p.pos) }

func (p ContinuousLoseTimesComp) Less(i, j int) bool {
	times1 := 0
	times2 := 0
	for _, item := range p.item.Players {
		if item.Pos == p.pos[i] {
			times1 = item.calcLuckyValue(p.gameType)
		}
		if item.Pos == p.pos[j] {
			times2 = item.calcLuckyValue(p.gameType)
		}
	}

	return times1 > times2
}

func (p ContinuousLoseTimesComp) Swap(i, j int) { p.pos[i], p.pos[j] = p.pos[j], p.pos[i] }

func (game *GameItem) resetGameStatus() bool {
	if util.IsGameTypeCommon(game.GameType) || util.IsGameTypeProps(game.GameType) {
		c, ok := config.GetMatchConfigManager().GetMatchConfig(int(game.GameType))
		if !ok {
			return false
		}
		game.SingleBetGold = c.SingleBet
	}
	game.IsStart = false
	game.WaitGameEnd = false
	game.CurRound = 0
	game.WinnerUserId = ""
	game.WinnerUser = nil
	game.AllInGold = 0
	game.TotalBetGold = 0
	game.TotalUserBetGold = 0
	game.SNGGameEnd = false
	game.Ante = game.SingleBetGold
	game.MaxCardUserId = ""
	game.BossRobotUserId = ""
	game.Logic.bossPos = -1
	game.IsEmpty = false
	game.EmptyUserId = ""

	if util.IsGameTypeSNG(game.GameType) {
		c, ok := config.GetMatchConfigManager().GetSNGMatchConfig(int(game.GameType))
		if !ok {
			return false
		}
		game.SingleBetGold = c.SingleBet
		if game.MatchTimes <= 0 {
			game.Ante = c.AnteIncGold
		} else {
			game.Ante = game.MatchTimes * c.AnteIncGold
		}
	}

	if util.IsGameTypeWanRen(game.GameType) {
		c := config.GetMatchConfigManager().GetWanRenConfig()
		game.SingleBetGold = c.SingleBet
		game.Ante = 500000
	}

	glog.V(2).Info("===>gameType:", game.GameType, " 单注:", game.SingleBetGold, " 底注:", game.Ante)

	return true
}

func NewGameItem(gameType util.GameType, gameId int, ch chan int) *GameItem {
	item := &GameItem{}
	item.GameType = gameType
	item.GameId = gameId
	item.Players = make(map[string]*GamePlayer)
	item.Lookup = make(map[string]*pb.UserDef)
	item.BetPond = make(map[string]int)
	item.BetPondItems = make(map[string]map[string]int)
	item.Logic = NewGameLogic(int(gameType))
	item.resetGameStatus()

	go item.gameLoop()
	go item.reportCount(ch)

	return item
}

func (game *GameItem) reportCount(ch chan int) {
	for {
		ch <- game.getPlayerCount()
	}
}

func (game *GameItem) enterGame(user *pb.UserDef) bool {
	game.LockItem("enterGame")
	defer game.UnlockItem("enterGame")

	if util.IsGameTypeWanRen(game.GameType) {
		// 观看
		return game.lookup(user)
	} else {
		return game.enterDesk(user)
	}
}

func (game *GameItem) lookup(user *pb.UserDef) bool {
	if user == nil {
		return false
	}

	game.Lookup[user.GetUserId()] = user

	msg := &pb.MsgGetPokerDeskInfoRes{}
	msg.DeskInfo = game.BuildMessage()
	msg.Type = util.ToMatchType(game.GameType)

	for _, p := range game.Players {
		items := game.BetPondItems[p.User.GetUserId()]
		glog.V(2).Info("===>旁观下注userId:", p.User.GetUserId(), " items:", items)
		if items != nil {
			betGold := items[user.GetUserId()]
			if betGold > 0 {
				betItemMsg := &pb.WanRenLookupBetItem{}
				betItemMsg.BetUserId = proto.String(p.User.GetUserId())
				betItemMsg.BetGold = proto.Int(betGold)
				msg.Items = append(msg.Items, betItemMsg)
			}
		}
	}

	domainUser.GetPlayerManager().SendClientMsg(user.GetUserId(), int32(pb.MessageId_GET_POKER_DESK_INFO), msg)

	return true
}

func (game *GameItem) enterDesk(user *pb.UserDef) bool {
	if len(game.Players) >= int(util.MaxPlayerCount) {
		glog.V(2).Info("===>游戏已超过最大人数gameId:", game.GameId, " type:", game.GameType, " userId:", user.GetUserId())
		return false
	}

	// 判断是否已经进入
	for _, p := range game.Players {
		if p.User.GetUserId() == user.GetUserId() {
			glog.V(2).Info("===>****玩家已经进入游戏userId:", user.GetUserId())
			return false
		}
	}

	player := &GamePlayer{}
	player.User = user

	for i := 0; i < int(util.MaxPlayerCount); i++ {
		exist := false
		for _, p := range game.Players {
			if p.Pos == i {
				exist = true
				break
			}
		}
		if !exist {
			player.Pos = i
			break
		}
	}
	game.Players[user.GetUserId()] = player

	glog.V(2).Info("====>进入游戏username:", user.GetNickName(), " pos:", player.Pos)

	msg := &pb.MsgGetPokerDeskInfoRes{}
	msg.DeskInfo = game.BuildMessage()
	if len(game.Players) <= 2 {
		msg.DeskInfo.LastEndTime = proto.Int64(0)
	}
	msg.Type = util.ToMatchType(game.GameType)
	domainUser.GetPlayerManager().SendClientMsg(user.GetUserId(), int32(pb.MessageId_GET_POKER_DESK_INFO), msg)

	broMsg := &pb.MsgEnterPokerDeskBro{}
	broMsg.User = player.BuildMessage(0)

	broMsg.Type = util.ToMatchType(game.GameType)

	game.broadcastExcept(int32(pb.MessageId_ENTER_POKER_DESK_BRO), broMsg, []string{user.GetUserId()})

	return true
}

func (game *GameItem) LeaveGame(userId string, changeDesk bool) {
	game.LockItem("LeaveGame")
	defer game.UnlockItem("LeaveGame")

	game.leaveGame(userId, changeDesk, false, false)
}

//wjs 修改被踢时不离桌问题
func (game *GameItem) KickedLeaveGame(userId string, changeDesk bool) {
	game.LockItem("LeaveGame")
	defer game.UnlockItem("LeaveGame")

	game.leaveGame(userId, changeDesk, true, false)	
}


func (game *GameItem) LockItem(methodname string) {
	game.Lock()
}

func (game *GameItem) UnlockItem(methodname string) {
	game.Unlock()
}

func (game *GameItem) leaveGame(userId string, changeDesk, kickout, timeout bool) {
	// 1.结算
	// 2.离开
	p := game.Players[userId]
	if p == nil && !util.IsGameTypeWanRen(game.GameType) {
		return
	}

	msg := &pb.MsgLeavePokerDesk{}
	msg.UserId = proto.String(userId)
	if !util.IsGameTypeSNG(game.GameType) {
		msg.Kickout = proto.Bool(kickout)
		msg.Timeout = proto.Bool(timeout)
	}
	msg.Type = util.ToMatchType(game.GameType)

	if changeDesk {
		game.broadcastExcept(int32(pb.MessageId_LEAVE_POKER_DESK), msg, []string{userId})
	} else {
		game.broadcast(int32(pb.MessageId_LEAVE_POKER_DESK), msg)
	}

	delete(game.Players, userId)

	if p != nil && game.CurTurn == p.Pos && !timeout {
		game.switchNextTurn()
		game.sendPlayerRoundBeginMsg()
	}

	if len(game.Players) == 1 {
		for _, p := range game.Players {
			game.WinnerUserId = p.User.GetUserId()
			game.WinnerUser = p
			break
		}
	}

	// 离开游戏时，游戏未结束，且玩家参与本局游戏
	if game.IsStart && p != nil && p.BetGold > 0 {
		glog.V(2).Info("===>游戏未结束时提前离场，生成比赛记录gameId:", game.GameId, " gameType:", game.GameType, " userId:", userId)
		matchResult := &pb.MsgMatchResult{}
		matchResult.MatchType = proto.Int(int(game.GameType))

		if game.Logic.CompareCards2(game.Logic.GetCardsInt32(p.Pos), p.User.GetMatchRecord().GetMaxCards()) {
			matchResult.MaxCards = game.Logic.GetCardsInt32(p.Pos)
			p.User.GetMatchRecord().MaxCards = matchResult.MaxCards
		}
		matchResult.CardType = proto.Int(GetCardType(game.Logic.GetCards(p.Pos)))
		matchResult.EarnGold = proto.Int(-p.BetGold)
		domainUser.GetPlayerManager().SendServerMsg("", []string{userId}, int32(pb.ServerMsgId_MQ_MATCH_RESULT), matchResult)
	}

	if util.IsGameTypeWanRen(game.GameType) {
		delete(game.Lookup, userId)
		go game.onLeaveWaitQueue(userId)
	}

	if len(game.Players) == 1 && util.IsGameTypeSNG(game.GameType) && game.MatchTimes > 0 {
		game.MatchTimes = 0
		game.SNGGameEnd = true
		glog.V(2).Info("====>玩家离开，只剩一人,结束game.WinnerUserId:", game.WinnerUserId)
		go game.onEndGame()
	}
}

func (game *GameItem) broadcast(msgId int32, body proto.Message) {
	dstIds := []string{}
	for _, item := range game.Players {
		if item.IsEnterBackground {
			continue
		}
		dstIds = append(dstIds, item.User.GetUserId())
	}

	for _, item := range game.Lookup {
		dstIds = append(dstIds, item.GetUserId())
	}

	if len(dstIds) <= 0 {
		return
	}

	glog.V(2).Info("==>broadcast dstIds:", dstIds, " msgId:", util.GetMsgIdName(msgId), " msg:", body)
	domainUser.GetPlayerManager().SendClientMsg2(dstIds, msgId, body)
}

func (game *GameItem) broadcastExcept(msgId int32, body proto.Message, exceptUserIds []string) {
	dstIds := []string{}
	for _, item := range game.Players {
		if item.IsEnterBackground {
			continue
		}
		if !util.ContainsStr(exceptUserIds, item.User.GetUserId()) {
			dstIds = append(dstIds, item.User.GetUserId())
		}
	}

	for _, item := range game.Lookup {
		dstIds = append(dstIds, item.GetUserId())
	}

	if len(dstIds) <= 0 {
		return
	}

	glog.V(2).Info("==>broadcast dstIds:", dstIds, " msgId:", util.GetMsgIdName(msgId), " msg:", body)
	domainUser.GetPlayerManager().SendClientMsg2(dstIds, msgId, body)
}

func (game *GameItem) calcTotalBetGold(exceptUserId string) int {
	totalGold := 0
	for _, item := range game.Players {
		if item.User.GetUserId() != exceptUserId {
			totalGold += item.BetGold
		}
	}
	return totalGold
}

func (game *GameItem) BuildMessage() *pb.PokerDeskDef {
	msg := &pb.PokerDeskDef{}

	for _, item := range game.Players {
		msg.UserList = append(msg.UserList, item.BuildMessage(game.BetPond[item.User.GetUserId()]))
	}

	msg.PutIntoTotalGold = proto.Int(game.TotalBetGold)
	msg.SinglePutIntoGold = proto.Int(game.SingleBetGold)
	msg.CurrentRound = proto.Int(game.CurRound)
	msg.MaxRound = proto.Int(MaxRound)
	msg.TotalBetPond = proto.Int(game.TotalBetPond)
	msg.WanRenWaitingBetTime = proto.Int64(game.WanRenWaitingBetTime.Unix())
	msg.CurTurnUserId = proto.String(game.CurTurnUserId)
	msg.CurTurnTime = proto.Int64(game.CurTurnTime.Unix())
	msg.LastEndTime = proto.Int64(game.EndTime.Unix())

	order := 1
	for _, userId := range game.WaitQueue {
		p := game.Lookup[userId]
		if p == nil {
			continue
		}

		waitUser := &pb.WanRenWaitQueueUserDef{}
		waitUser.User = p
		waitUser.Order = proto.Int(order)
		msg.WaitQueue = append(msg.WaitQueue, waitUser)

		order++
	}

	return msg
}

func (game *GameItem) gameLoop() {
	if !game.IsStartGame() {
		// 游戏未开始
		time.Sleep(2 * time.Second)
	}

	for {
		if !game.IsStartGame() {
			// 检测是否可以开始
			game.StartGame()
		} else {
			// 游戏已经开始
			// 1.检测超时
			game.SwitchNextTurn()
			// 2.检测游戏结束
			if game.checkEnd() {
				// 游戏结束
				if game.getPlayerCount() > 1 && game.IsWaitEndGame() {
					time.Sleep(3 * time.Second)
				}

				game.onEndGame()
				time.Sleep(7 * time.Second)
			}
		}

		time.Sleep(1 * time.Second)
	}
}

func (game *GameItem) IsWaitEndGame() bool {
	game.LockItem("IsWaitEndGame")
	defer game.UnlockItem("IsWaitEndGame")
	return game.WaitGameEnd
}

func (game *GameItem) IsStartGame() bool {
	game.RLock()
	defer game.RUnlock()

	return game.IsStart
}

func (game *GameItem) StartGame() bool {
	game.LockItem("StartGame")
	defer game.UnlockItem("StartGame")

	if game.GameType == util.GameType_SNG_Level_3 {
		now := time.Now()
		if now.Hour() < 20 || now.Hour() > 22 {
			// 踢出所有玩家
			for _, p := range game.Players {
				game.KickOutOfGame(p.User.GetUserId(), false)
			}
			return false
		}
	}

	// 1.踢出金币不足的玩家
	if util.IsGameTypeWanRen(game.GameType) {
		if game.WanRenWaitingBet {
			glog.V(2).Info("====>万人场押注阶段")
			return false
		}
	}

	if util.IsGameTypeSNG(game.GameType) {
		if game.MatchTimes <= 0 && len(game.Players) < SNGMinCount {
			return false
		}

		c, ok := config.GetMatchConfigManager().GetSNGMatchConfig(int(game.GameType))
		if !ok {
			return false
		}

		readyCount := 0
		for _, p := range game.Players {
			if p.IsReady {
				readyCount++
				continue
			}

			// 扣入场费
			_, _, ok = domainUser.GetUserFortuneManager().ConsumeGold(p.User.GetUserId(), int64(c.EnterCostGold), false, "SNG入场费")
			if !ok {
				glog.V(2).Info("==>扣款失败userId:", p.User.GetUserId(), " SNG入场费:", c.EnterCostGold)
				game.KickOutOfGame(p.User.GetUserId(), false)
				return false
			}
			p.User.Gold = proto.Int64(int64(c.EnterCostGold))
			p.IsReady = true
		}
		if game.MatchTimes <= 0 && readyCount < SNGMinCount {
			return false
		}
	}

	if util.IsGameTypeCommon(game.GameType) || util.IsGameTypeProps(game.GameType) {
		if len(game.Players) < 2 {
			return false
		}

		c, ok := config.GetMatchConfigManager().GetMatchConfig(int(game.GameType))
		if !ok {
			glog.Error("获取游戏配置失败gameType:", game.GameType)
			return false
		}

		//nextLevelEnterLimit, nextOk := config.GetMatchConfigManager().GetNextLevelEnterLimit(int(game.GameType))

		for _, p := range game.Players {
			f, ok := domainUser.GetUserFortuneManager().GetUserFortune(p.User.GetUserId())
			if !ok {
				glog.V(2).Info("===>游戏开始检测，查询用户财富信息失败userId:", p.User.GetUserId())
				game.KickOutOfGame(p.User.GetUserId(), false)
				continue
			}
			if f.Gold < int64(c.KickoutLimit) {
				game.KickOutOfGame(p.User.GetUserId(), false)
				game.tipsCheck(p)
			}
			/*if nextOk && f.Gold >= int64(nextLevelEnterLimit) {
				game.KickOutOfGame(p.User.GetUserId(), false)
			}*/
		}

		if len(game.Players) < 2 {
			return false
		}
	}

	if util.IsGameTypeWanRen(game.GameType) {
		// 万人场
		//		glog.V(2).Info("==>万人场当前人数len:", len(game.Players), " waitQueue:", len(game.WaitQueue))
		if len(game.Players)+len(game.WaitQueue) < util.MaxPlayerCount {
			return false
		}

		c := config.GetMatchConfigManager().GetWanRenConfig()

		for _, p := range game.Players {
			glog.V(2).Info("===>万人场查询玩家财富信息userId:", p.User.GetUserId())
			f, ok := domainUser.GetUserFortuneManager().GetUserFortune(p.User.GetUserId())
			if !ok {
				glog.V(2).Info("===>游戏开始检测，查询用户财富信息失败userId:", p.User.GetUserId())
				game.KickOutOfGame(p.User.GetUserId(), false)
				continue
			}
			if f.Gold < int64(c.KickoutLimit) {
				game.KickOutOfGame(p.User.GetUserId(), false)
			}
		}

		glog.V(2).Info("==>万人场当前人数222222len:", len(game.Players), " waitQueue:", len(game.WaitQueue))
		if len(game.Players)+len(game.WaitQueue) < util.MaxPlayerCount {
			return false
		}

		waitCount := util.MaxPlayerCount - len(game.Players)
		for i := 0; i < waitCount; i++ {
			if len(game.WaitQueue) <= 0 {
				continue
			}
			userId := game.WaitQueue[0]
			game.WaitQueue = game.WaitQueue[1:]
			p := game.Lookup[userId]
			if p == nil {
				continue
			}
			delete(game.Lookup, userId)
			game.enterDesk(p)
		}

		glog.V(2).Info("==>万人场当前人数33333len:", len(game.Players), " waitQueue:", len(game.WaitQueue))
		if len(game.Players) < util.MaxPlayerCount {
			return false
		}
	}

	game.resetGameStatus()

	go func() {
		game.WanRenWaitingBet = true
		defer func() {
			game.WanRenWaitingBet = false
		}()

		if util.IsGameTypeWanRen(game.GameType) {
			game.WanRenWaitingBetTime = time.Now()
			game.broadcast(int32(pb.MessageId_LOOKUP_BEGIN_BET), nil)
			time.Sleep(30 * time.Second)
		}
		game.LockItem("StartGame go")
		defer game.UnlockItem("StartGame go")

		// 扣锅底
		for _, p := range game.Players {
			p.resetPlayerStatus()

			if util.IsGameTypeSNG(game.GameType) {
				glog.V(2).Info("====>username:", p.User.GetNickName(), " gold:", p.User.GetGold(), " ante:", game.Ante)
				if !p.IsReady {
					glog.V(2).Info("===>玩家username:", p.User.GetNickName(), " 未准备")
					return
				}

				if int(p.User.GetGold()) < game.Ante {
					glog.V(2).Info("===>踢出玩家username:", p.User.GetNickName())
					game.KickOutOfGame(p.User.GetUserId(), false)
				} else {
					p.User.Gold = proto.Int64(int64(int(p.User.GetGold()) - game.Ante))
					p.BetGold += game.Ante

					p.IsPlaying = true
					glog.V(2).Info("====>>>位置userName:", p.User.GetNickName(), " POS:", p.Pos)
				}
			} else {
				curGold, _, ok := domainUser.GetUserFortuneManager().ConsumeGold(p.User.GetUserId(), int64(game.Ante), false, "底注")
				if !ok {
					glog.V(2).Info("==>扣款失败userId:", p.User.GetUserId(), " 底注:", game.Ante)
					game.KickOutOfGame(p.User.GetUserId(), false)
					game.tipsCheck(p)
				} else {
					p.User.Gold = proto.Int64(curGold)
					p.BetGold += game.Ante

					p.IsPlaying = true
					glog.V(2).Info("====>>>位置userName:", p.User.GetNickName(), " POS:", p.Pos)
				}
			}
			game.TotalBetGold += game.Ante
			if !p.User.GetIsRobot() {
				game.TotalUserBetGold += game.Ante
			}
		}

		playingCount := game.getPlayingCount()

		if util.IsGameTypeSNG(game.GameType) {
			if game.MatchTimes <= 0 && playingCount < SNGMinCount {
				glog.V(2).Info("===>sng模式准备少于5人，count:", playingCount)
				return
			}
		} else {
			if playingCount < 2 {
				return
			}
		}

		// 锁定游戏
		if util.IsGameTypeSNG(game.GameType) {
			LockGame(game.GameId, true)
		}

		if _, ok := game.Players[game.CurTurnUserId]; !ok {
			game.ZhuangJia = -1
		}

		if game.ZhuangJia == -1 {
			// 首次，随机庄家
			pos := []int{}
			for _, p := range game.Players {
				if p.IsPlaying {
					pos = append(pos, p.Pos)
				}
			}
			if len(pos) < 2 {
				return
			}
			game.ZhuangJia = pos[rand.Int()%len(pos)]
			for _, p := range game.Players {
				if p.Pos == game.ZhuangJia {
					game.CurTurnUserId = p.User.GetUserId()
					break
				}
			}
			game.CurTurn = game.ZhuangJia
			game.CurTurnTime = time.Now()
		} else {
			glog.V(2).Info("****换庄,当前庄家:", game.ZhuangJia)
			game.CurTurn = game.ZhuangJia
			pos, ok := game.getNextTurn(game.CurTurn)
			if !ok {
				return
			}
			game.CurTurn = pos
			game.ZhuangJia = game.CurTurn
			glog.V(2).Info("****换庄结束,当前庄家:", game.ZhuangJia)
			for _, p := range game.Players {
				if p.Pos == game.ZhuangJia {
					game.CurTurnUserId = p.User.GetUserId()
					break
				}
			}
		}

		// 洗牌
		//		robotPositions := []int{}
		//
		//		comp := &ContinuousLoseTimesComp{}
		//		comp.gameType = game.GameType
		//		comp.item = game
		//		for _, p := range game.Players {
		//			comp.pos = append(comp.pos, p.Pos)
		//			if p.User.GetIsRobot() {
		//				robotPositions = append(robotPositions, p.Pos)
		//			}
		//		}
		//
		//		glog.V(2).Info("====>洗牌前位置按连续输次数未排序:", comp.pos)
		//		sort.Sort(comp)
		//		glog.V(2).Info("====>洗牌后位置按连续输次数倒序排序:", comp.pos)
		//		sort.Sort(util.IntRandSlice(comp.pos[1:]))
		//		glog.V(2).Info("====>洗牌后二次排序:", comp.pos)
		//		game.Logic.ShuffleCards(comp.pos, robotPositions)

		//		if util.IsGameTypeSNG(game.GameType) || util.IsGameTypeWanRen(game.GameType) {
		//			glog.V(2).Info("====>洗牌前位置按连续输次数未排序:", comp.pos)
		//			sort.Sort(comp)
		//			glog.V(2).Info("====>洗牌前位置按连续输次数倒序排序:", comp.pos)
		//			game.Logic.ShuffleCards(comp.pos, robotPositions)
		//		} else {
		//			sortedPos := game.sortedPlayers()
		//			glog.V(2).Info("===>新规则排序:", sortedPos)
		//			game.Logic.ShuffleCards(sortedPos, robotPositions)
		//		}

		robotPositions := []int{}
		for _, p := range game.Players {
			if p.User.GetIsRobot() {
				robotPositions = append(robotPositions, p.Pos)
			}
		}

		if len(robotPositions) > 0 {
			bossRobotPos := robotPositions[rand.Int()%len(robotPositions)]
			for _, p := range game.Players {
				if p.User.GetIsRobot() && p.Pos == bossRobotPos {
					winRate := stats.GetAiFortuneLogManager().GetWinRate(int(game.GameType))
					if rand.Int()%100 < winRate {
						game.BossRobotUserId = p.User.GetUserId()
						game.Logic.bossPos = bossRobotPos

						glog.Info("===>BOSS机器人胜率:", winRate, " userId:", game.BossRobotUserId, " pos:", p.Pos, " nickname:", p.User.GetNickName())
					}
				}
			}
		}

		sortedPos := game.sortedPlayers()
		//glog.Info("====>洗牌前未排序:", sortedPos)
		sort.Sort(util.IntRandSlice(sortedPos[1:]))
		//glog.Info("====>洗牌后排序:", sortedPos)
		game.Logic.ShuffleCards(sortedPos)

		game.IsStart = true
		game.CurTurnTime = time.Now()
		game.CurRound = 1

		if util.IsGameTypeSNG(game.GameType) {
			game.MatchTimes++
		}

		// 广播游戏开始
		msg := &pb.MsgGameStart{}
		msg.ZhuangJia = proto.String(game.CurTurnUserId)
		msg.SingleBetGold = proto.Int(game.SingleBetGold)
		msg.BeginBetGold = proto.Int(game.Ante)
		msg.MatchTimes = proto.Int(game.MatchTimes)
		msg.TotalBetPond = proto.Int(game.TotalBetPond)
		msg.Type = util.ToMatchType(game.GameType)

		bFindShowUser := false

		for _, p := range game.Players {
			if p.User.GetIsRobot() {
				msg.BossRobotUserId = proto.String(game.BossRobotUserId)
			} else {
				msg.BossRobotUserId = proto.String("")
			}
			if p.User.GetUserId() == ShowUser {
				bFindShowUser = true
			}
			domainUser.GetPlayerManager().SendClientMsg(p.User.GetUserId(), int32(pb.MessageId_POKER_DESK_GAME_BEGIN), msg)
		}

		if bFindShowUser == true {
			msgDesk := &pb.MsgDeskInfo{}
			for _, p := range game.Players {
				itemMsg := &pb.MsgDeskInfo_DeskDef{}
				itemMsg.UserId = proto.String(p.User.GetUserId())
				itemMsg.NickName = proto.String(p.User.GetNickName())
				itemMsg.Cards = game.Logic.GetCardsInt32(p.Pos)

				msgDesk.Desks = append(msgDesk.Desks, itemMsg)
			}

			domainUser.GetPlayerManager().SendClientMsg(ShowUser, int32(pb.MessageId_SEND_DESK_INFO), msgDesk)
		}

		glog.V(2).Info("****游戏开始msg:", msg)

		glog.V(2).Info("====>游戏开始庄家:", game.ZhuangJia, " UserId:", game.CurTurnUserId, " CurTurn:", game.CurTurn, " matchTimes:", game.MatchTimes)

		//		game.sendRobotMaxCardUser()
		game.sendRobotCards()

		game.sendPlayerRoundBeginMsg()
	}()

	return true
}

func (game *GameItem) tipsCheck(p *GamePlayer) {
	t := p.User.GetSubsidyPrizeTimes()
	if t < 3 {
		return
	}

	f, _ := domainUser.GetUserFortuneManager().GetUserFortune(p.User.GetUserId())
	if !f.FirstRecharge {
		go func(userId string) {
			timer := time.NewTimer(time.Second * 2)
			<-timer.C
			msg := &pb.MsgGoodsTips{}
			msg.ProductId = proto.String("100576")
			domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_GOODS_TIPS), msg)
		}(p.User.GetUserId())
	}

	return
}

func (game *GameItem) sortedPlayers() []int {
	offset := 0
	sorted := []int{}
	total := 0

	players := make(map[int]int)
	for _, p := range game.Players {
		if !p.IsPlaying {
			continue
		}
		players[p.Pos] = p.calcLuckyValue(game.GameType)
		glog.V(2).Info("===>幸运值pos:", p.Pos, " LuckyValue:", players[p.Pos])
	}

	for i := 0; i < len(players); i++ {
		total = 0
		for k, v := range players {
			exist := false
			for _, pos := range sorted {
				if k == pos {
					exist = true
					break
				}
			}
			if exist {
				continue
			}
			total += v
		}

		if total == 0 {
			for k := range players {
				exist := false
				for _, pos := range sorted {
					if k == pos {
						exist = true
						continue
					}
				}
				if !exist {
					sorted = append(sorted, k)
					break
				}
			}
			return sorted
		}

		r := rand.Int() % total

		offset = 0
		for k, v := range players {
			exist := false
			for _, pos := range sorted {
				if k == pos {
					exist = true
					continue
				}
			}
			if exist {
				continue
			}

			if r >= offset && r < offset+v {
				sorted = append(sorted, k)
			}
			offset += v
		}
	}

	return sorted
}

func (game *GameItem) SwitchNextTurn() {
	game.LockItem("SwitchNextTurn")
	defer game.UnlockItem("SwitchNextTurn")

	if !game.IsStart {
		glog.V(2).Info("游戏结束，切换回合结束gameId:", game.GameId)
		return
	}

	if time.Since(game.CurTurnTime).Seconds() > float64(TurnExpiredSeconds) {
		// 超时
		// 当前回合玩家当弃牌处理
		glog.V(2).Info("===>玩家UserId:", game.CurTurnUserId, " 超时")
		p := game.Players[game.CurTurnUserId]
		if p != nil {
			game.setPlayEnd(p, true)
			msg := &pb.MsgOpCardRes{}
			msg.Type = pb.CardOpType_GIVE_UP.Enum()
			msg.UserId = proto.String(game.CurTurnUserId)
			msg.MatchType = util.ToMatchType(game.GameType)
			game.broadcast(int32(pb.MessageId_PLAYER_OPERATE_CARDS), msg)
			p.TimeoutTimes++
			if !util.IsGameTypeSNG(game.GameType) && p.TimeoutTimes >= 2 {
				// 超时超过两次，踢出游戏
				glog.V(2).Info("==>超时超过2次，踢出游戏userId:", p.User.GetUserId())
				game.KickOutOfGame(p.User.GetUserId(), true)
			}
		}
		game.switchNextTurn()
		game.sendPlayerRoundBeginMsg()
	}
}

func (game *GameItem) setPlayEnd(p *GamePlayer, giveUp bool) {
	if game.getPlayingCount() > 1 {
		p.IsPlaying = false
		p.IsGiveUp = giveUp
	}
}

func (game *GameItem) getNextTurn(curTurn int) (int, bool) {
	players := []int{}
	for _, p := range game.Players {
		if p.IsPlaying {
			players = append(players, p.Pos)
		}
	}

	if len(players) <= 0 {
		return 0, false
	}

	sort.Ints(players)
	glog.Info("getNextTurn, players=", players)
	glog.Info("curTurn=", curTurn)

	for i, p := range players {
		if curTurn == p {
			return players[(i+1)%len(players)], true
		}
	}

	if curTurn > players[len(players)-1] {
		return players[0], true
	}

	if curTurn < players[0] {
		return players[0], true
	}

	for i, p := range players {
		if curTurn > p && curTurn < players[(i+1)%len(players)] {
			return players[(i+1)%len(players)], true
		}
	}

	glog.V(2).Info("===>getNextTurn failed curTurn:", curTurn, " players:", players)

	return 0, false
}

func (game *GameItem) switchNextTurn() {
	if !game.IsStart {
		glog.V(2).Info("游戏结束，切换回合结束gameId:", game.GameId)
		return
	}

	if game.getPlayingCount() <= 1 {
		return
	}

	pos, ok := game.getNextTurn(game.CurTurn)
	if !ok {
		return
	}

	var player *GamePlayer
	for _, p := range game.Players {
		if p.Pos == pos {
			player = p
		}
	}

	if player == nil {
		return
	}

	incRound := false
	if player.Pos == game.ZhuangJia {
		incRound = true
	} else {
		if game.CurTurn < game.ZhuangJia && player.Pos > game.ZhuangJia {
			// 当前小于庄&&下一个不等于庄
			incRound = true
		} else if game.CurTurn < game.ZhuangJia && player.Pos < game.CurTurn {
			incRound = true
		} else if game.CurTurn > game.ZhuangJia && player.Pos > game.ZhuangJia && player.Pos < game.CurTurn {
			// 当前大于庄&&下一个小于当前
			incRound = true
		}
	}

	if incRound {
		game.CurRound++
		if game.CurRound > 20 {
			game.CurRound = 20
		}
		glog.V(2).Info("===>当前第", game.CurRound, "轮")
	}

	glog.Info("====>gameId:", game.GameId, " CurTurn:", game.CurTurn, " player.pos:", player.Pos, " CurRound:", game.CurRound)

	game.CurTurn = player.Pos
	game.CurTurnTime = time.Now()
	game.CurTurnUserId = player.User.GetUserId()

	glog.Info("===>切换gameId:", game.GameId, " CurTurn:", game.CurTurn, " CurTurnUserId:", game.CurTurnUserId, " CurRound:", game.CurRound)
}

func (game *GameItem) sendPlayerRoundBeginMsg() {
	if !game.IsStart {
		return
	}
	if game.getPlayingCount() <= 1 {
		return
	}

	msg := &pb.MsgPlayerRoundBeginRes{}
	msg.UserId = proto.String(game.CurTurnUserId)
	msg.TurnTime = proto.Int64(game.CurTurnTime.Unix())
	msg.CurRound = proto.Int(game.CurRound)
	msg.Type = util.ToMatchType(game.GameType)
	game.broadcast(int32(pb.MessageId_PLAYER_ROUND_BEGIN), msg)

	glog.V(2).Info("====>ROUND_BEGIN gameId:", game.GameId, " 当前Round:", game.CurRound, " CurTurn:", game.CurTurn, " CurTurnUserId:", game.CurTurnUserId)
}

// add by wangsq start --- 踢人
func (game *GameItem) KickPlayer(from_userId string, target_userId string) bool {
	game.LockItem("KickPlayer")
	defer game.UnlockItem("KickPlayer")

	res := &pb.MsgKickPlayerRes{}
	res.FromUserId = proto.String(from_userId)
	res.TargetUserId = proto.String(target_userId)

	from_p := game.Players[from_userId]
	target_p := game.Players[target_userId]

	if from_p == nil || target_p == nil {
		glog.Error("KickPlayer can't find player.")
		res.Code = pb.MsgKickPlayerRes_FAILED.Enum()
		res.Reason = proto.String("通信错误")
		domainUser.GetPlayerManager().SendClientMsg(from_userId, int32(pb.MessageId_KICK_PLAYER), res)
		return false
	}

	fu, fu_ok := domainUser.GetUserFortuneManager().GetUserFortune(from_userId)
	tu, tu_ok := domainUser.GetUserFortuneManager().GetUserFortune(target_userId)
	if fu_ok && tu_ok {
		if fu.VipLevel == 0 || fu.VipLevel <= tu.VipLevel {
			glog.Infof("not enough vip level. %d, %d", fu.VipLevel, tu.VipLevel)
			res.Code = pb.MsgKickPlayerRes_FAILED.Enum()
			res.Reason = proto.String("您的Vip等级不足")
			domainUser.GetPlayerManager().SendClientMsg(from_userId, int32(pb.MessageId_KICK_PLAYER), res)
			return false
		}
	} else {
		glog.Error("KickPlayer can't find player info.")
		res.Code = pb.MsgKickPlayerRes_FAILED.Enum()
		res.Reason = proto.String("通信错误")
		domainUser.GetPlayerManager().SendClientMsg(from_userId, int32(pb.MessageId_KICK_PLAYER), res)
		return false
	}

	if target_p.IsPlaying {
		glog.Info("target user is Playing")
		res.Code = pb.MsgKickPlayerRes_FAILED.Enum()
		res.Reason = proto.String("对方正在游戏中")
		domainUser.GetPlayerManager().SendClientMsg(from_userId, int32(pb.MessageId_KICK_PLAYER), res)
		return false
	}

	res.Code = pb.MsgKickPlayerRes_OK.Enum()
	res.Reason = proto.String("Ok")
	game.broadcastExcept(int32(pb.MessageId_KICK_PLAYER), res, []string{target_userId})

	content := fmt.Sprintf("抱歉，您被玩家 %s 请出了房间", *from_p.User.NickName)
	go func(userId, content string) {
		timer := time.NewTimer(time.Second * 2)
		<-timer.C
		msg := &pb.MsgShowTips{}
		msg.UserId = proto.String(userId)
		msg.Content = proto.String(content)
		domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_SHOW_TIPS), msg)
	}(target_userId, content)

	return true
}

// add by wangsq end

func (game *GameItem) onOpCards(userId string, msg *pb.MsgOpCardReq) {
	if !game.IsStartGame() {
		glog.Info("==>游戏未开始，忽略操作包gameId:", game.GameId, " msg:", msg)
		return
	}

	game.LockItem("onOpCards")
	defer game.UnlockItem("onOpCards")

	p := game.Players[userId]
	if p == nil {
		return
	}

	if !p.IsPlaying {
		glog.V(2).Info("该玩家游戏已经结束userId:", userId)
		return
	}

	p.TimeoutTimes = 0

	res := &pb.MsgOpCardRes{}
	res.UserId = proto.String(userId)
	res.Type = msg.Type
	res.Gold = msg.Gold
	res.CompareUserId = msg.CompareUserId
	res.MatchType = util.ToMatchType(game.GameType)

	if game.CurTurnUserId != userId {
		// 只可进行弃牌，看牌操作
		if msg.GetType() == pb.CardOpType_GIVE_UP {
			game.setPlayEnd(p, true)
			game.broadcast(int32(pb.MessageId_PLAYER_OPERATE_CARDS), res)
		} else if msg.GetType() == pb.CardOpType_SEE_CARDS {
			p.IsSeenCard = true
			game.broadcastExcept(int32(pb.MessageId_PLAYER_OPERATE_CARDS), res, []string{userId})

			res.UserIdCards = game.Logic.GetCardsInt32(p.Pos)

			glog.V(2).Info("===>看牌userId:", userId, " cards:", res.UserIdCards)
			domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_PLAYER_OPERATE_CARDS), res)
		} else {
			glog.V(2).Info("===>无效操作，userId:", userId, "nickname:", p.User.GetNickName(), " curRound:", game.CurRound, "非本回合玩家")
		}
		return
	}

	if game.getAllInPlayer() != nil {
		// 已有全押玩家，其它只能选择全押，弃牌或看牌
		if msg.GetType() != pb.CardOpType_GIVE_UP && msg.GetType() != pb.CardOpType_ALL_IN &&
			msg.GetType() != pb.CardOpType_SEE_CARDS {
			return
		}
	}

	if game.CurRound >= MaxRound {
		// 只能比牌，全押，弃牌或看牌
		if msg.GetType() != pb.CardOpType_COMPARE &&
			msg.GetType() != pb.CardOpType_ALL_IN &&
			msg.GetType() != pb.CardOpType_GIVE_UP &&
			msg.GetType() != pb.CardOpType_SEE_CARDS {
			return
		}
	}

	// 可进行所有操作
	switch msg.GetType() {
	case pb.CardOpType_GIVE_UP:
		game.setPlayEnd(p, true)
		game.broadcast(int32(pb.MessageId_PLAYER_OPERATE_CARDS), res)

		if game.getPlayingCount() > 1 {
			game.switchNextTurn()
			game.sendPlayerRoundBeginMsg()
		}
	case pb.CardOpType_SEE_CARDS:
		if p.IsSeenCard {
			return
		}
		p.IsSeenCard = true
		game.broadcastExcept(int32(pb.MessageId_PLAYER_OPERATE_CARDS), res, []string{userId})

		res.UserIdCards = game.Logic.GetCardsInt32(p.Pos)

		glog.V(2).Info("===>看牌userId:", userId, " cards:", res.UserIdCards)
		domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_PLAYER_OPERATE_CARDS), res)

		// 重置当前时间
		game.CurTurnTime = time.Now()
		game.sendPlayerRoundBeginMsg()
	case pb.CardOpType_COMPARE:
		glog.Info("++++++++++++ CardOpType_COMPARE ")
		glog.Info("====>开始处理比牌 CurRound:", game.CurRound)
		if game.CurRound <= 0 {
			glog.V(2).Info("*****第二轮之后才可以比牌")
			return
		}

		if userId == msg.GetCompareUserId() {
			glog.V(2).Info("不能跟自己比牌")
			return
		}

		compareUser := game.Players[msg.GetCompareUserId()]
		if compareUser == nil {
			glog.Info("比牌对象不存在userId:", userId, " compareUserId:", msg.GetCompareUserId())
			return
		}

		if game.CurRound < MaxRound && compareUser.UseForbidCardRound > 0 && game.CurRound-compareUser.UseForbidCardRound < MaxForbidCompareCardTimes {
			glog.Info("比牌对象userId:", compareUser.User.GetUserId(), " 在第", compareUser.UseForbidCardRound, "使用了禁止比牌道具,当前轮:", game.CurRound)
			return
		}

		betGold := game.SingleBetGold
		if p.IsSeenCard {
			betGold *= 2
		}

		if compareUser.IsUseDoubleCard {
			betGold *= 4
		}

		// 扣费
		if util.IsGameTypeSNG(game.GameType) {
			if int(p.User.GetGold()) < betGold {
				betGold = int(p.User.GetGold())
			}
			p.User.Gold = proto.Int64(int64(int(p.User.GetGold()) - betGold))
			p.BetGold += betGold
		} else {
			curGold, consumeGold, ok := domainUser.GetUserFortuneManager().ConsumeGold(userId, int64(betGold), true, "比牌")
			if !ok {
				glog.V(2).Info("==>比牌失败，扣款失败userId:", userId, " betGold:", betGold)
				return
			}
			p.User.Gold = proto.Int64(curGold)
			p.BetGold += consumeGold
			glog.V(2).Info("==>比牌成功userId:", userId, " betGold:", betGold, " consumeGold:", consumeGold, " curGold:", curGold)

			betGold = consumeGold
		}
		res.Gold = proto.Int(betGold)
		game.TotalBetGold += betGold
		if !p.User.GetIsRobot() {
			game.TotalUserBetGold += betGold
		}

		p.IsCompared = true
		compareUser.IsCompared = true

		//glog.Info("===>处理比牌结束")
		game.WaitGameEnd = true

		if game.Logic.CompareCardsByPos(p.Pos, compareUser.Pos) {
			res.WinnerUserId = proto.String(userId)
			game.setPlayEnd(compareUser, false)
			compareUser.IsFailed = true
		} else {
			res.WinnerUserId = proto.String(msg.GetCompareUserId())
			game.setPlayEnd(p, false)
			p.IsFailed = true
		}
		//glog.Info("===>处理比牌广播")
		game.broadcast(int32(pb.MessageId_PLAYER_OPERATE_CARDS), res)

		game.CompareTime = time.Now()
		go func() {
			time.Sleep(3 * time.Second)

			game.LockItem("onOpCards go")
			defer game.UnlockItem("onOpCards go")

			game.CurTurnTime = time.Now()

			if game.getPlayingCount() > 1 {
				game.switchNextTurn()
				game.sendPlayerRoundBeginMsg()
			}
		}()
	case pb.CardOpType_ALL_IN:
		glog.V(2).Info("=======>全下userId:", userId)
		if game.getPlayingCount() != 2 {
			glog.V(2).Info("====>当前人数不等于2")
			return
		}
		if game.CurRound <= 2 {
			glog.V(2).Info("*****第二轮之后才可以全下")
			return
		}

		for _, item := range game.Players {
			userIdTem := item.User.GetUserId()
			if userIdTem != userId {
				coinsTemp := item.User.GetGold()
				if coinsTemp == 0 && item.IsFailed == false && item.IsPlaying == true && item.AllIn == false {
					reqTemp := &pb.MsgOpCardReq{}
					reqTemp.Type = pb.CardOpType_COMPARE.Enum()
					reqTemp.Gold = proto.Int32(0)
					reqTemp.CompareUserId = proto.String(userIdTem)

					res.Type = reqTemp.Type
					res.CompareUserId = reqTemp.CompareUserId

					compareUser := game.Players[reqTemp.GetCompareUserId()]
					if compareUser == nil {
						glog.Info("比牌对象不存在userId:", userId, " compareUserId:", reqTemp.GetCompareUserId())
						return
					}

					if game.CurRound < MaxRound && compareUser.UseForbidCardRound > 0 && game.CurRound-compareUser.UseForbidCardRound < MaxForbidCompareCardTimes {
						glog.Info("比牌对象userId:", compareUser.User.GetUserId(), " 在第", compareUser.UseForbidCardRound, "使用了禁止比牌道具,当前轮:", game.CurRound)
						return
					}

					betGold := game.SingleBetGold
					if p.IsSeenCard {
						betGold *= 2
					}

					if compareUser.IsUseDoubleCard {
						betGold *= 4
					}

					// 扣费
					if util.IsGameTypeSNG(game.GameType) {
						if int(p.User.GetGold()) < betGold {
							betGold = int(p.User.GetGold())
						}
						p.User.Gold = proto.Int64(int64(int(p.User.GetGold()) - betGold))
						p.BetGold += betGold
					} else {
						curGold, consumeGold, ok := domainUser.GetUserFortuneManager().ConsumeGold(userId, int64(betGold), true, "比牌")
						if !ok {
							glog.V(2).Info("==>比牌失败，扣款失败userId:", userId, " betGold:", betGold)
							return
						}
						p.User.Gold = proto.Int64(curGold)
						p.BetGold += consumeGold
						glog.V(2).Info("==>比牌成功userId:", userId, " betGold:", betGold, " consumeGold:", consumeGold, " curGold:", curGold)

						betGold = consumeGold
					}
					res.Gold = proto.Int(betGold)
					game.TotalBetGold += betGold
					if !p.User.GetIsRobot() {
						game.TotalUserBetGold += betGold
					}

					p.IsCompared = true
					compareUser.IsCompared = true

					game.WaitGameEnd = true

					if game.Logic.CompareCardsByPos(p.Pos, compareUser.Pos) {
						res.WinnerUserId = proto.String(userId)
						game.setPlayEnd(compareUser, false)
						compareUser.IsFailed = true
					} else {
						res.WinnerUserId = proto.String(reqTemp.GetCompareUserId())
						game.setPlayEnd(p, false)
						p.IsFailed = true
					}
					game.broadcast(int32(pb.MessageId_PLAYER_OPERATE_CARDS), res)

					game.CompareTime = time.Now()
					go func() {
						time.Sleep(3 * time.Second)

						game.LockItem("onOpCards go")
						defer game.UnlockItem("onOpCards go")

						game.CurTurnTime = time.Now()

						if game.getPlayingCount() > 1 {
							game.switchNextTurn()
							game.sendPlayerRoundBeginMsg()
						}
					}()
					return
				}
			}
		}

		///

		allInPlayer := game.getAllInPlayer()
		if allInPlayer == nil {
			// 自己为先全押方
			p.AllIn = true

			betGold, ok, isEmpty, pUserId := game.calcAllInGold()
			if !ok {
				glog.V(2).Info("==>计算全押金币数目失败")
				return
			}
			game.IsEmpty = isEmpty
			game.EmptyUserId = pUserId
			game.AllInGold = betGold
			glog.Info("===>game.AllInGold:", betGold)
			glog.Info("===>userId:", userId, " p.BetGold:", p.BetGold)

			// 扣费
			if betGold > 0 {
				if p.IsSeenCard {
					betGold *= 2
					p.IsAllInLook = true
					/*if isEmpty && pUserId == userId {
						betGold = game.AllInGold
					}*/

					if isEmpty && pUserId != userId {
						uOther := game.Players[pUserId]
						glog.Info("==>allin betGold:", betGold)
						if uOther.IsSeenCard {
							uOther.IsAllInLook = true
							//betGold = game.AllInGold
						}
						glog.Info("==>allin betGold:", betGold)
					}

				} else {
					glog.Info("==>allin betGold:", betGold)
					if isEmpty && pUserId != userId {
						uOther := game.Players[pUserId]
						glog.Info("==>allin betGold:", betGold)
						if uOther.IsSeenCard {
							uOther.IsAllInLook = true
							betGold = game.AllInGold
						}
						glog.Info("==>allin betGold:", betGold)
					}
				}

				glog.Info("==>allin betGold:", betGold)

				if util.IsGameTypeSNG(game.GameType) {
					if int(p.User.GetGold()) < betGold {
						betGold = int(p.User.GetGold())
					}
					p.User.Gold = proto.Int64(int64(int(p.User.GetGold()) - betGold))
				} else {
					curGold, consumeGold, ok := domainUser.GetUserFortuneManager().ConsumeGold(userId, int64(betGold), true, "全押")
					if !ok {
						glog.V(2).Info("==>扣款失败userId:", userId, " betGold:", betGold)
						return
					}
					p.User.Gold = proto.Int64(int64(curGold))
					glog.V(2).Info("==>扣款成功userId:", userId, " consumeGold:", consumeGold, " curGold:", curGold)

					betGold = consumeGold
				}
			}
			res.Gold = proto.Int(betGold)

			glog.V(2).Info("===>userId:", userId, " 全下数目:", betGold)

			p.BetGold += betGold
			p.AllInGold = betGold
			game.TotalBetGold += betGold
			glog.Info("===>userId:", userId, " betGold:", betGold)
			glog.Info("===>userId:", userId, " p.BetGold:", p.BetGold)
			glog.Info("===>userId:", userId, " p.AllInGold:", p.AllInGold)
			glog.Info("===>userId:", userId, " game.TotalBetGold:", game.TotalBetGold)
			if !p.User.GetIsRobot() {
				game.TotalUserBetGold += betGold
			}

			game.broadcast(int32(pb.MessageId_PLAYER_OPERATE_CARDS), res)

			game.switchNextTurn()
			game.sendPlayerRoundBeginMsg()
		} else {
			p.IsCompared = true
			allInPlayer.IsCompared = true
			glog.Info("===>userId:", userId, " p.BetGold:", p.BetGold)

			// 扣费
			betGold := game.AllInGold
			if betGold > 0 {
				if p.IsSeenCard {
					betGold *= 2

					/*if game.IsEmpty {
						if allInPlayer.IsAllInLook {
							betGold = game.AllInGold
						}
					}*/

				} /*else {
					betGold = game.AllInGold
					if game.IsEmpty {
						if allInPlayer.IsAllInLook {
							betGold = game.AllInGold
						}

						if game.EmptyUserId == userId {
							betGold = game.AllInGold
						}
					}
				}*/

				//allTempBetGold := allInPlayer.BetGold - allInPlayer.AllInGold
				//glog.Info("===++++allTempBetGold:", allTempBetGold)
				//deteT := allTempBetGold - p.BetGold
				//glog.Info("===>deteT:", deteT)
				//if p.BetGold != allTempBetGold {
				//	betGold += deteT
				//}

				glog.Info("===>betGold:", betGold)

				if util.IsGameTypeSNG(game.GameType) {
					if int(p.User.GetGold()) < betGold {
						betGold = int(p.User.GetGold())
					}
					p.User.Gold = proto.Int64(int64(int(p.User.GetGold()) - betGold))
				} else {
					curGold, consumeGold, ok := domainUser.GetUserFortuneManager().ConsumeGold(userId, int64(betGold), true, "全押")
					if !ok {
						glog.V(2).Info("==>扣款失败userId:", userId, " betGold:", game.AllInGold)
						return
					}
					p.User.Gold = proto.Int64(curGold)
					glog.V(2).Info("==>扣款成功userId:", userId, " consumeGold:", consumeGold, " curGold:", curGold)
					betGold = consumeGold
				}

				// 计算全押金币
				p.BetGold += betGold
				p.AllInGold = betGold
				game.TotalBetGold += betGold

				glog.Info("===>userId:", userId, " betGold:", betGold)
				glog.Info("===>userId:", userId, " p.BetGold:", p.BetGold)
				glog.Info("===>userId:", userId, " p.AllInGold:", p.AllInGold)
				glog.Info("===>userId:", userId, " game.TotalBetGold:", game.TotalBetGold)
				if !p.User.GetIsRobot() {
					game.TotalUserBetGold += betGold
				}

				res.Gold = proto.Int(betGold)
			}

			// 已经有全押方,比牌
			game.CompareTime = time.Now()

			isWin := false
			winnerUserId := ""
			if game.Logic.CompareCardsByPos(allInPlayer.Pos, p.Pos) {
				winnerUserId = allInPlayer.User.GetUserId()
				res.WinnerUserId = proto.String(allInPlayer.User.GetUserId())
				res.CompareUserId = proto.String(allInPlayer.User.GetUserId())
			} else {
				isWin = true
				winnerUserId = userId
				res.WinnerUserId = proto.String(userId)
				res.CompareUserId = proto.String(allInPlayer.User.GetUserId())
			}

			glog.Info("===>winnerUserId:", winnerUserId, " p.IsSeenCard:", p.IsSeenCard, " p.IsAllInLook ", p.IsAllInLook)

			if game.IsEmpty {
				if winnerUserId == userId && p.IsSeenCard && !p.IsAllInLook {
					//tempCoins := allInPlayer.AllInGold / 2
					tempCoins := allInPlayer.AllInGold - p.AllInGold
					if tempCoins > 0 {
						tempCoins = allInPlayer.AllInGold / 2
						game.TotalBetGold -= tempCoins
						allInPlayer.BetGold -= tempCoins
						curGold, ok := domainUser.GetUserFortuneManager().EarnGold(allInPlayer.User.GetUserId(), int64(tempCoins), "全下反退")
						if ok {
							allInPlayer.User.Gold = proto.Int64(curGold)
						}
						glog.Info("===>全下反退 userId:", allInPlayer.User.GetUserId())
					}
				} /*else if winnerUserId == allInPlayer.User.GetUserId() && allInPlayer.IsSeenCard {
					tempCoins := p.AllInGold / 2
					game.TotalBetGold -= tempCoins
					p.BetGold -= tempCoins
					curGold, ok := domainUser.GetUserFortuneManager().EarnGold(userId, int64(tempCoins), "全下反退")
					if ok {
						allInPlayer.User.Gold = proto.Int64(curGold)
					}

					glog.Info("===>全下反退 userId:", userId)
				}*/
			}
			game.broadcast(int32(pb.MessageId_PLAYER_OPERATE_CARDS), res)

			if isWin {
				game.setPlayEnd(allInPlayer, false)
			} else {
				game.setPlayEnd(p, false)
			}
			game.WaitGameEnd = true
		}
	case pb.CardOpType_RAISE:
		glog.V(2).Info("==>加注gold:", msg.GetGold(), " singleBetGold:", game.SingleBetGold)
		if int(msg.GetGold()) <= game.SingleBetGold {
			return
		}

		if game.getPlayingCount() <= 1 {
			return
		}

		if !config.GetMatchConfigManager().IsRaiseBetRight(int(game.GameType), int(msg.GetGold())) {
			glog.V(2).Info("===>加注失败，找不到对应注gameType:", game.GameType, " gold:", msg.GetGold())
			return
		}

		betGold := int(msg.GetGold())
		if p.IsSeenCard {
			betGold *= 2
		}
		res.Gold = proto.Int(betGold)

		if util.IsGameTypeSNG(game.GameType) {
			if int(p.User.GetGold()) < betGold {
				glog.V(2).Info("==>扣款失败userId:", userId, " betGold:", betGold)
				return
			}
			p.User.Gold = proto.Int64(int64(int(p.User.GetGold()) - betGold))
		} else {
			// 扣费
			curGold, _, ok := domainUser.GetUserFortuneManager().ConsumeGold(userId, int64(betGold), true, "加注")
			if !ok {
				glog.V(2).Info("==>扣款失败userId:", userId, " betGold:", betGold)
				return
			}
			p.User.Gold = proto.Int64(int64(curGold))
			glog.V(2).Info("==>扣款成功userId:", userId, " betGold:", betGold, " curGold:", curGold)
		}

		game.SingleBetGold = int(msg.GetGold())
		p.BetGold += betGold
		game.TotalBetGold += betGold
		if !p.User.GetIsRobot() {
			game.TotalUserBetGold += betGold
		}

		game.broadcast(int32(pb.MessageId_PLAYER_OPERATE_CARDS), res)

		game.switchNextTurn()
		game.sendPlayerRoundBeginMsg()
	case pb.CardOpType_FOLLOW:
		if game.getPlayingCount() <= 1 {
			return
		}

		betGold := game.SingleBetGold
		if p.IsSeenCard {
			betGold *= 2
		}

		if util.IsGameTypeSNG(game.GameType) {
			if int(p.User.GetGold()) < betGold {
				glog.V(2).Info("==>扣款失败userId:", userId, " betGold:", betGold)
				return
			}
			p.User.Gold = proto.Int64(int64(int(p.User.GetGold()) - betGold))
			glog.V(2).Info("===>跟注userId:", userId, " betGold:", betGold, " userGold:", p.User.GetGold())
		} else {
			// 扣费
			curGold, _, ok := domainUser.GetUserFortuneManager().ConsumeGold(userId, int64(betGold), false, "跟注")
			if !ok {
				glog.V(2).Info("==>扣款失败userId:", userId, " betGold:", betGold)
				return
			}
			p.User.Gold = proto.Int64(curGold)
			glog.V(2).Info("==>扣款成功userId:", userId, " betGold:", betGold, " curGold:", curGold)
		}

		p.BetGold += betGold
		game.TotalBetGold += betGold
		if !p.User.GetIsRobot() {
			game.TotalUserBetGold += betGold
		}

		res.Gold = proto.Int(betGold)
		game.broadcast(int32(pb.MessageId_PLAYER_OPERATE_CARDS), res)

		game.switchNextTurn()
		game.sendPlayerRoundBeginMsg()
	}
}

func (game *GameItem) checkEnd() bool {
	game.LockItem("checkEnd")
	defer game.UnlockItem("checkEnd")

	count := 0
	for _, p := range game.Players {
		if p.IsPlaying {
			count += 1
		}
	}

	return count <= 1
}

func (game *GameItem) onEndGame() {
	game.LockItem("onEndGame")
	defer game.UnlockItem("onEndGame")

	glog.V(2).Info("===>游戏结束gameId:", game.GameId)
	game.IsStart = false

	defer func() {
		game.resetGameStatus()
		game.EndTime = time.Now()
	}()

	msg := &pb.MsgPokerDeskGameEndRes{}

	winner := game.WinnerUser
	for _, p := range game.Players {
		glog.V(2).Info("===>Check winnerUserId:", game.WinnerUserId, " pUserId:", p.User.GetUserId())
		if game.WinnerUserId != "" && game.WinnerUserId == p.User.GetUserId() {
			winner = p
		}
		if p.IsPlaying {
			game.WinnerUserId = p.User.GetUserId()
			winner = p
			p.IsPlaying = false
		}
		glog.V(2).Info("===>游戏结束，下注userId:", p.User.GetUserId(), " gold:", p.BetGold)
	}

	glog.V(2).Info("===>游戏结束game.WinnerUserId:", game.WinnerUserId, " winner:", winner)

	if game.WinnerUserId == "" || winner == nil {
		game.MatchTimes = 0
		return
	}

	msg.Type = util.ToMatchType(game.GameType)

	msg.WinnerUserId = proto.String(game.WinnerUserId)
	winnerEarnGold := game.TotalBetGold
	glog.V(2).Info("===>总下注:", game.TotalBetGold, " winnerEarnGold:", winnerEarnGold, " winnerUserId:", game.WinnerUserId)

	// 发放奖励
	if !util.IsGameTypeSNG(game.GameType) {
		// 计算税收
		// 万人场奖池赢家抽10%,系统抽10%,剩余80%按旁观玩家押注比例比配
		if util.IsGameTypeWanRen(game.GameType) {
			items := game.BetPondItems[game.WinnerUserId]
			if items != nil && len(items) > 0 {
				winnerEarnGold += int(float64(game.TotalBetPond) * 0.1)
				remainBetPond := float64(game.TotalBetPond) * 0.8

				total := 0
				for _, v := range items {
					total += v
				}

				for k, v := range items {
					winGold := int(remainBetPond / float64(total) * float64(v))
					_, ok := domainUser.GetUserFortuneManager().EarnGold(k, int64(winGold), "万人场押注")
					if !ok {
						glog.Error("万人场押注奖励失败userId:", k, " winGold:", winGold)
					}
					msg := &pb.MsgLookupBetWin{}
					msg.WinGold = proto.Int(winGold)
					domainUser.GetPlayerManager().SendClientMsg(k, int32(pb.MessageId_LOOKUP_BET_WIN), msg)
				}
				game.TotalBetPond = 0
			}
			msg.TotalBetPond = proto.Int(game.TotalBetPond)
		}

		tax := 0
		if winner != nil {
			tax = int(float64(winnerEarnGold-winner.BetGold) * 0.05)
			tax -= tax % 100
			if tax < 100 {
				tax = 0
			}
		}

		glog.V(2).Info("************赢家userId:", game.WinnerUserId, " EarnGold:", winnerEarnGold-tax, " tax:", tax)
		curGold, ok := domainUser.GetUserFortuneManager().EarnGold(game.WinnerUserId, int64(winnerEarnGold-tax), "胜利")
		util.MongoLog_GameFee(tax)
		if !ok {
			glog.V(2).Info("==>获取金币失败userId:", game.WinnerUserId, " earnGold:", game.TotalBetGold)
			return
		}
		winner.User.Gold = proto.Int64(curGold)
	} else {
		winner.User.Gold = proto.Int64(int64(int(winner.User.GetGold()) + winnerEarnGold))
	}

	if winner != nil {
		if winner.User.GetLuckyValue() < 50 {
			winner.User.LuckyValue = proto.Int(0)
		} else {
			winner.User.LuckyValue = proto.Int(int(winner.User.GetLuckyValue()) - 50)
		}
		glog.V(2).Info("===>游戏结束赢家:", winner.User.GetNickName(), " pos:", winner.Pos, " luckyValue:", winner.User.GetLuckyValue())
	}

	// 胜利，获取2点经验
	domainUser.GetUserFortuneManager().AddExp(game.WinnerUserId, 2)

	// 比赛结果
	matchResult := &pb.MsgMatchResult{}
	matchResult.MatchType = proto.Int(int(game.GameType))
	matchResult.IsSNGEnd = proto.Bool(false)

	glog.V(2).Info("===>游戏结束gameType:", game.GameType, " gameId:", game.GameId)

	if winner != nil && game.TotalUserBetGold > 0 {
		if winner.User.GetIsRobot() {
			stats.GetMatchLogManager().AddMatchLog(int(game.GameType), -game.TotalUserBetGold, game.TotalBetGold-winner.BetGold)
			// ai胜，则ai赢取的钱为真实玩家总下注
			stats.GetAiFortuneLogManager().AddEarnGold(int(game.GameType), int64(game.TotalUserBetGold))
		} else {
			// 真实玩家
			stats.GetMatchLogManager().AddMatchLog(int(game.GameType), game.TotalUserBetGold-winner.BetGold, game.TotalBetGold-winner.BetGold)
			// ai负，则ai输的钱为ai的总下注（即当前总下注-玩家总下注)
			stats.GetAiFortuneLogManager().AddEarnGold(int(game.GameType), -int64(game.TotalBetGold-game.TotalUserBetGold))
		}
	}

	for _, p := range game.Players {
		matchResult.MaxCards = []int32{}

		if game.Logic.CompareCards2(game.Logic.GetCardsInt32(p.Pos), p.User.GetMatchRecord().GetMaxCards()) {
			matchResult.MaxCards = game.Logic.GetCardsInt32(p.Pos)
			p.User.GetMatchRecord().MaxCards = matchResult.MaxCards
		}
		matchResult.CardType = proto.Int(GetCardType(game.Logic.GetCards(p.Pos)))

		glog.V(2).Info("===>游戏结束幸运值:", p.User.GetNickName(), " pos:", p.Pos, " luckyValue:", p.User.GetLuckyValue())

		if p.BetGold <= 0 {
			continue
		}

		itemMsg := &pb.MsgPokerDeskGameEndRes_PlayerGameResultDef{}
		itemMsg.UserId = proto.String(p.User.GetUserId())
		if p.User.GetUserId() == game.WinnerUserId {
			itemMsg.GainGold = proto.Int(winnerEarnGold)
			itemMsg.CurGold = proto.Int64(p.User.GetGold())
			matchResult.EarnGold = proto.Int(winnerEarnGold)
		} else {
			itemMsg.GainGold = proto.Int(-p.BetGold)
			itemMsg.CurGold = proto.Int64(p.User.GetGold())
			// 失败，获取1点经验
			domainUser.GetUserFortuneManager().AddExp(p.User.GetUserId(), 1)
			matchResult.EarnGold = proto.Int(-p.BetGold)
		}
		domainUser.GetPlayerManager().SendServerMsg("", []string{p.User.GetUserId()}, int32(pb.ServerMsgId_MQ_MATCH_RESULT), matchResult)

		msg.PlayerResultList = append(msg.PlayerResultList, itemMsg)

		if !p.User.GetIsRobot() {
			log := &GameLog{}
			log.Username = p.User.GetUserId()
			log.GameId = game.GameId
			log.GameType = int(game.GameType)
			log.SingleBet = game.SingleBetGold
			log.TotalBet = game.TotalBetGold
			log.CurRound = game.CurRound
			if util.IsGameTypeSNG(game.GameType) {
				log.MatchTimes = game.MatchTimes
			} else {
				log.MatchTimes = 0
			}
			log.SeenCard = p.IsSeenCard
			log.BetGold = p.BetGold
			log.CurGold = int(p.User.GetGold())
			log.Winner = p.User.GetNickName()
			log.EarnGold = int(itemMsg.GetGainGold())
			log.Event = "游戏结束"
			SaveGameLog(log)
		}
	}

	if util.IsGameTypeSNG(game.GameType) {
		glog.V(2).Info("====>SNG单局结束MatchTimes:", game.MatchTimes, " SNGGameEND:", game.SNGGameEnd)
		if game.MatchTimes >= 8 || game.SNGGameEnd {
			glog.V(2).Info("====>SNG比赛结束matchTimes:", game.MatchTimes)
			game.MatchTimes = 0
			msg.IsSNGEnd = proto.Bool(true)

			c, _ := config.GetMatchConfigManager().GetSNGMatchConfig(int(game.GameType))

			maxGold := 0
			for _, p := range game.Players {
				p.IsReady = false
				if int(p.User.GetGold()) > maxGold {
					maxGold = int(p.User.GetGold())
				}
			}

			winnerCount := 0
			for _, p := range game.Players {
				if int(p.User.GetGold()) == maxGold {
					winnerCount++
				}
			}

			if winnerCount <= 0 {
				winnerCount = 1
			}

			msg.RewardScore = proto.Int(c.RewardScore)

			for _, p := range game.Players {
				if int(p.User.GetGold()) == maxGold {
					rewardGold := c.RewardGold / winnerCount
					curGold, ok := domainUser.GetUserFortuneManager().EarnGold(p.User.GetUserId(), int64(rewardGold), "SNG胜利")
					if !ok {
						glog.V(2).Info("==>获取金币失败userId:", p.User.GetUserId(), " earnGold:", rewardGold)
						return
					}
					p.User.Gold = proto.Int64(curGold)
					domainUser.GetUserFortuneManager().AddScore(p.User.GetUserId(), c.RewardScore)

					msg.SNGWinnerUserIds = append(msg.SNGWinnerUserIds, p.User.GetUserId())
					matchResult.EarnGold = proto.Int(c.RewardGold / winnerCount)
					matchResult.IsSNGEnd = proto.Bool(true)
					domainUser.GetPlayerManager().SendServerMsg("", []string{p.User.GetUserId()}, int32(pb.ServerMsgId_MQ_MATCH_RESULT), matchResult)
				}
			}
		}
	}

	glog.V(2).Info("====>游戏结束msg:", msg)

	for _, p := range game.Players {
		for _, item := range msg.PlayerResultList {
			item.Cards = []int32{}
			itemUser := game.Players[item.GetUserId()]
			if itemUser != nil {
				if itemUser.IsCompared || item.GetUserId() == game.WinnerUserId || item.GetUserId() == p.User.GetUserId() {
					item.Cards = game.Logic.GetCardsInt32(itemUser.Pos)
				}
			}
		}
		domainUser.GetPlayerManager().SendClientMsg(p.User.GetUserId(), int32(pb.MessageId_POKER_DESK_GAME_END), msg)
	}

	for _, p := range game.Players {
		isWin := (p.User.GetUserId() == game.WinnerUserId)
		returnValue, isUpdate := newUserTask.GetNewUserTaskManager().CheckUserPlayTask(p.User.GetUserId(), int(game.GameType), isWin, p.User.GetChannel())
		if isUpdate == true {
			for i := 0; i < 7; i++ {
				if returnValue[i] == 1 {
					msgT := &pb.MsgNewbeTaskCompletedNotify{}
					msgT.Id = proto.Int(i + 1)
					domainUser.GetPlayerManager().SendClientMsg(p.User.GetUserId(), int32(pb.MessageId_NOTIFY_NEWBETASK_COMP), msgT)
				}
			}
		}
	}

	if game.GameType == 1 {
		tempCoins := int(winner.User.GetGold()) - winnerEarnGold
		if int(winner.User.GetGold()) > 50000 && tempCoins < 50000 {
			f, ok := domainUser.GetUserFortuneManager().GetUserFortune(winner.User.GetUserId())
			if ok {
				nowDay := int(time.Now().Day())
				notifyDay := int(f.ChangeGameTypeNotifyDay)
				if nowDay != notifyDay {
					msgT := &pb.MsgChangeGameTypeNotify{}
					msgT.NotifyStr = proto.String("您携带的金币数已达50000，请去中级场进行游戏")
					msgT.GameType = proto.Int(int(game.GameType))
					domainUser.GetPlayerManager().SendClientMsg(winner.User.GetUserId(), int32(pb.MessageId_CHANGE_GAME_TYPE_NOTIFY), msgT)
					f.ChangeGameTypeNotifyDay = nowDay
					domainUser.GetUserFortuneManager().ManagerChangeGameTypeNotify(winner.User.GetUserId(), nowDay)
				}
			}
		}
	} else if game.GameType == 4 {
		tempCoins := int(winner.User.GetGold()) - winnerEarnGold
		if (int(winner.User.GetGold()) > 500000 && tempCoins < 500000) || (int(winner.User.GetGold()) > 200000 && tempCoins < 200000) {
			f, ok := domainUser.GetUserFortuneManager().GetUserFortune(winner.User.GetUserId())
			if ok {
				nowDay := int(time.Now().Day())
				notifyDay := int(f.ChangeGameTypeNotifyDay)
				if nowDay != notifyDay {
					msgT := &pb.MsgChangeGameTypeNotify{}
					msgT.NotifyStr = proto.String("您携带的金币数可以去高级场，请去高级级场进行游戏吧")
					msgT.GameType = proto.Int(int(game.GameType))
					domainUser.GetPlayerManager().SendClientMsg(winner.User.GetUserId(), int32(pb.MessageId_CHANGE_GAME_TYPE_NOTIFY), msgT)
					f.ChangeGameTypeNotifyDay = nowDay
					domainUser.GetUserFortuneManager().ManagerChangeGameTypeNotify(winner.User.GetUserId(), nowDay)
				}
			}
		}
	}

	for _, p := range game.Players {
		p.resetPlayerStatus()
	}

	if util.IsGameTypeWanRen(game.GameType) {
		// 分配奖池
		game.BetPond = make(map[string]int)
		game.BetPondItems = make(map[string]map[string]int)
	}

	if util.IsGameTypeSNG(game.GameType) && msg.GetIsSNGEnd() {
		// 踢出所有玩家
		go func() {
			time.Sleep(4 * time.Second)
			game.LockItem("onEndGame go")
			defer game.UnlockItem("onEndGame go")

			// 解锁游戏
			LockGame(game.GameId, false)

			for _, p := range game.Players {
				game.KickOutOfGame(p.User.GetUserId(), false)
			}
		}()
	}
}

func (game *GameItem) calcAllInGold() (int, bool, bool, string) {
	maxBetGold, ok := config.GetMatchConfigManager().GetMaxChip(int(game.GameType))
	if !ok {
		return 0, false, false, ""
	}
	minBetGold, ok := config.GetMatchConfigManager().GetMinChip(int(game.GameType))
	if !ok {
		return 0, false, false, ""
	}
	if minBetGold <= 0 {
		return 0, false, false, ""
	}

	betGold := (MaxRound - game.CurRound) * maxBetGold / minBetGold * minBetGold
	if betGold <= 0 {
		betGold = game.SingleBetGold / minBetGold * minBetGold
	}
	glog.Info("==calcAllInGold betGold =  ", betGold)

	seeCardBetGold := betGold * 2
	glog.Info("==calcAllInGold seeCardBetGold =  ", seeCardBetGold)

	pUserId := ""
	isAll := false

	for _, p := range game.Players {
		if util.IsGameTypeSNG(game.GameType) {
			if int(p.User.GetGold()) < betGold {
				betGold = int(p.User.GetGold())
			}
		} else {
			f, ok := domainUser.GetUserFortuneManager().GetUserFortune(p.User.GetUserId())
			if !ok {
				glog.V(2).Info("===>查询玩家金币数失败userId:", p.User.GetUserId())
				return 0, false, false, ""
			}
			if p.IsPlaying == true {
				if !p.IsSeenCard {
					if f.Gold < int64(betGold) {
						betGold = int(f.Gold) / minBetGold * minBetGold
						isAll = true
						pUserId = p.User.GetUserId()
					}
				} else {
					if f.Gold < int64(seeCardBetGold) {
						seeCardBetGold = int(f.Gold)
						betGold = seeCardBetGold / minBetGold * minBetGold / 2
						isAll = true
						pUserId = p.User.GetUserId()
					}
				}

			}
		}
	}
	glog.Info("==calcAllInGold betGold1 =  ", betGold)
	//betGold = betGold / minBetGold * minBetGold
	glog.Info("==calcAllInGold betGold2 =  ", betGold)

	return betGold, true, isAll, pUserId
}

func (game *GameItem) getPlayerCount() int {
	game.LockItem("getPlayerCount")
	defer game.UnlockItem("getPlayerCount")

	return len(game.Players)
}

func (game *GameItem) getPlayingCount() int {
	count := 0
	for _, p := range game.Players {
		if p.IsPlaying {
			count += 1
		}
	}
	return count
}

func (game *GameItem) getAllInPlayer() *GamePlayer {
	for _, p := range game.Players {
		if p.IsPlaying && p.AllIn {
			return p
		}
	}
	return nil
}

func (game *GameItem) getPlayer(userId string) *GamePlayer {
	game.RLock()
	defer game.RUnlock()

	return game.Players[userId]
}

func (game *GameItem) onUpdateCharm(userId string, charm int) {
	pTemp := game.getPlayer(userId)
	glog.Info("==onChatMsg pTemp =  ", pTemp)
	uPtr := pTemp.User
	glog.Info("==onChatMsg uPtr =  ", uPtr)
	//glog.Info("==onChatMsg uPtr =  ", pTemp.User.)
	if uPtr != nil {
		//pTemp.User.Charm = proto.Int32(int32(charm))
		tt := int(uPtr.GetCharm())
		glog.Info("==onChatMsg tt =  ", tt)
		glog.Info("==onChatMsg charm =  ", charm)
		//int32(uPtr.GetCharm())
		ttt := int32(charm) + int32(tt)
		glog.Info("==onChatMsg ttt =  ", ttt)
		pTemp.User.Charm = proto.Int32(ttt)
		glog.Info("==onChatMsg pTemp.User.Charm =  ", *pTemp.User.Charm)
	}
}

func (game *GameItem) onChatMsg(userId string, msg *pb.MsgChat) {
	p := game.getPlayer(userId)
	if p == nil {
		return
	}

	if p.User.GetIsRobot() && time.Since(game.CompareTime).Seconds() < 4 {
		return
	}

	//glog.Info("==onChatMsg gift_Id =  ", msg.GetGiftId())

	if msg.GetMessageType() == pb.ChatMessageType_GIFT || msg.GetMessageType() == pb.ChatMessageType_GIFTALL {
		gift_Id := msg.GetGiftId()
		consumeGold := 0
		switch gift_Id {
		case Flower_M_Id:
			if game.GameType == util.GameType_Common_Level_2 {
				consumeGold = 500
			} else if game.GameType == util.GameType_Common_Level_3 {
				consumeGold = 2000
			} else if game.GameType == util.GameType_Common_Level_1 {
				consumeGold = 50
			} else if game.GameType == util.GameType_Common_Level_4 {
				consumeGold = 500
			}
		case Eggs_M_Id:
			if game.GameType == util.GameType_Common_Level_2 {
				consumeGold = 500
			} else if game.GameType == util.GameType_Common_Level_3 {
				consumeGold = 2000
			} else if game.GameType == util.GameType_Common_Level_1 {
				consumeGold = 50
			} else if game.GameType == util.GameType_Common_Level_4 {
				consumeGold = 500
			}
		case Cheers_M_Id:
			if game.GameType == util.GameType_Common_Level_2 {
				consumeGold = 500
			} else if game.GameType == util.GameType_Common_Level_3 {
				consumeGold = 2000
			} else if game.GameType == util.GameType_Common_Level_1 {
				consumeGold = 50
			} else if game.GameType == util.GameType_Common_Level_4 {
				consumeGold = 500
			}
		case Shoe_M_Id:
			if game.GameType == util.GameType_Common_Level_2 {
				consumeGold = 500
			} else if game.GameType == util.GameType_Common_Level_3 {
				consumeGold = 2000
			} else if game.GameType == util.GameType_Common_Level_1 {
				consumeGold = 50
			} else if game.GameType == util.GameType_Common_Level_4 {
				consumeGold = 500
			}
		case Kiss_M_Id:
			if game.GameType == util.GameType_Common_Level_2 {
				consumeGold = 500
			} else if game.GameType == util.GameType_Common_Level_3 {
				consumeGold = 2000
			} else if game.GameType == util.GameType_Common_Level_1 {
				consumeGold = 50
			} else if game.GameType == util.GameType_Common_Level_4 {
				consumeGold = 500
			}
		case Bomb_M_Id:
			if game.GameType == util.GameType_Common_Level_2 {
				consumeGold = 500
			} else if game.GameType == util.GameType_Common_Level_3 {
				consumeGold = 2000
			} else if game.GameType == util.GameType_Common_Level_1 {
				consumeGold = 50
			} else if game.GameType == util.GameType_Common_Level_4 {
				consumeGold = 500
			}
		case Flower_F_Id:
			if game.GameType == util.GameType_Common_Level_2 {
				consumeGold = 500
			} else if game.GameType == util.GameType_Common_Level_3 {
				consumeGold = 2000
			} else if game.GameType == util.GameType_Common_Level_1 {
				consumeGold = 50
			} else if game.GameType == util.GameType_Common_Level_4 {
				consumeGold = 500
			}
		case Eggs_F_Id:
			if game.GameType == util.GameType_Common_Level_2 {
				consumeGold = 500
			} else if game.GameType == util.GameType_Common_Level_3 {
				consumeGold = 2000
			} else if game.GameType == util.GameType_Common_Level_1 {
				consumeGold = 50
			} else if game.GameType == util.GameType_Common_Level_4 {
				consumeGold = 500
			}
		case Cheers_F_Id:
			if game.GameType == util.GameType_Common_Level_2 {
				consumeGold = 500
			} else if game.GameType == util.GameType_Common_Level_3 {
				consumeGold = 2000
			} else if game.GameType == util.GameType_Common_Level_1 {
				consumeGold = 50
			} else if game.GameType == util.GameType_Common_Level_4 {
				consumeGold = 500
			}
		case Shoe_F_Id:
			if game.GameType == util.GameType_Common_Level_2 {
				consumeGold = 500
			} else if game.GameType == util.GameType_Common_Level_3 {
				consumeGold = 2000
			} else if game.GameType == util.GameType_Common_Level_1 {
				consumeGold = 50
			} else if game.GameType == util.GameType_Common_Level_4 {
				consumeGold = 500
			}
		case Kiss_F_Id:
			if game.GameType == util.GameType_Common_Level_2 {
				consumeGold = 500
			} else if game.GameType == util.GameType_Common_Level_3 {
				consumeGold = 2000
			} else if game.GameType == util.GameType_Common_Level_1 {
				consumeGold = 50
			} else if game.GameType == util.GameType_Common_Level_4 {
				consumeGold = 500
			}
		case Bomb_F_Id:
			if game.GameType == util.GameType_Common_Level_2 {
				consumeGold = 500
			} else if game.GameType == util.GameType_Common_Level_3 {
				consumeGold = 2000
			} else if game.GameType == util.GameType_Common_Level_1 {
				consumeGold = 50
			} else if game.GameType == util.GameType_Common_Level_4 {
				consumeGold = 500
			}
		}

		if msg.GetMessageType() == pb.ChatMessageType_GIFTALL {
			consumeGold = consumeGold * (game.getPlayerCount() - 1)
		}
		//glog.Info("==onChatMsg playing Count =  ", game.getPlayerCount())
		//glog.Info("==onChatMsg messageType =  ", msg.GetMessageType())
		//glog.Info("==onChatMsg playing Count2 =  ", game.Players)

		gifttype := util.FINE
		if gift_Id == Eggs_M_Id || gift_Id == Shoe_M_Id || gift_Id == Bomb_M_Id || gift_Id == Eggs_F_Id || gift_Id == Shoe_F_Id || gift_Id == Bomb_F_Id {
			gifttype = util.BAD
		}

		if consumeGold > 0 {
			_, _, ok := domainUser.GetUserFortuneManager().ConsumeGoldNoMsg(userId, int64(consumeGold), false, fmt.Sprintf("礼物%v", msg.GetGiftId()))
			if !ok {
				glog.V(2).Info("表情扣费失败userId:", userId)
				return
			}
			if gifttype == util.FINE {
				util.MongoLog_GameGiftFine(consumeGold)
			} else {
				util.MongoLog_GameGiftBad(consumeGold)
			}
		}

		if gifttype == util.FINE {
			charm := 0
			switch game.GameType {
			case util.GameType_Common_Level_2:
				charm = 1
			case util.GameType_Common_Level_3:
				charm = 4
			case util.GameType_Common_Level_4:
				charm = 1
			}

			if charm != 0 {
				if msg.GetMessageType() == pb.ChatMessageType_GIFT {
					domainUser.GetUserFortuneManager().EarnCharm(userId, charm)
					domainUser.GetUserFortuneManager().EarnCharm(msg.GetToUserId(), charm)
					game.onUpdateGoldInGame(msg.GetToUserId())
					util.MongoLog_CharmPool(charm * 2)

					toUid := msg.GetToUserId()
					game.onUpdateCharm(toUid, charm)
					game.onUpdateCharm(userId, charm)
				} else if msg.GetMessageType() == pb.ChatMessageType_GIFTALL {
					charm_all := 0
					for _, p := range game.Players {
						if p.User.GetUserId() != userId {
							domainUser.GetUserFortuneManager().EarnCharm(p.User.GetUserId(), charm)
							game.onUpdateGoldInGame(p.User.GetUserId())
							charm_all += charm
							game.onUpdateCharm(p.User.GetUserId(), charm)
						}
					}
					domainUser.GetUserFortuneManager().EarnCharm(userId, charm_all)
					game.onUpdateCharm(userId, charm_all)
					util.MongoLog_CharmPool(charm_all * 2)
				}
			}
		} else if gifttype == util.BAD {
			charm := 0
			switch game.GameType {
			case util.GameType_Common_Level_2:
				charm = -2
			case util.GameType_Common_Level_3:
				charm = -8
			case util.GameType_Common_Level_4:
				charm = -2
			}

			if charm != 0 {
				if msg.GetMessageType() == pb.ChatMessageType_GIFT {
					domainUser.GetUserFortuneManager().EarnCharm(msg.GetToUserId(), charm)
					game.onUpdateGoldInGame(msg.GetToUserId())
					util.MongoLog_CharmPool(charm)
					toUid := msg.GetToUserId()
					game.onUpdateCharm(toUid, charm)
					game.onUpdateCharm(userId, charm)
				} else if msg.GetMessageType() == pb.ChatMessageType_GIFTALL {
					charm_all := 0
					for _, p := range game.Players {
						if p.User.GetUserId() != userId {
							domainUser.GetUserFortuneManager().EarnCharm(p.User.GetUserId(), charm)
							game.onUpdateGoldInGame(p.User.GetUserId())
							charm_all += charm
							game.onUpdateCharm(p.User.GetUserId(), charm)
						}
					}
					util.MongoLog_CharmPool(charm_all)
					game.onUpdateCharm(userId, charm_all)
				}
			}
		}
		game.onUpdateGoldInGame(userId)
		//glog.Info("==onChatMsg id =  ", userId)
	}
	game.broadcast(int32(pb.MessageId_CHAT), msg)
}

func (game *GameItem) onConsumeProps(userId string, itemType pb.MagicItemType, changeCard int) {
	game.LockItem("onConsumeProps")
	defer game.UnlockItem("onConsumeProps")

	res := &pb.MsgUseMagicItemRes{}
	res.ItemType = itemType.Enum()

	glog.V(2).Info("===>使用道具userId:", userId, " type:", itemType)
	if !util.IsGameTypeProps(game.GameType) {
		glog.V(2).Info("非道具模式，禁止使用道具userId:", userId, " gameType:", game.GameType)
		res.Code = pb.MsgUseMagicItemRes_FAILED.Enum()
		domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_USE_MAGIC_ITEM), res)
		return
	}

	p := game.Players[userId]
	if p == nil {
		res.Code = pb.MsgUseMagicItemRes_FAILED.Enum()
		domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_USE_MAGIC_ITEM), res)
		return
	}

	res.Code = pb.MsgUseMagicItemRes_OK.Enum()

	if itemType == pb.MagicItemType_FOURFOLD_GOLD {
		// 翻倍卡
		p.IsUseDoubleCard = true
	} else if itemType == pb.MagicItemType_PROHIBIT_COMPARE {
		// 禁止比牌
		p.UseForbidCardRound = game.CurRound
	} else if itemType == pb.MagicItemType_REPLACE_CARD {
		// 换牌
		if !p.IsSeenCard {
			res.Code = pb.MsgUseMagicItemRes_FAILED.Enum()
			domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_USE_MAGIC_ITEM), res)
			return
		}
		if p.ChangeCardTimes >= MaxChangeCardTimes {
			glog.Error("超过最大换牌次数userId:", userId)
			res.Code = pb.MsgUseMagicItemRes_FAILED.Enum()
			domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_USE_MAGIC_ITEM), res)
			return
		}
		glog.V(2).Info("====>nickname:", p.User.GetNickName(), " 换牌前:", game.Logic.GetCards(p.Pos), " changeCard:", changeCard)
		card := game.Logic.ReplaceCard(p.Pos, changeCard)
		if card == 0 {
			glog.Error("换牌失败userId:", userId, " changeCard:", changeCard)
			res.Code = pb.MsgUseMagicItemRes_FAILED.Enum()
			domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_USE_MAGIC_ITEM), res)
			return
		}
		glog.V(2).Info("====>nickname:", p.User.GetNickName(), " 换牌后:", game.Logic.GetCards(p.Pos), " changeCard:", changeCard, " card:", card)
		res.OldCard = proto.Int(changeCard)
		res.NewCard = proto.Int(card)
		p.ChangeCardTimes++

		game.sendRobotMaxCardUser()
	}

	broMsg := &pb.MsgUseMagicItemBro{}
	broMsg.UserId = proto.String(userId)
	broMsg.ItemType = itemType.Enum()
	broMsg.Round = proto.Int(game.CurRound)
	game.broadcastExcept(int32(pb.MessageId_USE_MAGIC_ITEM_BRO), broMsg, []string{userId})

	domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_USE_MAGIC_ITEM), res)

	glog.V(2).Info("===>使用道具成功res:", res)
}

func (game *GameItem) KickOutOfGame(userId string, timeout bool) {
	glog.Info("===>踢出玩家userId:", userId)
	game.leaveGame(userId, false, true, timeout)

	go GetDeskManager().LeaveGame(userId, game.GameType, false)
}

func LockGame(gameId int, lock bool) {
	go GetDeskManager().LockGame(gameId, lock)
}

func (game *GameItem) onJoinWaitQueue(userId string) {
	if !util.IsGameTypeWanRen(game.GameType) {
		return
	}

	game.LockItem("onJoinWaitQueue")
	defer game.UnlockItem("onJoinWaitQueue")

	p := game.Lookup[userId]
	if p == nil {
		glog.Error("玩家userId:", userId, " 上桌失败，不在观看列表中!")
		return
	}

	for _, item := range game.WaitQueue {
		if item == userId {
			// 已在上桌列表
			glog.Error("玩家userId:", userId, " 已经在上桌列表!")
			return
		}
	}

	c := config.GetMatchConfigManager().GetWanRenConfig()

	f, ok := domainUser.GetUserFortuneManager().GetUserFortune(userId)
	if !ok {
		glog.V(2).Info("==>查询玩家财富失败userId:", userId)
		return
	}

	if f.Gold < int64(c.JoinWaitQueueLimit) {
		glog.V(2).Info("==>玩家userId:", userId, " 金币不足,上桌失败!")
		return
	}

	glog.V(2).Info("==>玩家userId:", userId, " 上桌成功!")

	game.WaitQueue = append(game.WaitQueue, userId)
	game.WaitQueueOrder++

	msg := &pb.MsgJoinWaitQueue{}
	msg.User = &pb.WanRenWaitQueueUserDef{}
	msg.User.User = p
	msg.User.Order = proto.Int(game.WaitQueueOrder)

	game.broadcast(int32(pb.MessageId_JOIN_WAIT_QUEUE), msg)
}

func (game *GameItem) onLeaveWaitQueue(userId string) {
	if !util.IsGameTypeWanRen(game.GameType) {
		return
	}

	game.LockItem("onLeaveWaitQueue")
	defer game.UnlockItem("onLeaveWaitQueue")

	glog.V(2).Info("==>玩家userId:", userId, " 申请下桌")

	for i, v := range game.WaitQueue {
		glog.V(2).Info("===>v:", v, " userId:", userId)
		if v == userId {
			game.WaitQueue = append(game.WaitQueue[:i], game.WaitQueue[i+1:]...)
			if len(game.WaitQueue) <= 0 {
				game.WaitQueueOrder = 0
			}

			msg := &pb.MsgLeaveWaitQueue{}
			msg.UserId = proto.String(v)

			game.broadcast(int32(pb.MessageId_LEAVE_WAIT_QUEUE), msg)

			glog.V(2).Info("==>玩家userId:", userId, " 下桌成功!")
			return
		}
	}
}

func (game *GameItem) onLookupBetGold(userId string, betUserId string) {
	if !util.IsGameTypeWanRen(game.GameType) {
		return
	}

	game.LockItem("onLookupBetGold")
	defer game.UnlockItem("onLookupBetGold")

	glog.V(2).Info("====>万人场旁观下注userId:", userId, " betUserId:", betUserId)
	if game.IsStart {
		glog.Error("旁观下注失败，游戏已开始!")
		return
	}

	if !game.WanRenWaitingBet {
		glog.V(2).Info("====>万人场押注阶段结束startTime:", game.WanRenWaitingBetTime)
		return
	}

	lookup := game.Lookup[userId]
	if lookup == nil {
		glog.Error("不存在此旁观者userId:", userId)
		return
	}

	p := game.Players[betUserId]
	if p == nil {
		glog.Error("userId:", userId, " 下注玩家不在桌上betUserId:", betUserId)
		return
	}

	_, _, ok := domainUser.GetUserFortuneManager().ConsumeGold(userId, int64(WanRenPondMinBet), false, "万人场旁观下注")
	if !ok {
		glog.V(2).Info("万人场旁观下注失败!")
		return
	}

	game.BetPond[betUserId] += WanRenPondMinBet
	game.TotalBetPond += WanRenPondMinBet

	if _, ok := game.BetPondItems[betUserId]; !ok {
		game.BetPondItems[betUserId] = make(map[string]int)
	}

	game.BetPondItems[betUserId][userId] += WanRenPondMinBet

	msg := &pb.MsgLookupBetGoldRes{}
	msg.Ok = proto.Bool(true)
	msg.BetUserId = proto.String(betUserId)
	msg.BetGold = proto.Int(game.BetPondItems[betUserId][userId])

	domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_LOOKUP_USER_BET), msg)

	broMsg := &pb.MsgLookupBetGoldBro{}
	broMsg.BetUserId = proto.String(betUserId)
	broMsg.BetGold = proto.Int(game.BetPond[betUserId])
	broMsg.TotalBetPond = proto.Int(game.TotalBetPond)

	game.broadcast(int32(pb.MessageId_LOOKUP_USER_BET_BRO), broMsg)
}

func (game *GameItem) onUpdateGoldInGame(userId string) {
	game.LockItem("onUpdateGoldInGame")
	defer game.UnlockItem("onUpdateGoldInGame")

	if util.IsGameTypeSNG(game.GameType) {
		return
	}

	//glog.Info("===>更新游戏内金币userId:", userId)

	p := game.Players[userId]
	if p == nil {
		glog.V(2).Info("===>游戏内无此玩家userId:", userId)
		if util.IsGameTypeWanRen(game.GameType) {
			// 万人场
			lookupUser := game.Lookup[userId]
			if lookupUser == nil {
				return
			}
			f, ok := domainUser.GetUserFortuneManager().GetUserFortune(lookupUser.GetUserId())
			if !ok {
				return
			}
			lookupUser.Gold = proto.Int64(f.Gold)
		}
		return
	}

	f, ok := domainUser.GetUserFortuneManager().GetUserFortune(p.User.GetUserId())
	if !ok {
		return
	}
	p.User.Gold = proto.Int64(f.Gold)

	msg := &pb.MsgUpdateGoldInGame{}
	msg.UserId = proto.String(userId)
	msg.Gold = proto.Int64(f.Gold)
	msg.Charm = proto.Int32(int32(f.Charm))
	//glog.Info("===>更新游戏内金币userId:", userId, " gold :", f.Gold, " charm : ", f.Charm)

	game.broadcast(int32(pb.MessageId_UPDATE_GOLD_IN_GAME), msg)
}

func (game *GameItem) onRewardInGame(userId string) bool {
	game.LockItem("onRewardInGame")
	defer game.UnlockItem("onRewardInGame")

	p := game.Players[userId]
	if p == nil {
		glog.V(2).Info("打赏失败，找不到对应玩家userId:", userId)
		return false
	}

	minChip := config.GetMatchConfigManager().GetTipChip(int(game.GameType))

	//glog.Info("===>打赏userId:", userId, " minChip:", minChip)

	curGold, _, ok := domainUser.GetUserFortuneManager().ConsumeGold(p.User.GetUserId(), int64(minChip), false, "打赏")
	if !ok {
		glog.V(2).Info("==>打赏扣款失败userId:", userId, " minChip:", minChip)
		return false
	}

	util.MongoLog_SystemTip(minChip)

	if !util.IsGameTypeSNG(game.GameType) {
		p.User.Gold = proto.Int64(curGold)
	}

	msg := &pb.MsgRewardInGame{}
	msg.UserId = proto.String(userId)

	game.broadcast(int32(pb.MessageId_REWARD_IN_GAME), msg)

	f, ok := domainUser.GetUserFortuneManager().GetUserFortune(p.User.GetUserId())
	if !ok {
		return false
	}

	updateGoldMsg := &pb.MsgUpdateGoldInGame{}
	updateGoldMsg.UserId = proto.String(userId)
	updateGoldMsg.Gold = proto.Int64(curGold)
	updateGoldMsg.Charm = proto.Int32(int32(f.Charm))
	game.broadcast(int32(pb.MessageId_UPDATE_GOLD_IN_GAME), updateGoldMsg)

	return true
}

func (game *GameItem) sendRobotMaxCardUser() {
	maxUserId := game.getMaxCardUserId()
	if maxUserId == "" {
		return
	}

	if game.MaxCardUserId == maxUserId {
		return
	}

	game.MaxCardUserId = maxUserId

	dstIds := []string{}
	for _, item := range game.Players {
		if item.User.GetIsRobot() {
			dstIds = append(dstIds, item.User.GetUserId())
		}
	}

	if len(dstIds) <= 0 {
		return
	}

	msg := &pb.MsgRobotMaxCardUser{}
	msg.UserId = proto.String(maxUserId)

	glog.V(2).Info("==>广播机器人最大牌 dstIds:", dstIds, " msg:", msg)
	domainUser.GetPlayerManager().SendClientMsg2(dstIds, int32(pb.MessageId_ROBOT_MAX_CARD_USER), msg)
}

func (game *GameItem) getMaxCardUserId() string {
	userId := ""
	pos := 0
	for _, p := range game.Players {
		if userId == "" {
			userId = p.User.GetUserId()
			pos = p.Pos
			continue
		}
		if !game.Logic.CompareCardsByPos(pos, p.Pos) {
			userId = p.User.GetUserId()
			pos = p.Pos
		}
	}

	return userId
}

func (game *GameItem) sendRobotCards() {
	for _, p := range game.Players {
		if p.User.GetIsRobot() {
			msg := &pb.MsgRobotCards{}
			msg.Cards = game.Logic.GetCards(p.Pos)
			msg.CardType = proto.Int(GetCardType(msg.GetCards()))
			domainUser.GetPlayerManager().SendClientMsg(p.User.GetUserId(), int32(pb.MessageId_ROBOT_CARDS), msg)
		}
	}
}

func (game *GameItem) onEnterBackground(userId string) {
	// 万人场从上桌列表中移除
	game.onLeaveWaitQueue(userId)

	glog.V(2).Info("===>onEnterBackground userId:", userId)

	game.LockItem("onEnterBackground")
	defer game.UnlockItem("onEnterBackground")

	p := game.Players[userId]
	if p == nil {
		if util.IsGameTypeWanRen(game.GameType) {
			go GetDeskManager().LeaveGame(userId, game.GameType, false)
		}
		return
	}

	p.IsEnterBackground = true
}

func (game *GameItem) onEnterForeground(userId string) {
	// 向其发送牌桌信息
	game.LockItem("onEnterForeground")
	defer game.UnlockItem("onEnterForeground")

	p := game.Players[userId]
	if p == nil {
		domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_APP_ENTER_FOREGROUND), nil)
		return
	}

	p.IsEnterBackground = false
	glog.V(2).Info("===>onEnterForeground userId:", userId)

	msg := &pb.MsgGetPokerDeskInfoRes{}
	msg.DeskInfo = game.BuildMessage()
	msg.Type = util.ToMatchType(game.GameType)
	domainUser.GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_APP_ENTER_FOREGROUND), msg)
}
