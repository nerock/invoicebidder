package broker

import "github.com/bojanz/currency"

type EventType string

const (
	TypeTradeEvent     EventType = "TradeEvent"
	TypeFailedBidEvent EventType = "TypeFailedBidEvent"
)

type Event interface {
	Type() EventType
	Resend()
	Retries() int
}

type TradeEvent struct {
	InvoiceID string
	Bids      []string
	Approved  bool
	r         int
}

func (te *TradeEvent) Type() EventType {
	return TypeTradeEvent
}

func (te *TradeEvent) Resend() {
	te.r++
}

func (te *TradeEvent) Retries() int {
	return te.r
}

type FailedBidEvent struct {
	InvestorID string
	Amount     currency.Amount
	r          int
}

func (be *FailedBidEvent) Type() EventType {
	return TypeFailedBidEvent
}

func (be *FailedBidEvent) Resend() {
	be.r++
}

func (be *FailedBidEvent) Retries() int {
	return be.r
}
