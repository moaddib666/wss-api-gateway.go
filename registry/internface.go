package registry

type Registry interface {
	Add(connection *Connection) error
	Del(connection *Connection) error
	Get(connectionId string) (*Connection, error)
}
