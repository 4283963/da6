package device

import (
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

func (c *DeviceController) ListDevices(ctx *gin.Context) {
	deviceType := ctx.Query("type")
	devices, err := c.service.ListDevices(deviceType)
	if err != nil {
		common.InternalError(ctx, "Failed to list devices")
		return
	}
	common.Success(ctx, devices)
}

func (c *DeviceController) GetDevice(ctx *gin.Context) {
	deviceType := ctx.Param("type")
	deviceName := ctx.Param("name")
	device, err := c.service.GetDevice(deviceType, deviceName)
	if err != nil {
		common.NotFound(ctx, err.Error())
		return
	}
	common.Success(ctx, device)
}

func (c *DeviceController) ToggleDevice(ctx *gin.Context) {
	deviceType := ctx.Param("type")
	deviceName := ctx.Param("name")

	var req ToggleDeviceRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.BadRequest(ctx, "Invalid request: "+err.Error())
		return
	}

	device, err := c.service.ToggleDevice(deviceType, deviceName, *req.Status)
	if err != nil {
		common.BadRequest(ctx, err.Error())
		return
	}
	common.Success(ctx, device)
}

func (c *DeviceController) UpdateValue(ctx *gin.Context) {
	deviceType := ctx.Param("type")
	deviceName := ctx.Param("name")

	var req UpdateValueRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.BadRequest(ctx, "Invalid request: "+err.Error())
		return
	}

	device, err := c.service.UpdateValue(deviceType, deviceName, *req.CurrentValue)
	if err != nil {
		common.BadRequest(ctx, err.Error())
		return
	}
	common.Success(ctx, device)
}

func (c *DeviceController) SetManualMode(ctx *gin.Context) {
	deviceType := ctx.Param("type")
	deviceName := ctx.Param("name")

	var req SetManualModeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		common.BadRequest(ctx, "Invalid request: "+err.Error())
		return
	}

	device, err := c.service.SetManualMode(deviceType, deviceName, *req.ManualMode)
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
