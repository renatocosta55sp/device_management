package adddevice

import (
	"context"
	"time"

	http_ "net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/renatocosta55sp/device_management/internal/domain"
	"github.com/renatocosta55sp/device_management/internal/domain/commands"
	"github.com/renatocosta55sp/device_management/internal/infra/adapters/persistence"

	"github.com/renatocosta55sp/modeling/infra/bus"
	"github.com/renatocosta55sp/modeling/slice"
	"github.com/sirupsen/logrus"
)

type HttpServer struct {
	Db *pgxpool.Pool
}

const requestDataKey = "requestData"

type AddDeviceRequest struct {
	Name  string `json:"name"  binding:"required"`
	Brand string `json:"brand"  binding:"required"`
}

func AddDeviceValidator(ctx *gin.Context) {

	var requestData AddDeviceRequest

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http_.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	ctx.Set("requestData", requestData)

	ctx.Next()

}

func (h HttpServer) AddDevice(ctx *gin.Context) {

	requestData, exists := ctx.Get(requestDataKey)
	if !exists {
		ctx.JSON(http_.StatusInternalServerError, gin.H{"error": "Request data not found"})
		return
	}

	data := requestData.(AddDeviceRequest)

	aggregateIdentifier := uuid.New()
	command := commands.AddDeviceCommand{
		AggregateID: aggregateIdentifier,
		Name:        data.Name,
		Brand:       data.Brand,
	}

	device := domain.NewDevice(aggregateIdentifier)

	commandResult, err := device.HandleAdd(command)
	if err != nil {
		logrus.WithError(err).Error("failed to validate device on command creation")
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
		logrus.WithError(err).Error("failed to validate device on command creation")
		ctx.JSON(http_.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http_.StatusCreated, gin.H{"result": commandResult})

}
