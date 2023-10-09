package cronConfigProcessor

import (
	"bytes"
	cmdArgs "cmd-args"
	"config"
	"fmt"
	"os"
	publicIp "public-ip"

	_ "embed"

	"text/template"

	"github.com/oleksii-honchar/blablo"
	c "github.com/oleksii-honchar/coteco"
)

//go:embed crontab.tmpl
var CronConfTmpl []byte

var f = fmt.Sprintf
var logger *blablo.Logger

func GenerateConfig(nrpConfig *config.NrpConfig) bool {
	logger = blablo.NewLogger("cron-cfg", string(cmdArgs.LogLevel))
	// logger.Debug(f("Processing 'cron' config %s%+v%s", c.Yellow, config.Cron, c.Reset))
	logger.Debug(f("Processing 'cron' config"))

	confTemplate, err := template.New("cron-config").Parse(string(CronConfTmpl))
	if err != nil {
		logger.Error(f("Error creating 'cron' config tmpl: %s", c.WithRed(err.Error())))
		return false
	}

	var publicIpCheckAndUpdateTask = publicIp.GenerateCronTask(nrpConfig)
	var cronSchedule = generateCronSchedule(nrpConfig.PublicIp.Schedule)
	var tmplData = struct {
		CheckAndUpdate string
		CronSchedule   string
		Schedule       string
		Task           string
	}{
		CheckAndUpdate: nrpConfig.PublicIp.CheckAndUpdate,
		CronSchedule:   cronSchedule,
		Schedule:       nrpConfig.PublicIp.Schedule,
		Task:           publicIpCheckAndUpdateTask,
	}
	var content bytes.Buffer
	err = confTemplate.Execute(&content, tmplData)
	if err != nil {
		logger.Error(f("Failed to generate 'cron' config: %s", c.WithRed(err.Error())))
		return false
	}
	// logger.Debug(f("Generated (%s) bytes of config data", c.WithGreen(fmt.Sprint(content.Len()))))

	if err := os.WriteFile(nrpConfig.Cron.ConfigPath, content.Bytes(), 0644); err != nil {
		logger.Error(f("Saving content to file (%s): %s", c.WithCyan(nrpConfig.Cron.ConfigPath), c.WithRed(err.Error())))
		return false
	} else {
		logger.Debug(f("Saved (%s) bytes to file: %s", c.WithCyan(f("%v", content.Len())), c.WithGreen(nrpConfig.Cron.ConfigPath)))
	}

	logger.Info(c.WithGreen(f("Config generation completed for %s", c.WithCyan("'cron'"))))

	return true
}

func generateCronSchedule(schedule string) string {
	switch schedule {
	case "1m":
		return "0 0 1 * *"
	case "1w":
		return "0 0 * * 0"
	case "1d":
		return "0 0 * * *"
	case "1h":
		return "0 */1 * * *"
	case "1min":
		return "* * * * *"
	default:
		logger.Warn(f("Unknown cron schedule: %s. Using default: '1h'", c.WithRed(schedule)))
		return "* * * * *"
	}
}
