package user

import (
	"code.google.com/p/goprotobuf/proto"
	"config"
	"fmt"
	domainActive "game/domain/iosActive"
	"game/domain/rankingList"
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"pb"
	"strconv"
	"sync"
	"time"
	"util"
)

type UserFortuneManager struct {
	sync.RWMutex
	fortune              map[string]*UserFortune
	UpdateGoldInGameFunc func(userId string)
}

var userFortuneManager *UserFortuneManager

func init() {
	userFortuneManager = &UserFortuneManager{}
	userFortuneManager.fortune = make(map[string]*UserFortune)
}

func GetUserFortuneManager() *UserFortuneManager {
	return userFortuneManager
}

func (m *UserFortuneManager) LoadUserFortune(userId string, isRobot bool) bool {
	f, err := FindUserFortune(userId)
	if err != nil && err != mgo.ErrNotFound {
		glog.Error(err)
		return false
	}

	m.Lock()
	defer m.Unlock()

	f.IsRobot = isRobot
	m.fortune[userId] = f

	if !util.CompareDate(f.LastPayTime, time.Now()) {
		f.GainedFirstRechargeBonus = false
	}

	return true
}

func (m *UserFortuneManager) UnloadUserFortune(userId string) {
	m.Lock()
	defer m.Unlock()

	f := m.fortune[userId]
	delete(m.fortune, userId)

	if f.IsRobot {
		return
	}

	if f != nil {
		err := SaveFortune(f)
		if err != nil {
			glog.Error("保存用户财富信息失败userId:", userId, " f:", f)
		}
	}
}

func (m *UserFortuneManager) GetUserFortune(userId string) (UserFortune, bool) {
	m.RLock()
	defer m.RUnlock()

	f := m.fortune[userId]
	if f == nil {
		return UserFortune{}, false
	}

	return *f, true
}

func (m *UserFortuneManager) checkUserFortune(userId string) {
	_, ok := m.GetUserFortune(userId)
	if !ok {
		m.LoadUserFortune(userId, false)
	}
}

func (m *UserFortuneManager) SaveUserFortune(userId string) {
	m.Lock()
	f := m.fortune[userId]
	m.Unlock()

	if f == nil {
		return
	}

	SaveFortune(f)
}

func (m *UserFortuneManager) EarnFortune(userId string, gold int64, diamond, score int, isRecharge bool, reason string) bool {
	m.checkUserFortune(userId)

	m.Lock()
	defer m.Unlock()

	glog.Info("===>EarnFortune userId:", userId, " gold:", gold, " diamond:", diamond, " reason:", reason)

	fortune := m.fortune[userId]
	if fortune == nil {
		glog.Error("EarnFortune failed userId:", userId, " gold:", gold, " reason:", reason)
		return false
	}

	useGiftBag := false
	now := time.Now()
	if isRecharge && diamond == 10 && fortune.BuyDailyGiftBag && !util.CompareDate(fortune.DailyGiftBagUseTime, now) {
		fortune.Gold += 150000
		useGiftBag = true
		fortune.BuyDailyGiftBag = false
		fortune.DailyGiftBagUseTime = now
		fortune.isDailyGiftBagPay = true
	} else {
		fortune.Diamond += diamond
		fortune.Gold += int64(gold)
		fortune.BuyDailyGiftBag = false
	}

	fortune.Score += score
	if reason == "邮件奖励" {
		fortune.Charm += score
	}

	if fortune.Gold < 0 {
		fortune.Gold = 0
	}
	if fortune.Diamond < 0 {
		fortune.Diamond = 0
	}

	if isRecharge && diamond > 0 {
		if !util.CompareDate(fortune.LastPayTime, now) {
			fortune.GainedFirstRechargeBonus = false
		}
		fortune.LastPayTime = now
	}

	if !fortune.IsRobot {
		l := &FortuneLog{}
		l.UserId = userId
		l.Gold = int(gold)
		l.CurGold = fortune.Gold
		l.Diamond = diamond
		l.CurDiamond = fortune.Diamond
		l.Reason = reason
		SaveEarnFortuneLog(l)
	}

	if useGiftBag {
		GetPlayerManager().SendClientMsg(userId, int32(pb.MessageId_BUY_DAILY_GIFT_BAG_OK), nil)
		return true
	}

	rechargeBonus := 0
	isFirstRecharge := false
	if isRecharge && diamond > 0 {
		switch diamond {
		case 10:
			if !fortune.FirstRecharge10 {
				fortune.FirstRecharge10 = true
				rechargeBonus = 0
				isFirstRecharge = true
			} else {
				rechargeBonus = 0
			}
		case 30:
			if !fortune.FirstRecharge30 {
				fortune.FirstRecharge30 = true
				rechargeBonus = 0
				isFirstRecharge = true
			} else {
				rechargeBonus = 0
			}
		case 50:
			if !fortune.FirstRecharge50 {
				fortune.FirstRecharge50 = true
				isFirstRecharge = true
				rechargeBonus = 0
			} else {
				rechargeBonus = 0
			}
		case 100:
			if !fortune.FirstRecharge100 {
				fortune.FirstRecharge100 = true
				isFirstRecharge = true
				rechargeBonus = 0
			} else {
				rechargeBonus = 0
			}
		case 500:
			if !fortune.FirstRecharge500 {
				fortune.FirstRecharge500 = true
				isFirstRecharge = true
				rechargeBonus = 0
			} else {
				rechargeBonus = 0
			}
		}
	}

	if rechargeBonus > 0 {

		prizeMail := &pb.PrizeMailDef{}

		prizeMail.MailId = proto.String(bson.NewObjectId().Hex())
		if isFirstRecharge {
			prizeMail.Content = proto.String(fmt.Sprintf("首次购买%v钻石的加赠奖励", diamond))
		} else {
			prizeMail.Content = proto.String(fmt.Sprintf("购买%v钻石的加赠奖励", diamond))
		}
		prizeMail.Prize = &pb.PrizeDef{}
		prizeMail.Prize.Gold = proto.Int(rechargeBonus)

		GetPlayerManager().SendServerMsg("", []string{userId}, int32(pb.ServerMsgId_MQ_PRIZE_MAIL), prizeMail)
	}

	if gold != 0 {
		m.updateGoldInGame(userId)
		SaveGold(userId, fortune.Gold, fortune.IsRobot)
	}

	if diamond != 0 {
		SaveDiamond(userId, fortune.Diamond, fortune.IsRobot)
	}

	return true
}

