// Package controld provides a client for interacting with the ControlD API.
package controld

const (
	NetworkEndpoint = "/network" // API endpoint for network information
)

// NetworkResponse represents the response structure for the /network endpoint.
type NetworkResponse struct {
	Body struct {
		Network []struct {
			IATACode    string `json:"iata_code"`    // IATA code of the network location
			CityName    string `json:"city_name"`    // City name of the network location
			CountryName string `json:"country_name"` // Country name of the network location
			Location    struct {
				Lat  float64 `json:"lat"`  // Latitude of the location
				Long float64 `json:"long"` // Longitude of the location
			} `json:"location"`
			Status struct {
				API int `json:"api"` // API status
				DNS int `json:"dns"` // DNS status
				PXY int `json:"pxy"` // Proxy status
			} `json:"status"`
		} `json:"network"`
		Time       int    `json:"time"`        // Response time in milliseconds
		CurrentPOP string `json:"current_pop"` // Current point of presence
	} `json:"body"`
	Success bool `json:"success"` // Indicates if the request was successful
}

// GetNetwork retrieves network information.
func (t *Client) GetNetwork() (*NetworkResponse, error) {
	var data NetworkResponse
	err := t.sendAPIRequest(NetworkEndpoint, nil, &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
