package collector

func GetItems() map[string][]string {
	items := make(map[string][]string)
	items["stocks"] = getStocks()
	items["fxs"] = getFxs()
	items["commodities"] = getCommodities()
	items["cryptos"] = getCryptos()
	return items
}

func getFxs() []string {
	var fxs []string
	fxs = append(fxs, "EURUSD")
	fxs = append(fxs, "GBPUSD")
	fxs = append(fxs, "CNYUSD")
	return fxs
}

func getStocks() []string {
	var stocks []string
	stocks = append(stocks, "NVDA")
	stocks = append(stocks, "AMZN")
	stocks = append(stocks, "TSLA")
	return stocks
}

func getCommodities() []string {
	var commodities []string
	commodities = append(commodities, "CLUSD")
	commodities = append(commodities, "KCUSX")
	commodities = append(commodities, "GCUSD")
	return commodities
}

func getCryptos() []string {
	var cryptos []string
	cryptos = append(cryptos, "BTCUSD")
	cryptos = append(cryptos, "ETHUSD")
	cryptos = append(cryptos, "BNBUSD")
	return cryptos
}
