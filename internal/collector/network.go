// Package collector contains Prometheus metric collectors for the exporter.
package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	networkHealthLogPrefix = "networkHealth"
)

// collectNetworkMetrics collects all network-related metrics.
func (c *Collector) collectNetworkMetrics(ch chan<- prometheus.Metric) {
	c.collectNetworkHealthStatus(ch)
}

// collectNetworkHealthStatus collects metrics for network nodes.
func (c *Collector) collectNetworkHealthStatus(ch chan<- prometheus.Metric) {
	network, err := c.client.GetNetwork()
	if err != nil {
		c.log.error(networkHealthLogPrefix, errFetchingMetrics+"%v", err)
		return
	}

	for _, node := range network.Body.Network {
		ch <- prometheus.MustNewConstMetric(
			controld_network_health_code,
			prometheus.GaugeValue,
			float64(node.Status.API),
			node.CityName,
			node.IATACode,
			node.CountryName,
			"api",
		)

		ch <- prometheus.MustNewConstMetric(
			controld_network_health_code,
			prometheus.GaugeValue,
			float64(node.Status.DNS),
			node.CityName,
			node.IATACode,
			node.CountryName,
			"dns",
		)

		ch <- prometheus.MustNewConstMetric(
			controld_network_health_code,
			prometheus.GaugeValue,
			float64(node.Status.PXY),
			node.CityName,
			node.IATACode,
			node.CountryName,
			"proxy",
		)
	}
}
