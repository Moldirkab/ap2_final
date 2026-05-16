package tests

import (
	"car_service/internal/cache"
	"car_service/internal/domain"
	"car_service/internal/repository"
	"car_service/internal/usecase"
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func setupUsecase(t *testing.T) *usecase.CarUsecase {
	t.Helper()
	_ = godotenv.Load("../.env")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to connect db: %v", err)
	}

	if err = db.Ping(); err != nil {
		t.Fatalf("failed to ping db: %v", err)
	}

	repo := repository.NewCarRepository(db)
	redisCache := cache.NewRedisCache()
	return usecase.NewCarUsecase(repo, redisCache)
}

func createTestCar(t *testing.T, uc *usecase.CarUsecase) *domain.Car {
	t.Helper()
	car := &domain.Car{
		Brand:       "TestBrand",
		Model:       "TestModel",
		Year:        2023,
		PlateNumber: fmt.Sprintf("TEST%d", time.Now().UnixNano()),
		PricePerDay: 50.0,
		Status:      "available",
	}
	created, err := uc.CreateCar(context.Background(), car)
	if err != nil {
		t.Fatalf("failed to create test car: %v", err)
	}
	return created
}

// --- Integration Tests ---

func TestIntegration_CreateCar(t *testing.T) {
	uc := setupUsecase(t)

	car := &domain.Car{
		Brand:       "Toyota",
		Model:       "Camry",
		Year:        2023,
		PlateNumber: fmt.Sprintf("INT%d", os.Getpid()),
		PricePerDay: 100.0,
		Status:      "available",
	}

	created, err := uc.CreateCar(context.Background(), car)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if created.ID == 0 {
		t.Error("expected ID to be set")
	}
	if created.Brand != "Toyota" {
		t.Errorf("expected Toyota, got %s", created.Brand)
	}
}

func TestIntegration_GetCar(t *testing.T) {
	uc := setupUsecase(t)
	created := createTestCar(t, uc)

	car, err := uc.GetCar(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if car.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, car.ID)
	}
}

func TestIntegration_GetCar_NotFound(t *testing.T) {
	uc := setupUsecase(t)

	_, err := uc.GetCar(context.Background(), 999999)
	if err == nil {
		t.Fatal("expected error for nonexistent car")
	}
}

func TestIntegration_ListCars(t *testing.T) {
	uc := setupUsecase(t)
	createTestCar(t, uc)

	cars, err := uc.ListCars(context.Background(), "", "", 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(cars) == 0 {
		t.Error("expected at least one car")
	}
}

func TestIntegration_CheckAvailability_Available(t *testing.T) {
	uc := setupUsecase(t)
	created := createTestCar(t, uc)

	available, status, err := uc.CheckAvailability(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !available {
		t.Error("expected car to be available")
	}
	if status != "available" {
		t.Errorf("expected 'available', got %s", status)
	}
}

func TestIntegration_CheckAvailability_Booked(t *testing.T) {
	uc := setupUsecase(t)
	created := createTestCar(t, uc)

	_, err := uc.ChangeStatus(context.Background(), created.ID, "booked")
	if err != nil {
		t.Fatalf("failed to change status: %v", err)
	}

	available, _, err := uc.CheckAvailability(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if available {
		t.Error("expected car to not be available")
	}
}

func TestIntegration_ChangeStatus(t *testing.T) {
	uc := setupUsecase(t)
	created := createTestCar(t, uc)

	updated, err := uc.ChangeStatus(context.Background(), created.ID, "booked")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updated.Status != "booked" {
		t.Errorf("expected 'booked', got %s", updated.Status)
	}
}

func TestIntegration_UpdateCar(t *testing.T) {
	uc := setupUsecase(t)
	created := createTestCar(t, uc)

	created.Brand = "Honda"
	created.PricePerDay = 75.0

	updated, err := uc.UpdateCar(context.Background(), created)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if updated.Brand != "Honda" {
		t.Errorf("expected Honda, got %s", updated.Brand)
	}
	if updated.PricePerDay != 75.0 {
		t.Errorf("expected 75.0, got %f", updated.PricePerDay)
	}
}

func TestIntegration_DeleteCar(t *testing.T) {
	uc := setupUsecase(t)
	created := createTestCar(t, uc)

	err := uc.DeleteCar(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = uc.GetCar(context.Background(), created.ID)
	if err == nil {
		t.Fatal("expected error after deletion")
	}
}
