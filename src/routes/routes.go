package routes

import (
	"context"
	"hotel-system/src/services"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, service *services.Service) {
	mux.HandleFunc("/login", service.Login)
	mux.HandleFunc("/register", service.Register)

	mux.HandleFunc("/addHotel", Middleware(service.AddHotel))
	mux.HandleFunc("/getHotelsList", Middleware(service.GetHotelsList))
	mux.HandleFunc("/getHotelById", Middleware(service.GetHotelById))

	mux.HandleFunc("/bookHotel", Middleware(service.BookHotel))
	mux.HandleFunc("/getBookingById", Middleware(service.GetBookingDetailsById))
	mux.HandleFunc("/paymentStatus", Middleware(service.GetPaymentStatus))

	mux.HandleFunc("/payment-webhook", Middleware(service.PaymentWebhookHandler))

	// New route that redirects clients to the checkout URL
	mux.HandleFunc("/test/client/checkoutUrl", CORSMiddleware(service.ClientCheckoutRedirect))

}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//tokenStr := r.Header.Get("Authorization")
		//claims, err := utils.ValidateJWT(tokenStr)
		//if err != nil {
		//	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		//	return
		//}

		//ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		ctx := context.WithValue(r.Context(), "user_id", 6)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func CORSMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return CORSMiddleware(AuthMiddleware(next))
}
