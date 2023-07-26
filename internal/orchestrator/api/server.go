package api

import (
	"context"
	"fmt"

	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/bojanz/currency"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/nerock/invoicebidder/docs"
)

var currFmt = currency.NewFormatter(currency.NewLocale("fr"))

type AmountRequest struct {
	Amount   string `json:"amount" example:"1200.50"`
	Currency string `json:"currency" example:"EUR"`
}

type Broker interface {
	SendTradeEvent(string, []string, bool)
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

// New creates a new server
// @title Invoice Bidder API
// @version 0.1
// @description Create invoices and bid on them
// @contact.name Manuel Adalid
// @contact.url https://manueladalid.dev
// @contact.email manueladalidmoya@gmail.com
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
	s.e.GET("/swagger/*", echoSwagger.WrapHandler)

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
