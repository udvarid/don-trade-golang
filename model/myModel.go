package model

type Configuration struct {
	Mail_psw                string `json:"mail_psw"`
	Mail_from               string `json:"mail_from"`
	Mail_local_psw          string `json:"mail_local_psw"`
	Mail_local_from         string `json:"mail_local_from"`
	Environment             string `json:"environment"`
	Price_collector_api_key string `json:"price_collector_api_key"`
}
