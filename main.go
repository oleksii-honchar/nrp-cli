package main

import (
	"fmt"

	"github.com/oleksii-honchar/blablo"
	c "github.com/oleksii-honchar/coteco"

	cmdArgs "cmd-args"
	lv "latest-version"

	"config"
	cronCfgProc "cron-config-processor"
	dnsmasqCfgProc "dnsmasq-config-processor"
	nginxCfgProcessor "nginx-config-processor"
	publicIp "public-ip"
	squidCfgProc "squid-config-processor"
	supervisorCfgProc "supervisor-config-processor"
)

var f = fmt.Sprintf

func main() {
	if ok := cmdArgs.Init(); !ok {
		return
	}

	logger := blablo.NewLogger("main", cmdArgs.LogLevel)
	logger.Info(c.WithGreenCyan49(f("'Nginx Reverse Proxy' cli tool %s", c.WithCyan(lv.LatestVersion))))

	nrpConfig, err := config.Init()
	if err != nil {
		return
	}

	if ok := publicIp.Init(nrpConfig); !ok {
		return
	}
	if cmdArgs.CheckAndUpdatePublicIp {
		return
	}

	if ok := nginxCfgProcessor.Init(nrpConfig); !ok {
		return
	}

	if ok := nginxCfgProcessor.CreateDefaultConfFile(nrpConfig); !ok {
		return
	}

	// Process array of services
	for idx, svcCfg := range nrpConfig.Services {
		transportMode := c.WithGray247("[HTTP]")
		if svcCfg.HTTPS.Use == "yes" {
			transportMode = c.WithOrange("[HTTPS]")
		}

		logger.Info(f("%s processing service: %s %s",
			c.WithCyan(f("[%v/%v]", idx+1, len(nrpConfig.Services))),
			c.WithCyan(svcCfg.Name),
			transportMode),
		)

		//  Check/create certificates if HTTPS.Use = true
		if svcCfg.HTTPS.Use == "yes" {
			if ok := nginxCfgProcessor.CheckCertificateFiles(nrpConfig, svcCfg.Name); !ok {
				// need to create enw certs
				if ok := nginxCfgProcessor.CreateCertificateFiles(nrpConfig, &svcCfg); !ok {
					// something wrong with Letsencrypt certbot processing - turning off https
					svcCfg.HTTPS.Use = "no"
					logger.Info(f("HTTPS turned %s for service: %s", c.WithRed("off"), c.WithCyan(svcCfg.Name)))
				}
			}
			// if certs are in place or https turned off - continue to geneate nginx server config
		}

		nginxCfgProcessor.CreateServiceConfFile(nrpConfig, idx, &svcCfg)
	}

	nginxCfgProcessor.CopyConfFiles(nrpConfig)
	logger.Info(c.WithGreen(f("Config generation completed for %s", c.WithCyan("'nginx'"))))

	_ = squidCfgProc.GenerateConfig(nrpConfig)
	_ = dnsmasqCfgProc.GenerateConfig(nrpConfig)
	_ = supervisorCfgProc.GenerateConfig(nrpConfig)
	_ = cronCfgProc.GenerateConfig(nrpConfig)

	logger.Info(c.WithGreenCyan49("Done ✨"))
}
