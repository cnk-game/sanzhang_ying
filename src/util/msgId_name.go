package util

import (
	"fmt"
	"pb"
)

func GetMsgIdName(msgId int32) string {
	name := pb.ServerMsgId_name[msgId]
	if name != "" {
		return name
	}

	name = pb.MessageId_name[msgId]
	if name != "" {
		return name
	}

	return fmt.Sprintf("%v", msgId)
}
