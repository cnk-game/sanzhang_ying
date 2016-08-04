package models

import (
    "time"
    "fmt"
)




type ChannelSynthesizeInfo struct {
    Channel string
    NewUserCount int
    VaildNewUserCount int
    LoginUserCount int
    PayCount int
    PayPlayerCount int
    PayTotalAmount int
    AfterXDayRemainPlayerCountMap map[int]int
}

type SynthesizeCalc struct {
    year, month, day int
    // 新增用户
    totalNewUsers []*NewUserLog
    // 有效新增用户
    vaildNewUsers []*NewUserLog
    // 无效用户ID
    invaildNewUserIds map[string]int
    // 支付记录
    payRecords []*PayLog
    // 活跃用户
    loginUserCount int
    // 统计出的渠道列表
    channleIdMap map[string]int
    // 渠道的留存信息
    channelSynthesizeMap map[string]*ChannelSynthesizeInfo
    // X天留存人数
    afterXDayRemainPlayerCountMap map[int]int
    // 付费次数
    PayCount int
    // 付费人数
    PayPlayerCount int
    // 付费金额
    PayTotalAmount int
}

func (r *SynthesizeCalc)Init(year, month, day int) {
    r.year = year
    r.month = month
    r.day = day
    r.invaildNewUserIds = make(map[string]int)
    // 各渠道留存信息
    r.channelSynthesizeMap = make(map[string]*ChannelSynthesizeInfo)
    // 新增用户
    r.totalNewUsers = []*NewUserLog{}
    GetAllNewPlayers(year, month, day, &r.totalNewUsers, "")
    // 有效用户
    r.vaildNewUsers = []*NewUserLog{}
    r.channleIdMap = make(map[string]int)
    // 渠道列表
    adminUserList := []*AdminInfo{}
    AdminMgr.GetAllUser(&adminUserList)
    for _, item := range adminUserList {
        r.checkAndNewChannelMap(item.Channel)
    }
    for _, item := range r.totalNewUsers {
        // 新建各渠道留存信息
        r.checkAndNewChannelMap(item.Channel)
        r.channelSynthesizeMap[item.Channel].NewUserCount++
        // 有效用户
        if (item.TotalOnlineSeconds >= ConfigMgr.VaildNewPlayerOnlineSecond && item.MatchTimes >= ConfigMgr.VaildNewPlayerMatchCount) {
            r.vaildNewUsers = append(r.vaildNewUsers, item)
            r.channelSynthesizeMap[item.Channel].VaildNewUserCount++
        } else {
            r.invaildNewUserIds[item.UserId] = 1
        }
        // 渠道列表
        if len(item.Channel) > 0 {
            r.channleIdMap[item.Channel] = 1
        }
    }
    r.afterXDayRemainPlayerCountMap = make(map[int]int)
}

// 计算N天后的留存
func (r *SynthesizeCalc)CalcLastXDayRemain(xDay int) {
    curDay := time.Date(r.year, time.Month(r.month), r.day, 0, 0, 0, 0, time.Local)
    queryDay := curDay.AddDate(0, 0, xDay)
    loginRecords := []*LoginRecord{}
    GetLoginRecords(queryDay.Year(), int(queryDay.Month()), queryDay.Day(), &loginRecords)
    if len(loginRecords) > 0 {
        // 计算留存
        r.calcRemainPlayerCount(r.vaildNewUsers, loginRecords, xDay)
    }
}

// 计算留存玩家数量
func (r *SynthesizeCalc)calcRemainPlayerCount(users []*NewUserLog, loginRecords []*LoginRecord, xDay int) {
    for _, userItem := range users {
        for _, recordItem := range loginRecords {
            if userItem.UserId == recordItem.UserId {
                r.checkAndNewChannelMap(userItem.Channel)
                r.channelSynthesizeMap[userItem.Channel].AfterXDayRemainPlayerCountMap[xDay]++
                r.afterXDayRemainPlayerCountMap[xDay]++
                break
            }
        }
    }
}

// 计算活跃
func (r *SynthesizeCalc)CalcLoginPlayer() {
    loginRecords := []*LoginRecord{}
    GetLoginRecords(r.year, r.month, r.day, &loginRecords)
    userIdMap := make(map[string]int)
    for _, recordItem := range loginRecords {
        // 是否重复
        _, exist := userIdMap[recordItem.UserId]
        if exist {
            continue
        }
        // 是否无效
        _, exist = r.invaildNewUserIds[recordItem.UserId]
        if exist {
            continue
        }
        userIdMap[recordItem.UserId] = 1
        r.checkAndNewChannelMap(recordItem.Channel)
        r.channelSynthesizeMap[recordItem.Channel].LoginUserCount++
        r.loginUserCount++
    }
}

// 获取渠道列表
func (r *SynthesizeCalc)GetChannelMap() (map[string]int){
    return r.channleIdMap
}

// 获取信息
// 新用户数量、有效用户数量、活跃用户、支付次数、支付玩家人数、支付金额
func (r *SynthesizeCalc)GetSynthesizeInfo() (int, int, int, int, int, int){
    newUserCount := len(r.totalNewUsers)
    vaildNewUserCount := len(r.vaildNewUsers)
    return newUserCount, vaildNewUserCount, r.loginUserCount, r.PayCount, r.PayPlayerCount, r.PayTotalAmount
}
// 获取渠道信息
// 新用户数量、有效用户数量、活跃用户、支付次数、支付玩家人数、支付金额
func (r *SynthesizeCalc)GetChannelSynthesizeInfo(channel string) (int, int, int, int, int, int){
    remainInfo := r.channelSynthesizeMap[channel]
    if nil == remainInfo {
        return 0, 0, 0, 0, 0, 0
    }
    return remainInfo.NewUserCount, remainInfo.VaildNewUserCount, 
    remainInfo.LoginUserCount, remainInfo.PayCount, 
    remainInfo.PayPlayerCount, remainInfo.PayTotalAmount
}

