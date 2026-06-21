package router

import (
	"aquarium-control/internal/modules/device"
	"aquarium-control/internal/modules/lighting"
	"aquarium-control/internal/modules/oxygen"
	"aquarium-control/internal/modules/sensor"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		lighting.RegisterRoutes(api.Group("/lighting"))
		oxygen.RegisterRoutes(api.Group("/oxygen"))
		device.RegisterRoutes(api.Group("/device"))
		sensor.RegisterRoutes(api.Group("/sensor"))
	}
}
