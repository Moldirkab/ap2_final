CREATE TABLE IF NOT EXISTS cars (
                                    id SERIAL PRIMARY KEY,
                                    brand TEXT NOT NULL,
                                    model TEXT NOT NULL,
                                    year INT NOT NULL,
                                    plate_number TEXT UNIQUE NOT NULL,
                                    price_per_day DOUBLE PRECISION NOT NULL,
                                    status TEXT NOT NULL DEFAULT 'available',
                                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                    photo TEXT
);



