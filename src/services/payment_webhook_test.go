package services

//
//import (
//	"bytes"
//	"encoding/json"
//	"hotel-system/src/errorcodes"
//	"hotel-system/src/payments"
//	"hotel-system/src/pb"
//	"hotel-system/src/store"
//	"hotel-system/src/store/mocks"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//	"time"
//
//	"github.com/golang/mock/gomock"
//	"github.com/stretchr/testify/assert"
//)
//
//func TestPaymentWebhookHandler(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockStore := mocks.NewMockStorageService(ctrl)
//	service := &Service{
//		storageService: mockStore,
//	}
//
//	tests := []struct {
//		name           string
//		requestBody    interface{}
//		setupMocks     func()
//		expectedStatus int
//		expectedBody   string
//	}{
//		{
//			name: "successful payment webhook - payment success",
//			requestBody: pb.PaymentWebhookRequest{
//				Payload: &pb.PaymentWebhookRequest_Payload{
//					Status: string(store.PAYMENT_SUCCESS),
//					Metadata: func() []byte {
//						metadata := payments.PaymentsMetadata{BookingID: 1}
//						data, _ := json.Marshal(metadata)
//						return data
//					}(),
//				},
//			},
//			setupMocks: func() {
//				payment := &store.Payment{
//					ID:        1,
//					BookingID: 1,
//					OrderID:   "order_123",
//					Amount:    1000.0,
//					Status:    store.PAYMENT_PENDING,
//				}
//
//				booking := store.Booking{
//					BookingID:     1,
//					HotelID:       1,
//					UserID:        1,
//					NumberOfRooms: 2,
//					NumberOfDays:  3,
//					BookingTime:   time.Now(),
//					CheckInDate:   time.Now().AddDate(0, 0, 1),
//					CheckOutDate:  time.Now().AddDate(0, 0, 4),
//					Status:        store.PENDING,
//				}
//
//				mockStore.EXPECT().
//					GetPaymentByOrderId("order_123").
//					Return(payment, nil)
//
//				mockStore.EXPECT().
//					GetBookingById(int64(1)).
//					Return(booking, nil)
//
//				mockStore.EXPECT().
//					UpdateBookingStatus(int64(1), store.CONFIRMED).
//					Return(nil)
//
//				mockStore.EXPECT().
//					UpdatePaymentStatus(1, store.PAYMENT_SUCCESS).
//					Return(nil)
//			},
//			expectedStatus: http.StatusOK,
//		},
//		{
//			name: "successful payment webhook - payment failed",
//			requestBody: pb.PaymentWebhookRequest{
//				Payload: &pb.PaymentWebhookRequest_Payload{
//					OrderId: "order_124",
//					Status:  string(store.PAYMENT_FAILED),
//					Metadata: func() []byte {
//						metadata := payments.PaymentsMetadata{BookingID: 2}
//						data, _ := json.Marshal(metadata)
//						return data
//					}(),
//				},
//			},
//			setupMocks: func() {
//				payment := &store.Payment{
//					ID:        2,
//					BookingID: 2,
//					OrderID:   "order_124",
//					Amount:    1500.0,
//					Status:    store.PAYMENT_PENDING,
//				}
//
//				booking := store.Booking{
//					BookingID:     2,
//					HotelID:       2,
//					UserID:        2,
//					NumberOfRooms: 1,
//					NumberOfDays:  2,
//					BookingTime:   time.Now(),
//					CheckInDate:   time.Now().AddDate(0, 0, 1),
//					CheckOutDate:  time.Now().AddDate(0, 0, 3),
//					Status:        store.PENDING,
//				}
//
//				hotel := store.Hotel{
//					ID:             2,
//					Name:           "Test Hotel",
//					AvailableRooms: 5,
//					TotalRooms:     10,
//				}
//
//				mockStore.EXPECT().
//					GetPaymentByOrderId("order_124").
//					Return(payment, nil)
//
//				mockStore.EXPECT().
//					GetBookingById(int64(2)).
//					Return(booking, nil)
//
//				mockStore.EXPECT().
//					GetHotelById(int64(2)).
//					Return(hotel, nil)
//
//				mockStore.EXPECT().
//					UpdateHotelRooms(int64(2), 6). // 5 + 1 room returned
//					Return(nil)
//
//				mockStore.EXPECT().
//					UpdateBookingStatus(int64(2), store.FAILED).
//					Return(nil)
//
//				mockStore.EXPECT().
//					UpdatePaymentStatus(2, store.PAYMENT_FAILED).
//					Return(nil)
//			},
//			expectedStatus: http.StatusOK,
//		},
//		{
//			name:           "invalid JSON request",
//			requestBody:    "invalid json",
//			setupMocks:     func() {},
//			expectedStatus: http.StatusBadRequest,
//			expectedBody:   "invalid character",
//		},
//		{
//			name: "payment not found",
//			requestBody: pb.PaymentWebhookRequest{
//				Payload: &pb.PaymentWebhookRequest_Payload{
//					OrderId: "non_existent_order",
//					Status:  string(store.PAYMENT_SUCCESS),
//					Metadata: func() []byte {
//						metadata := payments.PaymentsMetadata{BookingID: 999}
//						data, _ := json.Marshal(metadata)
//						return data
//					}(),
//				},
//			},
//			setupMocks: func() {
//				mockStore.EXPECT().
//					GetPaymentByOrderId("non_existent_order").
//					Return(nil, assert.AnError)
//			},
//			expectedStatus: http.StatusNotFound,
//			expectedBody:   errorcodes.ErrPaymentNotFoundByOID,
//		},
//		{
//			name: "invalid metadata",
//			requestBody: pb.PaymentWebhookRequest{
//				Payload: &pb.PaymentWebhookRequest_Payload{
//					OrderId:  "order_125",
//					Status:   string(store.PAYMENT_SUCCESS),
//					Metadata: []byte("invalid json metadata"),
//				},
//			},
//			setupMocks: func() {
//				payment := &store.Payment{
//					ID:        3,
//					BookingID: 3,
//					OrderID:   "order_125",
//					Amount:    2000.0,
//					Status:    store.PAYMENT_PENDING,
//				}
//
//				mockStore.EXPECT().
//					GetPaymentByOrderId("order_125").
//					Return(payment, nil)
//			},
//			expectedStatus: http.StatusBadRequest,
//			expectedBody:   "Invalid metadata",
//		},
//		{
//			name: "booking not found",
//			requestBody: pb.PaymentWebhookRequest{
//				Payload: &pb.PaymentWebhookRequest_Payload{
//					OrderId: "order_126",
//					Status:  string(store.PAYMENT_SUCCESS),
//					Metadata: func() []byte {
//						metadata := payments.PaymentsMetadata{BookingID: 999}
//						data, _ := json.Marshal(metadata)
//						return data
//					}(),
//				},
//			},
//			setupMocks: func() {
//				payment := &store.Payment{
//					ID:        4,
//					BookingID: 999,
//					OrderID:   "order_126",
//					Amount:    1000.0,
//					Status:    store.PAYMENT_PENDING,
//				}
//
//				mockStore.EXPECT().
//					GetPaymentByOrderId("order_126").
//					Return(payment, nil)
//
//				mockStore.EXPECT().
//					GetBookingById(int64(999)).
//					Return(store.Booking{}, assert.AnError)
//			},
//			expectedStatus: http.StatusNotFound,
//			expectedBody:   errorcodes.ErrBookingNotFoundByID,
//		},
//		{
//			name: "update booking status fails",
//			requestBody: pb.PaymentWebhookRequest{
//				Payload: &pb.PaymentWebhookRequest_Payload{
//					OrderId: "order_127",
//					Status:  string(store.PAYMENT_SUCCESS),
//					Metadata: func() []byte {
//						metadata := payments.PaymentsMetadata{BookingID: 5}
//						data, _ := json.Marshal(metadata)
//						return data
//					}(),
//				},
//			},
//			setupMocks: func() {
//				payment := &store.Payment{
//					ID:        5,
//					BookingID: 5,
//					OrderID:   "order_127",
//					Amount:    1000.0,
//					Status:    store.PAYMENT_PENDING,
//				}
//
//				booking := store.Booking{
//					BookingID: 5,
//					HotelID:   1,
//					Status:    store.PENDING,
//				}
//
//				mockStore.EXPECT().
//					GetPaymentByOrderId("order_127").
//					Return(payment, nil)
//
//				mockStore.EXPECT().
//					GetBookingById(int64(5)).
//					Return(booking, nil)
//
//				mockStore.EXPECT().
//					UpdateBookingStatus(int64(5), store.CONFIRMED).
//					Return(assert.AnError)
//			},
//			expectedStatus: http.StatusInternalServerError,
//			expectedBody:   "Failed to update booking",
//		},
//		{
//			name: "update payment status fails",
//			requestBody: pb.PaymentWebhookRequest{
//				Payload: &pb.PaymentWebhookRequest_Payload{
//					OrderId: "order_128",
//					Status:  string(store.PAYMENT_SUCCESS),
//					Metadata: func() []byte {
//						metadata := payments.PaymentsMetadata{BookingID: 6}
//						data, _ := json.Marshal(metadata)
//						return data
//					}(),
//				},
//			},
//			setupMocks: func() {
//				payment := &store.Payment{
//					ID:        6,
//					BookingID: 6,
//					OrderID:   "order_128",
//					Amount:    1000.0,
//					Status:    store.PAYMENT_PENDING,
//				}
//
//				booking := store.Booking{
//					BookingID: 6,
//					HotelID:   1,
//					Status:    store.PENDING,
//				}
//
//				mockStore.EXPECT().
//					GetPaymentByOrderId("order_128").
//					Return(payment, nil)
//
//				mockStore.EXPECT().
//					GetBookingById(int64(6)).
//					Return(booking, nil)
//
//				mockStore.EXPECT().
//					UpdateBookingStatus(int64(6), store.CONFIRMED).
//					Return(nil)
//
//				mockStore.EXPECT().
//					UpdatePaymentStatus(6, store.PAYMENT_SUCCESS).
//					Return(assert.AnError)
//			},
//			expectedStatus: http.StatusInternalServerError,
//			expectedBody:   "Failed to update payment status",
//		},
//		{
//			name: "hotel not found for failed payment",
//			requestBody: pb.PaymentWebhookRequest{
//				Payload: &pb.PaymentWebhookRequest_Payload{
//					OrderId: "order_129",
//					Status:  string(store.PAYMENT_FAILED),
//					Metadata: func() []byte {
//						metadata := payments.PaymentsMetadata{BookingID: 7}
//						data, _ := json.Marshal(metadata)
//						return data
//					}(),
//				},
//			},
//			setupMocks: func() {
//				payment := &store.Payment{
//					ID:        7,
//					BookingID: 7,
//					OrderID:   "order_129",
//					Amount:    1000.0,
//					Status:    store.PAYMENT_PENDING,
//				}
//
//				booking := store.Booking{
//					BookingID: 7,
//					HotelID:   999,
//					Status:    store.PENDING,
//				}
//
//				mockStore.EXPECT().
//					GetPaymentByOrderId("order_129").
//					Return(payment, nil)
//
//				mockStore.EXPECT().
//					GetBookingById(int64(7)).
//					Return(booking, nil)
//
//				mockStore.EXPECT().
//					GetHotelById(int64(999)).
//					Return(store.Hotel{}, assert.AnError)
//			},
//			expectedStatus: http.StatusNotFound,
//			expectedBody:   "Hotel not found",
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			tt.setupMocks()
//
//			var reqBody []byte
//			var err error
//
//			if str, ok := tt.requestBody.(string); ok {
//				reqBody = []byte(str)
//			} else {
//				reqBody, err = json.Marshal(tt.requestBody)
//				assert.NoError(t, err)
//			}
//
//			req := httptest.NewRequest("POST", "/webhook", bytes.NewBuffer(reqBody))
//			req.Header.Set("Content-Type", "application/json")
//
//			recorder := httptest.NewRecorder()
//
//			service.PaymentWebhookHandler(recorder, req)
//
//			assert.Equal(t, tt.expectedStatus, recorder.Code)
//
//			if tt.expectedBody != "" {
//				assert.Contains(t, recorder.Body.String(), tt.expectedBody)
//			}
//		})
//	}
//}
//
//func TestPaymentWebhookHandler_EdgeCases(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	mockStore := mocks.NewMockStorageService(ctrl)
//	service := &Service{
//		storageService: mockStore,
//	}
//
//	t.Run("unknown payment status", func(t *testing.T) {
//		requestBody := pb.PaymentWebhookRequest{
//			Payload: &pb.PaymentWebhookRequest_Payload{
//				OrderId: "order_unknown",
//				Status:  "unknown_status",
//				Metadata: func() []byte {
//					metadata := payments.PaymentsMetadata{BookingID: 1}
//					data, _ := json.Marshal(metadata)
//					return data
//				}(),
//			},
//		}
//
//		payment := &store.Payment{
//			ID:        1,
//			BookingID: 1,
//			OrderID:   "order_unknown",
//			Amount:    1000.0,
//			Status:    store.PAYMENT_PENDING,
//		}
//
//		booking := store.Booking{
//			BookingID: 1,
//			HotelID:   1,
//			Status:    store.PENDING,
//		}
//
//		mockStore.EXPECT().
//			GetPaymentByOrderId("order_unknown").
//			Return(payment, nil)
//
//		mockStore.EXPECT().
//			GetBookingById(int64(1)).
//			Return(booking, nil)
//
//		reqBody, _ := json.Marshal(requestBody)
//		req := httptest.NewRequest("POST", "/webhook", bytes.NewBuffer(reqBody))
//		req.Header.Set("Content-Type", "application/json")
//
//		recorder := httptest.NewRecorder()
//
//		service.PaymentWebhookHandler(recorder, req)
//
//		assert.Equal(t, http.StatusOK, recorder.Code)
//	})
//}
