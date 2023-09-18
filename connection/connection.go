package connection

import "net"

type Connection struct {
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{}
}

func (c *Connection) Write(b []byte) (int, error) {
	return 0, nil
}

func (c *Connection) Read() ([]byte, error) {
	return nil, nil
}

func (c *Connection) Close() error {
	return nil
}
