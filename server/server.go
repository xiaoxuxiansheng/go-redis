package server

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/xiaoxuxiansheng/go-redis/connection"
	"github.com/xiaoxuxiansheng/go-redis/pkg/protocol"
)

var (
	unknownErrReplyBytes = []byte("-ERR unknown\r\n")
)

type DB interface {
	Exec(connection *connection.Connection, cmdLine [][]byte) protocol.Reply
	Close()
}

type Handler struct {
	closing     int32
	parser      *protocol.Parser
	db          DB
	activeConns sync.Map
}

func NewHandler(parser *protocol.Parser, db DB) *Handler {
	return &Handler{
		parser: parser,
		db:     db,
	}
}

func (h *Handler) Handle(ctx context.Context, conn net.Conn) {
	if atomic.LoadInt32(&h.closing) == 1 {
		_ = conn.Close()
		return
	}

	connection := connection.NewConnection(conn)
	h.activeConns.Store(connection, struct{}{})

	ch := h.parser.ParseStream(conn)
	for payload := range ch {
		if isClosedErr(payload.Err) {
			return
		}

		if payload.Err != nil {
			errReply := protocol.NewErrReply(payload.Err.Error())
			_, _ = connection.Write(errReply.ToBytes())
			continue
		}

		if payload.Reply == nil {
			continue
		}

		multiBuildReply, ok := payload.Reply.(*protocol.MultiBulkReply)
		if !ok {
			fmt.Println("not multi bulk reply..")
			continue
		}

		if reply := h.db.Exec(connection, multiBuildReply.Args); reply != nil {
			_, _ = connection.Write(reply.ToBytes())
			continue
		}

		_, _ = connection.Write(unknownErrReplyBytes)
	}
}

func (h *Handler) Close() error {
	atomic.StoreInt32(&h.closing, 1)
	h.activeConns.Range(func(rawConnection, _ any) bool {
		connection, _ := rawConnection.(*connection.Connection)
		_ = connection.Close()
		return true
	})
	h.db.Close()
	return nil
}

func isClosedErr(err error) bool {
	if err == nil {
		return false
	}

	if err == io.EOF || err == io.ErrUnexpectedEOF {
		return true
	}

	return strings.Contains(err.Error(), "use of closed network connection")
}
