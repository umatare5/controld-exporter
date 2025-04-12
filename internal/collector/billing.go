// Package collector contains Prometheus metric collectors for the exporter.
package collector

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	billingPaymentsLogPrefix      = "billingPayments"
	billingSubscriptionsLogPrefix = "billingSubscriptions"
)

// collectBillingMetrics collects billing-related metrics.
func (c *Collector) collectBillingMetrics(ch chan<- prometheus.Metric) {
	c.collectBillingPayments(ch)
	c.collectBillingSubscriptions(ch)
}

// collectBillingPayments collects metrics for billing payments.
func (c *Collector) collectBillingPayments(ch chan<- prometheus.Metric) {
	payments, err := c.client.GetBillingPayments()
	if err != nil {
		c.log.error(billingPaymentsLogPrefix, errFetchingMetrics+"%v", err)
		return
	}

	if isPaymentsEmpty(payments) {
		c.log.warn(billingPaymentsLogPrefix, warnSkipEmptyData+"%v", payments)
		return
	}

	for _, payment := range payments.Body.Payments {
		ch <- prometheus.MustNewConstMetric(
			controld_billing_status,
			prometheus.GaugeValue,
			float64(payment.Transaction.Status),
			payment.PK,
		)

		ch <- prometheus.MustNewConstMetric(
			controld_billing_refunded_status,
			prometheus.GaugeValue,
			float64(payment.Transaction.Refunded),
			payment.PK,
		)

		ch <- prometheus.MustNewConstMetric(
			controld_billing_subscription_amount_total,
			prometheus.GaugeValue,
			float64(payment.Amount),
			payment.PK,
			"USD",
		)

		ch <- prometheus.MustNewConstMetric(
			controld_billing_subscription_amount_total,
			prometheus.GaugeValue,
			float64(payment.CurrencyAmount),
			payment.PK,
			strings.ToUpper(payment.Currency),
		)
	}
}

// collectBillingSubscriptions collects metrics for billing subscriptions.
func (c *Collector) collectBillingSubscriptions(ch chan<- prometheus.Metric) {
	subscriptions, err := c.client.GetBillingSubscriptions()
	if err != nil {
		c.log.error(billingSubscriptionsLogPrefix, errFetchingMetrics+"%v", err)
		return
	}

	if isSubscriptionsEmpty(subscriptions) {
		c.log.warn(billingPaymentsLogPrefix, warnSkipEmptyData+"%v", subscriptions)
		return
	}

	for _, subscription := range subscriptions.Body.Subscriptions {
		ch <- prometheus.MustNewConstMetric(
			controld_billing_subscription_nextbill_timestamp,
			prometheus.GaugeValue,
			float64(subscription.NextBill),
			subscription.PK,
		)
	}
}
