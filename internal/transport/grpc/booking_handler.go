package grpc

import (
	"context"

	"booking-service/internal/domain"
	"booking-service/internal/usecase"

	bookingpb "github.com/Moldirkab/ap2_final_generated/bookingpb"
)

type BookingHandler struct {
	bookingpb.UnimplementedBookingServiceServer
	usecase *usecase.BookingUsecase
}

func NewBookingHandler(usecase *usecase.BookingUsecase) *BookingHandler {
	return &BookingHandler{usecase: usecase}
}

func (h *BookingHandler) CreateBooking(ctx context.Context, req *bookingpb.CreateBookingRequest) (*bookingpb.BookingResponse, error) {
	booking, err := h.usecase.CreateBooking(req.Token, req.CarId, req.StartDate, req.EndDate)
	if err != nil {
		return nil, err
	}

	return toResponse(booking), nil
}

func (h *BookingHandler) ConfirmBooking(
	ctx context.Context,
	req *bookingpb.ConfirmBookingRequest,
) (*bookingpb.BookingResponse, error) {

	booking, err := h.usecase.ConfirmBooking(req.BookingId)
	if err != nil {
		return nil, err
	}

	return toResponse(booking), nil
}

func (h *BookingHandler) CancelBooking(ctx context.Context, req *bookingpb.CancelBookingRequest) (*bookingpb.BookingResponse, error) {
	booking, err := h.usecase.CancelBooking(req.BookingId)
	if err != nil {
		return nil, err
	}

	return toResponse(booking), nil
}

func (h *BookingHandler) CompleteBooking(ctx context.Context, req *bookingpb.CompleteBookingRequest) (*bookingpb.BookingResponse, error) {
	booking, err := h.usecase.CompleteBooking(req.BookingId)
	if err != nil {
		return nil, err
	}

	return toResponse(booking), nil
}

func (h *BookingHandler) GetBooking(ctx context.Context, req *bookingpb.GetBookingRequest) (*bookingpb.BookingResponse, error) {
	booking, err := h.usecase.GetBooking(req.BookingId)
	if err != nil {
		return nil, err
	}

	return toResponse(booking), nil
}

func (h *BookingHandler) ListUserBookings(ctx context.Context, req *bookingpb.ListUserBookingsRequest) (*bookingpb.ListBookingsResponse, error) {
	bookings, err := h.usecase.ListUserBookings(req.Token)
	if err != nil {
		return nil, err
	}

	var response []*bookingpb.BookingResponse

	for _, booking := range bookings {
		response = append(response, toResponse(booking))
	}

	return &bookingpb.ListBookingsResponse{
		Bookings: response,
	}, nil
}

func toResponse(booking *domain.Booking) *bookingpb.BookingResponse {
	return &bookingpb.BookingResponse{
		Id:         booking.ID,
		UserId:     booking.UserID,
		CarId:      booking.CarID,
		StartDate:  booking.StartDate.Format("2006-01-02"),
		EndDate:    booking.EndDate.Format("2006-01-02"),
		TotalPrice: booking.TotalPrice,
		Status:     booking.Status,
		CreatedAt:  booking.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
