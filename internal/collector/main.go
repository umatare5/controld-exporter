// Package collector contains Prometheus metric collectors for the exporter.
package collector

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/umatare5/controld-exporter/internal/controld"
)

const (
	namespace = "controld"
	subsystem = ""

	labelID    = "id"
	labelName  = "name"
	labelOrgID = "orgId"
)

// Metrics descriptions.
var (
	controldBillingStatus = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "billing", "status"),
		"Transaction status of billing payments. ",
		[]string{labelID},
		nil,
	)

	controldBillingRefundedStatus = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "billing", "refunded"),
		"Refund status of billing payments.",
		[]string{labelID},
		nil,
	)

	controldBillingSubscriptionAmountTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "billing", "subscription_amount_total"),
		"Amount of a billing subscription in the specified currency.",
		[]string{labelID, "currency"},
		nil,
	)

	controldBillingSubscriptionNextbillTimestamp = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "billing", "subscription_nextbill_timestamp"),
		"Timestamp of the next billing date for a subscription.",
		[]string{labelID},
		nil,
	)

	controldEndpointClientsTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "endpoint", "clients_total"),
		"Number of clients connected to a device.",
		[]string{labelName, labelOrgID},
		nil,
	)

	controldNetworkHealthCode = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "network", "health_code"),
		"Health status of the network by city and service.",
		[]string{"city_name", "iata_code", "country_name", "service_name"},
		nil,
	)

	controldProfileContentFiltersTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "profile", "content_filters_total"),
		"Number of content filters applied to the profile.",
		[]string{labelName, labelOrgID},
		nil,
	)

	controldProfileEnabledOptionTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "profile", "enabled_option_total"),
		"Number of enabled options in the profile.",
		[]string{labelName, labelOrgID},
		nil,
	)

	controldProfileGroupsTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "profile", "groups_total"),
		"Number of group filters applied to the profile.",
		[]string{labelName, labelOrgID},
		nil,
	)

	controldProfileIPFiltersTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "profile", "ip_filters_total"),
		"Number of IP filters applied to the profile.",
		[]string{labelName, labelOrgID},
		nil,
	)

	controldProfilePresetFiltersTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "profile", "preset_filters_total"),
		"Number of preset filters applied to the profile.",
		[]string{labelName, labelOrgID},
		nil,
	)

	controldProfileRulesTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "profile", "rules_total"),
		"Number of rules applied to the profile.",
		[]string{labelName, labelOrgID},
		nil,
	)

	controldProfileServicesTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "profile", "services_total"),
		"Number of service filters applied to the profile.",
		[]string{labelName, labelOrgID},
		nil,
	)

	controldServiceCategoriesTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "service", "categories_total"),
		"Number of services in each category.",
		[]string{labelName, labelOrgID},
		nil,
	)

	controldStatsLastQueriesCount = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "stats", "last_queries_count"),
		"Count of DNS queries by type (redirect, success, blocked).",
		[]string{"type", labelOrgID},
		nil,
	)

	controldOrganizationMembersTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "organization", "members_total"),
		"Number of members in an organization.",
		[]string{labelName, labelOrgID},
		nil,
	)

	controldOrganizationProfilesTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "organization", "profiles_total"),
		"Number of profiles in an organization.",
		[]string{labelName, labelOrgID},
		nil,
	)

	controldOrganizationUsersTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "organization", "users_total"),
		"Number of users in an organization.",
		[]string{labelName, labelOrgID},
		nil,
	)

	controldOrganizationRoutersTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "organization", "routers_total"),
		"Number of routers in an organization.",
		[]string{labelName, labelOrgID},
		nil,
	)

	controldOrganizationSubOrgsTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "organization", "sub_orgs_total"),
		"Number of sub-organizations in an organization.",
		[]string{labelName, labelOrgID},
		nil,
	)

	controldSubOrganizationMembersTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "sub_organization", "members_total"),
		"Number of members in a sub-organization.",
		[]string{labelName, labelOrgID},
		nil,
	)

	controldSubOrganizationProfilesTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "sub_organization", "profiles_total"),
		"Number of profiles in a sub-organization.",
		[]string{labelName, labelOrgID},
		nil,
	)

	controldSubOrganizationUsersTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "sub_organization", "users_total"),
		"Number of users in a sub-organization.",
		[]string{labelName, labelOrgID},
		nil,
	)

	controldSubOrganizationRoutersTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "sub_organization", "routers_total"),
		"Number of routers in a sub-organization.",
		[]string{labelName, labelOrgID},
		nil,
	)
)

// Collector is responsible for collecting metrics from ControlD.
type Collector struct {
	client              *controld.Client                   // ControlD API client
	organizations       *controld.OrganizationResponse     // Cached organization data
	organizationsMu     sync.Mutex                         // Mutex to protect access to the cached data for main-organization
	subOrganizations    *controld.SubOrganizationsResponse // Cached sub-organization data
	subOrganizationsMu  sync.Mutex                         // Mutex to protect access to the cached data for sub-organization
	businessModeEnabled bool                               // Indicates if business features is enabled
	log                 *logger
}

// NewCollector initializes and returns a new Collector instance.
func NewCollector(client *controld.Client, businessMode bool) *Collector {
	return &Collector{
		client:              client,
		businessModeEnabled: businessMode,
	}
}

// Describe sends the descriptions of all metrics to the Prometheus channel.
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- controldBillingStatus
	ch <- controldBillingRefundedStatus
	ch <- controldBillingSubscriptionAmountTotal
	ch <- controldBillingSubscriptionNextbillTimestamp
	ch <- controldEndpointClientsTotal
	ch <- controldNetworkHealthCode
	ch <- controldProfileContentFiltersTotal
	ch <- controldProfileEnabledOptionTotal
	ch <- controldProfileGroupsTotal
	ch <- controldProfileIPFiltersTotal
	ch <- controldProfilePresetFiltersTotal
	ch <- controldProfileRulesTotal
	ch <- controldProfileServicesTotal
	ch <- controldServiceCategoriesTotal
	ch <- controldStatsLastQueriesCount
	ch <- controldOrganizationMembersTotal
	ch <- controldOrganizationProfilesTotal
	ch <- controldOrganizationRoutersTotal
	ch <- controldOrganizationSubOrgsTotal
	ch <- controldOrganizationUsersTotal
	ch <- controldSubOrganizationMembersTotal
	ch <- controldSubOrganizationProfilesTotal
	ch <- controldSubOrganizationRoutersTotal
	ch <- controldSubOrganizationUsersTotal
}

// Collect gathers metrics from ControlD and sends them to the Prometheus channel.
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.collectOrganizationMetrics(ch)
	c.collectBillingMetrics(ch)
	c.collectEndpointMetrics(ch)
	c.collectNetworkMetrics(ch)
	c.collectProfileMetrics(ch)
	c.collectServiceMetrics(ch)
	c.collectStatsMetrics(ch)
}

// isRunningInPersonalMode checks if the collector is running in personal mode.
func (c *Collector) isRunningInPersonalMode() bool {
	return !c.businessModeEnabled
}
