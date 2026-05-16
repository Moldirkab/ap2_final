package integration

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"booking-service/internal/domain"
	"booking-service/internal/repository"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sql.DB {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatal("failed to load .env:", err)
	}
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func TestCreateBookingIntegration(t *testing.T) {
	db := setupTestDB(t)

	repo := repository.NewBookingRepository(db)

	booking := &domain.Booking{
		ID:         "11111111-1111-1111-1111-111111111111",
		UserID:     "22222222-2222-2222-2222-222222222222",
		CarID:      "33333333-3333-3333-3333-333333333333",
		StartDate:  time.Now(),
		EndDate:    time.Now().Add(48 * time.Hour),
		TotalPrice: 500,
		Status:     domain.StatusPending,
		CreatedAt:  time.Now(),
	}

	err := repo.Create(booking)

	if err != nil {
		t.Fatal(err)
	}
}
