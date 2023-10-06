package nginxConfigProcessor

import (
	"fmt"
	c "github.com/oleksii-honchar/coteco"
	"net/http"
	"time"
)

func retryRequestUntil200(url string, retries int, timeWaitSec int) int {
	var lastStatusCode int = 0
	var res *http.Response
	var err error

	for i := 0; i < retries; i++ {
		res, err = http.Get(url)
		if err != nil {
			logger.Debug(f("Error making HTTP request: %s", c.WithRed(err.Error())))
			time.Sleep(1 * time.Second)
			continue
		}

		if res.StatusCode != 200 {
			logger.Debug(f("HTTP request failed with status code: %s", c.WithCyan(fmt.Sprint(res.StatusCode))))
			res.Body.Close()
			lastStatusCode = res.StatusCode
			time.Sleep(time.Duration(timeWaitSec) * time.Second)
			continue
		} else {
			logger.Debug(f("Succesfuly made HTTP request: %s", c.WithGreen(fmt.Sprint(res.StatusCode))))
		}

		lastStatusCode = res.StatusCode
		break
	}

	return lastStatusCode
}

func makeRequest(url string) (*http.Response, error) {
	res, err := http.Get(url)
	if err != nil {
		logger.Error("Error creating HTTP request:", "err", err)
		return nil, err
	}
	// defer res.Body.Close()

	// buf := new(bytes.Buffer)
	// io.Copy(buf, res.Body)
	// if err != nil {
	// 	logger.Error("Error reading HTTP response:", err)
	// 	return nil, err
	// }
	return res, nil
}
