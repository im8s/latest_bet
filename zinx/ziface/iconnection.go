package ziface

import (
	"context"
	"net"
)

type IConnection interface {
	Start()
	Stop()
	Context() context.Context

	GetTCPConnection() *net.TCPConn
	GetConnID() uint32
	RemoteAddr() net.Addr

	SendMsg(msgID uint32, data []byte) error
	// SendBuffMsg(msgID uint32, data []byte) error

	SetProperty(key string, value interface{})
	GetProperty(key string) (interface{}, error)
	RemoveProperty(key string)
}
