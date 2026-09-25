package controld

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// newTestClient points a client at srv, which is the only seam the package
// needs: baseURL is a field, so no production code changes for tests.
func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return &Client{baseURL: srv.URL, apiKey: "test-key"}
}

// serveFixture answers every request with the named testdata file.
func serveFixture(t *testing.T, name string) http.HandlerFunc {
	t.Helper()
	body, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatalf("read fixture %s: %v", name, err)
	}
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(body)
	}
}

func TestNewClient(t *testing.T) {
	t.Parallel()

	c := NewClient("key")
	if c.baseURL != "https://api.controld.com" {
		t.Errorf("baseURL = %q, want the vendor host", c.baseURL)
	}
	if c.apiKey != "key" {
		t.Errorf("apiKey = %q, want the key it was built with", c.apiKey)
	}
}

func TestCreateRequestSetsAuthAndHeaders(t *testing.T) {
	t.Parallel()

	var got *http.Request
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		got = r.Clone(r.Context())
		_, _ = w.Write([]byte(`{"success":true}`))
	})

	if _, err := c.GetSubOrgProfiles("org-1"); err != nil {
		t.Fatalf("GetSubOrgProfiles: %v", err)
	}
	if want := "Bearer test-key"; got.Header.Get("Authorization") != want {
		t.Errorf("Authorization = %q, want %q", got.Header.Get("Authorization"), want)
	}
	if got.Header.Get("Accept") != "application/json" {
		t.Errorf("Accept = %q, want application/json", got.Header.Get("Accept"))
	}
	if got.Header.Get("X-Force-Org-Id") != "org-1" {
		t.Errorf("X-Force-Org-Id = %q, want org-1", got.Header.Get("X-Force-Org-Id"))
	}
	if got.URL.Path != ProfilesEndpoint {
		t.Errorf("path = %q, want %q", got.URL.Path, ProfilesEndpoint)
	}
}

func TestHandleResponseErrors(t *testing.T) {
	t.Parallel()

	noData, err := os.ReadFile(filepath.Join("testdata", "error_no_data.json"))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	tests := []struct {
		name    string
		status  int
		body    string
		wantErr error
		wantMsg string
	}{
		{name: "no data", status: http.StatusNotFound, body: string(noData), wantErr: ErrNoData},
		{
			name:    "forbidden",
			status:  http.StatusForbidden,
			body:    `{"error":{"code":40301,"message":"forbidden"}}`,
			wantMsg: `unexpected status "403 Forbidden"`,
		},
		{
			name:    "success false",
			status:  http.StatusOK,
			body:    `{"success":false}`,
			wantMsg: "API response indicates failure",
		},
		{name: "malformed", status: http.StatusOK, body: `{`, wantMsg: "unexpected end of JSON input"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			})

			_, err := c.GetProfiles()
			if err == nil {
				t.Fatal("want an error, got nil")
			}
			if tt.wantErr != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantMsg != "" && !strings.Contains(err.Error(), tt.wantMsg) {
				t.Errorf("error = %v, want it to mention %q", err, tt.wantMsg)
			}
		})
	}
}

func TestSendRequestTransportError(t *testing.T) {
	t.Parallel()

	// A closed listener gives a dial error, the one failure sendRequest owns.
	srv := httptest.NewServer(http.NotFoundHandler())
	c := &Client{baseURL: srv.URL, apiKey: "k"}
	srv.Close()

	if _, err := c.GetNetwork(); err == nil {
		t.Fatal("want a transport error, got nil")
	}
}

func TestSummarizeErrorBodyBounds(t *testing.T) {
	t.Parallel()

	got := summarizeErrorBody([]byte("a\n  b\tc"))
	if got != "a b c" {
		t.Errorf("summarizeErrorBody collapsed to %q, want %q", got, "a b c")
	}

	long := summarizeErrorBody([]byte(strings.Repeat("x", maxErrorBodyLen+10)))
	if len(long) != maxErrorBodyLen+3 || !strings.HasSuffix(long, "...") {
		t.Errorf("long body = %d chars, want %d and an ellipsis", len(long), maxErrorBodyLen+3)
	}
}

func TestIsNoDataBody(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		body string
		want bool
	}{
		{name: "no data code", body: `{"error":{"code":40401}}`, want: true},
		{name: "other code", body: `{"error":{"code":40301}}`, want: false},
		{name: "not json", body: `<html>`, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := isNoDataBody([]byte(tt.body)); got != tt.want {
				t.Errorf("isNoDataBody = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSuccess(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		body string
		want bool
	}{
		{name: "true", body: `{"success":true}`, want: true},
		{name: "false", body: `{"success":false}`, want: false},
		{name: "absent", body: `{}`, want: false},
		{name: "wrong type", body: `{"success":1}`, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var m map[string]any
			if err := json.Unmarshal([]byte(tt.body), &m); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if got := isSuccess(m); got != tt.want {
				t.Errorf("isSuccess = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewClientWithBaseURL(t *testing.T) {
	t.Parallel()

	c := NewClientWithBaseURL("http://127.0.0.1:1", "key")
	if c.baseURL != "http://127.0.0.1:1" || c.apiKey != "key" {
		t.Errorf("client = %+v, want the base URL and key it was built with", c)
	}
}
