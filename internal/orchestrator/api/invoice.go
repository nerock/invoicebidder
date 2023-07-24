package api

import (
	"fmt"
	"github.com/bojanz/currency"
	"github.com/labstack/echo/v4"
	"net/http"
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
	InvestorID string `json:"investorId"`
	Amount     Amount `json:"amount"`
}

func (s *Server) invoiceRoutes(g *echo.Group) {
	g.POST("/", s.CreateInvoice)
	g.GET("/:id", s.RetrieveInvoice)
	g.POST("/:id/bid", s.Bid)
	g.POST("/:id/trade", s.ApproveTrade)
}

func (s *Server) CreateInvoice(c echo.Context) error {
	ctx := c.Request().Context()
	iss, err := s.issuerService.GetIssuer(ctx, c.FormValue("issuer_id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, fmt.Errorf("could not retrieve issuer"))
	}

	price := c.FormValue("price")
	curr := c.FormValue("currency")
	amount, err := currency.NewAmount(price, curr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Errorf("invalid price: %w", err))
	}

	formFile, err := c.FormFile("invoice")
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Errorf("could not read invoice file: %w", err))
	}
	file, err := formFile.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Errorf("could not open invoice file: %w", err))
	}
	defer func() {
		if err := file.Close(); err != nil {
			s.e.Logger.Errorf("could not close request file: %w", err)
		}
	}()

	inv, err := s.invoiceService.CreateInvoice(ctx, iss.ID, amount, file)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("could not create invoice: %w", err))
	}

	return c.JSON(http.StatusCreated, InvoiceResponse{
		ID:     inv.ID,
		Price:  inv.Price.String(),
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

	ctx := c.Request().Context()
	inv, err := s.invoiceService.GetInvoice(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	iss, err := s.issuerService.GetIssuer(ctx, inv.IssuerID)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	res := InvoiceResponse{
		ID:     inv.ID,
		Price:  inv.Price.String(),
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
			return c.JSON(http.StatusNotFound, err)
		}

		for _, b := range inv.Bids {
			inv, ok := investors[b.InvestorID]
			if !ok {
				return c.JSON(http.StatusNotFound, "missing investor info")
			}

			res.Bids = append(res.Bids, InvoiceBidResponse{
				ID:     b.ID,
				Amount: b.Amount.String(),
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
		return c.JSON(http.StatusBadRequest, err)
	}
	bidAmount, err := currency.NewAmount(req.Amount.Amount, req.Amount.Currency)
	if err != nil {
		return c.JSON(http.StatusBadRequest, fmt.Errorf("invalid bid amount: %w", err))
	}

	invoiceID := c.Param("id")
	if invoiceID == "" {
		return c.JSON(http.StatusBadRequest, fmt.Errorf("invoice id cannot be empty"))
	}

	ctx := c.Request().Context()
	investor, err := s.investorService.GetInvestor(ctx, req.InvestorID)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	var investorBalance *currency.Amount
	for _, b := range investor.Balances {
		if b.CurrencyCode() == bidAmount.CurrencyCode() {
			investorBalance = &b
			break
		}
	}
	if investorBalance == nil {
		return c.JSON(http.StatusBadRequest, fmt.Errorf("no funds in the bidding currency: %s", bidAmount.CurrencyCode()))
	}

	invoice, err := s.invoiceService.GetInvoice(ctx, invoiceID)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	remainingPrice := invoice.Price
	for _, b := range invoice.Bids {
		remainingPrice, _ = remainingPrice.Sub(b.Amount)
	}

	if cmp, _ := bidAmount.Cmp(remainingPrice); cmp < 0 {
		bidAmount = remainingPrice
	}

	if cmp, _ := bidAmount.Cmp(*investorBalance); cmp < 0 {
		return c.JSON(http.StatusBadRequest, fmt.Errorf("insufficient funds"))
	}

	if err := s.investorService.Bid(ctx, investor.ID, bidAmount); err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("could not retire funds from investor"))
	}

	id, err := s.invoiceService.PlaceBid(ctx, invoiceID, investor.ID, bidAmount)
	if err != nil {
		s.broker.SendFailedBidEvent(investor.ID, bidAmount)
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("could not place bid: %w", err))
	}

	return c.JSON(http.StatusCreated, InvoiceBidResponse{
		ID:     id,
		Amount: bidAmount.String(),
		Investor: BidInvestorResponse{
			ID:       investor.ID,
			FullName: investor.FullName,
		},
	})
}

func (s *Server) ApproveTrade(c echo.Context) error {
	var req ApproveTradeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	if err := s.invoiceService.ApproveTrade(ctx, req.InvoiceID, req.Approved); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	s.broker.SendTradeEvent(req.InvoiceID, req.Approved)

	return c.NoContent(http.StatusOK)
}
