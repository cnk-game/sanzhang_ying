package server

import (
	"code.google.com/p/go.net/websocket"
	"code.google.com/p/goprotobuf/proto"
	"github.com/golang/glog"
	"pb"
	"strings"
	"sync"
	"time"
	"util"
)

type Session struct {
	conn      *websocket.Conn
	IP        string
	mq        chan *pb.ServerMsg
	Data      interface{}
	LoggedIn  bool
	exitChan  chan bool
	cleanOnce sync.Once
	kickOnce  sync.Once
	OnLogout  func()
}

func newSess(conn *websocket.Conn) *Session {
	sess := &Session{}
	sess.conn = conn

	return sess
}

func (s *Session) cleanSess() {
	s.cleanOnce.Do(func() {
		glog.V(2).Info("===>清理session:", s)
		s.conn.Close()
		if s.mq != nil {
			close(s.mq)
		}
	})
}

func (s *Session) Kickout() {
	s.kickOnce.Do(func() {
		close(s.exitChan)
	})
}

func (s *Session) Run(dispatcher MsgDispatcher) {
	defer util.PrintPanicStack()
	s.IP = s.conn.Request().Header.Get("X-Real-Ip")
	if len(s.IP) == 0 {
		s.IP = strings.Split(s.conn.Request().RemoteAddr, ":")[0]
	}

	s.mq = make(chan *pb.ServerMsg, 100)
	s.exitChan = make(chan bool)

	GetServerInstance().waitGroup.Add(1)

	glog.Info("===>打开session:", s)

	defer func() {
		glog.Info("disconnected:", s)
		util.PrintPanicStack()
		s.cleanSess()
		s.logout()
		GetServerInstance().waitGroup.Done()
	}()

	for {
		select {
		case <-GetServerInstance().stopChan:
			return
		case msg, ok := <-s.mq:
			if !ok {
				return
			}

            if msg.GetClient() {
                glog.V(2).Info("收到客户端消息msgId:", util.GetMsgIdName(msg.GetMsgId()))
            } else {
                glog.V(2).Info("收到服务器消息msgId:", util.GetMsgIdName(msg.GetMsgId()))
            }

			res := dispatcher.DispatchMsg(msg, s)
			if res != nil {
				s.SendToClient(res)
			}
		case <-s.exitChan:
			glog.Info("==>Kickout sess:", s)
			return
		}
	}
}

func (s *Session) logout() {
	if s.LoggedIn && s.OnLogout != nil {
	    glog.Info("Session。logout, s=", s)
		s.OnLogout()
	}
	s.LoggedIn = false
}

func (s *Session) SendMQ(msg *pb.ServerMsg) bool {
	ret := true

	defer func() {
		if r := recover(); r != nil {
			ret = false
		}
	}()
	s.mq <- msg

	return ret
}

func (s *Session) SendToClient(msg []byte) {
	ret := true

	defer func() {
		if r := recover(); r != nil {
			ret = false
			clientMsg := &pb.ClientMsg{}
			proto.Unmarshal(msg, clientMsg)
			glog.Info("send on closed channel sess:", s.IP, " msgId:", util.GetMsgIdName(clientMsg.GetMsgId()))
		}
	}()

	if msg != nil {
		if glog.V(2) {
			clientMsg := &pb.ClientMsg{}
			err := proto.Unmarshal(msg, clientMsg)
			if err != nil {
				glog.Error("====>SendToClient unmarshal msg failed:", err)
			}

			glog.Info("==>向客户端发送消息:", util.GetMsgIdName(clientMsg.GetMsgId()))
		}
		s.conn.SetWriteDeadline(time.Now().Add(time.Second))
		err := websocket.Message.Send(s.conn, msg)
		if err != nil {
			glog.Info("===>发送失败err:", err, " sess:", s)
			s.cleanSess()
			return
		}
		s.conn.SetWriteDeadline(time.Time{})
	}
}

func BuildClientMsg(msgId int32, body proto.Message) []byte {
	if _, ok := pb.MessageId_name[msgId]; !ok {
		glog.Warning("build client msg failed client msgId:", msgId, " does not exist")
		return nil
	}

	msg := &pb.ClientMsg{}
	msg.MsgId = proto.Int32(int32(msgId))

	if body != nil {
		d, err := proto.Marshal(body)
		if err != nil {
			glog.Error("marshal msgId:", msgId, " err:", err)
			return nil
		}
		msg.MsgBody = d
	}

	res, err := proto.Marshal(msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	glog.V(2).Info("向客户端发送msgId:", util.GetMsgIdName(msgId), " 长度:", len(res))

	return res
}

func BuildClientMsg2(msgId int32) []byte {
	if _, ok := pb.MessageId_name[msgId]; !ok {
		glog.Warning("build client msg failed client msgId:", msgId, " does not exist")
		return nil
	}

	msg := &pb.ClientMsg{}
	msg.MsgId = proto.Int32(int32(msgId))

	res, err := proto.Marshal(msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	glog.V(1).Info("向客户端发送msgId:", util.GetMsgIdName(msgId), " 长度:", len(res))

	return res
}

func BuildClientMsg3(msgId int32, body []byte) []byte {
	if _, ok := pb.MessageId_name[msgId]; !ok {
		glog.Warning("build client msg failed client msgId:", msgId, " does not exist")
		return nil
	}

	msg := &pb.ClientMsg{}
	msg.MsgId = proto.Int32(int32(msgId))
	msg.MsgBody = body

	res, err := proto.Marshal(msg)
	if err != nil {
		glog.Error(err)
		return nil
	}

	glog.V(1).Info("向客户端发送msgId:", util.GetMsgIdName(msgId), " 长度:", len(res))

	return res
}
