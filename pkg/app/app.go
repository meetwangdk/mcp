package app

import (
	"context"
	"hcs-agent/pkg/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Client struct {
	name    string
	servers []server.Server
}

type Option func(a *Client)

func NewApp(opts ...Option) *Client {
	a := &Client{}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func WithServer(servers ...server.Server) Option {
	return func(a *Client) {
		a.servers = servers
	}
}

func WithName(name string) Option {
	return func(a *Client) {
		a.name = name
	}
}

func (a *Client) Run(ctx context.Context) error {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	for _, srv := range a.servers {
		go func(srv server.Server) {
			err := srv.Start(ctx)
			if err != nil {
				log.Printf("Server start err: %v", err)
			}
		}(srv)
	}

	select {
	case <-signals:
		// Received termination signal
		log.Println("Received termination signal")
	case <-ctx.Done():
		// Context canceled
		log.Println("Context canceled")
	}

	// Gracefully stop the servers
	for _, srv := range a.servers {
		err := srv.Stop(ctx)
		if err != nil {
			log.Printf("Server stop err: %v", err)
		}
	}

	return nil
}
