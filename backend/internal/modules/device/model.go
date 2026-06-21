package device

import (
	"time"

	"gorm.io/gorm"
)

type DeviceStatus struct {
	ID          uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	DeviceType  string         `gorm:"size:50;not null;uniqueIndex:uk_device_type_name" json:"device_type"`
	DeviceName  string         `gorm:"size:100;not null;uniqueIndex:uk_device_type_name" json:"device_name"`
	Status      bool           `gorm:"default:false" json:"status"`
	CurrentValue *int          `gorm:"column:current_value" json:"current_value"`
	ManualMode  bool           `gorm:"default:false" json:"manual_mode"`
	LastUpdated time.Time      `gorm:"autoUpdateTime" json:"last_updated"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type ToggleDeviceRequest struct {
	Status *bool `json:"status" binding:"required"`
}

type UpdateValueRequest struct {
	CurrentValue *int `json:"current_value" binding:"required"`
}

type SetManualModeRequest struct {
	ManualMode *bool `json:"manual_mode" binding:"required"`
}

func (DeviceStatus) TableName() string {
	return "device_status"
}
