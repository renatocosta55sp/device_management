package testsuite

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func InitTestContainer() (*pgxpool.Pool, testcontainers.Container, error) {

	ctx := context.Background()

	dbName := "device_management"
	dbUser := "postgres"
	dbPassword := "123456"

	sqlFiles, err := getSQLFilesFromDir("../../../database/migrations")
	if err != nil {
		log.Fatalf("failed to get SQL files: %v", err)
	}

	pgContainer, err := postgres.Run(ctx,
		"docker.io/postgres:16-alpine",
		postgres.WithInitScripts(sqlFiles...),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	host, err := pgContainer.Host(ctx)
	if err != nil {
		return nil, pgContainer, err
	}

	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, pgContainer, err
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, host, port.Port(), dbName)

	conn, err := pgxpool.New(ctx, connStr)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, pgContainer, err
	}

	return conn, pgContainer, nil
}

func getSQLFilesFromDir(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			sqlFiles = append(sqlFiles, filepath.Join(dir, file.Name()))
		}
	}
	return sqlFiles, nil
}
