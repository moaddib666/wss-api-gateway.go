package bus_transport

import "MargayGateway/protocol"

type Transport interface {
	Init() error
	GetMessage() *protocol.Message
	SendMessage(message *protocol.Message) error
	Destruct() error
	AckMessage(message *protocol.Message)
	String() string
}
