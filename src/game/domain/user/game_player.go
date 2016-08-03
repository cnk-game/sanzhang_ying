package user

import (
	"code.google.com/p/goprotobuf/proto"
	domainCdKey "game/domain/cdkey"
	domainUserIosAcive "game/domain/iosActive"
	newUserTask "game/domain/newusertask"
	domainPrize "game/domain/prize"
	domainSlots "game/domain/slots"
	"github.com/golang/glog"
	"time"
	"util"
)

type GamePlayer struct {
	User                   *User
	NewPlayer              bool
	SendToClientFunc       func(msgId int32, body proto.Message)
	OnlinePrizeGainRecords *domainPrize.OnlinePrizeGainRecords
	UserTasks              *domainPrize.UserTasks
	MatchRecord            *MatchRecord
	LastGameId             int
	LoginTime              time.Time
	UserLog                *UserLog
	SlotMachine            *domainSlots.SlotMachine
	PrizeMails             *domainPrize.PrizeMails
	LastChatTime           time.Time
	CDKeyGainRecord        *domainCdKey.CDKeyGainRecord
	SignInRecord           *domainPrize.SignInRecord
	LoginIP                string
	LoginDeviceId          string
	SessKey                string
	LastExp                int
	OnLogoutFunc           func(userId string)
}

func NewPlayer() *GamePlayer {
	p := &GamePlayer{}

	return p
}

func GetPlayer(p interface{}) *GamePlayer {
	if p == nil {
		return nil
	}
	switch player := p.(type) {
	case *GamePlayer:
		return player
	}
	return nil
}

func (p *GamePlayer) SendToClient(msgId int32, body proto.Message) {
	if p.SendToClientFunc != nil {
		p.SendToClientFunc(msgId, body)
	}
}

func (p *GamePlayer) OnLogin() bool {
	glog.V(2).Info("===>GamePlayer OnLogin")
	p.LoginTime = time.Now()

	if !GetUserFortuneManager().LoadUserFortune(p.User.UserId, p.User.IsRobot) {
		return false
	}

	p.OnlinePrizeGainRecords = domainPrize.NewOnlinePrizeGainRecords(p.User.UserId)
	p.UserTasks = domainPrize.NewUserTasks(p.User.UserId)
	p.MatchRecord = FindMatchRecord(p.User.UserId)
	if util.CompareDate(p.User.CreateTime, time.Now()) && p.UserLog == nil {
		p.UserLog, _ = FindUserLog(p.User.UserId)
	}
	p.SlotMachine = &domainSlots.SlotMachine{}
	p.PrizeMails = domainPrize.NewPrizeMails(p.User.UserId)
	p.ResetLuckyValue()
	p.CDKeyGainRecord, _ = domainCdKey.FindCdKeyGainRecord(p.User.UserId)
	p.SignInRecord, _ = domainPrize.FindSignInRecord(p.User.UserId)
	p.SignInRecord.ResetSignRecord()

	return true
}

func (p *GamePlayer) OnLogout() {
	defer func() {
		p.SendToClientFunc = nil
		GetPlayerManager().DelItem(p.User.UserId, p.User.IsRobot)
		GetBackgroundUserManager().DelUser(p.User.UserId)
		if p.OnLogoutFunc != nil {
			p.OnLogoutFunc(p.User.UserId)
			p.OnLogoutFunc = nil
		}
		if p.User.IsRobot {
			GetFakeRankingList().RemoveRobot(p.User.UserId)
		}
		domainUserIosAcive.GetUserIosActiveManager().SaveStatus(p.User.UserId)
		newUserTask.GetNewUserTaskManager().SaveUserTask(p.User.UserId)
	}()

	GetUserFortuneManager().UnloadUserFortune(p.User.UserId)

	if p.User.IsRobot {
		return
	}

	glog.Info("===>GamePlayer OnLogout userId:", p.User.UserId, " sessKey:", p.SessKey, " LoginIP:", p.LoginIP, " loginTime:", p.LoginTime)

	err := SaveUser(p.User)
	if err != nil {
		glog.Error("保存用户数据失败err:", err, " u:", p.User)
	}

	p.OnlinePrizeGainRecords.Save()
	p.UserTasks.SaveTasks()
	SaveMatchRecord(p.MatchRecord)
	if p.UserLog != nil && !p.User.IsRobot {
		p.UserLog.TotalOnlineSeconds += int(time.Since(p.LoginTime).Seconds())
		SaveUserLog(p.UserLog)
	}
	p.PrizeMails.SaveMails()

	p.saveLoginRecord()
	domainCdKey.SaveCdKeyGainRecord(p.CDKeyGainRecord)
	err = domainPrize.SaveSignInRecord(p.SignInRecord)
	if err != nil {
		glog.Error("保存签到记录失败err:", err, " r:", p.SignInRecord)
	}
}

func (p *GamePlayer) SetUserCache() {
	GetUserCache().SetUser(p.User.BuildMessage(p.MatchRecord.BuildMessage()))
}

func (p *GamePlayer) saveLoginRecord() {
	if p.User.IsRobot {
		return
	}
	loginRecord := &LoginRecord{}
	loginRecord.UserId = p.User.UserId
	loginRecord.UserName = p.User.UserName
	loginRecord.Channel = p.User.ChannelId
	loginRecord.LoginTime = p.LoginTime
	loginRecord.LogoutTime = time.Now()
	loginRecord.LoginIP = p.LoginIP
	loginRecord.DeviceId = p.LoginDeviceId
	InsertLoginRecord(loginRecord)
}

func (p *GamePlayer) ResetLuckyValue() {
	now := time.Now()
	if util.CompareDate(p.User.LuckyValueResetTime, now) {
		return
	}
	p.User.LuckyValueResetTime = now
	p.User.LuckyValue = 200
}
