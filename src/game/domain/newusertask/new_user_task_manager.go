package newusertask

import (
	"code.google.com/p/goprotobuf/proto"
	"fmt"
	"github.com/golang/glog"
	mgo "gopkg.in/mgo.v2"
	"pb"
	"strconv"
	"sync"
	"time"
)

type NewUserTaskManager struct {
	sync.RWMutex
	newUserTask       map[string]*NewUserTask
	newUserTaskConfig map[string]*NewUserTaskInfoConfig
}

var newUserTaskManager *NewUserTaskManager

func init() {
	newUserTaskManager = &NewUserTaskManager{}
	newUserTaskManager.newUserTask = make(map[string]*NewUserTask)
	newUserTaskManager.newUserTaskConfig = make(map[string]*NewUserTaskInfoConfig)
}

func GetNewUserTaskManager() *NewUserTaskManager {
	return newUserTaskManager
}

func (m *NewUserTaskManager) Init() {
	taskConfigs, err := GetNewUserTaskConfig()
	if err != nil {
		glog.Info("GetNewUserTaskConfig err ", err)
		return
	}

	for _, taskConfig := range taskConfigs {
		m.newUserTaskConfig[taskConfig.Id] = taskConfig
	}
	return
}

func (m *NewUserTaskManager) InitUserTask(userId string) bool {
	glog.Info("InitUserTask in userId ", userId)
	task, err := FindNewUserTask(userId)
	if err != nil && err != mgo.ErrNotFound {
		glog.Info("LoadUserTask err ", err)
		return false
	}

	if err == mgo.ErrNotFound {
		for _, taskConfig := range m.newUserTaskConfig {
			taskTemp := UserTaskInfo{}
			if taskConfig.Id == "1" {
				taskTemp.Status = 1
			} else {
				taskTemp.Status = 0
			}

			taskTemp.Id = taskConfig.Id
			taskTemp.Remains = taskConfig.PlayCount
			task.Tasks[taskTemp.Id] = taskTemp
		}

		m.Lock()
		defer m.Unlock()

		m.newUserTask[userId] = task
		return true
	} else {
		return false
	}
}

func (m *NewUserTaskManager) LoadUserTask(userId string) bool {
	glog.Info("LoadUserTask in userId ", userId)

	_, bResult := m.GetUserTask(userId)
	if bResult {
		return true
	}

	task, err := FindNewUserTask(userId)
	if err != nil && err != mgo.ErrNotFound {
		glog.Info("LoadUserTask err ", err)
		return false
	}

	if err == mgo.ErrNotFound {
		return false
	} else {
		taskDate := GetTaskDate(task.Year, task.YearDay)
		glog.Info("LoadUserTask in taskDate ", taskDate)
		for _, taskInfo := range task.Tasks {
			id, _ := strconv.ParseFloat(taskInfo.Id, 64)
			Id := int(id)
			if Id <= (taskDate+1) && taskInfo.Status == 0 {
				taskInfo.Status = 1
				task.Tasks[taskInfo.Id] = taskInfo
			}
		}
		m.Lock()
		defer m.Unlock()

		m.newUserTask[userId] = task

		glog.Info("LoadUserTask out ", task)
		return true
	}
}

func (m *NewUserTaskManager) GetUserTask(userId string) (NewUserTask, bool) {
	m.RLock()
	defer m.RUnlock()

	task := m.newUserTask[userId]
	if task == nil {
		return NewUserTask{}, false
	}

	return *task, true

}

func (m *NewUserTaskManager) SaveUserTask(userId string) bool {
	glog.Info("SaveUserTask in userId ", userId)
	m.Lock()
	defer m.Unlock()

	item := m.newUserTask[userId]
	if item == nil {
		glog.Info("SaveUserTask out fail")
		return false
	} else {
		SaveNewUserTask(item)
		delete(m.newUserTask, userId)
		glog.Info("SaveUserTask out true")
		return true
	}
}

func IsTaskExpired(year int, day int) bool {
	if year == time.Now().Year() {
		if (time.Now().YearDay() - day) > 10 {
			return true
		} else {
			return false
		}
	} else {
		if (time.Now().YearDay() + (365 - day)) > 10 {
			return true
		} else {
			return false
		}
	}
}

