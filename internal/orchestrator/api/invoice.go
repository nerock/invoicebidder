package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/bojanz/currency"
	"github.com/labstack/echo/v4"
	"github.com/nerock/invoicebidder/internal/invoice"
)

type InvoiceResponse struct {
	ID     string                `json:"id"`
	Price  string                `json:"price"`
	Status string                `json:"status"`
	Issuer InvoiceIssuerResponse `json:"issuer"`
	Bids   []InvoiceBidResponse  `json:"bids,omitempty"`
}

type InvoiceIssuerResponse struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
}

type InvoiceBidResponse struct {
	ID       string `json:"id"`
	Amount   string `json:"string"`
	Investor BidInvestorResponse
}

type BidInvestorResponse struct {
	ID       string `json:"id"`
	FullName string `json:"fullName"`
}

type ApproveTradeRequest struct {
	InvoiceID string `json:"invoiceId"`
	Approved  bool   `json:"approve"`
}

type BidRequest struct {
	InvestorID string                `json:"investorId"`
	Amount     InvestorAmountRequest `json:"amount"`
}

type InvoiceService interface {
	GetInvoice(context.Context, string) (invoice.Invoice, error)
	GetRemainingPrice(context.Context, string) (currency.Amount, error)
	GetByIssuerID(context.Context, string) ([]invoice.Invoice, error)
	CreateInvoice(context.Context, string, currency.Amount, io.Reader) (invoice.Invoice, error)
	PlaceBid(context.Context, string, string, currency.Amount) (string, error)
	ApproveTrade(context.Context, string, bool) error
}

func (s *Server) invoiceRoutes(g *echo.Group) {
	g.POST("/", s.CreateInvoice)
	g.GET("/:id", s.RetrieveInvoice)
	g.POST("/:id/bid", s.Bid)
	g.POST("/:id/trade", s.ApproveTrade)
}

func (s *Server) CreateInvoice(c echo.Context) error {
	issID := c.FormValue("issuer_id")
	if issID == "" {
		return errBadRequest(fmt.Errorf("issuer id cannot be empty"), c)
	}
	price := c.FormValue("price")
	if issID == "" {
		return errBadRequest(fmt.Errorf("price cannot be empty"), c)
	}
	curr := c.FormValue("currency")
	if issID == "" {
		return errBadRequest(fmt.Errorf("currency cannot be empty"), c)
	}

	amount, err := currency.NewAmount(price, curr)
	if err != nil {
		return errBadRequest(err, c)
	}

	formFile, err := c.FormFile("invoice")
	if err != nil {
		return errBadRequest(fmt.Errorf("could not read invoice file: %w", err), c)
	}
	file, err := formFile.Open()
	if err != nil {
		return errBadRequest(fmt.Errorf("could not open invoice file: %w", err), c)
	}
	defer func() {
		if err := file.Close(); err != nil {
			s.e.Logger.Errorf("could not close request file: %w", err)
		}
	}()

	ctx := c.Request().Context()
	iss, err := s.issuerService.GetIssuer(ctx, issID)
	if err != nil {
		return errHandler(err, c)
	}

	inv, err := s.invoiceService.CreateInvoice(ctx, iss.ID, amount, file)
	if err != nil {
		return errHandler(err, c)
	}

	return c.JSON(http.StatusCreated, InvoiceResponse{
		ID:     inv.ID,
		Price:  currFmt.Format(inv.Price),
		Status: string(inv.Status),
		Issuer: InvoiceIssuerResponse{
			ID:       iss.ID,
			FullName: iss.FullName,
		},
		Bids: nil,
	})
}

