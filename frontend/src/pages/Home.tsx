import React from "react";
import { Link } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

const Home: React.FC = () => {
  const { token } = useAuth();

  return (
    <div className="min-h-[80vh]">
      {/* Hero */}
      <section className="bg-gradient-to-br from-blue-700 to-blue-900 text-white">
        <div className="max-w-7xl mx-auto px-4 py-24 sm:py-32 text-center">
          <h1 className="text-4xl sm:text-5xl font-extrabold mb-6">
            Find Your Perfect Rental Car
          </h1>
          <p className="text-lg sm:text-xl text-blue-100 mb-8 max-w-2xl mx-auto">
            Browse our fleet, book in seconds, and hit the road. Affordable
            prices, top-quality vehicles.
          </p>
          <div className="flex justify-center gap-4 flex-wrap">
            <Link
              to="/cars"
              className="bg-white text-blue-700 px-8 py-3 rounded-lg font-semibold hover:bg-blue-50 transition shadow-lg"
            >
              Browse Cars
            </Link>
            {!token && (
              <Link
                to="/register"
                className="border-2 border-white text-white px-8 py-3 rounded-lg font-semibold hover:bg-white/10 transition"
              >
                Create Account
              </Link>
            )}
          </div>
        </div>
      </section>

      {/* Features */}
      <section className="max-w-7xl mx-auto px-4 py-16">
        <h2 className="text-3xl font-bold text-center text-gray-800 mb-12">
          Why Choose CarRental?
        </h2>
        <div className="grid md:grid-cols-3 gap-8">
          {features.map((f, i) => (
            <div
              key={i}
              className="bg-white rounded-xl shadow-md p-8 text-center hover:shadow-lg transition"
            >
              <div className="text-4xl mb-4">{f.icon}</div>
              <h3 className="text-xl font-semibold text-gray-800 mb-2">
                {f.title}
              </h3>
              <p className="text-gray-600">{f.desc}</p>
            </div>
          ))}
        </div>
      </section>

      {/* CTA */}
      <section className="bg-gray-100 py-16">
        <div className="max-w-3xl mx-auto text-center px-4">
          <h2 className="text-2xl font-bold text-gray-800 mb-4">
            Ready to get started?
          </h2>
          <p className="text-gray-600 mb-6">
            Sign up today and book your first ride in minutes.
          </p>
          <Link
            to={token ? "/cars" : "/register"}
            className="bg-blue-600 text-white px-8 py-3 rounded-lg font-semibold hover:bg-blue-700 transition"
          >
            {token ? "Browse Cars" : "Get Started"}
          </Link>
        </div>
      </section>
    </div>
  );
};

const features = [
  {
    icon: "\uD83D\uDE97",
    title: "Wide Selection",
    desc: "Choose from sedans, SUVs, trucks, and more at competitive prices.",
  },
  {
    icon: "\u26A1",
    title: "Instant Booking",
    desc: "Book a car in seconds with our streamlined reservation system.",
  },
  {
    icon: "\uD83D\uDD12",
    title: "Secure & Reliable",
    desc: "JWT-secured accounts and transparent booking history.",
  },
];

export default Home;
