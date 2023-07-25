package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/nerock/invoicebidder/internal/investor"
	investorStorage "github.com/nerock/invoicebidder/internal/investor/storage"

	"github.com/nerock/invoicebidder/internal/issuer"
	issuerStorage "github.com/nerock/invoicebidder/internal/issuer/storage"

	"github.com/nerock/invoicebidder/internal/invoice"
	invoiceStorage "github.com/nerock/invoicebidder/internal/invoice/storage"

	"github.com/jackc/pgx/v5/pgxpool"
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

	ctx := context.Background()
	invoiceDB, err := pgxpool.New(ctx, cfg.InvoiceDB)
	if err != nil {
		log.Fatal(err)
	}
	defer invoiceDB.Close()

	issuerDB, err := pgxpool.New(ctx, cfg.IssuerDB)
	if err != nil {
		log.Fatal(err)
	}
	defer issuerDB.Close()

	investorDB, err := pgxpool.New(ctx, cfg.InvestorDB)
	if err != nil {
		log.Fatal(err)
	}
	defer investorDB.Close()

	invoiceSvc := invoice.NewService(invoiceStorage.New(invoiceDB), invoiceStorage.NewFileStorage(cfg.BasePath))
	issuerSvc := issuer.NewService(issuerStorage.New(issuerDB))
	investorSvc := investor.NewService(investorStorage.New(investorDB))

	brk := broker.New(cfg.Broker.Handlers, cfg.Broker.Buffer, cfg.Broker.MaxRetries, invoiceSvc, investorSvc, issuerSvc)
	srv := api.New(cfg.Server.Port, invoiceSvc, investorSvc, issuerSvc, brk)

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
