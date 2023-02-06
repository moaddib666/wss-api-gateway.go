package backplane

import (
	"WSSFacade/backplane/bus_transport"
	"WSSFacade/registry"
)

type EventBus interface {
	ConnectClient(connection *registry.Connection)
	proxyClientMessage(from string, payload []byte) error
	subscribe(transport bus_transport.Transport)
}
