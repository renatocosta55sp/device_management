package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/renatocosta55sp/device_management/internal/infra"
	"github.com/sirupsen/logrus"
)

var (
	db *pgxpool.Pool
)

func configureDB(ctx context.Context, host, port, user, password, database string) {

	connectionString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, database)
	db, err = infra.NewPG(ctx, connectionString)

	handleError(err, "failed to create Postgresql connection")

	err = db.Ping(ctx)
	handleError(err, "failed to ping Postgresql connection")

	logrus.Info("Using Postgresql")

}
