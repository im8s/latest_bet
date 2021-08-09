package core

//"github.com/golang/protobuf/proto"

type ChatPlayer struct {
	ConnID uint32
	PID    string
	Name   string
	Type   int32

	// Conn ziface.IConnection
}

func NewChatPlayer(pid string, cid uint32, name string, tp int32) *ChatPlayer {
	p := &ChatPlayer{
		ConnID: cid,
		PID:    pid,
		Name:   name,
		Type:   tp,
	}

	return p
}
