package models

import (
    "time"
    "fmt"
)

var StatisticsMgr *StatisticsManager

type CacheGameInfo struct {
    data *GameInfo
    time time.Time
}

type CacheOnlineInfo struct {
    onlineLogs []*OnlineLog
    lastReportTime time.Time
}

type StatisticsManager struct {
    gameInfo *CacheGameInfo
    onlineInfo *CacheOnlineInfo
}

func init() {
    StatisticsMgr = &StatisticsManager{}
    StatisticsMgr.Init()
}

func (m *StatisticsManager)Init() {
    m.gameInfo = &CacheGameInfo{}
    m.gameInfo.data = &GameInfo{}
    GetNowGameInfo(m.gameInfo.data)
    m.gameInfo.time = time.Now()

    m.onlineInfo = &CacheOnlineInfo{}
    m.onlineInfo.onlineLogs = []*OnlineLog{}
    LoadTodayOnlineLog(true, time.Now(), &m.onlineInfo.onlineLogs)
    for _, log := range m.onlineInfo.onlineLogs {
        if m.onlineInfo.lastReportTime.Before(log.DateTime) {
            m.onlineInfo.lastReportTime = log.DateTime
        }
    }

    fmt.Println("StatisticsManager Init")
}

func (m *StatisticsManager)GetMainGameInfo(channleId string) (*GameInfo){
    now := time.Now()
    duration := now.Sub(m.gameInfo.time)
    if duration >= (time.Duration(60)*time.Second) {
        m.gameInfo.data = &GameInfo{}
        GetNowGameInfo(m.gameInfo.data)
        m.gameInfo.time = time.Now()
    }

    if IsAdminChannel(channleId) {
        return m.gameInfo.data
    } else {
        data := &GameInfo{}
        for _, item := range m.gameInfo.data.ChannelGameInfoList {
            if item.Channel == channleId {
                data.TotalPlayerCount = item.TotalPlayerCount
                data.TodayNewPlayerCount = item.TodayNewPlayerCount
                data.TodayTotalNewPlayerCount = item.TodayTotalNewPlayerCount
                data.TodayLoginPlayerCount = item.TodayLoginPlayerCount
                data.TodayPayPlayerCount = item.TodayPayPlayerCount
                data.TodayPayCount = item.TodayPayCount
                data.TodayTotalPay = item.TodayTotalPay
            }
        }
        return data
    }
}

func (m* StatisticsManager)GetOnlineInfo() ([]*OnlineLog) {
    now := time.Now()
    if m.onlineInfo.lastReportTime.Day() != now.Day() {
        m.onlineInfo.lastReportTime = now
        m.onlineInfo.onlineLogs = []*OnlineLog{}
    }
    duration := now.Sub(m.onlineInfo.lastReportTime)
    if duration >= (time.Duration(5)*time.Minute) {
        newOnlineLogs := []*OnlineLog{}
        LoadTodayOnlineLog(false, m.onlineInfo.lastReportTime, &newOnlineLogs)
        for _, log := range newOnlineLogs {
            m.onlineInfo.onlineLogs = append(m.onlineInfo.onlineLogs, log)
            if m.onlineInfo.lastReportTime.Before(log.DateTime) {
                m.onlineInfo.lastReportTime = log.DateTime
            }
        }
    }
    return m.onlineInfo.onlineLogs
}