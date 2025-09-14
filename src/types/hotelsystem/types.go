package hotelsystem

// GetHotelsListRequest corresponds to proto GetHotelsListRequest.
type GetHotelsListRequest struct {
	SearchQuery string `json:"search_query"`
	Limit       int32  `json:"limit"`
	Offset      int32  `json:"offset"`
}

// HotelData corresponds to proto HotelData.
type HotelData struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Images       []string `json:"images"`
	CostPerNight float32  `json:"cost_per_night"`
}

// GetImages returns a non-nil slice of images.
func (h *HotelData) GetImages() []string {
	if h == nil || h.Images == nil {
		return []string{}
	}
	return h.Images
}

// GetHotelsListResponse corresponds to proto GetHotelsListResponse.
type GetHotelsListResponse struct {
	TotalRecords int64        `json:"total_records"`
	HotelsList   []*HotelData `json:"hotels_list"`
}

// GetHotelsList returns a non-nil slice of hotels.
func (r *GetHotelsListResponse) GetHotelsList() []*HotelData {
	if r == nil || r.HotelsList == nil {
		return []*HotelData{}
	}
	return r.HotelsList
}

// GetHotelByIdRequest corresponds to proto GetHotelByIdRequest.
type GetHotelByIdRequest struct {
	HotelID int64 `json:"hotel_id"`
}

// GetHotelByIdResponse corresponds to proto GetHotelByIdResponse.
type GetHotelByIdResponse struct {
	ID             int64    `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Images         []string `json:"images"`
	CostPerNight   int64    `json:"cost_per_night"`
	AvailableRooms int64    `json:"available_rooms"`
	Address        string   `json:"address"`
}

// GetImages returns a non-nil slice of images.
func (r *GetHotelByIdResponse) GetImages() []string {
	if r == nil || r.Images == nil {
		return []string{}
	}
	return r.Images
}

// AddHotelRequest corresponds to proto AddHotelRequest.
type AddHotelRequest struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Images       []string `json:"images"`
	CostPerNight int64    `json:"cost_per_night"`
	TotalRooms   int64    `json:"total_rooms"`
	Street       string   `json:"street"`
	Landmark     string   `json:"landmark"`
	Locality     string   `json:"locality"`
	City         string   `json:"city"`
	State        string   `json:"state"`
	Pincode      string   `json:"pincode"`
}

// GetImages returns a non-nil slice of images.
func (r *AddHotelRequest) GetImages() []string {
	if r == nil || r.Images == nil {
		return []string{}
	}
	return r.Images
}

// BookHotelRequest corresponds to proto BookHotelRequest.
type BookHotelRequest struct {
	HotelID     int64  `json:"hotel_id"`
	NumRooms    int32  `json:"num_rooms"`
	NumDays     int32  `json:"num_days"`
	CheckInDate string `json:"check_in_date"` // YYYY-MM-DD
}

// BookHotelResponse corresponds to proto BookHotelResponse.
type BookHotelResponse struct {
	PaymentOrderID string                           `json:"payment_order_id"`
	Success        bool                             `json:"success"`
	Message        string                           `json:"message"`
	BookingDetails *BookHotelResponseBookingDetails `json:"booking_details"`
	CheckoutURL    string                           `json:"checkout_url"`
}

// GetBookingDetails returns a non-nil booking details struct.
func (r *BookHotelResponse) GetBookingDetails() *BookHotelResponseBookingDetails {
	if r == nil || r.BookingDetails == nil {
		return &BookHotelResponseBookingDetails{}
	}
	return r.BookingDetails
}

// BookHotelResponseBookingDetails mirrors nested BookHotelResponse.BookingDetails.
type BookHotelResponseBookingDetails struct {
	BookingID   int64   `json:"booking_id"`
	HotelID     int64   `json:"hotel_id"`
	NumRooms    int32   `json:"num_rooms"`
	NumDays     int32   `json:"num_days"`
	CheckInDate string  `json:"check_in_date"`
	TotalCost   float32 `json:"total_cost"`
}

// GenericSuccessResponse corresponds to proto GenericSuccessResponse.
type GenericSuccessResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// PaymentWebhookRequest corresponds to proto PaymentWebhookRequest.
type PaymentWebhookRequest struct {
	Entity  string                 `json:"entity"`
	Event   string                 `json:"event"`
	Payload *PaymentWebhookPayload `json:"payload"`
}

// GetPayload returns a non-nil payload struct.
func (r *PaymentWebhookRequest) GetPayload() *PaymentWebhookPayload {
	if r == nil || r.Payload == nil {
		return &PaymentWebhookPayload{}
	}
	return r.Payload
}

// PaymentWebhookPayload mirrors nested PaymentWebhookRequest.Payload.
type PaymentWebhookPayload struct {
	Payment *PaymentDetails `json:"payment"`
}

// GetPayment returns a non-nil payment details struct.
func (p *PaymentWebhookPayload) GetPayment() *PaymentDetails {
	if p == nil || p.Payment == nil {
		return &PaymentDetails{}
	}
	return p.Payment
}

// PaymentDetails mirrors nested PaymentWebhookRequest.PaymentDetails.
type PaymentDetails struct {
	ID        string            `json:"id"`
	OrderID   string            `json:"order_id"`
	Amount    int32             `json:"amount"`
	Currency  string            `json:"currency"`
	CreatedAt int64             `json:"createdAt"`
	Status    string            `json:"status"`
	Metadata  map[string]string `json:"metadata"`
	Entity    string            `json:"entity"`
	Method    string            `json:"method"`
}

// GetMetadata returns a non-nil metadata map.
func (p *PaymentDetails) GetMetadata() map[string]string {
	if p == nil || p.Metadata == nil {
		return map[string]string{}
	}
	return p.Metadata
}

// PaymentStatusRequest corresponds to proto PaymentStatusRequest.
type PaymentStatusRequest struct {
	OrderID string `json:"order_id"`
}

// PaymentStatusResponse corresponds to proto PaymentStatusResponse.
type PaymentStatusResponse struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// GetBookingDetailsRequest corresponds to proto GetBookingDetailsRequest.
type GetBookingDetailsRequest struct {
	BookingID int64 `json:"booking_id"`
}

// GetBookingDetailsResponse corresponds to proto GetBookingDetailsResponse.
type GetBookingDetailsResponse struct {
	BookingID     int64   `json:"booking_id"`
	NumRooms      int32   `json:"num_rooms"`
	NumDays       int32   `json:"num_days"`
	CheckInDate   string  `json:"check_in_date"`
	CheckOutDate  string  `json:"check_out_date"`
	TotalCost     float32 `json:"total_cost"`
	Status        string  `json:"status"`
	PaymentStatus string  `json:"payment_status"`
	BookingTime   string  `json:"booking_time"`
}
