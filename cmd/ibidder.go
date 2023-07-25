package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/nerock/invoicebidder/internal/config"
	"github.com/nerock/invoicebidder/internal/orchestrator/api"
	"github.com/nerock/invoicebidder/internal/orchestrator/broker"
)

type Server interface {
	Serve() error
	Shutdown(ctx context.Context) error
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	brk := broker.New(cfg.Broker.Handlers, cfg.Broker.Buffer, cfg.Broker.MaxRetries, nil, nil, nil)
	srv := api.New(cfg.Server.Port, nil, nil, nil, brk)

	run(srv, brk)
}

func run(servers ...Server) {
	for _, s := range servers {
		if err := s.Serve(); err != nil {
			log.Fatal(err)
		}
	}

	exit := make(chan os.Signal)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit

	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, s := range servers {
		go func(s Server) {
			wg.Add(1)
			defer wg.Done()

			if err := s.Shutdown(ctx); err != nil {
				log.Println(err)
			}
		}(s)
	}

	wg.Wait()
}
