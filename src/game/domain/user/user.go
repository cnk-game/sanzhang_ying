package user

import (
	"code.google.com/p/goprotobuf/proto"
	"config"
	"fmt"
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"pb"
	"strconv"
	"time"
	"util"
)

type User struct {
	UserId                      string    `bson:"userId"`
	UserName                    string    `bson:"userName"`
	Password                    string    `bson:"password"`
	Nickname                    string    `bson:"nickname"`
	Gender                      int       `bson:"gender"`
	Signiture                   string    `bson:"signiture"`
	PhotoUrl                    string    `bson:"photoUrl"`
	IsBind                      bool      `bson:"isBind"`
	TodayRecharge               int       `bson:"todayRecharge"`
	CurWeekRecharge             int       `bson:"curWeekRecharge"`
	TodayEarnings               int       `bson:"todayEarnings"`
	CurWeekEarnings             int       `bson:"curWeekEarnings"`
	LastEarnTime                int       `bson:"lastEarnTime"`
	LastGainVipPrizeTime        time.Time `bson:"lastGainVipPrizeTime"`
	TotalOnlineSeconds          int       `bson:"totalOnlineSeconds"`
	CreateTime                  time.Time `bson:"createTime"`
	ChannelId                   string    `bson:"channelId"`
	ShippingAddressName         string    `bson:"shippingAddressName"`
	ShippingAddressPhone        string    `bson:"shippingAddressPhone"`
	ShippingAddressAddress      string    `bson:"shippingAddressAddress"`
	ShippingAddressZipCode      string    `bson:"shippingAddressZipCode"`
	SubsidyPrizeTimes           int       `bson:"subsidyPrizeTimes"`
	LastSubsidyPrizeTime        time.Time `bson:"lastSubsidyPrizeTime"`
	IsBindShippingAddress       bool      `bson:"isBindShippingAddress"`
	YesterdayRechargeOrder1     bool      `bson:"yesterdayRechargeOrder1"`
	YesterdayRechargeOrder1Time time.Time `bson:"yesterdayRechargeOrder1Time"`
	YesterdayEarningOrder1      bool      `bson:"yesterdayEarningOrder1"`
	YesterdayEarningOrder1Time  time.Time `bson:"yesterdayEarningOrder1Time"`
	LastWeekRechargeOrder1      bool      `bson:"lastWeekRechargeOrder1"`
	LastWeekRechargeOrder1Time  time.Time `bson:"lastWeekRechargeOrder1Time"`
	RewardInGameTimes           int       `bson:"rewardInGameTimes"`
	LastRewardInGameTime        time.Time `bson:"lastRewardInGameTime"`
	IsRobot                     bool      `bson:"isRobot"`
	RobotVipLevel               int       `bson:"robotVipLevel"`
	IsChangedNickname           bool      `bson:"isChangedNickname"`
	IsChangedPhotoUrl           bool      `bson:"isChangedPhotoUrl"`
	Model                       string    `bson:"model"`
	LuckyValue                  int       `bson:"luckyValue"`
	LuckyValueResetTime         time.Time `bson:"luckyValueResetTime"`
	IsLocked                    bool      `bson:"isLocked"`
	UpgradePrizeVersion         string    `bson:"upgradePrizeVersion"`
	ShippingAddressQQ           string    `bson:"shippingAddressQQ"`
}

