package collector

import (
	"fmt"
	"log"
	"sync"
	"time"

	dataCollector "github.com/spacecodewor/fmpcloud-go"
	"github.com/spacecodewor/fmpcloud-go/objects"
	"github.com/udvarid/don-trade-golang/model"
)

// https://intelligence.financialmodelingprep.com/developer/docs/dashboard

func CollectData(config *model.Configuration) {
	APIClient, err := dataCollector.NewAPIClient(dataCollector.Config{APIKey: config.Price_collector_api_key})
	if err != nil {
		log.Println("Error init api client: " + err.Error())
	}

	dayParameter := 10

	items := GetItems()
	stockItems := items["stocks"]
	stockChannel := make(chan StockResult, len(stockItems))
	var wgStock sync.WaitGroup
	for _, stockItem := range stockItems {
		wgStock.Add(1)
		go collectStocks(&wgStock, stockChannel, stockItem, dayParameter, *APIClient)
	}

	go func() {
		wgStock.Wait()
		close(stockChannel)
	}()

	for stockElem := range stockChannel {
		for _, stockOnDay := range stockElem.result {
			fmt.Println(stockElem.name, stockOnDay.Date, stockOnDay.Close)
		}
	}
}

func collectStocks(wg *sync.WaitGroup, channel chan StockResult, stockItem string, days int, client dataCollector.APIClient) {
	defer wg.Done()

	param := objects.RequestStockCandleList{}
	param.Symbol = stockItem
	param.Period = "1day"
	til := time.Now()
	from := til.Add(time.Duration(-days) * 24 * time.Hour)
	param.To = &til
	param.From = &from
	candles, err := client.Stock.Candles(param)
	if err != nil {
		log.Println("Error get quote: " + err.Error())
	}
	var result StockResult
	result.name = stockItem
	result.result = candles

	channel <- result
}

type StockResult struct {
	name   string
	result []objects.StockCandle
}
