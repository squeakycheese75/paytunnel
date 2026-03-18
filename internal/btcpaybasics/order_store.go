package btcpaybasics

import (
	"sync"
)

type OrderStore struct {
	mu     sync.Mutex
	orders map[string]*Order
}

func NewOrderStore() *OrderStore {
	return &OrderStore{
		orders: map[string]*Order{
			"order-123": {
				ID:        "order-123",
				InvoiceID: "inv_123",
				Status:    "pending",
			},
		},
	}
}

func (s *OrderStore) List() map[string]*Order {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make(map[string]*Order, len(s.orders))
	for k, v := range s.orders {
		orderCopy := *v
		result[k] = &orderCopy
	}

	return result
}

func (s *OrderStore) MarkPaidByInvoiceID(invoiceID string) (*Order, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, order := range s.orders {
		if order.InvoiceID != invoiceID {
			continue
		}

		order.Status = "paid"
		orderCopy := *order
		return &orderCopy, true
	}

	return nil, false
}
