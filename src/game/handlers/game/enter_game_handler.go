package game

import (
	"code.google.com/p/goprotobuf/proto"
	domainGame "game/domain/game"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
	"util"
)

func EnterGameHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	glog.Info("EnterGameHandler in.")
	player := domainUser.GetPlayer(sess.Data)

	msg := &pb.MsgEnterPokerDeskReq{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		return nil
	}

	glog.V(2).Info("===>进入游戏msg:", msg)

	// 限制检查
	var gameType util.GameType

	switch msg.GetType() {
	case pb.MatchType_COMMON_LEVEL1:
		gameType = util.GameType_Common_Level_1
	case pb.MatchType_COMMON_LEVEL2:
		gameType = util.GameType_Common_Level_2
	case pb.MatchType_COMMON_LEVEL3:
		gameType = util.GameType_Common_Level_3
	case pb.MatchType_COMMON_LEVEL4:
		gameType = util.GameType_Common_Level_4
	/* del by wangsq --- 暂时取消道具场和淘汰赛工具
	case pb.MatchType_MAGIC_ITEM_LEVEL1:
		gameType = util.GameType_Props_Level_1
	case pb.MatchType_MAGIC_ITEM_LEVEL2:
		gameType = util.GameType_Props_Level_2
	case pb.MatchType_MAGIC_ITEM_LEVEL3:
		gameType = util.GameType_Props_Level_3
	case pb.MatchType_SNG_LEVEL1:
		gameType = util.GameType_SNG_Level_1
	case pb.MatchType_SNG_LEVEL2:
		gameType = util.GameType_SNG_Level_2
	case pb.MatchType_SNG_LEVEL3:
		gameType = util.GameType_SNG_Level_3
		now := time.Now()
		if now.Hour() < 20 || now.Hour() > 22 {
			glog.V(2).Info("===>淘汰赛大师场不在比赛时间!")
			return nil
		}
	case pb.MatchType_WAN_REN_GAME:
		gameType = util.GameType_WAN_REN
	*/
	default:
		glog.V(2).Info("非法游戏类型!")
		return nil
	}

	player.ResetLuckyValue()

	lastGameId := 0
	if msg.GetChangeDesk() {
		// 离开游戏
		lastGameId = player.LastGameId
	}

	if domainGame.GetDeskManager().IsPlayingWanRen(player.User.UserId) {
		domainGame.GetDeskManager().LeaveGame(player.User.UserId, util.GameType_WAN_REN, false)
	} else {
		if domainGame.GetDeskManager().IsPlayingNormal(player.User.UserId) {
			if gameType != util.GameType_WAN_REN {
				// 非进入万人场直播,退出当前游戏
				domainGame.GetDeskManager().LeaveGame(player.User.UserId, 0, msg.GetChangeDesk())
			}
		}
	}

	res := &pb.MsgEnterPokerDeskRes{}
	res.Type = msg.Type

	gameId, ok := domainGame.GetDeskManager().EnterGame(gameType, player.User.BuildMessage(player.MatchRecord.BuildMessage()), lastGameId)
	if !ok {
		res.Code = pb.MsgEnterPokerDeskRes_FAILED.Enum()
		return server.BuildClientMsg(m.GetMsgId(), res)
	}

	player.LastGameId = gameId
	res.Code = pb.MsgEnterPokerDeskRes_OK.Enum()

	return server.BuildClientMsg(m.GetMsgId(), res)
}
