import React, { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { createCar, getCar, updateCar } from "../api/cars";
import { useAuth } from "../context/AuthContext";
import { useNotification } from "../context/NotificationContext";
import LoadingSpinner from "../components/LoadingSpinner";

const AdminCreateCar: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const isEdit = Boolean(id);

  const { user, loading: authLoading } = useAuth();
  const { notify } = useNotification();
  const navigate = useNavigate();

  const [brand, setBrand] = useState("");
  const [model, setModel] = useState("");
  const [year, setYear] = useState("");
  const [plateNumber, setPlateNumber] = useState("");
  const [pricePerDay, setPricePerDay] = useState("");

  const [photo, setPhoto] = useState("");
  const [imagePreviewError, setImagePreviewError] = useState(false);

  const [submitting, setSubmitting] = useState(false);
  const [loadingCar, setLoadingCar] = useState(false);

  useEffect(() => {
    if (isEdit && id) {
      setLoadingCar(true);
      getCar(parseInt(id))
          .then((res) => {
            const c = res.data;
            setBrand(c.brand);
            setModel(c.model);
            setYear(c.year.toString());
            setPlateNumber(c.plate_number);
            setPricePerDay(c.price_per_day.toString());
            setPhoto(c.photo || "");
          })
          .catch(() => notify("Failed to load car for editing", "error"))
          .finally(() => setLoadingCar(false));
    }
  }, [id, isEdit, notify]);

  if (authLoading || loadingCar) return <LoadingSpinner />;

  if (!user || user.role !== "admin") {
    return (
        <div className="max-w-2xl mx-auto px-4 py-16 text-center">
          <div className="bg-white rounded-xl shadow-lg p-8">
            <h2 className="text-2xl font-bold text-red-600 mb-4">
              Access Denied
            </h2>
            <p className="text-gray-600 mb-6">
              This page is only accessible to administrators.
            </p>
            <button
                onClick={() => navigate("/")}
                className="bg-blue-600 text-white px-6 py-2 rounded-lg font-semibold hover:bg-blue-700 transition"
            >
              Go Home
            </button>
          </div>
        </div>
    );
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!brand || !model || !year || !plateNumber || !pricePerDay) {
      notify("Please fill in all required fields", "error");
      return;
    }

    setSubmitting(true);

    try {
      const payload = {
        brand,
        model,
        year: parseInt(year),
        plate_number: plateNumber,
        price_per_day: parseFloat(pricePerDay),
        photo: photo || undefined,
      };

      if (isEdit && id) {
        await updateCar(parseInt(id), payload);
        notify("Car updated successfully!", "success");
      } else {
        await createCar(payload);
        notify("Car created successfully!", "success");
      }

      navigate("/cars");
    } catch (err: any) {
      notify(
          err.response?.data?.error ||
          (isEdit ? "Failed to update car" : "Failed to create car"),
          "error"
      );
    } finally {
      setSubmitting(false);
    }
  };

  return (
      <div className="max-w-2xl mx-auto px-4 py-8">
        <button
            onClick={() => navigate(isEdit ? `/cars/${id}` : "/cars")}
            className="text-blue-600 hover:underline mb-6 inline-block"
        >
          &larr; {isEdit ? "Back to Car" : "Back to Cars"}
        </button>

        <div className="bg-white rounded-xl shadow-lg p-8">
          <h2 className="text-2xl font-bold text-gray-800 mb-6">
            {isEdit ? "Update Car" : "Add New Car"}
          </h2>

          <form onSubmit={handleSubmit} className="space-y-4">
            <input
                className="w-full border p-2 rounded"
                placeholder="Brand"
                value={brand}
                onChange={(e) => setBrand(e.target.value)}
            />

            <input
                className="w-full border p-2 rounded"
                placeholder="Model"
                value={model}
                onChange={(e) => setModel(e.target.value)}
            />

            <input
                className="w-full border p-2 rounded"
                placeholder="Year"
                type="number"
                value={year}
                onChange={(e) => setYear(e.target.value)}
            />

            <input
                className="w-full border p-2 rounded"
                placeholder="Plate Number"
                value={plateNumber}
                onChange={(e) => setPlateNumber(e.target.value)}
            />

            <input
                className="w-full border p-2 rounded"
                placeholder="Price Per Day"
                type="number"
                step="0.01"
                value={pricePerDay}
                onChange={(e) => setPricePerDay(e.target.value)}
            />

            {/* IMAGE URL */}
            <input
                className="w-full border p-2 rounded"
                placeholder="Image URL"
                value={photo}
                onChange={(e) => {
                  setPhoto(e.target.value);
                  setImagePreviewError(false);
                }}
            />

            {photo && (
                <div className="mt-3 border rounded overflow-hidden">
                  {!imagePreviewError ? (
                      <img
                          src={photo}
                          className="w-full h-48 object-cover"
                          onError={() => setImagePreviewError(true)}
                      />
                  ) : (
                      <div className="h-48 flex items-center justify-center text-red-500">
                        Invalid image URL
                      </div>
                  )}
                </div>
            )}

            <button
                disabled={submitting}
                className="w-full bg-blue-600 text-white p-3 rounded"
            >
              {submitting
                  ? isEdit
                      ? "Updating..."
                      : "Creating..."
                  : isEdit
                      ? "Update Car"
                      : "Create Car"}
            </button>
          </form>
        </div>
      </div>
  );
};

export default AdminCreateCar;