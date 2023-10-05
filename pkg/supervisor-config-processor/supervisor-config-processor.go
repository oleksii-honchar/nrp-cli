package supervisorConfigProcessor

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

//go:embed supervisord.conf.tmpl
var SupervisorConfTmpl []byte

var f = fmt.Sprintf
var logger *blablo.Logger

func GenerateConfig(config *cfgProc.NrpConfig) bool {
	logger = blablo.NewLogger("supvsr-cfg", string(cmdArgs.LogLevel))
	logger.Debug(f("Processing 'supervisor' config %s%+v%s", c.Yellow, config.Supervisor, c.Reset))

	confTemplate, err := template.New("supervisor-config").Parse(string(SupervisorConfTmpl))
	if err != nil {
		logger.Error(f("Error creating 'supervisor' config tmpl: %s", c.WithRed(err.Error())))
		return false
	}

	var content bytes.Buffer
	err = confTemplate.Execute(&content, config)
	if err != nil {
		logger.Error(f("Failed to generate 'supervisor' config: %s", c.WithRed(err.Error())))
		return false
	}
	// logger.Debug(f("Generated (%s) bytes of config data", c.WithGreen(fmt.Sprint(content.Len()))))

	filePath := filepath.Join(config.Supervisor.ConfigPath, "supervisord.conf")
	if err := os.WriteFile(filePath, content.Bytes(), 0644); err != nil {
		logger.Error(f("Saving content to file (%s): %s", c.WithCyan(filePath), c.WithRed(err.Error())))
		return false
	} else {
		logger.Debug(f("Saved (%s) bytes to file: %s", c.WithCyan(f("%v", content.Len())), c.WithGreen(filePath)))
	}

	logger.Info(c.WithGreen(f("Supervisor config generation completed")))

	return true
}
