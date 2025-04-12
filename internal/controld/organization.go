// Package controld provides a client for interacting with the ControlD API.
package controld

const (
	OrganizationEndpoint     = "/organizations/organization"      // Endpoint for retrieving main organization data
	SubOrganizationsEndpoint = "/organizations/sub_organizations" // Endpoint for retrieving sub organization data
)

// OrganizationResponse represents the response structure for the /organizations/organization endpoint.
type OrganizationResponse struct {
	Success bool `json:"success"` // Indicates if the API request was successful
	Body    struct {
		Organization struct {
			Learned              int      `json:"learned"`                // Number of learned items
			SiemEnabled          int      `json:"siem_enabled"`           // Indicates if SIEM is enabled
			ContactFirstName     string   `json:"contact_first_name"`     // First name of the contact person
			ContactLastName      string   `json:"contact_last_name"`      // Last name of the contact person
			Name                 string   `json:"name"`                   // Name of the organization
			Date                 string   `json:"date"`                   // Creation date of the organization
			MaxProfiles          int      `json:"max_profiles"`           // Maximum number of profiles allowed
			MaxSubOrgs           int      `json:"max_sub_orgs"`           // Maximum number of sub-organizations allowed
			PriceUsers           int      `json:"price_users"`            // Price per user
			MaxLegacyResolvers   int      `json:"max_legacy_resolvers"`   // Maximum number of legacy resolvers
			Website              string   `json:"website"`                // Website of the organization
			OktaClientSecret     string   `json:"okta_client_secret"`     // Okta client secret
			TwofaReq             int      `json:"twofa_req"`              // Indicates if 2FA is required
			Type                 string   `json:"type"`                   // Type of the organization
			BillingMethod        int      `json:"billing_method"`         // Billing method
			Status               int      `json:"status"`                 // Status of the organization
			StatsEndpoint        string   `json:"stats_endpoint"`         // Endpoint for statistics
			ContactEmail         string   `json:"contact_email"`          // Contact email address
			OktaDomain           string   `json:"okta_domain"`            // Okta domain
			TrialEnd             string   `json:"trial_end"`              // Trial end date
			HubspotCompanyURL    string   `json:"hubspot_company_url"`    // HubSpot company URL
			OktaClientID         string   `json:"okta_client_id"`         // Okta client ID
			MaxUsers             int      `json:"max_users"`              // Maximum number of users allowed
			PK                   string   `json:"PK"`                     // Primary key of the organization
			StatusPrinted        string   `json:"status_printed"`         // Human-readable status
			BillingMethodPrinted string   `json:"billing_method_printed"` // Human-readable billing method
			SsoProvider          string   `json:"sso_provider"`           // SSO provider
			SsoEmailDomains      []string `json:"sso_email_domains"`      // Email domains for SSO
			Members              struct {
				Count int `json:"count"` // Number of members
			} `json:"members"`
			Profiles struct {
				Count int `json:"count"` // Number of profiles
				Max   int `json:"max"`   // Maximum number of profiles
			} `json:"profiles"`
			Users struct {
				Count int `json:"count"` // Number of users
				Price int `json:"price"` // Price per user
				Max   int `json:"max"`   // Maximum number of users
			} `json:"users"`
			Routers struct {
				Count int `json:"count"` // Number of routers
				Max   int `json:"max"`   // Maximum number of routers
				Price int `json:"price"` // Price per router
			} `json:"routers"`
			SubOrganizations struct {
				Count int `json:"count"` // Number of sub-organizations
				Max   int `json:"max"`   // Maximum number of sub-organizations
			} `json:"sub_organizations"`
		} `json:"organization"`
	} `json:"body"`
}

// SubOrganizationsResponse represents the response structure for the /organizations/sub_organizations endpoint.
type SubOrganizationsResponse struct {
	Success bool `json:"success"` // Indicates if the API request was successful
	Body    struct {
		SubOrganizations []struct {
			ParentProfile        string `json:"parent_profile"`         // Parent profile ID
			ContactName          string `json:"contact_name"`           // Name of the contact person
			StatsEndpoint        string `json:"stats_endpoint"`         // Endpoint for statistics
			SiemEnabled          int    `json:"siem_enabled"`           // Indicates if SIEM is enabled
			AllowOverrides       string `json:"allow_overrides"`        // Indicates if overrides are allowed
			MaxLegacyResolvers   int    `json:"max_legacy_resolvers"`   // Maximum number of legacy resolvers
			MaxProfiles          int    `json:"max_profiles"`           // Maximum number of profiles
			ParentOrg            string `json:"parent_org"`             // Parent organization ID
			TwofaReq             int    `json:"twofa_req"`              // Indicates if 2FA is required
			ContactEmail         string `json:"contact_email"`          // Contact email address
			Status               int    `json:"status"`                 // Status of the sub-organization
			Date                 string `json:"date"`                   // Creation date of the sub-organization
			Name                 string `json:"name"`                   // Name of the sub-organization
			MaxUsers             int    `json:"max_users"`              // Maximum number of users
			PK                   string `json:"PK"`                     // Primary key of the sub-organization
			StatusPrinted        string `json:"status_printed"`         // Human-readable status
			BillingMethodPrinted string `json:"billing_method_printed"` // Human-readable billing method
			SubOrganizations     struct {
				Count int `json:"count"` // Number of sub-organizations
				Max   int `json:"max"`   // Maximum number of sub-organizations
			} `json:"sub_organizations"`
			Members struct {
				Count int `json:"count"` // Number of members
			} `json:"members"`
			Profiles struct {
				Count int `json:"count"` // Number of profiles
				Max   int `json:"max"`   // Maximum number of profiles
			} `json:"profiles"`
			Users struct {
				Count int `json:"count"` // Number of users
				Price int `json:"price"` // Price per user
				Max   int `json:"max"`   // Maximum number of users
			} `json:"users"`
			Routers struct {
				Count int `json:"count"` // Number of routers
				Max   int `json:"max"`   // Maximum number of routers
				Price int `json:"price"` // Price per router
			} `json:"routers"`
		} `json:"sub_organizations"`
	} `json:"body"`
}

// GetMainOrganization fetches the list of main organization.
func (t *Client) GetMainOrganization() (*OrganizationResponse, error) {
	var data OrganizationResponse
	err := t.sendAPIRequest(OrganizationEndpoint, nil, &data)
	if err != nil {
		return nil, err
	}

	if err := t.handleAPIError(OrganizationEndpoint, data.Success); err != nil {
		return nil, err
	}

	return &data, nil
}

// GetSubOrganizations fetches the list of sub organization.
func (t *Client) GetSubOrganizations() (*SubOrganizationsResponse, error) {
	var data SubOrganizationsResponse
	err := t.sendAPIRequest(SubOrganizationsEndpoint, nil, &data)
	if err != nil {
		return nil, err
	}

	if err := t.handleAPIError(SubOrganizationsEndpoint, data.Success); err != nil {
		return nil, err
	}

	return &data, nil
}
