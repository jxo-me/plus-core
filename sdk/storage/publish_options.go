package storage

// PublishOptions are used to control how data is published
type PublishOptions struct {
	Exchange string
	// MIME content type
	ContentType string
	// address to to reply to (ex: RPC)
	ReplyTo string
	// message identifier
	MessageID string
	// creating user id - ex: "guest"
	UserID string
	// creating application id
	AppID string
}

// WithPublishOptionsExchange returns a function that sets the exchange to publish to
func WithPublishOptionsExchange(exchange string) func(*PublishOptions) {
	return func(options *PublishOptions) {
		options.Exchange = exchange
	}
}

// WithPublishOptionsContentType returns a function that sets the content type, i.e. "application/json"
func WithPublishOptionsContentType(contentType string) func(*PublishOptions) {
	return func(options *PublishOptions) {
		options.ContentType = contentType
	}
}

// WithPublishOptionsReplyTo returns a function that sets the reply to field
func WithPublishOptionsReplyTo(replyTo string) func(*PublishOptions) {
	return func(options *PublishOptions) {
		options.ReplyTo = replyTo
	}
}

// WithPublishOptionsMessageID returns a function that sets the message identifier
func WithPublishOptionsMessageID(messageID string) func(*PublishOptions) {
	return func(options *PublishOptions) {
		options.MessageID = messageID
	}
}

// WithPublishOptionsUserID returns a function that sets the user id i.e. "user"
func WithPublishOptionsUserID(userID string) func(*PublishOptions) {
	return func(options *PublishOptions) {
		options.UserID = userID
	}
}

// WithPublishOptionsAppID returns a function that sets the application id
func WithPublishOptionsAppID(appID string) func(*PublishOptions) {
	return func(options *PublishOptions) {
		options.AppID = appID
	}
}
