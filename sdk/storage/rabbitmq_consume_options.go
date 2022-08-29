package storage

// GetDefaultConsumeOptions descibes the options that will be used when a value isn't provided
func GetDefaultConsumeOptions() ConsumeOptions {
	return ConsumeOptions{
		BindingExchange: &BindingExchangeOptions{
			Name: "",
			Kind: "direct",
		},
		Concurrency:     1,
		QOSPrefetch:     0,
		ConsumerName:    "",
		ConsumerAutoAck: false,
	}
}

// WithRabbitMqConsumeOptionsBindingRoutingKeys returns a function that sets the exchange name the RoutingKeys will be bound to
func WithRabbitMqConsumeOptionsBindingRoutingKeys(keys []string) func(*ConsumeOptions) {
	return func(options *ConsumeOptions) {
		options.BindingRoutingKeys = keys
	}
}

// getBindingExchangeOptionsOrSetDefault returns pointer to current BindingExchange options. if no BindingExchange options are set yet, it will set it with default values.
func getBindingExchangeOptionsOrSetDefault(options *ConsumeOptions) *BindingExchangeOptions {
	if options.BindingExchange == nil {
		options.BindingExchange = &BindingExchangeOptions{
			Name: "",
			Kind: "direct",
		}
	}
	return options.BindingExchange
}

// BindingExchangeOptions are used when binding to an exchange.
// it will verify the exchange is created before binding to it.
type BindingExchangeOptions struct {
	Name string
	Kind string
}

// WithRabbitMqConsumeOptionsBindingExchangeName returns a function that sets the exchange name the queue will be bound to
func WithRabbitMqConsumeOptionsBindingExchangeName(name string) func(*ConsumeOptions) {
	return func(options *ConsumeOptions) {
		getBindingExchangeOptionsOrSetDefault(options).Name = name
	}
}

// WithRabbitMqConsumeOptionsBindingExchangeType returns a function that sets the binding exchange kind/type
func WithRabbitMqConsumeOptionsBindingExchangeType(kind string) func(*ConsumeOptions) {
	return func(options *ConsumeOptions) {
		getBindingExchangeOptionsOrSetDefault(options).Kind = kind
	}
}

// WithRabbitMqConsumeOptionsConcurrency returns a function that sets the concurrency, which means that
// many goroutines will be spawned to run the provided handler on messages
func WithRabbitMqConsumeOptionsConcurrency(concurrency int) func(*ConsumeOptions) {
	return func(options *ConsumeOptions) {
		options.Concurrency = concurrency
	}
}

// WithRabbitMqConsumeOptionsQOSPrefetch returns a function that sets the prefetch count, which means that
// many messages will be fetched from the server in advance to help with throughput.
// This doesn't affect the handler, messages are still processed one at a time.
func WithRabbitMqConsumeOptionsQOSPrefetch(prefetchCount int) func(*ConsumeOptions) {
	return func(options *ConsumeOptions) {
		options.QOSPrefetch = prefetchCount
	}
}

// WithRabbitMqConsumeOptionsConsumerName returns a function that sets the name on the server of this consumer
// if unset a random name will be given
func WithRabbitMqConsumeOptionsConsumerName(consumerName string) func(*ConsumeOptions) {
	return func(options *ConsumeOptions) {
		options.ConsumerName = consumerName
	}
}

// WithRabbitMqConsumeOptionsConsumerAutoAck returns a function that sets the auto acknowledge property on the server of this consumer
// if unset the default will be used (false)
func WithRabbitMqConsumeOptionsConsumerAutoAck(autoAck bool) func(*ConsumeOptions) {
	return func(options *ConsumeOptions) {
		options.ConsumerAutoAck = autoAck
	}
}
