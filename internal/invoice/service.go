package invoice

import (
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"

	"github.com/bojanz/currency"
)

type Storage interface {
	SaveInvoice(context.Context, Invoice) error
	RetrieveInvoice(context.Context, string) (Invoice, error)
	RetrieveInvoicesByIssuerID(context.Context, string) ([]Invoice, error)
	UpdateStatus(context.Context, string, string) error

	SaveBid(context.Context, Bid) error
	RetrieveActiveBidsByInvoiceID(context.Context, string) ([]Bid, error)
	DisableBidsByInvoiceID(context.Context, string) error
}

type FileStorage interface {
	SaveFile(string, io.Reader) error
}

type Service struct {
	st  Storage
	fst FileStorage
}

func NewService(st Storage, fst FileStorage) *Service {
	return &Service{}
}

func (s *Service) GetInvoice(ctx context.Context, id string) (Invoice, error) {
	return s.st.RetrieveInvoice(ctx, id)
}

func (s *Service) GetByIssuerID(ctx context.Context, issID string) ([]Invoice, error) {
	return s.st.RetrieveInvoicesByIssuerID(ctx, issID)
}

func (s *Service) CreateInvoice(ctx context.Context, issuerID string, price currency.Amount, file io.Reader) (Invoice, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return Invoice{}, fmt.Errorf("could not generate id: %w", err)
	}

	invoice := Invoice{
		ID:       id.String(),
		IssuerID: issuerID,
		Price:    price,
		Status:   OPEN,
	}
	if err := s.st.SaveInvoice(ctx, invoice); err != nil {
		return Invoice{}, err
	}

	if err := s.fst.SaveFile(id.String(), file); err != nil {
		return Invoice{}, err
	}

	return invoice, nil
}

func (s *Service) PlaceBid(ctx context.Context, s3 string, s2 string, amount currency.Amount) (string, error) {

	panic("implement me")
}

func (s *Service) ApproveTrade(ctx context.Context, id string, approved bool) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetRemainingPrice(ctx context.Context, id string) (currency.Amount, error) {
	invoice, err := s.GetInvoice(ctx, id)
	if err != nil {
		return currency.Amount{}, err
	}

	remainingPrice := invoice.Price
	for _, b := range invoice.Bids {
		remainingPrice, _ = remainingPrice.Sub(b.Amount)
	}

	return remainingPrice, nil
}