func (u *User) BuildMessage(matchRecord *pb.UserMatchRecordDef) *pb.UserDef {
	msg := &pb.UserDef{}
	msg.UserId = proto.String(u.UserId)
	msg.NickName = proto.String(u.Nickname)
	if u.Gender == 1 {
		msg.Gender = pb.Gender_BOY.Enum()
	} else {
		msg.Gender = pb.Gender_GIRL.Enum()
	}
	msg.Signiture = proto.String(u.Signiture)
	msg.PhotoUrl = proto.String(u.PhotoUrl)

	f, _ := GetUserFortuneManager().GetUserFortune(u.UserId)

	msg.Gold = proto.Int64(f.Gold)
	msg.Diamond = proto.Int(f.Diamond)
	msg.GameScore = proto.Int(f.Score)
	msg.IsBind = proto.Bool(u.IsBind)
	msg.Charm = proto.Int(f.Charm)

	if u.IsRobot {
		msg.VipLevel = proto.Int(u.RobotVipLevel)
		msg.VipStartTime = proto.Int64(time.Now().Unix())
		msg.VipValidDays = proto.Int(30)
	} else {
		msg.VipLevel = proto.Int(f.VipLevel)
		msg.VipStartTime = proto.Int64(f.VipStartTime.Unix())
		msg.VipValidDays = proto.Int(f.VipValidDays)
	}

	msg.LastGainVipPrizeTime = proto.Int64(u.LastGainVipPrizeTime.Unix())

	msg.MatchRecord = matchRecord
	if u.IsRobot {
		msg.MatchRecord.PlayTotalCount = proto.Int32(msg.MatchRecord.GetPlayTotalCount() / 5)
		msg.MatchRecord.PlayWinCount = proto.Int32(msg.MatchRecord.GetPlayWinCount() / 5)
	}

	msg.Exp = proto.Int(f.Exp)
	msg.LuckyValue = proto.Int(u.LuckyValue)
	msg.IsBindAddress = proto.Bool(u.IsBindShippingAddress)
	msg.IsRobot = proto.Bool(u.IsRobot)

	msg.DailyGiftBagUseTime = proto.Int64(f.DailyGiftBagUseTime.Unix())
	msg.SubsidyPrizeTimes = proto.Int(3 - u.SubsidyPrizeTimes)
	msg.LastPayTime = proto.Int64(f.LastPayTime.Unix())
	msg.GainedFirstRechargeBonus = proto.Bool(f.GainedFirstRechargeBonus)
	msg.Horn = proto.Int(f.Horn)

	boxmsg := &pb.SafeBoxDef{}
	boxmsg.IsOpen = proto.Bool(f.SafeBox.IsOpen)
	boxmsg.Savings = proto.Int64(f.SafeBox.Savings)
	boxmsg.IsSetPwd = proto.Bool(f.SafeBox.IsSetPwd)
	boxLen := len(f.SafeBox.BoxLogs)
	boxIndex := 0
	if boxLen > 20 {
		boxIndex = boxLen - 20
	}
	for logId, v := range f.SafeBox.BoxLogs {
		ii, _ := strconv.Atoi(logId)
		if ii > boxIndex {
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
	}
	msg.SafeBox = boxmsg

	vip_configs := config.GetVipPriceConfigManager().GetVipConfig()
	for _, value := range f.VipTaskStates {
		submsg := &pb.UserVipTaskDef{}
		submsg.Level = proto.Int(value.VipTaskId)
		if value.StartTime != 0 {
			edTime := value.StartTime + int64(vip_configs[value.VipTaskId].PrizeDays*86400)
			if edTime < util.GetDayZero() {
				submsg.State = proto.Int(0)
				submsg.StartTime = proto.Int64(0)
				submsg.EndTime = proto.Int64(0)
			} else if value.LastGainTime != util.GetDayZero() {
				submsg.State = proto.Int(2)
				submsg.StartTime = proto.Int64(value.StartTime)
				submsg.EndTime = proto.Int64(value.StartTime + int64(vip_configs[value.VipTaskId].PrizeDays*86400))
			} else {
				submsg.State = proto.Int(1)
				submsg.StartTime = proto.Int64(value.StartTime)
				submsg.EndTime = proto.Int64(value.StartTime + int64(vip_configs[value.VipTaskId].PrizeDays*86400))
			}
		} else {
			submsg.State = proto.Int(0)
			submsg.StartTime = proto.Int64(0)
			submsg.EndTime = proto.Int64(0)
		}
		submsg.PrizeGold = proto.Int(vip_configs[value.VipTaskId].PrizeGold)
		msg.VipTaskList = append(msg.VipTaskList, submsg)
	}

	msg.Channel = proto.String(u.ChannelId)

	return msg
}

const (
	userNameIdC  = "user_name_id"
	userC        = "user"
	iosuserC     = "ios_user"
	phoneUserC   = "phone_user"
	SafeBoxUserC = "safebox_user"
	pushUserC    = "push_user"
	caoHuaUserC  = "caohua_user"
	leshiUserC   = "leshi_user"
	xunleiUserC  = "xunlei_user"
	hmIosUserC   = "hmios_user"
	JinLiUserC   = "jinli_user"
	KupaiUserC   = "kupai_user"
	MeizuUserC   = "Meizu_user"
	Ios51UserC   = "ios51_user"
)

type UserNameId struct {
	UserName string `bson:"userName"`
	UserId   string `bson:"userId"`
}

func FindUserNameIdByUserName(userName string) (*UserNameId, error) {
	u := &UserNameId{}
	err := util.WithUserCollection(userNameIdC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"userName": userName}).One(u)
	})
	return u, err
}

func FindUserNameIdByUserId(userId string) (*UserNameId, error) {
	u := &UserNameId{}
	err := util.WithUserCollection(userNameIdC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"userId": userId}).One(u)
	})
	return u, err
}

func SaveUserNameIdByUserName(userId string, userName string) error {
	u := &UserNameId{}
	u.UserId = userId
	u.UserName = userName
	return util.WithSafeUserCollection(userNameIdC, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"userName": u.UserName}, u)
		return err
	})
}

