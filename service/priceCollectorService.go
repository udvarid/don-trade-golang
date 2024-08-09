package collector

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/udvarid/don-trade-golang/model"
	"github.com/udvarid/don-trade-golang/repository/candleRepository"
)

/* Used provider
https://intelligence.financialmodelingprep.com/developer/docs/dashboard

Alternatives:
https://polygon.io/dashboard
https://www.alphavantage.co/
*/

func CollectData(config *model.Configuration) {

	items := GetItems()
	itemCounts := 0
	for _, v := range items {
		itemCounts += len(v)
	}
	channel := make(chan CandleResult, itemCounts)
	var wgStock sync.WaitGroup
	for _, items := range items {
		dayParameter := 15 // ennek értéke legyen 15 / 700, attól függően, hogy van e már ilyen az adatbázisban
		for _, item := range items {
			url := getUrl(dayParameter, item.Name, config.Price_collector_api_key)
			if url != "" {
				wgStock.Add(1)
				go collectCandles(&wgStock, channel, url, item.Name)
			}
		}
	}

	go func() {
		wgStock.Wait()
		close(channel)
	}()

	candlesPersisted := candleRepository.GetAllCandles()
	for result := range channel {
		persisted := 0
		for _, candleDto := range result.result {
			candle := mapCandleDtoToCandle(candleDto)
			if isCandleNew(&candlesPersisted, &candle) {
				candleRepository.AddCandle(candle)
				persisted++
			}
		}
		if persisted > 0 {
			fmt.Println("Persisted: ", result.name, persisted)
		}
	}

	candlesPersisted = candleRepository.GetAllCandles()
	itemCountMap := make(map[string]int)
	for _, candle := range candlesPersisted {
		itemCountMap[candle.Item]++
	}
	var candleSummary model.CandleSummary
	pureToday, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-01"))
	candleSummary.Date = pureToday
	candleSummary.Summary = itemCountMap
	// candleSummary lementése

	// 2 évnél régebbi candle-k törlése db-ből
}

func isCandleNew(candlesPersisted *[]model.Candle, candle *model.Candle) bool {
	for _, candlePersisted := range *candlesPersisted {
		if candlePersisted.Item == candle.Item && candlePersisted.Date == candle.Date {
			return false
		}
	}
	return true
}

func mapCandleDtoToCandle(candleDto model.CandleDto) model.Candle {
	var candle model.Candle
	candle.Item = candleDto.Item
	candle.High = candleDto.High
	candle.Low = candleDto.Low
	candle.Open = candleDto.Open
	candle.Close = candleDto.Close
	candle.Volume = candleDto.Volume
	candleDate, _ := time.Parse("2006-01-02", candleDto.Date)
	candle.Date = candleDate
	return candle
}

func getUrl(days int, item string, apiKey string) string {
	urlBase := "https://financialmodelingprep.com/api/v3/historical-chart/1day/"
	til := time.Now()
	from := til.Add(time.Duration(-days) * 24 * time.Hour)
	til_st := til.Format("2006-01-02")
	from_st := from.Format("2006-01-02")
	return urlBase + item + "?from=" + from_st + "&to=" + til_st + "&apikey=" + apiKey
}

func collectCandles(wg *sync.WaitGroup, channel chan CandleResult, url string, item string) {
	defer wg.Done()
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var candlesInit []model.CandleDto
	json.Unmarshal([]byte(responseData), &candlesInit)
	var candles []model.CandleDto
	for _, candleInit := range candlesInit {
		candle := candleInit
		candle.Item = item
		candle.Date = candleInit.Date[0:10]
		candles = append(candles, candle)
	}

	var result CandleResult
	result.name = item
	result.result = candles

	channel <- result
}

type CandleResult struct {
	name   string
	result []model.CandleDto
}
