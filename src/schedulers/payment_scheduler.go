package scheduler

import (
	"hotel-system/src/constants"
	"hotel-system/src/store"
	"log"

	"github.com/robfig/cron/v3"
)

type PaymentScheduler struct {
	storage store.StorageService
}

func NewPaymentScheduler(storage store.StorageService) *PaymentScheduler {
	return &PaymentScheduler{
		storage: storage,
	}
}

// TODO: Implement payment scheduler to mark pending payments as failed after certain time period

func (ps *PaymentScheduler) Start() {
	c := cron.New()
	// Runs every 10 minutes
	c.AddFunc("*/10 * * * *", func() {
		err := ps.MarkPendingPaymentsAsFailure()
		if err != nil {
			log.Printf("Failed to release payments: %v", err)
		}
	})

	c.Start()
}

func (ps *PaymentScheduler) MarkPendingPaymentsAsFailure() error {
	expiredPayments, err := ps.storage.GetExpiredPayments()
	if err != nil {
		return err
	}
	for _, payment := range expiredPayments {
		updateErr := ps.storage.UpdatePaymentStatus(payment.ID, constants.PAYMENT_FAILED)
		if updateErr != nil {

			log.Printf("Failed to update payment %v to failed: %v", payment.ID, updateErr)
		}
	}
	return nil
}
