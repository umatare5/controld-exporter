// Package controld provides a client for interacting with the ControlD API.
package controld

import (
	"fmt"
	"time"
)

const (
	DNSQueriesReportEndpoint = "/reports/dns-queries/all-by-verdict/time-series" // Endpoint for DNS query statistics
)

// QueryStatsResponse represents the response structure for DNS query statistics.
type QueryStatsResponse struct {
	Success bool `json:"success"`
	Body    struct {
		EndTS       int    `json:"endTs"`
		StartTS     int    `json:"startTs"`
		Granularity string `json:"granularity"`
		Tz          string `json:"tz"`
		Queries     []struct {
			TS    string         `json:"ts"`
			Count map[string]int `json:"count"`
		} `json:"queries"`
	} `json:"body"`
}

// GetDNSQueriesReport fetches DNS query statistics without additional headers.
func (t *Client) GetDNSQueriesReport(statsEndpoint string) (*QueryStatsResponse, error) {
	return t.sendDNSQueriesReportRequest(
		statsEndpoint, t.buildDNSQueriesReportURI(DNSQueriesReportEndpoint), nil,
	)
}

// GetSubOrgDNSQueriesReport fetches DNS query statistics with additional headers for a specific organization.
func (t *Client) GetSubOrgDNSQueriesReport(statsEndpoint, orgID string) (*QueryStatsResponse, error) {
	return t.sendDNSQueriesReportRequest(
		statsEndpoint, t.buildDNSQueriesReportURI(DNSQueriesReportEndpoint), t.buildOrgIDHeader(orgID),
	)
}

// sendDNSQueriesReportRequest sends a request to fetch DNS query statistics.
func (t *Client) sendDNSQueriesReportRequest(
	statsEndpoint, uri string,
	headers map[string]string,
) (*QueryStatsResponse, error) {
	var data QueryStatsResponse
	if err := t.sendReportAPIRequest(statsEndpoint, uri, headers, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// buildDNSQueriesReportURI constructs the URI for the DNS queries report.
func (t *Client) buildDNSQueriesReportURI(baseEndpoint string) string {
	return fmt.Sprintf(
		"%s?startTs=%d&granularity=%s&tz=%s",
		baseEndpoint,
		time.Now().Add(-1*time.Minute).Unix(),
		"minute",
		time.Now().Location().String(),
	)
}
