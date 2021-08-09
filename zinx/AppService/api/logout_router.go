package api

import (

	// "github.com/golang/protobuf/proto"

	// "zinx/AppService/core"
	// "zinx/AppService/pb"
	"zinx/ziface"
	"zinx/znet"

	"fmt"
	// "google.golang.org/protobuf/proto"
)

type LogoutRouter struct {
	znet.BaseRouter
}

func (*LogoutRouter) Handle(req ziface.IRequest) {
	id, err := req.GetConnection().GetProperty("PID")
	if err != nil {
		fmt.Println("req.GetConnection().GetProperty error: ", err)
		return
	}

	pid := id.(string)
	// cid := req.GetConnection().GetConnID()

	// fmt.Println("LogoutRouter::Handle: pid = ", pid)

	// GetConnMgr().Remove(cid)
	GetPMgr().RemovePlayerByPID(pid)
}
