package api

import (
	"context"
	"fmt"

	"github.com/bojanz/currency"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

var currFmt = currency.NewFormatter(currency.NewLocale("fr"))

type Broker interface {
	SendTradeEvent(string, bool)
	SendFailedBidEvent(string, currency.Amount)
}

type Server struct {
	e    *echo.Echo
	port int

	invoiceService  InvoiceService
	investorService InvestorService
	issuerService   IssuerService

	broker Broker
}

func New(port int, invoiceService InvoiceService, investorService InvestorService, issuerService IssuerService, broker Broker) *Server {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Pre(middleware.RemoveTrailingSlash())

	return &Server{
		e:               e,
		port:            port,
		invoiceService:  invoiceService,
		investorService: investorService,
		issuerService:   issuerService,
		broker:          broker,
	}
}

func (s *Server) Serve() error {
	s.issuerRoutes(s.e.Group("/issuer"))
	s.investorRoutes(s.e.Group("/investor"))
	s.invoiceRoutes(s.e.Group("/invoice"))

	go func() {
		if err := s.e.Start(fmt.Sprintf(":%d", s.port)); err != nil {
			log.Fatal(err)
		}
	}()

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.e.Shutdown(ctx)
}

func fmtBalance(balance currency.Amount) string {
	if balance.IsZero() {
		return ""
	}

	return currFmt.Format(balance)
}
