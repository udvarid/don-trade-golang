package collector

import (
	"fmt"
	"log"
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

	param := objects.RequestStockCandleList{}
	param.Symbol = "AAPL"
	param.Period = "1day"
	til := time.Now()
	from := til.Add(-10 * 24 * time.Hour)
	param.To = &til
	param.From = &from

	candles, err := APIClient.Stock.Candles(param)
	if err != nil {
		log.Println("Error get quote: " + err.Error())
	}

	for _, candle := range candles {
		fmt.Println(candle.Date, candle.Open, candle.Close, candle.Volume)
	}

}
