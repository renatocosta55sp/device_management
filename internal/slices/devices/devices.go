package devices

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/renatocosta55sp/device_management/internal/infra/adapters/persistence"
)

type HttpServer struct {
	Db *pgxpool.Pool
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
