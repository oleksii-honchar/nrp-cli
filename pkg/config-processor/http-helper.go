package configProcessor

import (
	"net/http"
)

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
