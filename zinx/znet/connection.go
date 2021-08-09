package znet

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"time"

	// "zinx/utils"
	"zinx/ziface"
)

type Connection struct {
	CnnMgr ziface.IConnManager
	Conn   *net.TCPConn
	ConnID uint32

	ctx    context.Context
	cancel context.CancelFunc

	msgQ [][]byte

	sync.RWMutex

	property     map[string]interface{}
	propertyLock sync.Mutex

	IsClosed bool
}

func NewConntion(cmgr ziface.IConnManager, conn *net.TCPConn, connID uint32) *Connection {
	c := &Connection{
		CnnMgr:   cmgr,
		Conn:     conn,
		ConnID:   connID,
		property: nil,

		IsClosed: false,
	}

	// if err := c.Conn.SetReadDeadline(time.Time{}); err != nil {
	// 	fmt.Println("c.Conn.SetReadDeadline error")
	// }

	c.CnnMgr.Add(c)

	return c
}

func (c *Connection) LoginoutRequest() {
	c.Stop()

	c.CnnMgr.Remove(c.ConnID)

	msg := &Message{
		DataLen: 8,
		ID:      19,
	}
	req := &Request{
		conn: c,
		msg:  msg,
	}

	// c.CnnMgr.GetWorkerPool(1).FetchRequestTo(&req)

	msgID := msg.ID

	for k := 0; k < 3; k++ {
		if c.CnnMgr.GetWorkerPool(k).GetRouterHandle().MatchRouter(msgID) {
			c.CnnMgr.GetWorkerPool(k).FetchRequestTo(req)
			break
		}
	}
}

func (c *Connection) StartWriter() {

	// defer c.LoginoutRequest()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			isQuit := false

			c.Lock()

			for _, data := range c.msgQ {
				if _, err := c.Conn.Write(data); err != nil {
					fmt.Println("[ConnID", c.ConnID, "] Send Data error:, ", err, " Conn Writer exit")
					isQuit = true
					break
				}
			}

			c.msgQ = c.msgQ[0:0]

			c.Unlock()

			if isQuit {
				return
			}

			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (c *Connection) StartReader() {

	defer c.LoginoutRequest()

	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			dp := NewDataPack()

			headData := make([]byte, dp.GetHeadLen())
			if _, err := io.ReadFull(c.Conn, headData); err != nil {
				fmt.Println("[ConnID", c.ConnID, "] read msg head error ", err)
				return
			}
			//fmt.Printf("read headData %+v\n", headData)

			msg, err := dp.Unpack(headData)
			if err != nil {
				fmt.Println("[ConnID", c.ConnID, "] unpack error ", err)
				return
			}

			var data []byte
			if msg.GetDataLen() > 0 {
				data = make([]byte, msg.GetDataLen())
				if _, err := io.ReadFull(c.Conn, data); err != nil {
					fmt.Println("[ConnID", c.ConnID, "] read msg data error ", err)
					return
				}
			}
			msg.SetData(data)

			req := &Request{
				conn: c,
				msg:  msg,
			}

			msgID := msg.GetMsgID()

			for k := 0; k < 3; k++ {
				if c.CnnMgr.GetWorkerPool(k).GetRouterHandle().MatchRouter(msgID) {
					c.CnnMgr.GetWorkerPool(k).FetchRequestTo(req)
					// fmt.Println("request msgID = ", msgID, " be dispatched to workpool: ", k)
					break
				}
			}

			// if msgID == 0 || msgID == 1 {
			// 	c.CnnMgr.GetWorkerPool(0).FetchRequestTo(req)
			// } else if msgID == 2 || msgID == 19 || msgID == 3 {
			// 	c.CnnMgr.GetWorkerPool(1).FetchRequestTo(req)
			// } else if msgID == 2 || msgID == 19 || msgID == 3 {
			// 	c.CnnMgr.GetWorkerPool(2).FetchRequestTo(req)
			// }
		}
	}
}

func (c *Connection) Start() {
	c.ctx, c.cancel = context.WithCancel(context.Background())

	go c.StartReader()
	go c.StartWriter()
}

func (c *Connection) Stop() {
	c.Lock()
	defer c.Unlock()

	if !c.IsClosed {
		c.Conn.Close()
		c.cancel()

		c.IsClosed = true

		c.msgQ = c.msgQ[0:0]
	}

	fmt.Println("Connection Stop()...ConnID = ", c.ConnID)
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(msgID uint32, data []byte) error {

	// fmt.Println("Connection::SendMsg: ConnID = ", c.ConnID)

	c.Lock()
	defer c.Unlock()

	if c.IsClosed {
		// fmt.Println("Connection::SendMsg: c.isClosed == true  ConnID = ", c.ConnID)
		return errors.New("connection closed when send msg")
	}

	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgID, data))
	if err != nil {
		fmt.Println("Connection::SendMsg: Pack error ConnID = ", c.ConnID, ", msgID = ", msgID)
		return errors.New("Pack error msg ")
	}

	c.msgQ = append(c.msgQ, msg)

	return nil
}

func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	if c.property == nil {
		c.property = make(map[string]interface{})
	}

	c.property[key] = value
}

func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	}

	return nil, errors.New("no property found")
}

func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}

func (c *Connection) Context() context.Context {
	return c.ctx
}
