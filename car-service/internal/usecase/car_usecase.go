package usecase

import (
	"car_service/internal/cache"
	"car_service/internal/domain"
	"car_service/internal/repository"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type CarUsecase struct {
	repo  *repository.CarRepository
	cache *cache.RedisCache
}

func NewCarUsecase(repo *repository.CarRepository, cache *cache.RedisCache) *CarUsecase {
	return &CarUsecase{
		repo:  repo,
		cache: cache,
	}
}

// Create Car
func (u *CarUsecase) CreateCar(ctx context.Context, car *domain.Car) (*domain.Car, error) {
	if car.Status == "" {
		car.Status = domain.CarStatusAvailable
	}

	if car.Photo == "" {
		car.Photo = ""
	}

	return u.repo.Create(ctx, car)
}

// Get Car
func (u *CarUsecase) GetCar(ctx context.Context, id int64) (*domain.Car, error) {
	cacheKey := fmt.Sprintf("car:%d", id)

	cachedCar, err := u.cache.Get(ctx, cacheKey)
	if err == nil {
		var car domain.Car
		if json.Unmarshal([]byte(cachedCar), &car) == nil {
			fmt.Println("Loaded from Redis")
			return &car, nil
		}
	}

	car, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	fmt.Println("Loaded from PostgreSQL")

	carJSON, _ := json.Marshal(car)
	_ = u.cache.Set(ctx, cacheKey, string(carJSON), 10*time.Minute)

	return car, nil
}

// List Cars
func (u *CarUsecase) ListCars(
	ctx context.Context,
	brand string,
	status string,
	maxPrice float64,
) ([]*domain.Car, error) {
	return u.repo.List(ctx, brand, status, maxPrice)
}

// Update Car
func (u *CarUsecase) UpdateCar(ctx context.Context, car *domain.Car) (*domain.Car, error) {
	existing, err := u.repo.GetByID(ctx, car.ID)
	if err != nil {
		return nil, err
	}

	if car.Photo == "" {
		car.Photo = existing.Photo
	}
	if car.Status == "" {
		car.Status = existing.Status
	}
	if car.Brand == "" {
		car.Brand = existing.Brand
	}
	if car.Model == "" {
		car.Model = existing.Model
	}
	if car.PlateNumber == "" {
		car.PlateNumber = existing.PlateNumber
	}
	if car.Year == 0 {
		car.Year = existing.Year
	}
	if car.PricePerDay == 0 {
		car.PricePerDay = existing.PricePerDay
	}

	updatedCar, err := u.repo.Update(ctx, car)
	if err != nil {
		return nil, err
	}

	cacheKey := fmt.Sprintf("car:%d", car.ID)
	_ = u.cache.Delete(ctx, cacheKey)

	return updatedCar, nil
}

// Delete Car
func (u *CarUsecase) DeleteCar(ctx context.Context, id int64) error {
	err := u.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("car:%d", id)
	_ = u.cache.Delete(ctx, cacheKey)

	return nil
}

// Check Availability
func (u *CarUsecase) CheckAvailability(ctx context.Context, id int64) (bool, string, error) {
	car, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return false, "", err
	}

	return car.Status == domain.CarStatusAvailable, car.Status, nil
}

// Change Status
func (u *CarUsecase) ChangeStatus(ctx context.Context, id int64, status string) (*domain.Car, error) {
	if status == "" {
		status = domain.CarStatusAvailable
	}

	car, err := u.repo.ChangeStatus(ctx, id, status)
	if err != nil {
		return nil, err
	}

	cacheKey := fmt.Sprintf("car:%d", id)
	_ = u.cache.Delete(ctx, cacheKey)

	return car, nil
}
