package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/bojanz/currency"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nerock/invoicebidder/internal/issuer"
)

type Storage struct {
	c *pgxpool.Pool
}

func New(c *pgxpool.Pool) *Storage {
	return &Storage{c}
}

func (s *Storage) CreateIssuer(ctx context.Context, issuer issuer.Issuer) error {
	const query = `INSERT INTO issuers (id, name) VALUES ($1, $2)`

	if _, err := s.c.Exec(ctx, query, issuer.ID, issuer.FullName); err != nil {
		return fmt.Errorf("could not save issuer in db: %w", err)
	}

	return nil
}

func (s *Storage) RetrieveIssuer(ctx context.Context, id string) (issuer.Issuer, error) {
	const query = `SELECT i.name FROM issuers i WHERE i.id = $1`

	iss := issuer.Issuer{ID: id}
	err := s.c.QueryRow(ctx, query, id).Scan(&iss.FullName)
	if err != nil {
		return iss, fmt.Errorf("could not retrieve invoice: %w", err)
	}

	iss.Balances, err = s.RetrieveBalances(ctx, id)
	if err != nil {
		return iss, fmt.Errorf("could not retrieve issuer balances: %w", err)
	}

	return iss, nil
}

func (s *Storage) UpdateBalances(ctx context.Context, id string, balances []currency.Amount) error {
	const removeQuery = `DELETE FROM issuer_balances WHERE issuer_id = $1`
	const insertQuery = `INSERT INTO issuer_balances (issuer_id, amount) VALUES ($1, $2)`

	tx, err := s.c.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("could not start transaction")
	}

	if _, err := tx.Exec(ctx, removeQuery, id); err != nil {
		err = fmt.Errorf("could not replace current balances in db: %w", err)
		if errRB := tx.Rollback(ctx); err != nil {
			err = fmt.Errorf("%w: %w", err, errRB)
		}

		return err
	}

	batch := &pgx.Batch{}
	for _, b := range balances {
		batch.Queue(insertQuery, id, b)
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

func (s *Storage) RetrieveBalances(ctx context.Context, id string) ([]currency.Amount, error) {
	const query = `SELECT b.amount FROM balances b WHERE b.issuer_id = $1`

	rows, err := s.c.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve balances: %w", err)
	}

	var balances []currency.Amount
	for rows.Next() {
		var b currency.Amount
		if err := rows.Scan(&b); err != nil {
			return nil, fmt.Errorf("could not scan balances: %w", err)
		}

		balances = append(balances, b)
	}

	return balances, nil
}
