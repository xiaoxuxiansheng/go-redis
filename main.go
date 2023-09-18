package main

import (
	"fmt"

	"github.com/xiaoxuxiansheng/go-redis/database"
	"github.com/xiaoxuxiansheng/go-redis/pkg/protocol"
	"github.com/xiaoxuxiansheng/go-redis/server"
	"github.com/xiaoxuxiansheng/go-redis/tcp"
)

func main() {
	parser := protocol.NewParser()
	db := database.NewDatabase()
	handler := server.NewHandler(parser, db)

	if err := tcp.ListAndServce("localhost:6399", handler); err != nil {
		fmt.Print(err)
	}
}
