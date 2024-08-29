package model

import "time"

type Configuration struct {
	Mail_psw                string `json:"mail_psw"`
	Mail_from               string `json:"mail_from"`
	Mail_local_psw          string `json:"mail_local_psw"`
	Mail_local_from         string `json:"mail_local_from"`
	Admin_user              string `json:"admin_user"`
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

type UserConfig struct {
	NotifyDaily         bool `json:"notify_daily"`
	NotifyAtTransaction bool `json:"notify_at_transaction"`
}

type Transaction struct {
	Asset  string    `json:"asset"`
	Date   time.Time `json:"date"`
	Volume float64   `json:"volume"`
}

type User struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Config       UserConfig         `json:"config"`
	Assets       map[string]float64 `json:"assets"`
	Transactions []Transaction      `json:"transactions"`
}

type AssetWithValue struct {
	Item   string  `json:"item"`
	Volume float64 `json:"volume"`
	Price  float64 `json:"price"`
	Value  float64 `json:"value"`
}

type AssetWithValueInString struct {
	Item   string `json:"item"`
	Volume string `json:"volume"`
	Price  string `json:"price"`
	Value  string `json:"value"`
}

type UserStatistic struct {
	ID           string           `json:"id"`
	Name         string           `json:"name"`
	Assets       []AssetWithValue `json:"assets"`
	Transactions []Transaction    `json:"transactions"`
}

type HistoryElement struct {
	Date  time.Time          `json:"date"`
	Items map[string]float64 `json:"items"`
}
