package configProcessor

import (
	"fmt"
	"os"

	"beaver/blablo"
	"beaver/blablo/color"

	"gopkg.in/yaml.v2"
)

type NrpServiceConfig struct {
	Name          string `yaml:"name"`
	ServiceIP     string `yaml:"serviceIp"`
	ServicePort   int    `yaml:"servicePort"`
	DomainName    string `yaml:"domainName"`
	CORS          bool   `yaml:"cors,omitempty"`
	BlockExploits bool   `yaml:"blockExploits,omitempty"`
	HTTPS         struct {
		Force bool `yaml:"force,omitempty"`
		HSTS  bool `yaml:"hsts,omitempty"`
	} `yaml:"https,omitempty"`
}
type NrpConfig struct {
	Services []NrpServiceConfig `yaml:"services"`
}

func LoadBaseConfig(configPath string) (*NrpConfig, error) {
	logger := blablo.NewLogger("cfg-prcsr")
	logger.Info("Loading nrp.yaml")

	file, err := os.Open(configPath)
	if err != nil {
		logger.Error("Failed to open config file:", "err", err)
		return nil, err
	}
	defer file.Close()

	var config NrpConfig
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		logger.Error("Failed to parse config file:", "err", err)
		return nil, err
	}

	logger.Info(fmt.Sprintf("Found %s(%v)%s services configuration", color.Green, len(config.Services), color.Reset))
	return &config, nil
}

func GenerateNginxServerConfig(svcConfig *NrpServiceConfig) (string, error) {
	return "pending", nil
}
