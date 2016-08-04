package GTGlobal

// service for log
// service for global mutex

import (
    "sync"
    "fmt"
    "os"
    "time"
    "strings"
)

var _instance *gtGlobal

type gtLog struct {
    file *os.File
    ch   chan string
    path string
    name string
}

type gtGlobal struct {
    sync.RWMutex
    GTLoger    *gtLog
}


func instance() *gtGlobal {
   if _instance == nil {
       _instance = new(gtGlobal)
   }
   return _instance
}


func Lock() {
    gt := instance()
    gt.Lock()
}


func Unlock() {
    gt := instance()
    gt.Unlock()
}


func Init(path, filename string) {
    gt := instance()
    gt.GTLoger = new(gtLog)
    gt.GTLoger.Init(path, filename)
}


func (this *gtLog)Init(path, filename string) {
    var fullname string
    if strings.LastIndex(path, "/") == len(path)-1 {
        fullname = path + filename
    } else {
        fullname = path + "/" + filename
    }
    
    f, err := os.OpenFile(fullname, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
    if(err != nil){
        fmt.Println(err)
    }
    this.file = f
    this.ch = make(chan string)
    this.path = path
    this.name = filename

    go this.writing()
    go this.updateing()
}


func GTLog() *gtLog {
    gt := instance()
    return gt.GTLoger
}


func (this *gtLog)Debug(format string, args ...interface{}) {
    var line string
    s := this.format_time()
    format = s + " [Debug] " + format + "\n"
    if len(args) > 0 {
        line = fmt.Sprintf(format, args...)
    } else {
        line = fmt.Sprintf(format)
    }
    
    this.ch <- line
}

func (this *gtLog)Info(format string, args ...interface{}) {
    var line string
    s := this.format_time()
    format = s + " [Info] " + format + "\n"
    if len(args) > 0 {
        line = fmt.Sprintf(format, args...)
    } else {
        line = fmt.Sprintf(format)
    }
    
    this.ch <- line
}

func (this *gtLog)Error(format string, args ...interface{}) {
    var line string
    s := this.format_time()
    format = s + " [Error] " + format + "\n"
    if len(args) > 0 {
        line = fmt.Sprintf(format, args...)
    } else {
        line = fmt.Sprintf(format)
    }
    
    this.ch <- line
}

func (this *gtLog)Warn(format string, args ...interface{}) {
    var line string
    s := this.format_time()
    format = s + " [Warn] " + format + "\n"
    if len(args) > 0 {
        line = fmt.Sprintf(format, args...)
    } else {
        line = fmt.Sprintf(format)
    }
    
    this.ch <- line
}

func (this *gtLog)format_time() string {
    t := time.Now()
    return fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d.%d", t.Year(),t.Month(),t.Day(),t.Hour(),t.Minute(),t.Second(),t.Nanosecond())
}

func (this *gtLog)writing() {
    for {
        select {
        case line,_ := <-this.ch:
            this.file.WriteString(line)
        }
    }
}

func (this *gtLog)updateing() {
    start := time.Now()
    for {
        time.Sleep(10)
        if time.Now().Day() != start.Day(){
            this.renamefile()
            start = time.Now()
        }
    }
    
}

func (this *gtLog)renamefile() {
    oldf := this.file
    
    newf, err := os.OpenFile(this.name, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
    if(err != nil){
        fmt.Println(err)
    }
    
    this.file = newf
    oldf.Close()
    
    t := time.Now()
    t = t.AddDate(0,0,-1)
    bakname := fmt.Sprintf("_%04d_%02d_%02d", t.Year(),t.Month(),t.Day())
    bakname = this.name + bakname
    os.Rename(this.name, bakname)
}
