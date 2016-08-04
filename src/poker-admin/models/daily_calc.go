package models

import (
    "fmt"
    "time"
)


var lastCalcDay int = 0


// 检查历史5天丢失的记录
func CheckOldLogs() {
    now := time.Now()
    for i := 1; i < 5; i++ {
        queryDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
        queryDay = queryDay.AddDate(0, 0, -i)
        //fmt.Println("===============>check old log, time : ", queryDay)
        
        count := GetSynthesizeLogCount(queryDay.Year(), int(queryDay.Month()), queryDay.Day())
        //fmt.Println("===============>query old log, time : ", queryDay, " count : ", count)
        if 0 == count {
            calcDaySynthesizeInfo(&queryDay)
        }
    }
}


// 每日凌晨1点计算
func DailyCalc() {
    for {
        time.Sleep(time.Minute * 1)
        now := time.Now()
        fmt.Println("==========================>check daily calc, time : ", now)
        // 凌晨1点
        if now.Hour() != 1 {
            continue
        }
        // 今日已经计算
        if lastCalcDay == now.Day() {
            continue
        }
        calcSyntheizeInfo()
        lastCalcDay = now.Day()
        fmt.Println("==========================>daily calc, time : ", now)
    }
}

// 计算综合信息
func calcSyntheizeInfo() {
    now := time.Now()
    fmt.Println("calcSyntheizeInfo", now)
    calcBefore1DaySynthesizeInfo(&now)
    calcBeforeXDayRemainInfo(2, 1, &now)
    calcBeforeXDayRemainInfo(3, 2, &now)
    calcBeforeXDayRemainInfo(4, 3, &now)
    calcBeforeXDayRemainInfo(6, 5, &now)
    calcBeforeXDayRemainInfo(8, 7, &now)
}

func calcBefore1DaySynthesizeInfo(now *time.Time) {
    // 1天前 => 计算新增、活跃、付费
    queryDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
    queryDay = queryDay.AddDate(0, 0, -1)
    calc := &SynthesizeCalc{}
    calc.Init(queryDay.Year(), int(queryDay.Month()), queryDay.Day())
    calc.CalcLoginPlayer()
    calc.CalcPayInfo()
    logs := calc.BuildSynthesizeInfoLogs()
    for _, log := range logs {
        InsertSynthesizeInfoLog(log)
    }
}

func calcDaySynthesizeInfo(queryDay *time.Time) {
    calc := &SynthesizeCalc{}
    calc.Init(queryDay.Year(), int(queryDay.Month()), queryDay.Day())
    calc.CalcLoginPlayer()
    calc.CalcPayInfo()
    logs := calc.BuildSynthesizeInfoLogs()
    for _, log := range logs {
        InsertSynthesizeInfoLog(log)
    }
    calcRemainInfo(calc, 1, queryDay)
    calcRemainInfo(calc, 2, queryDay)
    calcRemainInfo(calc, 3, queryDay)
    calcRemainInfo(calc, 5, queryDay)
    calcRemainInfo(calc, 7, queryDay)
}

func calcBeforeXDayRemainInfo(beforeXDay, afterXDayRemain int, now *time.Time) {
    // X天前 => X日留存
    queryDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
    queryDay = queryDay.AddDate(0, 0, -beforeXDay)
    calc := &SynthesizeCalc{}
    calc.Init(queryDay.Year(), int(queryDay.Month()), queryDay.Day())
    calc.CalcLastXDayRemain(afterXDayRemain)

    UpdateXDayRemainBySynthesizeInfoLog(afterXDayRemain, calc.GetRemainInfo(afterXDayRemain), "", queryDay)
    channelMap := calc.GetChannelMap()
    for channel, _ := range channelMap {
        UpdateXDayRemainBySynthesizeInfoLog(afterXDayRemain, calc.GetChannelRemainInfo(channel, afterXDayRemain), channel, queryDay)
    }
}

func calcRemainInfo(calc *SynthesizeCalc, afterXDayRemain int, queryDay *time.Time) {
    calc.CalcLastXDayRemain(afterXDayRemain)
    UpdateXDayRemainBySynthesizeInfoLog(afterXDayRemain, calc.GetRemainInfo(afterXDayRemain), "", *queryDay)
    channelMap := calc.GetChannelMap()
    for channel, _ := range channelMap {
        UpdateXDayRemainBySynthesizeInfoLog(afterXDayRemain, calc.GetChannelRemainInfo(channel, afterXDayRemain), channel, *queryDay)
    }
}