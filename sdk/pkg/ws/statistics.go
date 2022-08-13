package ws

import (
	"encoding/json"
	"sync/atomic"
)

var insStats = Statistics{}

type Statistics struct {
	// 反馈在线长连接的数量
	OnlineConnections int64 `json:"online_connections"`

	// 反馈客户端的推送压力
	SendMessageTotal int64 `json:"send_message_total"`
	SendMessageFail  int64 `json:"send_message_fail"`

	// 反馈ConnMgr消息分发模块的压力
	// 推送失败次数
	PushFail        int64 `json:"push_fail"`
	DispatchPending int64 `json:"dispatch_pending"`
	PushJobPending  int64 `json:"push_job_pending"`
	// 分发总消息数
	DispatchTotal int64 `json:"dispatch_total"`
	// 分发丢弃消息数
	DispatchFail int64 `json:"dispatch_fail"`
	// 返回出房间在线的总数, 有利于分析内存上涨的原因
	RoomCount int64 `json:"room_count"`

	// Merger模块处理队列, 反馈出消息合并的压力情况
	MergerPending int64 `json:"merger_pending"`

	// Merger模块合并发送的消息总数与失败总数
	MergerRoomTotal int64 `json:"merger_room_total"`
	MergerAllTotal  int64 `json:"merger_all_total"`
	MergerRoomFail  int64 `json:"merger_room_fail"`
	MergerAllFail   int64 `json:"merger_all_fail"`
}

func Stats() *Statistics {
	return &insStats
}

func (s *Statistics) GetStats() *Statistics {
	return s
}

func (s *Statistics) DispatchPendingIncr() {
	atomic.AddInt64(&s.DispatchPending, 1)
}

func (s *Statistics) DispatchPendingDesc() {
	atomic.AddInt64(&s.DispatchPending, -1)
}

func (s *Statistics) PushJobPendingIncr() {
	atomic.AddInt64(&s.PushJobPending, 1)
}

func (s *Statistics) PushJobPendingDesc() {
	atomic.AddInt64(&s.PushJobPending, -1)
}

func (s *Statistics) OnlineConnectionsIncr() {
	atomic.AddInt64(&s.OnlineConnections, 1)
}

func (s *Statistics) OnlineConnectionsDesc() {
	atomic.AddInt64(&s.OnlineConnections, -1)
}

func (s *Statistics) RoomCountIncr() {
	atomic.AddInt64(&s.RoomCount, 1)
}

func (s *Statistics) RoomCountDesc() {
	atomic.AddInt64(&s.RoomCount, -1)
}

func (s *Statistics) MergerPendingIncr() {
	atomic.AddInt64(&s.MergerPending, 1)
}

func (s *Statistics) MergerPendingDesc() {
	atomic.AddInt64(&s.MergerPending, -1)
}

func (s *Statistics) MergerAllTotalIncr(batchSize int64) {
	atomic.AddInt64(&s.MergerAllTotal, batchSize)
}

func (s *Statistics) MergerAllFailIncr(batchSize int64) {
	atomic.AddInt64(&s.MergerAllFail, batchSize)
}

func (s *Statistics) DispatchFailIncr() {
	atomic.AddInt64(&s.DispatchFail, 1)
}

func (s *Statistics) DispatchTotalIncr(batchSize int64) {
	atomic.AddInt64(&s.DispatchTotal, batchSize)
}

func (s *Statistics) DispatchFailNumIncr(batchSize int64) {
	atomic.AddInt64(&s.DispatchFail, batchSize)
}

func (s *Statistics) SendMessageFailIncr() {
	atomic.AddInt64(&s.SendMessageFail, 1)
}

func (s *Statistics) SendMessageTotalIncr() {
	atomic.AddInt64(&s.SendMessageTotal, 1)
}

func (s *Statistics) Dump() (data []byte, err error) {
	return json.Marshal(s)
}
