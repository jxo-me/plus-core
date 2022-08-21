package storage

// GetDefaultConsumeOptions descibes the options that will be used when a value isn't provided
func GetDefaultConsumeOptions() ConsumeOptions {
	return ConsumeOptions{
		BindingExchange: &BindingExchangeOptions{
			Name:       "",
			Kind:       "direct",
			Durable:    false,
			AutoDelete: false,
			Internal:   false,
			NoWait:     false,
			Declare:    true,
		},
		Concurrency:     1,
		QOSPrefetch:     0,
		ConsumerName:    "",
		ConsumerAutoAck: false,
		ConsumerNoLocal: false,
	}
}

// ConsumeOptions are used to describe how a new consumer will be created.
type ConsumeOptions struct {
	BindingRoutingKeys []string
	BindingExchange    *BindingExchangeOptions
	Concurrency        int
	QOSPrefetch        int
	ConsumerName       string
	ConsumerAutoAck    bool
	ConsumerNoLocal    bool
}

// WithConsumeOptionsBindingRoutingKeys returns a function that sets the exchange name the RoutingKeys will be bound to
func WithConsumeOptionsBindingRoutingKeys(keys []string) func(*ConsumeOptions) {
	return func(options *ConsumeOptions) {
		options.BindingRoutingKeys = keys
	}
}

// getBindingExchangeOptionsOrSetDefault returns pointer to current BindingExchange options. if no BindingExchange options are set yet, it will set it with default values.
func getBindingExchangeOptionsOrSetDefault(options *ConsumeOptions) *BindingExchangeOptions {
	if options.BindingExchange == nil {
		options.BindingExchange = &BindingExchangeOptions{
			Name:       "",
			Kind:       "direct",
			Durable:    false,
			AutoDelete: false,
			Internal:   false,
			NoWait:     false,
			Declare:    true,
		}
	}
	return options.BindingExchange
}

// BindingExchangeOptions are used when binding to an exchange.
// it will verify the exchange is created before binding to it.
type BindingExchangeOptions struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Declare    bool
}

// WithConsumeOptionsBindingExchangeName returns a function that sets the exchange name the queue will be bound to
func WithConsumeOptionsBindingExchangeName(name string) func(*ConsumeOptions) {
	return func(options *ConsumeOptions) {
		getBindingExchangeOptionsOrSetDefault(options).Name = name
	}
}

// WithConsumeOptionsBindingExchangeKind returns a function that sets the binding exchange kind/type
func WithConsumeOptionsBindingExchangeKind(kind string) func(*ConsumeOptions) {
	return func(options *ConsumeOptions) {
		getBindingExchangeOptionsOrSetDefault(options).Kind = kind
	}
}

// WithConsumeOptionsBindingExchangeDurable returns a function that sets the binding exchange durable flag
func WithConsumeOptionsBindingExchangeDurable(options *ConsumeOptions) {
	getBindingExchangeOptionsOrSetDefault(options).Durable = true
}

// WithConsumeOptionsBindingExchangeAutoDelete returns a function that sets the binding exchange autoDelete flag
func WithConsumeOptionsBindingExchangeAutoDelete(options *ConsumeOptions) {
	getBindingExchangeOptionsOrSetDefault(options).AutoDelete = true
}

// WithConsumeOptionsBindingExchangeInternal returns a function that sets the binding exchange internal flag
func WithConsumeOptionsBindingExchangeInternal(options *ConsumeOptions) {
	getBindingExchangeOptionsOrSetDefault(options).Internal = true
}

// WithConsumeOptionsBindingExchangeNoWait returns a function that sets the binding exchange noWait flag
func WithConsumeOptionsBindingExchangeNoWait(options *ConsumeOptions) {
	getBindingExchangeOptionsOrSetDefault(options).NoWait = true
}

// WithConsumeOptionsBindingExchangeSkipDeclare returns a function that skips the declaration of the
// binding exchange. Use this setting if the exchange already exists and you don't need to declare
// it on consumer start.
func WithConsumeOptionsBindingExchangeSkipDeclare(options *ConsumeOptions) {
	getBindingExchangeOptionsOrSetDefault(options).Declare = false
}

// WithConsumeOptionsConcurrency returns a function that sets the concurrency, which means that
// many goroutines will be spawned to run the provided handler on messages
func WithConsumeOptionsConcurrency(concurrency int) func(*ConsumeOptions) {
	return func(options *ConsumeOptions) {
		options.Concurrency = concurrency
	}
}

// WithConsumeOptionsQOSPrefetch returns a function that sets the prefetch count, which means that
// many messages will be fetched from the server in advance to help with throughput.
// This doesn't affect the handler, messages are still processed one at a time.
func WithConsumeOptionsQOSPrefetch(prefetchCount int) func(*ConsumeOptions) {
	return func(options *ConsumeOptions) {
		options.QOSPrefetch = prefetchCount
	}
}

// WithConsumeOptionsConsumerName returns a function that sets the name on the server of this consumer
// if unset a random name will be given
func WithConsumeOptionsConsumerName(consumerName string) func(*ConsumeOptions) {
	return func(options *ConsumeOptions) {
		options.ConsumerName = consumerName
	}
}

// WithConsumeOptionsConsumerAutoAck returns a function that sets the auto acknowledge property on the server of this consumer
// if unset the default will be used (false)
func WithConsumeOptionsConsumerAutoAck(autoAck bool) func(*ConsumeOptions) {
	return func(options *ConsumeOptions) {
		options.ConsumerAutoAck = autoAck
	}
}
