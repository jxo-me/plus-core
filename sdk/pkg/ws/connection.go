package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

type Connection struct {
	mutex             sync.Mutex
	connId            uint64
	wsSocket          *ghttp.WebSocket
	inChan            chan *Message
	outChan           chan *Message
	closeChan         chan byte
	isClosed          bool
	lastHeartbeatTime time.Time       // 最近一次心跳时间
	rooms             map[string]bool // 加入了哪些房间
	heartbeatInterval int
	inChannelSize     int
	outChannelSize    int
}

func NewConnection(
	ctx context.Context, connId uint64,
	wsSocket *ghttp.WebSocket,
	heartbeat, inChannelSize, outChannelSize int,
) (c *Connection) {
	c = &Connection{
		wsSocket:          wsSocket,
		connId:            connId,
		inChan:            make(chan *Message, inChannelSize),
		outChan:           make(chan *Message, outChannelSize),
		closeChan:         make(chan byte),
		lastHeartbeatTime: time.Now(),
		rooms:             make(map[string]bool),
		heartbeatInterval: heartbeat,
		inChannelSize:     inChannelSize,
		outChannelSize:    outChannelSize,
	}

	wsSocket.SetCloseHandler(func(code int, text string) error {
		glog.Debug(ctx, fmt.Sprintf("websocket Connection close code: %d, message: %s", code, text))
		c.Close(ctx)
		return nil
	})
	go c.readLoop(ctx)
	go c.writeLoop(ctx)

	return c
}

// 读websocket
func (conn *Connection) readLoop(ctx context.Context) {
	var (
		msgType int
		msgData []byte
		message *Message
		err     error
	)

	for {
		if msgType, msgData, err = conn.wsSocket.ReadMessage(); err != nil {
			goto ERR
		}

		message = BuildWsMessage(msgType, msgData)
		select {
		case conn.inChan <- message:
		case <-conn.closeChan:
			goto CLOSED
		}
	}

ERR:
	conn.Close(ctx)
CLOSED:
}

// 写websocket
func (conn *Connection) writeLoop(ctx context.Context) {
	var (
		message *Message
		err     error
	)
	for {
		select {
		case message = <-conn.outChan:
			if err = conn.wsSocket.WriteMessage(message.MsgType, message.MsgData); err != nil {
				goto ERR
			}
		case <-conn.closeChan:
			goto CLOSED
		}
	}
ERR:
	conn.Close(ctx)
CLOSED:
}

// SendMessage 发送消息
func (conn *Connection) SendMessage(message *Message) (err error) {
	select {
	case conn.outChan <- message:
		// 统计计数
		GetStats().SendMessageTotalIncr()
	case <-conn.closeChan:
		err = ErrConnectionLoss
	default: // 写操作不会阻塞, 因为channel已经预留给websocket一定的缓冲空间
		err = ErrSendMessageFull
		// 统计计数
		GetStats().SendMessageFailIncr()
	}
	return
}

// ReadMessage 读取消息
func (conn *Connection) ReadMessage() (message *Message, err error) {
	select {
	case message = <-conn.inChan:
	case <-conn.closeChan:
		err = ErrConnectionLoss
	}
	return
}

// Close 关闭连接
func (conn *Connection) Close(context.Context) {
	_ = conn.wsSocket.Close()

	conn.mutex.Lock()
	defer conn.mutex.Unlock()

	if !conn.isClosed {
		conn.isClosed = true
		close(conn.closeChan)
	}
}

// IsAlive 检查心跳（不需要太频繁）
func (conn *Connection) IsAlive() bool {
	var (
		now = time.Now()
	)

	conn.mutex.Lock()
	defer conn.mutex.Unlock()

	// 连接已关闭 或者 太久没有心跳
	if conn.isClosed || now.Sub(conn.lastHeartbeatTime) > time.Duration(conn.heartbeatInterval)*time.Second {
		return false
	}
	return true
}

// KeepAlive 更新心跳
func (conn *Connection) KeepAlive() {
	var (
		now = time.Now()
	)

	conn.mutex.Lock()
	defer conn.mutex.Unlock()

	conn.lastHeartbeatTime = now
}

