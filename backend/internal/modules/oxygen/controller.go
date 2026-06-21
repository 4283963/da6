package oxygen

import (
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

func (c *OxygenController) CreateConfig(ctx *gin.Context) {
	var req CreateConfigRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.BadRequest(ctx, "Invalid request: "+err.Error())
		return
	}

	config, err := c.service.CreateConfig(&req)
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

	config, err := c.service.UpdateConfig(id, &req)
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

	result, err := c.service.CalculateMatch(req.LightWattage, req.Temperature)
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
