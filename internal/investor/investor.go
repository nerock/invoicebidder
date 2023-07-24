package investor

import (
	"github.com/bojanz/currency"
)

type Investor struct {
	ID       string
	FullName string
	Bids     []string
	Balances []currency.Amount
}