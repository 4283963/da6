package sensor

import (
	"fmt"
	"time"

	"aquarium-control/internal/database"
	"gorm.io/gorm"
)

type SensorService struct {
	db *gorm.DB
}

func NewSensorService() *SensorService {
	return &SensorService{
		db: database.GetDB(),
	}
}

func (s *SensorService) CreateData(req *CreateSensorDataRequest) (*SensorData, error) {
	data := &SensorData{
		Temperature:     req.Temperature,
		LightWattage:    req.LightWattage,
		DissolvedOxygen: req.DissolvedOxygen,
	}

	if err := s.db.Create(data).Error; err != nil {
		return nil, fmt.Errorf("failed to create sensor data: %w", err)
	}

	return data, nil
}

func (s *SensorService) ListData(limit int, hours int) ([]SensorData, error) {
	var data []SensorData
	query := s.db.Order("recorded_at DESC")

	if hours > 0 {
		since := time.Now().Add(-time.Duration(hours) * time.Hour)
		query = query.Where("recorded_at >= ?", since)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func (s *SensorService) GetLatest() (*SensorData, error) {
	var data SensorData
	if err := s.db.Order("recorded_at DESC").First(&data).Error; err != nil {
		return nil, fmt.Errorf("no sensor data found: %w", err)
	}
	return &data, nil
}

func (s *SensorService) GetStats(hours int) (map[string]interface{}, error) {
	var data []SensorData
	query := s.db.Model(&SensorData{})

	if hours > 0 {
		since := time.Now().Add(-time.Duration(hours) * time.Hour)
		query = query.Where("recorded_at >= ?", since)
	}

	if err := query.Order("recorded_at DESC").Find(&data).Error; err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return map[string]interface{}{
			"count":              0,
			"avg_temperature":    0,
			"avg_light_wattage":  0,
			"avg_do":             nil,
			"current_temperature": nil,
			"current_light_wattage": 0,
		}, nil
	}

	totalTemp := 0.0
	totalLight := 0
	totalDO := 0.0
	doCount := 0

	for _, d := range data {
		totalTemp += d.Temperature
		totalLight += d.LightWattage
		if d.DissolvedOxygen != nil {
			totalDO += *d.DissolvedOxygen
			doCount++
		}
	}

	avgTemp := totalTemp / float64(len(data))
	avgLight := totalLight / len(data)
	avgDO := (*float64)(nil)
	if doCount > 0 {
		avg := totalDO / float64(doCount)
		avgDO = &avg
	}

	current := data[0]

	return map[string]interface{}{
		"count":                 len(data),
		"avg_temperature":       roundFloat(avgTemp, 1),
		"avg_light_wattage":     avgLight,
		"avg_do":                avgDO,
		"current_temperature":   current.Temperature,
		"current_light_wattage": current.LightWattage,
		"current_do":            current.DissolvedOxygen,
		"recorded_at":           current.RecordedAt,
	}, nil
}

func roundFloat(val float64, precision int) float64 {
	ratio := 1.0
	for i := 0; i < precision; i++ {
		ratio *= 10
	}
	return float64(int(val*ratio+0.5)) / ratio
}
