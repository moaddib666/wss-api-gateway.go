package registry

import (
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
	return wsutil.WriteServerMessage(c.WebSocket, code, msg)
}

func (c *Connection) GetMessage() ([]byte, error) {
	// OP Code omitted
	msg, _, err := wsutil.ReadClientData(c.WebSocket)
	if err != nil {
		return []byte{}, err
	}
	return msg, nil
}
