package queue

import "github.com/jxo-me/rabbitmq-go"

// GetDefaultClientOptions describes the options that will be used when a value isn't provided
func GetDefaultClientOptions(queueName string) ClientOptions {
	return ClientOptions{
		ConsumeOptions: ConsumeOptions{
			ConsumerName:    "",
			ConsumerAutoAck: false,
			Exclusive:       false,
		},
		QueueOptions: QueueOptions{
			Name:       queueName,
			Durable:    false,
			AutoDelete: false,
			Exclusive:  false,
			NoWait:     false,
			Passive:    false,
			Declare:    true,
			Args: map[string]interface{}{
				// Ensure the queue is deleted automatically when it's unused for
				// more than the set time. This is to ensure that messages that
				// are in flight during a reconnecting don't get lost (which might
				// happen when using `DeleteWhenUnused`).
				"x-expires": 1 * 60 * 1000, // 1 minute.
			},
		},
		ExchangeOptions: ExchangeOptions{
			Name:       "",
			Durable:    false,
			AutoDelete: false,
			Internal:   false,
			NoWait:     false,
			Passive:    false,
			Declare:    false,
			Args:       rabbitmq.Table{},
		},
	}
}

// WithClientPublishOptionsMandatory makes the publishing mandatory, which means when a queue is not
// bound to the routing key, a message will be sent back on the return channel for you to handle
func WithClientPublishOptionsMandatory(options *ClientOptions) {
	options.PublishOptions.Mandatory = true
}

// WithClientPublishOptionsImmediate makes the publishing immediate, which means when a consumer is not available
// to immediately handle the new message, a message will be sent back on the return channel for you to handle
func WithClientPublishOptionsImmediate(options *ClientOptions) {
	options.PublishOptions.Immediate = true
}

// WithClientOptionsConsumerName returns a function that sets the name on the server of this consumer
// if unset a random name will be given
func WithClientOptionsConsumerName(consumerName string) func(*ClientOptions) {
	return func(options *ClientOptions) {
		options.ConsumeOptions.ConsumerName = consumerName
	}
}

// WithClientOptionsConsumerAutoAck returns a function
// that sets the auto acknowledge property on the server of this consumer
// if unset, the default will be used (false)
func WithClientOptionsConsumerAutoAck(autoAck bool) func(*ClientOptions) {
	return func(options *ClientOptions) {
		options.ConsumeOptions.ConsumerAutoAck = autoAck
	}
}

// WithClientOptionsConsumerExclusive sets the consumer to exclusive, which means
// the server will ensure that this is the sole consumer
// from this queue. When exclusive is false, the server will fairly distribute
// deliveries across multiple consumers.
func WithClientOptionsConsumerExclusive(options *ClientOptions) {
	options.ConsumeOptions.Exclusive = true
}

// WithClientOptionsQueueDurable ensures the queue is a durable queue
func WithClientOptionsQueueDurable(options *ClientOptions) {
	options.QueueOptions.Durable = true
}

// WithClientOptionsQueueAutoDelete ensures the queue is an auto-delete queue
func WithClientOptionsQueueAutoDelete(options *ClientOptions) {
	options.QueueOptions.AutoDelete = true
}

// WithClientOptionsQueueExclusive ensures the queue is an exclusive queue
func WithClientOptionsQueueExclusive(options *ClientOptions) {
	options.QueueOptions.Exclusive = true
}

// WithClientOptionsQueueArgs adds optional args to the queue
func WithClientOptionsQueueArgs(args rabbitmq.Table) func(*ClientOptions) {
	return func(options *ClientOptions) {
		options.QueueOptions.Args = args
	}
}

// WithRabbitMqClientGroupName set group name address
func WithRabbitMqClientGroupName(name string) func(*ClientOptions) {
	return func(options *ClientOptions) {
		if name == "" {
			return
		}
		options.GroupName = name
	}
}