func (m *UserFortuneManager) updateGoldInGame(userId string) {
	if m.UpdateGoldInGameFunc != nil {
		go m.UpdateGoldInGameFunc(userId)
	}
}

func (m *UserFortuneManager) BuyDoubleCard(userId string, diamond int, count int) bool {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return false
	}

	if fortune.Diamond < diamond {
		return false
	}
	fortune.Diamond -= diamond

	fortune.DoubleCardCount += count

	return true
}

func (m *UserFortuneManager) BuyForbidCard(userId string, diamond int, count int) bool {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return false
	}

	if fortune.Diamond < diamond {
		return false
	}

	fortune.Diamond -= diamond
	fortune.ForbidCardCount += count

	return true
}

func (m *UserFortuneManager) BuyChangeCard(userId string, diamond int, count int) bool {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return false
	}

	if fortune.Diamond < diamond {
		return false
	}

	fortune.Diamond -= diamond

	fortune.ChangeCardCount += count

	m.updateUserFortune(userId, 0)

	return true
}

func (m *UserFortuneManager) ConsumeDoubleCard(userId string, count int) bool {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return false
	}

	if fortune.DoubleCardCount < count {
		return false
	}

	fortune.DoubleCardCount -= count

	return true
}

func (m *UserFortuneManager) ConsumeForbidCard(userId string, count int) bool {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return false
	}

	if fortune.ForbidCardCount < count {
		return false
	}

	fortune.ForbidCardCount -= count

	return true
}

func (m *UserFortuneManager) ConsumeChangeCard(userId string, count int) bool {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return false
	}

	if fortune.ChangeCardCount < count {
		return false
	}

	fortune.ChangeCardCount -= count

	return true
}

func (m *UserFortuneManager) ConsumeGold(userId string, gold int64, consumeAllIfNotEnough bool, reason string) (int64, int, bool) {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return 0, 0, false
	}

	if gold > 0 {
		if fortune.Gold < int64(gold) {
			if consumeAllIfNotEnough {
				oldGold := fortune.Gold
				fortune.Gold = 0
				gold = int64(oldGold)
			} else {
				return 0, 0, false
			}
		} else {
			fortune.Gold -= int64(gold)
		}
	}

	if !fortune.IsRobot {
		l := &FortuneLog{}
		l.UserId = userId
		l.Gold = int(gold)
		l.CurGold = fortune.Gold
		l.Diamond = 0
		l.CurDiamond = fortune.Diamond
		l.Reason = reason
		SaveConsumeFortuneLog(l)
	}

	if gold > 0 {
		m.updateGoldInGame(userId)
		SaveGold(userId, fortune.Gold, fortune.IsRobot)
	}

	m.updateUserFortune(userId, 0)

	return fortune.Gold, int(gold), true
}

