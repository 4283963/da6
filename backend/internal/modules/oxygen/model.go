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
	MinLightWattage int     `json:"min_light_wattage" binding:"required,min=0"`
	MaxLightWattage int     `json:"max_light_wattage" binding:"required,gtefield=MinLightWattage"`
	MinTemperature  float64 `json:"min_temperature" binding:"required,min=0,max=40"`
	MaxTemperature  float64 `json:"max_temperature" binding:"required,gtefield=MinTemperature,min=0,max=40"`
	PumpLevel       int     `json:"pump_level" binding:"required,min=1,max=5"`
	Description     string  `json:"description" binding:"max=255"`
}

type UpdateConfigRequest struct {
	MinLightWattage *int     `json:"min_light_wattage" binding:"omitempty,min=0"`
	MaxLightWattage *int     `json:"max_light_wattage" binding:"omitempty,min=0"`
	MinTemperature  *float64 `json:"min_temperature" binding:"omitempty,min=0,max=40"`
	MaxTemperature  *float64 `json:"max_temperature" binding:"omitempty,min=0,max=40"`
	PumpLevel       *int     `json:"pump_level" binding:"omitempty,min=1,max=5"`
	Description     *string  `json:"description" binding:"omitempty,max=255"`
}

type MatchRequest struct {
	LightWattage int     `json:"light_wattage" binding:"required,min=0"`
	Temperature  float64 `json:"temperature" binding:"required,min=0,max=40"`
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
