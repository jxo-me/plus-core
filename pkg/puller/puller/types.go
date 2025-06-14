package puller

import "time"

// PullTaskConfig 单一厂商拉单配置
type PullTaskConfig struct {
	Vendor           string        // 厂商名称 (如: JILI, PG, RSG)
	StartTime        time.Time     // 拉单起始时间
	EndTime          time.Time     // 拉单截止时间
	WindowSize       time.Duration // 每个分片的时间长度
	Concurrency      int           // 并发度
	RetryCount       int           // 重试次数
	TimeoutPerWindow time.Duration // 每个分片的超时
}

// WindowPullResult 单个窗口拉取任务结果
type WindowPullResult struct {
	Vendor      string
	WindowStart time.Time
	WindowEnd   time.Time
	Records     []interface{}
}
