import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { listCars, Car, CarFilters } from "../api/cars";
import { useNotification } from "../context/NotificationContext";
import LoadingSpinner from "../components/LoadingSpinner";

const CarList: React.FC = () => {
  const [cars, setCars] = useState<Car[]>([]);
  const [loading, setLoading] = useState(true);
  const [filters, setFilters] = useState<CarFilters>({});
  const [brandInput, setBrandInput] = useState("");
  const [statusFilter, setStatusFilter] = useState("");
  const [maxPriceInput, setMaxPriceInput] = useState("");
  const { notify } = useNotification();

  const fetchCars = async (f?: CarFilters) => {
    setLoading(true);
    try {
      const res = await listCars(f || filters);
      setCars(res.data || []);
    } catch (err: any) {
      notify(err.response?.data?.error || "Failed to load cars", "error");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCars();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const applyFilters = () => {
    const f: CarFilters = {};
    if (brandInput) f.brand = brandInput;
    if (statusFilter) f.status = statusFilter;
    if (maxPriceInput) f.max_price = parseFloat(maxPriceInput);
    setFilters(f);
    fetchCars(f);
  };

  const clearFilters = () => {
    setBrandInput("");
    setStatusFilter("");
    setMaxPriceInput("");
    setFilters({});
    fetchCars({});
  };

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold text-gray-800 mb-6">Available Cars</h1>

      {/* Filters */}
      <div className="bg-white rounded-xl shadow p-6 mb-8">
        <div className="grid sm:grid-cols-4 gap-4">
          <input
            type="text"
            placeholder="Brand"
            value={brandInput}
            onChange={(e) => setBrandInput(e.target.value)}
            className="border border-gray-300 rounded-lg px-4 py-2 focus:ring-2 focus:ring-blue-500 outline-none"
          />
          <select
            value={statusFilter}
            onChange={(e) => setStatusFilter(e.target.value)}
            className="border border-gray-300 rounded-lg px-4 py-2 focus:ring-2 focus:ring-blue-500 outline-none"
          >
            <option value="">All Status</option>
            <option value="available">Available</option>
            <option value="booked">Booked</option>
          </select>
          <input
            type="number"
            placeholder="Max Price/Day"
            value={maxPriceInput}
            onChange={(e) => setMaxPriceInput(e.target.value)}
            className="border border-gray-300 rounded-lg px-4 py-2 focus:ring-2 focus:ring-blue-500 outline-none"
          />
          <div className="flex gap-2">
            <button
              onClick={applyFilters}
              className="bg-blue-600 text-white px-4 py-2 rounded-lg font-semibold hover:bg-blue-700 transition flex-1"
            >
              Search
            </button>
            <button
              onClick={clearFilters}
              className="bg-gray-200 text-gray-700 px-4 py-2 rounded-lg hover:bg-gray-300 transition"
            >
              Clear
            </button>
          </div>
        </div>
      </div>

      {loading ? (
        <LoadingSpinner />
      ) : cars.length === 0 ? (
        <div className="text-center py-16 text-gray-500">
          <p className="text-lg">No cars found</p>
          <p className="text-sm mt-1">Try adjusting your filters</p>
        </div>
      ) : (
        <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {cars.map((car) => (
            <CarCard key={car.id} car={car} />
          ))}
        </div>
      )}
    </div>
  );
};

const CarCard: React.FC<{ car: Car }> = ({ car }) => (
  <Link
    to={`/cars/${car.id}`}
    className="bg-white rounded-xl shadow hover:shadow-lg transition overflow-hidden group"
  >
    <div className="h-40 flex items-center justify-center overflow-hidden">
      {car.photo? (
        <img
          src={car.photo}
          alt={`${car.brand} ${car.model}`}
          className="w-full h-full object-cover group-hover:scale-105 transition-transform"
          onError={(e) => {
            (e.target as HTMLImageElement).style.display = "none";
            (e.target as HTMLImageElement).nextElementSibling?.classList.remove("hidden");
          }}
        />
      ) : null}
      <div className={`bg-gradient-to-br from-blue-50 to-blue-100 w-full h-full flex items-center justify-center ${car.photo? "hidden" : ""}`}>
        <span className="text-6xl opacity-60 group-hover:scale-110 transition-transform">
          {"\uD83D\uDE97"}
        </span>
      </div>
    </div>
    <div className="p-5">
      <div className="flex justify-between items-start mb-2">
        <h3 className="text-lg font-semibold text-gray-800">
          {car.brand} {car.model}
        </h3>
        <span
          className={`text-xs px-2 py-1 rounded-full font-medium ${
            car.status === "available"
              ? "bg-green-100 text-green-700"
              : "bg-red-100 text-red-700"
          }`}
        >
          {car.status}
        </span>
      </div>
      <p className="text-sm text-gray-500 mb-3">
        {car.year} &middot; {car.plate_number}
      </p>
      <p className="text-xl font-bold text-blue-600">
        ${car.price_per_day.toFixed(2)}
        <span className="text-sm font-normal text-gray-400"> /day</span>
      </p>
    </div>
  </Link>
);

export default CarList;
