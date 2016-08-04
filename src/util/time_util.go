package util

import (
	"github.com/golang/glog"
	"time"
)

func ParseTime(s string) time.Time {
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
	return t
}

func ParseTime2(s string) (time.Time, error) {
	t, err := time.ParseInLocation("2006-1-2 15:4:5", s, time.Local)
	return t, err
}

func CompareDate(t1 time.Time, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}

func CompareMonth(t1 time.Time, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month()
}

func IsBigDate(t1 time.Time, t2 time.Time) bool {
	glog.Info("IsBigDate t1 year ", t1.Year())
	glog.Info("IsBigDate t2 year ", t2.Year())
	glog.Info("IsBigDate t1 Month ", t1.Month())
	glog.Info("IsBigDate t2 Month ", t2.Month())
	glog.Info("IsBigDate t1 Day ", t1.Day())
	glog.Info("IsBigDate t2 Day ", t2.Day())

	glog.Info("IsBigDate t1 Hour ", t1.Hour())
	glog.Info("IsBigDate t2 Hour ", t2.Hour())
	if t1.Year() > t2.Year() {
		return true
	}
	if t1.Month() > t2.Month() {
		return true
	}

	if t1.Day() > t2.Day() {
		return true
	}

	if t1.Hour() > t2.Hour() {
		return true
	}

	return false
}

// 获取零点时刻
func GetZeroTime() time.Time {
	//	now := time.Now()
	//	dur, _ := time.ParseDuration(fmt.Sprintf("%vh%vm%vs", now.Hour(), now.Minute(), now.Second()))
	//	return now.Add(-dur)
	return GetOClock(0)
}

func GetOClock(hour int) time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, time.Local)
}
