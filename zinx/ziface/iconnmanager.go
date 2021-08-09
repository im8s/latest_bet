package ziface

type IConnManager interface {
	// GetRouterHandle() IRouterHandle

	GetWorkerPool(index int) IWorkerPool

	Add(conn IConnection)
	Remove(connID uint32)
	Get(connID uint32) (IConnection, error)
	Len() int
	ClearConn()

	SendMsg(connID uint32, msgID uint32, data []byte) (err error)
	BroadCast(msgID uint32, data []byte) (err error)

	// SetOnConnStart(hookFunc func(connID uint32))
	// SetOnConnStop(hookFunc func(connID uint32))

	// CallOnConnStart(connID uint32)
	// CallOnConnStop(connID uint32)
}
