package btcpaybasics

import "sync"

type DeliveryStore struct {
	mu   sync.Mutex
	seen map[string]struct{}
}

func NewDeliveryStore() *DeliveryStore {
	return &DeliveryStore{
		seen: make(map[string]struct{}),
	}
}

func (s *DeliveryStore) MarkIfNew(deliveryID string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.seen[deliveryID]; exists {
		return false
	}

	s.seen[deliveryID] = struct{}{}
	return true
}
