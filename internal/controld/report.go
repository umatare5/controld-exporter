// Package controld provides a client for interacting with the ControlD API.
package controld

import (
	"fmt"
	"time"
)

const (
	DnsQueriesReportEndpoint = "/reports/dns-queries/all-by-verdict/time-series" // Endpoint for DNS query statistics
)

// QueryStatsResponse represents the response structure for DNS query statistics.
type QueryStatsResponse struct {
	Success bool `json:"success"`
	Body    struct {
		EndTs       int    `json:"endTs"`
		StartTs     int    `json:"startTs"`
		Granularity string `json:"granularity"`
		Tz          string `json:"tz"`
		Queries     []struct {
			Ts    string         `json:"ts"`
			Count map[string]int `json:"count"`
		} `json:"queries"`
	} `json:"body"`
}

// GetDnsQueriesReport fetches DNS query statisticswithout additional headers.
func (t *Client) GetDnsQueriesReport(stats_endpoint string) (*QueryStatsResponse, error) {
	return t.sendDnsQueriesReportRequest(
		stats_endpoint, t.buildDnsQueriesReportUri(DnsQueriesReportEndpoint), nil,
	)
}

// GetSubOrgDnsQueriesReport fetches DNS query statistics with additional headers for a specific organization.
func (t *Client) GetSubOrgDnsQueriesReport(stats_endpoint string, orgID string) (*QueryStatsResponse, error) {
	return t.sendDnsQueriesReportRequest(
		stats_endpoint, t.buildDnsQueriesReportUri(DnsQueriesReportEndpoint), t.buildOrgIDHeader(orgID),
	)
}

// sendDnsQueriesReportRequest sends a request to fetch DNS query statistics.
func (t *Client) sendDnsQueriesReportRequest(stats_endpoint string, uri string, headers map[string]string) (*QueryStatsResponse, error) {
	var data QueryStatsResponse
	if err := t.sendReportAPIRequest(stats_endpoint, uri, headers, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

// buildDnsQueriesReportUri constructs the URI for the DNS queries report.
func (t *Client) buildDnsQueriesReportUri(baseEndpoint string) string {
	return fmt.Sprintf(
		"%s?startTs=%d&granularity=%s&tz=%s",
		baseEndpoint,
		time.Now().Add(-1*time.Minute).Unix(),
		"minute",
		time.Now().Location().String(),
	)
}
