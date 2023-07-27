package queue

// WithRocketMqPublishGroupName set group name address
func WithRocketMqPublishGroupName(name string) func(*PublishOptions) {
	return func(options *PublishOptions) {
		if name == "" {
			return
		}
		options.GroupName = name
	}
}

// WithRocketMqPublishRetry return an Option that specifies the retry times when send failed.
func WithRocketMqPublishRetry(times int) func(*PublishOptions) {
	return func(options *PublishOptions) {
		options.RetryTimes = times
	}
}
