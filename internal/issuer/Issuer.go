package issuer

import (
	"github.com/bojanz/currency"
)

type Issuer struct {
	ID       string
	FullName string
	Balances []currency.Amount
}
