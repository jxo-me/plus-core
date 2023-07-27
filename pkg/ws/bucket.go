package ws

import (
	"context"
	"errors"
	"sync"
)

type Bucket struct {
	rwMutex sync.RWMutex
	index   int                    // 我是第几个桶
	id2Conn map[uint64]*Connection // 连接列表(key=连接唯一ID)
	rooms   map[string]*Room       // 房间列表
}

func InitBucket(bucketIdx int) (bucket *Bucket) {
	bucket = &Bucket{
		index:   bucketIdx,
		id2Conn: make(map[uint64]*Connection),
		rooms:   make(map[string]*Room),
	}
	return bucket
}

func (bucket *Bucket) JoinRoom(roomId string, wsConn *Connection) (err error) {
	var (
		existed bool
		room    *Room
	)
	bucket.rwMutex.Lock()
	defer bucket.rwMutex.Unlock()

	// 找到房间
	if room, existed = bucket.rooms[roomId]; !existed {
		room = InitRoom(roomId)
		bucket.rooms[roomId] = room
		GetStats().RoomCountIncr()
	}
	// 加入房间
	err = room.Join(wsConn)
	return
}

func (bucket *Bucket) LeaveRoom(roomId string, wsConn *Connection) (err error) {
	var (
		existed bool
		room    *Room
	)
	bucket.rwMutex.Lock()
	defer bucket.rwMutex.Unlock()

	// 找到房间
	if room, existed = bucket.rooms[roomId]; !existed {
		err = ErrNotInRoom
		return
	}

	err = room.Leave(wsConn)

	// 房间为空, 则删除
	if room.Count() == 0 {
		delete(bucket.rooms, roomId)
		GetStats().RoomCountDesc()
	}
	return
}

func (bucket *Bucket) AddConn(wsConn *Connection) {
	bucket.rwMutex.Lock()
	defer bucket.rwMutex.Unlock()

	bucket.id2Conn[wsConn.connId] = wsConn
}

func (bucket *Bucket) DelConn(wsConn *Connection) {
	bucket.rwMutex.Lock()
	defer bucket.rwMutex.Unlock()

	delete(bucket.id2Conn, wsConn.connId)
}

func (bucket *Bucket) DelConnByCID(ctx context.Context, connId uint64) {
	var (
		wsConn *Connection
		ok     bool
	)
	bucket.rwMutex.Lock()
	defer bucket.rwMutex.Unlock()
	if wsConn, ok = bucket.id2Conn[connId]; ok {
		wsConn.Close(ctx)
		delete(bucket.id2Conn, wsConn.connId)
	}
}

// PushAll 推送给Bucket内所有用户
func (bucket *Bucket) PushAll(wsMsg *Message) {
	var (
		wsConn *Connection
	)

	// 锁Bucket
	bucket.rwMutex.RLock()
	defer bucket.rwMutex.RUnlock()

	// 全量非阻塞推送
	for _, wsConn = range bucket.id2Conn {
		err := wsConn.SendMessage(wsMsg)
		if err != nil {
			//return
		}
	}
}

// PushByCid 根据CID 给用户推送消息
func (bucket *Bucket) PushByCid(connId uint64, wsMsg *Message) (err error) {
	var (
		wsConn *Connection
		ok     bool
	)
	// 锁Bucket
	bucket.rwMutex.RLock()
	defer bucket.rwMutex.RUnlock()
	if wsConn, ok = bucket.id2Conn[connId]; ok {
		if err := wsConn.SendMessage(wsMsg); err != nil {
			return err
		}
		return nil
	}
	return errors.New("发送失败，connId不正确或websocket连接不存在!")
}
