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

func (s *Service) CreateInvestor(ctx context.Context, name string, balance currency.Amount) (Investor, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return Investor{}, fmt.Errorf("could not generate id: %w", err)
	}

	investor := Investor{
		ID:       id.String(),
		FullName: name,
		Balance:  balance,
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

func (s *Service) Bid(ctx context.Context, id string, amount currency.Amount) error {
	investor, err := s.GetInvestor(ctx, id)
	if err != nil {
		return err
	}

	amount, _ = amount.Mul("-1")
	newBalance, err := addBalance(investor.Balance, amount)
	if err != nil {
		return err
	}

	if newBalance.IsNegative() {
		return fmt.Errorf("insufficient funds: %w", err)
	}

	return s.st.UpdateBalance(ctx, id, newBalance)
}

func (s *Service) CancelBid(ctx context.Context, id string, amount currency.Amount) error {
	investor, err := s.GetInvestor(ctx, id)
	if err != nil {
		return err
	}

	newBalance, err := addBalance(investor.Balance, amount)
	if err != nil {
		return err
	}

	return s.st.UpdateBalance(ctx, id, newBalance)
}

func (s *Service) CancelTrade(ctx context.Context, bids []Bid) error {
	investorIDs := make([]string, 0, len(bids))
	for _, b := range bids {
		investorIDs = append(investorIDs, b.InvestorID)
	}

	investors, err := s.st.RetrieveInvestors(ctx, investorIDs)
	if err != nil {
		return err
	}

	newBalances := make(map[string]currency.Amount)
	for _, inv := range investors {
		newBalance := inv.Balance
		for _, b := range bids {
			if b.InvestorID == inv.ID {
				var err error
				newBalance, err = addBalance(newBalance, b.Amount)
				if err != nil {
					return err
				}
			}
		}

		newBalances[inv.ID] = newBalance
	}

	return s.st.UpdateBalances(ctx, newBalances)
}

func addBalance(current, delta currency.Amount) (currency.Amount, error) {
	if delta.CurrencyCode() != current.CurrencyCode() {
		var err error
		delta, err = delta.Convert(current.CurrencyCode(), "1")
		if err != nil {
			return current, err
		}
	}

	newBalance, err := current.Add(delta)
	if err != nil {
		return current, fmt.Errorf("could not perform currency operation: %w", err)
	}

	return newBalance, nil
}
