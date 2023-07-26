package api

import (
	"context"
	"errors"
	"net/http"

	"github.com/nerock/invoicebidder/internal/issuer"

	"github.com/labstack/echo/v4"
)

type CreateIssuerRequest struct {
	FullName string `json:"fullName"`
}

type IssuerResponse struct {
	ID       string                  `json:"id"`
	FullName string                  `json:"fullName"`
	Balance  string                  `json:"balance,omitempty"`
	Invoices []IssuerInvoiceResponse `json:"invoices,omitempty"`
}

type IssuerInvoiceResponse struct {
	ID     string              `json:"id"`
	Price  string              `json:"price"`
	Status string              `json:"status"`
	Bids   []IssuerBidResponse `json:"bids,omitempty"`
}

type IssuerBidResponse struct {
	ID     string `json:"id"`
	Amount string `json:"string"`
}

type IssuerService interface {
	GetIssuer(context.Context, string) (issuer.Issuer, error)
	CreateIssuer(context.Context, string) (issuer.Issuer, error)
}

func (s *Server) issuerRoutes(g *echo.Group) {
	g.POST("", s.CreateIssuer)
	g.GET("/:id", s.RetrieveIssuer)
}

func (s *Server) CreateIssuer(c echo.Context) error {
	var req CreateIssuerRequest
	if err := c.Bind(&req); err != nil {
		return errBadRequest(err, c)
	}
	if req.FullName == "" {
		return errBadRequest(errors.New("issuer name cannot be empty"), c)
	}

	ctx := c.Request().Context()
	iss, err := s.issuerService.CreateIssuer(ctx, req.FullName)
	if err != nil {
		return errHandler(err, c)
	}

	return c.JSON(http.StatusCreated, IssuerResponse{
		ID:       iss.ID,
		FullName: iss.FullName,
	})
}

func (s *Server) RetrieveIssuer(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return errBadRequest(errors.New("id cannot be empty"), c)
	}

	ctx := c.Request().Context()
	iss, err := s.issuerService.GetIssuer(ctx, id)
	if err != nil {
		return errHandler(err, c)
	}

	invoices, err := s.invoiceService.GetByIssuerID(ctx, id)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return errHandler(err, c)
	}

	invoicesRes := make([]IssuerInvoiceResponse, 0, len(invoices))
	for _, inv := range invoices {
		bidsRes := make([]IssuerBidResponse, 0, len(inv.Bids))
		for _, bid := range inv.Bids {
			bidsRes = append(bidsRes, IssuerBidResponse{
				ID:     bid.ID,
				Amount: currFmt.Format(bid.Amount),
			})
		}

		invoicesRes = append(invoicesRes, IssuerInvoiceResponse{
			ID:     inv.ID,
			Price:  currFmt.Format(inv.Price),
			Status: string(inv.Status),
			Bids:   bidsRes,
		})
	}

	return c.JSON(http.StatusOK, IssuerResponse{
		ID:       iss.ID,
		FullName: iss.FullName,
		Balance:  fmtBalance(iss.Balance),
		Invoices: invoicesRes,
	})
}
