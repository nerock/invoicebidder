package invoice

import "github.com/bojanz/currency"

type Bid struct {
	ID         string
	InvestorID string
	InvoiceID  string
	Amount     currency.Amount
	Active     bool
}
