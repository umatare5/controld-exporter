// Package collector contains Prometheus metric collectors for the exporter.
package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/umatare5/controld-exporter/internal/controld"
	"github.com/umatare5/controld-exporter/internal/log"
)

const (
	statsLogPrefix = "stats"
)

// collectStatsMetrics collects DNS query statistics metrics.
func (c *Collector) collectStatsMetrics(ch chan<- prometheus.Metric) {
	if c.isRunningInPersonalMode() {
		c.collectPersonalQueryStatsMetrics(ch)
		c.log.debug(statsLogPrefix, logSkipOrgScraping)
		return
	}

	// Organization metrics are only available in business mode.
	org, err := c.fetchMainOrganization()
	if err != nil {
		c.log.info(statsLogPrefix, logNotFoundMainOrg)
		return
	}
	c.collectMainOrgQueryStatsMetrics(ch, org, org.Body.Organization.StatsEndpoint)

	subOrgs, err := c.fetchSubOrganizations()
	if err != nil {
		c.log.info(statsLogPrefix, logNotFoundSubOrgs)
		return
	}
	c.collectSubOrgQueryStatsMetrics(ch, subOrgs, org.Body.Organization.StatsEndpoint)
}

// collectPersonalQueryStatsMetrics collects DNS query statistics for the personal instance.
func (c *Collector) collectPersonalQueryStatsMetrics(ch chan<- prometheus.Metric) {
	stats, err := c.client.GetDnsQueriesReport("america")
	if err != nil {
		log.Errorf("Error fetching stats for Business: %v", err)
		return
	}

	c.storeStatsMetrics(ch, stats, dummyOrgId)
}

// collectMainOrgQueryStatsMetrics collects DNS query statistics for the main organization.
func (c *Collector) collectMainOrgQueryStatsMetrics(ch chan<- prometheus.Metric, org *controld.OrganizationResponse, statsEndpoint string) {
	stats, err := c.client.GetDnsQueriesReport(statsEndpoint)
	if err != nil {
		c.log.error(statsLogPrefix, errFetchingMainOrgMetrics+"%v", err)
		return
	}

	c.storeStatsMetrics(ch, stats, org.Body.Organization.PK)
}

// collectSubOrgQueryStatsMetrics collects DNS query statistics for sub organizations.
func (c *Collector) collectSubOrgQueryStatsMetrics(ch chan<- prometheus.Metric, subOrgs *controld.SubOrganizationsResponse, statsEndpoint string) {
	subOrgIDs := extractSubOrganizationIDs(subOrgs)
	for _, subOrgID := range subOrgIDs {
		stats, err := c.client.GetSubOrgDnsQueriesReport(statsEndpoint, subOrgID)
		if err != nil {
			c.log.error(statsLogPrefix, errFetchingSubOrgMetrics+"%s: %v", subOrgID, err)
			continue
		}
		c.storeStatsMetrics(ch, stats, subOrgID)
	}
}

// storeStatsMetrics stores DNS query statistics metrics in the Prometheus channel.
func (c *Collector) storeStatsMetrics(ch chan<- prometheus.Metric, stats *controld.QueryStatsResponse, orgID string) {
	if isQueryStatsEmpty(stats) {
		c.log.warn(statsLogPrefix, warnSkipEmptyData+"%v", stats)
		return
	}

	query := stats.Body.Queries[0]
	for queryType, count := range query.Count {
		queryTypeLabel := mapQueryTypeToLabel(queryType)
		ch <- prometheus.MustNewConstMetric(
			controld_stats_last_queries_count,
			prometheus.CounterValue,
			float64(count),
			queryTypeLabel,
			orgID,
		)
	}
}

// mapQueryTypeToLabel maps query types to human-readable labels.
func mapQueryTypeToLabel(queryType string) string {
	switch queryType {
	case "0":
		return "blocked"
	case "1":
		return "bypassed"
	case "3":
		return "redirected"
	default:
		return "unknown"
	}
}
