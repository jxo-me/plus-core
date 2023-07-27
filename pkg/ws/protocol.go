package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

// 推送类型
const (
	ActionPost   = "post"
	ActionGet    = "get"
	ActionDelete = "delete"
	ActionUpdate = "update"
	ActionJoin   = "join"
	ActionLeave  = "leave"
	PushTypeUser = 0 // 推送房间
	PushTypeRoom = 1 // 推送房间
	PushTypeAll  = 2 // 推送在线
)

// Message websocket的Message对象
type Message struct {
	MsgType int
	MsgData []byte
}

// BizMessage 业务消息的固定格式(type+data)
type BizMessage struct {
	Type string          `json:"type"` // type消息类型: PING, PONG, JOIN, LEAVE, PUSH
	Data json.RawMessage `json:"data"` // data数据字段
}

// Data数据类型

// BizPushData PUSH
type BizPushData struct {
	Items []*json.RawMessage `json:"items"`
}

// PingData PING
type PingData struct{}

// PongData PONG
type PongData struct {
	SessionId uint64 `json:"session_id"`
	KeepAlive int    `json:"keep_alive"`
	DateTime  string `json:"date_time"`
	Expire    int64  `json:"expire"`
}

// JoinData JOIN
type JoinData struct {
	Room string `json:"room"`
}

// LeaveData LEAVE
type LeaveData struct {
	Room string `json:"room"`
}

func BuildWsMessage(msgType int, msgData []byte) (wsMessage *Message) {
	return &Message{
		MsgType: msgType,
		MsgData: msgData,
	}
}

func EncodeWsMessage(message *MessageReq) (wsMessage *Message, err error) {
	var (
		buf []byte
	)
	if buf, err = json.Marshal(Response{Code: 0, Message: "ok", Body: *message}); err != nil {
		return
	}
	wsMessage = &Message{websocket.TextMessage, buf}
	return
}

type MessageReq struct {
	Service string      `json:"service"`
	Action  string      `json:"action"`
	Data    interface{} `json:"data"`
}

type MessageRes struct {
	Service string      `json:"service"`
	Action  string      `json:"action"`
	Data    interface{} `json:"data"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Body    interface{} `json:"body"`
}

type NullResp struct{}

// DecodeMessage 解析{"type": "PING", "data": {...}}的包
func DecodeMessage(buf []byte) (message *MessageReq, err error) {
	var (
		msgObj MessageReq
	)

	if err = json.Unmarshal(buf, &msgObj); err != nil {
		return nil, err
	}

	message = &msgObj
	return
}
