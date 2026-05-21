package repository

import (
	"car_service/internal/domain"
	"context"
	"database/sql"
	"fmt"
)

type CarRepository struct {
	db *sql.DB
}

func NewCarRepository(db *sql.DB) *CarRepository {
	return &CarRepository{db: db}
}

func (r *CarRepository) Create(ctx context.Context, car *domain.Car) (*domain.Car, error) {
	query := `
		INSERT INTO cars (brand, model, year, plate_number, price_per_day, status, photo)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	if car.Status == "" {
		car.Status = "available"
	}

	err := r.db.QueryRowContext(
		ctx,
		query,
		car.Brand,
		car.Model,
		car.Year,
		car.PlateNumber,
		car.PricePerDay,
		car.Status,
		car.Photo,
	).Scan(&car.ID)

	if err != nil {
		return nil, err
	}

	return car, nil
}

func (r *CarRepository) GetByID(ctx context.Context, id int64) (*domain.Car, error) {
	query := `
		SELECT id, brand, model, year, plate_number, price_per_day, status, photo
		FROM cars
		WHERE id = $1
	`

	car := &domain.Car{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&car.ID,
		&car.Brand,
		&car.Model,
		&car.Year,
		&car.PlateNumber,
		&car.PricePerDay,
		&car.Status,
		&car.Photo,
	)

	if err != nil {
		return nil, err
	}

	return car, nil
}

func (r *CarRepository) List(
	ctx context.Context,
	brand string,
	status string,
	maxPrice float64,
) ([]*domain.Car, error) {

	query := `
		SELECT id, brand, model, year, plate_number, price_per_day, status, photo
		FROM cars
		WHERE 1=1
	`

	args := []interface{}{}
	argID := 1

	if brand != "" {
		query += ` AND brand = $` + fmt.Sprint(argID)
		args = append(args, brand)
		argID++
	}

	if status != "" {
		query += ` AND status = $` + fmt.Sprint(argID)
		args = append(args, status)
		argID++
	}

	if maxPrice > 0 {
		query += ` AND price_per_day <= $` + fmt.Sprint(argID)
		args = append(args, maxPrice)
		argID++
	}

	query += ` ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cars []*domain.Car

	for rows.Next() {
		car := &domain.Car{}

		err := rows.Scan(
			&car.ID,
			&car.Brand,
			&car.Model,
			&car.Year,
			&car.PlateNumber,
			&car.PricePerDay,
			&car.Status,
			&car.Photo,
		)

		if err != nil {
			return nil, err
		}

		cars = append(cars, car)
	}

	return cars, nil
}

func (r *CarRepository) Update(ctx context.Context, car *domain.Car) (*domain.Car, error) {
	query := `
		UPDATE cars
		SET brand = $1,
		    model = $2,
		    year = $3,
		    plate_number = $4,
		    price_per_day = $5,
		    status = $6,
		    photo = $7,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $8
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		car.Brand,
		car.Model,
		car.Year,
		car.PlateNumber,
		car.PricePerDay,
		car.Status,
		car.Photo,
		car.ID,
	)

	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, car.ID)
}

func (r *CarRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM cars WHERE id = $1`, id)
	return err
}

func (r *CarRepository) ChangeStatus(ctx context.Context, id int64, status string) (*domain.Car, error) {
	_, err := r.db.ExecContext(ctx, `
		UPDATE cars
		SET status = $1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`, status, id)

	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id)
}
