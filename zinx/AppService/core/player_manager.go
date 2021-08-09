package core

import (
	"fmt"
	"sync"

	// "time"

	"zinx/AppService/pb"
	// "zinx/utils"
	"zinx/znet"

	//"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/proto"
)

type PlayerManager struct {
	Players map[string]*ChatPlayer
	pLock   sync.RWMutex

	CID2PID map[uint32]string
	c2pLock sync.RWMutex

	Svr znet.Server
}

func NewPlayerManager(svr znet.Server) (pm *PlayerManager) {
	PMgr := &PlayerManager{
		Players: make(map[string]*ChatPlayer),
		CID2PID: make(map[uint32]string),
	}

	PMgr.Svr = svr

	return PMgr
}

func init() {
}

func (pm *PlayerManager) AddPlayer(player *ChatPlayer) {
	bNotFound := false

	pm.pLock.Lock()

	if _, ok := pm.Players[player.PID]; !ok {
		pm.Players[player.PID] = player
		bNotFound = true
	}

	pm.pLock.Unlock()

	if bNotFound {
		pm.NotifyPlayerJoin(player)
		pm.DoPlayerStateForPlayerResMsg(player.PID, player.ConnID, 0, 0)
	}
}

func (pm *PlayerManager) RemovePlayerByPID(pid string) {

	bExist := false
	// var cid uint32 = 0

	pm.pLock.Lock()

	if _, ok := pm.Players[pid]; ok {
		// cid = p.ConnID

		delete(pm.Players, pid)

		fmt.Println("Player remove from PlayerManager successfully: PID = ", pid, ",number of player = ", len(pm.Players))

		bExist = true
	}

	pm.pLock.Unlock()

	if bExist {
		// pm.Svr.GetConnMgr().Remove(cid)
		pm.NotifyPlayerLeave(pid, 0)
	}
}

func (pm *PlayerManager) PutPCIDInStaleQueue(pid string, cid uint32) {

	pm.c2pLock.Lock()
	pm.CID2PID[cid] = pid
	pm.c2pLock.Unlock()

	pm.Svr.GetConnMgr().Remove(cid)
}

func (pm *PlayerManager) DoPlayerStateForPlayerResMsg(pid string, cid uint32, tp int32, code int32) (err error) {

	msg := &pb.PlayerResMsg{
		PID: pid,
		Data: &pb.PlayerResMsg_PStatus{
			PStatus: &pb.PlayerState{
				PID:  pid,
				Type: tp,
				Code: code,
			},
		},
	}

	err = pm.sendMsgBy(cid, 99, msg)
	return err
}

func (pm *PlayerManager) FetchPlayerInfo(pid string, cid uint32) (err error) {
	var pi *pb.PlayerInfo

	pm.pLock.RLock()

	if p, ok := pm.Players[pid]; ok {
		pi = &pb.PlayerInfo{
			PID:  p.PID,
			Name: p.Name,
			Type: p.Type,
		}
	}

	pm.pLock.RUnlock()

	if pi != nil {
		msg := &pb.PlayerResMsg{
			PID: pid,
			Data: &pb.PlayerResMsg_PInfo{
				PInfo: pi,
			},
		}

		err = pm.sendMsgBy(cid, 99, msg)
		return
	}

	return nil
}

func (pm *PlayerManager) FetchPlayerList(pid string, cid uint32) (err error) {
	var pilists []*pb.PlayerList

	pm.pLock.RLock()

	pilist := &pb.PlayerList{}

	ind := 0
	sz := len(pm.Players)
	for _, player := range pm.Players {
		pi := &pb.PlayerInfo{
			PID:  player.PID,
			Name: player.Name,
			Type: player.Type,
		}

		pilist.Players = append(pilist.Players, pi)
		ind++

		if len(pilist.Players) >= 100 || ind == sz {
			pilists = append(pilists, pilist)
			pilist = &pb.PlayerList{}
		}
	}

	pm.pLock.RUnlock()

	fmt.Println("PlayerManager::FetchPlayerList::Players: size = ", sz)

	for _, plst := range pilists {
		msg := &pb.PlayerResMsg{
			PID: pid,
			Data: &pb.PlayerResMsg_PList{
				PList: plst,
			},
		}

		err = pm.sendMsgBy(cid, 99, msg)
		if err != nil {
			return
		}
	}

	return nil
}

