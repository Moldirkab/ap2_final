package grpc

import (
	"car_service/internal/domain"
	"car_service/internal/usecase"
	"context"
	"fmt"

	carpb "github.com/Moldirkab/ap2_final_car_service_generated/carpb"
)

type CarHandler struct {
	carpb.UnimplementedCarServiceServer
	usecase *usecase.CarUsecase
}

func NewCarHandler(usecase *usecase.CarUsecase) *CarHandler {
	return &CarHandler{usecase: usecase}
}

func toProtoCar(car *domain.Car) *carpb.Car {
	return &carpb.Car{
		Id:          car.ID,
		Brand:       car.Brand,
		Model:       car.Model,
		Year:        car.Year,
		PlateNumber: car.PlateNumber,
		PricePerDay: car.PricePerDay,
		Status:      car.Status,
		Photo:       car.Photo,
	}
}

func (h *CarHandler) CreateCar(ctx context.Context, req *carpb.CreateCarRequest) (*carpb.CarResponse, error) {
	car, err := h.usecase.CreateCar(ctx, &domain.Car{
		Brand:       req.Brand,
		Model:       req.Model,
		Year:        req.Year,
		PlateNumber: req.PlateNumber,
		PricePerDay: req.PricePerDay,
		Status:      "available",
		Photo:       req.Photo,
	})
	if err != nil {
		return nil, err
	}
	fmt.Println("PHOTO IN GRPC:", req.Photo)

	return &carpb.CarResponse{Car: toProtoCar(car)}, nil
}

func (h *CarHandler) GetCar(ctx context.Context, req *carpb.GetCarRequest) (*carpb.CarResponse, error) {
	car, err := h.usecase.GetCar(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &carpb.CarResponse{Car: toProtoCar(car)}, nil
}

func (h *CarHandler) ListCars(ctx context.Context, req *carpb.ListCarsRequest) (*carpb.ListCarsResponse, error) {
	cars, err := h.usecase.ListCars(ctx, req.Brand, req.Status, req.MaxPrice)
	if err != nil {
		return nil, err
	}

	var protoCars []*carpb.Car

	for _, car := range cars {
		protoCars = append(protoCars, toProtoCar(car))
	}

	return &carpb.ListCarsResponse{
		Cars: protoCars,
	}, nil
}

func (h *CarHandler) UpdateCar(ctx context.Context, req *carpb.UpdateCarRequest) (*carpb.CarResponse, error) {
	car, err := h.usecase.UpdateCar(ctx, &domain.Car{
		ID:          req.Id,
		Brand:       req.Brand,
		Model:       req.Model,
		Year:        req.Year,
		PlateNumber: req.PlateNumber,
		PricePerDay: req.PricePerDay,
		Status:      req.Status,
		Photo:       req.Photo,
	})
	if err != nil {
		return nil, err
	}
	fmt.Println("PHOTO IN UPDATE:", req.Photo)

	return &carpb.CarResponse{Car: toProtoCar(car)}, nil
}

func (h *CarHandler) DeleteCar(ctx context.Context, req *carpb.DeleteCarRequest) (*carpb.DeleteCarResponse, error) {
	err := h.usecase.DeleteCar(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &carpb.DeleteCarResponse{
		Success: true,
		Message: "Car deleted successfully",
	}, nil
}

func (h *CarHandler) CheckAvailability(ctx context.Context, req *carpb.CheckAvailabilityRequest) (*carpb.CheckAvailabilityResponse, error) {
	available, status, err := h.usecase.CheckAvailability(ctx, req.CarId)
	if err != nil {
		return nil, err
	}

	return &carpb.CheckAvailabilityResponse{
		Available: available,
		Status:    status,
	}, nil
}

func (h *CarHandler) ChangeCarStatus(ctx context.Context, req *carpb.ChangeCarStatusRequest) (*carpb.CarResponse, error) {
	car, err := h.usecase.ChangeStatus(ctx, req.CarId, req.Status)
	if err != nil {
		return nil, err
	}

	return &carpb.CarResponse{Car: toProtoCar(car)}, nil
}
