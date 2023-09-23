package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Store provides all functions  to run database queries and transactions
type Store interface {
	Querier
}

// SQLStore / SQLStore provides all functions to run database queries and transactions
type SQLStore struct {
	*Queries
	pool *pgxpool.Pool
}

// NewStore creates a new Store with connection pooling
func NewStore(pool *pgxpool.Pool) Store {
	return &SQLStore{
		Queries: New(pool),
		pool:    pool,
	}
}

// ExecTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	conn, err := store.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit(ctx)
}
