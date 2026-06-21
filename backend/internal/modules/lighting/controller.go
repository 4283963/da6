package lighting

import (
	"strconv"

	"aquarium-control/internal/common"

	"github.com/gin-gonic/gin"
)

type LightController struct {
	service *LightService
}

func NewLightController() *LightController {
	return &LightController{
		service: NewLightService(),
	}
}

func (c *LightController) CreateSchedule(ctx *gin.Context) {
	var req CreateScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.BadRequest(ctx, "Invalid request: "+err.Error())
		return
	}

	schedule, err := c.service.CreateSchedule(&req)
	if err != nil {
		common.BadRequest(ctx, err.Error())
		return
	}

	common.Success(ctx, schedule)
}

func (c *LightController) GetSchedule(ctx *gin.Context) {
	id, err := parseID(ctx.Param("id"))
	if err != nil {
		common.BadRequest(ctx, "Invalid schedule ID")
		return
	}

	schedule, err := c.service.GetScheduleByID(id)
	if err != nil {
		common.NotFound(ctx, err.Error())
		return
	}

	common.Success(ctx, schedule)
}

func (c *LightController) ListSchedules(ctx *gin.Context) {
	schedules, err := c.service.ListSchedules()
	if err != nil {
		common.InternalError(ctx, "Failed to list schedules")
		return
	}

	common.Success(ctx, schedules)
}

func (c *LightController) UpdateSchedule(ctx *gin.Context) {
	id, err := parseID(ctx.Param("id"))
	if err != nil {
		common.BadRequest(ctx, "Invalid schedule ID")
		return
	}

	var req UpdateScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.BadRequest(ctx, "Invalid request: "+err.Error())
		return
	}

	schedule, err := c.service.UpdateSchedule(id, &req)
	if err != nil {
		common.BadRequest(ctx, err.Error())
		return
	}

	common.Success(ctx, schedule)
}

func (c *LightController) DeleteSchedule(ctx *gin.Context) {
	id, err := parseID(ctx.Param("id"))
	if err != nil {
		common.BadRequest(ctx, "Invalid schedule ID")
		return
	}

	if err := c.service.DeleteSchedule(id); err != nil {
		common.NotFound(ctx, err.Error())
		return
	}

	common.Success(ctx, gin.H{"message": "Schedule deleted successfully"})
}

func (c *LightController) GetCurrentStatus(ctx *gin.Context) {
	status, err := c.service.GetCurrentStatus()
	if err != nil {
		common.InternalError(ctx, "Failed to get current status")
		return
	}

	common.Success(ctx, status)
}

func RegisterRoutes(r *gin.RouterGroup) {
	controller := NewLightController()

	r.POST("/schedules", controller.CreateSchedule)
	r.GET("/schedules", controller.ListSchedules)
	r.GET("/schedules/:id", controller.GetSchedule)
	r.PUT("/schedules/:id", controller.UpdateSchedule)
	r.DELETE("/schedules/:id", controller.DeleteSchedule)
	r.GET("/status", controller.GetCurrentStatus)
}

func parseID(idStr string) (uint64, error) {
	return strconv.ParseUint(idStr, 10, 64)
}
