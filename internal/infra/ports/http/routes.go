package http

import (
	"github.com/gin-gonic/gin"
)

func InitRoutes(
	r *gin.RouterGroup,
	controller HttpServer) {
	r.POST("/devices", addDeviceValidator, controller.AddDevice)
	r.PUT("/devices/:id", updateDeviceValidator, controller.UpdateDevice)
	patchDeviceRequest := PatchDeviceRequest{repo: controller.Db}
	r.PATCH("/devices/:id", patchDeviceRequest.updatePartiallyDeviceValidator, controller.UpdateDevice)
	r.DELETE("/devices/:id", removeDeviceValidator, controller.RemoveDevice)
	r.GET("/devices/:id", controller.GetDeviceById)
	r.GET("/devices/brand/:brand", controller.GetDeviceByBrand)
	r.GET("/devices", controller.GetDevices)

}
