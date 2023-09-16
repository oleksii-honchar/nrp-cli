package main

import (
	"beaver/blablo"
	c "beaver/blablo/color"
	"fmt"
	configProcessor "nrp/config-processor"
	"path/filepath"
)

func main() {
	f := fmt.Sprintf
	logger := blablo.NewLogger("nrp-cli")
	logger.Info(c.With(c.GreenCyan49, "'Nginx reverse Proxy' cli tool v0.1"))

	nrpConfig, err := configProcessor.LoadBaseConfig(filepath.Join(".", "configs", "nrp.yaml"))
	if err != nil {
		return
	}

	logger.Info("Generating nginx configs")
	for idx, svcCfg := range nrpConfig.Services {
		// nginxConfig, err := configProcessor.GenerateNginxConfig(svcCfg)
		logger.Info(f("Processing %s for service: %s",
			c.With(c.Cyan, f("[%v/%v]", idx+1, len(nrpConfig.Services))),
			c.With(c.Cyan, svcCfg.Name)),
		)

		content, _ := configProcessor.GenerateNginxServerConfig(&svcCfg)
		logger.Info("Saving content to file", "content", content)
	}

	logger.Info(c.With(c.Green, "Done ✨"))
}
