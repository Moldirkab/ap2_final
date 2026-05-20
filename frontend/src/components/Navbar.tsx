import React from "react";
import { Link, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

const Navbar: React.FC = () => {
  const { token, user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate("/login");
  };

  return (
    <nav className="bg-blue-700 text-white shadow-lg">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          <Link to="/" className="text-xl font-bold tracking-wide">
            CarRental
          </Link>
          <div className="hidden md:flex items-center space-x-6">
            <Link to="/" className="hover:text-blue-200 transition">
              Home
            </Link>
            <Link to="/cars" className="hover:text-blue-200 transition">
              Cars
            </Link>
            {token ? (
              <>
                <Link to="/bookings" className="hover:text-blue-200 transition">
                  My Bookings
                </Link>
                <Link to="/profile" className="hover:text-blue-200 transition">
                  Profile
                </Link>
                {user?.role === "admin" && (
                  <Link
                    to="/admin/cars/new"
                    className="bg-green-600 hover:bg-green-500 px-3 py-1.5 rounded-lg text-sm font-semibold transition"
                  >
                    + Add Car
                  </Link>
                )}
                <button
                  onClick={handleLogout}
                  className="bg-blue-600 hover:bg-blue-500 px-4 py-2 rounded-lg text-sm transition"
                >
                  Logout
                </button>
              </>
            ) : (
              <>
                <Link
                  to="/login"
                  className="hover:text-blue-200 transition"
                >
                  Login
                </Link>
                <Link
                  to="/register"
                  className="bg-white text-blue-700 px-4 py-2 rounded-lg text-sm font-semibold hover:bg-blue-50 transition"
                >
                  Sign Up
                </Link>
              </>
            )}
          </div>
          {/* Mobile menu button */}
          <MobileMenu token={token} user={user} onLogout={handleLogout} />
        </div>
      </div>
    </nav>
  );
};

const MobileMenu: React.FC<{
  token: string | null;
  user: { role: string } | null;
  onLogout: () => void;
}> = ({ token, user, onLogout }) => {
  const [open, setOpen] = React.useState(false);
  return (
    <div className="md:hidden">
      <button
        onClick={() => setOpen(!open)}
        className="text-white focus:outline-none"
        aria-label="Toggle menu"
      >
        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          {open ? (
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
          ) : (
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
          )}
        </svg>
      </button>
      {open && (
        <div className="absolute top-16 left-0 w-full bg-blue-700 z-50 px-4 pb-4 space-y-2">
          <Link to="/" className="block py-2 hover:text-blue-200" onClick={() => setOpen(false)}>Home</Link>
          <Link to="/cars" className="block py-2 hover:text-blue-200" onClick={() => setOpen(false)}>Cars</Link>
          {token ? (
            <>
              <Link to="/bookings" className="block py-2 hover:text-blue-200" onClick={() => setOpen(false)}>My Bookings</Link>
              <Link to="/profile" className="block py-2 hover:text-blue-200" onClick={() => setOpen(false)}>Profile</Link>
              {user?.role === "admin" && (
                <Link to="/admin/cars/new" className="block py-2 text-green-300 hover:text-green-200 font-semibold" onClick={() => setOpen(false)}>+ Add Car</Link>
              )}
              <button onClick={() => { onLogout(); setOpen(false); }} className="block w-full text-left py-2 hover:text-blue-200">Logout</button>
            </>
          ) : (
            <>
              <Link to="/login" className="block py-2 hover:text-blue-200" onClick={() => setOpen(false)}>Login</Link>
              <Link to="/register" className="block py-2 hover:text-blue-200" onClick={() => setOpen(false)}>Sign Up</Link>
            </>
          )}
        </div>
      )}
    </div>
  );
};

export default Navbar;
