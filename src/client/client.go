package main

import (
	"code.google.com/p/go.net/websocket"
	"code.google.com/p/goprotobuf/proto"
	//	"config"
	"log"
	"os"
	"os/signal"
	"pb"
	"syscall"
	"util"
)

var ws *websocket.Conn

func sendMsg(msgId pb.MessageId, body proto.Message) {
	if ws == nil {
		log.Fatal("ws nil")
	}

	log.Println("==>发送消息msgId:", msgId)

	msg := &pb.ClientMsg{}
	msg.MsgId = proto.Int(int(msgId))

	if body != nil {
		b, err := proto.Marshal(body)
		if err != nil {
			log.Fatal(err)
		}

		msg.MsgBody = b
	}

	b, err := proto.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}

	ws.Write(b)
}

func receive() {
	return
	for {
		var data []byte
		if ws == nil {
			log.Fatal("ws nil")
		}
		err := websocket.Message.Receive(ws, &data)
		if err != nil {
			log.Fatal("error receiving msg:", err)
			break
		}

		clientMsg := &pb.ClientMsg{}
		err = proto.Unmarshal(data, clientMsg)
		if err != nil {
			log.Fatal("unmarshal client msg failed!")
			break
		}

		handleMsg(clientMsg)
	}
}

func main() {
	//	util.PushNotification("6e462d36-a77a-4e07-a5c0-c17e81b9a307", bson.NewObjectId().Hex())
	//	util.PushNotificationEveryone("hello,大家好~")
	//	util.PushNotificationChannels(`["user"]`, "user频道用户大家好~")
	//	util.SubscribeChannels("fphjd9KLxk", `[]`)
	//	time.Sleep(3 * time.Second)
	//	flag.Parse()
	//	glog.V(2).Info("client启动...")
	//	glog.V(0).Info("V0消息----")
	//	return

	origin := "http://192.168.1.63/"
	url := "ws://192.168.1.63:8002/ws/"

	//	origin := "http://203.195.170.83/"
	//	url := "ws://203.195.170.83:8002/ws/"

	//	origin := "http://123.57.17.118/"
	//	url := "ws://123.57.17.118:8002/ws/"
	//url := "ws://203.195.147.120:8080/ws/"
	conn, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal("错误:", err)
	}
	defer conn.Close()

	ws = conn

	go receive()

	doLogic()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	<-c
}

func doLogic() {
	sendLoginReq()
}

func handleMsg(msg *pb.ClientMsg) {
	log.Println("收到消息 id:", util.GetMsgIdName(msg.GetMsgId()))
	switch pb.MessageId(msg.GetMsgId()) {
	case pb.MessageId_LOGIN:
		onLoginRes(msg.GetMsgBody())
	case pb.MessageId_GET_USER_INFO:
		onGetUserInfo(msg.GetMsgBody())
	case pb.MessageId_ENTER_POKER_DESK:
		onEnterGame(msg.GetMsgBody())
	case pb.MessageId_GET_POKER_DESK_INFO:
		onGetGameInfo(msg.GetMsgBody())
	case pb.MessageId_LEAVE_POKER_DESK:
		onLeaveGame(msg.GetMsgBody())
	case pb.MessageId_POKER_DESK_GAME_BEGIN:
		onGameStart(msg.GetMsgBody())
	case pb.MessageId_PLAYER_OPERATE_CARDS:
		onOpCardRes(msg.GetMsgBody())
	case pb.MessageId_GET_RANKING_LIST:
		onGetRankingList(msg.GetMsgBody())
	case pb.MessageId_EXCHANGE_GOLD:
		onExchangeGold(msg.GetMsgBody())
	case pb.MessageId_CHAT:
		onChatMsg(msg.GetMsgBody())
	}
}

func sendLoginReq() {
	msg := &pb.MsgLoginReq{}
	msg.Username = proto.String("j8")
	msg.Nickname = proto.String("穷人得瑟")
	msg.Userpwd = proto.String("穷人得瑟")
	//	msg.RobotKey = proto.String(config.RobotKey)
	sendMsg(pb.MessageId_LOGIN, msg)
}

