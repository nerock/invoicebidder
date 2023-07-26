package broker

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/nerock/invoicebidder/internal/investor"

	"github.com/bojanz/currency"
	"github.com/nerock/invoicebidder/internal/invoice"
)

type InvoiceService interface {
	GetInvoice(context.Context, string) (invoice.Invoice, error)
	ApproveTrade(context.Context, string, bool) error
}
type InvestorService interface {
	CancelTrade(context.Context, []investor.Bid) error
	CancelBid(context.Context, string, currency.Amount) error
}

type IssuerService interface {
	ApproveTrade(context.Context, string, currency.Amount) error
}

type Broker struct {
	maxRetries    int
	eventHandlers int
	events        chan Event
	wg            *sync.WaitGroup

	invoiceService  InvoiceService
	investorService InvestorService
	issuerService   IssuerService
}

func (b *Broker) SendTradeEvent(invoiceID string, approved bool) {
	b.events <- &TradeEvent{
		InvoiceID: invoiceID,
		Approved:  approved,
	}
}

func (b *Broker) SendFailedBidEvent(investorID string, amount currency.Amount) {
	b.events <- &FailedBidEvent{
		InvestorID: investorID,
		Amount:     amount,
	}
}

func New(eventHandlers, eventBuffer, maxRetries int, invoiceService InvoiceService, investorService InvestorService, issuerService IssuerService) *Broker {
	return &Broker{
		invoiceService:  invoiceService,
		investorService: investorService,
		issuerService:   issuerService,
		events:          make(chan Event, eventBuffer), // Random buffer number
		eventHandlers:   eventHandlers,
		maxRetries:      maxRetries,
		wg:              &sync.WaitGroup{},
	}
}

func (b *Broker) Serve() error {
	for i := 0; i < b.eventHandlers; i++ {
		go b.eventHandler(b.wg, b.events)
	}

	return nil
}

func (b *Broker) Shutdown(ctx context.Context) error {
	close(b.events)

	closeChan := make(chan struct{})
	go func() {
		b.wg.Wait()
		closeChan <- struct{}{}
	}()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for events to comple")
		case <-closeChan:
			return nil
		}
	}
}

func (b *Broker) eventHandler(wg *sync.WaitGroup, events chan Event) {
	wg.Add(1)
	defer wg.Done()

	for e := range events {
		var err error

		switch e.Type() {
		case TypeTradeEvent:
			err = b.tradeEventHandler(e.(*TradeEvent))
		case TypeFailedBidEvent:
			err = b.failedBidEventHandler(e.(*FailedBidEvent))
		}

		if err != nil {
			if e.Retries() > b.maxRetries {
				log.Printf("max retries exahusted: %s", err)
			} else {
				log.Println(err)
				events <- e
			}
		}
	}
}

func (b *Broker) failedBidEventHandler(be *FailedBidEvent) error {
	return b.investorService.CancelBid(context.Background(), be.InvestorID, be.Amount)
}

func (b *Broker) tradeEventHandler(te *TradeEvent) error {
	inv, err := b.invoiceService.GetInvoice(context.Background(), te.InvoiceID)
	if err != nil {
		return err
	}

	if te.Approved {
		err = b.approveTradeEvent(inv)
	} else {
		err = b.cancelTradeEvent(inv)
	}

	return err
}

func (b *Broker) approveTradeEvent(inv invoice.Invoice) error {
	return b.issuerService.ApproveTrade(context.Background(), inv.IssuerID, inv.Price)
}

func (b *Broker) cancelTradeEvent(inv invoice.Invoice) error {
	invBids := make([]investor.Bid, len(inv.Bids))
	for _, bid := range inv.Bids {
		invBids = append(invBids, investor.Bid{
			InvestorID: bid.InvestorID,
			Amount:     bid.Amount,
		})
	}

	return b.investorService.CancelTrade(context.Background(), invBids)
}
