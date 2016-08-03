package game

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	domainGame "game/domain/game"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
	"util"
)

func MatchResultHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	if m.GetClient() {
		return nil
	}

	msg := &pb.MsgMatchResult{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	glog.V(2).Info("userId:", player.User.UserId, " 比赛结果:", msg)

	player.MatchRecord.AddEarnGold(int(msg.GetEarnGold()))
	domainUser.GetRankingListUpdater().UpdateUser(player.User.UserId, int(msg.GetEarnGold()))

	// 任意比赛X场
	player.UserTasks.AccomplishTask(util.TaskAccomplishType_PLAY_X_TIMES, 1, player.SendToClientFunc)
	if msg.GetEarnGold() > 0 {
		// 任意比赛胜利X场
		player.UserTasks.AccomplishTask(util.TaskAccomplishType_WIN_X_TIMES, 1, player.SendToClientFunc)
		player.UserTasks.AccomplishTask(util.TaskAccomplishType_TOTAL_EARN_GOLD, int64(msg.GetEarnGold()), player.SendToClientFunc)
		player.UserTasks.AccomplishTask(util.TaskAccomplishType_SINGLE_EARN_GOLD, int64(msg.GetEarnGold()), player.SendToClientFunc)
		player.MatchRecord.WinTimes++

		player.User.LuckyValue -= 50
		if player.User.LuckyValue < 0 {
			player.User.LuckyValue = 0
		}

		if msg.GetMatchType() == int32(util.GameType_WAN_REN) {
			// 万人场胜利
			domainUser.GetPlayerManager().BroadcastClientMsg(int32(pb.MessageId_CHAT), util.BuildSysBugle(fmt.Sprintf("恭喜%v在万众瞩目的“万人场”上大获全胜，赢得了%v万，你不去试试？", player.User.Nickname, msg.GetEarnGold()/10000)))
		}

		if msg.GetEarnGold() > 1000000 {
			clientMsg := &pb.ClientMsg{}
			clientMsg.MsgId = proto.Int(int(pb.MessageId_CHAT))

			if int(msg.GetMatchType()) == int(util.GameType_Common_Level_1) {
				domainUser.GetPlayerManager().BroadcastClientMsg(int32(pb.MessageId_CHAT), util.BuildSysBugle(fmt.Sprintf("恭喜%v在菜鸟场中单场赢钱到%v万!", player.User.Nickname, msg.GetEarnGold()/10000)))
			} else if int(msg.GetMatchType()) == int(util.GameType_Common_Level_2) {
				domainUser.GetPlayerManager().BroadcastClientMsg(int32(pb.MessageId_CHAT), util.BuildSysBugle(fmt.Sprintf("恭喜%v在高手场中单场赢钱到%v万!", player.User.Nickname, msg.GetEarnGold()/10000)))
			} else if int(msg.GetMatchType()) == int(util.GameType_Common_Level_3) {
				domainUser.GetPlayerManager().BroadcastClientMsg(int32(pb.MessageId_CHAT), util.BuildSysBugle(fmt.Sprintf("恭喜%v在精英场中单场赢钱到%v万!", player.User.Nickname, msg.GetEarnGold()/10000)))
			} else if int(msg.GetMatchType()) == int(util.GameType_Common_Level_4) {
				domainUser.GetPlayerManager().BroadcastClientMsg(int32(pb.MessageId_CHAT), util.BuildSysBugle(fmt.Sprintf("恭喜%v在中级场中单场赢钱到%v万!", player.User.Nickname, msg.GetEarnGold()/10000)))
			}
		}
	} else {
		player.MatchRecord.LoseTimes++
	}

	switch int(msg.GetCardType()) {
	case domainGame.CARD_TYPE_BAO_ZI:
		player.UserTasks.AccomplishTask(util.TaskAccomplishType_GAIN_X_BAO_ZI, 1, player.SendToClientFunc)

		if !player.User.IsRobot {
			domainUser.GetPlayerManager().BroadcastClientMsg(int32(pb.MessageId_CHAT), util.BuildSysBugle(fmt.Sprintf("恭喜%v灵光一闪，幸运的在牌桌上翻出了“豹子”通杀全场！", player.User.Nickname)))
		}
	case domainGame.CARD_TYPE_SHUN_JIN:
		player.UserTasks.AccomplishTask(util.TaskAccomplishType_GAIN_X_TONG_HUA_SHUN, 1, player.SendToClientFunc)
	case domainGame.CARD_TYPE_JIN_HUA:
		player.UserTasks.AccomplishTask(util.TaskAccomplishType_GAIN_X_TONG_HUA, 1, player.SendToClientFunc)
	}

	if player.UserLog != nil {
		player.UserLog.MatchTimes++
	}

	if len(msg.GetMaxCards()) == 3 {
		player.MatchRecord.MaxCards = msg.GetMaxCards()
		glog.V(2).Info("==>最大牌maxCards:", player.MatchRecord.MaxCards)
	}

	gameType := util.GameType(int(msg.GetMatchType()))

	if util.IsGameTypeCommon(gameType) {
		// 常规赛X场
		player.UserTasks.AccomplishTask(util.TaskAccomplishType_PLAY_X_COMMON_TIMES, 1, player.SendToClientFunc)
		if msg.GetEarnGold() > 0 {
			// 常规赛胜利X场
			player.UserTasks.AccomplishTask(util.TaskAccomplishType_WIN_X_COMMON_TIMES, 1, player.SendToClientFunc)
		}

		if gameType == util.GameType_Common_Level_1 {
			player.UserTasks.AccomplishTask(util.TaskAccomplishType_PLAY_X_COMMON_LEVEL_1_TIMES, 1, player.SendToClientFunc)
		} else if gameType == util.GameType_Common_Level_2 {
			player.UserTasks.AccomplishTask(util.TaskAccomplishType_PLAY_X_COMMON_LEVEL_2_TIMES, 1, player.SendToClientFunc)
		} else if gameType == util.GameType_Common_Level_3 {
			player.UserTasks.AccomplishTask(util.TaskAccomplishType_PLAY_X_COMMON_LEVEL_3_TIMES, 1, player.SendToClientFunc)
		} else if gameType == util.GameType_Common_Level_4 {
			player.UserTasks.AccomplishTask(util.TaskAccomplishType_PLAY_X_COMMON_LEVEL_4_TIMES, 1, player.SendToClientFunc)
		}
	} else if util.IsGameTypeProps(gameType) {
		player.UserTasks.AccomplishTask(util.TaskAccomplishType_PLAY_X_PROPS_TIMES, 1, player.SendToClientFunc)
		if msg.GetEarnGold() > 0 {
			player.UserTasks.AccomplishTask(util.TaskAccomplishType_WIN_X_PROPS_TIMES, 1, player.SendToClientFunc)
		}

		if gameType == util.GameType_Props_Level_1 {
			player.UserTasks.AccomplishTask(util.TaskAccomplishType_PLAY_X_PROPS_LEVEL_1_TIMES, 1, player.SendToClientFunc)
		} else if gameType == util.GameType_Props_Level_2 {
			player.UserTasks.AccomplishTask(util.TaskAccomplishType_PLAY_X_PROPS_LEVEL_2_TIMES, 1, player.SendToClientFunc)
		} else if gameType == util.GameType_Props_Level_3 {
			player.UserTasks.AccomplishTask(util.TaskAccomplishType_PLAY_X_PROPS_LEVEL_3_TIMES, 1, player.SendToClientFunc)
		}
	} else if util.IsGameTypeSNG(gameType) {
		if msg.GetIsSNGEnd() && msg.GetEarnGold() > 0 {
			if gameType == util.GameType_SNG_Level_1 {
				player.UserTasks.AccomplishTask(util.TaskAccomplishType_WIN_X_SNG_LEVEL_1_TIMES, 1, player.SendToClientFunc)
			} else if gameType == util.GameType_SNG_Level_2 {
				player.UserTasks.AccomplishTask(util.TaskAccomplishType_WIN_X_SNG_LEVEL_2_TIMES, 1, player.SendToClientFunc)
			} else if gameType == util.GameType_SNG_Level_3 {
				player.UserTasks.AccomplishTask(util.TaskAccomplishType_WIN_X_SNG_LEVEL_3_TIMES, 1, player.SendToClientFunc)
			}
		}

	}

	player.SetUserCache()
	domainUser.GetUserFortuneManager().UpdateUserFortune(player.User.UserId)

	return nil
}
