package api

import (
	// "zinx/utils"
	// "zinx/AppService/core"

	"zinx/AppService/core"
	"zinx/utils"
	"zinx/ziface"
	"zinx/znet"
)

type AppServer struct {
	znet.Server

	PMgr *core.PlayerManager
}

func NewAppServer(tpNum int) (as *AppServer) {
	ms := &AppServer{
		Server: znet.Server{
			Name:      utils.GlobalObject.Name,
			IPVersion: "tcp",
			IP:        utils.GlobalObject.Host,
			Port:      utils.GlobalObject.TCPPort,
			ConnMgr:   znet.NewConnManager(tpNum),
		},
	}

	ms.PMgr = core.NewPlayerManager(ms.Server)

	return ms
}

var AS *AppServer

func GetAppServer() (as *AppServer) {
	return AS
}

func GetPMgr() (pmgr *core.PlayerManager) {
	return AS.PMgr
}

func GetConnMgr() ziface.IConnManager {
	return AS.GetConnMgr()
}

func RunAppServer() {
	AS.Serve()
}

func init() {

	AS = NewAppServer(3)

	AS.AddRouter(0, 0, &PlayerRouter{})
	// AS.AddRouter(0, 1, &ChatRouter{})
	// AS.AddRouter(0, 3, &SysMsgRouter{})
	// AS.AddRouter(0, 19, &LogoutRouter{})

	// AS.AddRouter(1, 0, &PlayerRouter{})
	AS.AddRouter(1, 1, &ChatRouter{})
	AS.AddRouter(1, 3, &SysMsgRouter{})
	// AS.AddRouter(1, 19, &LogoutRouter{})

	// AS.AddRouter(2, 0, &PlayerRouter{})
	// AS.AddRouter(2, 1, &ChatRouter{})
	// AS.AddRouter(2, 3, &SysMsgRouter{})
	AS.AddRouter(2, 19, &LogoutRouter{})
}

/////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////
