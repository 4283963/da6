package oxygen

import (
	"errors"
	"strconv"

	"aquarium-control/internal/common"

	"github.com/gin-gonic/gin"
)

type OxygenController struct {
	service *OxygenService
}

func NewOxygenController() *OxygenController {
	return &OxygenController{
		service: NewOxygenService(),
	}
}

func parseInt(val interface{}) (int, error) {
	switch v := val.(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	case string:
		if level, ok := common.ParsePumpLevel(v); ok {
			return level, nil
		}
		if num, err := strconv.Atoi(v); err == nil {
			return num, nil
		}
		return 0, errors.New("invalid integer value: " + v)
	default:
		return 0, errors.New("invalid integer type")
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

func (c *OxygenController) CreateConfig(ctx *gin.Context) {
	var req CreateConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.BadRequest(ctx, "Invalid request: "+err.Error())
		return
	}

	minLight, err := parseInt(req.MinLightWattage)
	if err != nil {
		common.BadRequest(ctx, "min_light_wattage: "+err.Error())
		return
	}

	maxLight, err := parseInt(req.MaxLightWattage)
	if err != nil {
		common.BadRequest(ctx, "max_light_wattage: "+err.Error())
		return
	}

	minTemp, err := parseFloat(req.MinTemperature)
	if err != nil {
		common.BadRequest(ctx, "min_temperature: "+err.Error())
		return
	}

	maxTemp, err := parseFloat(req.MaxTemperature)
	if err != nil {
		common.BadRequest(ctx, "max_temperature: "+err.Error())
		return
	}

	pumpLevel, err := parseInt(req.PumpLevel)
	if err != nil {
		common.BadRequest(ctx, "pump_level: "+err.Error())
		return
	}

	convertedReq := &createConfigDTO{
		MinLightWattage: minLight,
		MaxLightWattage: maxLight,
		MinTemperature:  minTemp,
		MaxTemperature:  maxTemp,
		PumpLevel:       pumpLevel,
		Description:     req.Description,
	}

	config, err := c.service.CreateConfig(convertedReq)
	if err != nil {
		common.InternalError(ctx, err.Error())
		return
	}

	common.Success(ctx, config)
}

func (c *OxygenController) GetConfig(ctx *gin.Context) {
	id, err := parseID(ctx.Param("id"))
	if err != nil {
		common.BadRequest(ctx, "Invalid config ID")
		return
	}

	config, err := c.service.GetConfigByID(id)
	if err != nil {
		common.NotFound(ctx, err.Error())
		return
	}

	common.Success(ctx, config)
}

func (c *OxygenController) ListConfigs(ctx *gin.Context) {
	configs, err := c.service.ListConfigs()
	if err != nil {
		common.InternalError(ctx, "Failed to list configs")
		return
	}

	common.Success(ctx, configs)
}

func (c *OxygenController) UpdateConfig(ctx *gin.Context) {
	id, err := parseID(ctx.Param("id"))
	if err != nil {
		common.BadRequest(ctx, "Invalid config ID")
		return
	}

	var req UpdateConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.BadRequest(ctx, "Invalid request: "+err.Error())
		return
	}

	convertedReq := &updateConfigDTO{
		Description: req.Description,
	}

	if req.MinLightWattage != nil {
		val, err := parseInt(req.MinLightWattage)
		if err != nil {
			common.BadRequest(ctx, "min_light_wattage: "+err.Error())
			return
		}
		convertedReq.MinLightWattage = &val
	}

	if req.MaxLightWattage != nil {
		val, err := parseInt(req.MaxLightWattage)
		if err != nil {
			common.BadRequest(ctx, "max_light_wattage: "+err.Error())
			return
		}
		convertedReq.MaxLightWattage = &val
	}

	if req.MinTemperature != nil {
		val, err := parseFloat(req.MinTemperature)
		if err != nil {
			common.BadRequest(ctx, "min_temperature: "+err.Error())
			return
		}
		convertedReq.MinTemperature = &val
	}

	if req.MaxTemperature != nil {
		val, err := parseFloat(req.MaxTemperature)
		if err != nil {
			common.BadRequest(ctx, "max_temperature: "+err.Error())
			return
		}
		convertedReq.MaxTemperature = &val
	}

	if req.PumpLevel != nil {
		val, err := parseInt(req.PumpLevel)
		if err != nil {
			common.BadRequest(ctx, "pump_level: "+err.Error())
			return
		}
		convertedReq.PumpLevel = &val
	}

	config, err := c.service.UpdateConfig(id, convertedReq)
	if err != nil {
		common.BadRequest(ctx, err.Error())
		return
	}

	common.Success(ctx, config)
}

func (c *OxygenController) DeleteConfig(ctx *gin.Context) {
	id, err := parseID(ctx.Param("id"))
	if err != nil {
		common.BadRequest(ctx, "Invalid config ID")
		return
	}

	if err := c.service.DeleteConfig(id); err != nil {
		common.NotFound(ctx, err.Error())
		return
	}

	common.Success(ctx, gin.H{"message": "Config deleted successfully"})
}

func (c *OxygenController) CalculateMatch(ctx *gin.Context) {
	var req MatchRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.BadRequest(ctx, "Invalid request: "+err.Error())
		return
	}

	lightWattage, err := parseInt(req.LightWattage)
	if err != nil {
		common.BadRequest(ctx, "light_wattage: "+err.Error())
		return
	}

	temperature, err := parseFloat(req.Temperature)
	if err != nil {
		common.BadRequest(ctx, "temperature: "+err.Error())
		return
	}

	dto := &matchRequestDTO{
		LightWattage: lightWattage,
		Temperature:  temperature,
	}

	result, err := c.service.CalculateMatch(dto)
	if err != nil {
		common.InternalError(ctx, "Failed to calculate match")
		return
	}

	common.Success(ctx, result)
}

func (c *OxygenController) GetMatchMatrix(ctx *gin.Context) {
	matrix, err := c.service.GetMatchMatrix()
	if err != nil {
		common.InternalError(ctx, "Failed to get match matrix")
		return
	}

	common.Success(ctx, matrix)
}

func RegisterRoutes(r *gin.RouterGroup) {
	controller := NewOxygenController()

	r.POST("/configs", controller.CreateConfig)
	r.GET("/configs", controller.ListConfigs)
	r.GET("/configs/:id", controller.GetConfig)
	r.PUT("/configs/:id", controller.UpdateConfig)
	r.DELETE("/configs/:id", controller.DeleteConfig)
	r.POST("/match", controller.CalculateMatch)
	r.GET("/matrix", controller.GetMatchMatrix)
}

func parseID(idStr string) (uint64, error) {
	return strconv.ParseUint(idStr, 10, 64)
}
