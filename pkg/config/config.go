package config

import (
	cmdArgs "cmd-args"
	cd "config-defaults"

	"fmt"
	"os"

	"github.com/coryb/figtree"

	"github.com/oleksii-honchar/blablo"
	c "github.com/oleksii-honchar/coteco"

	"gopkg.in/yaml.v3"
)

var f = fmt.Sprintf
var logger *blablo.Logger

var nrpConfig *NrpConfig
var nrpSvcDefaultConfig *NrpServiceConfig

// Should be called first before any other pkg function calls
func Init() (*NrpConfig, error) {
	logger = blablo.NewLogger("config", cmdArgs.LogLevel)
	logger.Debug("Init 'Config'")

	logger.Debug("Loading config and merging with defaults")

	var err error
	var nrpBaseConfig *NrpConfig
	nrpConfig, err = loadNrpConfig(cmdArgs.ConfigPath)
	if err != nil {
		return nil, err
	}

	nrpBaseConfig, err = loadNrpDefaultsConfig(cmdArgs.DefaultsMode)
	if err != nil {
		return nil, err
	}

	nrpSvcDefaultConfig, err = loadNrpSvcDefaultConfig()
	if err != nil {
		return nil, err
	}

	// kinda RIGHT JOIN, i.e adding to the left param missing parts from right param
	figtree.Merge(nrpConfig, nrpBaseConfig)
	logger.Debug(c.WithGreen("Base configs merged successfuly"))

	err = mergeNrpSvcDefaultsConfig(nrpConfig, nrpSvcDefaultConfig)
	if err != nil {
		return nil, err
	}

	logger.Debug("Init completed for 'Config'")

	return nrpConfig, nil
}

func loadNrpSvcDefaultConfig() (*NrpServiceConfig, error) {
	defaultsContent := cd.NrpSvcConfigDefaults

	logger.Debug(f("Parsing (%s) defaults config", c.WithCyan("service")))

	var nrpSvcDefaultConfig NrpServiceConfig
	err := yaml.Unmarshal(defaultsContent, &nrpSvcDefaultConfig)
	if err != nil {
		logger.Error(f("Failed to parse service defaults config content: %s", c.WithRed(err.Error())))
		return nil, err
	}
	return &nrpSvcDefaultConfig, nil
}

func loadNrpDefaultsConfig(defaultsMode string) (*NrpConfig, error) {
	var defaultsContent []byte
	if defaultsMode == cd.DefaultsDevMode {
		defaultsContent = cd.NrpConfigDevDefaults
	} else if defaultsMode == cd.DefaultsProdMode {
		defaultsContent = cd.NrpConfigProdDefaults
	}

	logger.Debug(f("Parsing [%s] defaults config", c.WithCyan(defaultsMode)))

	var nrpDefaultsConfig NrpConfig
	err := yaml.Unmarshal(defaultsContent, &nrpDefaultsConfig)
	if err != nil {
		logger.Error(f("Failed to parse defaults config content: %s", c.WithRed(err.Error())))
		return nil, err
	}
	return &nrpDefaultsConfig, nil
}

func loadNrpConfig(configPath string) (*NrpConfig, error) {
	file, err := os.Open(configPath)
	if err != nil {
		logger.Error(f("Failed to open config file: %s", c.WithRed(err.Error())))
		return nil, err
	}
	defer file.Close()

	var config NrpConfig
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		logger.Error(f("Failed to parse config file: %s", c.WithRed(err.Error())))
		return nil, err
	}

	logger.Debug(f("Loaded %s", c.WithCyan("nrp.yaml")))
	return &config, nil
}

// now let's loop services and merge with default svc cofig
func mergeNrpSvcDefaultsConfig(
	nrpConfig *NrpConfig,
	nrpSvcDefaultConfig *NrpServiceConfig,
) error {
	logger.Debug("Going to merge service defaults with each service")
	for idx := range nrpConfig.Services {
		figtree.Merge(&nrpConfig.Services[idx], nrpSvcDefaultConfig)
	}

	logger.Debug(f("Completed merging default config for (%s) services ", c.WithGreen(fmt.Sprint(len(nrpConfig.Services)))))
	return nil // kinda weird ...
}
