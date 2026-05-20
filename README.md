# Car Rental Microservices Project

## Project Overview

This project implements a **Car Rental System** using a microservices architecture with **gRPC**, **NATS**, **Redis**, and **PostgreSQL**. It consists of **three microservices** and an **API Gateway**, allowing users to register, manage cars, and handle bookings with secure authentication and caching mechanisms.

**Team Responsibilities:**

* **Person 1 (Asem):** User/Auth Service (gRPC, JWT, user management)
* **Person 2 (Gulzada):** Car/Fleet Service (gRPC, Redis caching, car availability)
* **Person 3 (Moldir):** Booking/Rental Service (gRPC, business logic, NATS events)
* **API Gateway:** Handled by all members

---

## Architecture

```
┌─────────────┐      gRPC      ┌──────────────┐
│  User/Auth  │ <────────────> │ Booking/Rental│
│   Service   │                │   Service    │
└─────────────┘                └──────────────┘
        │ gRPC/NATS                       │ gRPC/NATS
        │                                 │
        ▼                                 ▼
  PostgreSQL (Users)                 PostgreSQL (Bookings)
  
┌───────────────┐
│ Car/Fleet     │
│ Service       │
└───────────────┘
        │
        ▼
    PostgreSQL (Cars)
        │
      Redis Cache
```

### Inter-Service Communication

* **gRPC**: service-to-service calls
* **NATS**: event-driven notifications (`booking.created`, `booking.cancelled`, `booking.completed`)
* **Redis**: caching of cars list, popular cars, and availability
* **JWT**: authentication and authorization

---

## Microservices

### 1. User/Auth Service

**Responsibilities:**

* User registration and login
* JWT token issuance and validation
* Role management (`admin` / `customer`)
* Email notifications (welcome email)

**Database:**

* Table `users`:

  * `id`
  * `email`
  * `password_hash`
  * `role`
  * `created_at`

**gRPC API:**

* `Register(email, password)` → create user
* `Login(email, password)` → return JWT
* `GetUser(user_id)` → get user profile
* `ValidateToken(token)` → validate JWT

**Notes:**

* Only this service handles passwords.
* Provides `ValidateToken()` for other services.

---

### 2. Car/Fleet Service

**Responsibilities:**

* Manage cars (add, list, detail)
* Check and update availability
* Caching with Redis to reduce DB load

**Database:**

* Table `cars`:

  * `id`
  * `brand`
  * `model`
  * `price_per_day`
  * `is_available`
  * `created_at`

**gRPC API:**

* `AddCar(...)`
* `GetCar(car_id)`
* `ListCars()`
* `CheckAvailability(car_id)`
* `UpdateAvailability(car_id, status)`

**Redis Cache:**

* Stores car list, popular cars, and availability

**Notes:**

* Does not handle bookings or users.

---

### 3. Booking/Rental Service

**Responsibilities:**

* Create, cancel, and complete bookings
* Calculate total price
* Publish events via NATS

**Database:**

* Table `bookings`:

  * `id`
  * `user_id`
  * `car_id`
  * `start_date`
  * `end_date`
  * `total_price`
  * `status` (`pending` / `cancelled` / `completed`)

**gRPC API:**

* `CreateBooking(user_id, car_id, start_date, end_date)`
* `CancelBooking(booking_id)`
* `CompleteBooking(booking_id)`
* `GetBooking(booking_id)`
* `ListUserBookings(user_id)`

**Business Logic:**

* `total_price = days * price_per_day`
* Publishes NATS events for booking lifecycle

**Notes:**

* Does not manage users or car details
* Depends on User Service and Car Service for validation and pricing

---

## API Gateway

* Routes requests from clients to the appropriate microservice
* Handles authentication and authorization

---

## Technologies Used

* **gRPC**: service communication
* **NATS**: event bus
* **Redis**: caching layer
* **PostgreSQL**: database for users, cars, and bookings
* **JWT**: authentication
* **Email SMTP**: for sending emails

---

## Requirements Coverage

| Requirement                   | Coverage                                            |
| ----------------------------- | --------------------------------------------------- |
| Clean Architecture            | Implemented across all services                   |
| At least 12 gRPC Endpoints    | Provided in three microservices                   |
| Usage of Message Queue (NATS) | Booking events published                          |
| Usage of Databases & Caches   | PostgreSQL & Redis in Car Service                 |
| Sending Emails                | Welcome emails via SMTP in User Service           |
| Testing                       | Unit and Integration tests implemented            |
| Frontend (Bonus)              | JS/React frontend can interact via API Gateway |

---

## Example Workflow

1. User logs in → receives JWT
2. User requests a booking
3. Booking Service:

   * Validates JWT via User Service
   * Checks car availability via Car Service
   * Gets car price via Car Service
   * Creates booking and saves to DB
   * Publishes event to NATS
4. Car Service caches updates
5. User receives confirmation

---

## Installation & Running

1. Clone repository:

```bash
git clone https://github.com/Moldirkab/ap2_final.git
cd ap2_final
```

2. Start PostgreSQL, Redis, and NATS services

3. Run each microservice individually:

```bash
# User/Auth Service
cd user-service
go run main.go

# Car/Fleet Service
cd car-service
go run main.go

# Booking Service
cd booking-service
go run main.go
```

---

## Testing

* Unit tests and integration tests are provided in each service folder.
* Run tests using:

```bash
go test ./...
```

---

## Example gRPC Calls

```bash
# Register User
grpcurl -d '{"email": "user@test.com", "password": "pass123"}' localhost:50051 UserService/Register

# List Cars
grpcurl -d '{}' localhost:50052 CarService/ListCars

# Create Booking
grpcurl -d '{"user_id": 1, "car_id": 2, "start_date": "2026-05-22", "end_date": "2026-05-24"}' localhost:50053 BookingService/CreateBooking
```
