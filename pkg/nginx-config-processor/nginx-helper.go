package nginxConfigProcessor

import (
	"config"
	"os"
	"path/filepath"

	c "github.com/oleksii-honchar/coteco"

	"os/exec"
	"strings"
	"time"
)

func startNginx(cmd string) (bool, error) {
	logger.Debug(f("Starting Nginx with command: %s", c.WithYellow(cmd)))
	proc := exec.Command("bash", "-c", cmd)

	err := proc.Start()
	if err != nil {
		logger.Debug(f("Error starting Nginx: %s", err))
		return false, err
	}

	for {
		if proc.Process != nil {
			break
		}
		logger.Debug(f("Waiting 1 sec for process to start"))
		time.Sleep(1 * time.Second)
	}

	return true, nil
}

func getNginxStatus(cmd string) (int, error) {
	logger.Debug(f("Checking Nginx status with command: %s", c.WithYellow(cmd)))
	proc := exec.Command("bash", "-c", cmd)

	output, err := proc.CombinedOutput()
	if err != nil {
		logger.Debug(f("Error checking status from Nginx: %s", c.WithRed(err.Error())))
	} else {
		logger.Debug(f("Nginx status : %s", c.WithCyan(strings.TrimSpace(string(output)))))
		return 200, nil
	}

	statusCode := retryRequestUntil200("http://127.0.0.1", 10, 2)

	return statusCode, nil
}
func getNginxLogs(cmd string) (bool, error) {
	logger.Debug(f("Getting Nginx logs with command: %s", c.WithYellow(cmd)))
	proc := exec.Command("bash", "-c", cmd)

	output, err := proc.CombinedOutput()
	if err != nil {
		logger.Error(f("Error getting logs from Nginx: %s", c.WithRed(err.Error())))
		return false, err
	}
	logger.Debug(f("Nginx logs : \n%s", c.WithGray247((string(output)))))

	return true, nil
}

func stopNginx(cmd string) (bool, error) {
	logger.Debug(f("Stopping Nginx with command: %s", c.WithYellow(cmd)))
	proc := exec.Command("bash", "-c", cmd)

	_, err := proc.CombinedOutput()
	if err != nil {
		logger.Error(f("Error stopping NRP: %s", err))
		return false, err
	}

	return true, nil
}

func cleanNginxConfDPath(nrpConfig *config.NrpConfig) bool {
	confDPath := filepath.Join(nrpConfig.Nginx.ConfigPath, "conf.d")
	if err := os.RemoveAll(confDPath); err != nil {
		logger.Error(f("Failed to clean folder:", c.WithCyan(confDPath)), "err", err)
		return false
	}
	err := os.MkdirAll(confDPath, os.ModePerm)
	if err != nil {
		logger.Error(f("Failed to re-create folder (%s): %s", c.WithCyan(confDPath), c.WithRed(err.Error())))
		return false
	}
	return true
}
