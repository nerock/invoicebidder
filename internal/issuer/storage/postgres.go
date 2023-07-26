package storage

import (
	"context"
	"fmt"

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
	const query = `INSERT INTO issuers (id, name, balance) VALUES ($1, $2, $3)`

	if _, err := s.c.Exec(ctx, query, issuer.ID, issuer.FullName, issuer.Balance); err != nil {
		return fmt.Errorf("could not save issuer in db: %w", err)
	}

	return nil
}

func (s *Storage) RetrieveIssuer(ctx context.Context, id string) (issuer.Issuer, error) {
	const query = `SELECT i.name, i.balance FROM issuers i WHERE i.id = $1`

	iss := issuer.Issuer{ID: id}
	err := s.c.QueryRow(ctx, query, id).Scan(&iss.FullName, &iss.Balance)
	if err != nil {
		return iss, fmt.Errorf("could not retrieve invoice: %w", err)
	}

	return iss, nil
}

func (s *Storage) UpdateBalance(ctx context.Context, id string, balance currency.Amount) error {
	const query = `UPDATE issuers SET balance = $1 WHERE id = $2`

	if _, err := s.c.Exec(ctx, query, balance, id); err != nil {
		return fmt.Errorf("could not replace current balance in db: %w", err)
	}

	return nil
}
