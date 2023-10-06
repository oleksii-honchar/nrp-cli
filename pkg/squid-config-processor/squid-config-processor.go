package squidConfigProcessor

import (
	"bytes"
	cmdArgs "cmd-args"
	"config"
	"fmt"
	"os"

	_ "embed"

	"text/template"

	"github.com/oleksii-honchar/blablo"
	c "github.com/oleksii-honchar/coteco"
)

//go:embed squid.conf.tmpl
var SquidConfTmpl []byte

var f = fmt.Sprintf
var logger *blablo.Logger

func GenerateConfig(config *config.NrpConfig) bool {
	logger = blablo.NewLogger("squid-cfg", string(cmdArgs.LogLevel))
	logger.Debug(f("Processing Squid config: %s%+v%s", c.Yellow, config.Squid, c.Reset))
	if config.Squid.Use != "yes" {
		logger.Debug(f("Squid config is disabled. Skipping config generation."))
		return true
	}

	confTemplate, err := template.New("squid-config").Parse(string(SquidConfTmpl))
	if err != nil {
		logger.Error(f("Error creating squid config tmpl: %s", c.WithRed(err.Error())))
		return false
	}

	var content bytes.Buffer
	err = confTemplate.Execute(&content, config.Squid)
	if err != nil {
		logger.Error(f("Failed to generate squid config: %s", c.WithRed(err.Error())))
		return false
	}
	// logger.Debug(f("Generated (%s) bytes of config data", c.WithGreen(fmt.Sprint(content.Len()))))

	if err := os.WriteFile(config.Squid.ConfigPath, content.Bytes(), 0644); err != nil {
		logger.Error(f("Saving content to file (%s): %s", c.WithCyan(config.Squid.ConfigPath), c.WithRed(err.Error())))
		return false
	} else {
		logger.Debug(f("Saved (%s) bytes to file: %s", c.WithCyan(f("%v", content.Len())), c.WithGreen(config.Squid.ConfigPath)))
	}

	logger.Info(c.WithGreen(f("Config generation completed for %s", c.WithCyan("'squid'"))))

	return true
}
