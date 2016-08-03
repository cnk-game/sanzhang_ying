package chat

import (
	"code.google.com/p/goprotobuf/proto"
	"game/domain/forbidWords"
	domainGame "game/domain/game"
	domainUser "game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"pb"
	"time"
	"util"
	"fmt"
)

func ChatMsgHandler(m *pb.ServerMsg, sess *server.Session) []byte {
	player := domainUser.GetPlayer(sess.Data)

	msg := &pb.MsgChat{}
	err := proto.Unmarshal(m.GetMsgBody(), msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	glog.V(2).Info("****聊天消息msg:", msg)

	if time.Since(player.LastChatTime).Seconds() < 2 {
		msg := &pb.MsgShowTips{}
        content := fmt.Sprintf("您发言过于频繁，请休息一下")
        msg.UserId = proto.String(player.User.UserId)
        msg.Content = proto.String(content)
        domainUser.GetPlayerManager().SendClientMsg(player.User.UserId, int32(pb.MessageId_SHOW_TIPS), msg)
		return nil
	}

	if forbidWords.IsForbid(msg.GetContent()) {
		msg := &pb.MsgShowTips{}
        content := fmt.Sprintf("您的发言中包含敏感词汇，请注意")
        msg.UserId = proto.String(player.User.UserId)
        msg.Content = proto.String(content)
        domainUser.GetPlayerManager().SendClientMsg(player.User.UserId, int32(pb.MessageId_SHOW_TIPS), msg)
		return nil
	}

	if msg.GetMessageType() == pb.ChatMessageType_DESK || msg.GetMessageType() == pb.ChatMessageType_GIFT || msg.GetMessageType() == pb.ChatMessageType_GIFTALL {
		domainGame.GetDeskManager().OnChatMsg(player.User.UserId, msg)
	} else if msg.GetMessageType() == pb.ChatMessageType_BUGLE {
	    f, _ := domainUser.GetUserFortuneManager().GetUserFortune(player.User.UserId)
	    if f.Horn <= 0 {
	        msg := &pb.MsgShowTips{}
	        content := fmt.Sprintf("您的喇叭不足，请到商城兑换")
            msg.UserId = proto.String(player.User.UserId)
            msg.Content = proto.String(content)
            domainUser.GetPlayerManager().SendClientMsg(player.User.UserId, int32(pb.MessageId_SHOW_TIPS), msg)
	        return nil
	    }
	    domainUser.GetUserFortuneManager().EarnHorn(player.User.UserId, -1)
	    domainUser.GetUserFortuneManager().UpdateUserFortune(player.User.UserId)

		domainUser.GetPlayerManager().BroadcastClientMsg(m.GetMsgId(), msg)
		player.UserTasks.AccomplishTask(util.TaskAccomplishType_BUGLE_X_TIMES, 1, player.SendToClientFunc)
	}

	return nil
}
