package common

import (
	"strings"
	"time"
)

var (
	DeviceTypeMap = map[string]string{
		"light":      "light",
		"灯光":        "light",
		"主灯":        "light",
		"辅灯":        "light",
		"pump":       "pump",
		"气泵":        "pump",
		"主气泵":      "pump",
		"氧气泵":      "pump",
	}

	DeviceNameMap = map[string]string{
		"主灯":   "主灯",
		"main": "主灯",
		"辅灯":   "辅灯",
		"sub":  "辅灯",
		"主气泵": "主气泵",
		"main_pump": "主气泵",
	}

	ModeMap = map[string]bool{
		"manual": true,
		"手动":   true,
		"手动模式": true,
		"auto":   false,
		"自动":   false,
		"自动模式": false,
	}

	StatusMap = map[string]bool{
		"on":     true,
		"开":      true,
		"开启":     true,
		"打开":     true,
		"true":   true,
		"off":    false,
		"关":      false,
		"关闭":     false,
		"false":  false,
	}

	BrightnessMap = map[string]int{
		"低":   30,
		"弱":   30,
		"低光":  30,
		"中":   60,
		"中等": 60,
		"中光":  60,
		"高":   100,
		"强":   100,
		"高光":  100,
		"最高": 100,
	}

	PumpLevelMap = map[string]int{
		"1档": 1,
		"一档": 1,
		"最低": 1,
		"2档": 2,
		"二档": 2,
		"低":   2,
		"3档": 3,
		"三档": 3,
		"中":   3,
		"中等": 3,
		"4档": 4,
		"四档": 4,
		"高":   4,
		"5档": 5,
		"五档": 5,
		"最高": 5,
	}

	TimeRangeMap = map[string][2]string{
		"早上":    {"06:00:00", "09:00:00"},
		"上午":    {"09:00:00", "12:00:00"},
		"中午":    {"12:00:00", "14:00:00"},
		"下午":    {"14:00:00", "18:00:00"},
		"傍晚":    {"17:00:00", "19:00:00"},
		"晚上":    {"19:00:00", "22:00:00"},
		"夜间":    {"22:00:00", "06:00:00"},
		"全天":    {"00:00:00", "23:59:59"},
	}

	NightModeStartHour = 0
	NightModeEndHour   = 5

	SafeTempMin = 22.0
	SafeTempMax = 28.0

	PowerSavingFactor = 0.5
)

func NormalizeDeviceType(input string) string {
	if input == "" {
		return ""
	}
	key := strings.ToLower(strings.TrimSpace(input))
	if val, ok := DeviceTypeMap[key]; ok {
		return val
	}
	if val, ok := DeviceTypeMap[input]; ok {
		return val
	}
	return input
}

func NormalizeDeviceName(input string) string {
	if input == "" {
		return ""
	}
	key := strings.ToLower(strings.TrimSpace(input))
	if val, ok := DeviceNameMap[key]; ok {
		return val
	}
	return input
}

func ParseMode(input string) (bool, bool) {
	if input == "" {
		return false, false
	}
	key := strings.ToLower(strings.TrimSpace(input))
	if val, ok := ModeMap[key]; ok {
		return val, true
	}
	if val, ok := ModeMap[input]; ok {
		return val, true
	}
	return false, false
}

func ParseStatus(input string) (bool, bool) {
	if input == "" {
		return false, false
	}
	key := strings.ToLower(strings.TrimSpace(input))
	if val, ok := StatusMap[key]; ok {
		return val, true
	}
	if val, ok := StatusMap[input]; ok {
		return val, true
	}
	return false, false
}

func ParseBrightness(input string) (int, bool) {
	if input == "" {
		return 0, false
	}
	key := strings.TrimSpace(input)
	if val, ok := BrightnessMap[key]; ok {
		return val, true
	}
	return 0, false
}

func ParsePumpLevel(input string) (int, bool) {
	if input == "" {
		return 0, false
	}
	key := strings.TrimSpace(input)
	if val, ok := PumpLevelMap[key]; ok {
		return val, true
	}
	return 0, false
}

func ParseTimeRange(input string) (string, string, bool) {
	if input == "" {
		return "", "", false
	}
	key := strings.TrimSpace(input)
	if val, ok := TimeRangeMap[key]; ok {
		return val[0], val[1], true
	}
	return "", "", false
}

func IsNightMode() bool {
	now := time.Now()
	hour := now.Hour()
	return hour >= NightModeStartHour && hour < NightModeEndHour
}

func IsNightModeAt(t time.Time) bool {
	hour := t.Hour()
	return hour >= NightModeStartHour && hour < NightModeEndHour
}

func IsTemperatureSafe(temperature float64) bool {
	return temperature >= SafeTempMin && temperature <= SafeTempMax
}

func ShouldEnablePowerSaving(temperature float64) bool {
	return IsNightMode() && IsTemperatureSafe(temperature)
}

func ShouldEnablePowerSavingAt(t time.Time, temperature float64) bool {
	return IsNightModeAt(t) && IsTemperatureSafe(temperature)
}

func ApplyPowerSaving(value int) int {
	return int(float64(value) * PowerSavingFactor)
}

func ApplyPowerSavingWithMin(value int, minValue int) int {
	result := int(float64(value) * PowerSavingFactor)
	if result < minValue {
		return minValue
	}
	return result
}
