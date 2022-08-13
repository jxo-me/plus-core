package ws

import (
	"encoding/json"
	"fmt"
	"time"
)

type PushBatch struct {
	items       []*json.RawMessage
	commitTimer *time.Timer

	// union {
	room string // 按room合并
	// }
}

type PushContext struct {
	msg *json.RawMessage

	// union {
	room string // 按room合并
	// }
}

type MergeWorker struct {
	mergeType int // 合并类型: 广播, room, uid...

	contextChan chan *PushContext
	timeoutChan chan *PushBatch

	// union {
	room2Batch map[string]*PushBatch // room合并
	allBatch   *PushBatch            // 广播合并
	// }
}

// Merger 广播消息、消息的合并
type Merger struct {
	roomWorkers     []*MergeWorker // 房间合并
	broadcastWorker *MergeWorker   // 广播合并
}

func InitMerger(config *Config) *Merger {
	insMerger := Merger{
		roomWorkers: make([]*MergeWorker, config.MergerWorkerCount),
	}
	insMerger.broadcastWorker = initMergeWorker(PushTypeAll, config)

	return &insMerger
}

func (worker *MergeWorker) autoCommit(batch *PushBatch) func() {
	return func() {
		worker.timeoutChan <- batch
	}
}

func (worker *MergeWorker) commitBatch(batch *PushBatch) (err error) {
	var (
		bizPushData *BizPushData
		message     *MessageReq
		buf         []byte
	)

	bizPushData = &BizPushData{
		Items: batch.items,
	}
	if buf, err = json.Marshal(*bizPushData); err != nil {
		return
	}

	message = &MessageReq{
		Data: buf,
	}

	// 打包发送
	if worker.mergeType == PushTypeAll {
		worker.allBatch = nil
		err = ConnMgr().PushAll(message)
	}
	return
}

func (worker *MergeWorker) mergeWorkerMain(config *Config) {
	var (
		context      *PushContext
		batch        *PushBatch
		timeoutBatch *PushBatch
		isCreated    bool
		err          error
	)
	for {
		select {
		case context = <-worker.contextChan:
			Stats().MergerPendingDesc()

			isCreated = false
			// 按房间合并
			if worker.mergeType == PushTypeAll { // 广播合并
				batch = worker.allBatch
				if batch == nil {
					batch = &PushBatch{}
					worker.allBatch = batch
					isCreated = true
				}
			}

			// 合并消息
			batch.items = append(batch.items, context.msg)

			// 新建批次, 启动超时自动提交
			if isCreated {
				batch.commitTimer = time.AfterFunc(time.Duration(config.MaxMergerDelay)*time.Millisecond, worker.autoCommit(batch))
			}

			// 批次未满, 继续等待下次提交
			if len(batch.items) < config.MaxMergerBatchSize {
				continue
			}

			// 批次已满, 取消超时自动提交
			batch.commitTimer.Stop()
		case timeoutBatch = <-worker.timeoutChan:
			if worker.mergeType == PushTypeAll {
				batch = worker.allBatch
				// 定时器触发时, 批次已被提交
				if timeoutBatch != batch {
					continue
				}
			}
		}
		// 提交批次
		err = worker.commitBatch(batch)

		// 打点统计
		if worker.mergeType == PushTypeAll {
			Stats().MergerAllTotalIncr(int64(len(batch.items)))
			if err != nil {
				Stats().MergerAllFailIncr(int64(len(batch.items)))
			}
		}
	}
}

func initMergeWorker(mergeType int, config *Config) (worker *MergeWorker) {
	worker = &MergeWorker{
		mergeType:   mergeType,
		room2Batch:  make(map[string]*PushBatch),
		contextChan: make(chan *PushContext, config.MergerChannelSize),
		timeoutChan: make(chan *PushBatch, config.MergerChannelSize),
	}
	go worker.mergeWorkerMain(config)
	return
}

func (worker *MergeWorker) pushAll(msg *json.RawMessage) (err error) {
	var (
		context *PushContext
	)
	context = &PushContext{
		msg: msg,
	}
	select {
	case worker.contextChan <- context:
		Stats().MergerPendingIncr()
	default:
		err = ErrMergeChannelFull
	}
	return
}

// PushAll 广播合并推送
func (merger *Merger) PushAll(msg *json.RawMessage) (err error) {
	return merger.broadcastWorker.pushAll(msg)
}

func (merger *Merger) PushByCid(connId uint64, msg *json.RawMessage) (err error) {
	fmt.Println(connId, msg)
	return err
}
