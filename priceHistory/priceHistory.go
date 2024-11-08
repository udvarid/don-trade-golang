package priceHistory

import (
	"errors"
	"log"
	"slices"
	"time"

	"github.com/udvarid/don-trade-golang/model"
	"github.com/udvarid/don-trade-golang/repository/candleRepository"
)

func GetPriceHistory(pureToday time.Time, force bool, items map[string]model.Item) []model.HistoryElement {
	persistedPriceHistories := candleRepository.GetAllPriceHistory()
	if len(persistedPriceHistories) > 0 && !force {
		return persistedPriceHistories[0].Group
	}
	candles := candleRepository.GetAllCandles()
	var itemNames []string
	for item := range items {
		itemNames = append(itemNames, item)
	}
	firstDate, err := getFirstDate(candles, itemNames, pureToday)
	if err != nil {
		log.Default().Println("Error getting first date")
	}

	return getPriceHistory(candles, itemNames, firstDate, pureToday)
}

func getPriceHistory(candles []model.Candle, itemNames []string, firstDate time.Time, pureToday time.Time) []model.HistoryElement {
	var result []model.HistoryElement
	var firstElement model.HistoryElement
	firstElement.Date = firstDate
	firstElement.Items = make(map[string]float64)
	for _, candle := range candles {
		if candle.Date == firstDate {
			firstElement.Items[candle.Item] = candle.Close
		}
		if len(firstElement.Items) == len(itemNames) {
			break
		}
	}
	result = append(result, firstElement)
	for {
		firstDate = firstDate.AddDate(0, 0, 1)
		if firstDate.After(pureToday) {
			break
		}
		var nextElement model.HistoryElement
		nextElement.Date = firstDate
		nextElement.Items = make(map[string]float64)
		for _, item := range itemNames {
			foundItem := false
			for _, candle := range candles {
				if candle.Date == firstDate && candle.Item == item {
					foundItem = true
					nextElement.Items[item] = candle.Close
					break
				}
			}
			if !foundItem {
				priceYesterday := result[len(result)-1].Items[item]
				nextElement.Items[item] = priceYesterday
			}
		}
		result = append(result, nextElement)
	}
	return result
}

func getFirstDate(candles []model.Candle, itemNames []string, pureToday time.Time) (time.Time, error) {
	firstDate := time.Now()
	itemsByDate := make(map[time.Time][]string)
	for _, candle := range candles {
		if candle.Date.Before(firstDate) {
			firstDate = candle.Date
		}
		dayElement, exists := itemsByDate[candle.Date]
		if !exists {
			itemsByDate[candle.Date] = []string{candle.Item}
		} else {
			itemsByDate[candle.Date] = append(dayElement, candle.Item)
		}
	}
	foundDateWhenAllItemsExist := false
	for !foundDateWhenAllItemsExist {
		dayElement, exists := itemsByDate[firstDate]
		contains := true
		if exists {
			for _, itemName := range itemNames {
				if !slices.Contains(dayElement, itemName) {
					contains = false
					break
				}
			}
			if contains {
				foundDateWhenAllItemsExist = true
			}
		}
		if !foundDateWhenAllItemsExist {
			firstDate = firstDate.AddDate(0, 0, 1)
			if firstDate.After(pureToday) {
				break
			}
		}
	}
	if !foundDateWhenAllItemsExist {
		return time.Now(), errors.New("no first date found")
	}
	return firstDate, nil
}