func (pm *PlayerManager) NotifyPlayerJoin(p *ChatPlayer) (err error) {
	msg := &pb.PlayerResMsg{
		PID: p.PID,
		Data: &pb.PlayerResMsg_PJoin{
			PJoin: &pb.PlayerJoin{
				PInfo: &pb.PlayerInfo{
					PID:  p.PID,
					Name: p.Name,
					Type: p.Type,
				},
			},
		},
	}

	err = pm.broadCastByWithMe(99, msg)
	return err
}

func (pm *PlayerManager) NotifyPlayerLeave(pid string, Reason int32) (err error) {
	msg := &pb.PlayerResMsg{
		PID: pid,
		Data: &pb.PlayerResMsg_PLeave{
			PLeave: &pb.PlayerLeave{
				PID:    pid,
				Reason: Reason,
			},
		},
	}

	err = pm.broadCastBy(pid, 99, msg)
	return err
}

func (pm *PlayerManager) Talk(pid string, content string) (err error) {
	msg := &pb.TalkMsg{
		PID:     pid,
		Content: content,
	}

	err = pm.broadCastByWithMe(1, msg)
	return err
}

func (pm *PlayerManager) SysTalk(pid string, content string) (err error) {
	msg := &pb.SysMsg{
		PID:     pid,
		Content: content,
	}

	err = pm.broadCastByWithMe(3, msg)
	return err
}

func (pm *PlayerManager) sendMsg(cid uint32, msgID uint32, data []byte) (err error) {
	if err := pm.Svr.ConnMgr.SendMsg(cid, msgID, data); err != nil {
		//fmt.Println("Player SendMsg error !")
		return err
	}

	return nil
}

func (pm *PlayerManager) sendMsgBy(cid uint32, msgID uint32, msg proto.Message) (err error) {
	data, err := proto.Marshal(msg)
	if err != nil {
		fmt.Println("marshal msg err: ", err)
		return err
	}

	return pm.sendMsg(cid, msgID, data)
}

func (pm *PlayerManager) broadCast(pid string, msgID uint32, data []byte) (err error) {
	pcids := make(map[uint32]string)

	pm.pLock.RLock()
	for _, player := range pm.Players {
		if pid != player.PID {
			pcids[player.ConnID] = player.PID
		}
	}
	pm.pLock.RUnlock()

	for cid, pid2 := range pcids {
		err = pm.sendMsg(cid, msgID, data)
		if err != nil {
			pm.PutPCIDInStaleQueue(pid2, cid)

			fmt.Println("broadCast::SendMsg error: cid = ", cid, ",msgID = ", msgID, ", error = ", err)
		}
	}

	return nil
}

func (pm *PlayerManager) broadCastBy(pid string, msgID uint32, msg proto.Message) (err error) {
	data, err := proto.Marshal(msg)
	if err != nil {
		fmt.Println("marshal msg err: ", err)
		return err
	}

	err = pm.broadCast(pid, msgID, data)

	return err
}

func (pm *PlayerManager) broadCastWithMe(msgID uint32, data []byte) (err error) {

	pcids := make(map[uint32]string)

	pm.pLock.RLock()
	for _, player := range pm.Players {
		pcids[player.ConnID] = player.PID
	}
	pm.pLock.RUnlock()

	for cid, pid := range pcids {
		err = pm.sendMsg(cid, msgID, data)
		if err != nil {
			pm.PutPCIDInStaleQueue(pid, cid)

			fmt.Println("BroadCastWithMe::SendMsg error: cid = ", cid, ",msgID = ", msgID, ", error = ", err)
		}
	}

	return nil
}

func (pm *PlayerManager) broadCastByWithMe(msgID uint32, msg proto.Message) (err error) {
	data, err := proto.Marshal(msg)
	if err != nil {
		fmt.Println("marshal msg err: ", err)
		return err
	}

	err = pm.broadCastWithMe(msgID, data)

	return err
}
