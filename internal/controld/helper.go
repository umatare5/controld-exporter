// Package controld provides a client for interacting with the ControlD API.
package controld

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/umatare5/controld-exporter/internal/log"
)

// isSuccess checks if the "success" field in the response is true.
func isSuccess(response map[string]any) bool {
	success, ok := response["success"].(bool)
	return ok && success
}

// isSuccessStatus checks if the response status is in the 2xx range.
func isSuccessStatus(resp *http.Response) bool {
	return resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusMultipleChoices
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

// sendRequest performs an HTTP request, handles errors, and decodes the response into the result.
func (t *Client) sendRequest(uri string, headers map[string]string, result any) error {
	log.Debugf("Sending request to URI: %s, headers: %s", uri, headers) // Debug log for the request URI

	req, err := t.createRequest(uri, headers)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
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
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+t.apiKey)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	return req, nil
}

// maxErrorBodyLen bounds the vendor error envelope quoted into a returned error.
const maxErrorBodyLen = 512

// summarizeErrorBody renders a non-2xx body as a bounded single-line string.
func summarizeErrorBody(body []byte) string {
	s := strings.Join(strings.Fields(string(body)), " ")
	if len(s) > maxErrorBodyLen {
		return s[:maxErrorBodyLen] + "..."
	}
	return s
}

// codeNoData is the 404 Control D answers when a collection holds nothing.
const codeNoData = 40401

// ErrNoData reports an envelope that answers "nothing here" rather than a failure.
var ErrNoData = errors.New("endpoint holds no data")

// errorEnvelope is the body Control D returns beside a non-2xx status. The
// measured shape also carries error.date and error.message, both unread here.
type errorEnvelope struct {
	Error struct {
		Code int `json:"code"` // First three digits restate the HTTP status
	} `json:"error"`
}

// isNoDataBody reports whether a non-2xx body carries the empty-collection code.
func isNoDataBody(body []byte) bool {
	var envelope errorEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		return false
	}
	return envelope.Error.Code == codeNoData
}

func (t *Client) handleResponse(resp *http.Response, endpoint string, result any) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Debugf("Raw JSON response: %s", string(body))

	if !isSuccessStatus(resp) {
		if isNoDataBody(body) {
			return fmt.Errorf("%w: %s: %s", ErrNoData, endpoint, summarizeErrorBody(body))
		}
		return fmt.Errorf(
			"unexpected status %q from endpoint: %s: %s",
			resp.Status, endpoint, summarizeErrorBody(body),
		)
	}

	var rawResponse map[string]any
	if err := json.Unmarshal(body, &rawResponse); err != nil {
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
