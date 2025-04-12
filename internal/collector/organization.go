// Package collector contains Prometheus metric collectors for the exporter.
package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/umatare5/controld-exporter/internal/controld"
)

const (
	organizationLogPrefix    = "organization"
	subOrganizationLogPrefix = "sub-organization"
)

// collectOrganizationMetrics collects organization-related metrics.
func (c *Collector) collectOrganizationMetrics(ch chan<- prometheus.Metric) {
	if c.isRunningInPersonalMode() {
		c.log.debug(profileLogPrefix, logSkipOrgScraping)
		return
	}

	// organization metrics are only available in business mode.
	org, err := c.fetchMainOrganization()
	c.collectMainOrganizationMetrics(ch, org)
	if err != nil {
		c.log.error(organizationLogPrefix, errFetchingMainOrgMetrics+"%v", err)
		return
	}

	subOrgs, err := c.fetchSubOrganizations()
	c.collectSubOrganizationMetrics(ch, subOrgs)
	if err != nil {
		c.log.error(subOrganizationLogPrefix, errFetchingSubOrgMetrics+"%v", err)
		return
	}
}

// collectMainOrganizationMetrics collects metrics for main organization.
func (c *Collector) collectMainOrganizationMetrics(ch chan<- prometheus.Metric, org *controld.OrganizationResponse) {
	ch <- prometheus.MustNewConstMetric(
		controld_organization_members_total,
		prometheus.GaugeValue,
		float64(org.Body.Organization.Members.Count),
		org.Body.Organization.Name,
		org.Body.Organization.PK,
	)
	ch <- prometheus.MustNewConstMetric(
		controld_organization_profiles_total,
		prometheus.GaugeValue,
		float64(org.Body.Organization.Profiles.Count),
		org.Body.Organization.Name,
		org.Body.Organization.PK,
	)
	ch <- prometheus.MustNewConstMetric(
		controld_organization_users_total,
		prometheus.GaugeValue,
		float64(org.Body.Organization.Users.Count),
		org.Body.Organization.Name,
		org.Body.Organization.PK,
	)
	ch <- prometheus.MustNewConstMetric(
		controld_organization_routers_total,
		prometheus.GaugeValue,
		float64(org.Body.Organization.Routers.Count),
		org.Body.Organization.Name,
		org.Body.Organization.PK,
	)
	ch <- prometheus.MustNewConstMetric(
		controld_organization_sub_orgs_total,
		prometheus.GaugeValue,
		float64(org.Body.Organization.SubOrganizations.Count),
		org.Body.Organization.Name,
		org.Body.Organization.PK,
	)
}

// collectSubOrganizationMetrics collects metrics for sub organizations.
func (c *Collector) collectSubOrganizationMetrics(ch chan<- prometheus.Metric, subOrgs *controld.SubOrganizationsResponse) {
	for _, subOrg := range subOrgs.Body.SubOrganizations {
		ch <- prometheus.MustNewConstMetric(
			controld_sub_organization_members_total,
			prometheus.GaugeValue,
			float64(subOrg.Members.Count),
			subOrg.Name,
			subOrg.PK,
		)
		ch <- prometheus.MustNewConstMetric(
			controld_sub_organization_profiles_total,
			prometheus.GaugeValue,
			float64(subOrg.Profiles.Count),
			subOrg.Name,
			subOrg.PK,
		)
		ch <- prometheus.MustNewConstMetric(
			controld_sub_organization_users_total,
			prometheus.GaugeValue,
			float64(subOrg.Users.Count),
			subOrg.Name,
			subOrg.PK,
		)
		ch <- prometheus.MustNewConstMetric(
			controld_sub_organization_routers_total,
			prometheus.GaugeValue,
			float64(subOrg.Routers.Count),
			subOrg.Name,
			subOrg.PK,
		)
	}
}

// fetchMainOrganization fetches and caches main organization data.
func (c *Collector) fetchMainOrganization() (*controld.OrganizationResponse, error) {
	c.organizationsMu.Lock()
	defer c.organizationsMu.Unlock()

	// Return cached data if already fetched
	if c.organizations != nil {
		return c.organizations, nil
	}

	// Fetch organization data from the API
	orgs, err := c.client.GetMainOrganization()
	if err != nil {
		return nil, err
	}

	c.organizations = orgs
	return orgs, nil
}

// fetchSubrganizations fetches and caches sub organization data.
func (c *Collector) fetchSubOrganizations() (*controld.SubOrganizationsResponse, error) {
	c.subOrganizationsMu.Lock()
	defer c.subOrganizationsMu.Unlock()

	// Return cached data if already fetched
	if c.subOrganizations != nil {
		return c.subOrganizations, nil
	}

	// Fetch organization data from the API
	orgs, err := c.client.GetSubOrganizations()
	if err != nil {
		return nil, err
	}

	c.subOrganizations = orgs
	return orgs, nil
}

// extractSubOrganizationIDs extracts sub-organization IDs from the response.
func extractSubOrganizationIDs(orgs *controld.SubOrganizationsResponse) []string {
	subOrgs := make([]string, len(orgs.Body.SubOrganizations))
	for i, subOrg := range orgs.Body.SubOrganizations {
		subOrgs[i] = subOrg.PK
	}
	return subOrgs
}
