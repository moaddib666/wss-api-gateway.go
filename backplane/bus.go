package backplane

import (
	"MargayGateway/backplane/bus_transport"
	"MargayGateway/constants"
	"MargayGateway/protocol"
	"MargayGateway/registry"
	"log"
	"time"
)

type DefaultBus struct {
	clients   registry.Registry
	transport bus_transport.Transport
}

func NewBus(clients registry.Registry) EventBus {
	bus := &DefaultBus{clients: clients}
	transport := bus_transport.NewRMQTransport()
	bus.subscribe(transport)
	return bus
}

func (s *DefaultBus) ConnectClient(connection *registry.Connection) {
	if s.clients.Add(connection) != nil {
		log.Printf("It's not allowed to have several connection from 1 client `%s`", connection.ConnectionId)
		connection.WebSocket.Close()
		return
	}
	go s.handleClientConnection(connection)
}

func (s *DefaultBus) proxyClientMessage(from string, payload []byte) error {
	msg := protocol.NewMessage(from, "*", payload)
	return s.transport.SendMessage(msg)
}

func (s *DefaultBus) handleClientConnection(connection *registry.Connection) {
	defer s.clients.Del(connection)
	err := s.proxyClientMessage(
		constants.AppInternalName,
		protocol.ClientConnectedEvent(connection.ConnectionId).Dump())
	if err != nil {
		log.Printf("Failed to send onConnect message %-v", err)
		return
	}
	defer s.proxyClientMessage(
		constants.AppInternalName,
		protocol.ClientDisconnectedEvent(connection.ConnectionId).Dump())
	for {
		msg, err := connection.GetMessage()
		if err != nil {
			// TODO handle error
			log.Printf("Failed to get messages %-v", err)
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

func (s *DefaultBus) subscribe(transport bus_transport.Transport) {
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
			if msg.Metadata.Recipient == constants.AppInternalName {
				log.Printf("Reject message due to invalid recipient `%s` is internal client", constants.AppInternalName)
				s.transport.AckMessage(msg)
				continue
			}
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
