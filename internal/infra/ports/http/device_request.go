package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/renatocosta55sp/device_management/internal/infra/adapters/persistence"
)

const requestDataKey = "requestData"

type AddDeviceRequest struct {
	Name  string `json:"name"  binding:"required"`
	Brand string `json:"brand"  binding:"required"`
}

func addDeviceValidator(ctx *gin.Context) {

	var requestData AddDeviceRequest

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	ctx.Set("requestData", requestData)

	ctx.Next()

}

type UpdateDeviceRequest struct {
	Id    string
	Name  string `json:"name"  binding:"required"`
	Brand string `json:"brand"  binding:"required"`
}

func updateDeviceValidator(ctx *gin.Context) {

	id := ctx.Param("id")

	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		ctx.Abort()
		return
	}

	var requestData UpdateDeviceRequest

	if err := ctx.ShouldBindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"errorUpd": err.Error()})
		ctx.Abort()
		return
	}

	requestData.Id = id
	ctx.Set("requestData", requestData)

	ctx.Next()

}

type PatchDeviceRequest struct {
	repo *pgxpool.Pool
}

func (p *PatchDeviceRequest) updatePartiallyDeviceValidator(ctx *gin.Context) {

	id := ctx.Param("id")

	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		ctx.Abort()
		return
	}

	var patchData map[string]interface{}
	var requestData UpdateDeviceRequest

	if err := ctx.ShouldBindJSON(&patchData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	repo := persistence.NewDeviceRepository(p.repo, "public")
	device, err := repo.GetById(id, ctx)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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

type RemoveDeviceRequest struct {
	Id string
}

func removeDeviceValidator(ctx *gin.Context) {

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
