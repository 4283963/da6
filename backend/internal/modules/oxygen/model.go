package oxygen

import (
	"time"

	"gorm.io/gorm"
)

type OxygenConfig struct {
	ID              uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	MinLightWattage int            `gorm:"not null" json:"min_light_wattage"`
	MaxLightWattage int            `gorm:"not null" json:"max_light_wattage"`
	MinTemperature  float64        `gorm:"type:decimal(4,1);not null" json:"min_temperature"`
	MaxTemperature  float64        `gorm:"type:decimal(4,1);not null" json:"max_temperature"`
	PumpLevel       int            `gorm:"not null" json:"pump_level"`
	Description     string         `gorm:"size:255" json:"description"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreateConfigRequest struct {
	MinLightWattage interface{} `json:"min_light_wattage" binding:"required"`
	MaxLightWattage interface{} `json:"max_light_wattage" binding:"required"`
	MinTemperature  interface{} `json:"min_temperature" binding:"required"`
	MaxTemperature  interface{} `json:"max_temperature" binding:"required"`
	PumpLevel       interface{} `json:"pump_level" binding:"required"`
	Description     string      `json:"description" binding:"max=255"`
}

type UpdateConfigRequest struct {
	MinLightWattage interface{} `json:"min_light_wattage"`
	MaxLightWattage interface{} `json:"max_light_wattage"`
	MinTemperature  interface{} `json:"min_temperature"`
	MaxTemperature  interface{} `json:"max_temperature"`
	PumpLevel       interface{} `json:"pump_level"`
	Description     *string     `json:"description" binding:"omitempty,max=255"`
}

type MatchRequest struct {
	LightWattage interface{} `json:"light_wattage" binding:"required"`
	Temperature  interface{} `json:"temperature" binding:"required"`
}

type createConfigDTO struct {
	MinLightWattage int
	MaxLightWattage int
	MinTemperature  float64
	MaxTemperature  float64
	PumpLevel       int
	Description     string
}

type updateConfigDTO struct {
	MinLightWattage *int
	MaxLightWattage *int
	MinTemperature  *float64
	MaxTemperature  *float64
	PumpLevel       *int
	Description     *string
}

type matchRequestDTO struct {
	LightWattage int
	Temperature  float64
}

type MatchResult struct {
	PumpLevel   int    `json:"pump_level"`
	Description string `json:"description"`
	Formula     string `json:"formula"`
	Reason      string `json:"reason"`
}

func (OxygenConfig) TableName() string {
	return "oxygen_configs"
}
