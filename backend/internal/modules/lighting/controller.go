package lighting

import (
	"errors"
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

func parseBrightness(val interface{}) (int, error) {
	switch v := val.(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	case string:
		if result, ok := common.ParseBrightness(v); ok {
			return result, nil
		}
		if num, err := strconv.Atoi(v); err == nil {
			return num, nil
		}
		return 0, errors.New("invalid brightness value: " + v)
	default:
		return 0, errors.New("invalid brightness type")
	}
}

func parseEnabled(val interface{}) (bool, bool, error) {
	if val == nil {
		return false, false, nil
	}
	switch v := val.(type) {
	case bool:
		return v, true, nil
	case string:
		if result, ok := common.ParseStatus(v); ok {
			return result, true, nil
		}
		return false, false, errors.New("invalid enabled value: " + v)
	case float64:
		return v != 0, true, nil
	default:
		return false, false, errors.New("invalid enabled type")
	}
}

func (c *LightController) CreateSchedule(ctx *gin.Context) {
	var req CreateScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.BadRequest(ctx, "Invalid request: "+err.Error())
		return
	}

	brightness, err := parseBrightness(req.Brightness)
	if err != nil {
		common.BadRequest(ctx, err.Error())
		return
	}

	enabled, hasEnabled, err := parseEnabled(req.Enabled)
	if err != nil {
		common.BadRequest(ctx, err.Error())
		return
	}

	if req.TimeRange != "" {
		start, end, ok := common.ParseTimeRange(req.TimeRange)
		if ok {
			req.StartTime = start
			req.EndTime = end
		}
	}

	if req.StartTime == "" || req.EndTime == "" {
		common.BadRequest(ctx, "start_time and end_time are required, or use time_range")
		return
	}

	convertedReq := &createScheduleDTO{
		Name:       req.Name,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		Brightness: brightness,
		Enabled:    enabled,
		HasEnabled: hasEnabled,
	}

	schedule, err := c.service.CreateSchedule(convertedReq)
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

	convertedReq := &updateScheduleDTO{
		Name:      req.Name,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	if req.TimeRange != "" {
		start, end, ok := common.ParseTimeRange(req.TimeRange)
		if ok {
			convertedReq.StartTime = start
			convertedReq.EndTime = end
		}
	}

	if req.Brightness != nil {
		brightness, err := parseBrightness(req.Brightness)
		if err != nil {
			common.BadRequest(ctx, err.Error())
			return
		}
		convertedReq.Brightness = &brightness
		convertedReq.HasBrightness = true
	}

	if req.Enabled != nil {
		enabled, hasEnabled, err := parseEnabled(req.Enabled)
		if err != nil {
			common.BadRequest(ctx, err.Error())
			return
		}
		if hasEnabled {
			convertedReq.Enabled = &enabled
			convertedReq.HasEnabled = true
		}
	}

	schedule, err := c.service.UpdateSchedule(id, convertedReq)
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
