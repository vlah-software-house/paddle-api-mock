package store

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/vlah-software-house/paddle-api-mock/internal/models"
)

var idCounter uint64

func NextID(prefix string) string {
	n := atomic.AddUint64(&idCounter, 1)
	return fmt.Sprintf("%s_%08d", prefix, n)
}

// Store is the thread-safe in-memory data store for all Paddle resources.
type Store struct {
	mu sync.RWMutex

	Products             map[string]*models.Product
	Prices               map[string]*models.Price
	Customers            map[string]*models.Customer
	Subscriptions        map[string]*models.Subscription
	Transactions         map[string]*models.Transaction
	Events               []*models.Event
	NotificationSettings map[string]*models.NotificationSetting
}

func New() *Store {
	return &Store{
		Products:             make(map[string]*models.Product),
		Prices:               make(map[string]*models.Price),
		Customers:            make(map[string]*models.Customer),
		Subscriptions:        make(map[string]*models.Subscription),
		Transactions:         make(map[string]*models.Transaction),
		Events:               make([]*models.Event, 0),
		NotificationSettings: make(map[string]*models.NotificationSetting),
	}
}

// Reset clears all data from the store.
func (s *Store) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Products = make(map[string]*models.Product)
	s.Prices = make(map[string]*models.Price)
	s.Customers = make(map[string]*models.Customer)
	s.Subscriptions = make(map[string]*models.Subscription)
	s.Transactions = make(map[string]*models.Transaction)
	s.Events = make([]*models.Event, 0)
	s.NotificationSettings = make(map[string]*models.NotificationSetting)
}

// --- Products ---

func (s *Store) GetProduct(id string) (*models.Product, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.Products[id]
	return p, ok
}

func (s *Store) ListProducts() []*models.Product {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*models.Product, 0, len(s.Products))
	for _, p := range s.Products {
		result = append(result, p)
	}
	return result
}

func (s *Store) SetProduct(p *models.Product) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Products[p.ID] = p
}

// --- Prices ---

func (s *Store) GetPrice(id string) (*models.Price, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.Prices[id]
	return p, ok
}

func (s *Store) ListPrices() []*models.Price {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*models.Price, 0, len(s.Prices))
	for _, p := range s.Prices {
		result = append(result, p)
	}
	return result
}

func (s *Store) SetPrice(p *models.Price) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Prices[p.ID] = p
}

// --- Customers ---

func (s *Store) GetCustomer(id string) (*models.Customer, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.Customers[id]
	return c, ok
}

func (s *Store) ListCustomers() []*models.Customer {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*models.Customer, 0, len(s.Customers))
	for _, c := range s.Customers {
		result = append(result, c)
	}
	return result
}

func (s *Store) SetCustomer(c *models.Customer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Customers[c.ID] = c
}

// --- Subscriptions ---

func (s *Store) GetSubscription(id string) (*models.Subscription, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sub, ok := s.Subscriptions[id]
	return sub, ok
}

func (s *Store) ListSubscriptions() []*models.Subscription {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*models.Subscription, 0, len(s.Subscriptions))
	for _, sub := range s.Subscriptions {
		result = append(result, sub)
	}
	return result
}

func (s *Store) SetSubscription(sub *models.Subscription) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Subscriptions[sub.ID] = sub
}

// --- Transactions ---

func (s *Store) GetTransaction(id string) (*models.Transaction, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.Transactions[id]
	return t, ok
}

func (s *Store) ListTransactions() []*models.Transaction {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*models.Transaction, 0, len(s.Transactions))
	for _, t := range s.Transactions {
		result = append(result, t)
	}
	return result
}

func (s *Store) SetTransaction(t *models.Transaction) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Transactions[t.ID] = t
}

// --- Events ---

func (s *Store) AddEvent(e *models.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Events = append(s.Events, e)
}

func (s *Store) ListEvents() []*models.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*models.Event, len(s.Events))
	copy(result, s.Events)
	return result
}

// --- Notification Settings ---

func (s *Store) GetNotificationSetting(id string) (*models.NotificationSetting, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ns, ok := s.NotificationSettings[id]
	return ns, ok
}

func (s *Store) ListNotificationSettings() []*models.NotificationSetting {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]*models.NotificationSetting, 0, len(s.NotificationSettings))
	for _, ns := range s.NotificationSettings {
		result = append(result, ns)
	}
	return result
}

func (s *Store) SetNotificationSetting(ns *models.NotificationSetting) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.NotificationSettings[ns.ID] = ns
}
