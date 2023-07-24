package broker

import "github.com/bojanz/currency"

type EventType string

const (
	TypeTradeEvent     EventType = "TradeEvent"
	TypeFailedBidEvent EventType = "TypeFailedBidEvent"
)

type Event interface {
	Type() EventType
}

type TradeEvent struct {
	InvoiceID string
	Approved  bool
}

func (te TradeEvent) Type() EventType {
	return TypeTradeEvent
}

type FailedBidEvent struct {
	InvestorID string
	Amount     currency.Amount
}

func (be FailedBidEvent) Type() EventType {
	return TypeFailedBidEvent
}
