package znet

import (
	"fmt"
	"net"

	"zinx/utils"
	"zinx/ziface"
	// "github.com/golang/protobuf/descriptor"
)

var zinxLogo = `                                        
              ██                        
              ▀▀                        
 ████████   ████     ██▄████▄  ▀██  ██▀ 
     ▄█▀      ██     ██▀   ██    ████   
   ▄█▀        ██     ██    ██    ▄██▄   
 ▄██▄▄▄▄▄  ▄▄▄██▄▄▄  ██    ██   ▄█▀▀█▄  
 ▀▀▀▀▀▀▀▀  ▀▀▀▀▀▀▀▀  ▀▀    ▀▀  ▀▀▀  ▀▀▀ 
                                        `
var topLine = `┌───────────────────────────────────────────────────┐`
var borderLine = `│`
var bottomLine = `└───────────────────────────────────────────────────┘`

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int

	ConnMgr ziface.IConnManager
}

func NewServer(tpNum int) (svr *Server) {
	printLogo()

	svr = &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TCPPort,
		ConnMgr:   NewConnManager(tpNum),
	}

	return svr
}

func (s *Server) Start() {
	fmt.Printf("[START] Server name: %s,listenner at IP: %s, Port %d is starting\n", s.Name, s.IP, s.Port)

	go func() {
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err: ", err)
			return
		}

		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}

		fmt.Println("start Zinx server  ", s.Name, " succ, now listenning...")

		var cID uint32
		cID = 0

		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err ", err)
				continue
			}
			fmt.Println("Get conn remote addr = ", conn.RemoteAddr().String())

			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				conn.Close()
				continue
			}

			dealConn := NewConntion(s.ConnMgr, conn, cID)
			cID++

			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] Zinx server , name ", s.Name)

	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	s.Start()

	select {}
}

func (s *Server) AddRouter(index int, msgID uint32, router ziface.IRouter) {
	s.ConnMgr.GetWorkerPool(index).GetRouterHandle().AddRouter(msgID, router)
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

func printLogo() {
	fmt.Println(zinxLogo)
	fmt.Println(topLine)
	// fmt.Println(fmt.Sprintf("%s [Github] https://github.com/aceld                 %s", borderLine, borderLine))
	// fmt.Println(fmt.Sprintf("%s [tutorial] https://www.kancloud.cn/aceld/zinx     %s", borderLine, borderLine))
	fmt.Println(bottomLine)
	fmt.Printf("[AppService] Version: %s, MaxConn: %d, MaxPacketSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)
}

func init() {
}
