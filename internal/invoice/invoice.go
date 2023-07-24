package invoice

import (
	"github.com/bojanz/currency"
)

type Status string

const (
	OPEN   Status = "open"
	LOCKED Status = "locked"
	TRADED Status = "traded"
)

type Invoice struct {
	ID       string
	IssuerID string
	Price    currency.Amount
	Bids     []Bid
	Status   Status
}
