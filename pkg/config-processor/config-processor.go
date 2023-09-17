package configProcessor

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"beaver/blablo"
	c "beaver/blablo/color"

	"gopkg.in/yaml.v2"
)

var f = fmt.Sprintf
var logger *blablo.Logger

var nginxConfigBaseFolder string

type NrpServiceConfig struct {
	Name          string `yaml:"name"`
	ServiceIP     string `yaml:"serviceIp"`
	ServicePort   int    `yaml:"servicePort"`
	DomainName    string `yaml:"domainName"`
	CORS          bool   `yaml:"cors,omitempty"`
	BlockExploits bool   `yaml:"blockExploits,omitempty"`
	HTTPS         struct {
		Use   bool `yaml:"use,omitempty"`
		Force bool `yaml:"force,omitempty"`
		HSTS  bool `yaml:"hsts,omitempty"`
	} `yaml:"https,omitempty"`
}
type NrpConfig struct {
	Services []NrpServiceConfig `yaml:"services"`
}

func Init(nginxConfigPath string) (bool, error) {
	nginxConfigBaseFolder = nginxConfigPath
	logger = blablo.NewLogger("cfg-prcsr")
	logger.Info("Init 'Config Processor'")

	confAvailablePath := filepath.Join(nginxConfigBaseFolder, "conf.available")
	if err := os.RemoveAll(confAvailablePath); err != nil {
		logger.Error(f("Failed to clean folder:", c.WithCyan(confAvailablePath)), "err", err)
	}

	if err := os.Mkdir(confAvailablePath, os.ModePerm); err != nil {
		logger.Error(f("Failed to re-create folder:", c.WithCyan(confAvailablePath)), "err", err)
	}
	logger.Info(f("Folder cleaned: %s", c.WithGreen(confAvailablePath)))

	svcConfTmplPath := filepath.Join(nginxConfigBaseFolder, "/templates/service.conf.tmpl")
	res1, _ := loadConfTemplate(svcConfTmplPath)

	defaultSvcConfTmplPath := filepath.Join(nginxConfigBaseFolder, "/templates/default.conf.tmpl")
	res2, _ := loadDefaultConfTemplate(defaultSvcConfTmplPath)

	logger.Info("Init completed")

	return res1 && res2, nil
}

func LoadBaseConfig(configPath string) (*NrpConfig, error) {
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

	logger.Info(f("Loaded %s", c.WithCyan("nrp.yaml")))
	logger.Info(f("Found (%s) services configuration", c.WithGreen(fmt.Sprint(len(config.Services)))))
	return &config, nil
}

func GenerateDefaultNginxConfig() (*bytes.Buffer, error) {
	var content bytes.Buffer
	err := defaultConfTemplate.Execute(&content, nil)
	if err != nil {
		logger.Error(f("Failed to generate nginx config for service: %s", c.WithCyan("default")), "err", err)
		return nil, err
	}
	logger.Info(f("Generated (%s) bytes of config data", c.WithGreen(fmt.Sprint(content.Len()))))
	return &content, nil
}

func GenerateNginxServerConfig(svcConfig *NrpServiceConfig) (*bytes.Buffer, error) {
	var content bytes.Buffer
	err := confTemplate.Execute(&content, svcConfig)
	if err != nil {
		logger.Error(f("Failed to generate nginx config for service: %s", c.WithCyan(svcConfig.Name)), "err", err)
		return nil, err
	}
	logger.Info(f("Generated (%s) bytes of config data", c.WithGreen(fmt.Sprint(content.Len()))))
	return &content, nil
}
