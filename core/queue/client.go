package queue

import "github.com/jxo-me/rabbitmq-go"

// ExchangeOptions are used to configure an exchange.
// If the Passive flag is set, the client will only check if the exchange exists on the server
// and that the settings match, no creation attempt will be made.
type ExchangeOptions struct {
	Name       string
	Kind       string // possible values: empty string for default exchange or direct, topic, fanout
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Passive    bool // if false, a missing exchange will be created on the server
	Args       rabbitmq.Table
	Declare    bool
}

// QueueOptions are used to configure a queue.
// A passive queue is assumed by RabbitMQ to already exist, and attempting to connect
// to a non-existent queue will cause RabbitMQ to throw an exception.
type QueueOptions struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Passive    bool // if false, a missing queue will be created on the server
	Args       rabbitmq.Table
	Declare    bool
}

type ClientOptions struct {
	ConsumeOptions  ConsumeOptions
	QueueOptions    QueueOptions
	ExchangeOptions ExchangeOptions
	PublishOptions  PublishOptions
	// ConfirmMode puts the channel that messages are published over in
	// confirmation mode.
	// This makes sending requests more reliable at the cost
	// of some performance.
	// The server must confirm each publishing.
	// See https://www.rabbitmq.com/confirms.html#publisher-confirms
	ConfirmMode bool
	// rocketmq
	GroupName string
}
