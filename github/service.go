package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"ton-tg-bot/logger"
	"ton-tg-bot/models"
)

var httpTransport *http.Transport

func init() {
	httpTransport = newTransport()
}

func ValidateGitHub(userName string, userID int64) error {
	url := fmt.Sprintf(pathUser, userName)
	data, err := getGitHubData(url, userID)
	if err != nil {
		return err
	}

	if data.Message == models.NotFound {
		return fmt.Errorf("user %w", ErrNotFound)
	}

	if data.CreatedAt.After(time.Now().Add(-1 * month)) {
		return ErrNewUser
	}

	url = fmt.Sprintf(pathRepo, userName)
	data, err = getGitHubData(url, userID)
	if err != nil {
		return err
	}

	if data.Message == models.NotFound {
		return fmt.Errorf("repo %w", ErrNotFound)
	}

	return nil
}

func getGitHubData(url string, userID int64) (models.Data, error) {
	var data models.Data

	client := newClient(httpTransport)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return data, fmt.Errorf("data %d new request create: %w", userID, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return data, fmt.Errorf("data %d response: %w", userID, err)
	}
	defer func() {
		if errClose := resp.Body.Close(); errClose != nil {
			logger.LogWarn(userID, "body close:", err)
		}
	}()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return data, fmt.Errorf("data %d read body: %w", userID, err)
	}

	if err = json.Unmarshal(respBody, &data); err != nil {
		return data, fmt.Errorf("data %d unmarshal body: %w", userID, err)
	}

	return data, nil
}
