package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

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
	const query = `SELECT i.name, i.balance FROM investors i WHERE i.id = $1`

	inv := investor.Investor{ID: id}
	err := s.c.QueryRow(ctx, query, id).Scan(&inv.FullName, &inv.Balance)
	if err != nil {
		return inv, fmt.Errorf("could not retrieve investor: %w", err)
	}

	return inv, nil
}

func (s *Storage) RetrieveInvestors(ctx context.Context, ids []string) ([]investor.Investor, error) {
	const query = `SELECT i.id, i.name, i.balance FROM investors i`
	const withIDs = `WHERE i.id = any($1)`

	var rows pgx.Rows
	var err error
	if len(ids) == 0 {
		rows, err = s.c.Query(ctx, query)
	} else {
		rows, err = s.c.Query(ctx, fmt.Sprintf("%s %s", query, withIDs), ids)
	}

	if err != nil {
		return nil, fmt.Errorf("could not retrieve investors: %w", err)
	}

	var investors []investor.Investor
	for rows.Next() {
		var inv investor.Investor
		if err := rows.Scan(&inv.ID, &inv.FullName, &inv.Balance); err != nil {
			return nil, fmt.Errorf("could not scan investors: %w", err)
		}

		investors = append(investors, inv)
	}

	return investors, nil
}

func (s *Storage) UpdateBalance(ctx context.Context, id string, balance currency.Amount) error {
	const query = `UPDATE investors SET balance = $1 WHERE id = $2`

	if _, err := s.c.Exec(ctx, query, balance, id); err != nil {
		return fmt.Errorf("could not replace current investor balance in db: %w", err)
	}

	return nil
}

func (s *Storage) UpdateBalances(ctx context.Context, balances map[string]currency.Amount) error {
	const query = `UPDATE investors SET balance = $1 WHERE id = $2`

	tx, err := s.c.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("could not initialize transaction: %w", err)
	}

	batch := &pgx.Batch{}
	for i, b := range balances {
		batch.Queue(query, b, i)
	}

	if err := tx.SendBatch(ctx, batch).Close(); err != nil {
		err = fmt.Errorf("could not save investor balances in db: %w", err)
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
