package nginxConfigProcessor

import (
	"bytes"
	"config"

	"os"
	"path/filepath"

	c "github.com/oleksii-honchar/coteco"
)

func generateCorsServersConfig(corsEnabledDomains []string) (*bytes.Buffer, error) {
	var content bytes.Buffer
	err := corsServersTemplate.Execute(&content, corsEnabledDomains)
	if err != nil {
		logger.Error(f("Failed to generate CORS servers config: %s", c.WithRed(err.Error())))
		return nil, err
	}
	return &content, nil
}

func CreateCorsServersConfFile(nrpConfig *config.NrpConfig, corsEnabledDomains []string) bool {
	logger.Debug(f("Generating CORS servers configs in '%s'", c.WithCyan("cors-servers.conf")))
	content, err := generateCorsServersConfig(corsEnabledDomains)
	if err != nil {
		return false
	}
	filePath := filepath.Join(nrpConfig.Nginx.ConfigPath, "cors-servers.conf")
	if err := os.WriteFile(filePath, content.Bytes(), 0644); err != nil {
		logger.Error(f("Saving content to file: %s", c.WithCyan(filePath)))
		return false
	} else {
		logger.Debug(f("Saved (%s) bytes to file: %s", c.WithCyan(f("%v", content.Len())), c.WithGreen(filePath)))
	}
	logger.Info(c.WithGreen(f("Config generation completed for %s", c.WithCyan("'CORS servers'"))))

	return true
}
