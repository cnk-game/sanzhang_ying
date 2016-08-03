package prize

import (
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"sync"
	"util"
)

type TaskPrize struct {
	TaskId         int    `bson:"taskId"`
	TaskTitle      string `bson:"taskTitle"`
	IcoRes         string `bson:"icoRes"`
	TaskType       int    `bson:"taskType"`
	TargetType     int    `bson:"targetType"`
	TargetValue    int    `bson:"targetValue"`
	IsAccumulated  int    `bson:"isAccumulated"`
	PrizeGold      int    `bson:"prizeGold"`
	PrizeDiamond   int    `bson:"prizeDiamond"`
	PrizeExp       int    `bson:"prizeExp"`
	PrizeScore     int    `bson:"prizeScore"`
	PrizeItemType  int    `bson:"prizeItemType"`
	PrizeItemCount int    `bson:"prizeItemCount"`
}

const (
	taskPrizeC = "task_prize"
)

func FindTaskPrizes() ([]*TaskPrize, error) {
	prizes := []*TaskPrize{}

	err := util.WithGameCollection(taskPrizeC, func(c *mgo.Collection) error {
		return c.Find(nil).All(&prizes)
	})
	return prizes, err
}

type TaskPrizeManager struct {
	sync.RWMutex
	prizes map[int]*TaskPrize
}

var taskPrizeManager *TaskPrizeManager

func init() {
	taskPrizeManager = &TaskPrizeManager{}
	taskPrizeManager.prizes = make(map[int]*TaskPrize)
}

func GetTaskPrizeManager() *TaskPrizeManager {
	return taskPrizeManager
}

func (m *TaskPrizeManager) Init() {
	m.Lock()
	defer m.Unlock()

	tasks, err := FindTaskPrizes()
	if err != nil {
		glog.Fatal(err)
	}

	for _, task := range tasks {
		m.prizes[task.TaskId] = task
	}
}

func (m *TaskPrizeManager) GetTaskPrize(taskId int) (TaskPrize, bool) {
	m.Lock()
	defer m.Unlock()

	t := m.prizes[taskId]
	if t != nil {
		return *t, true
	}

	return TaskPrize{}, false
}

func (m *TaskPrizeManager) GetTaskPrizes() map[int]*TaskPrize {
	return m.prizes
}
