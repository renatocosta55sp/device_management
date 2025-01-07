package infra

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/renatocosta55sp/device_management/internal/slices/adddevice"
	"github.com/renatocosta55sp/device_management/internal/slices/devices"
	"github.com/renatocosta55sp/device_management/internal/slices/removedevice"
	"github.com/renatocosta55sp/device_management/internal/slices/updatedevice"
)

func InitRoutes(
	r *gin.RouterGroup,
	db *pgxpool.Pool) {

	res := adddevice.HttpServer{Db: db}
	r.POST("/devices", adddevice.ValidateRequest, res.AddDevice)

	resUpdateDevice := updatedevice.HttpServer{Db: db}
	r.PUT("/devices/:id", updatedevice.UpdateDeviceRequestValidator, resUpdateDevice.UpdateDevice)

	patchDeviceRequest := updatedevice.PatchDeviceRequest{Repo: db}
	r.PATCH("/devices/:id", patchDeviceRequest.UpdatePartiallyDeviceRequestValidator, resUpdateDevice.UpdateDevice)

	resRemoveDevice := removedevice.HttpServer{Db: db}
	r.DELETE("/devices/:id", removedevice.RemoveDeviceRequestValidator, resRemoveDevice.RemoveDevice)

	resDevices := devices.HttpServer{Db: db}
	r.GET("/devices/:id", resDevices.GetDeviceById)
	r.GET("/devices/brand/:brand", resDevices.GetDeviceByBrand)
	r.GET("/devices", resDevices.GetDevices)

}
