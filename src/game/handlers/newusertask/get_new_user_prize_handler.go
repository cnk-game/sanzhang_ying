package newUserTask

import (
	"code.google.com/p/goprotobuf/proto"
	domainUserTask "game/domain/newusertask"
	domainUser "game/domain/user"
	"game/server"
	"pb"
	"strconv"
)

func GetPrizeHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	msg := &pb.MsgGetNewbeTaskReward{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		return nil
	}

	id := int(msg.GetId())
	userId := player.User.UserId
	res := &pb.MsgGetNewbeTaskRewardRes{}

	result, userTask := domainUserTask.GetNewUserTaskManager().GetUserTaskPrize(userId, id)
	if result == 1 {
		res.Code = pb.MsgGetNewbeTaskRewardRes_OK.Enum()
		res.Reason = proto.String("领取成功")
	} else if result == 2 {
		res.Code = pb.MsgGetNewbeTaskRewardRes_FAILED_TASK_NO_COMPLETED.Enum()
		res.Reason = proto.String("任务未达成")
	} else if result == 3 {
		res.Code = pb.MsgGetNewbeTaskRewardRes_FAILED_TASK_REWARDED.Enum()
		res.Reason = proto.String("任务已经领取过")
	} else {
		res.Code = pb.MsgGetNewbeTaskRewardRes_FAILED.Enum()
		res.Reason = proto.String("任务信息错误")
	}

	res.Id = proto.Int(id)
	res.HfCount = proto.Int(userTask.HfCount)

	for _, v := range userTask.Tasks {
		tt := &pb.NewbeTask{}
		id, _ := strconv.ParseFloat(v.Id, 64)
		Id := int(id)
		tt.Id = proto.Int(Id)
		tt.Status = proto.Int(v.Status)
		tt.Remains = proto.Int(v.Remains)

		res.Tasks = append(res.Tasks, tt)
	}

	return server.BuildClientMsg(m.GetMsgId(), res)
}
