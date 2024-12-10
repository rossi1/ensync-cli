package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"go.uber.org/zap"
	"golang.org/x/time/rate"

	"github.com/ensync-cli/internal/domain"
)

type Client struct {
	baseURL     string
	apiKey      string
	httpClient  *http.Client
	rateLimiter *rate.Limiter
	logger      *zap.Logger
}

func NewClient(baseURL, apiKey string, opts ...ClientOption) *Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.RetryWaitMin = 1 * time.Second
	retryClient.RetryWaitMax = 5 * time.Second

	c := &Client{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: retryClient.StandardClient(),
		logger:     zap.NewNop(),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) doRequest(ctx context.Context, method, path string, query url.Values, body interface{}) ([]byte, error) {
	// Handle rate limiting
	if c.rateLimiter != nil {
		err := c.rateLimiter.Wait(ctx)
		if err != nil {
			return nil, fmt.Errorf("rate limit error: %w", err)
		}
	}

	// Prepare request URL
	reqURL := c.baseURL + path
	if query != nil {
		reqURL = fmt.Sprintf("%s?%s", reqURL, query.Encode())
	}

	// Marshal body if present
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, reqURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("X-API-KEY", c.apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	// Log request
	c.logger.Debug("Sending request",
		zap.String("method", method),
		zap.String("path", path),
		zap.String("url", reqURL),
	)

	// Execute request with timeout
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Log response
	c.logger.Debug("Received response",
		zap.Int("status", resp.StatusCode),
		zap.Int("body_size", len(respBody)),
	)

	// Handle non-200 responses
	if resp.StatusCode >= 400 {
		var apiErr struct {
			Message string `json:"message"`
			Code    string `json:"code"`
		}
		if err := json.Unmarshal(respBody, &apiErr); err != nil {
			// If can't unmarshal error response, return raw response
			return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(respBody))
		}
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    apiErr.Message,
			Code:       apiErr.Code,
		}
	}

	return respBody, nil
}

func (c *Client) ListEvents(ctx context.Context, params *ListParams) (*domain.EventList, error) {
	query := url.Values{}
	query.Set("pageIndex", fmt.Sprintf("%d", params.PageIndex))
	query.Set("limit", fmt.Sprintf("%d", params.Limit))
	query.Set("order", params.Order)
	query.Set("orderBy", params.OrderBy)

	data, err := c.doRequest(ctx, http.MethodGet, "/event", query, nil)
	if err != nil {
		return nil, err
	}

	var response domain.EventList
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

func (c *Client) CreateEvent(ctx context.Context, event *domain.Event) error {
	_, err := c.doRequest(ctx, http.MethodPost, "/event", nil, event)
	return err
}

func (c *Client) UpdateEvent(ctx context.Context, event *domain.Event) error {
	url := fmt.Sprintf("/event/%s", event.ID)
	payload := map[string]interface{}{
		"name":    event.Name,
		"payload": event.Payload,
	}

	_, err := c.doRequest(ctx, http.MethodPut, url, nil, payload)
	return err
}

func (c *Client) ListAccessKeys(ctx context.Context, params *ListParams) (*domain.AccessKeyList, error) {
	query := url.Values{}
	query.Set("pageIndex", fmt.Sprintf("%d", params.PageIndex))
	query.Set("limit", fmt.Sprintf("%d", params.Limit))
	query.Set("order", params.Order)
	query.Set("orderBy", params.OrderBy)

	if v, ok := params.Filter["accessKey"]; ok && v != "" {
		query.Set("accessKey", v)
	}

	data, err := c.doRequest(ctx, http.MethodGet, "/access-key", query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list access keys: %w", err)
	}

	var response domain.AccessKeyList
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

func (c *Client) CreateAccessKey(ctx context.Context, accessKey *domain.AccessKey) (*domain.AccessKey, error) {
	payload := map[string]interface{}{
		"permissions": accessKey.Permissions,
		"access":      accessKey.Permissions.Access,
	}

	data, err := c.doRequest(ctx, http.MethodPost, "/access-key", nil, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create access key: %w", err)
	}

	var response struct {
		AccessKey string `json:"accessKey"`
	}
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Return new access key with the generated key
	return &domain.AccessKey{
		Key:         response.AccessKey,
		Permissions: accessKey.Permissions,
		CreatedAt:   time.Now(), // Server might not return this
	}, nil
}

func (c *Client) VerifyAccessKey(ctx context.Context, key string) (bool, error) {
	url := fmt.Sprintf("/access/verify/%s", key)

	data, err := c.doRequest(ctx, http.MethodGet, url, nil, nil)
	if err != nil {
		return false, fmt.Errorf("failed to verify access key: %w", err)
	}

	var response struct {
		Status bool `json:"status"`
	}
	if err := json.Unmarshal(data, &response); err != nil {
		return false, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Status, nil
}

func (c *Client) GetAccessKeyPermissions(ctx context.Context, key string) (*domain.AccessKeyPermissions, error) {
	url := fmt.Sprintf("/access-key/permissions/%s", key)

	data, err := c.doRequest(ctx, http.MethodGet, url, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get access key permissions: %w", err)
	}

	var response struct {
		Permissions *domain.AccessKeyPermissions `json:"permissions"`
	}
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Permissions, nil
}

func (c *Client) SetAccessKeyPermissions(ctx context.Context, key string, permissions *domain.AccessKeyPermissions) error {
	url := fmt.Sprintf("/access-key/permissions/%s", key)

	payload := map[string]interface{}{
		"access":  permissions.Access,
		"send":    permissions.Send,
		"receive": permissions.Receive,
	}

	_, err := c.doRequest(ctx, http.MethodPost, url, nil, payload)
	if err != nil {
		return fmt.Errorf("failed to set access key permissions: %w", err)
	}

	return nil
}