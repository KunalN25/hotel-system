package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"

	"hotel-system/src/types/hotelsystem"
)

// ---------- Main ----------

func main() {
	AddMultipleHotels()
}

// func BookHotel() {

// }

func AddMultipleHotels() error {
	rand.Seed(time.Now().UnixNano())

	cities := []struct {
		Name  string
		State string
	}{
		{"Mumbai", "Maharashtra"},
		{"Bangalore", "Karnataka"},
	}

	// name parts to generate somewhat realistic hotel names
	prefixes := []string{"Azure", "Royal", "Lotus", "Ocean", "Urban", "Metro", "Grand", "Harbour", "Saffron", "Crown", "Serene", "Palace", "Sunrise", "Garden", "Cloud"}
	suffixes := []string{"Inn", "Suites", "Residency", "Hotel", "Lodge", "Villa", "Stay", "Retreat", "Haven", "Plaza"}

	streets := []string{"MG Road", "Brigade Road", "Linking Road", "SV Road", "Churchgate", "Colaba Causeway", "Juhu Tara Road", "Koramangala 5th Block", "Indiranagar 100 ft Road"}
	landmarks := []string{"Near Central Mall", "Opposite Metro Station", "Beside City Park", "Next to Stadium", "Close to Beach", "Near Railway Station"}
	localities := []string{"Andheri", "Bandra", "Juhu", "Churchgate", "Powai", "Koramangala", "Indiranagar", "Whitefield", "MG Road", "Brigade Road"}

	// For each hotel, pick images using Unsplash source that returns a relevant image
	imageFor := func(city string) []string {
		// Using source.unsplash to get random relevant images for the "hotel" + city query.
		// This gives you variable images without scraping Google.
		// You can build a small set per hotel to look nicer in UIs.
		return []string{
			fmt.Sprintf("https://source.unsplash.com/featured/?hotel,%s", city),
			fmt.Sprintf("https://source.unsplash.com/featured/?interior,%s", city),
			fmt.Sprintf("https://source.unsplash.com/featured/?lobby,%s", city),
		}
	}

	for i := 0; i < 10; i++ {
		c := cities[rand.Intn(len(cities))]

		name := fmt.Sprintf("%s %s", prefixes[rand.Intn(len(prefixes))], suffixes[rand.Intn(len(suffixes))])
		desc := fmt.Sprintf("%s located in the heart of %s. Comfortable rooms, free WiFi and great service.", name, c.Name)

		req := hotelsystem.AddHotelRequest{
			Name:         name,
			Description:  desc,
			Images:       imageFor(c.Name),
			CostPerNight: int64(1500 + rand.Intn(8500)), // ₹1500 - ₹10000
			TotalRooms:   int64(10 + rand.Intn(90)),     // 10 - 99 rooms
			Street:       streets[rand.Intn(len(streets))],
			Landmark:     landmarks[rand.Intn(len(landmarks))],
			Locality:     localities[rand.Intn(len(localities))],
			City:         c.Name,
			State:        c.State,
			Pincode:      fmt.Sprintf("%06d", 400000+rand.Intn(99999)), // string field in request is fine (struct has string pincode)
		}

		// Note: your AddHotelRequest declares Pincode as string in earlier snippet — above we produce string.
		// If your AddHotelRequest has Pincode as int, convert accordingly:
		// pincodeInt := 400000 + rand.Intn(99999)
		// req.Pincode = strconv.Itoa(pincodeInt) // or assign int if field is int

		// send request
		url := "/addHotel"
		var resp map[string]any
		headers := map[string]string{
			"Content-Type": "application/json",
			// add authorization header if needed, e.g. "Authorization": "Bearer <token>"
		}

		// call the generic helper - replace with your real postJSON implementation
		err := postJSON(url, req, &resp, headers)
		if err != nil {
			log.Printf("failed to add hotel %q (%s): %v", req.Name, req.City, err)
			// continue trying the remaining hotels; return error only if you prefer to abort
			continue
		}

		log.Printf("added hotel %q in %s -> response: %v", req.Name, req.City, resp)
	}

	return nil
}

func postJSON[Req any, Resp any](endpoint string, payload Req, out *Resp, headers map[string]string) error {
	baseUrl := "http://localhost:8080"
	url := baseUrl + endpoint
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
