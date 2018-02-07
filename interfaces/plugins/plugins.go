package plugins

import "github.com/unchainio/pkg/interfaces/logger"

type Message struct {
	Tag        uint64
	Body       []byte
	Attributes []interface{}
}

func (m *Message) AddAttribute(attr interface{}) {
	if len(m.Attributes) > 0 {
		for _, v := range m.Attributes {
			if v == attr {
				return
			}
		}
		m.Attributes = append(m.Attributes, attr)
	} else {
		m.Attributes = append(m.Attributes, attr)
	}
}

func (m *Message) RemoveAttribute(attr interface{}) {
	for i, v := range m.Attributes {
		if v == attr {
			m.Attributes = append(m.Attributes[:i], m.Attributes[i+1:]...)
		}
	}
}

type EndpointPlugin interface {
	Init(config []byte, log logger.Logger) (err error)
	Send(message *Message) (response *Message, err error)
	Receive() (message *Message, err error)
	Ack(tag uint64, response *Message) error
	Nack(tag uint64) error
	Close() error
}

type ActionPlugin interface {
	Init(config []byte, log logger.Logger) (err error)
	Handle(message *Message) (result *Message, err error)
}
