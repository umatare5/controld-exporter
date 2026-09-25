package controld

import (
	"net/http"
	"testing"
)

func TestGetNetwork(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, serveFixture(t, "network.json"))

	got, err := c.GetNetwork()
	if err != nil {
		t.Fatalf("GetNetwork: %v", err)
	}
	if len(got.Body.Network) != 3 {
		t.Fatalf("nodes = %d, want 3", len(got.Body.Network))
	}
	if got.Body.CurrentPOP != "NRT" {
		t.Errorf("current_pop = %q, want NRT", got.Body.CurrentPOP)
	}

	// Atlanta answers -1 for the proxy, which is how the API says "not offered here".
	atl := got.Body.Network[1]
	if atl.IATACode != "ATL" || atl.Status.PXY != -1 {
		t.Errorf("ATL = %+v, want the proxy at -1", atl.Status)
	}
	if atl.Location.Lat == 0 || atl.Location.Long == 0 {
		t.Errorf("ATL location = %+v, want the measured coordinates", atl.Location)
	}
}

func TestGetServiceCategories(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, serveFixture(t, "services_categories.json"))

	got, err := c.GetServiceCategories()
	if err != nil {
		t.Fatalf("GetServiceCategories: %v", err)
	}
	if len(got.Body.Categories) != 12 {
		t.Fatalf("categories = %d, want 12", len(got.Body.Categories))
	}
	if first := got.Body.Categories[0]; first.PK != "audio" || first.Count != 18 {
		t.Errorf("first category = %+v, want audio with 18 services", first)
	}
}

func TestGetProfiles(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, serveFixture(t, "profiles.json"))

	got, err := c.GetProfiles()
	if err != nil {
		t.Fatalf("GetProfiles: %v", err)
	}
	if len(got.Body.Profiles) != 1 {
		t.Fatalf("profiles = %d, want 1", len(got.Body.Profiles))
	}

	p := got.Body.Profiles[0].Profile
	// Every count feeds its own series, and cflt and ipflt differ, which is the
	// distinction the ip_filters series used to lose.
	for _, tc := range []struct {
		name string
		got  int
		want int
	}{
		{"flt", p.Flt.Count, 9},
		{"cflt", p.Cflt.Count, 0},
		{"ipflt", p.Ipflt.Count, 1},
		{"rule", p.Rule.Count, 3},
		{"svc", p.Svc.Count, 4},
		{"grp", p.Grp.Count, 0},
		{"opt", p.Opt.Count, 2},
	} {
		if tc.got != tc.want {
			t.Errorf("%s count = %d, want %d", tc.name, tc.got, tc.want)
		}
	}

	// opt values are numbers of either kind, so the field stays raw.
	if string(p.Opt.Data[0].Value) != "0.9" {
		t.Errorf("ai_malware value = %s, want 0.9", p.Opt.Data[0].Value)
	}
}

func TestGetDevices(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, serveFixture(t, "devices.json"))

	got, err := c.GetDevices()
	if err != nil {
		t.Fatalf("GetDevices: %v", err)
	}
	if len(got.Body.Devices) != 1 {
		t.Fatalf("devices = %d, want 1", len(got.Body.Devices))
	}

	d := got.Body.Devices[0]
	if d.Name != "Example Device" || d.ClientCount != 1 {
		t.Errorf("device = %q with %d clients, want Example Device with 1", d.Name, d.ClientCount)
	}
	if d.Ctrld.Version != "v1.5.0" {
		t.Errorf("ctrld version = %q, want v1.5.0", d.Ctrld.Version)
	}
}

