package dnsmasqConfigProcessor

import (
	"bytes"
	cmdArgs "cmd-args"
	"config"
	"fmt"
	"os"

	_ "embed"

	"text/template"

	"github.com/oleksii-honchar/blablo"
	c "github.com/oleksii-honchar/coteco"
)

//go:embed dnsmasq.conf.tmpl
var DnsmasqConfTmpl []byte

var f = fmt.Sprintf
var logger *blablo.Logger

func getUniqueDomains(config *config.NrpConfig) []string {
	uniqueMap := make(map[string]bool)
	uniqueStrings := []string{}

	for _, svc := range config.Services {
		if !uniqueMap[svc.DomainName] {
			uniqueMap[svc.DomainName] = true
			uniqueStrings = append(uniqueStrings, svc.DomainName)
		}
	}

	return uniqueStrings
}

func GenerateConfig(config *config.NrpConfig) bool {
	logger = blablo.NewLogger("dnsmsq-cfg", string(cmdArgs.LogLevel), false)
	logger.Debug(f("Processing Dnsmasq config: %s%+v%s", c.Yellow, config.Dnsmasq, c.Reset))
	if config.Squid.UseDnsmasq != "yes" {
		logger.Debug(f("Dnsmasq config is disabled in 'squid' section. Skipping config generation."))
		return true
	}

	confTemplate, err := template.New("dnsmsq-config").Parse(string(DnsmasqConfTmpl))
	if err != nil {
		logger.Error(f("Error creating dnsmsq config tmpl: %s", c.WithRed(err.Error())))
		return false
	}

	var tmplConfig = struct {
		SquidPort int
		Domains   []string
		Logs      string
	}{
		SquidPort: config.Squid.Port,
		Domains:   getUniqueDomains(config),
		Logs:      config.Dnsmasq.Logs,
	}

	var content bytes.Buffer
	err = confTemplate.Execute(&content, tmplConfig)
	if err != nil {
		logger.Error(f("Failed to generate %s config: %s", c.WithCyan("dnsmasq"), c.WithRed(err.Error())))
		return false
	}
	// logger.Debug(f("Generated (%s) bytes of config data", c.WithGreen(fmt.Sprint(content.Len()))))

	if err := os.WriteFile(config.Dnsmasq.ConfigPath, content.Bytes(), 0644); err != nil {
		logger.Error(f("Saving content to file (%s): %s", c.WithCyan(config.Dnsmasq.ConfigPath), c.WithRed(err.Error())))
		return false
	} else {
		logger.Debug(f("Saved (%s) bytes to file: %s", c.WithCyan(f("%v", content.Len())), c.WithGreen(config.Dnsmasq.ConfigPath)))
	}

	logger.Info(c.WithGreen(f("Config generation completed for %s", c.WithCyan("'dnsmasq'"))))

	return true
}
