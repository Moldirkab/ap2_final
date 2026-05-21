package domain

type Car struct {
	ID          int64
	Brand       string
	Model       string
	Year        int32
	PlateNumber string
	PricePerDay float64
	Status      string
	Photo       string
}

const (
	CarStatusAvailable = "available"
	CarStatusBooked    = "booked"
)