func (m *UserFortuneManager) ConsumeGoldNoMsg(userId string, gold int64, consumeAllIfNotEnough bool, reason string) (int64, int, bool) {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return 0, 0, false
	}

	if gold > 0 {
		if fortune.Gold < int64(gold) {
			if consumeAllIfNotEnough {
				oldGold := fortune.Gold
				fortune.Gold = 0
				gold = int64(oldGold)
			} else {
				return 0, 0, false
			}
		} else {
			fortune.Gold -= int64(gold)
		}
	}

	if !fortune.IsRobot {
		l := &FortuneLog{}
		l.UserId = userId
		l.Gold = int(gold)
		l.CurGold = fortune.Gold
		l.Diamond = 0
		l.CurDiamond = fortune.Diamond
		l.Reason = reason
		SaveConsumeFortuneLog(l)
	}

	if gold > 0 {
		m.updateGoldInGame(userId)
		SaveGold(userId, fortune.Gold, fortune.IsRobot)
	}

	return fortune.Gold, int(gold), true
}

func (m *UserFortuneManager) AddExp(userId string, exp int) {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return
	}

	fortune.Exp += exp
}

func (m *UserFortuneManager) AddScore(userId string, score int) {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return
	}

	fortune.Score += score
}

// add by wangsq start
func (m *UserFortuneManager) AddFish(userId string, fromId string, fishType int, count int) {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return
	}

	if fortune.FishInfo == nil {
		fortune.FishInfo = map[string]UserFish{}
	}

	isNew := true
	for key, value := range fortune.FishInfo {
		if value.FromUid == fromId && value.FishType == fishType {
			value.Count += count
			value.GetTime = time.Now()
			fortune.FishInfo[key] = value
			m.fortune[userId] = fortune
			isNew = false
			break
		}
	}

	if isNew {
		recordId := m.getFishRecordId(fortune.FishInfo)
		fish := UserFish{recordId, fromId, time.Now(), fishType, count}
		fortune.FishInfo[recordId] = fish
		m.fortune[userId] = fortune
	}
}

func (m *UserFortuneManager) AddFish2Mongo(userId string, fromId string, fishType int, count int) error {
	m.Lock()
	defer m.Unlock()

	fortune, err := FindUserFortune(userId)
	if err != nil {
		return err
	}

	isNew := true
	for key, value := range fortune.FishInfo {
		if value.FromUid == fromId && value.FishType == fishType {
			value.Count += count
			value.GetTime = time.Now()
			fortune.FishInfo[key] = value
			SaveFortune(fortune)
			isNew = false
			break
		}
	}

	if isNew {
		recordId := m.getFishRecordId(fortune.FishInfo)
		fish := UserFish{recordId, fromId, time.Now(), fishType, count}
		fortune.FishInfo[recordId] = fish
		SaveFortune(fortune)
	}

	return nil
}

func (m *UserFortuneManager) getFishRecordId(info map[string]UserFish) string {
	max := 0
	for k, _ := range info {
		ki, _ := strconv.Atoi(k)
		if ki > max {
			max = ki
		}
	}
	return strconv.Itoa(max + 1)
}

func (m *UserFortuneManager) AddGold2Mongo(userId string, gold int64, reason string) (int64, bool) {
	m.Lock()
	defer m.Unlock()

	fortune, err := FindUserFortune(userId)
	if err != nil {
		glog.Info("FindUserFortune error, err=", err)
		return 0, false
	}

	fortune.Gold += gold
	SaveFortune(fortune)

	l := &FortuneLog{}

	l.UserId = userId
	l.Gold = int(gold)
	l.CurGold = fortune.Gold
	l.Diamond = 0
	l.CurDiamond = fortune.Diamond
	l.Reason = reason
	SaveEarnFortuneLog(l)

	return fortune.Gold, true
}

func (m *UserFortuneManager) DelFish(userId string, recordId string, count int) bool {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	glog.Info("DelFish in.", fortune.FishInfo)
	if fortune == nil {
		return false
	}

	for key, value := range fortune.FishInfo {
		if value.RecordId == recordId {
			if value.Count == count {
				delete(fortune.FishInfo, key)
			} else {
				value.Count -= count
				fortune.FishInfo[key] = value
			}
			m.fortune[userId] = fortune
			glog.Info("DelFish true out.", m.fortune[userId].FishInfo)
			return true
		}
	}

	glog.Info("DelFish false out.", m.fortune[userId].FishInfo)
	return false
}

