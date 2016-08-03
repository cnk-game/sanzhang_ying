package newusertask

import (
	//"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	//"sync"
	"gopkg.in/mgo.v2/bson"
	"time"
	"util"
)

type UserTaskInfo struct {
	Id      string `bson:"id"`
	Status  int    `bson:"status"`
	Remains int    `bson:"remains"`
}

type NewUserTask struct {
	UserId  string                  `bson:"userId"`
	HfCount int                     `bson:"hfCount"`
	YearDay int                     `bson:"yearDay"`
	Year    int                     `bson:"year"`
	IsGet   bool                    `bson:"isGet"`
	Tasks   map[string]UserTaskInfo `bson:"tasks"`
}

type NewUserTaskInfoConfig struct {
	Id         string `bson:"id"`
	Type       int    `bson:"type"`
	GameType   int    `bson:"gameType"`
	PlayCount  int    `bson:"playCount"`
	PrizeCount int    `bson:"prizeCount"`
}

type NewUserTaskPhone struct {
	Phone  string    `bson:"phone"`
	UserId string    `bson:"userId"`
	Time   time.Time `bson:"time"`
}

type NewUserHuaFeiLog struct {
	Phone  string    `bson:"phone"`
	UserId string    `bson:"userId"`
	Time   time.Time `bson:"time"`
	Count  int       `bson:"count"`
}

const (
	newUserTaskCofnigC = "new_user_task_config"
	newUserTaskC       = "new_user_task"
	newUserTaskPhoneC  = "new_user_task_phone"
	newUserHuafeiLogC  = "new_user_huafei_log"
)

//type:
//1 play in gameType
//2 modify nickname
//3	play in any gameType
//4 play in gameType win
//5 recharge any

func GetNewUserTaskConfig() ([]*NewUserTaskInfoConfig, error) {
	tasks := []*NewUserTaskInfoConfig{}

	err := util.WithGameCollection(newUserTaskCofnigC, func(c *mgo.Collection) error {
		return c.Find(nil).All(&tasks)
	})
	return tasks, err
}

func FindNewUserTask(userId string) (*NewUserTask, error) {
	task := &NewUserTask{}
	task.UserId = userId
	task.HfCount = 0
	task.Year = time.Now().Year()
	task.YearDay = time.Now().YearDay()
	task.IsGet = false
	task.Tasks = map[string]UserTaskInfo{}

	err := util.WithGameCollection(newUserTaskC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"userId": userId}).One(task)
	})

	return task, err
}

func SaveNewUserTask(newUserTask *NewUserTask) error {
	return util.WithGameCollection(newUserTaskC, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"userId": newUserTask.UserId}, newUserTask)
		return err
	})
}

func SaveNewUserPhone(userId string, phone string) error {
	newUserPhone := &NewUserTaskPhone{}
	newUserPhone.UserId = userId
	newUserPhone.Time = time.Now()
	newUserPhone.Phone = phone
	return util.WithGameCollection(newUserTaskPhoneC, func(c *mgo.Collection) error {
		return c.Insert(newUserPhone)
	})
}

func GetNewUserPhone(phone string) error {
	newUserPhone := &NewUserTaskPhone{}
	newUserPhone.Phone = phone

	err := util.WithGameCollection(newUserTaskPhoneC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"phone": phone}).One(newUserPhone)
	})

	return err
}

func SaveGetHuafeiLog(phone string, userId string, huafei int) error {
	newUserHuaFeiLog := &NewUserHuaFeiLog{}
	newUserHuaFeiLog.UserId = userId
	newUserHuaFeiLog.Time = time.Now()
	newUserHuaFeiLog.Phone = phone
	newUserHuaFeiLog.Count = huafei
	return util.WithGameCollection(newUserHuafeiLogC, func(c *mgo.Collection) error {
		return c.Insert(newUserHuaFeiLog)
	})
}
