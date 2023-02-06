package bus_transport

import (
	"WSSFacade/protocol"
	"log"
)

type SampleTransport struct{}

func (s SampleTransport) AckMessage(message *protocol.Message) {
	//TODO implement me
	panic("implement me")
}

func (s SampleTransport) Init() error {
	//TODO implement me
	panic("implement me")
}

func (s SampleTransport) Destruct() error {
	//TODO implement me
	panic("implement me")
}

func (s SampleTransport) GetMessage() *protocol.Message {
	msg := protocol.NewMessage("*", "TestClient", []byte("Hello world !"))
	log.Printf("Found message in event bus for user `%s`", msg.Metadata.Recipient)
	return msg
}

func (s SampleTransport) SendMessage(message *protocol.Message) error {
	log.Printf("Sending message `%+v` to event bus", message)
	return nil
}

func NewSampleTransport() Transport {
	return &SampleTransport{}
}
