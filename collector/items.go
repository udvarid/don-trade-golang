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
	fxs = append(fxs, createItem("EURUSD", "EUR"))
	fxs = append(fxs, createItem("GBPUSD", "GBP"))
	fxs = append(fxs, createItem("CHFUSD", "CHF"))
	fxs = append(fxs, createItem("HUFUSD", "HUF"))
	return fxs
}

func getStocks() []model.Item {
	var stocks []model.Item
	stocks = append(stocks, createItem("NVDA", "Nvidia"))
	stocks = append(stocks, createItem("AMZN", "Amazon"))
	stocks = append(stocks, createItem("TSLA", "Tesla"))
	stocks = append(stocks, createItem("SHEL", "Shell"))
	return stocks
}

func getCommodities() []model.Item {
	var commodities []model.Item
	commodities = append(commodities, createItem("CLUSD", "Crude Oil"))
	commodities = append(commodities, createItem("KCUSX", "Coffee"))
	commodities = append(commodities, createItem("GCUSD", "Gold"))
	commodities = append(commodities, createItem("SBUSX", "Sugar"))
	return commodities
}

func getCryptos() []model.Item {
	var cryptos []model.Item
	cryptos = append(cryptos, createItem("BTCUSD", "Bitcoin"))
	cryptos = append(cryptos, createItem("ETHUSD", "Ethereum"))
	cryptos = append(cryptos, createItem("BNBUSD", "Binance"))
	cryptos = append(cryptos, createItem("SOLUSD", "Solana"))
	return cryptos
}

func createItem(name string, description string) model.Item {
	var item model.Item
	item.Name = name
	item.Description = description
	return item
}
