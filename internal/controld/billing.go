// Package controld provides a client for interacting with the ControlD API.
package controld

const (
	BillingPaymentsEndpoint      = "/billing/payments"      // Endpoint for retrieving billing payments
	BillingSubscriptionsEndpoint = "/billing/subscriptions" // Endpoint for retrieving billing subscriptions
)

// BillingPaymentsResponse represents the response structure for the /billing/payments endpoint.
type BillingPaymentsResponse struct {
	Success bool `json:"success"` // Indicates if the API request was successful
	Body    struct {
		Payments []struct {
			User           string `json:"user"`            // User associated with the payment
			Currency       string `json:"currency"`        // Currency of the payment
			SubID          string `json:"sub_id"`          // Subscription ID associated with the payment
			CurrencyAmount int    `json:"currency_amount"` // Amount in the specified currency
			Date           string `json:"date"`            // Payment date
			Product        struct {
				Type        string `json:"type"`         // Type of the product
				Priority    int    `json:"priority"`     // Priority of the product
				Name        string `json:"name"`         // Name of the product
				ProxyAccess int    `json:"proxy_access"` // Proxy access level
				PK          int    `json:"PK"`           // Primary key of the product
			} `json:"product"`
			Amount      int `json:"amount"`  // Payment amount
			Balance     int `json:"balance"` // Remaining balance after the payment
			Timestamp   int `json:"ts"`      // Timestamp of the payment
			Transaction struct {
				ID          string `json:"tx_id"`       // Transaction ID
				Status      int    `json:"tx_status"`   // Status of the transaction
				Refunded    int    `json:"tx_refunded"` // Refund status of the transaction
				Fingerprint string `json:"fingerprint"` // Fingerprint of the transaction
			} `json:"transaction"`
			PricePoint struct {
				ProductID int    `json:"product_id"` // Product ID associated with the price point
				Duration  int    `json:"duration"`   // Duration of the subscription in months
				JPYPrice  int    `json:"jpy_price"`  // Price in Japanese Yen
				EURPrice  int    `json:"eur_price"`  // Price in Euros
				GBPPrice  int    `json:"gbp_price"`  // Price in British Pounds
				AUDPrice  int    `json:"aud_price"`  // Price in Australian Dollars
				CADPrice  int    `json:"cad_price"`  // Price in Canadian Dollars
				CHFPrice  int    `json:"chf_price"`  // Price in Swiss Francs
				StripeID  string `json:"stripe_id"`  // Stripe ID for the price point
				Comment   string `json:"comment"`    // Additional comments about the price point
			} `json:"price_point"`
			Method string `json:"method"` // Payment method used
			PK     string `json:"PK"`     // Primary key of the payment
		} `json:"payments"`
	} `json:"body"`
}

// BillingSubscriptionsResponse represents the response structure for the /billing/subscriptions API.
type BillingSubscriptionsResponse struct {
	Success bool `json:"success"` // Indicates if the API request was successful
	Body    struct {
		Subscriptions []struct {
			Method  string `json:"method"` // Payment method for the subscription
			State   string `json:"state"`  // State of the subscription
			Product struct {
				Type        string `json:"type"`         // Type of the product
				Priority    int    `json:"priority"`     // Priority of the product
				Name        string `json:"name"`         // Name of the product
				ProxyAccess int    `json:"proxy_access"` // Proxy access level
				PK          int    `json:"PK"`           // Primary key of the product
			} `json:"product"`
			User           string `json:"user"`             // User associated with the subscription
			CurrencyAmount int    `json:"currency_amount"`  // Amount in the specified currency
			Currency       string `json:"currency"`         // Currency of the subscription
			NextBill       int    `json:"next_bill"`        // Timestamp of the next billing date
			PK             string `json:"PK"`               // Primary key of the subscription
			Status         int    `json:"status"`           // Status of the subscription
			NextRebillDate string `json:"next_rebill_date"` // Next rebill date
		} `json:"subscriptions"`
	} `json:"body"`
}

// GetBillingPayments fetches billing payments data from the API.
func (t *Client) GetBillingPayments() (*BillingPaymentsResponse, error) {
	var data BillingPaymentsResponse
	err := t.sendAPIRequest(BillingPaymentsEndpoint, nil, &data)
	if err != nil {
		return nil, err
	}

	if err := t.handleAPIError(BillingPaymentsEndpoint, data.Success); err != nil {
		return nil, err
	}

	return &data, nil
}

// GetBillingSubscriptions fetches billing subscriptions data from the API.
func (t *Client) GetBillingSubscriptions() (*BillingSubscriptionsResponse, error) {
	var data BillingSubscriptionsResponse
	err := t.sendAPIRequest(BillingSubscriptionsEndpoint, nil, &data)
	if err != nil {
		return nil, err
	}

	if err := t.handleAPIError(BillingSubscriptionsEndpoint, data.Success); err != nil {
		return nil, err
	}

	return &data, nil
}
