import api from "./client";

export interface Car {
  id: number;
  brand: string;
  model: string;
  year: number;
  plate_number: string;
  price_per_day: number;
  status: string;
  photo?: string;
}

export interface CarFilters {
  brand?: string;
  status?: string;
  max_price?: number;
}

export interface CreateCarPayload {
  brand: string;
  model: string;
  year: number;
  plate_number: string;
  price_per_day: number;
 photo?: string;
}

export const listCars = (filters?: CarFilters) =>
    api.get<Car[]>("/cars", { params: filters });

export const getCar = (id: number) =>
    api.get<Car>(`/cars/${id}`);

export const createCar = (car: CreateCarPayload) =>
    api.post<Car>("/cars", car);

export const updateCar = (id: number, car: Partial<Car>) =>
    api.put<Car>(`/cars/${id}`, car);

export const deleteCar = (id: number) =>
    api.delete(`/cars/${id}`);

export const checkAvailability = (id: number) =>
    api.get<{ available: boolean; status: string }>(`/cars/${id}/availability`);