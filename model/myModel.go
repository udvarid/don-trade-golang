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
	ID              int                        `json:"id"`
	Date            time.Time                  `json:"date"`
	Summary         map[string]CandleStatistic `json:"summary"`
	Persisted       []string                   `json:"persisted"`
	DailyStatusSent bool                       `json:"daily_status_sent"`
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

type TransactionWitString struct {
	Asset  string `json:"asset"`
	Date   string `json:"date"`
	Volume string `json:"volume"`
}

type User struct {
	ID           string                       `json:"id"`
	Name         string                       `json:"name"`
	Config       UserConfig                   `json:"config"`
	Assets       map[string][]VolumeWithPrice `json:"assets"`
	Debts        map[string][]VolumeWithPrice `json:"debts"`
	Transactions []Transaction                `json:"transactions"`
}

type VolumeWithPrice struct {
	Volume float64 `json:"volume"`
	Price  float64 `json:"price"`
}

type AssetWithValue struct {
	Item      string  `json:"item"`
	Volume    float64 `json:"volume"`
	Price     float64 `json:"price"`
	Value     float64 `json:"value"`
	BookValue float64 `json:"book_value"`
}

type AssetWithValueInString struct {
	Item      string `json:"item"`
	ItemPure  string `json:"item_pure"`
	Volume    string `json:"volume"`
	Price     string `json:"price"`
	Value     string `json:"value"`
	BookValue string `json:"book_value"`
	Profit    string `json:"profit"`
}

type UserStatistic struct {
	ID           string           `json:"id"`
	Name         string           `json:"name"`
	Assets       []AssetWithValue `json:"assets"`
	Transactions []Transaction    `json:"transactions"`
	CreditLimit  float64          `json:"credit_limit"`
}

type HistoryElement struct {
	Date  time.Time          `json:"date"`
	Items map[string]float64 `json:"items"`
}

type PriceChanges struct {
	Item   string `json:"item"`
	Change string `json:"daily"`
}

type Order struct {
	ID            int     `json:"id"`
	UserID        string  `json:"user_id"`
	Item          string  `json:"item"`
	Direction     string  `json:"direction"`
	Type          string  `json:"type"`
	LimitPrice    float64 `json:"limit_price"`
	NumberOfItems float64 `json:"number_of_items"`
	Usd           float64 `json:"usd"`
	AllIn         bool    `json:"all_in"`
	ValidDays     int     `json:"valid_days"`
}

type OrderModifyInString struct {
	OrderId    string `json:"order_id"`
	LimitPrice string `json:"limit_price"`
	ValidDays  string `json:"valid_days"`
}

type OrderInString struct {
	ID            int    `json:"id"`
	UserID        string `json:"user_id"`
	Item          string `json:"item"`
	Direction     string `json:"direction"`
	Type          string `json:"type"`
	LimitPrice    string `json:"limit_price"`
	NumberOfItems string `json:"number_of_items"`
	Usd           string `json:"usd"`
	AllIn         bool   `json:"all_in"`
	ValidDays     string `json:"valid_days"`
}

type UserSummary struct {
	UserID      string  `json:"user_id"`
	UserName    string  `json:"user_name"`
	Profit      float64 `json:"profit"`
	TraderSince int     `json:"trader_since"`
	Invested    float64 `json:"invested"`
	CreditLimit float64 `json:"credit_limit"`
}

type UserSummaryInString struct {
	UserID      string `json:"user_id"`
	UserName    string `json:"user_name"`
	Profit      string `json:"profit"`
	TraderSince int    `json:"trader_since"`
	Invested    string `json:"invested"`
	CreditLimit string `json:"credit_limit"`
}

type CompletedOrderToMail struct {
	Id      string `json:"id"`
	Type    string `json:"type"`
	Item    string `json:"item"`
	Volumen string `json:"volumen"`
	Price   string `json:"price"`
	Usd     string `json:"usd"`
}
