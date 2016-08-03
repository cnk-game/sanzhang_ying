package rankingList

import (
	"code.google.com/p/goprotobuf/proto"
	domainActive "game/domain/iosActive"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2/bson"
	"pb"
	"strconv"
	"sync"
	"time"
	"util"
)

type RankingType int

const (
	RankingType_RechargeToday     RankingType = 1  // 今日充值排行
	RankingType_RechargeYesterday RankingType = 2  // 昨日充值排行
	RankingType_RechargeCurWeek   RankingType = 3  // 本周充值排行
	RankingType_RechargeLastWeek  RankingType = 4  // 上周充值排行
	RankingType_EarningsToday     RankingType = 5  // 今日收入排行
	RankingType_EarningsYesterday RankingType = 6  // 昨日收入排行
	RankingType_EarningsCurWeek   RankingType = 7  // 本周收入排行
	RankingType_EarningsLastWeek  RankingType = 8  // 上周收入排行
	RankingType_Gold              RankingType = 9  // 财富排行
	RankingType_Charm             RankingType = 10 // 昨日魅力值
	RankingType_RechargeCurMonth  RankingType = 11 //本月充值
	RankingType_EarningsCurMonth  RankingType = 12 //本月盈利
	RankingType_WinCurMonth       RankingType = 13 //本月胜场
	RankingType_Competition       RankingType = 14 //奖金竞赛
)

type RankingList struct {
	sync.RWMutex
	items               map[RankingType]*RankingItem
	BroadcastClientFunc func(msgId int32, body proto.Message)
	SendToUser          func(srcId string, dstIds []string, msgId int32, body proto.Message)
}

var rankingList *RankingList

func init() {
	rankingList = &RankingList{}
	rankingList.items = make(map[RankingType]*RankingItem)
}

func GetRankingList() *RankingList {
	return rankingList
}

func (l *RankingList) Init() {
	l.Lock()
	defer l.Unlock()

	rankingList.items[RankingType_RechargeToday] = NewRankingItem(int(RankingType_RechargeToday), 30)
	rankingList.items[RankingType_RechargeYesterday] = NewRankingItem(int(RankingType_RechargeYesterday), 30)
	rankingList.items[RankingType_RechargeCurWeek] = NewRankingItem(int(RankingType_RechargeCurWeek), 30)
	rankingList.items[RankingType_RechargeLastWeek] = NewRankingItem(int(RankingType_RechargeLastWeek), 30)
	rankingList.items[RankingType_EarningsToday] = NewRankingItem(int(RankingType_EarningsToday), 30)
	rankingList.items[RankingType_EarningsYesterday] = NewRankingItem(int(RankingType_EarningsYesterday), 30)
	rankingList.items[RankingType_EarningsCurWeek] = NewRankingItem(int(RankingType_EarningsCurWeek), 30)
	rankingList.items[RankingType_EarningsLastWeek] = NewRankingItem(int(RankingType_EarningsLastWeek), 30)
	rankingList.items[RankingType_Gold] = NewRankingItem(int(RankingType_Gold), 30)
	rankingList.items[RankingType_Charm] = NewRankingItem(int(RankingType_Charm), 30)

	rankingList.items[RankingType_RechargeCurMonth] = NewRankingItem(int(RankingType_RechargeCurMonth), 30)
	rankingList.items[RankingType_EarningsCurMonth] = NewRankingItem(int(RankingType_EarningsCurMonth), 30)
	rankingList.items[RankingType_WinCurMonth] = NewRankingItem(int(RankingType_WinCurMonth), 30)
	rankingList.items[RankingType_Competition] = NewRankingItem(int(RankingType_Competition), 30)
}

func (l *RankingList) Save() {
	l.Lock()
	defer l.Unlock()

	glog.V(2).Info("====保存排行榜")
	for _, item := range l.items {
		item.SaveRankingItem()
	}
}

func (l *RankingList) UpdateRankingItem(rankingType RankingType, user *pb.UserDef, v int64) {
	glog.V(2).Info("====>更新排行榜rankingType:", rankingType, " userId:", user.GetUserId(), " v:", v)
	if l.items[rankingType] == nil {
		return
	}
	//	oldOrder1 := l.items[rankingType].GetOrder1UserNickname()
	l.resetRankingList()
	l.updateRankingItem(rankingType, user, v)
	//	order1 := l.items[rankingType].GetOrder1UserNickname()

	//	if oldOrder1 != order1 {
	//		// 第一名有变化
	//		switch rankingType {
	//		case RankingType_RechargeToday:
	//			if l.BroadcastClientFunc != nil {
	//				l.BroadcastClientFunc(int32(pb.MessageId_CHAT), util.BuildSysBugle(fmt.Sprintf("恭喜%v成为了今日充值排行榜第一名！", order1)))
	//			}
	//		case RankingType_EarningsToday:
	//			if l.BroadcastClientFunc != nil {
	//				l.BroadcastClientFunc(int32(pb.MessageId_CHAT), util.BuildSysBugle(fmt.Sprintf("恭喜%v成为了今日盈利排行榜第一名！", order1)))
	//			}
	//		case RankingType_Gold:
	//			if l.BroadcastClientFunc != nil {
	//				l.BroadcastClientFunc(int32(pb.MessageId_CHAT), util.BuildSysBugle(fmt.Sprintf("恭喜%v成为了总财富排行榜第一名！", order1)))
	//			}
	//		}
	//	}
}

