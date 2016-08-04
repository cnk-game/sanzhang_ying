package main

import (
	"code.google.com/p/go.net/websocket"
	"code.google.com/p/goprotobuf/proto"
	"config"
	"github.com/golang/glog"
	"math/rand"
	"pb"
	"sync"
	"time"
	"util"
)

var mu sync.RWMutex

func waitEnterGame() {
	return
	mu.Lock()
	defer mu.Unlock()

	time.Sleep(time.Duration(100+rand.Int()%250) * time.Millisecond)
}

type GamePlayer struct {
	UserId             string
	IsPlaying          bool
	ForbidCompareRound int
}

type Robot struct {
	sync.RWMutex
	Username            string
	Nickname            string
	Gender              int
	Sign                string
	Photo               string
	Gold                int
	Vip                 int
	conn                *websocket.Conn
	dur                 time.Duration
	userInfo            *pb.UserDef
	gameType            int
	GameStart           bool
	MaxCardUserId       string
	CurSingleBet        int
	CurRound            int
	CurRoundUserId      string
	AllIn               bool
	Cards               []byte
	CardType            int
	SeenCard            bool
	Players             map[string]*GamePlayer
	PlayTimes           int
	WinTimes            int
	LoseTimes           int
	CurDayEarnGold      int
	CurWeekEarnGold     int
	MaxCards            []int
	matchType           *pb.MatchType
	bossRobot           string
	useForbidCard       bool      // 是否使用了禁比卡
	useFourTimesCard    bool      // 是否使用了4倍卡
	isRewardInGame      bool      // 本局是否已打赏
	compareTime         time.Time // 比牌开始时间(3秒内不发送动画）
	lastExpressionRound int       // 上次判定发送表情回合数
	lastSendGiftTime    time.Time
	robotConfig 		*RobotConfig
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewRobot(robotConfig *RobotConfig, matchType *pb.MatchType) *Robot {
	robot := &Robot{}
	robot.robotConfig = robotConfig
	robot.Username = robotConfig.Username.Hex()
	robot.Nickname = robotConfig.Nickname
	robot.Gender = robotConfig.Gender
	robot.Sign = robotConfig.Sign
	robot.Photo = robotConfig.Photo
	robot.Gold = robotConfig.Gold / 10000 * 10000
	robot.Vip = robotConfig.Vip
	robot.WinTimes = robotConfig.WinTimes
	robot.LoseTimes = robotConfig.LoseTimes
	robot.CurDayEarnGold = robotConfig.CurDayEarnGold
	robot.CurWeekEarnGold = robotConfig.CurWeekEarnGold
	robot.MaxCards = robotConfig.MaxCards
	robot.matchType = matchType

	robot.Players = make(map[string]*GamePlayer)

	return robot
}

func (robot *Robot) reset() {
	robot.gameType = 0
	robot.GameStart = false
	robot.MaxCardUserId = ""
	robot.CurRound = 0
	robot.CurRoundUserId = ""
	robot.AllIn = false
	robot.Cards = []byte{}
	robot.CardType = 0
	robot.SeenCard = false
	robot.Players = make(map[string]*GamePlayer)
	robot.bossRobot = ""
	robot.useForbidCard = false
	robot.useFourTimesCard = false
	robot.isRewardInGame = false
	robot.lastExpressionRound = 0
}

func (robot *Robot) sendMsg(msgId pb.MessageId, body proto.Message, waitTime time.Duration) {
	if waitTime.Seconds() < 1 {
		robot.sendMsg2(msgId, body)
	} else {
		go func(waitTime time.Duration) {
			time.Sleep(waitTime)
			robot.sendMsg2(msgId, body)
		}(waitTime)
	}
}

func (robot *Robot) sendMsg2(msgId pb.MessageId, body proto.Message) {
	if robot.conn == nil {
		glog.V(2).Info("===>robot.conn nil")
		return
	}

	msg := &pb.ClientMsg{}
	msg.MsgId = proto.Int(int(msgId))

	if body != nil {
		b, err := proto.Marshal(body)
		if err != nil {
			glog.V(2).Info(err)
			return
		}

		msg.MsgBody = b
	}

	b, err := proto.Marshal(msg)
	if err != nil {
		glog.V(2).Info(err)
	}

	robot.conn.Write(b)

	if msg.GetMsgId() == int32(pb.MessageId_PLAYER_OPERATE_CARDS) {
		m := &pb.MsgOpCardReq{}
		proto.Unmarshal(msg.MsgBody, m)
		glog.Info("==>发送消息userId:", robot.userInfo.GetUserId(), " nickname:", robot.userInfo.GetNickName(), " round:", robot.CurRound, " msg:", m)
	}
}

func (robot *Robot) Login(url, origin string) {
	conn, err := websocket.Dial(url, "", origin)
	if err != nil {
		glog.V(2).Info("连接服务器失败username:", robot.Username)
		return
	}
	robot.conn = conn
	defer conn.Close()
	defer func() {
		glog.Info("连接断开conn:", robot.conn, "重连")
		r := NewRobot(robot.robotConfig, robot.matchType)
		go r.Login(url, origin)
	}()

	glog.V(2).Info("==>连接服务器成功conn:", robot.conn)
	defer glog.V(2).Info("===>连接断开conn:", robot.conn)

	robot.sendLoginMsg()

	for {
		var data []byte
		err := websocket.Message.Receive(robot.conn, &data)
		if err != nil {
			glog.V(2).Info("error receiving msg:", err)
			break
		}

		clientMsg := &pb.ClientMsg{}
		err = proto.Unmarshal(data, clientMsg)
		if err != nil {
			glog.V(2).Info("unmarshal client msg failed!")
			break
		}

		robot.handleMsg(clientMsg)
	}
}

func (robot *Robot) handleMsg(msg *pb.ClientMsg) {
	//	glog.V(2).Info("==>收到消息msgId:", util.GetMsgIdName(msg.GetMsgId()))
	switch pb.MessageId(msg.GetMsgId()) {
	case pb.MessageId_LOGIN:
		robot.onLoginMsg(msg.GetMsgBody())
		break
	case pb.MessageId_GET_USER_INFO:
		robot.onGetUserInfo(msg.GetMsgBody())
		break
	case pb.MessageId_ENTER_POKER_DESK:
		robot.onEnterGame(msg.GetMsgBody())
		break
	case pb.MessageId_ENTER_POKER_DESK_BRO:
		robot.onEnterGameBro(msg.GetMsgBody())
		break
	case pb.MessageId_LEAVE_POKER_DESK:
		robot.onLeaveGame(msg.GetMsgBody())
		break
	case pb.MessageId_GET_POKER_DESK_INFO:
		robot.onGetPokerDeskInfo(msg.GetMsgBody())
		break
	case pb.MessageId_POKER_DESK_GAME_BEGIN:
		// 游戏开始
		robot.onGameBegin(msg.GetMsgBody())
		break
	case pb.MessageId_POKER_DESK_GAME_END:
		// 游戏结束
		robot.onGameEnd(msg.GetMsgBody())
		break
	case pb.MessageId_PLAYER_ROUND_BEGIN:
		// 游戏回合开始
		robot.onRoundBegin(msg.GetMsgBody())
		break
	case pb.MessageId_PLAYER_OPERATE_CARDS:
		robot.onOpCards(msg.GetMsgBody())
		// 游戏操作
		break
	case pb.MessageId_UPDATE_GOLD:
		robot.onUpdateGold(msg.GetMsgBody())
		break
	case pb.MessageId_UPDATE_GOLD_IN_GAME:
		robot.onUpdateGoldInGame(msg.GetMsgBody())
		break
	case pb.MessageId_ROBOT_MAX_CARD_USER:
		robot.onRobotMaxCardUser(msg.GetMsgBody())
		break
	case pb.MessageId_ROBOT_CARDS:
		robot.onRobotCards(msg.GetMsgBody())
		break
	case pb.MessageId_CHAT:
		robot.onChat(msg.GetMsgBody())
		break
	case pb.MessageId_REWARD_IN_GAME:
		robot.onRewardInGame(msg.GetMsgBody())
		break
	}
}

func (robot *Robot) getServerTime() time.Time {
	now := time.Now()
	now.Add(-robot.dur)

	return now
}

// 登录
func (robot *Robot) sendLoginMsg() {
	msg := &pb.MsgLoginReq{}
	msg.Username = proto.String(robot.Username)
	msg.Nickname = proto.String(robot.Nickname)
	msg.RobotKey = proto.String(config.RobotKey)
	msg.RobotGender = proto.Int(robot.Gender)
	msg.RobotSign = proto.String(robot.Sign)
	msg.RobotPhoto = proto.String(robot.Photo)
	msg.RobotGold = proto.Int(robot.Gold)
	msg.RobotVip = proto.Int(robot.Vip)
	msg.RobotWinTimes = proto.Int(robot.WinTimes)
	msg.RobotLoseTimes = proto.Int(robot.LoseTimes)
	msg.RobotCurDayEarnGold = proto.Int(robot.CurDayEarnGold)
	msg.RobotCurWeekEarnGold = proto.Int(robot.CurWeekEarnGold)
	for _, card := range robot.MaxCards {
		msg.RobotMaxCards = append(msg.RobotMaxCards, int32(card))
	}

	robot.sendMsg(pb.MessageId_LOGIN, msg, time.Millisecond)
}

func (robot *Robot) onLoginMsg(m []byte) {
	msg := &pb.MsgLoginRes{}
	err := proto.Unmarshal(m, msg)
	if err != nil {
		glog.V(2).Info(err)
		return
	}

	if msg.GetCode() != pb.MsgLoginRes_OK {
		glog.V(2).Info("登录失败")
		robot.conn.Close()
		return
	}

	serverTime := time.Unix(msg.GetServerTime(), 0)
	robot.dur = time.Since(serverTime)

	// 登录成功，获取玩家信息
	robot.sendGetUserInfo()
}

func (robot *Robot) sendGetUserInfo() {
	robot.sendMsg(pb.MessageId_GET_USER_INFO, nil, time.Millisecond)
}

func (robot *Robot) onGetUserInfo(m []byte) {
	msg := &pb.MsgGetUserInfoRes{}
	err := proto.Unmarshal(m, msg)
	if err != nil {
		glog.V(2).Info(err)
		robot.conn.Close()
		return
	}

	robot.userInfo = msg.GetUserInfo()

	robot.checkRobotGold(robot.Gold)

	glog.V(2).Info("===>当前玩家userId:", msg.GetUserInfo().GetUserId())

	// 进入游戏
	robot.sendEnterGame(time.Duration(5+rand.Int()%5) * time.Second)
}

func (robot *Robot) checkRobotGold(gold int) bool {
	if robot.SeenCard {
		gold *= 2
	}

	if util.IsGameTypeSNG(util.GameType(robot.gameType)) {
		return false
	}

	if robot.userInfo.GetGold() < int64(gold) {
		gold *= (5 + rand.Int()%5)
		robot.sendRobotSetGold(gold)
	}
	return true
}

// 进入游戏
func (robot *Robot) sendEnterGame(waitTime time.Duration) {
	waitEnterGame()

	robot.reset()

	msg := &pb.MsgEnterPokerDeskReq{}
	msg.Type = robot.matchType
	msg.ChangeDesk = proto.Bool(true)

	enterLimit := 0
	randGold := 0
	if msg.GetType() == pb.MatchType_COMMON_LEVEL1 {
		enterLimit, _ = config.GetMatchConfigManager().GetEnterLimit(int(util.GameType_Common_Level_1))
		randGold = rand.Int() % 40000
	} else if msg.GetType() == pb.MatchType_COMMON_LEVEL2 {
		enterLimit, _ = config.GetMatchConfigManager().GetEnterLimit(int(util.GameType_Common_Level_2))
		randGold = rand.Int() % 40000
	} else if msg.GetType() == pb.MatchType_COMMON_LEVEL3 {
		enterLimit, _ = config.GetMatchConfigManager().GetEnterLimit(int(util.GameType_Common_Level_3))
		randGold = rand.Int() % 5000000
	} else if msg.GetType() == pb.MatchType_MAGIC_ITEM_LEVEL1 {
		enterLimit, _ = config.GetMatchConfigManager().GetEnterLimit(int(util.GameType_Props_Level_1))
		randGold = rand.Int() % 40000
	} else if msg.GetType() == pb.MatchType_MAGIC_ITEM_LEVEL2 {
		enterLimit, _ = config.GetMatchConfigManager().GetEnterLimit(int(util.GameType_Props_Level_2))
		randGold = rand.Int() % 40000
	} else if msg.GetType() == pb.MatchType_MAGIC_ITEM_LEVEL3 {
		enterLimit, _ = config.GetMatchConfigManager().GetEnterLimit(int(util.GameType_Props_Level_3))
		randGold = rand.Int() % 5000000
	} else if msg.GetType() == pb.MatchType_SNG_LEVEL1 {
		enterLimit, _ = config.GetMatchConfigManager().GetEnterLimit(int(util.GameType_SNG_Level_1))
		randGold = rand.Int() % 40000
	} else if msg.GetType() == pb.MatchType_SNG_LEVEL2 {
		enterLimit, _ = config.GetMatchConfigManager().GetEnterLimit(int(util.GameType_SNG_Level_2))
		randGold = rand.Int() % 40000
	} else if msg.GetType() == pb.MatchType_SNG_LEVEL3 {
		enterLimit, _ = config.GetMatchConfigManager().GetEnterLimit(int(util.GameType_SNG_Level_3))
		randGold = rand.Int() % 5000000
	} else if msg.GetType() == pb.MatchType_WAN_REN_GAME {
		enterLimit = 5000000
		randGold = rand.Int() % 5000000
	}

	robot.sendRobotSetGold(enterLimit - int(robot.userInfo.GetGold()) + randGold)

	robot.sendMsg(pb.MessageId_ENTER_POKER_DESK, msg, waitTime)
	glog.V(2).Info("====>进入桌子userId:", robot.userInfo.GetUserId(), " 类型:", msg.GetType())
}

func (robot *Robot) onEnterGame(m []byte) {
	msg := &pb.MsgEnterPokerDeskRes{}
	err := proto.Unmarshal(m, msg)
	if err != nil {
		glog.V(2).Info(err)
		robot.conn.Close()
		return
	}

	if msg.GetCode() == pb.MsgEnterPokerDeskRes_OK {
		// 进入游戏成功
		glog.V(2).Info("进入桌子成功")
		robot.gameType = int(msg.GetType())
	} else if msg.GetCode() == pb.MsgEnterPokerDeskRes_FAILED {
		// 进入游戏失败
		glog.V(2).Info("进入桌子失败")
	} else if msg.GetCode() == pb.MsgEnterPokerDeskRes_LACK_GOLD {
		// 进入游戏失败，缺少金币
		glog.V(2).Info("进入桌子失败，缺少金币")
	}

	if msg.GetCode() != pb.MsgEnterPokerDeskRes_OK {
		robot.sendEnterGame(time.Duration(1+rand.Int()%2) * time.Minute)
	}
}

// 牌桌信息
func (robot *Robot) onGetPokerDeskInfo(m []byte) {
	msg := &pb.MsgGetPokerDeskInfoRes{}
	err := proto.Unmarshal(m, msg)
	if err != nil {
		glog.Error(err)
		return
	}

	for _, item := range msg.GetDeskInfo().GetUserList() {
		p := &GamePlayer{}
		p.UserId = item.GetBaseInfo().GetUserId()
		p.IsPlaying = item.GetIsPlaying()
		p.ForbidCompareRound = int(item.GetForbidCompareCardRound())

		robot.Players[p.UserId] = p

		glog.V(2).Info("===>桌子玩家userId:", p.UserId)
	}

	if msg.GetType() == pb.MatchType_WAN_REN_GAME {
		// 万人场游戏
		if len(msg.GetDeskInfo().GetWaitQueue())+len(msg.GetDeskInfo().GetUserList()) < 8 {
			robot.sendJoinWaitQueue()
		}
	}
}

// 进入游戏广播
func (robot *Robot) onEnterGameBro(m []byte) {
	msg := &pb.MsgEnterPokerDeskBro{}
	err := proto.Unmarshal(m, msg)
	if err != nil {
		glog.Error(err)
		return
	}

	p := &GamePlayer{}
	p.UserId = msg.GetUser().GetBaseInfo().GetUserId()
	p.IsPlaying = msg.GetUser().GetIsPlaying()
	p.ForbidCompareRound = int(msg.GetUser().GetForbidCompareCardRound())

	robot.Players[p.UserId] = p

	glog.V(2).Info("===>桌子玩家22222userId:", p.UserId)
}

// 离开游戏
func (robot *Robot) sendLeaveGame(waitTime time.Duration) {
	msg := &pb.MsgLeavePokerDeskReq{}
	msg.Type = robot.matchType
	robot.sendMsg(pb.MessageId_LEAVE_POKER_DESK, msg, waitTime)
}

func (robot *Robot) onLeaveGame(m []byte) {
	msg := &pb.MsgLeavePokerDesk{}
	err := proto.Unmarshal(m, msg)
	if err != nil {
		glog.Error(err)
		robot.conn.Close()
		return
	}

	glog.V(2).Info("===>离开游戏msg:", msg)

	p := robot.Players[msg.GetUserId()]
	if p != nil {
		p.IsPlaying = false
	}

	delete(robot.Players, msg.GetUserId())

	if msg.GetUserId() == robot.userInfo.GetUserId() {
		robot.GameStart = false
		robot.sendEnterGame(time.Duration(1+rand.Int()%2) * time.Minute)
	}
}

func (robot *Robot) onGameBegin(m []byte) {
	robot.GameStart = true

	msg := &pb.MsgGameStart{}
	err := proto.Unmarshal(m, msg)
	if err != nil {
		glog.Error(err)
		return
	}

	robot.MaxCardUserId = ""
	robot.CurSingleBet = int(msg.GetSingleBetGold())
	robot.CurRound = 0
	robot.CurRoundUserId = ""
	robot.AllIn = false
	robot.SeenCard = false
	robot.PlayTimes++
	robot.bossRobot = msg.GetBossRobotUserId()
	robot.useFourTimesCard = false
	robot.useForbidCard = false
	robot.isRewardInGame = false
	robot.lastExpressionRound = 0

	for _, item := range robot.Players {
		item.IsPlaying = true
	}
}

func (robot *Robot) onGameEnd(m []byte) {
	robot.GameStart = false

	for _, item := range robot.Players {
		item.IsPlaying = false
	}

	robot.robotLeaveGame()
}

func (robot *Robot) robotLeaveGame() {
	if util.IsGameTypeSNG(util.GameType(robot.gameType)) {
		return
	}

	if rand.Int()%100 < 50 {
		return
	}

	if robot.PlayTimes >= 2+rand.Int()%4 {
		// 退出
		robot.sendLeaveGame(time.Duration(4+rand.Int()%2) * time.Second)
		robot.PlayTimes = 0
	}
}

func (robot *Robot) onRoundBegin(m []byte) {
	msg := &pb.MsgPlayerRoundBeginRes{}
	err := proto.Unmarshal(m, msg)
	if err != nil {
		glog.Error(err)
		robot.conn.Close()
		return
	}

	robot.CurRound = int(msg.GetCurRound())
	robot.CurRoundUserId = msg.GetUserId()

	if robot.bossRobot == robot.userInfo.GetUserId() {
		// boss机器人，10回合前30%机率看牌,10回合以后80%几率看牌
		if msg.GetCurRound() <= 10 {
			if rand.Int()%100 < 50 {
				robot.sendSeeCards()
			}
		} else {
			if rand.Int()%100 < 80 {
				robot.sendSeeCards()
			}
		}
	} else {
		// 常规ai每回合有50机率看牌
		if rand.Int()%100 < 30 {
			robot.sendSeeCards()
		}
	}

	if msg.GetCurRound() >= 3 && rand.Int()%100 < 50 {
		robot.sendUseMagicItems()
	}

	if rand.Int()%100 < 10 {
		robot.rewardInGame()
	}
	robot.isRewardInGame = true

	if msg.GetCurRound() > 2 && rand.Int()%100 < 10 {
		giftUserIds := []string{}
		for _, item := range robot.Players {
			if item.UserId == robot.userInfo.GetUserId() {
				continue
			}
			if item.IsPlaying {
				giftUserIds = append(giftUserIds, item.UserId)
			}
		}
		if len(giftUserIds) > 0 {
			robot.sendGift(giftUserIds[rand.Int()%len(giftUserIds)], 0)
		}
	}

	if robot.userInfo.GetUserId() != msg.GetUserId() {
		if msg.GetCurRound() >= 3 && rand.Int()%100 < 30 {
			robot.sendExpression()
		}
		return
	}

	var waitTime time.Duration
	if robot.AllIn {
		waitTime = time.Duration(3+rand.Int()%2) * time.Second
	} else {
		waitTime = time.Duration(1+rand.Int()%4) * time.Second
	}

	glog.V(2).Info("===>切换回合round:", msg.GetCurRound())

	if msg.GetCurRound() >= 20 {
		// 超过最大回合,比牌
		robot.sendCompareCard(waitTime)
		return
	}

	if robot.bossRobot == robot.userInfo.GetUserId() {
		// 自己是boss
		if robot.AllIn {
			robot.sendAllIn(waitTime)
			return
		}

		if msg.GetCurRound() >= 3 {
			if rand.Int()%100 < 30 {
				robot.sendCompareCard(waitTime)
				return
			}
		}

		if msg.GetCurRound() >= 8 {
			if rand.Int()%100 < 30 {
				// 全下
				if robot.sendAllIn(waitTime) {
					return
				}
			}
		}
		maxChip, _ := config.GetMatchConfigManager().GetMaxChip(robot.gameType)
		if robot.CurSingleBet >= maxChip {
			// 跟注
			robot.sendFollow(waitTime)
		} else {
			if rand.Int()%100 < 70 {
				// 跟注
				robot.sendFollow(waitTime)
			} else {
				// 加注到最高
				robot.sendRaise(waitTime)
			}
		}
		return
	}

	// 陪玩机器人
	cardType := getCardType(robot.Cards)
	if robot.AllIn {
		if cardType >= config.CARD_TYPE_JIN_HUA {
			robot.sendAllIn(waitTime)
		} else {
			robot.sendGiveUp(waitTime)
		}
		return
	}

	if cardType == CARD_TYPE_SINGLE {
		if rand.Int()%100 < 80 {
			robot.sendGiveUp(waitTime)
		} else {
			robot.sendFollow(waitTime)
		}
		if msg.GetCurRound() >= 3 {
			if rand.Int()%100 < 50 {
				robot.sendCompareCard(waitTime)
			} else {
				robot.sendGiveUp(waitTime)
			}
		}
	} else if cardType == CARD_TYPE_DOUBLE {
		if rand.Int()%100 < 50 {
			robot.sendGiveUp(waitTime)
		} else {
			robot.sendFollow(waitTime)
		}

		if msg.GetCurRound() >= 3 {
			if rand.Int()%100 < 50 {
				robot.sendCompareCard(waitTime)
			} else {
				robot.sendGiveUp(waitTime)
			}
		}
	} else {
		if msg.GetCurRound() >= 3 {
			if rand.Int()%100 < 50 {
				robot.sendCompareCard(waitTime)
				return
			}
		}

		if msg.GetCurRound() >= 5 {
			if rand.Int()%100 < 80 {
				robot.sendCompareCard(waitTime)
				return
			}
		}

		if msg.GetCurRound() >= 7 {
			robot.sendCompareCard(waitTime)
			return
		}

		if msg.GetCurRound() > 2 && robot.SeenCard && rand.Int()%100 < 50 {
			maxChip, _ := config.GetMatchConfigManager().GetMaxChip(robot.gameType)
			if robot.CurSingleBet < maxChip {
				robot.sendRaise(waitTime)
				return
			}
		}
		robot.sendFollow(waitTime)
	}
}

// 弃牌
func (robot *Robot) sendGiveUp(waitTime time.Duration) {
	if !robot.SeenCard {
		msg := &pb.MsgOpCardReq{}
		msg.Type = pb.CardOpType_SEE_CARDS.Enum()
		robot.sendMsg(pb.MessageId_PLAYER_OPERATE_CARDS, msg, time.Millisecond)
	}

	glog.V(2).Info("===>发送弃牌")
	msg := &pb.MsgOpCardReq{}
	msg.Type = pb.CardOpType_GIVE_UP.Enum()
	robot.sendMsg(pb.MessageId_PLAYER_OPERATE_CARDS, msg, waitTime)

	robot.robotLeaveGame()
}

// 看牌
func (robot *Robot) sendSeeCards() {
	if robot.SeenCard {
		return
	}
	glog.V(2).Info("===>发送看牌")
	msg := &pb.MsgOpCardReq{}
	msg.Type = pb.CardOpType_SEE_CARDS.Enum()
	robot.sendMsg(pb.MessageId_PLAYER_OPERATE_CARDS, msg, time.Duration(2+rand.Int()%3)*time.Second)
}

// 比牌
func (robot *Robot) sendCompareCard(waitTime time.Duration) {
	glog.V(2).Info("===>发送比牌")

	waitTime = time.Duration(4+rand.Int()%2) * time.Second

	userIds := []string{}
	for _, item := range robot.Players {
		glog.V(2).Info("===>玩家userId:", item.UserId)
		if item.UserId == robot.userInfo.GetUserId() {
			continue
		}
		if !item.IsPlaying {
			continue
		}
		if util.IsGameTypeProps(util.GameType(robot.gameType)) {
			if robot.CurRound < 20 && item.ForbidCompareRound > 0 && robot.CurRound-item.ForbidCompareRound < 5 {
				continue
			}
		}

		userIds = append(userIds, item.UserId)
	}

	glog.V(2).Info("===>比牌userIds:", userIds)

	if len(userIds) < 1 {
		robot.checkRobotGold(robot.CurSingleBet)

		glog.V(2).Info("==>发送跟注")
		msg := &pb.MsgOpCardReq{}
		msg.Type = pb.CardOpType_FOLLOW.Enum()
		robot.sendMsg(pb.MessageId_PLAYER_OPERATE_CARDS, msg, waitTime)
		return
	}

	msg := &pb.MsgOpCardReq{}
	msg.Type = pb.CardOpType_COMPARE.Enum()
	msg.CompareUserId = proto.String(userIds[rand.Int()%len(userIds)])
	robot.sendMsg(pb.MessageId_PLAYER_OPERATE_CARDS, msg, waitTime)

	glog.V(2).Info("===>比牌msg:", msg)
}

func (robot *Robot) getPlayingCount() int {
	count := 0
	for _, item := range robot.Players {
		if item.IsPlaying {
			count++
		}
	}
	return count
}

func (robot *Robot) isCanAllIn() bool {
	if robot.AllIn {
		return true
	}

	if robot.CurRound >= 20 {
		return true
	}

	// 判断场上对手人数为1时
	count := 0
	for _, item := range robot.Players {
		if item.UserId == robot.userInfo.GetUserId() {
			continue
		}

		if !item.IsPlaying {
			continue
		}

		count++
	}
	return count == 1
}

// 全下
func (robot *Robot) sendAllIn(waitTime time.Duration) bool {
	if !robot.isCanAllIn() {
		return false
	}

	glog.V(2).Info("==>发送全押")
	msg := &pb.MsgOpCardReq{}
	msg.Type = pb.CardOpType_ALL_IN.Enum()
	robot.sendMsg(pb.MessageId_PLAYER_OPERATE_CARDS, msg, waitTime)

	return true
}

// 加注
func (robot *Robot) sendRaise(waitTime time.Duration) {
	maxChip, _ := config.GetMatchConfigManager().GetMaxChip(robot.gameType)
	if robot.userInfo.GetGold() < int64(maxChip) {
		if robot.CurRound >= 3 {
			robot.sendCompareCard(waitTime)
		} else {
			robot.sendFollow(waitTime)
		}
		return
	}

	glog.V(2).Info("==>发送加注")
	msg := &pb.MsgOpCardReq{}
	msg.Type = pb.CardOpType_RAISE.Enum()
	msg.Gold = proto.Int(maxChip)
	robot.sendMsg(pb.MessageId_PLAYER_OPERATE_CARDS, msg, waitTime)
}

// 跟注
func (robot *Robot) sendFollow(waitTime time.Duration) {
	gold := robot.CurSingleBet
	if robot.SeenCard {
		gold *= 2
	}

	if robot.userInfo.GetGold() < int64(gold) {
		if robot.CurRound >= 3 {
			robot.sendCompareCard(waitTime)
		} else {
			robot.sendGiveUp(waitTime)
		}
		return
	}

	glog.V(2).Info("==>发送跟注")
	msg := &pb.MsgOpCardReq{}
	msg.Type = pb.CardOpType_FOLLOW.Enum()
	robot.sendMsg(pb.MessageId_PLAYER_OPERATE_CARDS, msg, waitTime)
}

func (robot *Robot) onOpCards(m []byte) {
	msg := &pb.MsgOpCardRes{}
	err := proto.Unmarshal(m, msg)
	if err != nil {
		glog.Error(err)
		return
	}

	if msg.GetType() == pb.CardOpType_SEE_CARDS {
		if msg.GetUserId() == robot.userInfo.GetUserId() {
			robot.SeenCard = true
		}
	} else if msg.GetType() == pb.CardOpType_GIVE_UP {
		// 弃牌
		p := robot.Players[msg.GetUserId()]
		if p == nil {
			return
		}
		p.IsPlaying = false
		glog.V(2).Info("===>userId:", msg.GetUserId(), " 弃牌")
	} else if msg.GetType() == pb.CardOpType_ALL_IN {
		// 全下
		robot.AllIn = true
		glog.V(2).Info("===>userId:", msg.GetUserId(), " 全下")
		robot.compareTime = time.Now()
	} else if msg.GetType() == pb.CardOpType_RAISE {
		// 加注
		robot.CurSingleBet = int(msg.GetGold())
		glog.V(2).Info("===>userId:", msg.GetUserId(), " 加注gold:", msg.GetGold())
	} else if msg.GetType() == pb.CardOpType_COMPARE {
		robot.compareTime = time.Now()
		loseUserId := ""
		if msg.GetWinnerUserId() == msg.GetUserId() {
			loseUserId = msg.GetCompareUserId()
		} else {
			loseUserId = msg.GetUserId()
		}

		p := robot.Players[loseUserId]
		if p != nil {
			p.IsPlaying = false
		}

		if msg.GetUserId() == robot.userInfo.GetUserId() && msg.GetWinnerUserId() != robot.userInfo.GetUserId() {
			robot.robotLeaveGame()
			return
		}
		if msg.GetCompareUserId() == robot.userInfo.GetUserId() && msg.GetWinnerUserId() != robot.userInfo.GetUserId() {
			robot.robotLeaveGame()
			return
		}
	}
}

func (robot *Robot) onUpdateGold(m []byte) {
	msg := &pb.MsgUpdateGold{}
	err := proto.Unmarshal(m, msg)
	if err != nil {
		glog.Error(err)
		return
	}

	glog.V(2).Info("==>更新金币gold:", msg.GetGold())
	if robot.userInfo == nil {
		return
	}
	robot.userInfo.Gold = msg.Gold
}

func (robot *Robot) onUpdateGoldInGame(m []byte) {

}

func (robot *Robot) sendRobotSetGold(gold int) {
	msg := &pb.MsgRobotSetGold{}
	msg.Gold = proto.Int(gold)
	robot.sendMsg(pb.MessageId_ROBOT_SET_GOLD, msg, time.Millisecond)
}

func (robot *Robot) onRobotMaxCardUser(m []byte) {
	msg := &pb.MsgRobotMaxCardUser{}
	err := proto.Unmarshal(m, msg)
	if err != nil {
		glog.Error(err)
		return
	}

	glog.V(2).Info("===>收到最大牌消息:", msg)
	robot.MaxCardUserId = msg.GetUserId()
}

func (robot *Robot) onRobotCards(m []byte) {
	msg := &pb.MsgRobotCards{}
	err := proto.Unmarshal(m, msg)
	if err != nil {
		glog.Error(err)
		return
	}

	glog.V(2).Info("==>机器人牌msg:", msg)
	robot.Cards = msg.GetCards()
}

func (robot *Robot) onChat(m []byte) {
	msg := &pb.MsgChat{}
	err := proto.Unmarshal(m, msg)
	if err != nil {
		glog.Error(err)
		return
	}

	if msg.GetMessageType() == pb.ChatMessageType_GIFT {
		// 表情
		if msg.GetToUserId() == robot.userInfo.GetUserId() {
			if msg.GetGiftId() == 5 || msg.GetGiftId() == 11 {
				if rand.Int()%100 < 10 {
					if robot.userInfo.GetGender() == pb.Gender_BOY {
						robot.sendGift(msg.GetFromUserId(), 5)
					} else {
						robot.sendGift(msg.GetFromUserId(), 11)
					}
				}
			}
		} else {
			if rand.Int()%100 < 10 {
				robot.sendGift(msg.GetFromUserId(), 0)
			}
		}
	}

	glog.V(2).Info("===>聊天消息msg:", msg)
}

func (robot *Robot) sendJoinWaitQueue() {
	p := robot.Players[robot.userInfo.GetUserId()]
	if p != nil && p.IsPlaying {
		return
	}
	robot.sendMsg(pb.MessageId_JOIN_WAIT_QUEUE, nil, time.Millisecond)
}

func (robot *Robot) sendUseMagicItems() {
	if robot.useFourTimesCard && robot.useForbidCard {
		return
	}

	if !robot.useForbidCard {
		// 没有使用过禁比卡
		msg := &pb.MsgUseMagicItemReq{}
		msg.ItemType = pb.MagicItemType_PROHIBIT_COMPARE.Enum()
		if rand.Int()%100 < 50 {
			robot.useForbidCard = true
			robot.sendMsg(pb.MessageId_USE_MAGIC_ITEM, msg, time.Duration(1+rand.Int()%2)*time.Second)
		}
	}

	if !robot.useFourTimesCard {
		msg := &pb.MsgUseMagicItemReq{}
		msg.ItemType = pb.MagicItemType_FOURFOLD_GOLD.Enum()
		if rand.Int()%100 < 50 {
			robot.useFourTimesCard = true
			robot.sendMsg(pb.MessageId_USE_MAGIC_ITEM, msg, time.Duration(1+rand.Int()%2)*time.Second)
		}
	}
}

func (robot *Robot) rewardInGame() {
	if robot.isRewardInGame {
		return
	}

	if time.Since(robot.compareTime).Seconds() <= 3 {
		return
	}

	msg := &pb.MsgRewardInGame{}
	msg.UserId = proto.String(robot.userInfo.GetUserId())

	robot.sendMsg(pb.MessageId_REWARD_IN_GAME, msg, time.Second)
}

func (robot *Robot) sendGift(userId string, giftId int) {
	if time.Since(robot.compareTime).Seconds() <= 3 {
		return
	}

	if time.Since(robot.lastSendGiftTime).Seconds() < 3 {
		return
	}

	robot.lastSendGiftTime = time.Now()

	msg := &pb.MsgChat{}
	msg.MessageType = pb.ChatMessageType_GIFT.Enum()
	msg.FromUserId = proto.String(robot.userInfo.GetUserId())
	msg.ToUserId = proto.String(userId)
	if giftId > 0 {
		msg.GiftId = proto.Int(giftId)
	} else {
		if robot.userInfo.GetGender() == pb.Gender_BOY {
			msg.GiftId = proto.Int([]int{3, 4, 5, 6}[rand.Int()%4])
		} else {
			msg.GiftId = proto.Int([]int{9, 10, 11, 12}[rand.Int()%4])
		}
	}
	glog.Info("===>发送礼物msg:", msg)
	robot.sendMsg(pb.MessageId_CHAT, msg, time.Duration(1+rand.Int()%3)*time.Second)
}

func (robot *Robot) onRewardInGame(m []byte) {
	msg := &pb.MsgRewardInGame{}
	err := proto.Unmarshal(m, msg)
	if err != nil {
		glog.Error(err)
		return
	}

	if robot.userInfo.GetUserId() == msg.GetUserId() {
		return
	}

	if robot.userInfo.GetUserId() == robot.CurRoundUserId {
		return
	}

	if rand.Int()%100 < 70 {
		robot.rewardInGame()
	}
}

func (robot *Robot) sendExpression() {
	if time.Since(robot.compareTime).Seconds() <= 3 {
		return
	}

	if robot.lastExpressionRound == 0 {
		robot.lastExpressionRound = robot.CurRound + rand.Int()%5
	}

	if robot.CurRound != robot.lastExpressionRound {
		return
	}

	robot.lastExpressionRound = robot.CurRound + rand.Int()%5

	msg := &pb.MsgChat{}
	msg.MessageType = pb.ChatMessageType_DESK.Enum()
	msg.FromUserId = proto.String(robot.userInfo.GetUserId())
	msg.ExpressionId = proto.Int(rand.Int()%20 + 1)

	robot.sendMsg(pb.MessageId_CHAT, msg, time.Duration(3+rand.Int()%2)*time.Second)
}
