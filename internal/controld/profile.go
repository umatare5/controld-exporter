// Package controld provides a client for interacting with the ControlD API.
package controld

const (
	ProfilesEndpoint = "/profiles" // Endpoint for retrieving profiles
)

// ProfilesResponse represents the response from the /profiles API.
type ProfilesResponse struct {
	Success bool `json:"success"` // Indicates if the API request was successful
	Body    struct {
		Profiles []struct {
			PK      string `json:"PK"`      // Primary key of the profile
			Updated int    `json:"updated"` // Last updated timestamp
			Name    string `json:"name"`    // Name of the profile
			Profile struct {
				Flt struct {
					Count int `json:"count"` // Count of filters
				} `json:"flt"`
				Cflt struct {
					Count int `json:"count"` // Count of custom filters
				} `json:"cflt"`
				Ipflt struct {
					Count int `json:"count"` // Count of IP filters
				} `json:"ipflt"`
				Rule struct {
					Count int `json:"count"` // Count of rules
				} `json:"rule"`
				Svc struct {
					Count int `json:"count"` // Count of services
				} `json:"svc"`
				Grp struct {
					Count int `json:"count"` // Count of groups
				} `json:"grp"`
				Opt struct {
					Count int `json:"count"` // Count of options
					Data  []struct {
						PK    string  `json:"PK"`    // Primary key of the option
						Value float64 `json:"value"` // Value of the option
					} `json:"data"`
				} `json:"opt"`
			} `json:"profile"`
		} `json:"profiles"`
	} `json:"body"`
}

// GetProfiles retrieves profiles without additional headers.
func (t *Client) GetProfiles() (*ProfilesResponse, error) {
	return t.sendProfilesRequest(nil)
}

// GetSubOrgProfiles retrieves profiles with additional headers for a specific organization.
func (t *Client) GetSubOrgProfiles(orgID string) (*ProfilesResponse, error) {
	return t.sendProfilesRequest(t.buildOrgIDHeader(orgID))
}

// sendProfilesRequest sends a request to fetch profiles.
func (t *Client) sendProfilesRequest(headers map[string]string) (*ProfilesResponse, error) {
	var data ProfilesResponse
	err := t.sendAPIRequest(ProfilesEndpoint, headers, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
