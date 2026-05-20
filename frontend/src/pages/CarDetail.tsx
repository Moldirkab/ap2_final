import React, { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { getCar, deleteCar, Car } from "../api/cars";
import { createBooking } from "../api/bookings";
import { useAuth } from "../context/AuthContext";
import { useNotification } from "../context/NotificationContext";
import LoadingSpinner from "../components/LoadingSpinner";

const CarDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [car, setCar] = useState<Car | null>(null);
  const [loading, setLoading] = useState(true);
  const [startDate, setStartDate] = useState("");
  const [endDate, setEndDate] = useState("");
  const [booking, setBooking] = useState(false);
  const [deleting, setDeleting] = useState(false);
  const { token, user } = useAuth();
  const { notify } = useNotification();
  const navigate = useNavigate();

  useEffect(() => {
    if (!id) return;
    getCar(parseInt(id))
      .then((res) => setCar(res.data))
      .catch(() => notify("Failed to load car details", "error"))
      .finally(() => setLoading(false));
  }, [id, notify]);

  const handleBook = async () => {
    if (!token) {
      notify("Please login to book a car", "error");
      navigate("/login");
      return;
    }
    if (!startDate || !endDate) {
      notify("Please select dates", "error");
      return;
    }
    if (new Date(endDate) <= new Date(startDate)) {
      notify("End date must be after start date", "error");
      return;
    }
    setBooking(true);
    try {
      await createBooking({
        car_id: id!,
        start_date: startDate,
        end_date: endDate,
      });
      notify("Booking created successfully!", "success");
      navigate("/bookings");
    } catch (err: any) {
      notify(err.response?.data?.error || "Booking failed", "error");
    } finally {
      setBooking(false);
    }
  };

  const handleDelete = async () => {
    if (!window.confirm("Are you sure you want to delete this car? This action cannot be undone.")) return;
    setDeleting(true);
    try {
      await deleteCar(parseInt(id!));
      notify("Car deleted successfully", "success");
      navigate("/cars");
    } catch (err: any) {
      notify(err.response?.data?.error || "Failed to delete car", "error");
    } finally {
      setDeleting(false);
    }
  };

  if (loading) return <LoadingSpinner />;
  if (!car) return <div className="text-center py-16 text-gray-500">Car not found</div>;

  const isAdmin = user?.role === "admin";

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <button
        onClick={() => navigate("/cars")}
        className="text-blue-600 hover:underline mb-6 inline-block"
      >
        &larr; Back to Cars
      </button>

      <div className="bg-white rounded-xl shadow-lg overflow-hidden">
        <div className="h-56 flex items-center justify-center overflow-hidden">
          {car.photo ? (
            <img
              src={car.photo}
              alt={`${car.brand} ${car.model}`}
              className="w-full h-full object-cover"
              onError={(e) => {
                (e.target as HTMLImageElement).style.display = "none";
                (e.target as HTMLImageElement).nextElementSibling?.classList.remove("hidden");
              }}
            />
          ) : null}
          <div className={`bg-gradient-to-br from-blue-100 to-blue-200 w-full h-full flex items-center justify-center ${car.photo ? "hidden" : ""}`}>
            <span className="text-8xl opacity-50">{"\uD83D\uDE97"}</span>
          </div>
        </div>

        <div className="p-8">
          <div className="flex justify-between items-start mb-4">
            <h1 className="text-3xl font-bold text-gray-800">
              {car.brand} {car.model}
            </h1>
            <span
              className={`text-sm px-3 py-1 rounded-full font-medium ${
                car.status === "available"
                  ? "bg-green-100 text-green-700"
                  : "bg-red-100 text-red-700"
              }`}
            >
              {car.status}
            </span>
          </div>

          <div className="grid sm:grid-cols-2 gap-6 mb-8">
            <div className="space-y-3">
              <DetailRow label="Year" value={car.year.toString()} />
              <DetailRow label="Plate Number" value={car.plate_number} />
              <DetailRow label="Brand" value={car.brand} />
              <DetailRow label="Model" value={car.model} />
            </div>
            <div className="bg-blue-50 rounded-xl p-6 text-center">
              <p className="text-sm text-gray-500 mb-1">Price per day</p>
              <p className="text-4xl font-bold text-blue-600">
                ${car.price_per_day.toFixed(2)}
              </p>
            </div>
          </div>

          {car.status === "available" && (
            <div className="border-t border-gray-200 pt-6">
              <h3 className="text-lg font-semibold text-gray-800 mb-4">
                Book this car
              </h3>
              <div className="grid sm:grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Start Date
                  </label>
                  <input
                    type="date"
                    value={startDate}
                    onChange={(e) => setStartDate(e.target.value)}
                    min={new Date().toISOString().split("T")[0]}
                    className="w-full border border-gray-300 rounded-lg px-4 py-2 focus:ring-2 focus:ring-blue-500 outline-none"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    End Date
                  </label>
                  <input
                    type="date"
                    value={endDate}
                    onChange={(e) => setEndDate(e.target.value)}
                    min={startDate || new Date().toISOString().split("T")[0]}
                    className="w-full border border-gray-300 rounded-lg px-4 py-2 focus:ring-2 focus:ring-blue-500 outline-none"
                  />
                </div>
              </div>
              {startDate && endDate && new Date(endDate) > new Date(startDate) && (
                <p className="text-sm text-gray-500 mb-4">
                  Estimated total:{" "}
                  <span className="font-semibold text-blue-600">
                    $
                    {(
                      car.price_per_day *
                      Math.ceil(
                        (new Date(endDate).getTime() -
                          new Date(startDate).getTime()) /
                          (1000 * 60 * 60 * 24)
                      )
                    ).toFixed(2)}
                  </span>
                </p>
              )}
              <button
                onClick={handleBook}
                disabled={booking}
                className="w-full sm:w-auto bg-blue-600 text-white px-8 py-3 rounded-lg font-semibold hover:bg-blue-700 transition disabled:opacity-50"
              >
                {booking ? "Booking..." : "Book Now"}
              </button>
            </div>
          )}

          {/* Admin actions */}
          {isAdmin && (
            <div className="border-t border-gray-200 pt-6 mt-6">
              <h3 className="text-lg font-semibold text-gray-800 mb-4">
                Admin Actions
              </h3>
              <div className="flex gap-3">
                <button
                  onClick={() => navigate(`/admin/cars/${car.id}/edit`)}
                  className="bg-yellow-500 text-white px-6 py-2.5 rounded-lg font-semibold hover:bg-yellow-600 transition"
                >
                  Edit Car
                </button>
                <button
                  onClick={handleDelete}
                  disabled={deleting}
                  className="bg-red-600 text-white px-6 py-2.5 rounded-lg font-semibold hover:bg-red-700 transition disabled:opacity-50"
                >
                  {deleting ? "Deleting..." : "Delete Car"}
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

const DetailRow: React.FC<{ label: string; value: string }> = ({
  label,
  value,
}) => (
  <div className="flex justify-between">
    <span className="text-gray-500 text-sm">{label}</span>
    <span className="text-gray-800 font-medium">{value}</span>
  </div>
);

export default CarDetail;
