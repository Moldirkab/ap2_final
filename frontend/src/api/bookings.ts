import api from "./client";

export interface Booking {
  id: string;
  user_id: string;
  car_id: string;
  start_date: string;
  end_date: string;
  total_price: number;
  status: string;
  created_at: string;
}

export interface CreateBookingPayload {
  car_id: string;
  start_date: string;
  end_date: string;
}

export const createBooking = (payload: CreateBookingPayload) =>
  api.post<Booking>("/bookings", payload);

export const listBookings = () => api.get<Booking[]>("/bookings");

export const getBooking = (id: string) => api.get<Booking>(`/bookings/${id}`);

export const cancelBooking = (id: string) =>
  api.post<Booking>(`/bookings/${id}/cancel`);

export const confirmBooking = (id: string) =>
  api.post<Booking>(`/bookings/${id}/confirm`);

export const completeBooking = (id: string) =>
  api.post<Booking>(`/bookings/${id}/complete`);
