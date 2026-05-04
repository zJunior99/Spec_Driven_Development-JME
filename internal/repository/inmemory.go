package repository

import (
	"fmt"
	"sync"

	"github.com/zJunior99/Spec_Driven_Development-JME/internal/payment"
)

// InMemoryRepository implementa payment.Repository usando un mapa en memoria.
type InMemoryRepository struct {
	mu       sync.RWMutex
	payments map[string]payment.Payment
}

// NewInMemoryRepository crea un repositorio en memoria vacío.
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		payments: make(map[string]payment.Payment),
	}
}

// Save persiste un pago en memoria.
func (r *InMemoryRepository) Save(p payment.Payment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.payments[p.ID] = p
	return nil
}

// FindByID recupera un pago por su ID.
func (r *InMemoryRepository) FindByID(id string) (payment.Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.payments[id]
	if !ok {
		return payment.Payment{}, fmt.Errorf("payment not found")
	}
	return p, nil
}

// Update actualiza un pago existente.
func (r *InMemoryRepository) Update(p payment.Payment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.payments[p.ID]; !ok {
		return fmt.Errorf("payment not found")
	}
	r.payments[p.ID] = p
	return nil
}
