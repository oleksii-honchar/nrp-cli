package squidConfigProcessor

import (
	"bytes"
	cmdArgs "cmd-args"
	cfgProc "config-processor"
	"fmt"
	"os"
	"path/filepath"

	_ "embed"

	"text/template"

	"github.com/oleksii-honchar/blablo"
	c "github.com/oleksii-honchar/coteco"
)

//go:embed squid.conf.tmpl
var SquidConfTmpl []byte

var f = fmt.Sprintf
var logger *blablo.Logger

func GenerateSquidConfig(config *cfgProc.NrpConfig) bool {
	logger = blablo.NewLogger("squid-cfg", string(cmdArgs.LogLevel))
	logger.Debug(f("Processing Squid config: %s%+v%s", c.Yellow, config.Squid, c.Reset))
	if !config.Squid.Use {
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

	filePath := filepath.Join(config.Squid.ConfigPath, "squid.conf")
	if err := os.WriteFile(filePath, content.Bytes(), 0644); err != nil {
		logger.Error(f("Saving content to file (%s): %s", c.WithCyan(filePath), c.WithRed(err.Error())))
		return false
	} else {
		logger.Debug(f("Saved (%s) bytes to file: %s", c.WithCyan(f("%v", content.Len())), c.WithGreen(filePath)))
	}

	logger.Info(c.WithGreen(f("Squid config generation completed")))

	return true
}
