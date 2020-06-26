package mining

type CMC struct {
	Status *Status `json:"status"`
	Data   *Data   `json:"data"`
}

type Status struct {
	Timestamp    string `json:"timestamp"`
	ErrorCode    int    `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Elapsed      int    `json:"elapsed"`
	CreditCount  int    `json:"credit_count"`
	Notice       string `json:"notice"`
}

type Data struct {
	Id          int    `json:"id"`
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	Amount      int    `json:"amount"`
	LastUpdated string `json:"last_updated"`
	Quote       struct {
		MXC struct {
			Price      float64 `json:"price"`
			LastUpdate string  `json:"last_update"`
		}
		USD struct {
			Price      float64 `json:"price"`
			LastUpdate string  `json:"last_update"`
		}
	}
}
