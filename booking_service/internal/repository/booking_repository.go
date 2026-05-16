package repository

import (
	"database/sql"

	"booking-service/internal/domain"
)

type BookingRepository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(booking *domain.Booking) error {
	query := `
		INSERT INTO bookings 
		(id, user_id, car_id, start_date, end_date, total_price, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(
		query,
		booking.ID,
		booking.UserID,
		booking.CarID,
		booking.StartDate,
		booking.EndDate,
		booking.TotalPrice,
		booking.Status,
		booking.CreatedAt,
	)

	return err
}

func (r *BookingRepository) GetByID(id string) (*domain.Booking, error) {
	query := `
		SELECT id, user_id, car_id, start_date, end_date, total_price, status, created_at
		FROM bookings
		WHERE id = $1
	`

	booking := &domain.Booking{}

	err := r.db.QueryRow(query, id).Scan(
		&booking.ID,
		&booking.UserID,
		&booking.CarID,
		&booking.StartDate,
		&booking.EndDate,
		&booking.TotalPrice,
		&booking.Status,
		&booking.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return booking, nil
}

func (r *BookingRepository) ListByUserID(userID string) ([]*domain.Booking, error) {
	query := `
		SELECT id, user_id, car_id, start_date, end_date, total_price, status, created_at
		FROM bookings
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*domain.Booking

	for rows.Next() {
		booking := &domain.Booking{}

		err := rows.Scan(
			&booking.ID,
			&booking.UserID,
			&booking.CarID,
			&booking.StartDate,
			&booking.EndDate,
			&booking.TotalPrice,
			&booking.Status,
			&booking.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		bookings = append(bookings, booking)
	}

	return bookings, nil
}

func (r *BookingRepository) UpdateStatus(id string, status string) (*domain.Booking, error) {
	query := `
		UPDATE bookings
		SET status = $1
		WHERE id = $2
		RETURNING id, user_id, car_id, start_date, end_date, total_price, status, created_at
	`

	booking := &domain.Booking{}

	err := r.db.QueryRow(query, status, id).Scan(
		&booking.ID,
		&booking.UserID,
		&booking.CarID,
		&booking.StartDate,
		&booking.EndDate,
		&booking.TotalPrice,
		&booking.Status,
		&booking.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return booking, nil
}
