// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: start.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createProcessedReceipt = `-- name: CreateProcessedReceipt :one
INSERT INTO receipts_processed (receipt_id, user_id, data)
VALUES ($1, $2, $3)
RETURNING processed_id, receipt_id, user_id, data
`

type CreateProcessedReceiptParams struct {
	ReceiptID *int32 `json:"receipt_id"`
	UserID    *int32 `json:"user_id"`
	Data      []byte `json:"data"`
}

func (q *Queries) CreateProcessedReceipt(ctx context.Context, arg CreateProcessedReceiptParams) (ReceiptsProcessed, error) {
	row := q.db.QueryRow(ctx, createProcessedReceipt, arg.ReceiptID, arg.UserID, arg.Data)
	var i ReceiptsProcessed
	err := row.Scan(
		&i.ProcessedID,
		&i.ReceiptID,
		&i.UserID,
		&i.Data,
	)
	return i, err
}

const createRawReceipt = `-- name: CreateRawReceipt :one
INSERT INTO receipts_raw (user_id, ocr_text)
VALUES ($1, $2)
RETURNING receipt_id, user_id, ocr_text
`

type CreateRawReceiptParams struct {
	UserID  *int32 `json:"user_id"`
	OcrText string `json:"ocr_text"`
}

func (q *Queries) CreateRawReceipt(ctx context.Context, arg CreateRawReceiptParams) (ReceiptsRaw, error) {
	row := q.db.QueryRow(ctx, createRawReceipt, arg.UserID, arg.OcrText)
	var i ReceiptsRaw
	err := row.Scan(&i.ReceiptID, &i.UserID, &i.OcrText)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (first_name, last_name, email, password)
VALUES ($1, $2, $3, $4)
RETURNING user_id, first_name, last_name, email, date_joined
`

type CreateUserParams struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type CreateUserRow struct {
	UserID     int32              `json:"user_id"`
	FirstName  string             `json:"first_name"`
	LastName   string             `json:"last_name"`
	Email      string             `json:"email"`
	DateJoined pgtype.Timestamptz `json:"date_joined"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.Password,
	)
	var i CreateUserRow
	err := row.Scan(
		&i.UserID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.DateJoined,
	)
	return i, err
}

const getProcessedReceiptByUserId = `-- name: GetProcessedReceiptByUserId :many
SELECT processed_id, receipt_id, user_id, data
FROM receipts_processed
WHERE user_id = $1
`

func (q *Queries) GetProcessedReceiptByUserId(ctx context.Context, userID *int32) ([]ReceiptsProcessed, error) {
	rows, err := q.db.Query(ctx, getProcessedReceiptByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ReceiptsProcessed{}
	for rows.Next() {
		var i ReceiptsProcessed
		if err := rows.Scan(
			&i.ProcessedID,
			&i.ReceiptID,
			&i.UserID,
			&i.Data,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRawReceiptByUserId = `-- name: GetRawReceiptByUserId :many
SELECT receipt_id, user_id, ocr_text
FROM receipts_raw
WHERE user_id = $1
`

func (q *Queries) GetRawReceiptByUserId(ctx context.Context, userID *int32) ([]ReceiptsRaw, error) {
	rows, err := q.db.Query(ctx, getRawReceiptByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ReceiptsRaw{}
	for rows.Next() {
		var i ReceiptsRaw
		if err := rows.Scan(&i.ReceiptID, &i.UserID, &i.OcrText); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT user_id, first_name, last_name, email, date_joined
FROM users
WHERE email = $1
`

type GetUserByEmailRow struct {
	UserID     int32              `json:"user_id"`
	FirstName  string             `json:"first_name"`
	LastName   string             `json:"last_name"`
	Email      string             `json:"email"`
	DateJoined pgtype.Timestamptz `json:"date_joined"`
}

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (GetUserByEmailRow, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i GetUserByEmailRow
	err := row.Scan(
		&i.UserID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.DateJoined,
	)
	return i, err
}