package main

import (
	"os"
	"path/filepath"

	"beaver/blablo"
	"beaver/blablo/color"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Configurations []struct {
		Name          string `yaml:"name"`
		ServerIP      string `yaml:"serverIp"`
		ServerPort    int    `yaml:"serverPort"`
		DomainName    string `yaml:"domainName"`
		CORS          bool   `yaml:"cors,omitempty"`
		BlockExploits bool   `yaml:"blockExploits,omitempty"`
		HTTPS         struct {
			Force bool `yaml:"force,omitempty"`
			HSTS  bool `yaml:"hsts,omitempty"`
		} `yaml:"https,omitempty"`
	} `yaml:"configurations"`
}

func main() {

	logger := blablo.NewLogger()
	logger.Info(color.GreenCyan49 + "'Nginx reverse Proxy' cli tool v0.1" + color.Reset)

	filePath := filepath.Join(".", "configs", "nrp.yaml")
	file, err := os.Open(filePath)
	if err != nil {
		logger.Error("Failed to open config file:", "err", err)
		return
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		logger.Error("Failed to parse config file:", "err", err)
		return
	}

	// Print the parsed configuration
	logger.Info("Parsed configuration:", "config", config)
}
