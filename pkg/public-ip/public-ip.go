package publicIp

import (
	cmdArgs "cmd-args"
	"config"

	"fmt"

	"github.com/oleksii-honchar/blablo"
	c "github.com/oleksii-honchar/coteco"
)

var f = fmt.Sprintf
var logger *blablo.Logger

func Init(nrpConfig *config.NrpConfig) bool {
	logger = blablo.NewLogger("public-ip", cmdArgs.LogLevel)
	logger.Debug("Init 'Public IP' tools")
	if nrpConfig.PublicIp.CheckAndUpdate != "yes" {
		logger.Debug(f("'Public IP' config is disabled. Skipping."))
		return true
	}

	logger.Info(c.WithGreen("'Public IP' tools configured successfuly"))
	return true
}

func GenerateCronTask(nrpConfig *config.NrpConfig) string {
	return "/usr/local/bin/nrp-cli -config=/etc/nrp.yaml -log-level=info -check-and-update-public-ip"
}
