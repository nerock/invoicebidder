package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/bojanz/currency"
	"github.com/labstack/echo/v4"
	"github.com/nerock/invoicebidder/internal/investor"
)

type CreateInvestorRequest struct {
	FullName string                `json:"fullName"`
	Balance  InvestorAmountRequest `json:"balance"`
}

type UpdateInvestorRequest struct {
	ID           string `json:"id"`
	BalanceDelta string `json:"balance"`
}

type InvestorBidResponse struct {
	ID      string             `json:"id"`
	Amount  string             `json:"amount"`
	Invoice InvoiceBidResponse `json:"invoice"`
}

type BidInvoiceResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type InvestorResponse struct {
	ID       string                `json:"id"`
	FullName string                `json:"fullName"`
	Balance  string                `json:"balance,omitempty"`
	Bids     []BidInvestorResponse `json:"bids,omitempty"`
}

type InvestorAmountRequest struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type InvestorService interface {
	GetInvestor(context.Context, string) (investor.Investor, error)
	ListInvestors(context.Context, []string) (map[string]investor.Investor, error)
	CreateInvestor(context.Context, string, currency.Amount) (investor.Investor, error)
	Bid(context.Context, string, currency.Amount, currency.Amount) error
}

func (s *Server) investorRoutes(g *echo.Group) {
	g.POST("/", s.CreateInvestor)
	g.GET("/:id", s.ListInvestors)
}

func (s *Server) CreateInvestor(c echo.Context) error {
	var req CreateInvestorRequest
	if err := c.Bind(&req); err != nil {
		return errBadRequest(err, c)
	}

	balance, err := currency.NewAmount(req.Balance.Amount, req.Balance.Currency)
	if err != nil {
		return errBadRequest(fmt.Errorf("invalid balance: %w", err), c)
	}

	ctx := c.Request().Context()
	inv, err := s.investorService.CreateInvestor(ctx, req.FullName, balance)
	if err != nil {
		return errHandler(err, c)
	}

	return c.JSON(http.StatusCreated, InvestorResponse{
		ID:       inv.ID,
		FullName: inv.FullName,
		Balance:  fmtBalance(balance),
	})
}

func (s *Server) ListInvestors(c echo.Context) error {
	param := c.QueryParam("ids")
	var ids []string
	if param != "" {
		ids = strings.Split(param, ",")
	}

	ctx := c.Request().Context()
	invMap, err := s.investorService.ListInvestors(ctx, ids)
	if err != nil {
		return errHandler(err, c)
	}

	res := make([]InvestorResponse, 0, len(invMap))
	for _, inv := range invMap {
		res = append(res, InvestorResponse{
			ID:       inv.ID,
			FullName: inv.FullName,
			Balance:  fmtBalance(inv.Balance),
		})
	}

	return c.JSON(http.StatusOK, res)
}
