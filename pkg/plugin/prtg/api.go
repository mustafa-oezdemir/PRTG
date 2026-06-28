package prtg

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func NewApi(baseURL, apiKey string, cacheTime, requestTimeout time.Duration) *Api {
	return &Api{
		baseURL:   baseURL,
		apiKey:    apiKey,
		timeout:   requestTimeout,
		cacheTime: cacheTime,
		cache:     make(map[string]cacheItem),
	}
}

func (a *Api) ClearCache() {
	a.cacheMu.Lock()
	defer a.cacheMu.Unlock()
	a.cache = make(map[string]cacheItem)
}

func (a *Api) buildApiUrl(method string, params map[string]string) (string, error) {
	baseUrl := fmt.Sprintf("%s/api/%s", a.baseURL, method)
	u, err := url.Parse(baseUrl)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	q := url.Values{}
	q.Set("apitoken", a.apiKey)

	for key, value := range params {
		q.Set(key, value)
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}

func (a *Api) SetTimeout(timeout time.Duration) {
	if timeout > 0 {
		if timeout < 20*time.Second {
			timeout = 20 * time.Second
		}
		a.timeout = timeout
	}
}

func (a *Api) baseExecuteRequest(endpoint string, params map[string]string) ([]byte, error) {
	apiUrl, err := a.buildApiUrl(endpoint, params)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL for endpoint '%s': %w", endpoint, err)
	}

	client := &http.Client{
		Timeout: a.timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for endpoint '%s': %w", endpoint, err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed for endpoint '%s': %w", endpoint, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
		log.DefaultLogger.Error("Access denied: please verify API token and permissions", "endpoint", endpoint)
		return nil, fmt.Errorf("access denied: please verify API token and permissions (endpoint: %s)", endpoint)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d for endpoint: %s", resp.StatusCode, endpoint)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body for endpoint '%s': %w", endpoint, err)
	}
	return body, nil
}

func (a *Api) GetCacheTime() time.Duration {
	return a.cacheTime
}
