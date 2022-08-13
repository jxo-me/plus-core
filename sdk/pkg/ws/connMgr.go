package ws

import (
	"context"
)

var (
	insConnManager = ConnManager{}
)

// PushJob 推送任务
type PushJob struct {
	pushType int    // 推送类型
	roomId   string // 房间ID
	// union {
	Msg   *MessageReq // 未序列化的业务消息
	wsMsg *Message    //  已序列化的业务消息
	// }
}

// ConnManager 连接管理器
type ConnManager struct {
	buckets []*Bucket
	jobChan []chan *PushJob // 每个Bucket对应一个Job Queue

	dispatchChan chan *PushJob // 待分发消息队列
}

func ConnMgr() *ConnManager {
	return &insConnManager
}

func InitConnMgr(config *Config) *ConnManager {
	var (
		bucketIdx         int
		jobWorkerIdx      int
		dispatchWorkerIdx int
	)

	insConnManager = ConnManager{
		buckets:      make([]*Bucket, config.BucketCount),
		jobChan:      make([]chan *PushJob, config.BucketCount),
		dispatchChan: make(chan *PushJob, config.DispatchChannelSize),
	}
	for bucketIdx = range insConnManager.buckets {
		insConnManager.buckets[bucketIdx] = InitBucket(bucketIdx)                            // 初始化Bucket
		insConnManager.jobChan[bucketIdx] = make(chan *PushJob, config.BucketJobChannelSize) // Bucket的Job队列
		// Bucket的Job worker
		for jobWorkerIdx = 0; jobWorkerIdx < config.BucketJobWorkerCount; jobWorkerIdx++ {
			go insConnManager.jobWorkerMain(bucketIdx)
		}
	}
	// 初始化分发协程, 用于将消息扇出给各个Bucket
	for dispatchWorkerIdx = 0; dispatchWorkerIdx < config.DispatchWorkerCount; dispatchWorkerIdx++ {
		go insConnManager.dispatchWorkerMain()
	}

	return &insConnManager
}

func (connMgr *ConnManager) JoinRoom(roomId string, wsConn *Connection) (err error) {
	var (
		bucket *Bucket
	)

	bucket = connMgr.GetBucket(wsConn)
	err = bucket.JoinRoom(roomId, wsConn)
	return
}

func (connMgr *ConnManager) LeaveRoom(roomId string, wsConn *Connection) (err error) {
	var (
		bucket *Bucket
	)

	bucket = connMgr.GetBucket(wsConn)
	err = bucket.LeaveRoom(roomId, wsConn)
	return
}

// 消息分发到Bucket
func (connMgr *ConnManager) dispatchWorkerMain() {
	var (
		bucketIdx int
		pushJob   *PushJob
		err       error
	)
	for {
		select {
		case pushJob = <-connMgr.dispatchChan:
			Stats().DispatchPendingDesc()

			// 序列化
			if pushJob.wsMsg, err = EncodeWsMessage(pushJob.Msg); err != nil {
				continue
			}
			// 分发给所有Bucket, 若Bucket拥塞则等待
			for bucketIdx = range connMgr.buckets {
				Stats().PushJobPendingIncr()
				connMgr.jobChan[bucketIdx] <- pushJob
			}
		}
	}
}

// Job负责消息广播给客户端
func (connMgr *ConnManager) jobWorkerMain(bucketIdx int) {
	var (
		bucket  = connMgr.buckets[bucketIdx]
		pushJob *PushJob
	)

	for {
		select {
		case pushJob = <-connMgr.jobChan[bucketIdx]: // 从Bucket的job queue取出一个任务
			Stats().PushJobPendingDesc()
			if pushJob.pushType == PushTypeAll {
				bucket.PushAll(pushJob.wsMsg)
			}
		}
	}
}

func (connMgr *ConnManager) GetBucket(wsConnection *Connection) (bucket *Bucket) {
	bucket = connMgr.buckets[wsConnection.connId%uint64(len(connMgr.buckets))]
	return
}

func (connMgr *ConnManager) GetBucketByCID(connId uint64) (bucket *Bucket) {
	bucket = connMgr.buckets[connId%uint64(len(connMgr.buckets))]
	return
}

func (connMgr *ConnManager) AddConn(wsConnection *Connection) {
	var (
		bucket *Bucket
	)

	bucket = connMgr.GetBucket(wsConnection)
	bucket.AddConn(wsConnection)

	Stats().OnlineConnectionsIncr()
}

func (connMgr *ConnManager) DelConn(wsConnection *Connection) {
	var (
		bucket *Bucket
	)

	bucket = connMgr.GetBucket(wsConnection)
	bucket.DelConn(wsConnection)

	Stats().OnlineConnectionsDesc()
}

func (connMgr *ConnManager) DelConnByCID(ctx context.Context, connId uint64) {
	var (
		bucket *Bucket
	)
	bucket = connMgr.GetBucketByCID(connId)
	bucket.DelConnByCID(ctx, connId)
}

// PushAll 向所有在线用户发送消息
func (connMgr *ConnManager) PushAll(msg *MessageReq) (err error) {
	var (
		pushJob *PushJob
	)

	pushJob = &PushJob{
		pushType: PushTypeAll,
		Msg:      msg,
	}

	select {
	case connMgr.dispatchChan <- pushJob:
		Stats().DispatchPendingIncr()
	default:
		err = ErrDispatchChannelFull
		Stats().DispatchFailIncr()
	}
	return
}

func (connMgr *ConnManager) PushRoom(roomId string, items *MessageReq) (err error) {
	var (
		pushJob *PushJob
	)
	pushJob = &PushJob{
		pushType: PushTypeRoom,
		roomId:   roomId,
		Msg:      items,
	}

	select {
	case connMgr.dispatchChan <- pushJob:
		Stats().DispatchTotalIncr(1)
	default:
		Stats().DispatchFailNumIncr(1)
		err = ErrDispatchChannelFull
	}
	return
}

// PushByCID 向指定用户发送消息
func (connMgr *ConnManager) PushByCID(connId uint64, message *MessageReq) (err error) {
	var (
		wsMsg  *Message
		bucket *Bucket
	)

	bucket = connMgr.GetBucketByCID(connId)
	// 序列化
	if wsMsg, err = EncodeWsMessage(message); err != nil {
		return err
	}
	err = bucket.PushByCid(connId, wsMsg)
	return
}
