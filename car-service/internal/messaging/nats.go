package messaging

import (
	"car_service/internal/usecase"
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/nats-io/nats.go"
)

type BookingEvent struct {
	CarID string `json:"car_id"`
}

func StartNATSListener(carUsecase *usecase.CarUsecase) {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("NATS connection error: %v", err)
	}

	log.Println("Connected to NATS")

	_, err = nc.Subscribe("booking.created", func(msg *nats.Msg) {
		var event BookingEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Println(err)
			return
		}

		id, err := strconv.ParseInt(event.CarID, 10, 64)
		if err != nil {
			log.Println("invalid car id:", err)
			return
		}

		_, err = carUsecase.ChangeStatus(context.Background(), id, "booked")
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("Car %d marked as booked", id)
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = nc.Subscribe("booking.cancelled", func(msg *nats.Msg) {
		var event BookingEvent
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Println(err)
			return
		}

		id, err := strconv.ParseInt(event.CarID, 10, 64)
		if err != nil {
			log.Println("invalid car id:", err)
			return
		}

		_, err = carUsecase.ChangeStatus(context.Background(), id, "available")
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("Car %d marked as booked", id)
	})
	if err != nil {
		log.Fatal(err)
	}
}
