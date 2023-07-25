package investor

import (
	"context"
	"fmt"

	"github.com/bojanz/currency"
	"github.com/google/uuid"
)

type Storage interface {
	CreateInvestor(context.Context, Investor) error
	RetrieveInvestor(context.Context, string) (Investor, error)
	RetrieveInvestors(context.Context, []string) ([]Investor, error)
	UpdateBalance(context.Context, string, currency.Amount) error
	UpdateBalances(context.Context, map[string]currency.Amount) error
}

type Service struct {
	st Storage
}

func NewService(st Storage) *Service {
	return &Service{
		st: st,
	}
}

func (s *Service) CreateInvestor(ctx context.Context, name string, balances []currency.Amount) (Investor, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return Investor{}, fmt.Errorf("could not generate id: %w", err)
	}

	investor := Investor{
		ID:       id.String(),
		FullName: name,
		Balances: balances,
	}
	if err := s.st.CreateInvestor(ctx, investor); err != nil {
		return Investor{}, nil
	}

	return investor, nil
}

func (s *Service) GetInvestor(ctx context.Context, id string) (Investor, error) {
	return s.st.RetrieveInvestor(ctx, id)
}

func (s *Service) ListInvestors(ctx context.Context, ids []string) (map[string]Investor, error) {
	investors, err := s.st.RetrieveInvestors(ctx, ids)
	if err != nil {
		return nil, err
	}

	investorsMap := make(map[string]Investor, len(investors))
	for _, inv := range investors {
		investorsMap[inv.ID] = inv
	}

	return investorsMap, nil
}

func (s *Service) Bid(ctx context.Context, id string, amount currency.Amount, amount2 currency.Amount) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) CancelBid(ctx context.Context, id string, amount currency.Amount) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) CancelTrade(ctx context.Context, investors map[string]currency.Amount) error {
	//TODO implement me
	panic("implement me")
}