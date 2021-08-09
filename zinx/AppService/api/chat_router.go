package api

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"zinx/AppService/pb"

	// "github.com/golang/protobuf/proto"
	// "google.com/golang/protobuf/proto"

	// "zinx/AppService/core"
	"zinx/ziface"
	"zinx/znet"
)

type ChatRouter struct {
	znet.BaseRouter
}

func (*ChatRouter) Handle(req ziface.IRequest) {

	msg := &pb.TalkMsg{}
	err := proto.Unmarshal(req.GetData(), msg)
	if err != nil {
		fmt.Println("Talk Unmarshal error ", err)
		return
	}

	id, err := req.GetConnection().GetProperty("PID")
	if err != nil {
		fmt.Println("request.GetConnection().GetProperty error: ", err)
		return
	}

	pid := id.(string)

	GetPMgr().Talk(pid, msg.Content)

	// fmt.Println("ChatRouter::Handle msg.Content = ", msg.Content)
}
