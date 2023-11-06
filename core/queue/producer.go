package queue

// PublishOptions are used to control how data is published
type PublishOptions struct {
	Exchange string
	// MIME content type
	ContentType string
	// address to reply to (ex: RPC)
	ReplyTo string
	// message identifier
	MessageID string
	// creating user id - ex: "guest"
	UserID string
	// creating application id
	AppID string

	Mandatory bool
	Immediate bool

	// rocketmq
	GroupName  string
	RetryTimes int
}
