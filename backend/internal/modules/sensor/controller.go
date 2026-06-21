package sensor

import (
	"errors"
	"strconv"

	"aquarium-control/internal/common"

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

func parseFloat(val interface{}) (float64, error) {
	switch v := val.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case string:
		if num, err := strconv.ParseFloat(v, 64); err == nil {
			return num, nil
		}
		return 0, errors.New("invalid float value: " + v)
	default:
		return 0, errors.New("invalid float type")
	}
}

func parseInt(val interface{}) (int, error) {
	switch v := val.(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	case string:
		if num, err := strconv.Atoi(v); err == nil {
			return num, nil
		}
		return 0, errors.New("invalid integer value: " + v)
	default:
		return 0, errors.New("invalid integer type")
	}
}

func (c *SensorController) CreateData(ctx *gin.Context) {
	var req CreateSensorDataRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.BadRequest(ctx, "Invalid request: "+err.Error())
		return
	}

	temperature, err := parseFloat(req.Temperature)
	if err != nil {
		common.BadRequest(ctx, "temperature: "+err.Error())
		return
	}

	lightWattage, err := parseInt(req.LightWattage)
	if err != nil {
		common.BadRequest(ctx, "light_wattage: "+err.Error())
		return
	}

	var dissolvedOxygen *float64
	if req.DissolvedOxygen != nil {
		do, err := parseFloat(req.DissolvedOxygen)
		if err != nil {
			common.BadRequest(ctx, "dissolved_oxygen: "+err.Error())
			return
		}
		dissolvedOxygen = &do
	}

	convertedReq := &createSensorDataDTO{
		Temperature:     temperature,
		LightWattage:    lightWattage,
		DissolvedOxygen: dissolvedOxygen,
	}

	data, err := c.service.CreateData(convertedReq)
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
