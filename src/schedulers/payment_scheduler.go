package scheduler

import (
	"hotel-system/src/store"
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

//func (ps *PaymentScheduler) Start() {
//	c := cron.New()
//	// Runs every 10 minutes
//	c.AddFunc("*/10 * * * *", func() {
//		err := ps.MarkPendingPaymentsAsFailure()
//		if err != nil {
//			log.Printf("Failed to release payments: %v", err)
//		}
//	})
//
//	c.Start()
//}
//
//func (ps *PaymentScheduler) MarkPendingPaymentsAsFailure() error {
//
//}
