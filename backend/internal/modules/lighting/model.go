package lighting

import (
	"time"

	"gorm.io/gorm"
)

type LightSchedule struct {
	ID         uint64         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string         `gorm:"size:100;not null" json:"name"`
	StartTime  string         `gorm:"type:time;not null" json:"start_time"`
	EndTime    string         `gorm:"type:time;not null" json:"end_time"`
	Brightness int            `gorm:"not null" json:"brightness"`
	Enabled    bool           `gorm:"default:true" json:"enabled"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreateScheduleRequest struct {
	Name       string      `json:"name" binding:"required,max=100"`
	StartTime  string      `json:"start_time"`
	EndTime    string      `json:"end_time"`
	TimeRange  string      `json:"time_range"`
	Brightness interface{} `json:"brightness" binding:"required"`
	Enabled    interface{} `json:"enabled"`
}

type UpdateScheduleRequest struct {
	Name       string      `json:"name" binding:"omitempty,max=100"`
	StartTime  string      `json:"start_time"`
	EndTime    string      `json:"end_time"`
	TimeRange  string      `json:"time_range"`
	Brightness interface{} `json:"brightness"`
	Enabled    interface{} `json:"enabled"`
}

type createScheduleDTO struct {
	Name       string
	StartTime  string
	EndTime    string
	Brightness int
	Enabled    bool
	HasEnabled bool
}

type updateScheduleDTO struct {
	Name         string
	StartTime    string
	EndTime      string
	Brightness   *int
	Enabled      *bool
	HasBrightness bool
	HasEnabled   bool
}

type CurrentLightStatus struct {
	IsOn              bool   `json:"is_on"`
	Brightness        int    `json:"brightness"`
	OriginalBrightness int   `json:"original_brightness,omitempty"`
	NightMode         bool   `json:"night_mode"`
	PowerSaving       bool   `json:"power_saving"`
	ScheduleID        uint64 `json:"schedule_id,omitempty"`
	ScheduleName      string `json:"schedule_name,omitempty"`
	NextAction        string `json:"next_action"`
	NextTime          string `json:"next_time"`
	CurrentTemp       float64 `json:"current_temp,omitempty"`
}

func (LightSchedule) TableName() string {
	return "light_schedules"
}
