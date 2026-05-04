package payment

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Status representa el estado de un pago.
type Status string

const (
	StatusPending   Status = "pending"
	StatusCompleted Status = "completed"
	StatusCancelled Status = "cancelled"
	StatusFailed    Status = "failed"
)

var supportedCurrencies = map[string]bool{
	"USD": true,
	"PEN": true,
	"EUR": true,
}

// Payment representa un pago en el sistema.
type Payment struct {
	ID        string
	OrderID   string
	Amount    float64
	Currency  string
	Status    Status
	CreatedAt time.Time
}

// CreateRequest contiene los datos necesarios para crear un pago.
type CreateRequest struct {
	OrderID  string
	Amount   float64
	Currency string
}

// Repository define el contrato de persistencia de pagos.
type Repository interface {
	Save(p Payment) error
	FindByID(id string) (Payment, error)
	Update(p Payment) error
}

// Service orquesta la lógica de negocio de pagos.
type Service struct {
	repo Repository
}

// NewService crea un nuevo servicio de pagos con el repositorio dado.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Create valida y crea un nuevo pago con estado pendiente.
func (s *Service) Create(req CreateRequest) (Payment, error) {
	if err := validateCreateRequest(req); err != nil {
		return Payment{}, err
	}

	p := Payment{
		ID:        uuid.NewString(),
		OrderID:   req.OrderID,
		Amount:    req.Amount,
		Currency:  req.Currency,
		Status:    StatusPending,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Save(p); err != nil {
		return Payment{}, fmt.Errorf("error saving payment: %w", err)
	}

	return p, nil
}

// GetByID recupera un pago por su ID.
func (s *Service) GetByID(id string) (Payment, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return Payment{}, fmt.Errorf("payment %s not found: %w", id, err)
	}
	return p, nil
}

// Cancel cancela un pago si está en estado pendiente.
func (s *Service) Cancel(id string) error {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("payment %s not found: %w", id, err)
	}

	if p.Status != StatusPending {
		return fmt.Errorf("cannot cancel payment with status %q", p.Status)
	}

	p.Status = StatusCancelled
	return s.repo.Update(p)
}

// validateCreateRequest verifica que los datos del request sean válidos.
func validateCreateRequest(req CreateRequest) error {
	if req.OrderID == "" {
		return fmt.Errorf("order_id is required")
	}
	if req.Amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	if !supportedCurrencies[req.Currency] {
		return fmt.Errorf("currency %q is not supported", req.Currency)
	}
	return nil
}
