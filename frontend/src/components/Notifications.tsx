import React from "react";
import { useNotification } from "../context/NotificationContext";

const bgMap = {
  success: "bg-green-500",
  error: "bg-red-500",
  info: "bg-blue-500",
};

const Notifications: React.FC = () => {
  const { notifications, remove } = useNotification();

  if (notifications.length === 0) return null;

  return (
    <div className="fixed top-4 right-4 z-50 space-y-2 max-w-sm">
      {notifications.map((n) => (
        <div
          key={n.id}
          className={`${bgMap[n.type]} text-white px-4 py-3 rounded-lg shadow-lg flex items-center justify-between`}
        >
          <span className="text-sm">{n.message}</span>
          <button
            onClick={() => remove(n.id)}
            className="ml-3 text-white/80 hover:text-white"
          >
            &times;
          </button>
        </div>
      ))}
    </div>
  );
};

export default Notifications;
