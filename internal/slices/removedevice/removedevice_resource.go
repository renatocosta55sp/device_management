package removedevice

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

type RemoveDeviceRequest struct {
	Id string
}

func RemoveDeviceRequestValidator(ctx *gin.Context) {

	id := ctx.Param("id")

	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"errorD": "ID parameter is required"})
		ctx.Abort()
		return
	}

	var requestData RemoveDeviceRequest

	requestData.Id = id
	ctx.Set("requestData", requestData)

	ctx.Next()

}

func (h HttpServer) RemoveDevice(ctx *gin.Context) {

	requestData, exists := ctx.Get(requestDataKey)
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Request data not found"})
		return
	}

	data := requestData.(RemoveDeviceRequest)
	aggregateIdentifier, err := uuid.Parse(data.Id)

	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	commandResult, _, err := CommandExecutor{
		persistence.NewPersistentEventStore(h.Db, h.DomainEventRegistry, "public"),
	}.Send(
		ctx,
		commands.RemoveDeviceCommand{
			AggregateID: aggregateIdentifier,
		},
	)

	if err != nil {
		logrus.WithError(err).Error("failed to validate device on command creation")
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"result": commandResult})

}
