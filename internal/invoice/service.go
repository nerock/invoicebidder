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
	RetrieveBidsByIDs(context.Context, []string) ([]Bid, error)
	UpdateStatus(context.Context, string, Status) error

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
	return &Service{
		st:  st,
		fst: fst,
	}
}

func (s *Service) GetInvoice(ctx context.Context, id string) (Invoice, error) {
	return s.st.RetrieveInvoice(ctx, id)
}

func (s *Service) GetByIssuerID(ctx context.Context, issID string) ([]Invoice, error) {
	return s.st.RetrieveInvoicesByIssuerID(ctx, issID)
}

func (s *Service) ListBidsByIDs(ctx context.Context, bidsIDs []string) ([]Bid, error) {
	return s.st.RetrieveBidsByIDs(ctx, bidsIDs)
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

func (s *Service) PlaceBid(ctx context.Context, invoiceID string, investorID string, amount currency.Amount) (string, error) {
	invoice, err := s.GetInvoice(ctx, invoiceID)
	if err != nil {
		return "", err
	}

	if invoice.Status != OPEN {
		return "", fmt.Errorf("can only place bids in open invoices")
	}

	id, err := uuid.NewUUID()
	if err != nil {
		return "", fmt.Errorf("could not generate id: %w", err)
	}

	if err := s.st.SaveBid(ctx, Bid{
		ID:         id.String(),
		InvestorID: investorID,
		InvoiceID:  invoiceID,
		Amount:     amount,
		Active:     true,
	}); err != nil {
		return "", err
	}

	if remaining, _ := s.getRemainingPrice(invoice).Sub(amount); remaining.IsZero() {
		if err := s.st.UpdateStatus(ctx, invoiceID, LOCKED); err != nil {
			return "", err
		}
	}

	return id.String(), nil
}

func (s *Service) ApproveTrade(ctx context.Context, id string, approved bool) ([]string, error) {
	invoice, err := s.GetInvoice(ctx, id)
	if err != nil {
		return nil, err
	}

	if invoice.Status != LOCKED {
		return nil, fmt.Errorf("cannot close trade if the status is not locked")
	}

	if approved {
		if err := s.st.UpdateStatus(ctx, id, TRADED); err != nil {
			return nil, err
		}
	} else {
		if err := s.st.UpdateStatus(ctx, id, OPEN); err != nil {
			return nil, err
		}

		if err := s.st.DisableBidsByInvoiceID(ctx, id); err != nil {
			return nil, err
		}
	}

	bidsIDs := make([]string, 0, len(invoice.Bids))
	for _, b := range invoice.Bids {
		bidsIDs = append(bidsIDs, b.ID)
	}

	return bidsIDs, nil
}

func (s *Service) GetRemainingPrice(ctx context.Context, id string) (currency.Amount, error) {
	invoice, err := s.GetInvoice(ctx, id)
	if err != nil {
		return currency.Amount{}, err
	}

	return s.getRemainingPrice(invoice), nil
}

func (s *Service) getRemainingPrice(invoice Invoice) currency.Amount {
	remainingPrice := invoice.Price
	for _, b := range invoice.Bids {
		remainingPrice, _ = remainingPrice.Sub(b.Amount)
	}

	return remainingPrice
}
