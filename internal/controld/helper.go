// Package controld provides a client for interacting with the ControlD API.
package controld

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/umatare5/controld-exporter/internal/log"
)

// isSuccess checks if the "success" field in the response is true.
func isSuccess(response map[string]any) bool {
	success, ok := response["success"].(bool)
	return ok && success
}

// buildOrgIDHeader creates a header map containing the "X-Force-Org-Id" field.
func (t *Client) buildOrgIDHeader(orgID string) map[string]string {
	return map[string]string{"X-Force-Org-Id": orgID}
}

// sendAPIRequest constructs the full URI and delegates the request to sendRequest.
func (t *Client) sendAPIRequest(endpoint string, headers map[string]string, result any) error {
	uri := t.baseURL + endpoint
	return t.sendRequest(uri, headers, result)
}

// sendReportAPIRequest constructs the full URI for Analytics API and delegates the request to sendRequest.
func (t *Client) sendReportAPIRequest(stats_endpoint, endpoint string, headers map[string]string, result any) error {
	uri := "https://" + stats_endpoint + ".analytics.controld.com" + endpoint
	return t.sendRequest(uri, headers, result)
}

// sendRequest performs an HTTP request, handles errors, and decodes the response into the result.
func (t *Client) sendRequest(uri string, headers map[string]string, result any) error {
	log.Debugf("Sending request to URI: %s, headers: %s", uri, headers) // Debug log for the request URI

	req, err := t.createRequest(uri, headers)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("Error sending request to %s: %s", uri, err)
		return err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Errorf("Error closing response body: %v", closeErr)
		}
	}()

	return t.handleResponse(resp, uri, result)
}

func (t *Client) createRequest(url string, headers map[string]string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.apiKey))
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	return req, nil
}

func (t *Client) handleResponse(resp *http.Response, endpoint string, result any) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error reading response body: %s", err)
		return err
	}
	log.Debugf("Raw JSON response: %s", string(body))

	var rawResponse map[string]any
	if err := json.Unmarshal(body, &rawResponse); err != nil {
		log.Errorf("Error parsing JSON: %s", err)
		return err
	}

	if err := t.handleAPIError(endpoint, isSuccess(rawResponse)); err != nil {
		return err
	}

	return json.Unmarshal(body, result)
}

func (t *Client) handleAPIError(endpoint string, success bool) error {
	if !success {
		return fmt.Errorf("API response indicates failure for endpoint: %s", endpoint)
	}
	return nil
}
