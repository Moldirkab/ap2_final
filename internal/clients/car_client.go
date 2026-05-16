package clients

import (
	"context"
	"strconv"

	"booking-service/internal/domain"

	carpb "github.com/Moldirkab/ap2_final_car_service_generated/carpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type CarClient struct {
	client carpb.CarServiceClient
}

func NewCarClient(address string) (*CarClient, error) {
	conn, err := grpc.Dial(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &CarClient{
		client: carpb.NewCarServiceClient(conn),
	}, nil
}

func (c *CarClient) CheckAvailability(carID string) (bool, error) {
	id, err := strconv.ParseInt(carID, 10, 64)
	if err != nil {
		return false, err
	}

	res, err := c.client.GetCar(
		context.Background(),
		&carpb.GetCarRequest{
			Id: id,
		},
	)
	if err != nil {
		return false, err
	}

	return res.Car.Status == "available", nil
}

func (c *CarClient) GetCar(carID string) (*domain.CarData, error) {
	id, err := strconv.ParseInt(carID, 10, 64)
	if err != nil {
		return nil, err
	}

	res, err := c.client.GetCar(
		context.Background(),
		&carpb.GetCarRequest{
			Id: id,
		},
	)
	if err != nil {
		return nil, err
	}

	return &domain.CarData{
		ID:          strconv.FormatInt(res.Car.Id, 10),
		PricePerDay: res.Car.PricePerDay,
	}, nil
}

func (c *CarClient) UpdateAvailability(carID string, available bool) error {
	id, err := strconv.ParseInt(carID, 10, 64)
	if err != nil {
		return err
	}

	status := "available"
	if !available {
		status = "booked"
	}

	_, err = c.client.ChangeCarStatus(
		context.Background(),
		&carpb.ChangeCarStatusRequest{
			CarId:  id,
			Status: status,
		},
	)

	return err
}
