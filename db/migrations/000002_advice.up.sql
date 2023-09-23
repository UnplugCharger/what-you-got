CREATE TABLE receipts_advice (
                                 advice_id SERIAL PRIMARY KEY,
                                 processed_id INT REFERENCES receipts_processed(processed_id) ON DELETE CASCADE,
                                 advice TEXT NOT NULL,
                                 date_generated DATE DEFAULT CURRENT_DATE,
                                 metadata JSONB
);
