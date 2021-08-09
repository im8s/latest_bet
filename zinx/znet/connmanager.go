package znet

import (
	"errors"
	"fmt"
	"sync"

	"zinx/utils"
	"zinx/ziface"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection
	connLock    sync.RWMutex

	WorkerPoolSize uint32

	MsgWP []ziface.IWorkerPool

	OnConnStart func(connID uint32)
	OnConnStop  func(connID uint32)
}

func NewConnManager(tpNum int) (cm *ConnManager) {
	cm = &ConnManager{
		connections:    make(map[uint32]ziface.IConnection),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		// MsgWP:          NewWorkerPool(utils.GlobalObject.WorkerPoolSize),
	}

	for k := 0; k < tpNum; k++ {
		wp := NewWorkerPool(utils.GlobalObject.WorkerPoolSize)
		cm.MsgWP = append(cm.MsgWP, wp)

		cm.MsgWP[k].Start()
	}

	return
}

// func (cm *ConnManager) GetRouterHandle() ziface.IRouterHandle {
// 	return cm.MsgRH
// }

func (cm *ConnManager) GetWorkerPool(index int) ziface.IWorkerPool {
	if index >= 0 && index < len(cm.MsgWP) {
		return cm.MsgWP[index]
	}

	return nil
}

func (cm *ConnManager) Add(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	connID := conn.GetConnID()

	if _, ok := cm.connections[connID]; !ok {
		cm.connections[connID] = conn

		fmt.Println("connection add to ConnManager successfully: ConnID = ", connID, ", conn num = ", cm.Len())
	}
}

func (cm *ConnManager) Remove(connID uint32) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	if cnn, ok := cm.connections[connID]; ok {
		cnn.Stop()

		delete(cm.connections, connID)

		fmt.Println("connection Remove from ConnManager successfully: ConnID = ", connID, ", conn num = ", cm.Len())
	}
}

func (cm *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	}

	return nil, errors.New("connection not found")
}

func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

func (cm *ConnManager) ClearConn() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	for connID, conn := range cm.connections {
		conn.Stop()
		delete(cm.connections, connID)
	}

	fmt.Println("Clear All Connections successfully: conn num = ", cm.Len())
}

func (cm *ConnManager) ClearOneConn(connID uint32) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	if conn, ok := cm.connections[connID]; ok {
		conn.Stop()
		delete(cm.connections, connID)
		fmt.Println("Clear Connections ID:  ", connID, "succeed")
		return
	}

	fmt.Println("Clear Connections ID:  ", connID, "err")
	return
}

func (cm *ConnManager) SendMsg(connID uint32, msgID uint32, data []byte) (err error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	if conn, ok := cm.connections[connID]; ok {
		if err := conn.SendMsg(msgID, data); err != nil {
			//fmt.Println("Player SendMsg error !")
			return err
		}
	}

	// fmt.Println("After ConnManager::SendMsg: connID = ", connID)

	return nil
}

func (cm *ConnManager) BroadCast(msgID uint32, data []byte) (err error) {

	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	for _, conn := range cm.connections {
		if err = conn.SendMsg(msgID, data); err != nil {
			fmt.Println("ConnManager::BroadCast:conn.SendMsg error: ", err)
		}
	}

	return nil
}
