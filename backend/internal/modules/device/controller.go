package device

import (
	"errors"

	"aquarium-control/internal/common"

	"github.com/gin-gonic/gin"
)

type DeviceController struct {
	service *DeviceService
}

func NewDeviceController() *DeviceController {
	return &DeviceController{
		service: NewDeviceService(),
	}
}

func parseBoolValue(val interface{}) (bool, error) {
	switch v := val.(type) {
	case bool:
		return v, nil
	case string:
		if result, ok := common.ParseStatus(v); ok {
			return result, nil
		}
		if result, ok := common.ParseMode(v); ok {
			return result, nil
		}
		return false, errors.New("invalid status/mode value: " + v)
	case float64:
		return v != 0, nil
	default:
		return false, errors.New("invalid value type")
	}
}

func parseIntValue(val interface{}, deviceType string) (int, error) {
	switch v := val.(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	case string:
		if deviceType == "light" {
			if result, ok := common.ParseBrightness(v); ok {
				return result, nil
			}
		} else if deviceType == "pump" {
			if result, ok := common.ParsePumpLevel(v); ok {
				return result, nil
			}
		}
		return 0, errors.New("invalid value: " + v)
	default:
		return 0, errors.New("invalid value type")
	}
}

func (c *DeviceController) ListDevices(ctx *gin.Context) {
	deviceType := common.NormalizeDeviceType(ctx.Query("type"))
	devices, err := c.service.ListDevices(deviceType)
	if err != nil {
		common.InternalError(ctx, "Failed to list devices")
		return
	}
	common.Success(ctx, devices)
}

func (c *DeviceController) GetDevice(ctx *gin.Context) {
	deviceType := common.NormalizeDeviceType(ctx.Param("type"))
	deviceName := common.NormalizeDeviceName(ctx.Param("name"))
	device, err := c.service.GetDevice(deviceType, deviceName)
	if err != nil {
		common.NotFound(ctx, err.Error())
		return
	}
	common.Success(ctx, device)
}

func (c *DeviceController) ToggleDevice(ctx *gin.Context) {
	deviceType := common.NormalizeDeviceType(ctx.Param("type"))
	deviceName := common.NormalizeDeviceName(ctx.Param("name"))

	var req ToggleDeviceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.BadRequest(ctx, "Invalid request: "+err.Error())
		return
	}

	status, err := parseBoolValue(req.Status)
	if err != nil {
		common.BadRequest(ctx, err.Error())
		return
	}

	device, err := c.service.ToggleDevice(deviceType, deviceName, status)
	if err != nil {
		common.BadRequest(ctx, err.Error())
		return
	}
	common.Success(ctx, device)
}

func (c *DeviceController) UpdateValue(ctx *gin.Context) {
	deviceType := common.NormalizeDeviceType(ctx.Param("type"))
	deviceName := common.NormalizeDeviceName(ctx.Param("name"))

	var req UpdateValueRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.BadRequest(ctx, "Invalid request: "+err.Error())
		return
	}

	value, err := parseIntValue(req.CurrentValue, deviceType)
	if err != nil {
		common.BadRequest(ctx, err.Error())
		return
	}

	device, err := c.service.UpdateValue(deviceType, deviceName, value)
	if err != nil {
		common.BadRequest(ctx, err.Error())
		return
	}
	common.Success(ctx, device)
}

func (c *DeviceController) SetManualMode(ctx *gin.Context) {
	deviceType := common.NormalizeDeviceType(ctx.Param("type"))
	deviceName := common.NormalizeDeviceName(ctx.Param("name"))

	var req SetManualModeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.BadRequest(ctx, "Invalid request: "+err.Error())
		return
	}

	manualMode, err := parseBoolValue(req.ManualMode)
	if err != nil {
		common.BadRequest(ctx, err.Error())
		return
	}

	device, err := c.service.SetManualMode(deviceType, deviceName, manualMode)
	if err != nil {
		common.BadRequest(ctx, err.Error())
		return
	}
	common.Success(ctx, device)
}

func (c *DeviceController) GetDashboardStatus(ctx *gin.Context) {
	status, err := c.service.GetDashboardStatus()
	if err != nil {
		common.InternalError(ctx, "Failed to get dashboard status")
		return
	}
	common.Success(ctx, status)
}

func RegisterRoutes(r *gin.RouterGroup) {
	controller := NewDeviceController()

	r.GET("/", controller.ListDevices)
	r.GET("/dashboard", controller.GetDashboardStatus)
	r.GET("/:type/:name", controller.GetDevice)
	r.PUT("/:type/:name/toggle", controller.ToggleDevice)
	r.PUT("/:type/:name/value", controller.UpdateValue)
	r.PUT("/:type/:name/manual", controller.SetManualMode)
}
