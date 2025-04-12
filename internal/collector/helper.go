// Package collector contains Prometheus metric collectors for the exporter.
package collector

import (
	"github.com/umatare5/controld-exporter/internal/controld"
)

const (
	dummyOrgId = "000000000" // Placeholder for the personal instance
)

// isDevicesEmpty checks if the devices array in the response is empty.
func isDevicesEmpty(devices *controld.DevicesResponse) bool {
	return isEmpty(devices) || isEmpty(devices.Body.Devices)
}

// isPaymentsEmpty checks if the payments array in the response is empty.
func isPaymentsEmpty(payments *controld.BillingPaymentsResponse) bool {
	return isEmpty(payments) || isEmpty(payments.Body.Payments)
}

// isSubscriptionsEmpty checks if the subscriptions array in the response is empty.
func isSubscriptionsEmpty(subscriptions *controld.BillingSubscriptionsResponse) bool {
	return isEmpty(subscriptions) || isEmpty(subscriptions.Body.Subscriptions)
}

// isServiceCategoriesEmpty checks if the service categories array in the response is empty.
func isServiceCategoriesEmpty(categories *controld.ServiceCategoriesResponse) bool {
	return isEmpty(categories) || isEmpty(categories.Body.Categories)
}

// isProfilesEmpty checks if the profiles array in the response is empty.
func isProfilesEmpty(profiles *controld.ProfilesResponse) bool {
	return isEmpty(profiles) || isEmpty(profiles.Body.Profiles)
}

// isQueryStatsEmpty checks if the devices array in the response is empty.
func isQueryStatsEmpty(stats *controld.QueryStatsResponse) bool {
	return isEmpty(stats) || isEmpty(stats.Body.Queries)
}

// isEmpty checks if the given data is empty or nil.
func isEmpty(data any) bool {
	switch v := data.(type) {
	case nil:
		return true
	case []any:
		return len(v) == 0
	case map[string]any:
		return len(v) == 0
	case string:
		return v == ""
	default:
		return false
	}
}
