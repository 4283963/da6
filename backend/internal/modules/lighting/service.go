package lighting

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"aquarium-control/internal/common"
	"aquarium-control/internal/database"
	"gorm.io/gorm"
)

type LightService struct {
	db *gorm.DB
}

func NewLightService() *LightService {
	return &LightService{
		db: database.GetDB(),
	}
}

func (s *LightService) CreateSchedule(req *createScheduleDTO) (*LightSchedule, error) {
	if !isValidTimeFormat(req.StartTime) || !isValidTimeFormat(req.EndTime) {
		return nil, errors.New("invalid time format, use HH:MM:SS")
	}

	if req.StartTime >= req.EndTime {
		return nil, errors.New("start time must be before end time")
	}

	enabled := true
	if req.HasEnabled {
		enabled = req.Enabled
	}

	schedule := &LightSchedule{
		Name:       req.Name,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		Brightness: req.Brightness,
		Enabled:    enabled,
	}

	if err := s.db.Create(schedule).Error; err != nil {
		return nil, fmt.Errorf("failed to create schedule: %w", err)
	}

	return schedule, nil
}

func (s *LightService) GetScheduleByID(id uint64) (*LightSchedule, error) {
	var schedule LightSchedule
	if err := s.db.First(&schedule, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("schedule not found")
		}
		return nil, err
	}
	return &schedule, nil
}

func (s *LightService) ListSchedules() ([]LightSchedule, error) {
	var schedules []LightSchedule
	if err := s.db.Order("start_time ASC").Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (s *LightService) UpdateSchedule(id uint64, req *updateScheduleDTO) (*LightSchedule, error) {
	schedule, err := s.GetScheduleByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		schedule.Name = req.Name
	}
	if req.StartTime != "" {
		if !isValidTimeFormat(req.StartTime) {
			return nil, errors.New("invalid start time format")
		}
		schedule.StartTime = req.StartTime
	}
	if req.EndTime != "" {
		if !isValidTimeFormat(req.EndTime) {
			return nil, errors.New("invalid end time format")
		}
		schedule.EndTime = req.EndTime
	}
	if req.HasBrightness {
		schedule.Brightness = *req.Brightness
	}
	if req.HasEnabled {
		schedule.Enabled = *req.Enabled
	}

	if schedule.StartTime >= schedule.EndTime {
		return nil, errors.New("start time must be before end time")
	}

	if err := s.db.Save(schedule).Error; err != nil {
		return nil, err
	}

	return schedule, nil
}

func (s *LightService) DeleteSchedule(id uint64) error {
	result := s.db.Delete(&LightSchedule{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("schedule not found")
	}
	return nil
}

type LatestSensorData struct {
	Temperature float64 `gorm:"column:temperature"`
}

func (s *LightService) getLatestTemperature() (float64, error) {
	var data LatestSensorData
	result := s.db.Table("sensor_data").Order("recorded_at DESC").Limit(1).Find(&data)
	if result.Error != nil {
		return 0, result.Error
	}
	if result.RowsAffected == 0 {
		return 25.0, nil
	}
	return data.Temperature, nil
}

func (s *LightService) GetCurrentStatus() (*CurrentLightStatus, error) {
	now := time.Now()
	currentTimeStr := now.Format("15:04:05")

	currentTemp, err := s.getLatestTemperature()
	if err != nil {
		currentTemp = 25.0
	}

	isNightMode := common.IsNightModeAt(now)
	isTempSafe := common.IsTemperatureSafe(currentTemp)
	powerSaving := isNightMode && isTempSafe

	var schedules []LightSchedule
	if err := s.db.Where("enabled = ?", true).Order("start_time ASC").Find(&schedules).Error; err != nil {
		return nil, err
	}

	status := &CurrentLightStatus{
		IsOn:        false,
		Brightness:  0,
		NightMode:   isNightMode,
		PowerSaving: powerSaving,
		CurrentTemp: currentTemp,
	}

	sort.Slice(schedules, func(i, j int) bool {
		return schedules[i].StartTime < schedules[j].StartTime
	})

	var activeSchedule *LightSchedule
	for i := range schedules {
		sched := &schedules[i]
		if sched.StartTime <= currentTimeStr && currentTimeStr < sched.EndTime {
			activeSchedule = sched
			break
		}
	}

	if activeSchedule != nil {
		status.IsOn = true
		originalBrightness := activeSchedule.Brightness
		status.OriginalBrightness = originalBrightness
		status.ScheduleID = activeSchedule.ID
		status.ScheduleName = activeSchedule.Name
		status.NextAction = "关灯"
		status.NextTime = activeSchedule.EndTime

		if powerSaving {
			status.Brightness = common.ApplyPowerSavingWithMin(originalBrightness, 0)
		} else {
			status.Brightness = originalBrightness
		}
	} else {
		for i := range schedules {
			sched := &schedules[i]
			if sched.StartTime > currentTimeStr {
				status.NextAction = "开灯"
				status.NextTime = sched.StartTime
				break
			}
		}
		if status.NextTime == "" && len(schedules) > 0 {
			status.NextAction = "开灯"
			status.NextTime = schedules[0].StartTime
		}
	}

	return status, nil
}

func (s *LightService) CalculateWattage(brightness int, maxWattage int) int {
	if brightness <= 0 {
		return 0
	}
	if brightness >= 100 {
		return maxWattage
	}
	return int(float64(maxWattage) * float64(brightness) / 100.0)
}

func isValidTimeFormat(t string) bool {
	if len(t) != 8 {
		return false
	}
	_, err := time.Parse("15:04:05", t)
	return err == nil
}
