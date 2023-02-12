package backplane

import (
	"MargayGateway/backplane/bus_transport"
	"MargayGateway/registry"
)

type EventBus interface {
	ConnectClient(connection *registry.Connection)
	proxyClientMessage(from string, payload []byte) error
	subscribe(transport bus_transport.Transport)
}
