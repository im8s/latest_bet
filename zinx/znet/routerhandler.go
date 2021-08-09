package znet

import (
	"fmt"
	"strconv"

	"zinx/ziface"
)

type RouterHandle struct {
	Apis map[uint32]ziface.IRouter
}

func NewRouterHandle() *RouterHandle {
	mh := &RouterHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}

	return mh
}

// func (mh *RouterHandle) FetchRequestTo(req ziface.IRequest) {
// 	mh.WP.FetchRequestTo(req)
// }

func (mh *RouterHandle) DoHandler(req ziface.IRequest) {
	handler, ok := mh.Apis[req.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", req.GetMsgID(), " is not FOUND!")
		return
	}

	handler.PreHandle(req)
	handler.Handle(req)
	handler.PostHandle(req)
}

func (mh *RouterHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeated api , msgID = " + strconv.Itoa(int(msgID)))
	}

	mh.Apis[msgID] = router
	fmt.Println("Add api msgID = ", msgID)
}

func (mh *RouterHandle) MatchRouter(msgID uint32) bool {
	if _, ok := mh.Apis[msgID]; ok {
		return true
	}

	return false
}
