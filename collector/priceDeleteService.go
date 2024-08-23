package collector

import (
	"fmt"

	"github.com/udvarid/don-trade-golang/model"
	"github.com/udvarid/don-trade-golang/repository/candleRepository"
)

func DeletePriceDatabase(config *model.Configuration) {
	candles := candleRepository.GetAllCandles()
	candleSummaries := candleRepository.GetAllCandleSummaries()
	fmt.Printf("Deleting %d candles and requesting new candles\n", len(candles))
	for _, candle := range candles {
		candleRepository.DeleteCandle(candle.ID)
	}
	for _, candleSummary := range candleSummaries {
		candleRepository.DeleteCandleSummary(candleSummary.ID)
	}
	CollectData(config)
}