func (l *RankingList) updateRankingItem(rankingType RankingType, user *pb.UserDef, v int64) {
	l.RLock()
	defer l.RUnlock()

	item := l.items[rankingType]
	if item == nil {
		glog.Error("排行类型不存在rankingType:", rankingType)
		return
	}

	item.UpdateRankingItem(user, v)
}

func (l *RankingList) resetRankingList() {
	l.Lock()
	defer l.Unlock()

	now := time.Now()

	rechargeToday := l.items[RankingType_RechargeToday]
	if !util.CompareDate(now, rechargeToday.ResetTime) {
		// 重置,复制今日记录到昨日记录
		glog.V(2).Info("****>重置今日充值记录rechargeToday.resetTime:", rechargeToday.ResetTime, " items:", rechargeToday.Items)
		rechargeYesterday := l.items[RankingType_RechargeYesterday]
		rechargeYesterday.Items = rechargeToday.Items

		rechargeToday.ResetRankingItem()
		glog.V(2).Info("====>rechargeToday.resetTime:", rechargeToday.ResetTime)

		// 向昨日充值排行榜前3名发放奖励
		/*if len(rechargeYesterday.Items) >= 1 {
			l.sendYesterdayRechargePrize(rechargeYesterday.Items[0].UserId, int(rechargeYesterday.Items[0].Value), "恭喜您成为昨日充值排行榜第1！")
		}

		if len(rechargeYesterday.Items) >= 2 {
			l.sendYesterdayRechargePrize(rechargeYesterday.Items[1].UserId, int(float64(rechargeYesterday.Items[1].Value)*0.5), "恭喜您成为昨日充值排行榜第2！")
		}

		if len(rechargeYesterday.Items) >= 3 {
			l.sendYesterdayRechargePrize(rechargeYesterday.Items[2].UserId, int(float64(rechargeYesterday.Items[2].Value)*0.2), "恭喜您成为昨日充值排行榜第3！")
		}*/
	}
	//重置本周充值
	rechargeCurWeek := l.items[RankingType_RechargeCurWeek]
	if !util.CompareDate(now, rechargeCurWeek.ResetTime) && now.Weekday() == time.Monday {
		rechargeLastWeek := l.items[RankingType_RechargeLastWeek]
		rechargeLastWeek.Items = rechargeCurWeek.Items

		rechargeCurWeek.ResetRankingItem()
		cur_C := strconv.Itoa(int(now.Year())) + strconv.Itoa(int(now.Month())) + strconv.Itoa(int(now.Day())) + "_" + "RechargeCurWeek"
		rechargeCurWeek.SaveAsRankingItem(cur_C)
	}

	earningsToday := l.items[RankingType_EarningsToday]
	if !util.CompareDate(now, earningsToday.ResetTime) {
		earningsYesterday := l.items[RankingType_EarningsYesterday]
		earningsYesterday.Items = earningsToday.Items

		earningsToday.ResetRankingItem()

		glog.V(2).Info("****>重置今日盈利记录earningsToday.resetTime:", earningsToday.ResetTime, " today.items:", earningsToday.Items, " yesterday:", earningsYesterday.Items)
	}

	//充值本周赢利
	earningsCurWeek := l.items[RankingType_EarningsCurWeek]
	if !util.CompareDate(now, earningsCurWeek.ResetTime) && now.Weekday() == time.Monday {
		glog.V(2).Info("=====>本周盈利重置")
		earningsLastWeek := l.items[RankingType_EarningsLastWeek]
		earningsLastWeek.Items = earningsCurWeek.Items
		cur_C := strconv.Itoa(int(now.Year())) + strconv.Itoa(int(now.Month())) + strconv.Itoa(int(now.Day())) + "_" + "EarningsCurWeek"
		earningsCurWeek.SaveAsRankingItem(cur_C)

		earningsCurWeek.ResetRankingItem()
	}

	//重置本月充值
	rechargeCurMonth := l.items[RankingType_RechargeCurMonth]
	if !util.CompareDate(now, rechargeCurMonth.ResetTime) && now.Day() == 1 {
		glog.V(2).Info("=====>本月充值重置")

		cur_C := strconv.Itoa(int(now.Year())) + strconv.Itoa(int(now.Month())) + strconv.Itoa(int(now.Day())) + "_" + "RechargeCurMonth"
		rechargeCurWeek.SaveAsRankingItem(cur_C)
		rechargeCurMonth.ResetRankingItem()
	}
	//重置本月盈利
	earningsCurMonth := l.items[RankingType_EarningsCurMonth]
	if !util.CompareDate(now, earningsCurMonth.ResetTime) && now.Day() == 1 {
		glog.V(2).Info("=====>本月充值重置")

		cur_C := strconv.Itoa(int(now.Year())) + strconv.Itoa(int(now.Month())) + strconv.Itoa(int(now.Day())) + "_" + "EarningsCurMonth"
		earningsCurMonth.SaveAsRankingItem(cur_C)
		earningsCurMonth.ResetRankingItem()
	}
	//重置本月胜场
	winsCurMonth := l.items[RankingType_WinCurMonth]
	if !util.CompareDate(now, winsCurMonth.ResetTime) && now.Day() == 1 {
		glog.V(2).Info("=====>本月充值重置")

		cur_C := strconv.Itoa(int(now.Year())) + strconv.Itoa(int(now.Month())) + strconv.Itoa(int(now.Day())) + "_" + "WinCurMonth"
		winsCurMonth.SaveAsRankingItem(cur_C)
		winsCurMonth.ResetRankingItem()
	}
	//重置活动比赛
	competition := l.items[RankingType_Competition]
	isInActive := domainActive.GetUserIosActiveManager().IsActiveContinue()
	if isInActive == false && !util.CompareDate(now, competition.ResetTime) {
		_, _, endTime := domainActive.GetUserIosActiveManager().GetIosActiveContent()
		dateEnd := util.ParseTime(endTime).Day() + 1
		if now.Day() == dateEnd {
			glog.V(2).Info("=====>重置活动比赛")

			cur_C := strconv.Itoa(int(now.Year())) + strconv.Itoa(int(now.Month())) + strconv.Itoa(int(now.Day())) + "_" + "COMPETITION"
			competition.SaveAsRankingItem(cur_C)
			competition.ResetRankingItem()
		}
	}
}

