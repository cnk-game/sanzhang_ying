package server

import (
	"code.google.com/p/go.net/websocket"
	"code.google.com/p/goprotobuf/proto"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"net/http"
	_ "net/http/pprof"
	"os"
	"pb"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type MsgDispatcher interface {
	RegisterHandlers(r *mux.Router)
	DispatchMsg(msg *pb.ServerMsg, sess *Session) []byte
}

type GameServer struct {
	dispatcher    MsgDispatcher
	sigChan       chan os.Signal
	waitGroup     *sync.WaitGroup
	stopChan      chan bool
	stopOnce      sync.Once
	refuseService int32
}

var s *GameServer

var bindHost string

func init() {
	flag.StringVar(&bindHost, "bindHost", ":8002", "bind server host.")

	s = &GameServer{}
	s.sigChan = make(chan os.Signal, 1)
	s.waitGroup = &sync.WaitGroup{}
	s.stopChan = make(chan bool)
}

func GetServerInstance() *GameServer {
	return s
}

func (s *GameServer) GetSigChan() chan os.Signal {
	return s.sigChan
}

func (s *GameServer) StartServer(dispatcher MsgDispatcher) {
	s.dispatcher = dispatcher

	r := mux.NewRouter()
	http.Handle("/", r)
	http.Handle("/ws/", websocket.Server{Handler: s.handleClient, Handshake: nil})

	s.dispatcher.RegisterHandlers(r)

	glog.Info("===>启动Game服务", bindHost)

	glog.Fatal(http.ListenAndServe(fmt.Sprintf("%v", bindHost), nil))
}

func (s *GameServer) StopServer() {
	s.stopOnce.Do(func() {
		go func() {
			s.sigChan <- syscall.SIGKILL
		}()
	})
}

func (s *GameServer) WaitStopServer() {
	glog.Info("==>Start WaitStopServer")
	defer glog.Info("==>WaitStopServer done.")

	close(s.stopChan)
	s.waitGroup.Wait()
}

func (s *GameServer) IsRefuseService() bool {
	return atomic.AddInt32(&s.refuseService, 0) > 0
}

func (s *GameServer) SetRefuseService() {
	atomic.AddInt32(&s.refuseService, 1)
}

func (s *GameServer) handleClient(conn *websocket.Conn) {
	if s.IsRefuseService() {
		glog.Info("==>正在停止服务，拒绝连接...")
		conn.Close()
		return
	}

	sess := newSess(conn)
	go sess.Run(s.dispatcher)

	defer sess.cleanSess()

	for {
		var data []byte
		conn.SetReadDeadline(time.Now().Add(time.Minute * 10))
		err := websocket.Message.Receive(conn, &data)
		if err != nil {
			glog.Info("error receiving msg:", err)
			break
		}

		conn.SetReadDeadline(time.Time{})

		clientMsg := &pb.ClientMsg{}
		err = proto.Unmarshal(data, clientMsg)
		if err != nil {
			glog.V(1).Info("unmarshal client msg failed!")
			break
		}

		msg := &pb.ServerMsg{}
		msg.Client = proto.Bool(true)
		msg.MsgId = clientMsg.MsgId
		msg.MsgBody = clientMsg.MsgBody

		sess.mq <- msg
	}

}
