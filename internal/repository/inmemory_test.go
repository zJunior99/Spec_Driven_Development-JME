package repository_test

import (
	"testing"

	"github.com/zJunior99/Spec_Driven_Development-JME/internal/payment"
	"github.com/zJunior99/Spec_Driven_Development-JME/internal/repository"
)

func newPayment(id string) payment.Payment {
	return payment.Payment{
		ID:       id,
		OrderID:  "ORD-001",
		Amount:   100.00,
		Currency: "PEN",
		Status:   payment.StatusPending,
	}
}

func TestSave_AndFindByID(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	p := newPayment("PAY-001")

	if err := repo.Save(p); err != nil {
		t.Fatalf("Save() error inesperado: %v", err)
	}

	got, err := repo.FindByID("PAY-001")
	if err != nil {
		t.Fatalf("FindByID() error inesperado: %v", err)
	}
	if got.ID != p.ID {
		t.Errorf("FindByID() = %q; quería %q", got.ID, p.ID)
	}
}

func TestFindByID_NotFound(t *testing.T) {
	repo := repository.NewInMemoryRepository()

	_, err := repo.FindByID("NO-EXISTE")
	if err == nil {
		t.Error("FindByID() con ID inexistente debería retornar error")
	}
}

func TestUpdate_ExistingPayment(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	p := newPayment("PAY-002")
	_ = repo.Save(p)

	p.Status = payment.StatusCancelled
	if err := repo.Update(p); err != nil {
		t.Fatalf("Update() error inesperado: %v", err)
	}

	got, _ := repo.FindByID("PAY-002")
	if got.Status != payment.StatusCancelled {
		t.Errorf("Update() status = %q; quería %q", got.Status, payment.StatusCancelled)
	}
}

func TestUpdate_NotFound(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	p := newPayment("PAY-999")

	err := repo.Update(p)
	if err == nil {
		t.Error("Update() con ID inexistente debería retornar error")
	}
}

func TestSave_Overwrites(t *testing.T) {
	repo := repository.NewInMemoryRepository()
	p := newPayment("PAY-003")
	_ = repo.Save(p)

	p.Amount = 999.00
	_ = repo.Save(p)

	got, _ := repo.FindByID("PAY-003")
	if got.Amount != 999.00 {
		t.Errorf("Save() debería sobreescribir; Amount = %.2f; quería 999.00", got.Amount)
	}
}
