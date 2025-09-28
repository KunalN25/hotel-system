# Hotel System

This project is a backend system for managing hotel operations, built in Go. It provides a modular architecture for handling user authentication, bookings, payments, scheduling, and more. The system is designed for extensibility and maintainability, using best practices in Go development.

## Features Implemented

- **User Authentication**: Secure JWT-based authentication and user management.
- **Booking Management**: Endpoints and services for creating, updating, and managing hotel bookings.
- **Payment Processing**: Payment handling, idempotency, and webhook support for payment events.
- **Scheduling**: Automated schedulers for bookings and payments.
- **Protobuf & gRPC**: Protocol Buffers definitions for user and hotel system, with generated Go code for gRPC services.
- **Validation**: Input validation for various API endpoints.
- **Error Handling**: Centralized error codes and error handling utilities.
- **Constants & Types**: Well-organized constants and type definitions for strong typing and maintainability.
- **Serializers**: Serialization logic for API responses.
- **Store Layer**: Abstractions for data storage, including payment and idempotency logic, with mock implementations for testing.
- **Testing**: Unit tests for utilities and payment webhooks.
- **AI Prompts**: (Directory present, details not specified in code structure.)

## Project Structure

- `src/` — Main application source code
  - `constants/` — Project-wide constants
  - `errorcodes/` — Error code definitions
  - `payments/` — Payment logic and types
  - `pb/` — Generated protobuf Go files
  - `protos/` — Protobuf definitions
  - `routes/` — API route definitions
  - `schedulers/` — Booking and payment schedulers
  - `serializers/` — Serialization logic
  - `services/` — Business logic for auth, booking, payment, etc.
  - `store/` — Data storage abstractions and mocks
  - `types/` — Core type definitions
  - `utils/` — Utility functions (e.g., JWT, helpers)
  - `validators/` — Input validation logic
- `client/` — Client code (details not specified)
- `ai_prompts/` — AI prompt files (details not specified)
- `schema.sql` — Database schema
- `generate_protobufs.sh` — Script to generate protobuf files
- `go.mod`, `go.sum` — Go module files

## Getting Started

1. **Install dependencies:**
   ```sh
   go mod tidy
   ```
2. **Generate protobufs:**
   ```sh
   ./generate_protobufs.sh
   ```
3. **Run the application:**
   ```sh
   go run src/main.go
   ```

## Testing

Run unit tests with:
```sh
go test ./src/...
```

## License

MIT License (add your license details here)
