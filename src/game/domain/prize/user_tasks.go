package prize

import (
	"code.google.com/p/goprotobuf/proto"
	"github.com/golang/glog"
	"pb"
)

type UserTasks struct {
	UserId string
	tasks  map[int]*UserTask
}

func NewUserTasks(userId string) *UserTasks {
	tasks := &UserTasks{}
	tasks.UserId = userId
	tasks.tasks = make(map[int]*UserTask)

	ts, err := FindUserTasks(userId)
	if err == nil {
		for _, t := range ts {
			tasks.tasks[t.TaskId] = t
		}
	}

	tasks.ResetDailyTasks()

	return tasks
}

func (tasks *UserTasks) ResetDailyTasks() {
	for _, item := range tasks.tasks {
		task, ok := GetTaskPrizeManager().GetTaskPrize(item.TaskId)
		if !ok {
			continue
		}
		if task.TaskType == 2 {
			// 日常任务
			item.ResetTask()
		}
	}
}

func (tasks *UserTasks) SaveTasks() {
	for _, t := range tasks.tasks {
		glog.V(2).Info("====>保存任务:", t)
		err := SaveUserTask(t)
		if err != nil {
			glog.Error(err)
		}
	}
}

func (tasks *UserTasks) IsTaskAccomplished(taskId int) bool {
	t := tasks.tasks[taskId]
	if t == nil {
		return false
	}

	taskTemplate, ok := GetTaskPrizeManager().GetTaskPrize(taskId)
	if !ok {
		glog.Error("找不到任务模板taskid:", taskId)
		return false
	}

	return t.CurValue >= int64(taskTemplate.TargetValue)
}

func (tasks *UserTasks) IsTaskGained(taskId int) bool {
	t := tasks.tasks[taskId]
	if t == nil {
		return false
	}
	return t.IsGained
}

func (tasks *UserTasks) SetGainTask(taskId int) {
	t := tasks.tasks[taskId]
	if t == nil {
		return
	}

	t.IsGained = true
}

func (tasks *UserTasks) AccomplishTask(accomplishType int, taskValue int64, f func(msgId int32, body proto.Message)) {
	updateTasks := []int{}
	taskTemplates := GetTaskPrizeManager().GetTaskPrizes()
	for _, t := range taskTemplates {
		if t.TargetType == accomplishType {
			glog.V(2).Info("=====>任务:", t)
			if t.IsAccumulated > 0 {
				task := tasks.tasks[t.TaskId]
				if task == nil {
					task = &UserTask{}
					task.UserId = tasks.UserId
					task.TaskId = t.TaskId
					tasks.tasks[t.TaskId] = task
				}
				task.CurValue += int64(taskValue)
				updateTasks = append(updateTasks, t.TaskId)
			} else {
				if taskValue >= int64(t.TargetValue) {
					// 任务完成
					task := tasks.tasks[t.TaskId]
					if task == nil {
						task = &UserTask{}
						task.UserId = tasks.UserId
						task.TaskId = t.TaskId
						tasks.tasks[t.TaskId] = task
						task.CurValue = taskValue
					}
					updateTasks = append(updateTasks, t.TaskId)
				}
			}
		}
	}

	if len(updateTasks) > 0 {
		updateMsg := &pb.MsgUpdateTaskRes{}
		for _, t := range updateTasks {
			msgItem := tasks.BuildMessage(t)
			if msgItem != nil {
				updateMsg.Task = append(updateMsg.Task, msgItem)
			}
		}
		glog.V(2).Info("====>更新任务:", updateMsg)
		if f != nil {
			f(int32(pb.MessageId_UPDATE_TASK_STATUS), updateMsg)
		}
	}
}

func (tasks *UserTasks) BuildMessage(taskId int) *pb.UserPrizeTaskDef {
	t := tasks.tasks[taskId]
	if t == nil {
		return nil
	}
	return t.BuildMessage()
}

func (tasks *UserTasks) BuildMessage2() []*pb.UserPrizeTaskDef {
	msg := []*pb.UserPrizeTaskDef{}

	for _, t := range tasks.tasks {
		msg = append(msg, t.BuildMessage())
	}

	return msg
}
