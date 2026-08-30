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
func (l *logger) info(prefix, message string) {
	log.Infof("%s: %s", prefix, message)
}

// error wraps the log.Errorf function to include a prefix in the log messages.
func (l *logger) error(prefix, format string, args ...any) {
	log.Errorf(prefix+": "+format, args...)
}

// warnEmptyData logs the skip-empty-data warning for the given payload.
func (l *logger) warnEmptyData(prefix string, data any) {
	log.Warnf("%s: %s%v", prefix, warnSkipEmptyData, data)
}

// debugSkipOrgScraping logs the personal-mode skip message.
func (l *logger) debugSkipOrgScraping(prefix string) {
	log.Debugf("%s: %s", prefix, logSkipOrgScraping)
}
