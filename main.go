package main

import (
	configProcessor "config-processor"
	"fmt"

	cmdArgs "cmd-args"

	lv "latest-version"

	squidCfgProc "squid-config-processor"

	"github.com/oleksii-honchar/blablo"
	c "github.com/oleksii-honchar/coteco"
)

var f = fmt.Sprintf

func main() {
	if ok := cmdArgs.Init(); !ok {
		return
	}

	logger := blablo.NewLogger("main", cmdArgs.LogLevel)
	logger.Info(c.WithGreenCyan49(f("'Nginx reverse Proxy' cli tool %s", c.WithCyan(lv.LatestVersion))))

	nrpConfig, err := configProcessor.Init()
	if err != nil {
		return
	}

	logger.Info(f("Generating nginx configs in '%s'", c.WithCyan("conf.available")))

	configProcessor.CreateDeafultConfFile()

	// Process array of services
	for idx, svcCfg := range nrpConfig.Services {
		transportMode := c.WithGray247("[HTTP]")
		if svcCfg.HTTPS.Use {
			transportMode = c.WithOrange("[HTTPS]")
		}

		logger.Info(f("%s processing service: %s %s",
			c.WithCyan(f("[%v/%v]", idx+1, len(nrpConfig.Services))),
			c.WithCyan(svcCfg.Name),
			transportMode),
		)

		//  Check/create certificates if HTTPS.Use = true
		if svcCfg.HTTPS.Use {
			if ok := configProcessor.CheckCertificateFiles(svcCfg.Name); !ok {
				// need to create enw certs
				if ok := configProcessor.CreateCertificateFiles(&svcCfg); !ok {
					// something wrong with Letsencrypt certbot processing - turning off https
					svcCfg.HTTPS.Use = false
					logger.Info(f("HTTPS turned %s for service: %s", c.WithRed("off"), c.WithCyan(svcCfg.Name)))
				}
			}
			// if certs are in place or https turned off - continue to geneate nginx server config
		}

		configProcessor.CreateServiceConfFile(idx, &svcCfg)
	}

	configProcessor.CopyConfFiles()

	_ = squidCfgProc.GenerateSquidConfig(nrpConfig)

	logger.Info(c.WithGreenCyan49("Done ✨"))
}
