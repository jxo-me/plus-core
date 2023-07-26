package queue

// WithRocketMqGroupName set group name address
func WithRocketMqGroupName(name string) func(*ConsumeOptions) {
	return func(options *ConsumeOptions) {
		if name == "" {
			return
		}
		options.GroupName = name
	}
}

// WithRocketMqMaxReconsumeTimes set MaxReconsumeTimes of options, if message reconsume greater than MaxReconsumeTimes, it will
// be sent to retry or dlq topic. more info reference by examples/consumer/retry.
func WithRocketMqMaxReconsumeTimes(times int32) func(*ConsumeOptions) {
	return func(options *ConsumeOptions) {
		options.MaxReconsumeTimes = times
	}
}

func WithRocketMqAutoCommit(auto bool) func(*ConsumeOptions) {
	return func(options *ConsumeOptions) {
		options.AutoCommit = auto
	}
}
