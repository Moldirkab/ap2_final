CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS bookings (
                                        id UUID PRIMARY KEY,
                                        user_id UUID NOT NULL,
                                        car_id UUID NOT NULL,
                                        start_date DATE NOT NULL,
                                        end_date DATE NOT NULL,
                                        total_price NUMERIC(10,2) NOT NULL,
                                        status VARCHAR(30) NOT NULL,
                                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE bookings ALTER COLUMN user_id TYPE VARCHAR(255);

ALTER TABLE bookings ALTER COLUMN car_id TYPE VARCHAR(255);