// 获取留存信息
func (r *SynthesizeCalc)GetRemainInfo(xDay int) (int){
    return r.afterXDayRemainPlayerCountMap[xDay]
}
// 获取渠道留存信息
func (r *SynthesizeCalc)GetChannelRemainInfo(channel string, xDay int) (int){
    remainInfo := r.channelSynthesizeMap[channel]
    if nil == remainInfo {
        return 0
    }
    return remainInfo.AfterXDayRemainPlayerCountMap[xDay]
}

// 计算付费信息
func (r *SynthesizeCalc)CalcPayInfo() {
    // 支付记录
    r.payRecords = []*PayLog{}
    GetAllPayLogList(r.year, r.month, r.day, &r.payRecords)
    payUserIdMap := make(map[string]int)
    for _, item := range r.payRecords {
        r.checkAndNewChannelMap(item.Channel)
        // 付费次数
        r.PayCount++
        r.channelSynthesizeMap[item.Channel].PayCount++
        // 付费金额
        r.PayTotalAmount += item.Amount
        r.channelSynthesizeMap[item.Channel].PayTotalAmount += item.Amount
        // 付费人数
        _, exist := payUserIdMap[item.UserId]
        if !exist {
            r.channelSynthesizeMap[item.Channel].PayPlayerCount++
            r.PayPlayerCount++
            payUserIdMap[item.UserId] = 1
        }
    }
}

// 检测渠道信息
func (r *SynthesizeCalc)checkAndNewChannelMap(channel string) {
    _, existChannel := r.channelSynthesizeMap[channel]
    if !existChannel {
        r.channelSynthesizeMap[channel] = &ChannelSynthesizeInfo{}
        r.channelSynthesizeMap[channel].Channel = channel
        r.channelSynthesizeMap[channel].AfterXDayRemainPlayerCountMap = make(map[int]int)
    }
}

// 生成日志
func (r *SynthesizeCalc)BuildSynthesizeInfoLogs() (map[string]*SynthesizeInfoLog){
    logMaps := make(map[string]*SynthesizeInfoLog)
    // 汇总日志
    totalLog := &SynthesizeInfoLog{}
    totalLog.Channel = ""
    totalLog.NewUserCount = len(r.totalNewUsers)
    totalLog.VaildNewUserCount = len(r.vaildNewUsers)
    totalLog.LoginUserCount = r.loginUserCount
    totalLog.PayCount = r.PayCount
    totalLog.PayTotalAmount = r.PayTotalAmount
    totalLog.PayPlayerCount = r.PayPlayerCount
    totalLog.After1DayRemainUserCount = r.afterXDayRemainPlayerCountMap[1]
    totalLog.After2DayRemainUserCount = r.afterXDayRemainPlayerCountMap[2]
    totalLog.After3DayRemainUserCount = r.afterXDayRemainPlayerCountMap[3]
    totalLog.After5DayRemainUserCount = r.afterXDayRemainPlayerCountMap[5]
    totalLog.After7DayRemainUserCount = r.afterXDayRemainPlayerCountMap[7]
    totalLog.DateTime = time.Date(r.year, time.Month(r.month), r.day, 0, 0, 0, 0, time.Local)
    logMaps[totalLog.Channel] = totalLog
    // 渠道日志
    for channel, info := range r.channelSynthesizeMap {
        if len(channel) == 0 {
            continue
        }
        log := &SynthesizeInfoLog{}
        log.Channel = channel
        log.NewUserCount = info.NewUserCount
        log.VaildNewUserCount = info.VaildNewUserCount
        log.LoginUserCount = info.LoginUserCount
        log.PayCount = info.PayCount
        log.PayTotalAmount = info.PayTotalAmount
        log.PayPlayerCount = info.PayPlayerCount
        log.After1DayRemainUserCount = info.AfterXDayRemainPlayerCountMap[1]
        log.After2DayRemainUserCount = info.AfterXDayRemainPlayerCountMap[2]
        log.After3DayRemainUserCount = info.AfterXDayRemainPlayerCountMap[3]
        log.After5DayRemainUserCount = info.AfterXDayRemainPlayerCountMap[5]
        log.After7DayRemainUserCount = info.AfterXDayRemainPlayerCountMap[7]
        log.DateTime = time.Date(r.year, time.Month(r.month), r.day, 0, 0, 0, 0, time.Local)
        logMaps[log.Channel] = log
    }
    return logMaps
}

func (r *SynthesizeCalc)PrintInfo() {
    fmt.Println("===============================================")
    fmt.Println("year:", r.year, "month:", r.month, "day:", r.day)
    fmt.Println("new user count:", len(r.totalNewUsers))
    fmt.Println("new vaild user count:", len(r.vaildNewUsers))
    fmt.Println("login player count:", r.loginUserCount)
    fmt.Println("pay count:", r.PayCount)
    fmt.Println("pay player count:", r.PayPlayerCount)
    fmt.Println("pay total amount", r.PayTotalAmount)
    fmt.Println("channel count:", len(r.channleIdMap))
    for _, remainInfo := range r.channelSynthesizeMap {
        printInfo := "channel remain info:" + remainInfo.Channel
        for day, count := range remainInfo.AfterXDayRemainPlayerCountMap {
            printInfo += "(day" + fmt.Sprintf("%d", day) + ":" + fmt.Sprintf("%d", count) + ") "
        }
        fmt.Println(printInfo)
    }
    fmt.Println("===============================================")
}