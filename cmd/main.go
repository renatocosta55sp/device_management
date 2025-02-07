package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/renatocosta55sp/device_management/internal/infra"
	"github.com/sirupsen/logrus"
)

var (
	err error
)

func handleError(err error, msg string) {
	if err != nil {
		logrus.WithError(err).Fatal(msg)
	}
}

func main() {

	ctx := context.Background()

	ctx, stopFn := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stopFn()

	logrus.Info(ctx, infra.Config.GetString("APP_NAME")+" - v.: "+infra.Config.GetString("APP_VERSION"))

	configureDB(ctx,
		infra.Config.GetString("DB_HOST"),
		infra.Config.GetString("DB_PORT"),
		infra.Config.GetString("DB_USER"),
		infra.Config.GetString("DB_PASS"),
		infra.Config.GetString("DB_NAME"),
	)

	defer db.Close()

	domainEventRegistry := configureDomainEventRegistry()

	server := gin.Default()
	infra.InitRoutes(&server.RouterGroup, db, domainEventRegistry)

	serverErr := make(chan error, 1)
	go func() {
		if err := server.Run(":8080"); err != nil {
			serverErr <- err
		}
	}()

	//Graceful shutdown
	select {
	case err = <-serverErr:
		handleError(err, "error on server execution")
	case <-ctx.Done():
		handleError(ctx.Err(), "ctx error")
		handleError(err, "failed to stop server gracefully")
	}

}
