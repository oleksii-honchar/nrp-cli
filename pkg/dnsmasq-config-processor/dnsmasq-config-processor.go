package dnsmasqConfigProcessor

import (
	"bytes"
	cmdArgs "cmd-args"
	cfgProc "config-processor"
	"fmt"
	"os"
	"path/filepath"

	_ "embed"

	"text/template"

	"github.com/oleksii-honchar/blablo"
	c "github.com/oleksii-honchar/coteco"
)

//go:embed dnsmasq.conf.tmpl
var DnsmasqConfTmpl []byte

var f = fmt.Sprintf
var logger *blablo.Logger

func getUniqueDomains(config *cfgProc.NrpConfig) []string {
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

func GenerateConfig(config *cfgProc.NrpConfig) bool {
	logger = blablo.NewLogger("dnsmsq-cfg", string(cmdArgs.LogLevel))
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
	}{
		SquidPort: config.Squid.Port,
		Domains:   getUniqueDomains(config),
	}

	var content bytes.Buffer
	err = confTemplate.Execute(&content, tmplConfig)
	if err != nil {
		logger.Error(f("Failed to generate %s config: %s", c.WithCyan("dnsmasq"), c.WithRed(err.Error())))
		return false
	}
	// logger.Debug(f("Generated (%s) bytes of config data", c.WithGreen(fmt.Sprint(content.Len()))))

	filePath := filepath.Join(config.Dnsmasq.ConfigPath, "dnsmasq.conf")
	if err := os.WriteFile(filePath, content.Bytes(), 0644); err != nil {
		logger.Error(f("Saving content to file (%s): %s", c.WithCyan(filePath), c.WithRed(err.Error())))
		return false
	} else {
		logger.Debug(f("Saved (%s) bytes to file: %s", c.WithCyan(f("%v", content.Len())), c.WithGreen(filePath)))
	}

	logger.Info(c.WithGreen(f("Dnsmasq config generation completed")))

	return true
}
