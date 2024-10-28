package collector

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
	"sort"
	"sync"
	"time"

	chart "github.com/udvarid/don-trade-golang/chartBuilder"
	"github.com/udvarid/don-trade-golang/model"
	"github.com/udvarid/don-trade-golang/orderManager"
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
	itemMap := GetItems()
	itemNamesWithItem := GetItemsFromItemMap(itemMap)
	candlesPersisted := candleRepository.GetAllCandles()
	if len(summaries) > 0 && summaries[0].Date == pureToday {
		log.Println("No more data collection today")
		var groupOfCandles model.GroupOfCandles
		groupOfCandles.Group = candlesPersisted
		orderCharts(&groupOfCandles, itemNamesWithItem)
		return
	}

	// starting data collection for each item in separated goroutins
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
	var persistedItems []string
	for result := range channel {
		persisted := 0
		for _, candleDto := range result.result {
			candle := mapCandleDtoToCandle(&candleDto)
			if isCandleNew(&candlesPersisted, &candle) {
				candleRepository.AddCandle(&candle)
				persisted++
			}
		}
		if persisted > 0 {
			persistedItems = append(persistedItems, result.name)
			fmt.Println("Persisted: ", result.name, persisted)
		}
	}

	// the old and unrelevant candles should be deleted
	timeSomeYearsBefore := pureToday.AddDate(-3, 0, 0)
	var itemNames []string
	for itemName := range itemNamesWithItem {
		itemNames = append(itemNames, itemName)
	}
	for _, candlePersisted := range candlesPersisted {
		if shouldBeDeleted(&candlePersisted, itemNames, timeSomeYearsBefore) {
			candleRepository.DeleteCandle(candlePersisted.ID)
		}
	}

	// creating new candleSummary statistics
	candlesPersisted = candleRepository.GetAllCandles()
	itemCountMap := make(map[string]model.CandleStatistic)
	for _, candle := range candlesPersisted {
		cs, exists := itemCountMap[candle.Item]
		if !exists {
			var candleStatistic model.CandleStatistic
			candleStatistic.Number = 1
			candleStatistic.LastPrice = candle.Close
			candleStatistic.LastDate = candle.Date
			itemCountMap[candle.Item] = candleStatistic
		} else {
			cs.Number++
			if cs.LastDate.Before(candle.Date) {
				cs.LastPrice = candle.Close
				cs.LastDate = candle.Date
			}
			itemCountMap[candle.Item] = cs
		}
	}

	var candleSummary model.CandleSummary
	candleSummary.Date = pureToday
	candleSummary.Summary = itemCountMap
	candleSummary.Persisted = persistedItems
	candleSummary.DailyStatusSent = false
	if len(summaries) == 0 {
		candleRepository.AddCandleSummary(&candleSummary)
	} else {
		candleSummaryToUpdate := summaries[0]
		candleSummaryToUpdate.Date = candleSummary.Date
		candleSummaryToUpdate.Summary = candleSummary.Summary
		candleSummaryToUpdate.Persisted = persistedItems
		candleSummaryToUpdate.DailyStatusSent = false
		candleRepository.UpdateCandleSummary(&candleSummaryToUpdate)
	}

	var groupOfCandles model.GroupOfCandles
	groupOfCandles.Group = candlesPersisted

	orderCharts(&groupOfCandles, itemNamesWithItem)
	if len(persistedItems) > 0 {
		orderManager.ServeOrders(true, "all users")
	}
}

func orderCharts(candles *model.GroupOfCandles, itemNames map[string]model.Item) {
	for itemName, item := range itemNames {
		var candlesToChart []model.Candle
		for _, candle := range candles.Group {
			if candle.Item == itemName {
				candlesToChart = append(candlesToChart, candle)
			}
		}
		sort.Slice(candlesToChart, func(i, j int) bool {
			return candlesToChart[i].Date.Before(candlesToChart[j].Date)
		})
		n := len(candlesToChart)
		if n > 100 {
			n = 100
		}
		var groupOfCandlesForSimpleChart model.GroupOfCandles
		var groupOfCandlesForDetailedChart model.GroupOfCandles
		groupOfCandlesForSimpleChart.Group = candlesToChart[len(candlesToChart)-n:]
		groupOfCandlesForDetailedChart.Group = candlesToChart
		chart.BuildSimpleCandleChart(&groupOfCandlesForSimpleChart, item.Description)
		chart.BuildDetailedChart(&groupOfCandlesForDetailedChart)
	}
}

func GetItemsFromItemMap(itemMap map[string][]model.Item) map[string]model.Item {
	itemNameSet := make(map[string]model.Item)
	for _, v := range itemMap {
		for _, item := range v {
			itemNameSet[item.Name] = item
		}
	}
	return itemNameSet
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
	candleStatistic, exists := summary[item]
	if !exists || candleStatistic.Number < 300 {
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

func mapCandleDtoToCandle(candleDto *model.CandleDto) model.Candle {
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
		log.Print("Couldn't get data")
		log.Print(err)
		retryCount := 1
		for err != nil && retryCount <= 10 {
			time.Sleep(5000 * time.Millisecond)
			log.Print("Trying to collect data again - ", retryCount)
			response, err = http.Get(url)
			retryCount++
		}
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
