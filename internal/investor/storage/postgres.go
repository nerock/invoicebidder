package storage

import (
	"context"
	"fmt"

	"github.com/bojanz/currency"
	"github.com/jackc/pgx/v5"

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
	const queryInv = `INSERT INTO investors (id, name) VALUES ($1, $2)`
	const queryBalance = `INSERT INTO investor_balances (investor_id, amount) VALUES ($1, $2)`

	tx, err := s.c.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("could not start transaction")
	}

	if _, err := tx.Exec(ctx, queryInv, inv.ID, inv.FullName); err != nil {
		err = fmt.Errorf("could not save investor in db: %w", err)
		if errRB := tx.Rollback(ctx); err != nil {
			err = fmt.Errorf("%w: %w", err, errRB)
		}

		return err
	}

	batch := &pgx.Batch{}
	for _, b := range inv.Balances {
		batch.Queue(queryBalance, inv.ID, b)
	}

	if err := tx.SendBatch(ctx, batch).Close(); err != nil {
		err = fmt.Errorf("could not save balances in db: %w", err)
		if errRB := tx.Rollback(ctx); err != nil {
			err = fmt.Errorf("%w: %w", err, errRB)
		}

		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("could not commit transaction: %w", err)
	}

	return nil
}

func (s *Storage) RetrieveInvestor(ctx context.Context, s2 string) (investor.Investor, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) RetrieveInvestors(ctx context.Context, strings []string) ([]investor.Investor, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) UpdateBalance(ctx context.Context, s2 string, amount currency.Amount) error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) UpdateBalances(ctx context.Context, m map[string]currency.Amount) error {
	//TODO implement me
	panic("implement me")
}
