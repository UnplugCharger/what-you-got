-- name: CreateUser :one
INSERT INTO users (first_name, last_name, email, password)
VALUES ($1, $2, $3, $4)
RETURNING user_id, first_name, last_name, email, date_joined;

-- name: GetUserByEmail :one
SELECT user_id, first_name, last_name, email, date_joined
FROM users
WHERE email = $1;


-- name: CreateRawReceipt :one
INSERT INTO receipts_raw (user_id, ocr_text)
VALUES ($1, $2)
RETURNING receipt_id, user_id, ocr_text;

-- name: GetRawReceiptByUserId :many
SELECT receipt_id, user_id, ocr_text
FROM receipts_raw
WHERE user_id = $1;

-- name: CreateProcessedReceipt :one
INSERT INTO receipts_processed (receipt_id, user_id, data)
VALUES ($1, $2, $3)
RETURNING processed_id, receipt_id, user_id, data;

-- name: GetProcessedReceiptByUserId :many
SELECT processed_id, receipt_id, user_id, data
FROM receipts_processed
WHERE user_id = $1;
