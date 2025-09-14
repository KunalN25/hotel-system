package payments

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type paymentsClient struct {
	httpClient *http.Client
}

func (p *paymentsClient) CreatePaymentOrder(req *CreatePaymentOrderRequest) (*CreatePaymentOrderResponse, error) {

	url := fmt.Sprintf("%s/order/create", BASE_URL)

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("payment gateway responded with status %d", resp.StatusCode)
	}

	var out CreatePaymentOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &out, nil
}

func NewPaymentsClient() Client {
	return &paymentsClient{
		httpClient: &http.Client{},
	}
}
