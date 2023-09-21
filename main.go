package main

import (
	configProcessor "config-processor"
	"fmt"

	"github.com/oleksii-honchar/blablo"
	c "github.com/oleksii-honchar/coteco"
)

var f = fmt.Sprintf

func main() {

	logger := blablo.NewLogger("main")
	logger.Info(c.WithGreenCyan49("'Nginx reverse Proxy' cli tool v0.1"))

	nrpConfig, err := configProcessor.Init()
	if err != nil {
		return
	}

	logger.Info("Generating nginx configs")

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

	logger.Info(c.WithGreenCyan49("Done ✨"))
}
