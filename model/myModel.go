package model

type Configuration struct {
	Mail_psw                string `json:"mail_psw"`
	Mail_from               string `json:"mail_from"`
	Mail_local_psw          string `json:"mail_local_psw"`
	Mail_local_from         string `json:"mail_local_from"`
	Environment             string `json:"environment"`
	Price_collector_api_key string `json:"price_collector_api_key"`
}

type CandleDto struct {
	Item   string  `json:"item"`
	Date   string  `json:"date"`
	Open   float64 `json:"open"`
	Close  float64 `json:"close"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Volume float64 `json:"volume"`
}

type Item struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
