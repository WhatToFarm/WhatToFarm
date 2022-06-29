package external

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"ton-tg-bot/core"
	"ton-tg-bot/logger"
)

func ExtService(fileName string, body io.Reader) (string, error) {
	client := &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   30 * time.Second,
	}

	path := core.BasePath + fileName
	extUrl := url.URL{
		Scheme: "http",
		Host:   core.Host,
		Path:   path,
	}
	req, err := http.NewRequest(http.MethodGet, extUrl.String(), body)
	if err != nil {
		return "", fmt.Errorf("file: %s, new request create: %w", fileName, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("file: %s, response: %w", fileName, err)
	}
	defer func() {
		if errClose := resp.Body.Close(); errClose != nil {
			logger.LogWarn("body close:", err)
		}
	}()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("file: %s read body: %w", fileName, err)
	}

	return string(respBody), nil
}
