package ws

type Config struct {
	// websocket HTTP握手读超时 单位毫秒
	ReadTimeout int `json:"readTimeout" yaml:"readTimeout"`
	// websocket HTTP握手写超时 单位毫秒
	WriteTimeout int `json:"writeTimeout" yaml:"writeTimeout"`
	// websocket读队列长度 一般不需要修改
	InChannelSize int `json:"inChannelSize" yaml:"inChannelSize"`
	// WebSocket写队列长度 一般不需要修改
	OutChannelSize int `json:"outChannelSize" yaml:"outChannelSize"`
	// WebSocket心跳检查间隔 单位秒, 超过时间没有收到心跳, 服务端将主动断开链接
	HeartbeatInterval int `json:"heartbeatInterval" yaml:"heartbeatInterval"`
	// 连接分桶的数量 桶越多, 推送的锁粒度越小, 推送并发度越高
	BucketCount int `json:"bucketCount" yaml:"bucketCount"`
	// 每个桶的处理协程数量 影响同一时刻可以有多少个不同消息被分发出去
	BucketWorkerCount int `json:"bucketWorkerCount" yaml:"bucketWorkerCount"`
	// bucket工作队列长度 每个bucket的分发任务放在一个独立队列中
	BucketJobChannelSize int `json:"bucketJobChannelSize" yaml:"bucketJobChannelSize"`
	// bucket发送协程的数量 每个bucket有多个协程并发的推送消息
	BucketJobWorkerCount int `json:"bucketJobWorkerCount" yaml:"bucketJobWorkerCount"`
	// 待分发队列的长度 分发队列缓冲所有待推送的消息, 等待被分发到bucket
	DispatchChannelSize int `json:"dispatchChannelSize" yaml:"dispatchChannelSize"`
	// 分发协程的数量 分发协程用于将待推送消息扇出给各个bucket
	DispatchWorkerCount int `json:"dispatchWorkerCount" yaml:"dispatchWorkerCount"`
	// 合并推送的最大延迟时间 单位毫秒, 在抵达maxPushBatchSize之前超时则发送
	MaxMergerDelay int `json:"maxMergerDelay" yaml:"maxMergerDelay"`
	// 合并最多消息条数 消息推送频次越高, 应该使用更大的合并批次, 得到更高的吞吐收益
	MaxMergerBatchSize int `json:"maxMergerBatchSize" yaml:"maxMergerBatchSize"`
	// 消息合并协程的数量 消息合并与json编码耗费CPU, 注意一个房间的消息只会由同一个协程处理.
	MergerWorkerCount int `json:"mergerWorkerCount" yaml:"mergerWorkerCount"`
	// 消息合并队列的容量 每个房间消息合并线程有一个队列, 推送量超过队列将被丢弃
	MergerChannelSize int `json:"mergerChannelSize" yaml:"mergerChannelSize"`
	// 每个房间连接最多加入数量
	MaxJoinRoom int `json:"maxJoinRoom" yaml:"maxJoinRoom"`
}
