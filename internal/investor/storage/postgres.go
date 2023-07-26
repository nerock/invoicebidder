package storage

import (
	"context"
	"fmt"

	"github.com/bojanz/currency"
	"github.com/nerock/invoicebidder/internal/investor"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	c *pgxpool.Pool
}

func New(c *pgxpool.Pool) *Storage {
	return &Storage{c}
}

func (s *Storage) CreateInvestor(ctx context.Context, inv investor.Investor) error {
	const query = `INSERT INTO investors (id, name, balance) VALUES ($1, $2, $3)`

	if _, err := s.c.Exec(ctx, query, inv.ID, inv.FullName, inv.Balance); err != nil {
		return fmt.Errorf("could not save investor in db: %w", err)
	}

	return nil
}

func (s *Storage) RetrieveInvestor(ctx context.Context, id string) (investor.Investor, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) RetrieveInvestors(ctx context.Context, ids []string) ([]investor.Investor, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) UpdateBalance(ctx context.Context, id string, balance currency.Amount) error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) UpdateBalances(ctx context.Context, m map[string]currency.Amount) error {
	//TODO implement me
	panic("implement me")
}
