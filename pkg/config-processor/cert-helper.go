package configProcessor

import (
	c "beaver/blablo/color"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

func CheckCertificateFiles(configName string) bool {
	configFilesFolder := filepath.Join(nrpConfig.Letsencrypt.CertFilesPath, configName)
	sslCertPath := filepath.Join(configFilesFolder, "cert.pem")
	sslKeyPath := filepath.Join(configFilesFolder, "privkey.pem")

	_, certErr := os.Stat(sslCertPath)
	_, keyErr := os.Stat(sslKeyPath)

	if certErr == nil && keyErr == nil {
		return true
	}

	logger.Warn(f("Missing certificates in path: %s", c.WithYellow(configFilesFolder)))
	return false
}

func requestCertificate(svcConfig *NrpServiceConfig) (bool, error) {
	var data = NewCertRequest{
		DryRun:   nrpConfig.Letsencrypt.DryRun,
		BaseDir:  nrpConfig.Letsencrypt.BasePath,
		CertName: svcConfig.Name,
		Email:    nrpConfig.Letsencrypt.Email,
		Domain:   svcConfig.DomainName,
	}

	// parse and prepare certbot cmd
	tmpl := template.Must(template.New("RequestCertCmd").Parse(nrpConfig.Letsencrypt.RequestCertCmdTmpl))
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, data)
	if err != nil {
		logger.Error(f("Error creating request certificate cmd: %s", err))
		return false, err
	}
	cmd := buf.String()

	logger.Debug(f("Requesting certificate for domain: %s", c.WithYellow(svcConfig.DomainName)))
	logger.Debug("With data", "data", data)
	logger.Debug(f("With cmd: %s", c.WithYellow(cmd)))

	proc := exec.Command("bash", "-c", cmd)

	output, err := proc.CombinedOutput()
	logger.Info(f("Certbot response: \n%s", c.WithYellow(string(output))))
	if err != nil {
		logger.Error(f("Error requesting certificate: %s", err))
		return false, err
	}
	return true, nil
}

/*
Request certificates using certbot
- Generate acme-challenge server.conf -> conf.d
- Start nginx with tmp server
- Request certificate using certbot
- Remove tmp server, stop nginx
- If success - return true, else -> false
*/
func CreateCertificateFiles(svcConfig *NrpServiceConfig) bool {
	// defer removeAcmeChallengeServerConfigFile(svcConfig)
	defer stopNginx(nrpConfig.Nginx.StopCmd)

	logger.Info(f("Creating certificates for: %s", c.WithCyan(svcConfig.Name)))

	if ok := createAcmeChallengeServerConfigFile(svcConfig); !ok {
		return false
	}

	// When starting nginx from dev: nrp-cli based docker config used
	// When starting nginx from prod: nrp-cli will be inside of container and nginx will be available directly
	if ok, _ := startNginx(nrpConfig.Nginx.StartCmd); !ok {
		return false
	}

	if status, err := getNginxStatus(nrpConfig.Nginx.StatusCmd); err != nil && status != 200 {
		return false
	} else {
		logger.Info(f("Nginx status: %s", c.WithCyan(fmt.Sprint(status))))
	}

	// let's make certbot to do its job
	_, _ = requestCertificate(svcConfig)

	logger.Info(f("Finished creating the certificates"))
	return true
}
