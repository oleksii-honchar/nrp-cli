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

// Should be called first before any other pkg function calls
func Init() (*NrpConfig, error) {
	logger = blablo.NewLogger("cfg-prcsr", cmdArgs.LogLevel)
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

	figtree.Merge(nrpConfig, nrpBaseConfig)
	// figtree.Merge(nrpBaseConfig, nrpConfig)
	logger.Debug(c.WithGreen("Configs merged successfuly"))

	logger.Debug("Init completed for 'Config'")

	return nrpConfig, nil
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
		logger.Error("Failed to parse defaults config content:", "err", err.Error())
		return nil, err
	}
	return &nrpDefaultsConfig, nil
}

func loadNrpConfig(configPath string) (*NrpConfig, error) {
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
		logger.Error("Failed to parse config file:", "err", err.Error())
		return nil, err
	}

	logger.Debug(f("Loaded %s", c.WithCyan("nrp.yaml")))
	return &config, nil
}
