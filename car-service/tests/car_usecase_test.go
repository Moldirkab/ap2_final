package tests

import (
	"car_service/internal/domain"
	"context"
	"testing"
)

type MockCarRepository struct{}

func (m *MockCarRepository) Create(ctx context.Context, car *domain.Car) (*domain.Car, error) {
	car.ID = 1
	return car, nil
}

func TestCreateCar(t *testing.T) {

	mockRepo := &MockCarRepository{}

	car := &domain.Car{
		Brand:       "Toyota",
		Model:       "Camry",
		Year:        2023,
		PlateNumber: "123AAA",
		PricePerDay: 100,
		Status:      "available",
	}

	createdCar, err := mockRepo.Create(context.Background(), car)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if createdCar.ID != 1 {
		t.Errorf("expected ID 1, got %d", createdCar.ID)
	}

	if createdCar.Brand != "Toyota" {
		t.Errorf("expected Toyota, got %s", createdCar.Brand)
	}
}

func TestCheckAvailability(t *testing.T) {

	car := &domain.Car{
		Status: "available",
	}

	if car.Status != "available" {
		t.Errorf("car should be available")
	}
}
