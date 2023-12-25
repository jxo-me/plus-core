package message

type IMessage interface {
	SetId(string)
	GetId() string
	SetExchange(string)
	GetExchange() string
	SetRoutingKey(string)
	GetRoutingKey() string
	SetValue(value any)
	GetValue() any
	SetErrorIncr()
	SetErrorCount(uint64)
	GetErrorCount() uint64
}
