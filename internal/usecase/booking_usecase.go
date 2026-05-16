package usecase

import (
	"errors"
	"time"

	"booking-service/internal/domain"

	"github.com/google/uuid"
)

type BookingUsecase struct {
	repo      BookingRepository
	user      UserClient
	car       CarClient
	publisher EventPublisher
}

func NewBookingUsecase(
	repo BookingRepository,
	user UserClient,
	car CarClient,
	publisher EventPublisher,
) *BookingUsecase {
	return &BookingUsecase{
		repo:      repo,
		user:      user,
		car:       car,
		publisher: publisher,
	}
}

func (u *BookingUsecase) CreateBooking(token string, carID string, startDate string, endDate string) (*domain.Booking, error) {
	userID, err := u.user.ValidateToken(token)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, errors.New("invalid start date")
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, errors.New("invalid end date")
	}

	if !end.After(start) {
		return nil, errors.New("end date must be after start date")
	}

	if start.Before(time.Now().Truncate(24 * time.Hour)) {
		return nil, errors.New("start date cannot be in the past")
	}

	available, err := u.car.CheckAvailability(carID)
	if err != nil {
		return nil, err
	}

	if !available {
		return nil, errors.New("car is not available")
	}

	carData, err := u.car.GetCar(carID)
	if err != nil {
		return nil, err
	}

	days := int(end.Sub(start).Hours()/24) + 1
	totalPrice := float64(days) * carData.PricePerDay

	booking := &domain.Booking{
		ID:         uuid.New().String(),
		UserID:     userID,
		CarID:      carID,
		StartDate:  start,
		EndDate:    end,
		TotalPrice: totalPrice,
		Status:     domain.StatusPending,
		CreatedAt:  time.Now(),
	}

	if err := u.repo.Create(booking); err != nil {
		return nil, err
	}

	if err := u.publisher.PublishBookingCreated(booking); err != nil {
		return nil, err
	}

	return booking, nil
}

func (u *BookingUsecase) ConfirmBooking(bookingID string) (*domain.Booking, error) {
	booking, err := u.repo.GetByID(bookingID)
	if err != nil {
		return nil, err
	}

	if booking.Status != domain.StatusPending {
		return nil, errors.New("only pending bookings can be confirmed")
	}

	booking, err = u.repo.UpdateStatus(bookingID, domain.StatusConfirmed)
	if err != nil {
		return nil, err
	}

	if err := u.car.UpdateAvailability(booking.CarID, false); err != nil {
		return nil, err
	}

	if err := u.publisher.PublishBookingConfirmed(booking); err != nil {
		return nil, err
	}

	return booking, nil
}

func (u *BookingUsecase) CancelBooking(bookingID string) (*domain.Booking, error) {
	booking, err := u.repo.GetByID(bookingID)
	if err != nil {
		return nil, err
	}

	if booking.Status == domain.StatusCancelled {
		return nil, errors.New("booking already cancelled")
	}

	if booking.Status == domain.StatusCompleted {
		return nil, errors.New("completed booking cannot be cancelled")
	}

	previousStatus := booking.Status

	booking, err = u.repo.UpdateStatus(bookingID, domain.StatusCancelled)
	if err != nil {
		return nil, err
	}

	if previousStatus == domain.StatusConfirmed {
		if err := u.car.UpdateAvailability(booking.CarID, true); err != nil {
			return nil, err
		}
	}

	if err := u.publisher.PublishBookingCancelled(booking); err != nil {
		return nil, err
	}

	return booking, nil
}

func (u *BookingUsecase) CompleteBooking(bookingID string) (*domain.Booking, error) {
	booking, err := u.repo.GetByID(bookingID)
	if err != nil {
		return nil, err
	}

	if booking.Status != domain.StatusConfirmed {
		return nil, errors.New("only confirmed bookings can be completed")
	}

	booking, err = u.repo.UpdateStatus(bookingID, domain.StatusCompleted)
	if err != nil {
		return nil, err
	}

	if err := u.car.UpdateAvailability(booking.CarID, true); err != nil {
		return nil, err
	}

	if err := u.publisher.PublishBookingCompleted(booking); err != nil {
		return nil, err
	}

	return booking, nil
}

func (u *BookingUsecase) GetBooking(bookingID string) (*domain.Booking, error) {
	return u.repo.GetByID(bookingID)
}

func (u *BookingUsecase) ListUserBookings(token string) ([]*domain.Booking, error) {
	userID, err := u.user.ValidateToken(token)
	if err != nil {
		return nil, errors.New("invalid token")
	}

	return u.repo.ListByUserID(userID)
}