func SaveUserNameIdByUserId(userId string, userName string) error {
	u := &UserNameId{}
	u.UserId = userId
	u.UserName = userName
	return util.WithSafeUserCollection(userNameIdC, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"userId": userId}, u)
		return err
	})
}

func FindByUserId(userId string) (*User, error) {
	u := &User{}
	userTableC, er := GetUserIdTable(userId)
	if er != nil {
		return u, er
	}

	err := util.WithUserCollection(userTableC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"userId": userId}).One(u)
	})
	return u, err
}

func FindByUserName(username string) (*User, error) {
	uNameId := &UserNameId{}
	u := &User{}
	uNameId, err1 := FindUserNameIdByUserName(username)
	if err1 != nil {
		return u, err1
	}

	glog.Info("++FindByUserName = ", uNameId.UserId)

	userId := uNameId.UserId
	userTableC, err2 := GetUserIdTable(userId)
	if err2 != nil {
		return u, err2
	}

	glog.Info("++FindByUserName = ", userTableC)

	err3 := util.WithUserCollection(userTableC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"userName": username}).One(u)
	})
	return u, err3
}

func SaveUser(u *User) error {
	userId := u.UserId
	userTableC, er := GetUserIdTable(userId)
	if er != nil {
		return er
	}

	return util.WithSafeUserCollection(userTableC, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"userId": u.UserId}, u)
		return err
	})
}

func FindByNickname(nickname string) ([]*User, error) {
	users := []*User{}
	err := util.WithUserCollection(userC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"nickname": nickname}).Limit(10).All(&users)
	})
	return users, err
}

func SetLocked(userId string, locked bool) error {
	userTableC, er := GetUserIdTable(userId)
	if er != nil {
		return er
	}
	return util.WithSafeUserCollection(userTableC, func(c *mgo.Collection) error {
		return c.Update(bson.M{"userId": userId}, bson.M{"$set": bson.M{"isLocked": locked}})
	})
}

func (u *User) ResetSubsidyPrizeTime() {
	now := time.Now()
	if !util.CompareDate(now, u.LastSubsidyPrizeTime) {
		u.SubsidyPrizeTimes = 0
		u.LastSubsidyPrizeTime = now
	}
}

type IOSUser struct {
	Uid      string `bson:"uid"`
	UserName string `bson:"userName"`
	Password string `bson:"password"`
}

func FindIOSUser(uid string) (string, string, error) {
	user := IOSUser{"", "", ""}
	err := util.WithUserCollection(iosuserC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"uid": uid}).One(&user)
	})
	return user.UserName, user.Password, err
}

func InsertIOSUser(uid, userName, pwd string) error {
	user := IOSUser{uid, userName, pwd}
	err := util.WithUserCollection(iosuserC, func(c *mgo.Collection) error {
		return c.Insert(&user)
	})
	return err
}

func RemoveIOSUser(uid string) {
	go util.WithUserCollection(iosuserC, func(c *mgo.Collection) error {
		return c.Remove(bson.M{"uid": uid})
	})
}

func FindCaoHuaUser(uid string) (string, string, error) {
	user := IOSUser{"", "", ""}
	err := util.WithUserCollection(caoHuaUserC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"uid": uid}).One(&user)
	})
	return user.UserName, user.Password, err
}

func InsertChaoHuaUser(uid, userName, pwd string) error {
	user := IOSUser{uid, userName, pwd}
	err := util.WithUserCollection(caoHuaUserC, func(c *mgo.Collection) error {
		return c.Insert(&user)
	})
	return err
}

func FindJinLiUser(uid string) (string, string, error) {
	user := IOSUser{"", "", ""}
	err := util.WithUserCollection(JinLiUserC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"uid": uid}).One(&user)
	})
	return user.UserName, user.Password, err
}

func InsertJinLiUser(uid, userName, pwd string) error {
	user := IOSUser{uid, userName, pwd}
	err := util.WithUserCollection(JinLiUserC, func(c *mgo.Collection) error {
		return c.Insert(&user)
	})
	return err
}

func FindKupaiUser(uid string) (string, string, error) {
	glog.Info("Find kupaiuserid=", uid)
	user := IOSUser{"", "", ""}
	err := util.WithUserCollection(KupaiUserC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"uid": uid}).One(&user)
	})
	return user.UserName, user.Password, err
}

func InsertKupaiUser(uid, userName, pwd string) error {
	glog.Info("insert kupaiuserid=", uid)

	user := IOSUser{uid, userName, pwd}
	err := util.WithUserCollection(KupaiUserC, func(c *mgo.Collection) error {
		return c.Insert(&user)
	})
	return err
}

