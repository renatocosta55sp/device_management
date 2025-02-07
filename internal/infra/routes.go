package infra

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/renatocosta55sp/device_management/internal/slices/adddevice"
	"github.com/renatocosta55sp/device_management/internal/slices/removedevice"
	"github.com/renatocosta55sp/modeling/infra/bus"
)

func InitRoutes(
	r *gin.RouterGroup,
	db *pgxpool.Pool,
	domainEventRegistry *bus.EventRegistry) {

	res := adddevice.HttpServer{Db: db, DomainEventRegistry: domainEventRegistry}
	r.POST("/devices", adddevice.ValidateRequest, res.AddDevice)

	resRemoveDevice := removedevice.HttpServer{Db: db, DomainEventRegistry: domainEventRegistry}
	r.DELETE("/devices/:id", removedevice.RemoveDeviceRequestValidator, resRemoveDevice.RemoveDevice)
}
