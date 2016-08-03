package slots

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	domainGame "game/domain/game"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
	"sync"
	"util"
	domainSlots "game/domain/slots"
)

var slotMachineGameLogic *domainGame.GameLogic
var mu sync.RWMutex

func init() {
	slotMachineGameLogic = domainGame.NewGameLogic(0)
	slotMachineGameLogic.ShuffleCards2()
	mu = sync.RWMutex{}
}

func GetRandomCards(m *domainSlots.SlotMachine) []byte {
	mu.Lock()
	defer mu.Unlock()

    poolValue := domainSlots.GetSlotGlobal().GetPoolValue()
    configs := m.UpdatePerByPoolValue(poolValue, m.Coin)
	card_type := m.RandomCardType(configs)

	return m.GetSlotCards(card_type)
}

func GetRandomCard(pos int, m *domainSlots.SlotMachine) (byte, bool) {
	mu.Lock()
	defer mu.Unlock()

	configs_card, ok := m.CheckUpdate(pos)
	if !ok {
	    glog.Error("GetRandomCard CheckUpdate error")
	    return 0, false
	}
	poolValue := domainSlots.GetSlotGlobal().GetPoolValue()
	configs_pool := m.UpdatePerByPoolValue(poolValue, m.Coin)

	configs := m.MergeCardPer(configs_pool, configs_card)

	card_type := m.RandomCardType(configs)
	card, isOk := m.UpdateCardByType(card_type, pos)

	return card, isOk
}

const MIN_GOLD = 10000

func PlaySlotMachineHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	msg := &pb.Msg_SlotMachinesPlayReq{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	if msg.GetCoin() < MIN_GOLD {
		glog.V(2).Info("userId:", player.User.UserId, " coin:", msg.GetCoin(), " 小于10000")
		return nil
	}

	if msg.GetCoin()%MIN_GOLD != 0 {
		glog.V(2).Info("userId:", player.User.UserId, " coin:", msg.GetCoin(), " 非10000整数倍")
		return nil
	}

	res := &pb.Msg_SlotMachinesPlayRes{}

	_, _, ok := domainUser.GetUserFortuneManager().ConsumeGold(player.User.UserId, int64(msg.GetCoin()), false, "老虎机")
	if !ok {
		glog.V(2).Info("老虎机扣钱失败userId:", player.User.UserId, " consumeGold:", msg.GetCoin())
		res.Code = pb.Msg_SlotMachinesPlayRes_LACK_GOLD.Enum()
		return server.BuildClientMsg(m.GetMsgId(), res)
	}

	domainSlots.GetSlotGlobal().InputPool(int(msg.GetCoin()))

	player.SlotMachine.Coin = int(msg.GetCoin())
	player.SlotMachine.SetCards(GetRandomCards(player.SlotMachine))

	cardType := domainGame.GetCardType(player.SlotMachine.GetCards())
	switch cardType {
	case domainGame.CARD_TYPE_BAO_ZI:
		player.UserTasks.AccomplishTask(util.TaskAccomplishType_SLOT_MACHINE_X_BAO_ZI, 1, player.SendToClientFunc)
		domainUser.GetPlayerManager().BroadcastClientMsg(int32(pb.MessageId_CHAT), util.BuildSysBugle(fmt.Sprintf("恭喜%v在一夜暴富中奖！豹子！豹子！豹子！", player.User.Nickname)))
	case domainGame.CARD_TYPE_SHUN_JIN:
		player.UserTasks.AccomplishTask(util.TaskAccomplishType_SLOT_MACHINE_X_TONG_HUA_SHUN, 1, player.SendToClientFunc)
	case domainGame.CARD_TYPE_JIN_HUA:
		player.UserTasks.AccomplishTask(util.TaskAccomplishType_SLOT_MACHINE_X_TONG_HUA, 1, player.SendToClientFunc)
	case domainGame.CARD_TYPE_SPECIAL:
		player.UserTasks.AccomplishTask(util.TaskAccomplishType_SLOT_MACHINE_X_SPECIAL, 1, player.SendToClientFunc)
	}

	domainUser.GetUserFortuneManager().UpdateUserFortune(player.User.UserId)

	res.Code = pb.Msg_SlotMachinesPlayRes_OK.Enum()
	res.Card1 = proto.Int(int(player.SlotMachine.GetCard(0)))
	res.Card2 = proto.Int(int(player.SlotMachine.GetCard(1)))
	res.Card3 = proto.Int(int(player.SlotMachine.GetCard(2)))

	glog.V(2).Info("====>老虎机牌userId:", player.User.UserId, " cards:", player.SlotMachine.GetCards(), " res:", res)

	return server.BuildClientMsg(m.GetMsgId(), res)
}
