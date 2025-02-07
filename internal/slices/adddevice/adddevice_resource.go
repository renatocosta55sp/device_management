package adddevice

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/renatocosta55sp/device_management/internal/domain/commands"
	"github.com/renatocosta55sp/device_management/internal/infra/adapters/persistence"
	"github.com/renatocosta55sp/modeling/infra/bus"
	"github.com/sirupsen/logrus"
)

type HttpServer struct {
	Db                  *pgxpool.Pool
	DomainEventRegistry *bus.EventRegistry
}

const requestDataKey = "requestData"

type DeviceRequest struct {
	Name  string `json:"name"  binding:"required"`
	Brand string `json:"brand"  binding:"required"`
}

func ValidateRequest(ctx *gin.Context) {

	var requestData DeviceRequest

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	ctx.Set("requestData", requestData)

	ctx.Next()

}

func (h HttpServer) AddDevice(ctx *gin.Context) {

	requestData, exists := ctx.Get(requestDataKey)
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Request data not found"})
		return
	}

	data := requestData.(DeviceRequest)

	commandResult, _, err := CommandExecutor{
		persistence.NewPersistentEventStore(h.Db, h.DomainEventRegistry, "public"),
	}.Send(
		ctx,
		commands.AddDeviceCommand{
			AggregateID: uuid.New(),
			Name:        data.Name,
			Brand:       data.Brand,
		},
	)

	if err != nil {
		logrus.WithError(err).Error("failed to validate device on command creation")
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"result": commandResult})

}
