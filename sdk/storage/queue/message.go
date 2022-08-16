package queue

import (
	"github.com/jxo-me/plus-core/sdk/storage"
)

type Message struct {
	Id         string
	RoutingKey string
	Values     map[string]interface{}
	GroupId    string
	ErrorCount int
}

func (m *Message) GetId() string {
	return m.Id
}

func (m *Message) GetRoutingKey() string {
	return m.RoutingKey
}

func (m *Message) GetValues() map[string]interface{} {
	return m.Values
}

func (m *Message) SetId(id string) {
	m.Id = id
}

func (m *Message) SetRoutingKey(routingKey string) {
	m.RoutingKey = routingKey
}

func (m *Message) SetValues(values map[string]interface{}) {
	m.Values = values
}

func (m *Message) SetGroupId(group string) {
	m.GroupId = group
}

func (m *Message) GetGroupId() string {
	return m.GroupId
}

func (m *Message) GetPrefix() (prefix string) {
	if m.Values == nil {
		return
	}
	v, _ := m.Values[storage.PrefixKey]
	prefix, _ = v.(string)
	return
}

func (m *Message) SetPrefix(prefix string) {
	if m.Values == nil {
		m.Values = make(map[string]interface{})
	}
	m.Values[storage.PrefixKey] = prefix
}

func (m *Message) SetErrorCount() {
	m.ErrorCount = m.ErrorCount + 1
}

func (m *Message) GetErrorCount() int {
	return m.ErrorCount
}
