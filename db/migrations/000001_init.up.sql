CREATE TABLE users (
user_id SERIAL PRIMARY KEY,
first_name VARCHAR(255) NOT NULL,
last_name VARCHAR(255) NOT NULL,
email VARCHAR(255) UNIQUE NOT NULL,
password VARCHAR(255) NOT NULL,
date_joined TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE receipts_raw (
receipt_id SERIAL PRIMARY KEY,
user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
ocr_text TEXT NOT NULL
);

CREATE TABLE receipts_processed (
processed_id SERIAL PRIMARY KEY,
receipt_id INT REFERENCES receipts_raw(receipt_id) ON DELETE CASCADE,
user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
data JSONB NOT NULL
);
