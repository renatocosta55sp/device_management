package http

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/renatocosta55sp/device_management/internal/domain"
	"github.com/renatocosta55sp/device_management/internal/domain/commands"
	"github.com/renatocosta55sp/device_management/internal/infra/adapters/persistence"
	"github.com/renatocosta55sp/device_management/internal/slices/adddevice"
	"github.com/renatocosta55sp/modeling/infra/bus"
	"github.com/renatocosta55sp/modeling/slice"
	"github.com/sirupsen/logrus"
)

type HttpServer struct {
	Db *pgxpool.Pool
}

func (h HttpServer) AddDevice(ctx *gin.Context) {

	requestData, exists := ctx.Get(requestDataKey)
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Request data not found"})
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
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	eventBus := bus.NewEventBus()
	_, ctxCancFunc := context.WithTimeout(context.Background(), 5*time.Second)

	eventResultChan := adddevice.WireApp(ctx,
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
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"result": commandResult})

}

func (h HttpServer) UpdateDevice(ctx *gin.Context) {

	requestData, exists := ctx.Get(requestDataKey)
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Request data not found"})
		return
	}

	data := requestData.(UpdateDeviceRequest)
	aggregateIdentifier, err := uuid.Parse(data.Id)

	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
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
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	eventBus := bus.NewEventBus()
	_, ctxCancFunc := context.WithTimeout(context.Background(), 5*time.Second)

	eventResultChan := adddevice.WireApp(ctx,
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
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"result": commandResult})

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

	command := commands.RemoveDeviceCommand{
		AggregateID: aggregateIdentifier,
	}

	device := domain.NewDevice(aggregateIdentifier)

	commandResult, err := device.HandleDelete(command)
	if err != nil {
		logrus.WithError(err).Error("failed to validate device on command remove")
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	eventBus := bus.NewEventBus()
	_, ctxCancFunc := context.WithTimeout(context.Background(), 5*time.Second)

	eventResultChan := adddevice.WireApp(ctx,
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
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"result": commandResult})

}

func (h HttpServer) GetDeviceById(ctx *gin.Context) {

	id := ctx.Param("id")

	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		ctx.Abort()
		return
	}

	repo := persistence.NewDeviceRepository(h.Db, "public")
	devices, err := repo.GetById(id, ctx)

	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": devices})

}

func (h HttpServer) GetDeviceByBrand(ctx *gin.Context) {

	brand := ctx.Param("brand")

	if brand == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Brand parameter is required"})
		ctx.Abort()
		return
	}

	repo := persistence.NewDeviceRepository(h.Db, "public")
	devices, err := repo.GetByBrand(brand, ctx)

	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": devices})

}

func (h HttpServer) GetDevices(ctx *gin.Context) {

	repo := persistence.NewDeviceRepository(h.Db, "public")
	devices, err := repo.GetAll(ctx)

	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": devices})

}
