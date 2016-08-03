package prize

import (
	"code.google.com/p/goprotobuf/proto"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"pb"
	"time"
	"util"
)

type UserTask struct {
	UserId    string    `bson:"userId"`
	TaskId    int       `bson:"taskId"`
	CurValue  int64     `bson:"curValue"`
	IsGained  bool      `bson:"isGained"`
	ResetTime time.Time `bson:"resetTime"`
	hashCode  *util.HashCode
}

const (
	userTaskC = "user_task"
)

func (task *UserTask) HashCode() *util.HashCode {
	return task.hashCode
}

func (task *UserTask) SetHashCode(hashCode *util.HashCode) {
	task.hashCode = hashCode
}

func (task *UserTask) BuildMessage() *pb.UserPrizeTaskDef {
	msg := &pb.UserPrizeTaskDef{}

	msg.TaskId = proto.Int(task.TaskId)
	msg.CompleteValue = proto.Int64(task.CurValue)
	msg.IsGain = proto.Bool(task.IsGained)

	return msg
}

func (task *UserTask) ResetTask() {
	now := time.Now()
	if !util.CompareDate(now, task.ResetTime) {
		task.CurValue = 0
		task.IsGained = false
		task.ResetTime = now
	}
}

func FindUserTasks(userId string) ([]*UserTask, error) {
	rs := []*UserTask{}
	err := util.WithUserCollection(userTaskC, func(c *mgo.Collection) error {
		return c.Find(bson.M{"userId": userId}).All(&rs)
	})

	if err == nil {
		for _, r := range rs {
			r.SetHashCode(util.NewHashCode(r))
		}
	}
	return rs, err
}

func SaveUserTask(t *UserTask) error {
	hashCode := util.NewHashCode(t)
	if t.HashCode() != nil && t.HashCode().Compare(hashCode) {
		return nil
	}

	return util.WithSafeUserCollection(userTaskC, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"userId": t.UserId, "taskId": t.TaskId}, t)
		if err == nil {
			// 保存成功
			t.SetHashCode(hashCode)
		}
		return err
	})
}
