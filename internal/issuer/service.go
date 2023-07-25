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
	UpdateBalances(context.Context, []currency.Amount) error
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

	issuer := Issuer{
		ID:       id.String(),
		FullName: name,
	}
	if err := s.st.CreateIssuer(ctx, issuer); err != nil {
		return Issuer{}, nil
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

	var ok bool
	for i, b := range issuer.Balances {
		if b.CurrencyCode() == amount.CurrencyCode() {
			ok = true
			newBalance, err := b.Add(amount)
			if err != nil {
				return fmt.Errorf("could not add funds to issuer balance: %w", err)
			}

			issuer.Balances[i] = newBalance
		}
	}

	if !ok {
		issuer.Balances = append(issuer.Balances, amount)
	}

	return s.st.UpdateBalances(ctx, issuer.Balances)
}
