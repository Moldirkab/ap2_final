import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { listBookings, cancelBooking, confirmBooking, completeBooking, Booking } from "../api/bookings";
import { useAuth } from "../context/AuthContext";
import { useNotification } from "../context/NotificationContext";
import LoadingSpinner from "../components/LoadingSpinner";

const statusColors: Record<string, string> = {
  PENDING: "bg-yellow-100 text-yellow-700",
  CONFIRMED: "bg-blue-100 text-blue-700",
  CANCELLED: "bg-red-100 text-red-700",
  COMPLETED: "bg-green-100 text-green-700",
};

const BookingHistory: React.FC = () => {
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [loading, setLoading] = useState(true);
  const { token } = useAuth();
  const { notify } = useNotification();
  const navigate = useNavigate();

  useEffect(() => {
    if (!token) {
      navigate("/login");
      return;
    }
    fetchBookings();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [token]);

  const fetchBookings = async () => {
    setLoading(true);
    try {
      const res = await listBookings();
      setBookings(res.data || []);
    } catch (err: any) {
      notify(err.response?.data?.error || "Failed to load bookings", "error");
    } finally {
      setLoading(false);
    }
  };

  const handleCancel = async (id: string) => {
    if (!window.confirm("Are you sure you want to cancel this booking?")) return;
    try {
      await cancelBooking(id);
      notify("Booking cancelled", "success");
      fetchBookings();
    } catch (err: any) {
      notify(err.response?.data?.error || "Cancel failed", "error");
    }
  };

  const handleConfirm = async (id: string) => {
    try {
      await confirmBooking(id);
      notify("Booking confirmed!", "success");
      fetchBookings();
    } catch (err: any) {
      notify(err.response?.data?.error || "Confirm failed", "error");
    }
  };

  const handleComplete = async (id: string) => {
    try {
      await completeBooking(id);
      notify("Booking completed!", "success");
      fetchBookings();
    } catch (err: any) {
      notify(err.response?.data?.error || "Complete failed", "error");
    }
  };

  if (loading) return <LoadingSpinner />;

  return (
    <div className="max-w-5xl mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold text-gray-800 mb-6">My Bookings</h1>

      {bookings.length === 0 ? (
        <div className="text-center py-16 bg-white rounded-xl shadow">
          <p className="text-lg text-gray-500 mb-4">No bookings yet</p>
          <button
            onClick={() => navigate("/cars")}
            className="bg-blue-600 text-white px-6 py-2 rounded-lg font-semibold hover:bg-blue-700 transition"
          >
            Browse Cars
          </button>
        </div>
      ) : (
        <div className="space-y-4">
          {bookings.map((b) => (
            <div
              key={b.id}
              className="bg-white rounded-xl shadow p-6 hover:shadow-md transition"
            >
              <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                <div className="flex-1">
                  <div className="flex items-center gap-3 mb-2">
                    <h3 className="font-semibold text-gray-800">
                      Booking #{b.id.slice(0, 8)}
                    </h3>
                    <span
                      className={`text-xs px-2 py-1 rounded-full font-medium ${
                        statusColors[b.status] || "bg-gray-100 text-gray-700"
                      }`}
                    >
                      {b.status}
                    </span>
                  </div>
                  <div className="grid sm:grid-cols-3 gap-2 text-sm text-gray-500">
                    <p>
                      Car ID: <span className="text-gray-700">{b.car_id}</span>
                    </p>
                    <p>
                      From:{" "}
                      <span className="text-gray-700">
                        {formatDate(b.start_date)}
                      </span>
                    </p>
                    <p>
                      To:{" "}
                      <span className="text-gray-700">
                        {formatDate(b.end_date)}
                      </span>
                    </p>
                  </div>
                </div>
                <div className="text-right">
                  <p className="text-xl font-bold text-blue-600">
                    ${b.total_price.toFixed(2)}
                  </p>
                  <div className="mt-2 flex gap-2 justify-end flex-wrap">
                    {b.status === "PENDING" && (
                      <button
                        onClick={() => handleConfirm(b.id)}
                        className="text-sm bg-green-600 text-white px-3 py-1 rounded-lg hover:bg-green-700 transition"
                      >
                        Confirm
                      </button>
                    )}
                    {b.status === "CONFIRMED" && (
                      <button
                        onClick={() => handleComplete(b.id)}
                        className="text-sm bg-blue-600 text-white px-3 py-1 rounded-lg hover:bg-blue-700 transition"
                      >
                        Complete
                      </button>
                    )}
                    {(b.status === "PENDING" || b.status === "CONFIRMED") && (
                      <button
                        onClick={() => handleCancel(b.id)}
                        className="text-sm bg-red-100 text-red-600 px-3 py-1 rounded-lg hover:bg-red-200 transition"
                      >
                        Cancel
                      </button>
                    )}
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

const formatDate = (d: string): string => {
  try {
    return new Date(d).toLocaleDateString();
  } catch {
    return d;
  }
};

export default BookingHistory;
