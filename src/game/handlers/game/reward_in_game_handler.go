package game

import (
	"fmt"
	domainGame "game/domain/game"
	domainUser "game/domain/user"
	"game/server"
	"pb"
	"time"
	"util"
)

func RewardInGameHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	if !domainGame.GetDeskManager().OnRewardInGame(player.User.UserId) {
		return nil
	}

	// 打赏成功
	player.UserTasks.AccomplishTask(util.TaskAccomplishType_REWARD_OTHER_X_TIMES, 1, player.SendToClientFunc)

	if !util.CompareDate(time.Now(), player.User.LastRewardInGameTime) {
		player.User.LastRewardInGameTime = time.Now()
		player.User.RewardInGameTimes = 1
	}

	if player.User.RewardInGameTimes == 100 || player.User.RewardInGameTimes == 200 ||
		player.User.RewardInGameTimes == 500 || player.User.RewardInGameTimes == 1000 {

		domainUser.GetPlayerManager().BroadcastClientMsg(int32(pb.MessageId_CHAT), util.BuildSysBugle(fmt.Sprintf("恭喜%v今日打赏了%v次！有钱！任性！", player.User.Nickname, player.User.RewardInGameTimes)))
	}

	return nil
}
