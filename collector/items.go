package collector

import "github.com/udvarid/don-trade-golang/model"

func GetItems() map[string][]model.Item {
	items := make(map[string][]model.Item)
	items["stocks"] = getStocks()
	items["fxs"] = getFxs()
	items["commodities"] = getCommodities()
	items["cryptos"] = getCryptos()
	return items
}

func getFxs() []model.Item {
	var fxs []model.Item
	fxs = append(fxs, createItem("EURUSD", "EUR-USD fx quote"))
	fxs = append(fxs, createItem("GBPUSD", "GBP-USD fx quote"))
	fxs = append(fxs, createItem("CHFUSD", "CHF-USD fx quote"))
	fxs = append(fxs, createItem("HUFUSD", "HUF-USD fx quote"))
	return fxs
}

func getStocks() []model.Item {
	var stocks []model.Item
	stocks = append(stocks, createItem("NVDA", "Nvidia stock price"))
	stocks = append(stocks, createItem("AMZN", "Amazon stock price"))
	stocks = append(stocks, createItem("TSLA", "Tesla stock price"))
	stocks = append(stocks, createItem("OTP.BD", "OTP stock price"))
	return stocks
}

func getCommodities() []model.Item {
	var commodities []model.Item
	commodities = append(commodities, createItem("CLUSD", "Crude Oil-USD price"))
	commodities = append(commodities, createItem("KCUSX", "Coffee-USD price"))
	commodities = append(commodities, createItem("GCUSD", "Gold-USD price"))
	commodities = append(commodities, createItem("SBUSX", "Sugar-USD price"))
	return commodities
}

func getCryptos() []model.Item {
	var cryptos []model.Item
	cryptos = append(cryptos, createItem("BTCUSD", "Bitcoin-USD price"))
	cryptos = append(cryptos, createItem("ETHUSD", "Ethereum-USD price"))
	cryptos = append(cryptos, createItem("BNBUSD", "Binance-USD price"))
	cryptos = append(cryptos, createItem("SOLUSD", "Solana-USD price"))
	return cryptos
}

func createItem(name string, description string) model.Item {
	var item model.Item
	item.Name = name
	item.Description = description
	return item
}
