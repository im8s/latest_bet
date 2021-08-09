package ziface

type IServer interface {
	Start()
	Stop()
	Serve()
	AddRouter(index int, msgID uint32, router IRouter)
	GetConnMgr() IConnManager
}
