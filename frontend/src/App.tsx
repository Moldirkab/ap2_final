import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { AuthProvider } from "./context/AuthContext";
import { NotificationProvider } from "./context/NotificationContext";
import Navbar from "./components/Navbar";
import Footer from "./components/Footer";
import Notifications from "./components/Notifications";
import Home from "./pages/Home";
import Login from "./pages/Login";
import Register from "./pages/Register";
import Profile from "./pages/Profile";
import CarList from "./pages/CarList";
import CarDetail from "./pages/CarDetail";
import BookingHistory from "./pages/BookingHistory";
import AdminCreateCar from "./pages/AdminCreateCar";

const App: React.FC = () => (
  <AuthProvider>
    <NotificationProvider>
      <Router>
        <div className="flex flex-col min-h-screen bg-gray-50">
          <Navbar />
          <Notifications />
          <main className="flex-1">
            <Routes>
              <Route path="/" element={<Home />} />
              <Route path="/login" element={<Login />} />
              <Route path="/register" element={<Register />} />
              <Route path="/profile" element={<Profile />} />
              <Route path="/cars" element={<CarList />} />
              <Route path="/cars/:id" element={<CarDetail />} />
              <Route path="/bookings" element={<BookingHistory />} />
              <Route path="/admin/cars/new" element={<AdminCreateCar />} />
              <Route path="/admin/cars/:id/edit" element={<AdminCreateCar />} />
            </Routes>
          </main>
          <Footer />
        </div>
      </Router>
    </NotificationProvider>
  </AuthProvider>
);

export default App;
