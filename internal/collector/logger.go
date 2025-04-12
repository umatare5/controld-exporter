package collector

import "github.com/umatare5/controld-exporter/internal/log"

const (
	logSkipOrgScraping         = "Running in personal mode. Skipping the scraping metrics of the organizations."
	logNotFoundMainOrg         = "Not found main organization. Skipping the scraping metrics of the main organization."
	logNotFoundSubOrgs         = "Not found sub organizations. Skipping the scraping metrics of the sub organizations."
	errFetchingMetrics         = "Error fetching metrics: "
	errFetchingPersonalMetrics = "Error fetching metrics for personal instance: "
	errFetchingMainOrgMetrics  = "Error fetching metrics for main organization: "
	errFetchingSubOrgMetrics   = "Error fetching metrics for sub organization ID: "
	warnSkipEmptyData          = "Skipping empty data: "
)

type logger struct{}

// info wraps the log.Infof function to include a prefix in the log messages.
func (l *logger) info(prefix string, format string, args ...interface{}) {
	log.Infof(prefix+": "+format, args...)
}

// error wraps the log.Errorf function to include a prefix in the log messages.
func (l *logger) error(prefix string, format string, args ...interface{}) {
	log.Errorf(prefix+": "+format, args...)
}

// warn wraps the log.Warnf function to include a prefix in the log messages.
func (l *logger) warn(prefix string, format string, args ...interface{}) {
	log.Warnf(prefix+": "+format, args...)
}

// debug wraps the log.Debugf function to include a prefix in the log messages.
func (l *logger) debug(prefix string, format string, args ...interface{}) {
	log.Debugf(prefix+": "+format, args...)
}
