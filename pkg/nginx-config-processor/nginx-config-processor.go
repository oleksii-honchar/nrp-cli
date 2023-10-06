package nginxConfigProcessor

import (
	"bytes"
	cmdArgs "cmd-args"
	"config"

	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/oleksii-honchar/blablo"
	c "github.com/oleksii-honchar/coteco"
)

var f = fmt.Sprintf
var logger *blablo.Logger

func Init(nrpConfig *config.NrpConfig) bool {
	logger = blablo.NewLogger("cfg-prcsr", cmdArgs.LogLevel)
	logger.Debug("Init 'Nginx Processor'")

	logger.Info(f("Found (%s) services configuration", c.WithGreen(fmt.Sprint(len(nrpConfig.Services)))))
	logger.Debug(f("Certificates base folder: %s", c.WithCyan(nrpConfig.Letsencrypt.CertFilesPath)))

	confAvailablePath := filepath.Join(nrpConfig.Nginx.ConfigPath, "conf.available")
	if err := os.RemoveAll(confAvailablePath); err != nil {
		logger.Error(f("Failed to clean folder:", c.WithCyan(confAvailablePath)), "err", err)
	}

	if err := os.Mkdir(confAvailablePath, os.ModePerm); err != nil {
		logger.Error(f("Failed to re-create folder:", c.WithCyan(confAvailablePath)), "err", err)
	}
	logger.Debug(f("Folder cleaned: %s", c.WithGreen(confAvailablePath)))

	svcConfTmplPath := filepath.Join(nrpConfig.Nginx.ConfigPath, "/templates/service.conf.tmpl")
	_, err1 := loadConfTemplate(svcConfTmplPath)

	defaultSvcConfTmplPath := filepath.Join(nrpConfig.Nginx.ConfigPath, "/templates/default.conf.tmpl")
	_, err2 := loadDefaultConfTemplate(defaultSvcConfTmplPath)

	acmeChallngeConfTmplPath := filepath.Join(nrpConfig.Nginx.ConfigPath, "/templates/acme-challenge.conf.tmpl")
	_, err3 := loadAcmeChallengeConfTemplate(acmeChallngeConfTmplPath)

	logger.Debug("Init completed for 'Nginx Processor'")

	return (err1 == nil && err2 == nil && err3 == nil)
}

func generateDefaultNginxServerConfig() (*bytes.Buffer, error) {
	var content bytes.Buffer
	err := defaultConfTemplate.Execute(&content, nil)
	if err != nil {
		logger.Error(f("Failed to generate nginx config for service: %s", c.WithCyan("default")), "err", err)
		return nil, err
	}
	// logger.Debug(f("Generated (%s) bytes of config data", c.WithGreen(fmt.Sprint(content.Len()))))
	return &content, nil
}

func generateNginxServerConfig(svcConfig *config.NrpServiceConfig) (*bytes.Buffer, error) {
	var content bytes.Buffer
	err := confTemplate.Execute(&content, svcConfig)
	if err != nil {
		logger.Error(f("Failed to generate nginx config for service: %s", c.WithCyan(svcConfig.Name)), "err", err)
		return nil, err
	}
	// logger.Debug(f("Generated (%s) bytes of config data", c.WithGreen(fmt.Sprint(content.Len()))))
	return &content, nil
}

func generateAcmeChallengeServerConfig(svcConfig *config.NrpServiceConfig) (*bytes.Buffer, error) {
	var content bytes.Buffer
	err := acmeChallengeTemplate.Execute(&content, svcConfig)
	if err != nil {
		logger.Error(f(
			"Failed to generate %s config for service: %s",
			c.WithCyan("acme-challenge"),
			c.WithCyan(svcConfig.Name),
		), "err", err)
		return nil, err
	}
	// logger.Debug(f("Generated (%s) bytes of config data", c.WithGreen(fmt.Sprint(content.Len()))))
	return &content, nil
}

func CreateDefaultConfFile(nrpConfig *config.NrpConfig) bool {
	logger.Info(f("Generating nginx configs in '%s'", c.WithCyan("conf.available")))
	// Generate default "welcome page" nginx server config
	content, err := generateDefaultNginxServerConfig()
	if err != nil {
		return false
	}
	filePath := filepath.Join(nrpConfig.Nginx.ConfigPath, "conf.available", f("%v-%s.conf", 0, "default"))
	if err := os.WriteFile(filePath, content.Bytes(), 0644); err != nil {
		logger.Error(f("Saving content to file: %s", c.WithCyan(filePath)))
		return false
	} else {
		logger.Info(f("Saved (%s) bytes to file: %s", c.WithCyan(f("%v", content.Len())), c.WithGreen(filePath)))
	}

	return true
}

