// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0

package db

import (
	"context"
)

type Querier interface {
	CreateProcessedReceipt(ctx context.Context, arg CreateProcessedReceiptParams) (ReceiptsProcessed, error)
	CreateRawReceipt(ctx context.Context, arg CreateRawReceiptParams) (ReceiptsRaw, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetProcessedReceiptByUserId(ctx context.Context, userID *int32) ([]ReceiptsProcessed, error)
	GetRawReceiptByUserId(ctx context.Context, userID *int32) ([]ReceiptsRaw, error)
	GetUserByEmail(ctx context.Context, email string) (GetUserByEmailRow, error)
}

var _ Querier = (*Queries)(nil)
