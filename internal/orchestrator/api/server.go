package api

import (
	"context"
	"fmt"
	"io"

	"github.com/bojanz/currency"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/nerock/invoicebidder/internal/investor"
	"github.com/nerock/invoicebidder/internal/invoice"
	"github.com/nerock/invoicebidder/internal/issuer"
)

var currFmt = currency.NewFormatter(currency.NewLocale("fr"))

type InvoiceService interface {
	GetInvoice(context.Context, string) (invoice.Invoice, error)
	GetRemainingPrice(context.Context, string) (currency.Amount, error)
	GetByIssuerID(context.Context, string) ([]invoice.Invoice, error)
	CreateInvoice(context.Context, string, currency.Amount, io.Reader) (invoice.Invoice, error)
	PlaceBid(context.Context, string, string, currency.Amount) (string, error)
	ApproveTrade(context.Context, string, bool) error
}
type InvestorService interface {
	GetInvestor(context.Context, string) (investor.Investor, error)
	ListInvestors(context.Context, []string) (map[string]investor.Investor, error)
	CreateInvestor(context.Context, string, []currency.Amount) (investor.Investor, error)
	Bid(context.Context, string, currency.Amount, currency.Amount) error
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
	e := echo.New()
	e.Use(middleware.Logger())
	e.Pre(middleware.AddTrailingSlash())

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

func balances(bs []currency.Amount) []string {
	balances := make([]string, 0, len(bs))
	for _, b := range bs {
		balances = append(balances, currFmt.Format(b))
	}

	return balances
}
