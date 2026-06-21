package sensor

import (
	"time"

	"gorm.io/gorm"
)

type SensorData struct {
	ID              uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Temperature     float64        `gorm:"type:decimal(4,1);not null" json:"temperature"`
	LightWattage    int            `gorm:"not null" json:"light_wattage"`
	DissolvedOxygen *float64       `gorm:"type:decimal(4,2)" json:"dissolved_oxygen"`
	RecordedAt      time.Time      `gorm:"autoCreateTime" json:"recorded_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreateSensorDataRequest struct {
	Temperature     interface{} `json:"temperature" binding:"required"`
	LightWattage    interface{} `json:"light_wattage" binding:"required"`
	DissolvedOxygen interface{} `json:"dissolved_oxygen"`
}

type createSensorDataDTO struct {
	Temperature     float64
	LightWattage    int
	DissolvedOxygen *float64
}

func (SensorData) TableName() string {
	return "sensor_data"
}
