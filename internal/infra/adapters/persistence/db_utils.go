package persistence

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/renatocosta55sp/modeling/domain"
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

func ExecuteResultSetBatch(results pgx.BatchResults, newEvents []domain.Event, streamID string, nextVersion int) error {
	for range newEvents {
		_, err := results.Exec()
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
				return fmt.Errorf("aggregateIdentifier %s and sequence number %s can not be duplicated for the same stream", streamID, nextVersion)
			}

			return fmt.Errorf("unable to insert row: %w", err)
		}
	}
	return nil
}