func (m *UserFortuneManager) AddToken(userId string, token string) {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return
	}

	fishToken := FToken{token, int(time.Now().Unix())}
	fortune.FishToken = fishToken
}

func (m *UserFortuneManager) UpdateVipTaskState(userId string, level int) {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return
	}

	levelKey := strconv.Itoa(level)
	vipState, ok := fortune.VipTaskStates[levelKey]
	if ok {
		sttime := util.GetDayZero()
		vipState.StartTime = sttime
		fortune.VipTaskStates[levelKey] = vipState
	} else {
		sttime := util.GetDayZero()
		state := UserVipState{level, sttime, 0}
		fortune.VipTaskStates[levelKey] = state
	}

	m.fortune[userId] = fortune
}

func (m *UserFortuneManager) GainVipTask(userId string, level int) {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return
	}

	levelKey := strconv.Itoa(level)
	vipState, ok := fortune.VipTaskStates[levelKey]
	if ok {
		gainTime := util.GetDayZero()
		vipState.LastGainTime = gainTime
		fortune.VipTaskStates[levelKey] = vipState
	}

	m.fortune[userId] = fortune
}

func (m *UserFortuneManager) InitVipTaskState(userId string) {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return
	}

	fortune.VipTaskStates = make(map[string]UserVipState)
	vip_configs := config.GetVipPriceConfigManager().GetVipConfig()
	for _, value := range vip_configs {
		level := value.Level
		sttime := 0
		lastGainTime := 0
		vipState := UserVipState{level, int64(sttime), int64(lastGainTime)}
		levelKey := strconv.Itoa(level)
		fortune.VipTaskStates[levelKey] = vipState
	}
}

func (m *UserFortuneManager) ConsumeCharm(userId string, charm int) bool {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return false
	}

	if fortune.Charm < charm {
		return false
	}

	fortune.Charm -= charm
	m.fortune[userId] = fortune

	u := GetUserCache().GetUser(userId)
	rankingList.GetRankingList().UpdateRankingItem(rankingList.RankingType_Charm, u, int64(fortune.Charm))

	return true
}

func (m *UserFortuneManager) LoginSafeBox(userId string, pwd string) *pb.Msg_LoginSafeBoxRes_ResCode {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return pb.Msg_LoginSafeBoxRes_FAILED.Enum()
	}

	if fortune.SafeBox.Pwd != pwd {
		return pb.Msg_LoginSafeBoxRes_PWDERR.Enum()
	}

	return pb.Msg_LoginSafeBoxRes_OK.Enum()
}

func (m *UserFortuneManager) ChangePwdSafeBox(userId, oldpwd, newpwd string) (int, string) {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return -1, "获取信息失败"
	}

	if fortune.SafeBox.Pwd != oldpwd {
		return -2, "保管箱密码错误"
	}

	fortune.SafeBox.Pwd = newpwd
	m.fortune[userId] = fortune

	return 0, ""
}

func (m *UserFortuneManager) ResetPwdSafeBox(userId, newpwd string) int {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return -1
	}

	if !fortune.SafeBox.IsSetPwd {
		fortune.SafeBox.IsSetPwd = true
	}

	glog.Info("ResetPwdSafeBox in,newpwd=", newpwd)
	fortune.SafeBox.Pwd = newpwd
	m.fortune[userId] = fortune
	glog.Info("m.fortune[userId].SafeBox.Pwd=", m.fortune[userId].SafeBox.Pwd)

	return 0
}

func (m *UserFortuneManager) OpenSafeBox(userId string, pwd string) bool {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return false
	}

	if fortune.SafeBox.IsOpen == true {
		return true
	}

	fortune.SafeBox.IsOpen = true
	fortune.SafeBox.Pwd = pwd
	m.fortune[userId] = fortune

	return true
}

func (m *UserFortuneManager) UpdateFirstRecharge(userId string) {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return
	}
	fortune.FirstRecharge = true
	m.fortune[userId] = fortune

	return
}

func (m *UserFortuneManager) GetFirstRecharge(userId string) bool {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return false
	}

	return fortune.FirstRecharge
}

