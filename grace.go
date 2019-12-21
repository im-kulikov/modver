package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// newGracefulContext returns graceful context
func newGracefulContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		sig := <-ch
		fmt.Println("received signal", sig.String())
		cancel()
	}()
	return ctx
}
