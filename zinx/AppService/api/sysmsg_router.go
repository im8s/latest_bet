package api

import (
	"fmt"

	"zinx/AppService/pb"
	// "github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/proto"

	"zinx/ziface"
	"zinx/znet"
)

type SysMsgRouter struct {
	znet.BaseRouter
}

func (*SysMsgRouter) Handle(req ziface.IRequest) {
	msg := &pb.SysMsg{}
	err := proto.Unmarshal(req.GetData(), msg)
	if err != nil {
		fmt.Println("SysMsgRouter Unmarshal error ", err)
		return
	}

	//str := msg.Content + "[" + fmt.Sprintf("%d", sn) + "]"

	id, err := req.GetConnection().GetProperty("PID")
	if err != nil {
		fmt.Println("request.GetConnection().GetProperty error: ", err)
		return
	}

	pid := id.(string)

	GetPMgr().SysTalk(pid, msg.Content)
}