func (m *UserFortuneManager) UpdateSavings(userId, pwd string, gold int64, reason string, toUid string) (int, string, int64, int64, *BoxLog) {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return -1, "用户信息错误", 0, 0, nil
	}

	if reason != "取款" && !fortune.SafeBox.IsOpen {
		return -2, "保管箱未激活，请到商场购买", 0, 0, nil
	}

	if reason == "取款" || reason == "赠送扣款" {
		if fortune.SafeBox.Pwd != pwd {
			return -3, "保管箱密码错误", 0, 0, nil
		}
	}

	if fortune.SafeBox.Savings+gold < 0 {
		return -4, "您的存款不足", 0, 0, nil
	}

	fortune.SafeBox.Savings += gold
	log := BoxLog{reason, gold, fortune.SafeBox.Savings, time.Now().Local(), toUid, "0"}
	if fortune.SafeBox.BoxLogs == nil {
		fortune.SafeBox.BoxLogs = map[string]BoxLog{}
	}
	logId := strconv.Itoa(len(fortune.SafeBox.BoxLogs) + 1)
	log.LogId = logId
	fortune.SafeBox.BoxLogs[logId] = log
	m.fortune[userId] = fortune

	return 0, "", fortune.SafeBox.Savings, fortune.Gold, &log
}

//存保管箱离线
func (m *UserFortuneManager) UpdateSavingsAddOffLine(fromUserId string, toUserId string, gold int64) bool {
	m.checkUserFortune(toUserId)
	m.Lock()
	defer m.Unlock()
	fortune := m.fortune[toUserId]
	if fortune == nil {
		glog.Error("赠送用户存保管箱离线失败usreId :", toUserId)
		return false
	}

	fortune.SafeBox.Savings += gold
	log := BoxLog{"赠送", gold, fortune.SafeBox.Savings, time.Now().Local(), fromUserId, "0"}
	if fortune.SafeBox.BoxLogs == nil {
		fortune.SafeBox.BoxLogs = map[string]BoxLog{}
	}
	logId := strconv.Itoa(len(fortune.SafeBox.BoxLogs) + 1)
	log.LogId = logId
	fortune.SafeBox.BoxLogs[logId] = log
	m.fortune[toUserId] = fortune
	err := SaveFortune(fortune)
	if err != nil {
		glog.Error("赠送用户存保管箱离线失败userId:", toUserId, " f:", fortune)
		return false
	}
	return true
}

//存保险保管箱在线
func (m *UserFortuneManager) UpdateSavingsAdd(fromUserId string, toUserId string, gold int64) bool {
	m.Lock()
	defer m.Unlock()
	fortune := m.fortune[toUserId]

	if fortune == nil {
		glog.Error("赠送用户财富信息失败usreId :", toUserId)
		return false
	}

	fortune.SafeBox.Savings += gold
	log := BoxLog{"赠送", gold, fortune.SafeBox.Savings, time.Now().Local(), fromUserId, "0"}
	if fortune.SafeBox.BoxLogs == nil {
		fortune.SafeBox.BoxLogs = map[string]BoxLog{}
	}
	logId := strconv.Itoa(len(fortune.SafeBox.BoxLogs) + 1)
	log.LogId = logId
	fortune.SafeBox.BoxLogs[logId] = log
	m.fortune[toUserId] = fortune
	err := SaveFortune(fortune)
	if err != nil {
		glog.Error("赠送用户财富信息失败userId:", toUserId, " f:", fortune)
		return false
	}

	return true
}

// add by wangsq end

func (m *UserFortuneManager) ConsumeScore(userId string, score int) bool {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return false
	}

	if fortune.Score < score {
		return false
	}

	fortune.Score -= score

	return true
}

func (m *UserFortuneManager) UpdateUserFortune(userId string) {
	m.Lock()
	defer m.Unlock()

	m.updateUserFortune(userId, 0)
}

func (m *UserFortuneManager) UpdateUserFortune2(userId string, rechargeDiamond int) {
	m.Lock()
	defer m.Unlock()
	glog.Info("UpdateUserFortune2 in,userId=", userId, "|rechargeDiamond=", rechargeDiamond)
	m.updateUserFortune(userId, rechargeDiamond)
}

