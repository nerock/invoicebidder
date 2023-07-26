package investor

import (
	"github.com/bojanz/currency"
)

type Investor struct {
	ID       string
	FullName string
	Bids     []string
	Balance  currency.Amount
}

type Bid struct {
	InvestorID string
	Amount     currency.Amount
}
