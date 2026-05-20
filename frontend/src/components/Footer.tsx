import React from "react";

const Footer: React.FC = () => (
  <footer className="bg-gray-800 text-gray-400 py-8 mt-auto">
    <div className="max-w-7xl mx-auto px-4 text-center">
      <p className="text-sm">
        &copy; {new Date().getFullYear()} CarRental. All rights reserved.
      </p>
      <p className="text-xs mt-1">Car Rental Microservices Platform</p>
    </div>
  </footer>
);

export default Footer;