func (m *UserFortuneManager) updateUserFortune(userId string, rechargeDiamond int) {
	f := m.fortune[userId]
	if f == nil {
		return
	}

	updateGoldMsg := &pb.MsgUpdateGold{}
	updateGoldMsg.Gold = proto.Int64(f.Gold)
	updateGoldMsg.Diamond = proto.Int(f.Diamond)
	updateGoldMsg.GameScore = proto.Int(f.Score)
	updateGoldMsg.VipLevel = proto.Int(f.VipLevel)
	updateGoldMsg.VipStartTime = proto.Int64(f.VipStartTime.Unix())
	updateGoldMsg.VipValidDays = proto.Int(f.VipValidDays)
	updateGoldMsg.Exp = proto.Int(f.Exp)
	/*if f.isDailyGiftBagPay {
		updateGoldMsg.RechargeDiamond = proto.Int(0)
		f.isDailyGiftBagPay = false
	} else {
		updateGoldMsg.RechargeDiamond = proto.Int(rechargeDiamond)
	}*/
	updateGoldMsg.DoubleCardCount = proto.Int(f.DoubleCardCount)
	updateGoldMsg.ForbidCardCount = proto.Int(f.ForbidCardCount)
	updateGoldMsg.ChangeCardCount = proto.Int(f.ChangeCardCount)
	updateGoldMsg.Charm = proto.Int(f.Charm)
	updateGoldMsg.Horn = proto.Int(f.Horn)

	boxmsg := &pb.SafeBoxDef{}
	boxmsg.IsOpen = proto.Bool(f.SafeBox.IsOpen)
	boxmsg.Savings = proto.Int64(f.SafeBox.Savings)
	boxmsg.IsSetPwd = proto.Bool(f.SafeBox.IsSetPwd)
	glog.Info("UpdateUserFortune1 userId =  ", userId, " gold = ", f.Gold)
	for logId, v := range f.SafeBox.BoxLogs {
		logmsg := &pb.BoxLogsDef{}
		logmsg.Reason = proto.String(v.Reason)
		logmsg.Gold = proto.Int64(v.Gold)
		logmsg.Savings = proto.Int64(v.Savings)
		logmsg.Datetime = proto.String(fmt.Sprintf("%d-%d-%d %02d:%02d:%02d", v.DateTime.Year(), v.DateTime.Month(), v.DateTime.Day(), v.DateTime.Hour(), v.DateTime.Minute(), v.DateTime.Second()))
		if v.ToUid != "" {
			logmsg.ToUid = proto.String(v.ToUid)
		}
		logmsg.LogId = proto.String(logId)
		boxmsg.Boxlogs = append(boxmsg.Boxlogs, logmsg)
	}
	updateGoldMsg.SafeBox = boxmsg

	GetPlayerManager().SendServerMsg("", []string{userId}, int32(pb.MessageId_UPDATE_GOLD), updateGoldMsg)
}

func (m *UserFortuneManager) EarnGold(userId string, gold int64, reason string) (int64, bool) {
	ok := m.EarnFortune(userId, gold, 0, 0, false, reason)
	if !ok {
		return 0, false
	}

	f, _ := m.GetUserFortune(userId)
	m.updateUserFortune(userId, 0)
	return f.Gold, true
}

func (m *UserFortuneManager) EarnGoldNoMsg(userId string, gold int64, reason string) (int64, bool) {
	ok := m.EarnFortune(userId, gold, 0, 0, false, reason)
	if !ok {
		return 0, false
	}

	f, _ := m.GetUserFortune(userId)
	return f.Gold, true
}

// add by wangsq start
func (m *UserFortuneManager) EarnCharm(userId string, charm int) (int, bool) {
	f, _ := m.fortune[userId]
	f.Charm += charm
	SaveCharm(userId, f.Charm, false)
	m.fortune[userId] = f

	u := GetUserCache().GetUser(userId)
	rankingList.GetRankingList().UpdateRankingItem(rankingList.RankingType_Charm, u, int64(f.Charm))

	return f.Charm, true
}

func (m *UserFortuneManager) EarnHorn(userId string, horn int) (int, bool) {
	glog.Info("EarnHorn in.", horn)
	f, _ := m.fortune[userId]
	f.Horn += horn
	SaveHorn(userId, f.Horn, false)
	m.fortune[userId] = f

	return f.Horn, true
}

func (m *UserFortuneManager) GetCharmExchangeInfo(userId string, itemId int) int {
	glog.Info("GetCharmExchangeInfo in.userId=", userId, "|itemId=", itemId)
	f, _ := m.fortune[userId]

	if f.CharmExchangeInfo == nil {
		f.CharmExchangeInfo = map[string]int{}
	}
	itemId_str := fmt.Sprintf("%d", itemId)
	_, infok := f.CharmExchangeInfo[itemId_str]
	if !infok {
		f.CharmExchangeInfo[itemId_str] = 0
	}

	m.fortune[userId] = f
	return f.CharmExchangeInfo[itemId_str]
}

func (m *UserFortuneManager) UpdateCharmExchangeInfo(userId string, itemId int, count int) {
	glog.Info("UpdateCharmExchangeInfo in.userId=", userId, "|itemId=", itemId)
	f, _ := m.fortune[userId]

	if f.CharmExchangeInfo == nil {
		f.CharmExchangeInfo = map[string]int{}
	}

	itemId_str := fmt.Sprintf("%d", itemId)
	_, infok := f.CharmExchangeInfo[itemId_str]
	if !infok {
		f.CharmExchangeInfo[itemId_str] = 0
	}

	f.CharmExchangeInfo[itemId_str] += count
	m.fortune[userId] = f
}

