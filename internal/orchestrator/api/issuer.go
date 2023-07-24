package api

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CreateIssuerRequest struct {
	FullName string `json:"fullName"`
}

type IssuerResponse struct {
	ID       string                  `json:"id"`
	FullName string                  `json:"fullName"`
	Balances []string                `json:"balances"`
	Invoices []IssuerInvoiceResponse `json:"invoices"`
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

func (s *Server) issuerRoutes(g *echo.Group) {
	g.POST("/", s.CreateIssuer)
	g.GET("/:id", s.RetrieveIssuer)
}

func (s *Server) CreateIssuer(c echo.Context) error {
	var req CreateIssuerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	ctx := c.Request().Context()
	iss, err := s.issuerService.CreateIssuer(ctx, req.FullName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("could not create investor: %w", err))
	}

	return c.JSON(http.StatusCreated, IssuerResponse{
		ID:       iss.ID,
		FullName: iss.FullName,
		Balances: balances(iss.Balances),
	})
}

func (s *Server) RetrieveIssuer(c echo.Context) error {
	id := c.Param("id")

	ctx := c.Request().Context()
	iss, err := s.issuerService.GetIssuer(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, err)
	}

	invoices, err := s.invoiceService.GetByIssuerID(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Errorf("could not retrieve issuer invoices: %w", err))
	}

	invoicesRes := make([]IssuerInvoiceResponse, 0, len(invoices))
	for _, inv := range invoices {
		bidsRes := make([]IssuerBidResponse, 0, len(inv.Bids))
		for _, bid := range inv.Bids {
			bidsRes = append(bidsRes, IssuerBidResponse{
				ID:     bid.ID,
				Amount: bid.Amount.String(),
			})
		}

		invoicesRes = append(invoicesRes, IssuerInvoiceResponse{
			ID:     inv.ID,
			Price:  inv.Price.String(),
			Status: string(inv.Status),
			Bids:   bidsRes,
		})
	}

	return c.JSON(http.StatusOK, IssuerResponse{
		ID:       iss.ID,
		FullName: iss.FullName,
		Balances: balances(iss.Balances),
		Invoices: invoicesRes,
	})
}
