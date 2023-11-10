package message

type Message struct {
	Id         string
	RoutingKey string
	Values     any
	GroupId    string
	ErrorCount uint64
}

func (m *Message) GetId() string {
	return m.Id
}

func (m *Message) GetRoutingKey() string {
	return m.RoutingKey
}

func (m *Message) GetValue() any {
	return m.Values
}

func (m *Message) SetId(id string) {
	m.Id = id
}

func (m *Message) SetRoutingKey(routingKey string) {
	m.RoutingKey = routingKey
}

func (m *Message) SetValue(values any) {
	m.Values = values
}

func (m *Message) SetErrorIncr() {
	m.ErrorCount = m.ErrorCount + 1
}

func (m *Message) SetErrorCount(count uint64) {
	m.ErrorCount = count
}

func (m *Message) GetErrorCount() uint64 {
	return m.ErrorCount
}
