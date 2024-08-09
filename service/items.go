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
	fxs = append(fxs, model.Item{"EURUSD", "EUR-USD fx quote"})
	fxs = append(fxs, model.Item{"GBPUSD", "GBP-USD fx quote"})
	fxs = append(fxs, model.Item{"CHFUSD", "CHF-USD fx quote"})
	return fxs
}

func getStocks() []model.Item {
	var stocks []model.Item
	stocks = append(stocks, model.Item{"NVDA", "Nvidia stock price"})
	stocks = append(stocks, model.Item{"AMZN", "Amazon stock price"})
	stocks = append(stocks, model.Item{"TSLA", "Tesla stock price"})
	return stocks
}

func getCommodities() []model.Item {
	var commodities []model.Item
	commodities = append(commodities, model.Item{"CLUSD", "Crude Oil-USD price"})
	commodities = append(commodities, model.Item{"KCUSX", "Coffee-USD price"})
	commodities = append(commodities, model.Item{"GCUSD", "Gold-USD price"})
	return commodities
}

func getCryptos() []model.Item {
	var cryptos []model.Item
	cryptos = append(cryptos, model.Item{"BTCUSD", "Bitcoin-USD price"})
	cryptos = append(cryptos, model.Item{"ETHUSD", "Ethereum-USD price"})
	cryptos = append(cryptos, model.Item{"BNBUSD", "Binance-USD price"})
	return cryptos
}
