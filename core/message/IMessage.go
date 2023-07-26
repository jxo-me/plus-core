package message

const (
	PrefixKey = "__host"
)

type IMessage interface {
	SetId(string)
	GetId() string
	SetRoutingKey(string)
	GetRoutingKey() string
	SetValues(map[string]interface{})
	GetValues() map[string]interface{}
	GetPrefix() string
	SetPrefix(string)
	SetErrorIncr()
	SetErrorCount(uint64)
	GetErrorCount() uint64
}