// heartbeatChecker 每隔1秒, 检查一次连接是否健康
func (conn *Connection) heartbeatChecker(ctx context.Context) {
	var (
		timer *time.Timer
	)
	timer = time.NewTimer(time.Duration(conn.heartbeatInterval) * time.Second)
	for {
		select {
		case <-timer.C:
			if !conn.IsAlive() {
				conn.Close(ctx)
				// 更新统计

				goto EXIT
			}
			timer.Reset(time.Duration(conn.heartbeatInterval) * time.Second)
		case <-conn.closeChan:
			timer.Stop()
			goto EXIT
		}
	}

EXIT:
	// 确保连接被关闭
}

// 处理PING请求
func (conn *Connection) handlePing() (resp *MessageRes, err error) {
	conn.KeepAlive()
	t := gtime.Now().Add(time.Duration(conn.heartbeatInterval) * time.Second)

	resp = &MessageRes{
		Service: "ping",
		Action:  "pong",
		Data: PongData{
			SessionId: conn.connId,
			KeepAlive: conn.heartbeatInterval,
			Expire:    t.Unix(),
			DateTime:  t.Format("Y-m-d H:i:s"),
		},
	}
	return resp, nil
}

// 处理JOIN请求
func (conn *Connection) handleJoin(connMgr *ConnManager, maxJoinRoom int, req *MessageReq) (resp *MessageRes, err error) {
	var (
		joinData JoinData
		existed  bool
	)
	err = gconv.Struct(req.Data, &joinData)
	if err != nil {
		return nil, err
	}
	if len(joinData.Room) == 0 {
		return nil, ErrRoomIdInvalid
	}
	if len(conn.rooms) >= maxJoinRoom {
		// 房间超过了数量限制, 忽略这个请求
		return
	}
	// 已加入过
	if _, existed = conn.rooms[joinData.Room]; existed {
		// 忽略掉这个请求
		return
	}
	// 建立房间 -> 连接的关系
	if err = connMgr.JoinRoom(joinData.Room, conn); err != nil {
		return
	}
	// 建立连接 -> 房间的关系
	conn.rooms[joinData.Room] = true
	return
}

// 处理LEAVE请求
func (conn *Connection) handleLeave(connMgr *ConnManager, req *MessageReq) (resp *MessageRes, err error) {
	var (
		leaveData LeaveData
		existed   bool
	)
	err = gconv.Struct(req.Data, &leaveData)
	if err != nil {
		return nil, err
	}
	if len(leaveData.Room) == 0 {
		err = ErrRoomIdInvalid
		return
	}
	// 未加入过
	if _, existed = conn.rooms[leaveData.Room]; !existed {
		// 忽略掉这个请求
		return
	}
	// 删除房间 -> 连接的关系
	if err = connMgr.LeaveRoom(leaveData.Room, conn); err != nil {
		return
	}
	// 删除连接 -> 房间的关系
	delete(conn.rooms, leaveData.Room)
	return
}

func (conn *Connection) leaveAll(connMgr *ConnManager, ctx context.Context) {
	var (
		roomId string
	)
	// 从所有房间中退出
	for roomId = range conn.rooms {
		err := connMgr.LeaveRoom(roomId, conn)
		if err != nil {
			glog.Warning(ctx, "connMgr LeaveRoom error:", err.Error())
		}
		delete(conn.rooms, roomId)
	}
}

