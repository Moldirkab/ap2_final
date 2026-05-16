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
	return u.repo.Create(ctx, car)
}

// Get Car
func (u *CarUsecase) GetCar(ctx context.Context, id int64) (*domain.Car, error) {

	cacheKey := fmt.Sprintf("car:%d", id)

	// Try Redis
	cachedCar, err := u.cache.Get(ctx, cacheKey)

	if err == nil {
		var car domain.Car

		if json.Unmarshal([]byte(cachedCar), &car) == nil {
			fmt.Println("Loaded from Redis")
			return &car, nil
		}
	}

	// DB
	car, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Save to Redis
	carJSON, _ := json.Marshal(car)

	_ = u.cache.Set(ctx, cacheKey, string(carJSON), 10*time.Minute)

	fmt.Println("Loaded from PostgreSQL")

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

	return car.Status == "available", car.Status, nil
}

// Change Status
func (u *CarUsecase) ChangeStatus(ctx context.Context, id int64, status string) (*domain.Car, error) {

	car, err := u.repo.ChangeStatus(ctx, id, status)
	if err != nil {
		return nil, err
	}

	cacheKey := fmt.Sprintf("car:%d", id)

	_ = u.cache.Delete(ctx, cacheKey)

	return car, nil
}