// add by wangsq end

func (m *UserFortuneManager) ExchangeGold(userId string, diamond int) (bool, int64, int) {
	glog.V(2).Info("ExchangeGold userId:", userId, " diamond:", diamond)

	if diamond <= 0 {
		return false, 0, 0
	}

	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		glog.V(2).Info("===>userId:", userId, " fortune nil")
		return false, 0, 0
	}

	if fortune.Diamond < diamond {
		glog.V(2).Info("==>fortune diamond:", fortune.Diamond, " diamond:", diamond, " <")
		return false, 0, 0
	}

	fortune.Diamond -= diamond
	spGold := 0
	if diamond == 50 {
		spGold = 30000
	} else if diamond == 100 {
		spGold = 100000
	} else if diamond == 300 {
		spGold = 500000
	} else if diamond == 500 {
		spGold = 1000000
	} else if diamond == 1000 {
		spGold = 3000000
	} else if diamond == 108 {
		spGold = 120000
	} else if diamond == 298 {
		spGold = 520000
	} else if diamond == 518 {
		spGold = 1020000
	} else if diamond == 998 {
		spGold = 3020000
	}
	fortune.Gold += int64(diamond*exchangeGoldRate + spGold)

	if !fortune.IsRobot {
		l := &FortuneLog{}
		l.UserId = userId
		l.Gold = 0
		l.CurGold = fortune.Gold
		l.Diamond = diamond
		l.CurDiamond = fortune.Diamond
		l.Reason = "钻石兑换金币"
		SaveConsumeFortuneLog(l)

		l.Gold = diamond*exchangeGoldRate + spGold
		l.Diamond = 0
		SaveEarnFortuneLog(l)
	}

	m.updateGoldInGame(userId)

	SaveDiamond(userId, fortune.Diamond, fortune.IsRobot)
	SaveGold(userId, fortune.Gold, fortune.IsRobot)

	return true, int64(fortune.Gold), fortune.Diamond
}

func (m *UserFortuneManager) UpdateCurMonWinTimesRankingList(userId string, winTimes int) {
	u := GetUserCache().GetUser(userId)
	if u == nil {
		glog.V(2).Info("===>获取缓存用户信息失败userId:", userId)
		return
	}

	rankingList.GetRankingList().UpdateRankingItem(rankingList.RankingType_WinCurMonth, u, int64(winTimes))
}

func (m *UserFortuneManager) UpdateRechargeRankingList(userId string, diamond int) {
	u := GetUserCache().GetUser(userId)
	if u == nil {
		glog.V(2).Info("===>获取缓存用户信息失败userId:", userId)
		return
	}

	ok, todayRecharge, curWeekRecharge, curMonRecharge := func() (bool, int64, int64, int64) {
		m.Lock()
		defer m.Unlock()

		fortune := m.fortune[userId]
		if fortune == nil {
			return false, 0, 0, 0
		}

		now := time.Now()
		if !util.CompareDate(now, fortune.LastRechargeTime) {
			// 非同一天，清零
			fortune.TodayRecharge = 0
			if now.Weekday() != fortune.LastRechargeTime.Weekday() && now.Weekday() == time.Monday {
				fortune.CurWeekRecharge = 0
			}
			fortune.LastRechargeTime = now

			if now.Day() == 1 {
				fortune.CurMonRec = 0
			}
		}
		fortune.TodayRecharge += diamond
		fortune.CurWeekRecharge += diamond
		fortune.CurMonRec += diamond

		return true, int64(fortune.TodayRecharge), int64(fortune.CurWeekRecharge), int64(fortune.CurMonRec)
	}()

	if !ok {
		return
	}

	if u.GetIsRobot() {
		vipLevel := GetVipLevel(int(todayRecharge))
		if vipLevel != int(u.GetVipLevel()) {
			u.VipLevel = proto.Int(vipLevel)
			u.VipStartTime = proto.Int64(time.Now().Unix())
			u.VipValidDays = proto.Int(30)
		}
	}
	rankingList.GetRankingList().UpdateRankingItem(rankingList.RankingType_RechargeToday, u, todayRecharge)
	rankingList.GetRankingList().UpdateRankingItem(rankingList.RankingType_RechargeCurWeek, u, curWeekRecharge)
	rankingList.GetRankingList().UpdateRankingItem(rankingList.RankingType_RechargeCurMonth, u, curMonRecharge)
}

