package bid

import "github.com/bojanz/currency"

type Status string

const (
	WAITING  Status = "waiting"
	APPROVED Status = "approved"
)

type Bid struct {
	ID       string
	Investor string
	Invoice  string
	Status   Status
	Amount   currency.Amount
}
