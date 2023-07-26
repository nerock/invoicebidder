package issuer

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/bojanz/currency"
)

type Storage interface {
	CreateIssuer(context.Context, Issuer) error
	RetrieveIssuer(context.Context, string) (Issuer, error)
	UpdateBalance(context.Context, string, currency.Amount) error
}

func NewService(st Storage) *Service {
	return &Service{
		st: st,
	}
}

type Service struct {
	st Storage
}

func (s *Service) CreateIssuer(ctx context.Context, name string) (Issuer, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return Issuer{}, fmt.Errorf("could not generate id: %w", err)
	}

	b, _ := currency.NewAmountFromInt64(0, "EUR")
	issuer := Issuer{
		ID:       id.String(),
		FullName: name,
		Balance:  b,
	}
	if err := s.st.CreateIssuer(ctx, issuer); err != nil {
		return Issuer{}, err
	}

	return issuer, nil
}

func (s *Service) GetIssuer(ctx context.Context, id string) (Issuer, error) {
	return s.st.RetrieveIssuer(ctx, id)
}

func (s *Service) ApproveTrade(ctx context.Context, id string, amount currency.Amount) error {
	issuer, err := s.GetIssuer(ctx, id)
	if err != nil {
		return err
	}

	if issuer.Balance.IsZero() {
		return s.st.UpdateBalance(ctx, id, amount)
	}

	convertedAmount, err := amount.Convert(issuer.Balance.CurrencyCode(), "1")
	if err != nil {
		return fmt.Errorf("could not convert amount: %w", err)
	}

	total, _ := convertedAmount.Add(issuer.Balance)
	return s.st.UpdateBalance(ctx, id, total)
}