func (m *UserFortuneManager) UpdateEarningsRankingList(userId string, gold int) {
	u := GetUserCache().GetUser(userId)
	if u == nil {
		glog.V(2).Info("===>获取缓存用户信息失败userId:", userId)
		return
	}

	ok, todayEarnings, curWeekEarnings, curGold, curMonEarnings := func() (bool, int64, int64, int64, int64) {
		m.Lock()
		defer m.Unlock()

		glog.V(2).Info("==>UpdateEarningsRankingList userId:", userId, " gold:", gold)

		fortune := m.fortune[userId]
		if fortune == nil {
			glog.V(2).Info("===>更新金币排行失败，玩家财富信息不存在userId:", userId)
			return false, 0, 0, 0, 0
		}

		now := time.Now()
		if !util.CompareDate(now, fortune.LastEarnTime) {
			// 非同一天，清零
			fortune.TodayEarnings = 0
			if now.Weekday() != fortune.LastEarnTime.Weekday() && now.Weekday() == time.Monday {
				fortune.CurWeekEarnings = 0
			}
			fortune.LastEarnTime = now

			if now.Day() == 1 {
				fortune.CurMonEarn = 0
			}
		}
		fortune.TodayEarnings += gold
		fortune.CurWeekEarnings += gold
		fortune.CurMonEarn += gold

		return true, int64(fortune.TodayEarnings), int64(fortune.CurWeekEarnings), fortune.Gold, int64(fortune.CurMonEarn)
	}()

	if !ok {
		return
	}

	isActiveIn := domainActive.GetUserIosActiveManager().IsActiveContinue()
	//glog.Info("==>UpdateEarningsRankingList IsActiveContinue isActiveIn:", isActiveIn)
	if isActiveIn {
		activeGold := domainActive.GetUserIosActiveManager().AddGold(userId, int64(gold))
		rankingList.GetRankingList().UpdateRankingItem(rankingList.RankingType_Competition, u, int64(activeGold))
		//glog.Info("==>UpdateEarningsRankingList IsActiveContinue activeGold:", activeGold)
	}

	rankingList.GetRankingList().UpdateRankingItem(rankingList.RankingType_EarningsToday, u, todayEarnings)
	rankingList.GetRankingList().UpdateRankingItem(rankingList.RankingType_EarningsCurWeek, u, curWeekEarnings)
	rankingList.GetRankingList().UpdateRankingItem(rankingList.RankingType_Gold, u, curGold)
	rankingList.GetRankingList().UpdateRankingItem(rankingList.RankingType_EarningsCurMonth, u, curMonEarnings)

}

func (m *UserFortuneManager) UpdateVipLevel(userId string, level int) {
	m.setVipLevel(userId, level)
	m.UpdateUserFortune(userId)
}

func (m *UserFortuneManager) CheckVipLevel(userId string) {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return
	}

	if time.Since(fortune.VipStartTime).Hours() > 30*24 {
		fortune.VipLevel = 0
		m.fortune[userId] = fortune
	}
	return
}

func (m *UserFortuneManager) setVipLevel(userId string, level int) {
	vipLevel := level
	if vipLevel <= 0 {
		return
	}

	glog.V(2).Info("===>更新玩家VIP等级userId:", userId, " vipLevel:", vipLevel)

	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return
	}

	if time.Since(fortune.VipStartTime).Hours() > 30*24 {
		// 过期
		fortune.VipStartTime = time.Now()
		fortune.VipValidDays = 30
	} else {
		if vipLevel < fortune.VipLevel {
			return
		}
		// 没有过期
		if fortune.VipLevel == vipLevel {
			fortune.VipValidDays += 30
		} else {
			fortune.VipStartTime = time.Now()
			fortune.VipValidDays = 30
		}
	}
	fortune.VipLevel = vipLevel
}

func (m *UserFortuneManager) ResetDailyGiftBag(userId string) {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return
	}

	now := time.Now()
	if !util.CompareDate(fortune.DailyGiftBagUseTime, now) {
		fortune.BuyDailyGiftBag = false
		fortune.DailyGiftBagUseTime = now
	}
}

func (m *UserFortuneManager) SetBuyDailyGiftBag(userId string) bool {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return false
	}

	if util.CompareDate(fortune.DailyGiftBagUseTime, time.Now()) {
		return false
	}

	fortune.BuyDailyGiftBag = true

	return true
}

func (m *UserFortuneManager) ManagerChangeGameTypeNotify(userId string, date int) bool {
	m.Lock()
	defer m.Unlock()

	fortune := m.fortune[userId]
	if fortune == nil {
		return false
	}

	fortune.ChangeGameTypeNotifyDay = date

	return true
}
