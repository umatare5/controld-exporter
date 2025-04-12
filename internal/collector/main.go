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
)

// Metrics descriptions
var (
	controld_billing_status = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "billing", "status"),
		"Transaction status of billing payments. ",
		[]string{"id"},
		nil,
	)

	controld_billing_refunded_status = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "billing", "refunded"),
		"Refund status of billing payments.",
		[]string{"id"},
		nil,
	)

	controld_billing_subscription_amount_total = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "billing", "subscription_amount_total"),
		"Amount of a billing subscription in the specified currency.",
		[]string{"id", "currency"},
		nil,
	)

	controld_billing_subscription_nextbill_timestamp = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "billing", "subscription_nextbill_timestamp"),
		"Timestamp of the next billing date for a subscription.",
		[]string{"id"},
		nil,
	)

	controld_endpoint_clients_total = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "endpoint", "clients_total"),
		"Number of clients connected to a device.",
		[]string{"name", "orgId"},
		nil,
	)

	controld_network_health_code = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "network", "health_code"),
		"Health status of the network by city and service.",
		[]string{"city_name", "iata_code", "country_name", "service_name"},
		nil,
	)

	controld_profile_content_filters_total = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "profile", "content_filters_total"),
		"Number of content filters applied to the profile.",
		[]string{"name", "orgId"},
		nil,
	)

	controld_profile_enabled_option_total = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "profile", "enabled_option_total"),
		"Number of enabled options in the profile.",
		[]string{"name", "orgId"},
		nil,
	)

	controld_profile_groups_total = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "profile", "groups_total"),
		"Number of group filters applied to the profile.",
		[]string{"name", "orgId"},
		nil,
	)

	controld_profile_ip_filters_total = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "profile", "ip_filters_total"),
		"Number of IP filters applied to the profile.",
		[]string{"name", "orgId"},
		nil,
	)

	controld_profile_preset_filters_total = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "profile", "preset_filters_total"),
		"Number of preset filters applied to the profile.",
		[]string{"name", "orgId"},
		nil,
	)

	controld_profile_rules_total = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "profile", "rules_total"),
		"Number of rules applied to the profile.",
		[]string{"name", "orgId"},
		nil,
	)

	controld_profile_services_total = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "profile", "services_total"),
		"Number of service filters applied to the profile.",
		[]string{"name", "orgId"},
		nil,
	)

	controld_service_categories_total = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "service", "categories_total"),
		"Number of services in each category.",
		[]string{"name", "orgId"},
		nil,
	)

	controld_stats_last_queries_count = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "stats", "last_queries_count"),
		"Count of DNS queries by type (redirect, success, blocked).",
		[]string{"type", "orgId"},
		nil,
	)

	controld_organization_members_total = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "organization", "members_total"),
		"Number of members in an organization.",
		[]string{"name", "orgId"},
		nil,
	)

	controld_organization_profiles_total = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "organization", "profiles_total"),
		"Number of profiles in an organization.",
		[]string{"name", "orgId"},
		nil,
	)

	controld_organization_users_total = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "organization", "users_total"),
		"Number of users in an organization.",
		[]string{"name", "orgId"},
		nil,
	)

	controld_organization_routers_total = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "organization", "routers_total"),
		"Number of routers in an organization.",
		[]string{"name", "orgId"},
		nil,
	)

	controld_organization_sub_orgs_total = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "organization", "sub_orgs_total"),
		"Number of sub-organizations in an organization.",
		[]string{"name", "orgId"},
		nil,
	)

	controld_sub_organization_members_total = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "sub_organization", "members_total"),
		"Number of members in a sub-organization.",
		[]string{"name", "orgId"},
		nil,
	)

	controld_sub_organization_profiles_total = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "sub_organization", "profiles_total"),
		"Number of profiles in a sub-organization.",
		[]string{"name", "orgId"},
		nil,
	)

	controld_sub_organization_users_total = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "sub_organization", "users_total"),
		"Number of users in a sub-organization.",
		[]string{"name", "orgId"},
		nil,
	)

	controld_sub_organization_routers_total = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "sub_organization", "routers_total"),
		"Number of routers in a sub-organization.",
		[]string{"name", "orgId"},
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
	ch <- controld_billing_status
	ch <- controld_billing_refunded_status
	ch <- controld_billing_subscription_amount_total
	ch <- controld_billing_subscription_nextbill_timestamp
	ch <- controld_endpoint_clients_total
	ch <- controld_network_health_code
	ch <- controld_profile_content_filters_total
	ch <- controld_profile_enabled_option_total
	ch <- controld_profile_groups_total
	ch <- controld_profile_ip_filters_total
	ch <- controld_profile_preset_filters_total
	ch <- controld_profile_rules_total
	ch <- controld_profile_services_total
	ch <- controld_service_categories_total
	ch <- controld_stats_last_queries_count
	ch <- controld_organization_members_total
	ch <- controld_organization_profiles_total
	ch <- controld_organization_routers_total
	ch <- controld_organization_sub_orgs_total
	ch <- controld_organization_users_total
	ch <- controld_sub_organization_members_total
	ch <- controld_sub_organization_profiles_total
	ch <- controld_sub_organization_routers_total
	ch <- controld_sub_organization_users_total
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
	if c.businessModeEnabled {
		return false
	}
	return true
}
