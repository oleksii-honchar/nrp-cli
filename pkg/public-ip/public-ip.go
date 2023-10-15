package publicIp

import (
	cmdArgs "cmd-args"
	"config"
	"io"
	"net/http"
	"os"
	"public-ip/route53"
	"time"

	"fmt"

	"github.com/oleksii-honchar/blablo"
	c "github.com/oleksii-honchar/coteco"
	"gopkg.in/yaml.v3"
)

var f = fmt.Sprintf
var logger *blablo.Logger

func Init(nrpConfig *config.NrpConfig) bool {
	if !cmdArgs.CheckAndUpdatePublicIp {
		return true
	}
	logger = blablo.NewLogger("public-ip", cmdArgs.LogLevel)
	logger.Debug("Init 'Public IP' tools")
	if nrpConfig.PublicIp.CheckAndUpdate != "yes" {
		logger.Debug(f("'Public IP' config is disabled. Skipping."))
		return true
	}

	logger.Debug(f("Loading %s", c.WithCyan("'public-ip.yaml'")))

	var err error
	var data *PublicIpData
	data, err = loadData(nrpConfig.PublicIp.DataPath)
	if err != nil {
		data = &PublicIpData{
			PublicIp:    "",
			LastUpdated: "",
		}
	}
	logger.Debug(f("Loaded data: %s%+v%s", c.Yellow, data, c.Reset))

	publicIp, err := getPublicIp()
	if err != nil {
		return false
	}

	if publicIp == data.PublicIp {
		logger.Info(
			f("Public IP is up to date: %s = %s",
				c.WithCyan(publicIp),
				c.WithCyan(data.PublicIp)),
		)

		if cmdArgs.Force {
			logger.Info(c.WithYellow("Force update public IP"))
		} else {
			return true
		}
	}

	logger.Info(f("Updating public IP: %s -> %s", c.WithRed(data.PublicIp), c.WithYellow(publicIp)))

	if ok := checkAndUpdateDomains(nrpConfig, publicIp); !ok {
		return false
	}

	if ok := saveData(nrpConfig.PublicIp.DataPath, &PublicIpData{
		PublicIp:    publicIp,
		LastUpdated: time.Now().Format("2006-01-02 15:04:05"),
	}, nrpConfig.PublicIp.DryRun); !ok {
		return false
	}

	if nrpConfig.PublicIp.DryRun != "yes" {
		logger.Info(c.WithGreen("'Public IP' update executed successfuly"))
	}

	return true
}

func GenerateCronTask(nrpConfig *config.NrpConfig) string {
	return "/usr/local/bin/nrp-cli -config=/etc/nrp.yaml -log-level=info -check-and-update-public-ip"
}

func loadData(path string) (*PublicIpData, error) {
	file, err := os.Open(path)
	if err != nil {
		logger.Error(f("Failed to open config file: %s", c.WithRed(err.Error())))
		return nil, err
	}
	defer file.Close()

	var data PublicIpData
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		logger.Error(f("Failed to parse data file: %s", c.WithRed(err.Error())))
		return nil, err
	}

	logger.Debug(f("Loaded %s", c.WithCyan(path)))
	return &data, nil
}

func saveData(filePath string, data *PublicIpData, dryRun string) bool {
	if dryRun == "yes" {
		logger.Info(f("Dry run mode. Skipping saving data to file: %s", c.WithCyan(filePath)))
		return true
	}

	yamlData, err := yaml.Marshal(data)
	if err != nil {
		logger.Error(f("Failed to save data file: %s", c.WithRed(err.Error())))
		return false
	}

	if err := os.WriteFile(filePath, yamlData, 0644); err != nil {
		logger.Error(f("Saving content to file: %s", c.WithCyan(filePath)))
		return false
	} else {
		logger.Debug(f(
			"Saved (%s) bytes to file: %s",
			c.WithCyan(f("%v", len(yamlData))),
			c.WithGreen(filePath)),
		)
		return true
	}
}

func getPublicIp() (string, error) {
	var url = "https://api.ipify.org?format=text"
	logger.Debug(f("Making request to get public IP: %s", c.WithYellow(url)))
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error(f("Error receiving public ip: %s", c.WithRed(err.Error())))
	}
	return string(data), nil
}

// iterate trough domains which has configured domainRegistrant and update A record using CDK
func checkAndUpdateDomains(nrpConfig *config.NrpConfig, publicIp string) bool {
	var processedDomains []string
	for _, service := range nrpConfig.Services {
		if service.DomainRegistrant != "" {
			logger.Debug(f("Updating domain: %s", c.WithCyan(service.DomainName)))
			if ok := updateDomainRecord(service.DomainName, service.DomainRegistrant, publicIp, nrpConfig.PublicIp.DryRun); !ok {
				return false
			}
			processedDomains = append(processedDomains, service.DomainName)
		}
	}

	logger.Debug(f("Updated domains: %s%+v%s", c.Cyan, processedDomains, c.Reset))

	return true
}

func updateDomainRecord(domain, domainRegistrant, publicIp, dryRun string) bool {
	switch domainRegistrant {
	case "route53":
		return route53.UpdateDomainIp(domain, publicIp, dryRun)
	}
	return false
}
