package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hotel-system/src/pb"
	"io"
	"net/http"

	"github.com/google/uuid"
)

// ---------- Main ----------

func main() {
	BookHotel()
}

func BookHotel() {
	url := "http://localhost:8080/bookHotel"
	req := &pb.BookHotelRequest{
		HotelId:     3,
		NumRooms:    2,
		NumDays:     2,
		CheckInDate: "2025-08-24",
	}

	ik := uuid.New().String()
	headers := map[string]string{
		"Idempotency-Key": ik,
	}

	var resp pb.BookHotelResponse
	err := postJSON(url, req, &resp, headers)
	if err != nil {
		fmt.Println("BookHotel API error:", err)
		return
	}

	fmt.Println("Resp:", resp)

}

func postJSON[Req any, Resp any](url string, payload Req, out *Resp, headers map[string]string) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("create request error: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("http post error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("non-200: %d — %s", resp.StatusCode, string(b))
	}

	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decode error: %w", err)
	}
	return nil
}
