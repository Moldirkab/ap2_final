package usecase

import (
	"testing"
	"time"

	"booking-service/internal/domain"
)

type mockRepo struct{}

func (m *mockRepo) Create(b *domain.Booking) error {
	return nil
}

func (m *mockRepo) GetByID(id string) (*domain.Booking, error) {
	return &domain.Booking{
		ID:     id,
		Status: domain.StatusPending,
		CarID:  "1",
	}, nil
}

func (m *mockRepo) UpdateStatus(id string, status string) (*domain.Booking, error) {
	return &domain.Booking{
		ID:     id,
		Status: status,
		CarID:  "1",
	}, nil
}

func (m *mockRepo) ListByUserID(userID string) ([]*domain.Booking, error) {
	return []*domain.Booking{}, nil
}

type mockUserClient struct{}

func (m *mockUserClient) ValidateToken(token string) (string, error) {
	return "user-1", nil
}

type mockCarClient struct{}

func (m *mockCarClient) CheckAvailability(carID string) (bool, error) {
	return true, nil
}

func (m *mockCarClient) GetCar(carID string) (*domain.CarData, error) {
	return &domain.CarData{
		ID:          carID,
		PricePerDay: 100,
	}, nil
}

func (m *mockCarClient) UpdateAvailability(carID string, available bool) error {
	return nil
}

type mockPublisher struct{}

func (m *mockPublisher) PublishBookingCreated(b *domain.Booking) error {
	return nil
}

func (m *mockPublisher) PublishBookingCancelled(b *domain.Booking) error {
	return nil
}

func (m *mockPublisher) PublishBookingCompleted(b *domain.Booking) error {
	return nil
}

func (m *mockPublisher) PublishBookingConfirmed(b *domain.Booking) error {
	return nil
}

func TestCreateBooking_Success(t *testing.T) {
	u := &BookingUsecase{
		repo:      &mockRepo{},
		user:      &mockUserClient{},
		car:       &mockCarClient{},
		publisher: &mockPublisher{},
	}

	start := time.Now().Add(24 * time.Hour)
	end := start.Add(48 * time.Hour)

	booking, err := u.CreateBooking(
		"token",
		"1",
		start.Format("2006-01-02"),
		end.Format("2006-01-02"),
	)

	if err != nil {
		t.Fatal(err)
	}

	if booking.Status != domain.StatusPending {
		t.Error("expected pending status")
	}

	if booking.TotalPrice != 300 {
		t.Error("expected total price 300")
	}
}

func TestCreateBooking_InvalidDates(t *testing.T) {
	u := &BookingUsecase{
		user: &mockUserClient{},
		car:  &mockCarClient{},
	}

	_, err := u.CreateBooking(
		"token",
		"1",
		"2026-05-10",
		"2026-05-01",
	)

	if err == nil {
		t.Error("expected invalid date error")
	}
}

func TestConfirmBooking(t *testing.T) {
	u := &BookingUsecase{
		repo:      &mockRepo{},
		car:       &mockCarClient{},
		publisher: &mockPublisher{},
	}

	booking, err := u.ConfirmBooking("booking-1")

	if err != nil {
		t.Fatal(err)
	}

	if booking.Status != domain.StatusConfirmed {
		t.Error("expected confirmed status")
	}
}

func TestCancelBooking(t *testing.T) {
	u := &BookingUsecase{
		repo:      &mockRepo{},
		car:       &mockCarClient{},
		publisher: &mockPublisher{},
	}

	booking, err := u.CancelBooking("booking-1")

	if err != nil {
		t.Fatal(err)
	}

	if booking.Status != domain.StatusCancelled {
		t.Error("expected cancelled status")
	}
}
