package updatedevice

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

type UpdateDeviceRequest struct {
	Id    string
	Name  string `json:"name"  binding:"required"`
	Brand string `json:"brand"  binding:"required"`
}

func UpdateDeviceValidator(ctx *gin.Context) {

	id := ctx.Param("id")

	if id == "" {
		ctx.JSON(http_.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		ctx.Abort()
		return
	}

	var requestData UpdateDeviceRequest

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http_.StatusBadRequest, gin.H{"errorUpd": err.Error()})
		ctx.Abort()
		return
	}

	requestData.Id = id
	ctx.Set("requestData", requestData)

	ctx.Next()

}

type PatchDeviceRequest struct {
	Repo *pgxpool.Pool
}

func (p *PatchDeviceRequest) UpdatePartiallyDeviceValidator(ctx *gin.Context) {

	id := ctx.Param("id")

	if id == "" {
		ctx.JSON(http_.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		ctx.Abort()
		return
	}

	var patchData map[string]interface{}
	var requestData UpdateDeviceRequest

	if err := ctx.ShouldBindJSON(&patchData); err != nil {
		ctx.JSON(http_.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	repo := persistence.NewDeviceRepository(p.Repo, "public")
	device, err := repo.GetById(id, ctx)

	if err != nil {
		ctx.JSON(http_.StatusNotFound, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	if name, exists := patchData["name"]; exists {
		requestData.Name = name.(string)
	} else {
		requestData.Name = device.Name
	}

	if brand, exists := patchData["brand"]; exists {
		requestData.Brand = brand.(string)
	} else {
		requestData.Brand = device.Brand
	}

	requestData.Id = id
	ctx.Set("requestData", requestData)

	ctx.Next()

}

func (h HttpServer) UpdateDevice(ctx *gin.Context) {

	requestData, exists := ctx.Get(requestDataKey)
	if !exists {
		ctx.JSON(http_.StatusInternalServerError, gin.H{"error": "Request data not found"})
		return
	}

	data := requestData.(UpdateDeviceRequest)
	aggregateIdentifier, err := uuid.Parse(data.Id)

	if err != nil {
		ctx.JSON(http_.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	command := commands.UpdateDeviceCommand{
		AggregateID: aggregateIdentifier,
		Name:        data.Name,
		Brand:       data.Brand,
	}

	device := domain.NewDevice(aggregateIdentifier)

	commandResult, err := device.HandleUpdate(command)
	if err != nil {
		logrus.WithError(err).Error("failed to validate device on command update")
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
		logrus.WithError(err).Error("failed to validate device on command update")
		ctx.JSON(http_.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http_.StatusNoContent, gin.H{"result": commandResult})

}
