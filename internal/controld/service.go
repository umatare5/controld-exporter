// Package controld provides a client for interacting with the ControlD API.
package controld

const (
	ServiceCategoriesEndpoint = "/services/categories" // Endpoint for retrieving service categories
)

// ServiceCategoriesResponse represents the response from the /services/categories API.
type ServiceCategoriesResponse struct {
	Success bool `json:"success"` // Indicates if the API request was successful
	Body    struct {
		Categories []struct {
			PK          string `json:"PK"`          // Primary key of the category
			Name        string `json:"name"`        // Name of the category
			Description string `json:"description"` // Description of the category
			Count       int    `json:"count"`       // Number of items in the category
		} `json:"categories"`
	} `json:"body"`
}

// GetServiceCategories retrieves service categories without additional headers.
func (t *Client) GetServiceCategories() (*ServiceCategoriesResponse, error) {
	return t.sendServiceCategoriesRequest(nil)
}

// GetSubOrgServiceCategories retrieves service categories with additional headers for a specific organization.
func (t *Client) GetSubOrgServiceCategories(orgID string) (*ServiceCategoriesResponse, error) {
	return t.sendServiceCategoriesRequest(t.buildOrgIDHeader(orgID))
}

// sendServiceCategoriesRequest sends a request to fetch service categories.
func (t *Client) sendServiceCategoriesRequest(headers map[string]string) (*ServiceCategoriesResponse, error) {
	var data ServiceCategoriesResponse
	err := t.sendAPIRequest(ServiceCategoriesEndpoint, headers, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
