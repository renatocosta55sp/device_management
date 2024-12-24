package persistence

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/renatocosta55sp/device_management/internal/domain"
	"github.com/renatocosta55sp/device_management/internal/events"
	"github.com/sirupsen/logrus"
)

const DeviceTableName = "devices"

type RepoDevice struct {
	Conn     *pgxpool.Pool
	DBSchema string
}

func (r *RepoDevice) Add(entity *events.DeviceAdded, ctx context.Context) (int, error) {

	query := fmt.Sprintf(`INSERT INTO %s.%s (aggregate_identifier, name, brand) VALUES (@aggregateIdentifier, @name, @brand) RETURNING device_id`, r.DBSchema, DeviceTableName)

	args := pgx.NamedArgs{
		"aggregateIdentifier": entity.AggregateId,
		"name":                entity.Name,
		"brand":               entity.Brand,
	}

	logrus.Info("args", args)

	// Use context with timeout for query
	queryCtx, queryCancel := context.WithTimeout(ctx, 5*time.Second)
	defer queryCancel()

	var device_id int
	err := r.Conn.QueryRow(queryCtx, query, args).Scan(&device_id)

	if err != nil {
		return 0, fmt.Errorf("unable to insert row: %w", err)
	}

	return device_id, nil
}

func (r *RepoDevice) Update(entity *events.DeviceUpdated, ctx context.Context) error {

	command, err := r.Conn.Exec(ctx, "update devices set name=$1, brand=$2, updated_at=$3 where aggregate_identifier=$4", entity.Name, entity.Brand, time.Now(), entity.AggregateId)
	if err != nil {
		return err
	}

	if command.RowsAffected() != 1 {
		return errors.New("no row affected to update")
	}

	return nil
}

func (r *RepoDevice) Remove(entity *events.DeviceRemoved, ctx context.Context) error {

	command, err := r.Conn.Exec(ctx, "update devices set deleted_at=$1 where aggregate_identifier=$2", time.Now(), entity.AggregateId)
	if err != nil {
		return err
	}

	if command.RowsAffected() != 1 {
		return errors.New("no row affected to delete")
	}

	return nil
}

func (r *RepoDevice) GetById(id string, ctx context.Context) (*domain.DeviceAggregate, error) {
	var d domain.DeviceAggregate
	err := r.Conn.QueryRow(ctx, "select aggregate_identifier, name, brand from devices where aggregate_identifier=$1 and deleted_at is null", id).Scan(&d.AggregateID, &d.Name, &d.Brand)

	if err != nil {
		return nil, err
	}

	return &d, nil
}

func (r *RepoDevice) GetByBrand(brand string, ctx context.Context) (*domain.DeviceAggregate, error) {
	var d domain.DeviceAggregate
	err := r.Conn.QueryRow(ctx, "select aggregate_identifier, name, brand from devices where brand ILIKE $1 and deleted_at is null", brand).Scan(&d.Brand, &d.Name, &d.Brand)

	if err != nil {
		return nil, err
	}

	return &d, nil
}

func (r *RepoDevice) GetAll(ctx context.Context) (*[]domain.DeviceAggregate, error) {

	rows, err := r.Conn.Query(ctx, "SELECT aggregate_identifier, name, brand from devices where deleted_at is null")
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var devices []domain.DeviceAggregate
	for rows.Next() {
		var d domain.DeviceAggregate
		if err := rows.Scan(&d.AggregateID, &d.Name, &d.Brand); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		devices = append(devices, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return &devices, nil
}

func NewDeviceRepository(conn *pgxpool.Pool, dbSchema string) *RepoDevice {
	return &RepoDevice{Conn: conn, DBSchema: dbSchema}
}