func (s *Server) RetrieveInvoice(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return errBadRequest(errors.New("id cannot be empty"), c)
	}

	ctx := c.Request().Context()
	inv, err := s.invoiceService.GetInvoice(ctx, id)
	if err != nil {
		return errHandler(err, c)
	}

	iss, err := s.issuerService.GetIssuer(ctx, inv.IssuerID)
	if err != nil {
		return errHandler(err, c)
	}

	res := InvoiceResponse{
		ID:     inv.ID,
		Price:  currFmt.Format(inv.Price),
		Status: string(inv.Status),
		Issuer: InvoiceIssuerResponse{
			ID:       iss.ID,
			FullName: iss.FullName,
		},
	}

	if len(inv.Bids) > 0 {
		res.Bids = make([]InvoiceBidResponse, 0, len(inv.Bids))

		investorsIDs := make([]string, 0, len(inv.Bids))
		for _, b := range inv.Bids {
			investorsIDs = append(investorsIDs, b.InvestorID)
		}

		investors, err := s.investorService.ListInvestors(ctx, investorsIDs)
		if err != nil {
			return errHandler(err, c)
		}

		for _, b := range inv.Bids {
			inv, ok := investors[b.InvestorID]
			if !ok {
				return fmt.Errorf("missing investor info: %w", ErrNotFound)
			}

			res.Bids = append(res.Bids, InvoiceBidResponse{
				ID:     b.ID,
				Amount: currFmt.Format(b.Amount),
				Investor: BidInvestorResponse{
					ID:       inv.ID,
					FullName: inv.FullName,
				},
			})
		}
	}

	return c.JSON(http.StatusOK, res)
}

func (s *Server) Bid(c echo.Context) error {
	var req BidRequest
	if err := c.Bind(&req); err != nil {
		return errBadRequest(err, c)
	}
	bidAmount, err := currency.NewAmount(req.Amount.Amount, req.Amount.Currency)
	if err != nil {
		return errBadRequest(err, c)
	}

	invoiceID := c.Param("id")
	if invoiceID == "" {
		return errBadRequest(errors.New("id cannot be empty"), c)
	}

	ctx := c.Request().Context()
	investor, err := s.investorService.GetInvestor(ctx, req.InvestorID)
	if err != nil {
		return errHandler(err, c)
	}

	remainingPrice, err := s.invoiceService.GetRemainingPrice(ctx, invoiceID)
	if err != nil {
		return errHandler(err, c)
	}

	if cmp, err := compareAmounts(remainingPrice, bidAmount); err != nil {
		return errHandler(err, c)
	} else if cmp < 0 {
		bidAmount = remainingPrice
	}

	if err := s.investorService.Bid(ctx, investor.ID, bidAmount); err != nil {
		return errHandler(err, c)
	}

	id, err := s.invoiceService.PlaceBid(ctx, invoiceID, investor.ID, bidAmount)
	if err != nil {
		s.broker.SendFailedBidEvent(investor.ID, bidAmount)
		return errHandler(err, c)
	}

	return c.JSON(http.StatusCreated, InvoiceBidResponse{
		ID:     id,
		Amount: currFmt.Format(bidAmount),
		Investor: BidInvestorResponse{
			ID:       investor.ID,
			FullName: investor.FullName,
		},
	})
}

func (s *Server) ApproveTrade(c echo.Context) error {
	var req ApproveTradeRequest
	if err := c.Bind(&req); err != nil {
		return errBadRequest(err, c)
	}

	ctx := c.Request().Context()
	if err := s.invoiceService.ApproveTrade(ctx, req.InvoiceID, req.Approved); err != nil {
		return errHandler(err, c)
	}

	s.broker.SendTradeEvent(req.InvoiceID, req.Approved)

	return c.NoContent(http.StatusOK)
}

func compareAmounts(a, b currency.Amount) (int, error) {
	if a.CurrencyCode() != b.CurrencyCode() {
		var err error
		b, err = b.Convert(b.CurrencyCode(), "1")
		if err != nil {
			return 0, err
		}
	}

	cmp, err := a.Cmp(b)
	if err != nil {
		return 0, err
	}

	return cmp, nil
}
