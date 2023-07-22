package invoice

import (
	"github.com/bojanz/currency"
)

type Status string

const (
	BIDDING Status = "bidding"
	LOCKED  Status = "locked"
	TRADED  Status = "traded"
)

type Invoice struct {
	ID       string
	IssuerID string
	Price    currency.Amount
	Status   Status
}
