package configProcessor

import (
	c "beaver/blablo/color"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

func startNginx(cmd string) (bool, error) {
	logger.Debug(f("Starting NRP with command: %s", c.WithYellow(cmd)))
	proc := exec.Command("bash", "-c", cmd)

	err := proc.Start()
	if err != nil {
		logger.Error(f("Error starting NRP: %s", err))
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
		logger.Error(f("Error starting NRP: %s", c.WithRed(err.Error())))
		return 500, err
	}
	logger.Debug(f("Nginx status : %s", c.WithCyan(strings.TrimSpace(string(output)))))

	var lastStatusCode int
	var res *http.Response
	retries := 3

	for i := 0; i < retries; i++ {
		res, err = http.Get("http://127.0.0.1")
		if err != nil {
			logger.Error("Error making HTTP request", "err", err)
			time.Sleep(1 * time.Second)
			continue
		}

		if res.StatusCode != 200 {
			logger.Debug(f("Nginx http request status: %s", c.WithCyan(fmt.Sprint(res.StatusCode))))
			res.Body.Close()
			lastStatusCode = res.StatusCode
			time.Sleep(1 * time.Second)
			continue
		}

		lastStatusCode = res.StatusCode
		break
	}

	return lastStatusCode, nil
}

func stopNginx(cmd string) (bool, error) {
	logger.Debug(f("Stopping NRP with command: %s", c.WithYellow(cmd)))
	proc := exec.Command("bash", "-c", cmd)

	_, err := proc.CombinedOutput()
	if err != nil {
		logger.Error(f("Error stopping NRP: %s", err))
		return false, err
	}

	return true, nil
}
