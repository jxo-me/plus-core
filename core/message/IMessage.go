package message

type IMessage interface {
	SetId(string)
	GetId() string
	SetRoutingKey(string)
	GetRoutingKey() string
	SetValue(value any)
	GetValue() any
	SetErrorIncr()
	SetErrorCount(uint64)
	GetErrorCount() uint64
}
