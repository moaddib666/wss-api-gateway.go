package registry

import (
	"bytes"
	"net"
	"testing"
	"time"
)

func TestCreateRegistryItem(t *testing.T) {
	connection := CreateRegistryItem("test", newMockConn([]byte("test")))
	if connection.ConnectionId != "test" {
		t.Errorf("Error creating registry item: %s", connection.ConnectionId)
	}
}

type MockConnection struct {
	buf *bytes.Buffer
}

func newMockConn(data []byte) *MockConnection {
	return &MockConnection{buf: bytes.NewBuffer(data)}
}

func (c *MockConnection) Read(b []byte) (n int, err error) {
	return c.buf.Read(b)
}

func (c *MockConnection) Write(b []byte) (n int, err error) {
	return len(b), nil
}

func (c *MockConnection) Close() error {
	return nil
}

func (c *MockConnection) LocalAddr() net.Addr {
	return nil
}

func (c *MockConnection) RemoteAddr() net.Addr {
	return nil
}

func (c *MockConnection) SetDeadline(t time.Time) error {
	return nil
}

func (c *MockConnection) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *MockConnection) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestConnection_SendMessage(t *testing.T) {
	mockCon := newMockConn([]byte("test"))
	connection := CreateRegistryItem("test", mockCon)
	err := connection.SendMessage([]byte("test"))
	if err != nil {
		t.Errorf("Error sending message: %s", err)
	}
}

func TestConnection_GetMessage(t *testing.T) {
	connection := CreateRegistryItem("test", newMockConn([]byte("test")))
	_, err := connection.GetMessage()
	if err == nil {
		t.Errorf("Error getting message: %s", err)
	}
}
