package config

import (
	"bytes"
	cmdArgs "cmd-args"
	cd "config-defaults"
	"regexp"

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
	logger = blablo.NewLogger("config", cmdArgs.LogLevel, false)
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
	var contentWithValues []byte
	content, err := os.ReadFile(configPath)

	if err != nil {
		logger.Error(f("Failed to open config file: %s", c.WithRed(err.Error())))
		return nil, err
	}

	contentWithValues, err = replaceEnvVarsWithValues(content)
	if err != nil {
		return nil, err
	}

	var config NrpConfig
	decoder := yaml.NewDecoder(bytes.NewReader(contentWithValues))
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

/*
Yaml keys values can be set to environment variable name, e.g. $SERVICE_HOST
So here we take yaml string content and replacing via regex env vars name ina form of "$ENV_VAR" to it's possible value. If no value exists, set "" - empty string
*/
func replaceEnvVarsWithValues(content []byte) ([]byte, error) {
	logger.Debug("Going to replace env vars with their values if any")
	var hasErrors bool = false
	re := regexp.MustCompile(`\$(\w+)`)

	// Replace environment variable names with their values
	result := re.ReplaceAllStringFunc(string(content), func(match string) string {
		envVar := match[1:] // Remove the leading "$" of the environment variable name
		value := os.Getenv(envVar)
		if value == "" {
			hasErrors = true
			logger.Error(f(
				"Can't get value for env var mentioned in config: %s",
				c.WithRed(envVar),
			))
		}
		return value
	})

	if !hasErrors {
		logger.Debug(f(c.WithGreen("Completed env vars values replacement")))
	} else {
		logger.Info(f(c.WithRed("Some environment variables values missing, please check your setup")))
	}

	return []byte(result), nil
}
