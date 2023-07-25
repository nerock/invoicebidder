package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nerock/invoicebidder/internal/invoice"
)

type Storage struct {
	c *pgxpool.Pool
}

func New(c *pgxpool.Pool) *Storage {
	return &Storage{c}
}

func (s *Storage) SaveInvoice(ctx context.Context, i invoice.Invoice) error {
	const query = `INSERT INTO invoices (id, issuer_id, price, status) VALUES ($1, $2, $3, $4)`

	if _, err := s.c.Exec(ctx, query, i.ID, i.IssuerID, i.Price, i.Status); err != nil {
		return fmt.Errorf("could not save invoice in db: %w", err)
	}

	return nil
}

func (s *Storage) RetrieveInvoice(ctx context.Context, id string) (invoice.Invoice, error) {
	const query = `SELECT i.issuer_id, i.price, i.status FROM invoices i WHERE i.id = $1`

	inv := invoice.Invoice{ID: id}
	err := s.c.QueryRow(ctx, query, id).Scan(&inv.IssuerID, &inv.Price, &inv.Status)
	if err != nil {
		return inv, fmt.Errorf("could not retrieve invoice: %w", err)
	}

	inv.Bids, err = s.RetrieveActiveBidsByInvoiceID(ctx, id)
	if err != nil {
		return inv, fmt.Errorf("could not retrieve invoice bids: %w", err)
	}

	return inv, nil
}

func (s *Storage) RetrieveInvoicesByIssuerID(ctx context.Context, issID string) ([]invoice.Invoice, error) {
	const query = `SELECT i.id, i.price, i.status FROM invoices i WHERE i.issuer_id = $1`

	rows, err := s.c.Query(ctx, query, issID)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve bids: %w", err)
	}

	var invoices []invoice.Invoice
	for rows.Next() {
		inv := invoice.Invoice{IssuerID: issID}

		err := rows.Scan(&inv.ID, &inv.Price, &inv.Status)
		if err != nil {
			return nil, fmt.Errorf("could not scan bids: %w", err)
		}

		inv.Bids, err = s.RetrieveActiveBidsByInvoiceID(ctx, inv.ID)
		if err != nil {
			return nil, fmt.Errorf("could not retrieve invoice bids: %w", err)
		}

		invoices = append(invoices, inv)
	}

	return invoices, nil
}

func (s *Storage) UpdateStatus(ctx context.Context, id string, status string) error {
	const query = `UPDATE invoices SET status = $2 WHERE i.id = $1`

	if _, err := s.c.Exec(ctx, query, id, status); err != nil {
		return fmt.Errorf("could not update invoice status in db: %w", err)
	}

	return nil
}

func (s *Storage) SaveBid(ctx context.Context, b invoice.Bid) error {
	const query = `INSERT INTO bids (id, invoice_id, investor_id, amount, active) VALUES ($1, $2, $3, $4, $5)`

	if _, err := s.c.Exec(ctx, query, b.ID, b.InvoiceID, b.InvestorID, b.Amount, b.Active); err != nil {
		return fmt.Errorf("could not save bid in db: %w", err)
	}

	return nil
}

func (s *Storage) RetrieveActiveBidsByInvoiceID(ctx context.Context, invoiceID string) ([]invoice.Bid, error) {
	const query = `SELECT b.id, b.investor_id, b.amount FROM bids b WHERE b.invoice_id = $1 AND b.active = true`

	rows, err := s.c.Query(ctx, query, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve bids: %w", err)
	}

	var bids []invoice.Bid
	for rows.Next() {
		bid := invoice.Bid{InvoiceID: invoiceID, Active: true}
		if err := rows.Scan(&bid.ID, &bid.InvestorID, &bid.Amount); err != nil {
			return nil, fmt.Errorf("could not scan bids: %w", err)
		}

		bids = append(bids, bid)
	}

	return bids, nil
}

func (s *Storage) DisableBidsByInvoiceID(ctx context.Context, id string) error {
	const query = `UPDATE bids SET active = false WHERE invoice_id = $1 AND active = true`

	if _, err := s.c.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("could not disable bids in db: %w", err)
	}

	return nil
}
