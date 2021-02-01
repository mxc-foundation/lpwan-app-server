package user

// ShopifyAdminAPI defines shopify admin api configuration
type ShopifyAdminAPI struct {
	Hostname   string `mapstructure:"hostname"`
	APIKey     string `mapstructure:"api_key"`
	Secret     string `mapstructure:"secret"`
	APIVersion string `mapstructure:"api_version"`
	StoreName  string `mapstructure:"store_name"`
}

// ShopifyCustomer includes a part of attributes of customer
type ShopifyCustomer struct {
	ID          int64  `json:"id"`
	Email       string `json:"email"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	OrdersCount int    `json:"orders_count"`
	State       string `json:"state"`
	LastOrderID int64  `json:"last_order_id"`
}

// ShopifyCustomerList maps response of api
// https://apikey:secret@{hostname}/admin/api/2021-01/customers/search.json\?query\=email:{email}
type ShopifyCustomerList struct {
	Customers []ShopifyCustomer `json:"customers"`
}
