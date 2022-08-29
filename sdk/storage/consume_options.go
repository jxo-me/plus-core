package storage

// ConsumeOptions are used to describe how a new consumer will be created.
type ConsumeOptions struct {
	// rabbitmq
	BindingRoutingKeys []string
	BindingExchange    *BindingExchangeOptions
	Concurrency        int
	QOSPrefetch        int
	ConsumerName       string
	ConsumerAutoAck    bool
	// rocketmq
	GroupName         string
	MaxReconsumeTimes int32
	AutoCommit        bool
}
