package ws

import (
	"sync"
)

// Room 房间
type Room struct {
	rwMutex sync.RWMutex
	roomId  string
	id2Conn map[uint64]*Connection
}

func InitRoom(roomId string) (room *Room) {
	room = &Room{
		roomId:  roomId,
		id2Conn: make(map[uint64]*Connection),
	}
	return
}

func (room *Room) Join(wsConn *Connection) (err error) {
	var (
		existed bool
	)

	room.rwMutex.Lock()
	defer room.rwMutex.Unlock()

	if _, existed = room.id2Conn[wsConn.connId]; existed {
		err = ErrJoinRoomTwice
		return
	}

	room.id2Conn[wsConn.connId] = wsConn
	return
}

func (room *Room) Leave(wsConn *Connection) (err error) {
	var (
		existed bool
	)

	room.rwMutex.Lock()
	defer room.rwMutex.Unlock()

	if _, existed = room.id2Conn[wsConn.connId]; !existed {
		err = ErrNotInRoom
		return
	}

	delete(room.id2Conn, wsConn.connId)
	return
}

func (room *Room) Count() int {
	room.rwMutex.RLock()
	defer room.rwMutex.RUnlock()

	return len(room.id2Conn)
}

func (room *Room) Push(wsMsg *Message) {
	var (
		wsConn *Connection
	)
	room.rwMutex.RLock()
	defer room.rwMutex.RUnlock()

	for _, wsConn = range room.id2Conn {
		err := wsConn.SendMessage(wsMsg)
		if err != nil {
			//return
		}
	}
}
