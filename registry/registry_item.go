package registry

import (
	"MargayGateway/constants"
	"MargayGateway/monitoring"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"net"
)

type Connection struct {
	ConnectionId string
	WebSocket    net.Conn
}

func CreateRegistryItem(id string, con net.Conn) *Connection {
	return &Connection{ConnectionId: id, WebSocket: con}
}

func (c *Connection) SendMessage(msg []byte) error {
	code := ws.OpBinary
	err := wsutil.WriteServerMessage(c.WebSocket, code, msg)
	if err != nil {
		monitoring.IncrementMessageCount(constants.DefaultRoute, constants.OutboundMessage)
	}
	return err
}

func (c *Connection) GetMessage() ([]byte, error) {
	// Operation code is ignored
	msg, _, err := wsutil.ReadClientData(c.WebSocket)
	if err != nil {
		return []byte{}, err
	}
	monitoring.IncrementMessageCount(constants.DefaultRoute, constants.InboundMessage)
	return msg, nil
}
