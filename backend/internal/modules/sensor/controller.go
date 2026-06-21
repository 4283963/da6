package sensor

import (
	"aquarium-control/internal/common"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SensorController struct {
	service *SensorService
}

func NewSensorController() *SensorController {
	return &SensorController{
		service: NewSensorService(),
	}
}

func (c *SensorController) CreateData(ctx *gin.Context) {
	var req CreateSensorDataRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.BadRequest(ctx, "Invalid request: "+err.Error())
		return
	}

	data, err := c.service.CreateData(&req)
	if err != nil {
		common.InternalError(ctx, err.Error())
		return
	}

	common.Success(ctx, data)
}

func (c *SensorController) ListData(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "100"))
	hours, _ := strconv.Atoi(ctx.DefaultQuery("hours", "24"))

	data, err := c.service.ListData(limit, hours)
	if err != nil {
		common.InternalError(ctx, "Failed to list sensor data")
		return
	}

	common.Success(ctx, data)
}

func (c *SensorController) GetLatest(ctx *gin.Context) {
	data, err := c.service.GetLatest()
	if err != nil {
		common.NotFound(ctx, err.Error())
		return
	}

	common.Success(ctx, data)
}

func (c *SensorController) GetStats(ctx *gin.Context) {
	hours, _ := strconv.Atoi(ctx.DefaultQuery("hours", "24"))

	stats, err := c.service.GetStats(hours)
	if err != nil {
		common.InternalError(ctx, "Failed to get sensor stats")
		return
	}

	common.Success(ctx, stats)
}

func RegisterRoutes(r *gin.RouterGroup) {
	controller := NewSensorController()

	r.POST("/data", controller.CreateData)
	r.GET("/data", controller.ListData)
	r.GET("/data/latest", controller.GetLatest)
	r.GET("/stats", controller.GetStats)
}
