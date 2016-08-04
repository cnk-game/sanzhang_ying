package GTTcp

import(  
    "net"
    "GT/GTGlobal"
)

type GTTcpSock struct {
    Port      string
}

var GTTcp *GTTcpSock

func Init() {
    GTTcp = new(GTTcpSock)
    GTTcp.Port = "8887"
    GTTcp.Start()
}

func (this *GTTcpSock)Start() {
    service := ":" + this.Port

    tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
    if err != nil {
        GTGlobal.GTLog().Error("ResolveTCPAddr error. %s", err.Error())
    }

    listener, err := net.ListenTCP("tcp", tcpAddr)
    if err != nil {
        GTGlobal.GTLog().Error("ListenTCP error. %s", err.Error())
    }

    for {
        conn, err := listener.AcceptTCP() // type is *net.TCPConn
        if err != nil {
            GTGlobal.GTLog().Error("Accept error. %s", err.Error())
            continue
        }
        go this.handle(conn)
    }
} 

func (this *GTTcpSock)handle(conn *net.TCPConn) {
}