func onLoginRes(m []byte) {
	msg := &pb.MsgLoginRes{}
	err := proto.Unmarshal(m, msg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("登录结果:", msg.GetCode())

	//	sendGetUserInfo()
	//
	//	sendEnterGame()
	//	sendGetRankingList()
	//
	//	sendChatMsg()
}

func sendGetUserInfo() {
	log.Println("===>获取用户信息:")
	sendMsg(pb.MessageId_GET_USER_INFO, nil)
}

func onGetUserInfo(m []byte) {
	msg := &pb.MsgGetUserInfoRes{}
	err := proto.Unmarshal(m, msg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("获取玩家信息成功:")
}

func sendEnterGame() {
	msg := &pb.MsgEnterPokerDeskReq{}
	msg.Type = pb.MatchType_COMMON_LEVEL1.Enum()
	sendMsg(pb.MessageId_ENTER_POKER_DESK, msg)
	//	msg.Type = pb.MatchType_WAN_REN_GAME.Enum()
	//	sendMsg(pb.MessageId_ENTER_POKER_DESK, msg)
}

func onEnterGame(m []byte) {
	msg := &pb.MsgEnterPokerDeskRes{}
	err := proto.Unmarshal(m, msg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("进入游戏结果:", msg)
	//	sendGetGameInfo()

	if msg.GetType() == pb.MatchType_WAN_REN_GAME {
		leaveMsg := &pb.MsgLeavePokerDeskReq{}
		leaveMsg.Type = pb.MatchType_WAN_REN_GAME.Enum()
		sendMsg(pb.MessageId_LEAVE_POKER_DESK, leaveMsg)
		log.Println("离开万人场")
	}
}

func sendGetGameInfo() {
	sendMsg(pb.MessageId_GET_POKER_DESK_INFO, nil)
}

func onGetGameInfo(m []byte) {
	//	sendLeaveGame()
	if m != nil {
		msg := &pb.MsgGetPokerDeskInfoRes{}
		err := proto.Unmarshal(m, msg)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("获取游戏信息成功:")
		return
	}
	log.Println("获取游戏信息返回nil")
}

func sendLeaveGame() {
	log.Println("==>发送离开游戏消息")
	sendMsg(pb.MessageId_LEAVE_POKER_DESK, nil)
}

func onLeaveGame(m []byte) {
	msg := &pb.MsgLeavePokerDesk{}
	err := proto.Unmarshal(m, msg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("收到玩家离开游戏消息:", msg)
}

func onGameStart(m []byte) {
	msg := &pb.MsgGameStart{}
	err := proto.Unmarshal(m, msg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("onGameStart:", msg)
}

func sendRaiseMsg(userId string) {
	msg := &pb.MsgOpCardReq{}
	msg.Type = pb.CardOpType_COMPARE.Enum()
	msg.CompareUserId = proto.String(userId)
	sendMsg(pb.MessageId_PLAYER_OPERATE_CARDS, msg)
}

func onOpCardRes(m []byte) {
	msg := &pb.MsgOpCardRes{}
	err := proto.Unmarshal(m, msg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("==>收到游戏操作包:", msg)
}

func sendGetRankingList() {
	sendExchangeGold()

	msg := &pb.MsgGetRankingListReq{}
	msg.Types = append(msg.Types, pb.RankingType_RECHARGE_TODAY)
	msg.Types = append(msg.Types, pb.RankingType_RECHARGE_YESTERDAY)
	msg.Types = append(msg.Types, pb.RankingType_RECHARGE_LAST_WEEK)
	msg.Types = append(msg.Types, pb.RankingType_EARNINGS_TODAY)
	msg.Types = append(msg.Types, pb.RankingType_EARNINGS_YESTERDAY)
	msg.Types = append(msg.Types, pb.RankingType_EARNINGS_LAST_WEEK)
	sendMsg(pb.MessageId_GET_RANKING_LIST, msg)
}

func onGetRankingList(m []byte) {
	msg := &pb.MsgGetRankingListRes{}
	err := proto.Unmarshal(m, msg)
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range msg.GetItems() {
		if item.GetType() == pb.RankingType_RECHARGE_TODAY {
			log.Println("排行榜:", item)
		}
	}
}

func sendExchangeGold() {
	msg := &pb.MsgExchangeGoldReq{}
	msg.Diamond = proto.Int(30)
	sendMsg(pb.MessageId_EXCHANGE_GOLD, msg)
}

func onExchangeGold(m []byte) {
	msg := &pb.MsgExchangeGoldRes{}
	err := proto.Unmarshal(m, msg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("==>兑换金币成功:", msg)
}

func sendChatMsg() {
	msg := &pb.MsgChat{}
	msg.MessageType = pb.ChatMessageType_BUGLE.Enum()
	sendMsg(pb.MessageId_CHAT, msg)
}

func onChatMsg(m []byte) {
	msg := &pb.MsgChat{}
	err := proto.Unmarshal(m, msg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("聊天消息:", msg.GetContent())
}
