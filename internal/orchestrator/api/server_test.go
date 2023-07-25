package api

import (
	"context"
	"io"

	"github.com/bojanz/currency"
	"github.com/nerock/invoicebidder/internal/invoice"
	"github.com/nerock/invoicebidder/internal/issuer"
)

type mockIssuerService struct {
	getIssuerFunc    func(context.Context, string) (issuer.Issuer, error)
	createIssuerFunc func(context.Context, string) (issuer.Issuer, error)
}

func (m *mockIssuerService) GetIssuer(ctx context.Context, id string) (issuer.Issuer, error) {
	return m.getIssuerFunc(ctx, id)
}

func (m *mockIssuerService) CreateIssuer(ctx context.Context, name string) (issuer.Issuer, error) {
	return m.createIssuerFunc(ctx, name)
}

type mockInvoiceService struct {
	getByIssuerIDFunc func(context.Context, string) ([]invoice.Invoice, error)
}

func (m *mockInvoiceService) GetInvoice(ctx context.Context, s string) (invoice.Invoice, error) {
	panic("implement me")
}

func (m *mockInvoiceService) GetRemainingPrice(ctx context.Context, s string) (currency.Amount, error) {
	panic("implement me")
}

func (m *mockInvoiceService) CreateInvoice(ctx context.Context, s string, amount currency.Amount, reader io.Reader) (invoice.Invoice, error) {
	panic("implement me")
}

func (m *mockInvoiceService) PlaceBid(ctx context.Context, s string, s2 string, amount currency.Amount) (string, error) {
	panic("implement me")
}

func (m *mockInvoiceService) ApproveTrade(ctx context.Context, s string, b bool) error {
	panic("implement me")
}

func (m *mockInvoiceService) GetByIssuerID(ctx context.Context, id string) ([]invoice.Invoice, error) {
	return m.getByIssuerIDFunc(ctx, id)
}
