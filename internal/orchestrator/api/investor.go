package api

import (
	"fmt"
	"github.com/bojanz/currency"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type CreateInvestorRequest struct {
	FullName string   `json:"fullName"`
	Balances []Amount `json:"balances"`
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
	Balances []string              `json:"balances,omitempty"`
	Bids     []BidInvestorResponse `json:"bids,omitempty"`
}

type Amount struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

func (s *Server) investorRoutes(g *echo.Group) {
	g.POST("/", s.CreateInvestor)
	g.GET("/:id", s.ListInvestors)
}

func (s *Server) CreateInvestor(c echo.Context) error {
	var req CreateInvestorRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	bs := make([]currency.Amount, 0, len(req.Balances))
	for _, balance := range req.Balances {
		b, err := currency.NewAmount(balance.Amount, balance.Currency)
		if err != nil {
			return c.JSON(http.StatusBadRequest, fmt.Errorf("invalid bid amount: %w", err))
		}

		bs = append(bs, b)
	}

	ctx := c.Request().Context()
	inv, err := s.investorService.CreateInvestor(ctx, req.FullName, bs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("could not create investor: %w", err))
	}

	return c.JSON(http.StatusCreated, InvestorResponse{
		ID:       inv.ID,
		FullName: inv.FullName,
		Balances: balances(bs),
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
		return c.JSON(http.StatusInternalServerError, err)
	}

	res := make([]InvestorResponse, 0, len(invMap))
	for _, inv := range invMap {
		res = append(res, InvestorResponse{
			ID:       inv.ID,
			FullName: inv.FullName,
			Balances: balances(inv.Balances),
		})
	}

	return c.JSON(http.StatusOK, res)
}