func (l *RankingList) BuildMessage(rankingType pb.RankingType) *pb.RankingList {
	l.Lock()
	defer l.Unlock()

	switch rankingType {
	case pb.RankingType_RECHARGE_TODAY:
		return l.items[RankingType_RechargeToday].BuildMessage(rankingType)
	case pb.RankingType_RECHARGE_YESTERDAY:
		return l.items[RankingType_RechargeYesterday].BuildMessage(rankingType)
	case pb.RankingType_RECHARGE_LAST_WEEK:
		return l.items[RankingType_RechargeLastWeek].BuildMessage(rankingType)
	case pb.RankingType_EARNINGS_TODAY:
		return l.items[RankingType_EarningsToday].BuildMessage(rankingType)
	case pb.RankingType_EARNINGS_YESTERDAY:
		return l.items[RankingType_EarningsYesterday].BuildMessage(rankingType)
	case pb.RankingType_EARNINGS_LAST_WEEK:
		return l.items[RankingType_EarningsLastWeek].BuildMessage(rankingType)
	case pb.RankingType_GOLD:
		return l.items[RankingType_Gold].BuildMessage(rankingType)
	case pb.RankingType_CHARM:
		return l.items[RankingType_Charm].BuildMessage(rankingType)
	case pb.RankingType_RECHARGE_THIS_MONTH:
		return l.items[RankingType_RechargeCurMonth].BuildMessage(rankingType)
	case pb.RankingType_EARNINGS_THIS_MONTH:
		return l.items[RankingType_EarningsCurMonth].BuildMessage(rankingType)
	case pb.RankingType_RECHARGE_THIS_WEEK:
		return l.items[RankingType_RechargeCurWeek].BuildMessage(rankingType)
	case pb.RankingType_EARNINGS_THIS_WEEK:
		return l.items[RankingType_EarningsCurWeek].BuildMessage(rankingType)
	case pb.RankingType_WINS_THIS_MONTH:
		return l.items[RankingType_WinCurMonth].BuildMessage(rankingType)
	case pb.RankingType_COMPETITION:
		glog.Info("=====>RankingType_COMPETITION")
		return l.items[RankingType_Competition].BuildMessage(rankingType)
	}
	return nil
}

func (l *RankingList) sendYesterdayRechargePrize(userId string, prizeDiamond int, content string) {
	prizeMail := &pb.PrizeMailDef{}
	prizeMail.MailId = proto.String(bson.NewObjectId().Hex())
	prizeMail.Content = proto.String(content)
	prizeMail.Prize = &pb.PrizeDef{}
	prizeMail.Prize.Diamond = proto.Int(prizeDiamond)

	if l.SendToUser != nil {
		l.SendToUser("", []string{userId}, int32(pb.ServerMsgId_MQ_PRIZE_MAIL), prizeMail)
	}

	glog.V(2).Info("====>发放昨日充值排行榜前3名奖励prizeMail:", prizeMail)
}