func GetTaskDate(year int, day int) int {
	if year == time.Now().Year() {
		return (time.Now().YearDay() - day)
	} else {
		return (time.Now().YearDay() + (365 - day))
	}
}

func (m *NewUserTaskManager) BuildUserTask(userId string) *pb.MsgGetNewbeTaskListRes {
	msg := &pb.MsgGetNewbeTaskListRes{}
	task, _ := m.GetUserTask(userId)
	bIsExpired := false
	if task.IsGet == true {
		bIsExpired = true
	} else {
		bIsExpired = IsTaskExpired(task.Year, task.YearDay)
	}

	msg.IsExpired = proto.Bool(bIsExpired)
	msg.HfCount = proto.Int(task.HfCount)
	taskDate := GetTaskDate(task.Year, task.YearDay)
	msg.LeftDate = proto.Int(10 - taskDate)

	if bIsExpired == true {
		return msg
	}

	for _, v := range task.Tasks {
		tt := &pb.NewbeTask{}
		id, _ := strconv.ParseFloat(v.Id, 64)
		Id := int(id)
		tt.Id = proto.Int(Id)
		tt.Status = proto.Int(v.Status)
		tt.Remains = proto.Int(v.Remains)

		msg.Tasks = append(msg.Tasks, tt)
	}

	glog.Info("new user task BuildUserTask ", msg)

	return msg
}

func (m *NewUserTaskManager) GetUserTaskPrize(userId string, id int) (int, NewUserTask) {
	m.RLock()
	defer m.RUnlock()

	Id := fmt.Sprintf("%v", id)

	task := m.newUserTask[userId]
	if task == nil {
		return 4, *task
	}

	if id > 7 {
		return 4, *task
	}

	taskDate := GetTaskDate(task.Year, task.YearDay)
	taskDate = taskDate + 1
	if taskDate > 10 {
		return 4, *task
	}

	taskTemp := task.Tasks[Id]

	if taskTemp.Status == 3 {
		return 3, *task
	}

	if taskTemp.Status != 2 {
		return 2, *task
	}

	taskTemp.Status = 3
	task.Tasks[Id] = taskTemp
	task.HfCount += m.newUserTaskConfig[Id].PrizeCount

	SaveNewUserTask(task)

	return 1, *task
}

func (m *NewUserTaskManager) GetUserTaskHuafei(userId string, phone string) int {
	m.RLock()
	defer m.RUnlock()

	task := m.newUserTask[userId]
	if task == nil {
		return 4
	}

	err := GetNewUserPhone(phone)
	if err != mgo.ErrNotFound {
		return 3
	}

	taskDate := GetTaskDate(task.Year, task.YearDay)
	taskDate = taskDate + 1
	if taskDate > 10 {
		return 4
	}

	if task.IsGet == true {
		return 3
	}

	if task.HfCount < 10 {
		return 2
	}

	SaveNewUserTask(task)

	SaveNewUserPhone(userId, phone)
	SaveGetHuafeiLog(phone, userId, task.HfCount)

	task.HfCount -= 10
	task.IsGet = true

	return 1
}

func (m *NewUserTaskManager) CheckUserRechargeTask(userId string) int {
	m.RLock()
	defer m.RUnlock()

	task := m.newUserTask[userId]
	if task == nil {
		return 0
	}

	taskDate := GetTaskDate(task.Year, task.YearDay)
	taskDate = taskDate + 1
	if taskDate > 10 {
		return 0
	}

	if taskDate < 7 {
		return 0
	}

	taskInfo := task.Tasks["7"]
	if taskInfo.Status == 2 {
		return 0
	}

	taskInfo.Status = 2
	task.Tasks["7"] = taskInfo

	SaveNewUserTask(task)

	return 7
}

func (m *NewUserTaskManager) CheckUserCompleteInfoTask(userId string, channel string) int {
	glog.Info("CheckUserCompleteInfoTask in", userId, ", channel = ", channel)
	m.RLock()
	defer m.RUnlock()

	chlTep, _ := strconv.ParseFloat(channel, 64)
	chl := int(chlTep)

	if chl == 178 || chl == 184 || chl == 173 || chl == 186 || chl == 187 {
		return 0
	}

	task := m.newUserTask[userId]
	if task == nil {
		return 0
	}

	taskDate := GetTaskDate(task.Year, task.YearDay)
	taskDate = taskDate + 1
	if taskDate > 10 {
		return 0
	}

	taskInfo := task.Tasks["1"]
	if taskInfo.Status == 2 {
		return 0
	}

	taskInfo.Status = 2
	task.Tasks["1"] = taskInfo

	SaveNewUserTask(task)
	return 1
}

