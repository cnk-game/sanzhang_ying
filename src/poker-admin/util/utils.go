package util

import (
    "fmt"
    "log"
    "runtime"
    "time"
)

func PrintPanicStack() {
    if r := recover(); r != nil {
        log.Printf("%v", r)
        for i := 0; i < 10; i++ {
            funcName, file, line, ok := runtime.Caller(i)
            if ok {
                log.Printf("frame %v:[func:%v, file:%vf, line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
            }
        }
    }
}

func GetFunCaller(frame int) string {
    funcName, file, line, ok := runtime.Caller(frame)
    if ok {
        return fmt.Sprintf("frame %v:[func:%v, file:%vf, line:%v]\n", frame, runtime.FuncForPC(funcName).Name(), file, line)
    }
    return ""
}

func ParseTime(s string) time.Time {
    t, _ := time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
    return t
}
func ParseTimeUTC(s string) time.Time{
    t, _ := time.ParseInLocation("2006-01-02 15:04:05", s, time.UTC)
    return t
}