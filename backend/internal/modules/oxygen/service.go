package oxygen

import (
	"errors"
	"fmt"
	"math"

	"aquarium-control/internal/database"
	"gorm.io/gorm"
)

type OxygenService struct {
	db *gorm.DB
}

func NewOxygenService() *OxygenService {
	return &OxygenService{
		db: database.GetDB(),
	}
}

func (s *OxygenService) CreateConfig(req *createConfigDTO) (*OxygenConfig, error) {
	if req.MinLightWattage < 0 {
		return nil, errors.New("min_light_wattage must be >= 0")
	}
	if req.MaxLightWattage < req.MinLightWattage {
		return nil, errors.New("max_light_wattage must be >= min_light_wattage")
	}
	if req.MinTemperature < 0 || req.MinTemperature > 40 {
		return nil, errors.New("min_temperature must be between 0 and 40")
	}
	if req.MaxTemperature < req.MinTemperature || req.MaxTemperature > 40 {
		return nil, errors.New("max_temperature must be >= min_temperature and <= 40")
	}
	if req.PumpLevel < 1 || req.PumpLevel > 5 {
		return nil, errors.New("pump_level must be between 1 and 5")
	}

	config := &OxygenConfig{
		MinLightWattage: req.MinLightWattage,
		MaxLightWattage: req.MaxLightWattage,
		MinTemperature:  req.MinTemperature,
		MaxTemperature:  req.MaxTemperature,
		PumpLevel:       req.PumpLevel,
		Description:     req.Description,
	}

	if err := s.db.Create(config).Error; err != nil {
		return nil, fmt.Errorf("failed to create config: %w", err)
	}

	return config, nil
}

func (s *OxygenService) GetConfigByID(id uint64) (*OxygenConfig, error) {
	var config OxygenConfig
	if err := s.db.First(&config, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("config not found")
		}
		return nil, err
	}
	return &config, nil
}

func (s *OxygenService) ListConfigs() ([]OxygenConfig, error) {
	var configs []OxygenConfig
	if err := s.db.Order("min_light_wattage ASC, min_temperature ASC").Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

func (s *OxygenService) UpdateConfig(id uint64, req *updateConfigDTO) (*OxygenConfig, error) {
	config, err := s.GetConfigByID(id)
	if err != nil {
		return nil, err
	}

	if req.MinLightWattage != nil {
		if *req.MinLightWattage < 0 {
			return nil, errors.New("min_light_wattage must be >= 0")
		}
		config.MinLightWattage = *req.MinLightWattage
	}
	if req.MaxLightWattage != nil {
		if *req.MaxLightWattage < 0 {
			return nil, errors.New("max_light_wattage must be >= 0")
		}
		config.MaxLightWattage = *req.MaxLightWattage
	}
	if req.MinTemperature != nil {
		if *req.MinTemperature < 0 || *req.MinTemperature > 40 {
			return nil, errors.New("min_temperature must be between 0 and 40")
		}
		config.MinTemperature = *req.MinTemperature
	}
	if req.MaxTemperature != nil {
		if *req.MaxTemperature < 0 || *req.MaxTemperature > 40 {
			return nil, errors.New("max_temperature must be between 0 and 40")
		}
		config.MaxTemperature = *req.MaxTemperature
	}
	if req.PumpLevel != nil {
		if *req.PumpLevel < 1 || *req.PumpLevel > 5 {
			return nil, errors.New("pump_level must be between 1 and 5")
		}
		config.PumpLevel = *req.PumpLevel
	}
	if req.Description != nil {
		config.Description = *req.Description
	}

	if config.MinLightWattage > config.MaxLightWattage {
		return nil, errors.New("min_light_wattage must be <= max_light_wattage")
	}
	if config.MinTemperature > config.MaxTemperature {
		return nil, errors.New("min_temperature must be <= max_temperature")
	}

	if err := s.db.Save(config).Error; err != nil {
		return nil, err
	}

	return config, nil
}

func (s *OxygenService) DeleteConfig(id uint64) error {
	result := s.db.Delete(&OxygenConfig{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("config not found")
	}
	return nil
}

func (s *OxygenService) CalculateMatch(req *matchRequestDTO) (*MatchResult, error) {
	lightWattage := req.LightWattage
	temperature := req.Temperature

	if lightWattage < 0 {
		return nil, errors.New("light_wattage must be >= 0")
	}
	if temperature < 0 || temperature > 40 {
		return nil, errors.New("temperature must be between 0 and 40")
	}

	var configs []OxygenConfig
	if err := s.db.Find(&configs).Error; err != nil {
		return nil, err
	}

	var matchedConfig *OxygenConfig
	for i := range configs {
		cfg := &configs[i]
		if lightWattage >= cfg.MinLightWattage && lightWattage <= cfg.MaxLightWattage &&
			temperature >= cfg.MinTemperature && temperature <= cfg.MaxTemperature {
			matchedConfig = cfg
			break
		}
	}

	if matchedConfig == nil {
		level := s.calculateFallbackLevel(lightWattage, temperature)
		return &MatchResult{
			PumpLevel:   level,
			Description: "自动计算",
			Formula:     fmt.Sprintf("泵档位 = clamp(round((灯光瓦数/100 * 2) + (温度/40 * 3)), 1, 5)"),
			Reason:      fmt.Sprintf("无匹配配置，使用公式计算：灯光瓦数=%d, 温度=%.1f°C", lightWattage, temperature),
		}, nil
	}

	reason := fmt.Sprintf("匹配配置[%s]：灯光瓦数 %d 在 [%d-%d] 区间，温度 %.1f°C 在 [%.1f-%.1f]°C 区间",
		matchedConfig.Description,
		lightWattage, matchedConfig.MinLightWattage, matchedConfig.MaxLightWattage,
		temperature, matchedConfig.MinTemperature, matchedConfig.MaxTemperature)

	return &MatchResult{
		PumpLevel:   matchedConfig.PumpLevel,
		Description: matchedConfig.Description,
		Formula:     "查表匹配",
		Reason:      reason,
	}, nil
}

func (s *OxygenService) calculateFallbackLevel(lightWattage int, temperature float64) int {
	lightFactor := float64(lightWattage) / 100.0 * 2.0
	tempFactor := temperature / 40.0 * 3.0
	level := math.Round(lightFactor + tempFactor)
	if level < 1 {
		level = 1
	}
	if level > 5 {
		level = 5
	}
	return int(level)
}

func (s *OxygenService) GetMatchMatrix() (map[string]interface{}, error) {
	configs, err := s.ListConfigs()
	if err != nil {
		return nil, err
	}

	matrix := make(map[string]map[string]int)
	lightRanges := []string{"0-30W", "31-60W", "61-100W"}
	tempRanges := []string{"18-24°C", "24.1-28°C", "28.1-32°C"}

	for _, lr := range lightRanges {
		matrix[lr] = make(map[string]int)
		for _, tr := range tempRanges {
			matrix[lr][tr] = 0
		}
	}

	for _, cfg := range configs {
		lightKey := fmt.Sprintf("%d-%dW", cfg.MinLightWattage, cfg.MaxLightWattage)
		tempKey := fmt.Sprintf("%.1f-%.1f°C", cfg.MinTemperature, cfg.MaxTemperature)
		if matrix[lightKey] != nil {
			matrix[lightKey][tempKey] = cfg.PumpLevel
		}
	}

	return map[string]interface{}{
		"configs": configs,
		"matrix":  matrix,
		"light_ranges": lightRanges,
		"temp_ranges":  tempRanges,
	}, nil
}
