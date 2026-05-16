package usecase

import "booking-service/internal/domain"

type BookingRepository interface {
	Create(booking *domain.Booking) error
	GetByID(id string) (*domain.Booking, error)
	ListByUserID(userID string) ([]*domain.Booking, error)
	UpdateStatus(id string, status string) (*domain.Booking, error)
}

type UserClient interface {
	ValidateToken(token string) (string, error)
}

type CarData struct {
	ID          string
	PricePerDay float64
}

type CarClient interface {
	CheckAvailability(carID string) (bool, error)
	GetCar(carID string) (*domain.CarData, error)
	UpdateAvailability(carID string, available bool) error
}

type EventPublisher interface {
	PublishBookingCreated(booking *domain.Booking) error
	PublishBookingCancelled(booking *domain.Booking) error
	PublishBookingCompleted(booking *domain.Booking) error
	PublishBookingConfirmed(booking *domain.Booking) error
}