// RouterHandle 处理websocket请求
func (conn *Connection) RouterHandle(ctx context.Context, ins *Instance, routers *map[string]Service) {
	var (
		message *Message
		req     *MessageReq
		res     *MessageRes
		resp    *Response
		err     error
		buf     []byte
	)

	// 心跳检测线程
	go conn.heartbeatChecker(ctx)

	// 请求处理协程
	for {
		if message, err = conn.ReadMessage(); err != nil {
			goto ERR
		}

		// 只处理文本消息
		if message.MsgType != websocket.TextMessage {
			continue
		}

		// 解析消息体
		if req, err = DecodeMessage(message.MsgData); err != nil {
			goto ERR
		}

		// 1,收到PING则响应PONG: {"action": "ping"}, {"type": "pong"}
		// 2,收到JOIN则加入ROOM: {"action": "join", "data": {"room": "chrome-plugin"}}
		// 3,收到LEAVE则离开ROOM: {"action": "leave", "data": {"room": "chrome-plugin"}}
		//fmt.Println("收到新消息:", req.Data)
		switch req.Service {
		case "ping":
			res, err = conn.handlePing()
			//if err = conn.SendMessage(&Message{websocket.PongMessage, buf}); err != nil {
			//	if err != ErrSendMessageFull {
			//		goto ERR
			//	} else {
			//		err = nil
			//	}
			//}
			//continue
		//case "join":
		//	res, err = conn.handleJoin(ins.ConnMgr(), 100, req)
		//case "leave":
		//	res, err = conn.handleLeave(ins.ConnMgr(), req)
		default:
			res, err = conn.dispatcher(ctx, req, routers)
		}
		if err != nil {
			resp = &Response{Code: http.StatusInternalServerError, Message: err.Error(), Body: NullResp{}}
		} else {
			resp = &Response{Code: 0, Message: "ok", Body: res}
		}
		if buf, err = json.Marshal(resp); err != nil {
			goto ERR
		}
		// socket缓冲区写满不是致命错误
		if err = conn.SendMessage(&Message{websocket.TextMessage, buf}); err != nil {
			if err != ErrSendMessageFull {
				goto ERR
			} else {
				err = nil
			}
		}
	}

ERR:
	// 确保连接关闭
	conn.Close(ctx)

	// 离开所有房间
	conn.leaveAll(ins.ConnManager, ctx)

	// 从连接池中移除
	ins.ConnManager.DelConn(conn)
	return
}

func (conn *Connection) dispatcher(ctx context.Context, msg *MessageReq, routers *map[string]Service) (*MessageRes, error) {
	var res = &MessageRes{}
	var err error
	// action = join or leave ?
	if s, y := (*routers)[msg.Service]; y {
		// validation
		// before hook
		if h, ok := s.(SrvBefore); ok {
			err = h.Before(ctx, msg)
			if err != nil {
				glog.Warning(ctx, fmt.Sprintf("%s service handle after error:", msg.Service), err)
				return nil, err
			}
		}
		if h, ok := s.(SrvHandler); ok {
			res, err = h.Handle(ctx, msg)
			if err != nil {
				glog.Warning(ctx, fmt.Sprintf("%s service handle error:", msg.Service), err)
				return nil, err
			}
		} else {
			res, err = conn.Handle(ctx, s, msg)
			if err != nil {
				return nil, err
			}
		}
		// after hook
		if h, ok := s.(SrvAfter); ok {
			defer func(r *MessageRes, srvName string) {
				if r == nil {
					r = &MessageRes{}
				}
				err = h.After(ctx, r)
				if err != nil {
					glog.Warning(ctx, fmt.Sprintf("%s service handle after error:", srvName), err)
				}
			}(res, msg.Service)
		}

		return res, nil
	}
	// not found
	return nil, gerror.Newf("Not found service: %s", msg.Service)
}

func (conn *Connection) Handle(ctx context.Context, h Service, req *MessageReq) (res *MessageRes, err error) {
	glog.Infof(ctx, "default service %s handler.. ", req.Service)
	switch req.Action {
	case ActionJoin:
		if o, ok := h.Action().(ActJoin); ok {
			res, err = o.Join(ctx, req)
		}
	case ActionLeave:
		if o, ok := h.Action().(ActLeave); ok {
			res, err = o.Leave(ctx, req)
		}
	case ActionPost:
		if o, ok := h.Action().(ActPost); ok {
			res, err = o.Post(ctx, req)
		}
	case ActionDelete:
		if o, ok := h.Action().(ActDelete); ok {
			res, err = o.Delete(ctx, req)
		}
	case ActionUpdate:
		if o, ok := h.Action().(ActUpdate); ok {
			res, err = o.Update(ctx, req)
		}
	case ActionGet:
		if o, ok := h.Action().(ActGet); ok {
			res, err = o.Get(ctx, req)
		}
	default:
		return nil, gerror.Newf("default service %s not found action: %s ", req.Service, req.Action)
	}

	if res == nil {
		res = &MessageRes{}
	}

	return res, err
}
