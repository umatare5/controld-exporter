package collector

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"

	"github.com/umatare5/controld-exporter/internal/controld"
)

// fixtureDir holds the measured API bodies; the client package owns them
// because it owns the response shapes.
const fixtureDir = "../controld/testdata"

// routes answers each endpoint with a fixture name, or a bare 404 when absent.
type routes map[string]string

func newCollector(t *testing.T, r routes, businessMode bool) *Collector {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		name, ok := r[req.URL.Path]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"body":[],"success":false,"error":{"code":40401}}`))
			return
		}
		body, err := os.ReadFile(filepath.Join(fixtureDir, name))
		if err != nil {
			t.Errorf("read fixture %s: %v", name, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}))
	t.Cleanup(srv.Close)

	return NewCollector(controld.NewClientWithBaseURL(srv.URL, "test-key"), businessMode)
}

// gather runs one scrape through a registry, which is the path promhttp takes.
func gather(t *testing.T, c *Collector) string {
	t.Helper()
	reg := prometheus.NewRegistry()
	reg.MustRegister(c)

	var sb strings.Builder
	families, err := reg.Gather()
	if err != nil {
		t.Fatalf("gather: %v", err)
	}
	for _, mf := range families {
		sb.WriteString(mf.GetName())
		sb.WriteByte('\n')
	}
	return sb.String()
}

// personalRoutes is what an account without organizations answers.
var personalRoutes = routes{
	controld.NetworkEndpoint:              "network.json",
	controld.ServiceCategoriesEndpoint:    "services_categories.json",
	controld.ProfilesEndpoint:             "profiles.json",
	controld.DevicesEndpoint:              "devices.json",
	controld.BillingPaymentsEndpoint:      "billing_payments.json",
	controld.BillingSubscriptionsEndpoint: "billing_subscriptions.json",
}

func TestDescribeCoversEveryMetric(t *testing.T) {
	t.Parallel()

	ch := make(chan *prometheus.Desc, 64)
	NewCollector(nil, false).Describe(ch)
	close(ch)

	seen := map[string]bool{}
	for d := range ch {
		seen[d.String()] = true
	}
	if len(seen) != 23 {
		t.Errorf("Describe sent %d distinct descriptions, want 23", len(seen))
	}
}

func TestCollectPersonalMode(t *testing.T) {
	t.Parallel()

	c := newCollector(t, personalRoutes, false)

	got := gather(t, c)
	for _, name := range []string{
		"controld_billing_status",
		"controld_billing_refunded",
		"controld_billing_subscription_amount_total",
		"controld_billing_subscription_nextbill_timestamp",
		"controld_endpoint_clients_total",
		"controld_network_health_code",
		"controld_profile_rules_total",
		"controld_service_categories_total",
	} {
		if !strings.Contains(got, name) {
			t.Errorf("%s is absent from a personal-mode scrape", name)
		}
	}

	// Organization series need business mode, so personal mode must withhold them.
	if strings.Contains(got, "controld_organization_") {
		t.Error("organization series published in personal mode")
	}
}

func TestCollectPersonalModeValues(t *testing.T) {
	t.Parallel()

	c := newCollector(t, personalRoutes, false)
	reg := prometheus.NewRegistry()
	reg.MustRegister(c)

	// tx_status is 1 in every measured payment; a nested read left it at zero.
	const wantStatus = `
# HELP controld_billing_status Transaction status of billing payments.
# TYPE controld_billing_status gauge
controld_billing_status{id="in_000000000000000"} 1
controld_billing_status{id="in_000000000000001"} 1
`
	if err := testutil.GatherAndCompare(reg, strings.NewReader(wantStatus), "controld_billing_status"); err != nil {
		t.Error(err)
	}

	// cflt and ipflt differ in the measured profile, so the two series must too.
	const wantFilters = `
# HELP controld_profile_ip_filters_total Number of IP filters applied to the profile.
# TYPE controld_profile_ip_filters_total gauge
controld_profile_ip_filters_total{name="Example Profile",orgId="000000000"} 1
`
	if err := testutil.GatherAndCompare(
		reg, strings.NewReader(wantFilters), "controld_profile_ip_filters_total",
	); err != nil {
		t.Error(err)
	}

	// A node that does not offer the proxy reports -1 rather than withholding it.
	// The whole family is compared, because GatherAndCompare admits no subset.
	const wantHealth = `
