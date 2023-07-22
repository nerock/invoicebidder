package issuer

import (
	"github.com/bojanz/currency"
)

type Issuer struct {
	ID       string
	FullName string
	Balance  currency.Amount
}
