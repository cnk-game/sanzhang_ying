package handlers

import (
	"game/domain/user"
	"game/server"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"pb"
	"sync"
	"time"
	"util"
)

type MsgRegistry struct {
	registry           map[int32]func(msg *pb.ServerMsg, sess *server.Session) []byte
	unLoginMsgRegistry map[int32]bool
	mu                 sync.RWMutex
}

func (registry *MsgRegistry) RegisterMsg(msgId int32, f func(msg *pb.ServerMsg, sess *server.Session) []byte) {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	registry.registry[msgId] = f
}

func (registry *MsgRegistry) RegisterUnLoginMsg(msgId int32) {
	registry.mu.Lock()
	defer registry.mu.Unlock()

	registry.unLoginMsgRegistry[msgId] = true
}

func (registry *MsgRegistry) isUnLoginMsg(msgId int32) bool {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	_, ok := registry.unLoginMsgRegistry[msgId]
	return ok
}

func (registry *MsgRegistry) getHandler(msgId int32) func(msg *pb.ServerMsg, sess *server.Session) []byte {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	return registry.registry[msgId]
}

func (registry *MsgRegistry) DispatchMsg(msg *pb.ServerMsg, sess *server.Session) []byte {
	glog.V(2).Info("===>收到客户端消息msgId:", util.GetMsgIdName(msg.GetMsgId()))

	start := time.Now()
	defer func() {
		elapseTime := time.Since(start)
		if elapseTime.Seconds() > 0.1 {
			p := user.GetPlayer(sess.Data)
			if p != nil {
				user.SaveSlowMsg(p.User.UserId, util.GetMsgIdName(msg.GetMsgId()), start, elapseTime.String())
			}
		}
	}()

	if !registry.isUnLoginMsg(msg.GetMsgId()) {
		if user.GetPlayer(sess.Data) == nil {
			glog.Info("===>玩家未登录msgId:", util.GetMsgIdName(msg.GetMsgId()))
			return nil
		}
	}

	f := registry.getHandler(msg.GetMsgId())
	if f == nil {
		glog.Error("msgId: ", util.GetMsgIdName(msg.GetMsgId()), " has no handler")
		return nil
	}

	return f(msg, sess)
}

func (registry *MsgRegistry) RegisterHandlers(r *mux.Router) {
	registerHandlers(r)
}

var registry *MsgRegistry

func init() {
	registry = &MsgRegistry{
		registry:           make(map[int32]func(msg *pb.ServerMsg, sess *server.Session) []byte),
		unLoginMsgRegistry: make(map[int32]bool),
		mu:                 sync.RWMutex{},
	}
}

func GetMsgRegistry() *MsgRegistry {
	return registry
}
