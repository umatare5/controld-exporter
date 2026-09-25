// Package controld provides a client for interacting with the ControlD API.
package controld

// Client represents a client for making requests to the ControlD API.
type Client struct {
	baseURL string // Base URL of the ControlD API
	apiKey  string // API key for authentication
}

// NewClient initializes and returns a new ControlD API client.
func NewClient(apiKey string) *Client {
	return &Client{
		baseURL: "https://api.controld.com",
		apiKey:  apiKey,
	}
}

// NewClientWithBaseURL points a client at baseURL instead of the vendor host,
// which is the seam a caller outside this package needs to serve its own doubles.
func NewClientWithBaseURL(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
	}
}