func TestGetBillingPayments(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, serveFixture(t, "billing_payments.json"))

	got, err := c.GetBillingPayments()
	if err != nil {
		t.Fatalf("GetBillingPayments: %v", err)
	}
	if len(got.Body.Payments) != 2 {
		t.Fatalf("payments = %d, want 2", len(got.Body.Payments))
	}

	p := got.Body.Payments[0]
	// tx_status sits beside the payment; reading it from a nested object left
	// the status series at zero for every payment.
	if p.TxStatus != 1 || p.TxRefunded != 0 {
		t.Errorf("tx = %d/%d, want 1/0", p.TxStatus, p.TxRefunded)
	}
	if p.Amount != 2 || p.CurrencyAmount != 280 || p.Currency != "jpy" {
		t.Errorf("payment = %d USD / %d %s, want 2 USD / 280 jpy", p.Amount, p.CurrencyAmount, p.Currency)
	}
}

func TestGetBillingSubscriptions(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, serveFixture(t, "billing_subscriptions.json"))

	got, err := c.GetBillingSubscriptions()
	if err != nil {
		t.Fatalf("GetBillingSubscriptions: %v", err)
	}
	if len(got.Body.Subscriptions) != 1 {
		t.Fatalf("subscriptions = %d, want 1", len(got.Body.Subscriptions))
	}

	s := got.Body.Subscriptions[0]
	if s.State != "active" || s.NextBill != 1791189786 {
		t.Errorf("subscription = %q next billing at %d, want active at 1791189786", s.State, s.NextBill)
	}
}

// Billing answers 404 with the no-data code on an account that has never paid,
// which both getters absorb into an empty response rather than an error.
func TestGetBillingAbsorbsNoData(t *testing.T) {
	t.Parallel()

	handler := func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"body":[],"success":false,"error":{"code":40401}}`))
	}

	c := newTestClient(t, handler)

	payments, err := c.GetBillingPayments()
	if err != nil || len(payments.Body.Payments) != 0 {
		t.Errorf("payments = %+v, %v; want an empty response and no error", payments, err)
	}

	subs, err := c.GetBillingSubscriptions()
	if err != nil || len(subs.Body.Subscriptions) != 0 {
		t.Errorf("subscriptions = %+v, %v; want an empty response and no error", subs, err)
	}
}

func TestGetOrganizations(t *testing.T) {
	t.Parallel()

	c := newTestClient(t, serveFixture(t, "organization.json"))

	org, err := c.GetMainOrganization()
	if err != nil {
		t.Fatalf("GetMainOrganization: %v", err)
	}
	o := org.Body.Organization
	if o.PK != "000000org" || o.Members.Count != 4 || o.SubOrganizations.Count != 2 {
		t.Errorf("organization = %+v, want 4 members and 2 sub-organizations", o)
	}

	c = newTestClient(t, serveFixture(t, "sub_organizations.json"))
	subs, err := c.GetSubOrganizations()
	if err != nil {
		t.Fatalf("GetSubOrganizations: %v", err)
	}
	if len(subs.Body.SubOrganizations) != 2 {
		t.Fatalf("sub-organizations = %d, want 2", len(subs.Body.SubOrganizations))
	}
	if first := subs.Body.SubOrganizations[0]; first.PK != "000000sub1" || first.Users.Count != 5 {
		t.Errorf("first sub-organization = %+v, want 000000sub1 with 5 users", first)
	}
}

// The sub-organization getters differ from their parents only by the header,
// so one table covers every endpoint that takes one.
func TestSubOrgGettersSendTheHeader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		call func(*Client) error
		path string
	}{
		{"devices", func(c *Client) error { _, err := c.GetSubOrgDevices("org-9"); return err }, DevicesEndpoint},
		{"profiles", func(c *Client) error { _, err := c.GetSubOrgProfiles("org-9"); return err }, ProfilesEndpoint},
		{
			"service categories",
			func(c *Client) error { _, err := c.GetSubOrgServiceCategories("org-9"); return err },
			ServiceCategoriesEndpoint,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var path, org string
			c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
				path, org = r.URL.Path, r.Header.Get("X-Force-Org-Id")
				_, _ = w.Write([]byte(`{"success":true}`))
			})

			if err := tt.call(c); err != nil {
				t.Fatalf("call: %v", err)
			}
			if path != tt.path || org != "org-9" {
				t.Errorf("request = %s with org %q, want %s with org-9", path, org, tt.path)
			}
		})
	}
}
