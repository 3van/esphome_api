package types

import (
	"time"
)

type Config struct {
	Active  string         `yaml:"active"`
	Devices []DeviceConfig `yaml:"devices"`
}

type DeviceInfo struct {
	Name            string    `yaml:"name"`
	Model           string    `yaml:"model"`
	MacAddress      string    `yaml:"macAddress"`
	EsphomeVersion  string    `yaml:"esphomeVersion"`
	CompilationTime string    `yaml:"compilationTime"`
	HasDeepSleep    bool      `yaml:"hasDeepSleep"`
	StatusOn        time.Time `yaml:"statusOn"`
}

func (di *DeviceInfo) Clone() DeviceInfo {
	return DeviceInfo{
		Name:            di.Name,
		Model:           di.Model,
		MacAddress:      di.MacAddress,
		EsphomeVersion:  di.EsphomeVersion,
		CompilationTime: di.CompilationTime,
		HasDeepSleep:    di.HasDeepSleep,
	}
}

type DeviceConfig struct {
	Address       string        `yaml:"address"`
	EncryptionKey string        `yaml:"encryptionKey"`
	Timeout       time.Duration `yaml:"timeout"`
	Info          DeviceInfo    `yaml:"info"`
}

func (dc *DeviceConfig) Clone() DeviceConfig {
	return DeviceConfig{
		Address:       dc.Address,
		EncryptionKey: dc.EncryptionKey,
		Timeout:       dc.Timeout,
		Info:          dc.Info.Clone(),
	}
}
