package api

import (
	"fmt"
	// "reflect"

	"zinx/AppService/pb"

	// "github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/proto"

	"zinx/AppService/core"
	"zinx/ziface"
	"zinx/znet"
	// "google.golang.org/protobuf/proto"
	// "google.golang.org/protobuf/reflect/protoreflect"
	// "google.golang.org/protobuf/reflect/protoregistry"
)

type PlayerRouter struct {
	znet.BaseRouter
}

func (*PlayerRouter) Handle(req ziface.IRequest) {
	// type_name := "pb.LoginMsg"
	// msg, err := parseRequest(type_name, req.GetData())
	// if err != nil {
	// 	fmt.Println("parseRequest error ", err)
	// 	// return
	// } else {
	// 	lmsg := msg.(*pb.LoginMsg)
	// 	fmt.Println("msg Name = ", lmsg.Name)
	// }
	// fullName := loginMsg.ProtoReflect().Descriptor().FullName()
	// fmt.Println("fullName = ", fullName)

	msg := &pb.PlayerReqMsg{}
	err := proto.Unmarshal(req.GetData(), msg)
	if err != nil {
		fmt.Println("PlayerRouter::Handle: PlayerReqMsg Unmarshal error ", err)
		return
	}

	plogin := msg.GetPLogin()
	if plogin != nil {
		pid := plogin.Name
		cid := req.GetConnection().GetConnID()

		req.GetConnection().SetProperty("PID", pid)

		player := core.NewChatPlayer(pid, cid, plogin.Name, plogin.Type)
		GetPMgr().AddPlayer(player)

		// fmt.Println("Login success: pid = ", pid)

		return
	}

	pother := msg.GetPOther()
	if pother != nil {
		id, err := req.GetConnection().GetProperty("PID")
		if err != nil {
			fmt.Println("request.GetConnection().GetProperty error: ", err)
			return
		}
		pid := id.(string)
		cid := req.GetConnection().GetConnID()

		rtp := pother.RType
		if rtp == 2 {
			GetPMgr().FetchPlayerList(pid, cid)

			// fmt.Println("PlayerRouter fetch player list success: pid = ", pid)
		}

		return
	}
}
