// Package collector contains Prometheus metric collectors for the exporter.
package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/umatare5/controld-exporter/internal/controld"
)

const (
	profileLogPrefix = "profile"
)

// collectProfileMetrics collects profile-related metrics.
func (c *Collector) collectProfileMetrics(ch chan<- prometheus.Metric) {
	if c.isRunningInPersonalMode() {
		c.collectPersonalProfileMetrics(ch)
		c.log.debug(profileLogPrefix, logSkipOrgScraping)
		return
	}

	// Organization metrics are only available in business mode.
	org, err := c.fetchMainOrganization()
	if err != nil {
		c.log.info(profileLogPrefix, logNotFoundMainOrg)
		return
	}
	c.collectMainOrgProfileMetrics(ch, org)

	subOrgs, err := c.fetchSubOrganizations()
	if err != nil {
		c.log.info(profileLogPrefix, logNotFoundSubOrgs)
		return
	}
	c.collectSubOrgProfileMetrics(ch, subOrgs)
}

// collectPersonalProfileMetrics collects metrics for profiles in the personal instance.
func (c *Collector) collectPersonalProfileMetrics(ch chan<- prometheus.Metric) {
	profiles, err := c.client.GetProfiles()
	if err != nil {
		c.log.error(profileLogPrefix, errFetchingPersonalMetrics+"%v", err)
		return
	}

	c.storeProfileMetrics(ch, profiles, dummyOrgId)
}

// collectMainOrgProfileMetrics collects metrics for profiles in the main organization.
func (c *Collector) collectMainOrgProfileMetrics(ch chan<- prometheus.Metric, org *controld.OrganizationResponse) {
	profiles, err := c.client.GetProfiles()
	if err != nil {
		c.log.error(profileLogPrefix, errFetchingMainOrgMetrics+"%v", err)
		return
	}
	c.storeProfileMetrics(ch, profiles, org.Body.Organization.PK)
}

// collectSubOrgProfileMetrics collects metrics for profiles in sub organizations.
func (c *Collector) collectSubOrgProfileMetrics(ch chan<- prometheus.Metric, orgs *controld.SubOrganizationsResponse) {
	subOrgIDs := extractSubOrganizationIDs(orgs)
	for _, subOrgID := range subOrgIDs {
		profiles, err := c.client.GetSubOrgProfiles(subOrgID)
		if err != nil {
			c.log.error(profileLogPrefix, errFetchingSubOrgMetrics+"%s: %v", subOrgID, err)
			continue
		}
		c.storeProfileMetrics(ch, profiles, subOrgID)
	}
}

// storeProfileMetrics stores profile metrics in the Prometheus channel.
func (c *Collector) storeProfileMetrics(ch chan<- prometheus.Metric, profiles *controld.ProfilesResponse, orgID string) {
	if isProfilesEmpty(profiles) {
		c.log.warn(profileLogPrefix, warnSkipEmptyData+"%v", profiles)
		return
	}

	for _, profile := range profiles.Body.Profiles {
		ch <- prometheus.MustNewConstMetric(
			controld_profile_preset_filters_total,
			prometheus.GaugeValue,
			float64(profile.Profile.Flt.Count),
			profile.Name,
			orgID,
		)
		ch <- prometheus.MustNewConstMetric(
			controld_profile_content_filters_total,
			prometheus.GaugeValue,
			float64(profile.Profile.Cflt.Count),
			profile.Name,
			orgID,
		)
		ch <- prometheus.MustNewConstMetric(
			controld_profile_ip_filters_total,
			prometheus.GaugeValue,
			float64(profile.Profile.Cflt.Count),
			profile.Name,
			orgID,
		)
		ch <- prometheus.MustNewConstMetric(
			controld_profile_rules_total,
			prometheus.GaugeValue,
			float64(profile.Profile.Rule.Count),
			profile.Name,
			orgID,
		)
		ch <- prometheus.MustNewConstMetric(
			controld_profile_services_total,
			prometheus.GaugeValue,
			float64(profile.Profile.Svc.Count),
			profile.Name,
			orgID,
		)
		ch <- prometheus.MustNewConstMetric(
			controld_profile_groups_total,
			prometheus.GaugeValue,
			float64(profile.Profile.Grp.Count),
			profile.Name,
			orgID,
		)
		ch <- prometheus.MustNewConstMetric(
			controld_profile_enabled_option_total,
			prometheus.GaugeValue,
			float64(profile.Profile.Opt.Count),
			profile.Name,
			orgID,
		)
	}
}
