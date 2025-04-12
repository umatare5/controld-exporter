// Package controld provides a client for interacting with the ControlD API.
package controld

const (
	DevicesEndpoint = "/devices" // Endpoint for retrieving device information
)

// DevicesResponse represents the response from the /devices API.
type DevicesResponse struct {
	Success bool `json:"success"` // Indicates if the API request was successful
	Body    struct {
		Devices []struct {
			PK          string `json:"PK"`           // Primary key of the device
			Timestamp   int    `json:"ts"`           // Timestamp of the device data
			Name        string `json:"name"`         // Name of the device
			Org         string `json:"org"`          // Organization associated with the device
			Stats       int    `json:"stats"`        // Statistics related to the device
			DeviceID    string `json:"device_id"`    // Unique identifier for the device
			Status      int    `json:"status"`       // Current status of the device
			ClientCount int    `json:"client_count"` // Number of clients connected to the device
			LearnIP     int    `json:"learn_ip"`     // Indicates if the device learns IP addresses
			Ctrld       struct {
				Status    int    `json:"status"`     // Status of the ControlD service
				LastFetch int    `json:"last_fetch"` // Last fetch timestamp for ControlD data
				Version   string `json:"version"`    // Version of the ControlD service
			} `json:"ctrld"`
			Resolvers struct {
				UID string   `json:"uid"` // Unique identifier for the resolver
				DOH string   `json:"doh"` // DNS-over-HTTPS endpoint
				DOT string   `json:"dot"` // DNS-over-TLS endpoint
				V6  []string `json:"v6"`  // IPv6 addresses for the resolver
			} `json:"resolvers"`
			Icon    string `json:"icon"` // Icon representing the device
			Profile struct {
				PK      string `json:"PK"`      // Primary key of the profile
				Updated int    `json:"updated"` // Last updated timestamp for the profile
				Name    string `json:"name"`    // Name of the profile
			} `json:"profile"`
			IPCount      int `json:"ip_total"`      // Number of IP addresses associated with the device
			LastActivity int `json:"last_activity"` // Last activity timestamp of the device
			Clients      map[string]struct {
				Timestamp int      `json:"ts"`   // Timestamp of the client data
				Host      string   `json:"host"` // Hostname of the client
				MAC       string   `json:"mac"`  // MAC address of the client
				IP        string   `json:"ip"`   // IP address of the client
				OS        []string `json:"os"`   // Operating systems of the client
			} `json:"clients"`
		} `json:"devices"`
	} `json:"body"`
}

// GetDevices retrieves devices without additional headers.
func (t *Client) GetDevices() (*DevicesResponse, error) {
	return t.sendDevicesRequest(nil)
}

// GetSubOrgDevices retrieves devices with additional headers for a specific organization.
func (t *Client) GetSubOrgDevices(orgID string) (*DevicesResponse, error) {
	return t.sendDevicesRequest(t.buildOrgIDHeader(orgID))
}

// sendDevicesRequest sends a request to fetch devices.
func (t *Client) sendDevicesRequest(headers map[string]string) (*DevicesResponse, error) {
	var data DevicesResponse
	err := t.sendAPIRequest(DevicesEndpoint, headers, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
