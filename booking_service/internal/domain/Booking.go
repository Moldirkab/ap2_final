package domain

import "time"

type Booking struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	CarID      string    `json:"car_id"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

const (
	StatusPending   = "PENDING"
	StatusConfirmed = "CONFIRMED"
	StatusCancelled = "CANCELLED"
	StatusCompleted = "COMPLETED"
)
