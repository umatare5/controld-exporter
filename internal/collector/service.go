// Package collector contains Prometheus metric collectors for the exporter.
package collector

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/umatare5/controld-exporter/internal/controld"
)

const (
	serviceLogPrefix = "service"
)

// collectServiceMetrics collects service-related metrics.
func (c *Collector) collectServiceMetrics(ch chan<- prometheus.Metric) {
	if c.isRunningInPersonalMode() {
		c.collectPersonalServicesCategoryMetrics(ch)
		c.log.debugSkipOrgScraping(serviceLogPrefix)
		return
	}

	// Organization metrics are only available in business mode.
	org, err := c.fetchMainOrganization()
	if err != nil {
		c.log.info(serviceLogPrefix, logNotFoundMainOrg)
		return
	}
	c.collectMainOrgServicesCategoryMetrics(ch, org)

	subOrgs, err := c.fetchSubOrganizations()
	if err != nil {
		c.log.info(serviceLogPrefix, logNotFoundSubOrgs)
		return
	}
	c.collectSubOrgServicesCategoryMetrics(ch, subOrgs)
}

// collectPersonalServicesCategoryMetrics collects metrics for serviceCategories in the personal instance.
func (c *Collector) collectPersonalServicesCategoryMetrics(ch chan<- prometheus.Metric) {
	serviceCategories, err := c.client.GetServiceCategories()
	if err != nil {
		c.log.error(serviceLogPrefix, errFetchingPersonalMetrics+"%v", err)
		return
	}

	c.storeServicesCategoryMetrics(ch, serviceCategories, dummyOrgID)
}

// collectMainOrgServicesCategoryMetrics collects metrics for serviceCategories in the main organization.
func (c *Collector) collectMainOrgServicesCategoryMetrics(
	ch chan<- prometheus.Metric,
	org *controld.OrganizationResponse,
) {
	serviceCategories, err := c.client.GetServiceCategories()
	if err != nil {
		c.log.error(serviceLogPrefix, errFetchingMainOrgMetrics+"%v", err)
		return
	}

	c.storeServicesCategoryMetrics(ch, serviceCategories, org.Body.Organization.PK)
}

// collectSubOrgServicesCategoryMetrics collects metrics for serviceCategories in sub organizations.
func (c *Collector) collectSubOrgServicesCategoryMetrics(
	ch chan<- prometheus.Metric,
	subOrgs *controld.SubOrganizationsResponse,
) {
	subOrgIDs := extractSubOrganizationIDs(subOrgs)
	for _, subOrgID := range subOrgIDs {
		serviceCategories, err := c.client.GetSubOrgServiceCategories(subOrgID)
		if err != nil {
			c.log.error(serviceLogPrefix, errFetchingSubOrgMetrics+"%s: %v", subOrgID, err)
			continue
		}
		c.storeServicesCategoryMetrics(ch, serviceCategories, subOrgID)
	}
}

// storeServicesCategoryMetrics stores servicesCategory metrics in the Prometheus channel.
func (c *Collector) storeServicesCategoryMetrics(
	ch chan<- prometheus.Metric,
	serviceCategories *controld.ServiceCategoriesResponse,
	orgID string,
) {
	if isServiceCategoriesEmpty(serviceCategories) {
		c.log.warnEmptyData(serviceLogPrefix, serviceCategories)
		return
	}

	for _, servicesCategory := range serviceCategories.Body.Categories {
		ch <- prometheus.MustNewConstMetric(
			controldServiceCategoriesTotal,
			prometheus.GaugeValue,
			float64(servicesCategory.Count),
			servicesCategory.PK,
			orgID,
		)
	}
}
