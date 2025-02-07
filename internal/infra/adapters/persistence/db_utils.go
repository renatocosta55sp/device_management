package persistence

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetNextVal(ctx context.Context, db *pgxpool.Pool, sequenceName string) (int64, error) {
	var nextVal int64
	query := "SELECT nextval($1)"
	err := db.QueryRow(ctx, query, sequenceName).Scan(&nextVal)
	if err != nil {
		return 0, err
	}
	return nextVal, nil
}
