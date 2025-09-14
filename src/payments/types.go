package payments

type Client interface {
	CreatePaymentOrder(req *CreatePaymentOrderRequest) (*CreatePaymentOrderResponse, error)
}

type CreatePaymentOrderRequest struct {
	Amount   float32           `json:"amount"`
	Currency string            `json:"currency"`
	Receipt  string            `json:"receipt"`
	Metadata map[string]string `json:"metadata,omitempty"` // optional metadata for additional info
}

type CreatePaymentOrderResponse struct {
	ID          string            `json:"id"`
	Amount      float64           `json:"amount"`
	Currency    string            `json:"currency"`
	Receipt     string            `json:"receipt"`
	Status      string            `json:"status"`
	CreatedAt   int64             `json:"created_at"`
	PaymentID   string            `json:"payment_id"`
	WebhookSent bool              `json:"webhook_sent"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

type PaymentsMetadata struct {
	BookingID int64 `json:"booking_id"`
}

const BASE_URL = "http://localhost:4000"
