package main

import (
	"beaver/blablo"
	c "beaver/blablo/color"
	"fmt"
	configProcessor "nrp/config-processor"
	"os"
	"path/filepath"
)

var f = fmt.Sprintf

func main() {

	logger := blablo.NewLogger("main")
	logger.Info(c.WithGreenCyan49("'Nginx reverse Proxy' cli tool v0.1"))

	ok, _ := configProcessor.Init("./nginx-config")
	if !ok {
		return
	}

	nrpConfig, err := configProcessor.LoadBaseConfig("./configs/nrp.yaml")
	if err != nil {
		return
	}

	logger.Info("Generating nginx configs")

	content, err := configProcessor.GenerateDefaultNginxConfig()
	if err != nil {
		return
	}
	filePath := filepath.Join(".", "nginx-config/conf.available", f("%v-%s.conf", 0, "default"))
	if err := os.WriteFile(filePath, content.Bytes(), 0644); err != nil {
		logger.Error(f("Saving content to file: %s", c.WithCyan(filePath)))
	} else {
		logger.Info(f("Saved (%s) bytes to file: %s", c.WithCyan(f("%v", content.Len())), c.WithGreen(filePath)))
	}

	for idx, svcCfg := range nrpConfig.Services {
		logger.Info(f("Processing %s for service: %s",
			c.WithCyan(f("[%v/%v]", idx+1, len(nrpConfig.Services))),
			c.WithCyan(svcCfg.Name)),
		)

		content, err := configProcessor.GenerateNginxServerConfig(&svcCfg)
		if err != nil {
			continue
		}

		filePath := filepath.Join(".", "nginx-config/conf.available", f("%v-%s.conf", idx+1, svcCfg.Name))
		if err := os.WriteFile(filePath, content.Bytes(), 0644); err != nil {
			logger.Error(f("Saving content to file: %s", c.WithCyan(filePath)))
		} else {
			logger.Info(f("Saved (%s) bytes to file: %s", c.WithCyan(f("%v", content.Len())), c.WithGreen(filePath)))
		}
	}

	logger.Info(c.WithGreenCyan49("Done ✨"))
}
