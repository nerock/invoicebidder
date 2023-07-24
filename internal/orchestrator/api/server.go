package api

import (
	"context"
	"fmt"
	"github.com/bojanz/currency"
	"github.com/labstack/echo/v4"
	"github.com/nerock/invoicebidder/internal/investor"
	"github.com/nerock/invoicebidder/internal/invoice"
	"github.com/nerock/invoicebidder/internal/issuer"
	"io"
	"log"
)

type InvoiceService interface {
	GetInvoice(context.Context, string) (invoice.Invoice, error)
	GetByIssuerID(context.Context, string) ([]invoice.Invoice, error)
	CreateInvoice(context.Context, string, currency.Amount, io.Reader) (invoice.Invoice, error)
	PlaceBid(context.Context, string, string, currency.Amount) (string, error)
	ApproveTrade(context.Context, string, bool) error
}
type InvestorService interface {
	GetInvestor(context.Context, string) (investor.Investor, error)
	ListInvestors(context.Context, []string) (map[string]investor.Investor, error)
	CreateInvestor(context.Context, string, []currency.Amount) (investor.Investor, error)
	Bid(context.Context, string, currency.Amount) error
}

type IssuerService interface {
	GetIssuer(context.Context, string) (issuer.Issuer, error)
	CreateIssuer(context.Context, string) (issuer.Issuer, error)
}

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
	return &Server{
		e:               echo.New(),
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

func balances(bs []currency.Amount) []string {
	balances := make([]string, 0, len(bs))
	for _, b := range bs {
		balances = append(balances, b.String())
	}

	return balances
}
