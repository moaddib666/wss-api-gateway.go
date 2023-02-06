package backplane

import (
	"WSSFacade/backplane/bus_transport"
	"WSSFacade/constants"
	"WSSFacade/protocol"
	"WSSFacade/registry"
	"log"
	"time"
)

type SampleBus struct {
	clients   registry.Registry
	transport bus_transport.Transport
}

func NewSampleBus(clients registry.Registry) EventBus {
	bus := &SampleBus{clients: clients}
	transport := bus_transport.NewRMQTransport()
	bus.subscribe(transport)
	return bus
}

func (s *SampleBus) ConnectClient(connection *registry.Connection) {
	if s.clients.Add(connection) != nil {
		log.Printf("It's not allowed to have several connection from 1 client `%s`", connection.ConnectionId)
		connection.WebSocket.Close()
		return
	}
	go s.handleClientConnection(connection)
}

func (s *SampleBus) proxyClientMessage(from string, payload []byte) error {
	msg := protocol.NewMessage(from, "*", payload)
	return s.transport.SendMessage(msg)
}

func (s *SampleBus) handleClientConnection(connection *registry.Connection) {
	defer s.clients.Del(connection)
	err := s.proxyClientMessage(
		constants.AppInternalName,
		protocol.ClientConnectedEvent(connection.ConnectionId).Dump())
	if err != nil {
		return
	}
	defer s.proxyClientMessage(
		constants.AppInternalName,
		protocol.ClientDisconnectedEvent(connection.ConnectionId).Dump())
	for {
		msg, err := connection.GetMessage()
		if err != nil {
			// TODO handle error
			return
		}
		err = s.proxyClientMessage(connection.ConnectionId, msg)
		if err != nil {
			log.Printf("Failed to send message from %s - %v", connection.ConnectionId, err)
			return
		}
		time.Sleep(500)
	}
}

func (s *SampleBus) subscribe(transport bus_transport.Transport) {
	s.transport = transport
	err := s.transport.Init()
	if err != nil {
		log.Fatalf("can't setup transport %v", err)
	}

	go func() {
		defer s.transport.Destruct()
		for {
			time.Sleep(500)
			msg := s.transport.GetMessage()
			connection, err := s.clients.Get(msg.Metadata.Recipient)
			if err != nil {
				s.transport.AckMessage(nil)
				continue
			}
			log.Printf("Sent message to user `%s`", msg.Metadata.Recipient)
			err = connection.SendMessage(msg.Payload)
			if err != nil {
				log.Printf("Failed to send message to user `%s`", msg.Metadata.Recipient)
				continue
			}
			// TODO handle error
			s.transport.AckMessage(msg)
		}

	}()
}
