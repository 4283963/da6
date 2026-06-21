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
	Temperature     float64  `json:"temperature" binding:"required,min=0,max=40"`
	LightWattage    int      `json:"light_wattage" binding:"required,min=0,max=1000"`
	DissolvedOxygen *float64 `json:"dissolved_oxygen" binding:"omitempty,min=0,max=20"`
}

func (SensorData) TableName() string {
	return "sensor_data"
}
