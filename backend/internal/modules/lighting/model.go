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
	Name       string `json:"name" binding:"required,max=100"`
	StartTime  string `json:"start_time" binding:"required"`
	EndTime    string `json:"end_time" binding:"required"`
	Brightness int    `json:"brightness" binding:"required,min=0,max=100"`
	Enabled    bool   `json:"enabled"`
}

type UpdateScheduleRequest struct {
	Name       string `json:"name" binding:"omitempty,max=100"`
	StartTime  string `json:"start_time" binding:"omitempty"`
	EndTime    string `json:"end_time" binding:"omitempty"`
	Brightness *int   `json:"brightness" binding:"omitempty,min=0,max=100"`
	Enabled    *bool  `json:"enabled"`
}

type CurrentLightStatus struct {
	IsOn       bool   `json:"is_on"`
	Brightness int    `json:"brightness"`
	ScheduleID uint64 `json:"schedule_id,omitempty"`
	ScheduleName string `json:"schedule_name,omitempty"`
	NextAction string `json:"next_action"`
	NextTime   string `json:"next_time"`
}

func (LightSchedule) TableName() string {
	return "light_schedules"
}
