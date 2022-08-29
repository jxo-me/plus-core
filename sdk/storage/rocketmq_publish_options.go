package storage

// WithRocketMqPublishGroupName set group name address
func WithRocketMqPublishGroupName(name string) func(*PublishOptions) {
	return func(options *PublishOptions) {
		if name == "" {
			return
		}
		options.GroupName = name
	}
}

// WithRocketMqPublishRetry return a Option that specifies the retry times when send failed.
func WithRocketMqPublishRetry(times int) func(*PublishOptions) {
	return func(options *PublishOptions) {
		options.RetryTimes = times
	}
}
