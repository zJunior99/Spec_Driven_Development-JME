package main

import (
	"fmt"

	"github.com/zJunior99/Spec_Driven_Development-JME/internal/payment"
	"github.com/zJunior99/Spec_Driven_Development-JME/internal/repository"
)

func main() {
	fmt.Println("=== Spec Driven Development — Payment System ===")
	fmt.Println("Spec first, code second.")

	repo := repository.NewInMemoryRepository()
	svc := payment.NewService(repo)

	requests := []payment.CreateRequest{
		{OrderID: "ORD-001", Amount: 150.00, Currency: "PEN"},
		{OrderID: "ORD-002", Amount: 99.99, Currency: "USD"},
		{OrderID: "ORD-003", Amount: 0, Currency: "USD"},
		{OrderID: "ORD-004", Amount: 200.00, Currency: "JPY"},
	}

	var created []payment.Payment

	for _, req := range requests {
		p, err := svc.Create(req)
		if err != nil {
			fmt.Printf("❌ [%s] Error: %s\n", req.OrderID, err)
			continue
		}
		fmt.Printf("✅ [%s] Payment created — ID: %s | Amount: %.2f %s | Status: %s\n",
			p.OrderID, p.ID[:8]+"...", p.Amount, p.Currency, p.Status)
		created = append(created, p)
	}

	if len(created) > 0 {
		fmt.Printf("\n--- Cancelling payment %s ---\n", created[0].ID[:8]+"...")
		if err := svc.Cancel(created[0].ID); err != nil {
			fmt.Printf("❌ Cancel error: %s\n", err)
		} else {
			p, _ := svc.GetByID(created[0].ID)
			fmt.Printf("✅ Payment cancelled — Status: %s\n", p.Status)
		}

		fmt.Printf("\n--- Trying to cancel again ---\n")
		if err := svc.Cancel(created[0].ID); err != nil {
			fmt.Printf("❌ Expected conflict: %s\n", err)
		}
	}

	fmt.Println("\n=== Run 'go test ./... -v' to see all specs ===")
}