func CreateServiceConfFile(nrpConfig *config.NrpConfig, idx int, svcConfig *config.NrpServiceConfig) bool {
	// Continue to generate nginx  server config
	content, err := generateNginxServerConfig(svcConfig)
	if err != nil {
		return false
	}

	filePath := filepath.Join(nrpConfig.Nginx.ConfigPath, "conf.available", f("%v-%s.conf", idx+1, svcConfig.Name))
	if err := os.WriteFile(filePath, content.Bytes(), 0644); err != nil {
		logger.Error(f("Saving content to file: %s", c.WithCyan(filePath)))
		return false
	} else {
		logger.Debug(f("Saved (%s) bytes to file: %s", c.WithCyan(f("%v", content.Len())), c.WithGreen(filePath)))
		return true
	}

}

func createAcmeChallengeServerConfigFile(nrpConfig *config.NrpConfig, svcConfig *config.NrpServiceConfig) bool {
	content, err := generateAcmeChallengeServerConfig(svcConfig)
	if err != nil {
		return false
	}

	filePath := filepath.Join(nrpConfig.Nginx.ConfigPath, "conf.d", f("%s-acme-challenge.conf", svcConfig.Name))
	if err := os.WriteFile(filePath, content.Bytes(), 0644); err != nil {
		logger.Error(f("Saving content to file: %s", c.WithCyan(filePath)))
		return false
	} else {
		logger.Info(f("Saved (%s) bytes to file: %s", c.WithCyan(f("%v", content.Len())), c.WithGreen(filePath)))
		return true
	}
}

func removeAcmeChallengeServerConfigFile(nrpConfig *config.NrpConfig, svcConfig *config.NrpServiceConfig) bool {
	filePath := filepath.Join(nrpConfig.Nginx.ConfigPath, "conf.d", f("%s-acme-challenge.conf", svcConfig.Name))

	_, err := os.Stat(filePath)

	if err != nil {
		logger.Error(f("File not found: %s", c.WithCyan(filePath)))
		return false
	}

	err = os.Remove(filePath)
	if err != nil {
		logger.Error(f("Can't delete file: %s", c.WithCyan(filePath)), "err", err)
		return false
	}

	logger.Debug(f("Successfully deleted file: %s", c.WithCyan(filePath)))
	return true
}

/*
Assume *.conf files already in conf.available.
Now they will be copied to conf.d for nginx to use them
*/
func CopyConfFiles(nrpConfig *config.NrpConfig) bool {
	logger.Debug(f("Copying service conf files to folder: %s", c.WithCyan("conf.d")))

	if ok := cleanNginxConfDPath(nrpConfig); !ok {
		return false
	}

	confDPath := filepath.Join(nrpConfig.Nginx.ConfigPath, "conf.d")
	confAvailablePath := filepath.Join(nrpConfig.Nginx.ConfigPath, "conf.available")
	files, err := os.ReadDir(confAvailablePath)
	if err != nil {
		logger.Error(f("Failed to read conf.available directory: %v", c.WithRed(err.Error())))
		return false
	}

	for _, file := range files {
		src := filepath.Join(confAvailablePath, file.Name())
		dst := filepath.Join(confDPath, file.Name())

		in, err := os.Open(src)
		if err != nil {
			logger.Error(f("Failed to open source file: %v", c.WithRed(err.Error())))
		}
		defer in.Close()

		out, err := os.Create(dst)
		if err != nil {
			logger.Error(f("Failed to create destination file: %v", c.WithRed(err.Error())))
		}
		defer out.Close()

		_, err = io.Copy(out, in)
		if err != nil {
			logger.Error(f("Failed to copy file: %v", c.WithRed(err.Error())))
		}

		if err == nil {
			logger.Debug(f("Succesfuly copied file: %s", c.WithGreen(dst)))
		}
	}

	logger.Info(f("Finished copying files to folder: %s", c.WithCyan(confDPath)))
	return true
}
