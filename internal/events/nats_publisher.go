package events

import (
	"encoding/json"

	"booking-service/internal/domain"

	"github.com/nats-io/nats.go"
)

type NatsPublisher struct {
	conn *nats.Conn
}

func NewNatsPublisher(url string) (*NatsPublisher, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return &NatsPublisher{conn: conn}, nil
}

func (p *NatsPublisher) PublishBookingCreated(booking *domain.Booking) error {
	return p.publish("booking.created", booking)
}

func (p *NatsPublisher) PublishBookingCancelled(booking *domain.Booking) error {
	return p.publish("booking.cancelled", booking)
}

func (p *NatsPublisher) PublishBookingCompleted(booking *domain.Booking) error {
	return p.publish("booking.completed", booking)
}

func (p *NatsPublisher) publish(subject string, booking *domain.Booking) error {
	data, err := json.Marshal(booking)
	if err != nil {
		return err
	}

	return p.conn.Publish(subject, data)
}

func (p *NatsPublisher) PublishBookingConfirmed(booking *domain.Booking) error {
	return p.publish("booking.confirmed", booking)
}
