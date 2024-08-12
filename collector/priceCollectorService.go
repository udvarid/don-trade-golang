package collector

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
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

	// If today there was already a data collection, then we quit
	summaries := candleRepository.GetAllCandleSummaries()
	pureToday, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	if len(summaries) > 0 && summaries[0].Date == pureToday {
		log.Println("No more data collection today")
		return
	}

	// starting data collection for each item in separated goroutins
	itemMap := GetItems()
	itemCounts := 0
	for _, v := range itemMap {
		itemCounts += len(v)
	}
	channel := make(chan CandleResult, itemCounts)
	var wgStock sync.WaitGroup
	for _, items := range itemMap {
		for _, item := range items {
			dayParameter := getDayParameter(summaries, item.Name, 15)
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

	// waiting the results and the new candles will be presisted
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

	// the old and unrelevant candles should be deleted
	timeTwoYearsBefore := pureToday.AddDate(-2, 0, 0)
	itemNames := getItemsFromItemMap(itemMap)
	for _, candlePersisted := range candlesPersisted {
		if shouldBeDeleted(&candlePersisted, itemNames, timeTwoYearsBefore) {
			candleRepository.DeleteCandle(candlePersisted.ID)
		}
	}

	// creating new candleSummary statistics
	candlesPersisted = candleRepository.GetAllCandles()
	itemCountMap := make(map[string]int)
	for _, candle := range candlesPersisted {
		itemCountMap[candle.Item]++
	}

	var candleSummary model.CandleSummary
	candleSummary.Date = pureToday
	candleSummary.Summary = itemCountMap
	if len(summaries) == 0 {
		candleRepository.AddCandleSummary(candleSummary)
	} else {
		candleSummaryToUpdate := summaries[0]
		candleSummaryToUpdate.Date = candleSummary.Date
		candleSummaryToUpdate.Summary = candleSummary.Summary
		candleRepository.UpdateCandleSummary(candleSummaryToUpdate)
	}
}

func getItemsFromItemMap(itemMap map[string][]model.Item) []string {
	itemNameSet := make(map[string]bool)
	for _, v := range itemMap {
		for _, item := range v {
			itemNameSet[item.Name] = true
		}
	}
	var result []string
	for k := range itemNameSet {
		result = append(result, k)
	}
	return result
}

func shouldBeDeleted(candlePersisted *model.Candle, itemNames []string, date time.Time) bool {
	exists := slices.Contains(itemNames, candlePersisted.Item)
	return !exists || candlePersisted.Date.Before(date)
}

func getDayParameter(summaries []model.CandleSummary, item string, defaultDays int) int {
	longDayParam := 700
	if len(summaries) == 0 {
		return longDayParam
	}
	summary := summaries[0].Summary
	numberOfItem, exists := summary[item]
	if !exists || numberOfItem < 300 {
		return longDayParam
	}
	return defaultDays
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
