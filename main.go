package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"

	"github.com/xiaoxuxiansheng/go-redis/tcp"
)

type myHandler struct {
}

func (m *myHandler) Handle(ctx context.Context, conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("eof")
			} else {
				fmt.Print(err)
			}
			return
		}

		conn.Write([]byte(msg))
	}
}

func (m *myHandler) Close() error {
	return nil
}

func main() {
	if err := tcp.ListAndServce("localhost:6399", &myHandler{}); err != nil {
		fmt.Print(err)
	}
}
