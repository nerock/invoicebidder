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
	FullName string        `json:"fullName" example:"Manuel Adalid"`
	Balance  AmountRequest `json:"balance"`
}

type InvestorBidResponse struct {
	ID      string             `json:"id" example:"343abd7a-874c-4bb7-ba7b-81e9c71cf1b0"`
	Amount  string             `json:"amount" example:"1 230,45 €"`
	Invoice InvoiceBidResponse `json:"invoice"`
}

type InvestorResponse struct {
	ID       string                `json:"id" example:"343abd7a-874c-4bb7-ba7b-81e9c71cf1b0"`
	FullName string                `json:"fullName" example:"Manuel Adalid"`
	Balance  string                `json:"balance,omitempty" example:"1 230,45 €"`
	Bids     []BidInvestorResponse `json:"bids,omitempty"`
}

type InvestorService interface {
	GetInvestor(context.Context, string) (investor.Investor, error)
	ListInvestors(context.Context, []string) (map[string]investor.Investor, error)
	CreateInvestor(context.Context, string, currency.Amount) (investor.Investor, error)
	Bid(context.Context, string, currency.Amount) error
}

func (s *Server) investorRoutes(g *echo.Group) {
	g.POST("", s.CreateInvestor)
	g.GET("", s.ListInvestors)
}

// CreateInvestor creates a new investor
// @Summary      New investor
// @Description  Create a new investor to bid on invoices
// @Tags         investor
// @Accept       json
// @Produce      json
// @Param request body CreateInvestorRequest true "Issuer request"
// @Success      201  {object}  InvestorResponse
// @Failure      400  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /investor [post]
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

// ListInvestors retrieves investors
// @Summary      List investors
// @Description  Retrieve investors optionally filtering by ids
// @Tags         investor
// @Accept       json
// @Produce      json
// @Param ids query []string false "list of comma separated ids for filtering"
// @Success      200  {array}   InvestorResponse
// @Failure      500  {object}  HTTPError
// @Router       /investor [get]
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