func (m *NewUserTaskManager) CheckUserChangeNicknameTask(userId string) int {
	glog.Info("CheckUserChangeNicknameTask in ")
	m.RLock()
	defer m.RUnlock()

	task := m.newUserTask[userId]
	if task == nil {
		return 0
	}

	taskDate := GetTaskDate(task.Year, task.YearDay)
	taskDate = taskDate + 1
	if taskDate > 10 {
		return 0
	}

	taskInfo := task.Tasks["2"]
	if taskInfo.Status == 2 {
		return 0
	}

	taskInfo.Status = 2
	task.Tasks["2"] = taskInfo

	glog.Info("CheckUserChangeNicknameTask out ", task.Tasks["2"])
	SaveNewUserTask(task)
	return 2
}

func (m *NewUserTaskManager) CheckUserPlayTask(userId string, gameType int, isWin bool, userChannel string) ([7]int, bool) {
	m.RLock()
	defer m.RUnlock()

	returnValue := [7]int{0, 0, 0, 0, 0, 0, 0}

	task := m.newUserTask[userId]
	if task == nil {
		return returnValue, false
	}

	taskDate := GetTaskDate(task.Year, task.YearDay)
	taskDate = taskDate + 1
	if taskDate > 10 {
		return returnValue, false
	}

	chlTep, _ := strconv.ParseFloat(userChannel, 64)

	chl := int(chlTep)

	bUpdate := false

	if gameType == 1 && task.Tasks["1"].Status != 2 && (chl == 178 || chl == 184 || chl == 173 || chl == 186 || chl == 187) {
		if task.Tasks["1"].Remains > 0 {
			bUpdate = true
			taskInfo := task.Tasks["1"]
			taskInfo.Remains -= 1
			taskInfo.Status = 1
			if taskInfo.Remains == 0 {
				taskInfo.Status = 2
				returnValue[0] = 1
			}

			task.Tasks["1"] = taskInfo
		}
	}

	if taskDate < 3 {
		return returnValue, bUpdate
	}

	if task.Tasks["3"].Remains > 0 {
		bUpdate = true
		taskInfo := task.Tasks["3"]
		taskInfo.Remains -= 1
		taskInfo.Status = 1
		if taskInfo.Remains == 0 {
			taskInfo.Status = 2
			returnValue[2] = 1
		}
		task.Tasks["3"] = taskInfo
	}

	if taskDate < 4 {
		return returnValue, bUpdate
	}

	if gameType == 4 {
		if task.Tasks["4"].Remains > 0 {
			bUpdate = true
			taskInfo := task.Tasks["4"]
			taskInfo.Remains -= 1
			taskInfo.Status = 1
			if taskInfo.Remains == 0 {
				taskInfo.Status = 2
				returnValue[3] = 1
			}
			task.Tasks["4"] = taskInfo
		}
	}

	if taskDate < 5 {
		return returnValue, bUpdate
	}

	if gameType == 2 {
		if task.Tasks["5"].Remains > 0 {
			bUpdate = true
			taskInfo := task.Tasks["5"]
			taskInfo.Remains -= 1
			taskInfo.Status = 1
			if taskInfo.Remains == 0 {
				taskInfo.Status = 2
				returnValue[4] = 1
			}

			task.Tasks["5"] = taskInfo
		}
	}

	if taskDate < 6 {
		return returnValue, bUpdate
	}

	if gameType == 2 && isWin {
		if task.Tasks["6"].Remains > 0 {
			bUpdate = true
			taskInfo := task.Tasks["6"]
			taskInfo.Remains -= 1
			taskInfo.Status = 1
			if taskInfo.Remains == 0 {
				taskInfo.Status = 2
				returnValue[5] = 1
			}

			task.Tasks["6"] = taskInfo
		}
	}

	return returnValue, bUpdate
}
