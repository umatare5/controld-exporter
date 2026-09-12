// Package collector contains Prometheus metric collectors for the exporter.
package collector

import (
	"github.com/umatare5/controld-exporter/internal/controld"
)

const (
	dummyOrgID = "000000000" // Placeholder for the personal instance
)

// isDevicesEmpty checks if the devices array in the response is empty.
func isDevicesEmpty(devices *controld.DevicesResponse) bool {
	return devices == nil || len(devices.Body.Devices) == 0
}

// isPaymentsEmpty checks if the payments array in the response is empty.
func isPaymentsEmpty(payments *controld.BillingPaymentsResponse) bool {
	return payments == nil || len(payments.Body.Payments) == 0
}

// isSubscriptionsEmpty checks if the subscriptions array in the response is empty.
func isSubscriptionsEmpty(subscriptions *controld.BillingSubscriptionsResponse) bool {
	return subscriptions == nil || len(subscriptions.Body.Subscriptions) == 0
}

// isServiceCategoriesEmpty checks if the service categories array in the response is empty.
func isServiceCategoriesEmpty(categories *controld.ServiceCategoriesResponse) bool {
	return categories == nil || len(categories.Body.Categories) == 0
}

// isProfilesEmpty checks if the profiles array in the response is empty.
func isProfilesEmpty(profiles *controld.ProfilesResponse) bool {
	return profiles == nil || len(profiles.Body.Profiles) == 0
}

// isQueryStatsEmpty checks if the queries array in the response is empty.
func isQueryStatsEmpty(stats *controld.QueryStatsResponse) bool {
	return stats == nil || len(stats.Body.Queries) == 0
}
