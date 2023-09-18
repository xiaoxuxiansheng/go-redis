package database

import (
	"fmt"

	"github.com/xiaoxuxiansheng/go-redis/connection"
	"github.com/xiaoxuxiansheng/go-redis/pkg/protocol"
)

type Database struct {
}

func NewDatabase() *Database {
	return &Database{}
}

func (d *Database) Exec(connection *connection.Connection, cmdLine [][]byte) protocol.Reply {
	fmt.Printf("exec cmd line: %v", cmdLine)
	return &protocol.NullReply{}
}

func (d *Database) Close() {
}