func FindMeizuUser(uid string) (string, string, error) {
	glog.Info("Find Meizuuserid=", uid)
	user := IOSUser{"", "", ""}
	err := util.WithUserCollection(MeizuUserC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"uid": uid}).One(&user)
	})
	return user.UserName, user.Password, err
}

func InsertMeizuUser(uid, userName, pwd string) error {
	glog.Info("insert Meizuuserid=", uid)

	user := IOSUser{uid, userName, pwd}
	err := util.WithUserCollection(MeizuUserC, func(c *mgo.Collection) error {
		return c.Insert(&user)
	})
	return err
}

func FindIos51User(uid string) (string, string, error) {
	glog.Info("Find Ios51userid=", uid)
	user := IOSUser{"", "", ""}
	err := util.WithUserCollection(Ios51UserC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"uid": uid}).One(&user)
	})
	return user.UserName, user.Password, err
}

func InsertIos51User(uid, userName, pwd string) error {
	glog.Info("insert Ios51userid=", uid)

	user := IOSUser{uid, userName, pwd}
	err := util.WithUserCollection(Ios51UserC, func(c *mgo.Collection) error {
		return c.Insert(&user)
	})
	return err
}

func FindLeshiUser(uid string) (string, string, error) {
	user := IOSUser{"", "", ""}
	err := util.WithUserCollection(leshiUserC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"uid": uid}).One(&user)
	})
	return user.UserName, user.Password, err
}

func InsertLeshiUser(uid, userName, pwd string) error {
	user := IOSUser{uid, userName, pwd}
	err := util.WithUserCollection(leshiUserC, func(c *mgo.Collection) error {
		return c.Insert(&user)
	})
	return err
}

func FindXunleiUser(uid string) (string, string, error) {
	user := IOSUser{"", "", ""}
	err := util.WithUserCollection(xunleiUserC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"uid": uid}).One(&user)
	})
	return user.UserName, user.Password, err
}

func InsertXunleiUser(uid, userName, pwd string) error {
	user := IOSUser{uid, userName, pwd}
	err := util.WithUserCollection(xunleiUserC, func(c *mgo.Collection) error {
		return c.Insert(&user)
	})
	return err
}

func FindHMIosUser(uid string) (string, string, error) {
	user := IOSUser{"", "", ""}
	err := util.WithUserCollection(hmIosUserC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"uid": uid}).One(&user)
	})
	return user.UserName, user.Password, err
}

func InsertHMIosUser(uid, userName, pwd string) error {
	user := IOSUser{uid, userName, pwd}
	err := util.WithUserCollection(hmIosUserC, func(c *mgo.Collection) error {
		return c.Insert(&user)
	})
	return err
}

type PhoneUser struct {
	Phone    string    `bson:"phone"`
	UserId   string    `bson:"userId"`
	BindTime time.Time `bson:"bindTime"`
}

func GetPhoneIsBind(phone string) (string, error) {
	userTemp := PhoneUser{"", "", time.Now()}
	err := util.WithUserCollection(phoneUserC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"phone": phone}).One(&userTemp)
	})
	glog.Info("++GetPhoneIsBind = ", userTemp)
	return userTemp.UserId, err
}

func SavePhoneUser(phone string, userId string) error {
	u := PhoneUser{phone, userId, time.Now()}
	err := util.WithUserCollection(phoneUserC, func(c *mgo.Collection) error {
		return c.Insert(&u)
	})
	return err
}

func GetPhoneIsBindSafeBox(phone string) (string, error) {
	userTemp := PhoneUser{"", "", time.Now()}
	err := util.WithUserCollection(SafeBoxUserC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"phone": phone}).One(&userTemp)
	})
	glog.Info("++GetPhoneIsBindSafeBox = ", userTemp)
	return userTemp.UserId, err
}

func SaveSafeBoxUser(phone string, userId string) error {
	u := PhoneUser{phone, userId, time.Now()}
	err := util.WithUserCollection(SafeBoxUserC, func(c *mgo.Collection) error {
		return c.Insert(&u)
	})
	return err
}

type PushUser struct {
	PushToken string    `bson:"pushToken"`
	UserId    string    `bson:"userId"`
	SetTime   time.Time `bson:"setTime"`
}

func SavePushUser(umPushToken string, userId string) error {
	/*u := PushUser{umPushToken, userId, time.Now()}
	err := util.WithUserCollection(pushUserC, func(c *mgo.Collection) error {
		return c.Upsert(&u)
	})
	return err*/

	return util.WithUserCollection(pushUserC, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"userId": userId}, bson.M{"$set": bson.M{"pushToken": umPushToken}})
		return err
	})
}
