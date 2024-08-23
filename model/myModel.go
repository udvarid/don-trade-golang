package model

import "time"

type Configuration struct {
	Mail_psw                string `json:"mail_psw"`
	Mail_from               string `json:"mail_from"`
	Mail_local_psw          string `json:"mail_local_psw"`
	Mail_local_from         string `json:"mail_local_from"`
	Environment             string `json:"environment"`
	RemoteAddress           string `json:"remote_address"`
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

type BollingerBand struct {
	Item       string    `json:"item"`
	Date       time.Time `json:"date"`
	UpperBand  float64   `json:"upper_band"`
	LowerBand  float64   `json:"lower_band"`
	CenterBand float64   `json:"center_band"`
}

type Ma struct {
	Item     string    `json:"item"`
	Date     time.Time `json:"date"`
	MaLong   float64   `json:"ma_long"`
	MaMedium float64   `json:"ma_medium"`
	MaShort  float64   `json:"ma_short"`
}

type Macd struct {
	Item   string    `json:"item"`
	Date   time.Time `json:"date"`
	Macd   float64   `json:"macd"`
	Signal float64   `json:"signal"`
}

type TrendPoint struct {
	TrendPoint float64 `json:"trend_point"`
	TrendFlag  bool    `json:"trend_flag"`
}

type Rsi struct {
	Item string    `json:"item"`
	Date time.Time `json:"date"`
	RSI  float64   `json:"rsi"`
}

type Obv struct {
	Item string    `json:"item"`
	Date time.Time `json:"date"`
	Obv  float64   `json:"obv"`
}

type Adx struct {
	Item string    `json:"item"`
	Date time.Time `json:"date"`
	ADX  float64   `json:"adx"`
	PDI  float64   `json:"pdi"`
	MDI  float64   `json:"mdi"`
}

type Candle struct {
	ID     int       `json:"id"`
	Item   string    `json:"item"`
	Date   time.Time `json:"date"`
	Open   float64   `json:"open"`
	Close  float64   `json:"close"`
	High   float64   `json:"high"`
	Low    float64   `json:"low"`
	Volume float64   `json:"volume"`
}

type CandleSummary struct {
	ID        int                        `json:"id"`
	Date      time.Time                  `json:"date"`
	Summary   map[string]CandleStatistic `json:"summary"`
	Persisted []string                   `json:"persisted"`
}

type CandleStatistic struct {
	Number    int       `json:"number"`
	LastPrice float64   `json:"last_price"`
	LastDate  time.Time `json:"last_date"`
}

type Item struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SessionWithTime struct {
	ID        string `json:"id"`
	Session   string
	SessDate  time.Time
	IsChecked bool
}
