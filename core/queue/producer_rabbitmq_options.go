package queue

// WithRabbitMqPublishOptionsExchange returns a function that sets the exchange to publish to
func WithRabbitMqPublishOptionsExchange(exchange string) func(*PublishOptions) {
	return func(options *PublishOptions) {
		options.Exchange = exchange
	}
}

// WithRabbitMqPublishOptionsContentType returns a function that sets the content type, i.e. "application/json"
func WithRabbitMqPublishOptionsContentType(contentType string) func(*PublishOptions) {
	return func(options *PublishOptions) {
		options.ContentType = contentType
	}
}

// WithRabbitMqPublishOptionsReplyTo returns a function that sets the reply to field
func WithRabbitMqPublishOptionsReplyTo(replyTo string) func(*PublishOptions) {
	return func(options *PublishOptions) {
		options.ReplyTo = replyTo
	}
}

// WithRabbitMqPublishOptionsMessageID returns a function that sets the message identifier
func WithRabbitMqPublishOptionsMessageID(messageID string) func(*PublishOptions) {
	return func(options *PublishOptions) {
		options.MessageID = messageID
	}
}

// WithRabbitMqPublishOptionsUserID returns a function that sets the user id i.e. "user"
func WithRabbitMqPublishOptionsUserID(userID string) func(*PublishOptions) {
	return func(options *PublishOptions) {
		options.UserID = userID
	}
}

// WithRabbitMqPublishOptionsAppID returns a function that sets the application id
func WithRabbitMqPublishOptionsAppID(appID string) func(*PublishOptions) {
	return func(options *PublishOptions) {
		options.AppID = appID
	}
}
