package types

const(
	INVALID_TICKER = iota
	EUR_GBP = iota+1
	EUR_USD = iota+2
)

type ExchangeLatestResponse struct {
	Rates map[string]float64 `json:"rates"`
	Base string `json:"base"`
	Date string `json:"date"`
}

type ExchangeHistoryResponse struct {
	Rates map[string]map[string]float64 `json:"rates"`
}

type HTTPResponse struct {
	Ticker     string  `json:"ticker"`
	Rate       float32 `json:"rate"`
	WeekRate   float32 `json:"weeklyrate"`
	Prediction bool    `json:"buy"`
}
