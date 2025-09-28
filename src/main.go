package main

import (
	"database/sql"
	"fmt"
	"hotel-system/src/routes"
	"hotel-system/src/services"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/stripe/stripe-go/v82"
)

func main() {
	mux := http.NewServeMux()

	//connStr := os.Getenv("SQL_CONNECTION_STRING")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error while loading env data: ", err)
		return
	}
	connStr := os.Getenv("SQL_CONNECTION_STRING")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	stripeKey := os.Getenv("STRIPE_API_KEY")
	if stripeKey == "" {
		log.Fatal("STRIPE_API_KEY environment variable not set")
	}
	stripe.Key = stripeKey

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	service := services.NewService(db)
	routes.RegisterRoutes(mux, service)

	server := http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       15 * time.Second,
	}

	fmt.Println("Server starting on port 8080...")
	err = server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