# HELP controld_network_health_code Health status of the network by city and service.
# TYPE controld_network_health_code gauge
controld_network_health_code{city_name="Amsterdam",country_name="NL",iata_code="AMS",service_name="api"} 1
controld_network_health_code{city_name="Amsterdam",country_name="NL",iata_code="AMS",service_name="dns"} 1
controld_network_health_code{city_name="Amsterdam",country_name="NL",iata_code="AMS",service_name="proxy"} 1
controld_network_health_code{city_name="Atlanta",country_name="US",iata_code="ATL",service_name="api"} 1
controld_network_health_code{city_name="Atlanta",country_name="US",iata_code="ATL",service_name="dns"} 1
controld_network_health_code{city_name="Atlanta",country_name="US",iata_code="ATL",service_name="proxy"} -1
controld_network_health_code{city_name="Tokyo",country_name="JP",iata_code="NRT",service_name="api"} 1
controld_network_health_code{city_name="Tokyo",country_name="JP",iata_code="NRT",service_name="dns"} 1
controld_network_health_code{city_name="Tokyo",country_name="JP",iata_code="NRT",service_name="proxy"} 1
`
	if err := testutil.GatherAndCompare(
		reg, strings.NewReader(wantHealth), "controld_network_health_code",
	); err != nil {
		t.Error(err)
	}
}

func TestCollectBusinessMode(t *testing.T) {
	t.Parallel()

	r := routes{controld.OrganizationEndpoint: "organization.json"}
	for k, v := range personalRoutes {
		r[k] = v
	}
	r[controld.SubOrganizationsEndpoint] = "sub_organizations.json"

	c := newCollector(t, r, true)

	got := gather(t, c)
	for _, name := range []string{
		"controld_organization_members_total",
		"controld_organization_sub_orgs_total",
		"controld_sub_organization_members_total",
		"controld_sub_organization_routers_total",
	} {
		if !strings.Contains(got, name) {
			t.Errorf("%s is absent from a business-mode scrape", name)
		}
	}
}

// An organization-scoped read that fails must leave the scrape standing: the
// exporter answers with whatever the other collectors produced.
func TestCollectSurvivesOrganizationFailure(t *testing.T) {
	t.Parallel()

	c := newCollector(t, personalRoutes, true) // business mode, no organization routes

	got := gather(t, c)
	if !strings.Contains(got, "controld_network_health_code") {
		t.Error("network series lost when the organization read failed")
	}
	if strings.Contains(got, "controld_organization_") {
		t.Error("organization series published although the read failed")
	}
}

// Sub-organization reads are per-organization, so one failure must not stop the rest.
func TestCollectSurvivesSubOrgFailure(t *testing.T) {
	t.Parallel()

	r := routes{
		controld.OrganizationEndpoint:      "organization.json",
		controld.SubOrganizationsEndpoint:  "sub_organizations.json",
		controld.NetworkEndpoint:           "network.json",
		controld.ServiceCategoriesEndpoint: "services_categories.json",
	}
	c := newCollector(t, r, true)

	got := gather(t, c)
	if !strings.Contains(got, "controld_organization_members_total") {
		t.Error("main organization series lost when a sub-organization read failed")
	}
}

// Every collector guards an empty collection, which the API answers with a 200
// carrying an empty array rather than an error.
func TestCollectEmptyCollections(t *testing.T) {
	t.Parallel()

	empty := map[string]string{
		controld.NetworkEndpoint:              `{"success":true,"body":{"network":[]}}`,
		controld.ServiceCategoriesEndpoint:    `{"success":true,"body":{"categories":[]}}`,
		controld.ProfilesEndpoint:             `{"success":true,"body":{"profiles":[]}}`,
		controld.DevicesEndpoint:              `{"success":true,"body":{"devices":[]}}`,
		controld.BillingPaymentsEndpoint:      `{"success":true,"body":{"payments":[]}}`,
		controld.BillingSubscriptionsEndpoint: `{"success":true,"body":{"subscriptions":[]}}`,
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		body, ok := empty[req.URL.Path]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)

	c := NewCollector(controld.NewClientWithBaseURL(srv.URL, "k"), false)
	if got := gather(t, c); got != "" {
		t.Errorf("families = %q, want none when every collection is empty", got)
	}
}

// A dead endpoint is a transport error rather than a status, and it must be
// logged and skipped rather than panicking the scrape.
func TestCollectSurvivesTransportFailure(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.NotFoundHandler())
	srv.Close()

	c := NewCollector(controld.NewClientWithBaseURL(srv.URL, "k"), false)
	if got := gather(t, c); got != "" {
		t.Errorf("families = %q, want none when the API is unreachable", got)
	}
}

func TestIsRunningInPersonalMode(t *testing.T) {
	t.Parallel()

	if !NewCollector(nil, false).isRunningInPersonalMode() {
		t.Error("business mode off must read as personal mode")
	}
	if NewCollector(nil, true).isRunningInPersonalMode() {
		t.Error("business mode on must not read as personal mode")
	}
}

func TestFetchOrganizationCaches(t *testing.T) {
	t.Parallel()

	var calls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls++
		_, _ = w.Write([]byte(`{"success":true,"body":{"organization":{"PK":"o"}}}`))
	}))
	t.Cleanup(srv.Close)

	c := NewCollector(controld.NewClientWithBaseURL(srv.URL, "k"), true)
	for range 3 {
		if _, err := c.fetchMainOrganization(); err != nil {
			t.Fatalf("fetchMainOrganization: %v", err)
		}
	}
	if calls != 1 {
		t.Errorf("organization fetched %d times, want 1 across a scrape", calls)
	}
}

func TestExtractSubOrganizationIDs(t *testing.T) {
	t.Parallel()

	// The response type nests anonymous structs, so the fixture builds it.
	body, err := os.ReadFile(filepath.Join(fixtureDir, "sub_organizations.json"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	var orgs controld.SubOrganizationsResponse
	if err := json.Unmarshal(body, &orgs); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	got := extractSubOrganizationIDs(&orgs)
	if len(got) != 2 || got[0] != "000000sub1" || got[1] != "000000sub2" {
		t.Errorf("ids = %v, want [000000sub1 000000sub2]", got)
	}
}
