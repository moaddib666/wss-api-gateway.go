package registry

import (
	"fmt"
	"log"
	"sync"
)

type ConnectionRegistry struct {
	connections map[string]*Connection
	mu          sync.RWMutex
}

func (c *ConnectionRegistry) Add(connection *Connection) error {

	_, err := c.Get(connection.ConnectionId)
	if err == nil {
		return fmt.Errorf("connection for `%s` already exist", connection.ConnectionId)
	}
	c.mu.Lock()
	c.connections[connection.ConnectionId] = connection
	c.mu.Unlock()
	log.Printf("Connection `%s` added to cennection pool", connection.ConnectionId)
	return nil
}

func (c *ConnectionRegistry) Del(connection *Connection) error {
	_, err := c.Get(connection.ConnectionId)
	if err == nil {
		c.mu.Lock()
		delete(c.connections, connection.ConnectionId)
		c.mu.Unlock()
	}
	log.Printf("Connection `%s` removed from cennection pool", connection.ConnectionId)
	return connection.WebSocket.Close()
}

func (c *ConnectionRegistry) Get(connectionId string) (*Connection, error) {
	c.mu.RLock()
	conn, exist := c.connections[connectionId]
	c.mu.RUnlock()
	if exist {
		return conn, nil
	}
	return nil, fmt.Errorf("connection for `%s` does not exist", connectionId)
}

func GetConnectionRegistry() Registry {
	return &ConnectionRegistry{
		connections: map[string]*Connection{},
		mu:          sync.RWMutex{},
	}
}
