package retry

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

// RetryConfig defines the configuration for the retry logic
type RetryConfig struct {
	MaxRetries   int           // Maximum number of retries
	WaitDuration time.Duration // Duration to wait between retries
}

var defaultConfig = &RetryConfig{
	MaxRetries:   3,
	WaitDuration: 3 * time.Second,
}

func HttpRetry(client *http.Client, req *http.Request) (*http.Response, error) {
	var err error
	var resp *http.Response

	config := defaultConfig
	for attempt := 0; attempt < config.MaxRetries; attempt++ {
		resp, err = client.Do(req)

		if shouldRetry(resp, err) {
			if err != nil {
				fmt.Printf("Request error: %v", err)
			} else {
				fmt.Printf("Http status: %v. Retrying...\n", resp.Status)
			}

			resp.Body.Close()
		} else {
			return resp, err
		}

		// Wait for the specified duration before retrying
		time.Sleep(config.WaitDuration)
	}

	// Return the last error if all retries failed
	return nil, errors.New(fmt.Sprintf("Retry failed to make %v", req.URL))
}

func shouldRetry(resp *http.Response, err error) bool {
	return err != nil || !(resp.StatusCode >= 200 && resp.StatusCode < 300)
}
