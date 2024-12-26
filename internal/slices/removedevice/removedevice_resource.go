package removedevice

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/renatocosta55sp/device_management/internal/domain"
	"github.com/renatocosta55sp/device_management/internal/domain/commands"
	"github.com/renatocosta55sp/device_management/internal/infra/adapters/persistence"
	"github.com/renatocosta55sp/modeling/infra/bus"
	"github.com/renatocosta55sp/modeling/slice"
	"github.com/sirupsen/logrus"

	http_ "net/http"
)

type HttpServer struct {
	Db *pgxpool.Pool
}

const requestDataKey = "requestData"

type RemoveDeviceRequest struct {
	Id string
}

func RemoveDeviceValidator(ctx *gin.Context) {

	id := ctx.Param("id")

	if id == "" {
		ctx.JSON(http_.StatusBadRequest, gin.H{"errorD": "ID parameter is required"})
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
		ctx.JSON(http_.StatusInternalServerError, gin.H{"error": "Request data not found"})
		return
	}

	data := requestData.(RemoveDeviceRequest)
	aggregateIdentifier, err := uuid.Parse(data.Id)

	if err != nil {
		ctx.JSON(http_.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	command := commands.RemoveDeviceCommand{
		AggregateID: aggregateIdentifier,
	}

	device := domain.NewDevice(aggregateIdentifier)

	commandResult, err := device.HandleDelete(command)
	if err != nil {
		logrus.WithError(err).Error("failed to validate device on command remove")
		ctx.JSON(http_.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	eventBus := bus.NewEventBus()
	_, ctxCancFunc := context.WithTimeout(context.Background(), 5*time.Second)

	eventResultChan := WireApp(ctx,
		eventBus,
		*persistence.NewDeviceRepository(h.Db, "public"),
	)

	err = (&slice.CommandExecutionResult{
		EventBus:        eventBus,
		CtxCancFunc:     ctxCancFunc,
		EventResultChan: eventResultChan,
	}).Execute(device.Events)

	if err != nil {
		logrus.WithError(err).Error("failed to validate device on command remove")
		ctx.JSON(http_.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http_.StatusNoContent, gin.H{"result": commandResult})

}
