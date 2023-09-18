package tcp

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type TCPHandler interface {
	Handle(ctx context.Context, conn net.Conn)
	Close() error
}

func ListAndServce(address string, handler TCPHandler) error {
	closeCh := make(chan struct{})
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		for {
			signal := <-sigCh
			switch signal {
			case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
				closeCh <- struct{}{}
				return
			default:
			}
		}
	}()
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	return listenerAndService(listener, handler, closeCh)
}

func listenerAndService(listener net.Listener, handler TCPHandler, closeCh chan struct{}) error {
	errCh := make(chan error, 1)
	defer close(errCh)
	go func() {
		select {
		case <-closeCh:
			fmt.Println("server closing...")
		case err := <-errCh:
			fmt.Printf("server err: %v", err)
		}
		_ = listener.Close()
		_ = handler.Close()
	}()

	ctx := context.Background()
	var wg sync.WaitGroup
	for {
		conn, err := listener.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				time.Sleep(5 * time.Millisecond)
				continue
			}

			errCh <- err
			break
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			handler.Handle(ctx, conn)
		}()
	}

	wg.Wait()
	return nil
}
