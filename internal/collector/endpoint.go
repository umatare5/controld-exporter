// Package collector contains Prometheus metric collectors for the exporter.
package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/umatare5/controld-exporter/internal/controld"
)

const (
	endpointLogPrefix = "endpoint"
)

// collectEndpointMetrics collects endpoint-related metrics.
func (c *Collector) collectEndpointMetrics(ch chan<- prometheus.Metric) {
	if c.isRunningInPersonalMode() {
		c.collectPersonalEndpointMetrics(ch)
		c.log.debug(endpointLogPrefix, logSkipOrgScraping)
		return
	}

	// Organization metrics are only available in business mode.
	org, err := c.fetchMainOrganization()
	if err != nil {
		c.log.info(endpointLogPrefix, logNotFoundMainOrg)
		return
	}
	c.collectMainOrgEndpointMetrics(ch, org)

	subOrgs, err := c.fetchSubOrganizations()
	if err != nil {
		c.log.info(endpointLogPrefix, logNotFoundSubOrgs)
		return
	}
	c.collectSubOrgEndpointMetrics(ch, subOrgs)
}

// collectPersonalEndpointMetrics collects metrics for endpoints in the personal instance.
func (c *Collector) collectPersonalEndpointMetrics(ch chan<- prometheus.Metric) {
	endpoints, err := c.client.GetDevices()
	if err != nil {
		c.log.error(endpointLogPrefix, errFetchingPersonalMetrics+"%v", err)
		return
	}

	c.storeEndpointMetrics(ch, endpoints, dummyOrgId)
}

// collectMainOrgEndpointMetrics collects metrics for endpoints in the main organization.
func (c *Collector) collectMainOrgEndpointMetrics(ch chan<- prometheus.Metric, org *controld.OrganizationResponse) {
	endpoints, err := c.client.GetDevices()
	if err != nil {
		c.log.error(endpointLogPrefix, errFetchingMainOrgMetrics+"%v", err)
		return
	}
	c.storeEndpointMetrics(ch, endpoints, org.Body.Organization.PK)
}

// collectSubOrgEndpointMetrics collects metrics for endpoints in sub organizations.
func (c *Collector) collectSubOrgEndpointMetrics(ch chan<- prometheus.Metric, subOrgs *controld.SubOrganizationsResponse) {
	subOrgIDs := extractSubOrganizationIDs(subOrgs)
	for _, subOrgID := range subOrgIDs {
		endpoints, err := c.client.GetSubOrgDevices(subOrgID)
		if err != nil {
			c.log.error(endpointLogPrefix, errFetchingSubOrgMetrics+"%s: %v", subOrgID, err)
			continue
		}
		c.storeEndpointMetrics(ch, endpoints, subOrgID)
	}
}

// storeEndpointMetrics stores endpoint metrics in the Prometheus channel.
func (c *Collector) storeEndpointMetrics(ch chan<- prometheus.Metric, endpoints *controld.DevicesResponse, orgID string) {
	if isDevicesEmpty(endpoints) {
		c.log.warn(endpointLogPrefix, warnSkipEmptyData+"%v", endpoints)
		return
	}

	for _, endpoint := range endpoints.Body.Devices {
		ch <- prometheus.MustNewConstMetric(
			controld_endpoint_clients_total,
			prometheus.GaugeValue,
			float64(endpoint.ClientCount),
			endpoint.Name,
			orgID,
		)
	}
}
