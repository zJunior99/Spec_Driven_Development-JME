package payment_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/zJunior99/Spec_Driven_Development-JME/internal/payment"
	"github.com/zJunior99/Spec_Driven_Development-JME/internal/repository"
)

var _ = Describe("Payment Service", func() {

	var (
		svc  *payment.Service
		repo *repository.InMemoryRepository
	)

	BeforeEach(func() {
		repo = repository.NewInMemoryRepository()
		svc = payment.NewService(repo)
	})

	Describe("Creating a payment", func() {
		Context("with valid data", func() {
			It("should create a payment with pending status", func() {
				req := payment.CreateRequest{
					OrderID:  "ORD-001",
					Amount:   150.00,
					Currency: "PEN",
				}

				p, err := svc.Create(req)

				Expect(err).NotTo(HaveOccurred())
				Expect(p.ID).NotTo(BeEmpty())
				Expect(p.Status).To(Equal(payment.StatusPending))
				Expect(p.Amount).To(Equal(150.00))
				Expect(p.Currency).To(Equal("PEN"))
			})

			It("should create a payment with a valid USD amount", func() {
				req := payment.CreateRequest{
					OrderID:  "ORD-002",
					Amount:   99.99,
					Currency: "USD",
				}

				p, err := svc.Create(req)

				Expect(err).NotTo(HaveOccurred())
				Expect(p.Currency).To(Equal("USD"))
				Expect(p.Amount).To(Equal(99.99))
			})

			It("should assign a unique ID to each payment", func() {
				req := payment.CreateRequest{OrderID: "ORD-003", Amount: 10.00, Currency: "EUR"}

				p1, _ := svc.Create(req)
				p2, _ := svc.Create(req)

				Expect(p1.ID).NotTo(Equal(p2.ID))
			})
		})

		Context("with invalid data", func() {
			It("should reject a payment with zero amount", func() {
				req := payment.CreateRequest{
					OrderID:  "ORD-004",
					Amount:   0,
					Currency: "USD",
				}

				_, err := svc.Create(req)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("amount"))
			})

			It("should reject a payment with negative amount", func() {
				req := payment.CreateRequest{
					OrderID:  "ORD-005",
					Amount:   -50.00,
					Currency: "USD",
				}

				_, err := svc.Create(req)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("amount"))
			})

			It("should reject a payment with unsupported currency", func() {
				req := payment.CreateRequest{
					OrderID:  "ORD-006",
					Amount:   100.00,
					Currency: "JPY",
				}

				_, err := svc.Create(req)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("currency"))
			})

			It("should reject a payment without order ID", func() {
				req := payment.CreateRequest{
					Amount:   100.00,
					Currency: "USD",
				}

				_, err := svc.Create(req)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("order_id"))
			})
		})
	})

	Describe("Querying a payment", func() {
		Context("when the payment exists", func() {
			It("should return the payment with correct data", func() {
				req := payment.CreateRequest{OrderID: "ORD-007", Amount: 200.00, Currency: "USD"}
				created, _ := svc.Create(req)

				found, err := svc.GetByID(created.ID)

				Expect(err).NotTo(HaveOccurred())
				Expect(found.ID).To(Equal(created.ID))
				Expect(found.OrderID).To(Equal("ORD-007"))
				Expect(found.Amount).To(Equal(200.00))
			})
		})

		Context("when the payment does not exist", func() {
			It("should return an error", func() {
				_, err := svc.GetByID("non-existent-id")

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Cancelling a payment", func() {
		Context("when the payment is pending", func() {
			It("should cancel the payment successfully", func() {
				req := payment.CreateRequest{OrderID: "ORD-008", Amount: 300.00, Currency: "EUR"}
				created, _ := svc.Create(req)

				err := svc.Cancel(created.ID)

				Expect(err).NotTo(HaveOccurred())

				cancelled, _ := svc.GetByID(created.ID)
				Expect(cancelled.Status).To(Equal(payment.StatusCancelled))
			})
		})

		Context("when the payment is already cancelled", func() {
			It("should return a conflict error", func() {
				req := payment.CreateRequest{OrderID: "ORD-009", Amount: 100.00, Currency: "PEN"}
				created, _ := svc.Create(req)
				_ = svc.Cancel(created.ID)

				err := svc.Cancel(created.ID)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot cancel"))
			})
		})

		Context("when the payment does not exist", func() {
			It("should return an error", func() {
				err := svc.Cancel("non-existent-id")

				Expect(err).To(HaveOccurred())
			})
		})
	})
})
