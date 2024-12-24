package infra

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pgInstance *pgxpool.Pool
	pgOnce     sync.Once
)

func NewPG(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	var err error

	pgOnce.Do(func() {
		var db *pgxpool.Pool
		db, err = pgxpool.New(ctx, connString)
		if err == nil {
			pgInstance = db
		}
	})

	return pgInstance, err
}
