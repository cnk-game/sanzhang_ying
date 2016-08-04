package util

import (
	"code.google.com/p/goprotobuf/proto"
	"pb"
)

func BuildSysBugle(content string) *pb.MsgChat {
	msg := &pb.MsgChat{}

	msg.MessageType = pb.ChatMessageType_SYSTEM.Enum()
	msg.Content = proto.String(content)

	return msg
}
