package device

import (
	"errors"
	"fmt"
	"time"

	"aquarium-control/internal/common"
	"aquarium-control/internal/database"
	"gorm.io/gorm"
)

type DeviceService struct {
	db *gorm.DB
}

func NewDeviceService() *DeviceService {
	return &DeviceService{
		db: database.GetDB(),
	}
}

func (s *DeviceService) ListDevices(deviceType string) ([]DeviceStatus, error) {
	var devices []DeviceStatus
	query := s.db.Model(&DeviceStatus{})
	if deviceType != "" {
		query = query.Where("device_type = ?", deviceType)
	}
	if err := query.Order("device_type, device_name").Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

func (s *DeviceService) GetDevice(deviceType, deviceName string) (*DeviceStatus, error) {
	var device DeviceStatus
	if err := s.db.Where("device_type = ? AND device_name = ?", deviceType, deviceName).
		First(&device).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("device not found")
		}
		return nil, err
	}
	return &device, nil
}

func (s *DeviceService) ToggleDevice(deviceType, deviceName string, status bool) (*DeviceStatus, error) {
	device, err := s.GetDevice(deviceType, deviceName)
	if err != nil {
		return nil, err
	}

	device.Status = status
	if !status {
		zero := 0
		device.CurrentValue = &zero
	}

	if err := s.db.Save(device).Error; err != nil {
		return nil, fmt.Errorf("failed to update device status: %w", err)
	}

	return device, nil
}

func (s *DeviceService) UpdateValue(deviceType, deviceName string, value int) (*DeviceStatus, error) {
	device, err := s.GetDevice(deviceType, deviceName)
	if err != nil {
		return nil, err
	}

	if deviceType == "light" {
		if value < 0 || value > 100 {
			return nil, errors.New("brightness must be between 0 and 100")
		}
	} else if deviceType == "pump" {
		if value < 1 || value > 5 {
			return nil, errors.New("pump level must be between 1 and 5")
		}
	}

	device.CurrentValue = &value
	if value > 0 {
		device.Status = true
	}

	if err := s.db.Save(device).Error; err != nil {
		return nil, fmt.Errorf("failed to update device value: %w", err)
	}

	return device, nil
}

func (s *DeviceService) SetManualMode(deviceType, deviceName string, manualMode bool) (*DeviceStatus, error) {
	device, err := s.GetDevice(deviceType, deviceName)
	if err != nil {
		return nil, err
	}

	device.ManualMode = manualMode
	if !manualMode {
		zero := 0
		device.Status = false
		device.CurrentValue = &zero
	}

	if err := s.db.Save(device).Error; err != nil {
		return nil, fmt.Errorf("failed to update manual mode: %w", err)
	}

	return device, nil
}

type LatestSensorTemp struct {
	Temperature float64 `gorm:"column:temperature"`
}

func (s *DeviceService) getLatestTemperature() (float64, error) {
	var data LatestSensorTemp
	result := s.db.Table("sensor_data").Order("recorded_at DESC").Limit(1).Find(&data)
	if result.Error != nil {
		return 0, result.Error
	}
	if result.RowsAffected == 0 {
		return 25.0, nil
	}
	return data.Temperature, nil
}

func (s *DeviceService) GetDashboardStatus() (map[string]interface{}, error) {
	now := time.Now()

	currentTemp, err := s.getLatestTemperature()
	if err != nil {
		currentTemp = 25.0
	}

	isNightMode := common.IsNightModeAt(now)
	isTempSafe := common.IsTemperatureSafe(currentTemp)
	powerSaving := isNightMode && isTempSafe

	var lights, pumps []DeviceStatus

	if err := s.db.Where("device_type = ?", "light").Find(&lights).Error; err != nil {
		return nil, err
	}
	if err := s.db.Where("device_type = ?", "pump").Find(&pumps).Error; err != nil {
		return nil, err
	}

	totalLights := len(lights)
	lightsOn := 0
	totalBrightness := 0
	for _, l := range lights {
		if l.Status {
			lightsOn++
			if l.CurrentValue != nil {
				totalBrightness += *l.CurrentValue
			}
		}
	}

	avgBrightness := 0
	if lightsOn > 0 {
		avgBrightness = totalBrightness / lightsOn
	}

	pumpsOn := 0
	var currentPumpLevel int
	for _, p := range pumps {
		if p.Status {
			pumpsOn++
			if p.CurrentValue != nil {
				currentPumpLevel = *p.CurrentValue
			}
		}
	}

	return map[string]interface{}{
		"lights":             lights,
		"pumps":              pumps,
		"total_lights":       totalLights,
		"lights_on":          lightsOn,
		"avg_brightness":     avgBrightness,
		"total_pumps":        len(pumps),
		"pumps_on":           pumpsOn,
		"current_pump_level": currentPumpLevel,
		"night_mode":         isNightMode,
		"power_saving":       powerSaving,
		"current_temp":       currentTemp,
		"temp_safe":          isTempSafe,
	}, nil
}
