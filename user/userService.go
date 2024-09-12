package user

import (
	"errors"
	"fmt"
	"log"
	"math"
	"slices"
	"time"

	"github.com/udvarid/don-trade-golang/collector"
	"github.com/udvarid/don-trade-golang/model"
	"github.com/udvarid/don-trade-golang/orderService"
	"github.com/udvarid/don-trade-golang/repository/candleRepository"
	"github.com/udvarid/don-trade-golang/repository/userRepository"
)

func ChangeName(id string, name string) {
	user, err := userRepository.FindUser(id)
	if err == nil {
		user.Name = name
		userRepository.UpdateUser(user)
	}
}

func DeleteUser(id string) {
	orders := orderService.GetOrdersByUserId(id)
	for _, order := range orders {
		orderService.DeleteOrder(order.ID, id)
	}
	userRepository.DeleteUser(id)
}

func createUserIfNotExists(id string) {
	_, err := userRepository.FindUser(id)
	if err != nil {
		var newUser model.User
		newUser.ID = id
		newUser.Name = GetRandomUniqueName(getUserNames())
		newUser.Config = getDefaultUserConfig()
		newUser.Transactions = getInitTransactions()
		newUser.Assets = getInitAssets()
		userRepository.AddUser(newUser)
	}
}

func GetUser(id string) model.UserStatistic {
	user, _ := userRepository.FindUser(id)

	var userStatistic model.UserStatistic
	userStatistic.ID = user.ID
	userStatistic.Name = user.Name
	if len(user.Transactions) > 10 {
		userStatistic.Transactions = user.Transactions[len(user.Transactions)-10:]
	} else {
		userStatistic.Transactions = user.Transactions
	}
	candleSummary := candleRepository.GetAllCandleSummaries()[0]
	userStatistic.Assets = getAssetsWithValue(user.Assets, candleSummary)
	return userStatistic
}

func GetPriceChanges() []model.PriceChanges {
	pureToday, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	priceHistory := createPriceHistory(pureToday)
	items := collector.GetItemsFromItemMap(collector.GetItems())
	itemsBaseDate := make(map[time.Time]map[string]float64)
	for _, historyElement := range priceHistory {
		itemsBaseDate[historyElement.Date] = historyElement.Items
	}
	candleSummaries := candleRepository.GetAllCandleSummaries()[0]
	var changes []model.PriceChanges
	for _, item := range items {
		candleLastDate := candleSummaries.Summary[item.Name].LastDate
		itemsForToday := itemsBaseDate[candleLastDate]
		itemsForYesterday := itemsBaseDate[candleLastDate.AddDate(0, 0, -1)]
		itemsForWeekAgo := itemsBaseDate[candleLastDate.AddDate(0, 0, -7)]
		var change model.PriceChanges
		change.Item = item.Description
		dailyChange := math.Round((itemsForToday[item.Name]/itemsForYesterday[item.Name]-1)*1000) / 10
		dailyChangeStr := fmt.Sprintf("%.1f", dailyChange)
		if dailyChange > 0.0 {
			dailyChangeStr = "+" + dailyChangeStr
		}
		weeklyChange := math.Round((itemsForToday[item.Name]/itemsForWeekAgo[item.Name]-1)*1000) / 10
		weeklyChangeStr := fmt.Sprintf("%.1f", weeklyChange)
		if weeklyChange > 0.0 {
			weeklyChangeStr = "+" + weeklyChangeStr
		}
		change.Change = dailyChangeStr + "%/" + weeklyChangeStr + "%"
		changes = append(changes, change)
	}
	return changes
}

func createPriceHistory(pureToday time.Time) []model.HistoryElement {
	candles := candleRepository.GetAllCandles()
	items := collector.GetItemsFromItemMap(collector.GetItems())
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

func GetUserHistory(id string, days int) []model.HistoryElement {
	createUserIfNotExists(id)
	var result []model.HistoryElement
	pureToday, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	priceHistory := createPriceHistory(pureToday)
	user, _ := userRepository.FindUser(id)
	transactions := user.Transactions
	periodStart := pureToday.AddDate(0, 0, -days)

	assetsOfUser := make(map[string]bool)
	assetsVolumeHistory := make(map[time.Time]map[string]float64)
	firstTransactionDate := transactions[0].Date
	for _, transaction := range transactions {
		assetsOfUser[transaction.Asset] = true
		today, exists := assetsVolumeHistory[transaction.Date]
		if !exists {
			today = make(map[string]float64)
		}
		today[transaction.Asset] += transaction.Volume
		assetsVolumeHistory[transaction.Date] = today
	}
	var assetsVolumeHistoryDaily []model.HistoryElement
	for !firstTransactionDate.After(pureToday) {
		var historyElement model.HistoryElement
		historyElement.Date = firstTransactionDate
		historyElement.Items = assetsVolumeHistory[firstTransactionDate]
		assetsVolumeHistoryDaily = append(assetsVolumeHistoryDaily, historyElement)
		firstTransactionDate = firstTransactionDate.AddDate(0, 0, 1)
	}
	for i := 1; i < len(assetsVolumeHistoryDaily); i++ {
		element_current := assetsVolumeHistoryDaily[i].Items
		element_previous := assetsVolumeHistoryDaily[i-1].Items
		element_refreshed := make(map[string]float64)
		for item, volume := range element_previous {
			element_refreshed[item] = volume
		}
		for item, volume := range element_current {
			element_refreshed[item] += volume
		}
		assetsVolumeHistoryDaily[i].Items = element_refreshed
	}

	for !periodStart.After(pureToday) {
		priceElement := getElementByDate(priceHistory, periodStart)
		volumenElement := getElementByDate(assetsVolumeHistoryDaily, periodStart)
		if len(priceElement.Items) > 0 {
			priceElement.Items["USD"] = 1.0
			var historyElement model.HistoryElement
			historyElement.Date = periodStart
			historyElement.Items = make(map[string]float64)
			for item, price := range priceElement.Items {
				muliplier := 0.0
				if len(volumenElement.Items) > 0 {
					muliplier = volumenElement.Items[item]
				}
				historyElement.Items[item] = price * muliplier
			}

			result = append(result, historyElement)
		}
		periodStart = periodStart.AddDate(0, 0, 1)
	}

	return result
}

func getElementByDate(history []model.HistoryElement, date time.Time) model.HistoryElement {
	for _, element := range history {
		if element.Date == date {
			return element
		}
	}
	return model.HistoryElement{}
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

func getAssetsWithValue(assets map[string][]model.VolumeWithPrice, candleSummary model.CandleSummary) []model.AssetWithValue {
	var result []model.AssetWithValue
	totalValue := 0.0
	for asset, volumes := range assets {
		if asset != "USD" {
			price := candleSummary.Summary[asset].LastPrice
			if len(volumes) > 0 {
				volume := 0.0
				for _, volumeWithPrice := range volumes {
					volume += volumeWithPrice.Volume
				}
				value := price * volume
				totalValue += value
				result = append(result, model.AssetWithValue{
					Item:   asset,
					Volume: volume,
					Price:  price,
					Value:  value,
				})
			}
		}
	}
	usd := assets["USD"]
	if len(usd) > 0 {
		result = append(result, model.AssetWithValue{
			Item:   "USD",
			Volume: usd[0].Volume,
			Price:  1.0,
			Value:  usd[0].Volume,
		})
		totalValue += usd[0].Volume
	}
	result = append(result, model.AssetWithValue{
		Item:  "Total",
		Value: totalValue,
	})
	return result
}

func getInitAssets() map[string][]model.VolumeWithPrice {
	assets := make(map[string][]model.VolumeWithPrice)
	var usdList []model.VolumeWithPrice
	var initUsd model.VolumeWithPrice
	initUsd.Volume = 1000000.0
	initUsd.Price = 1.0
	usdList = append(usdList, initUsd)
	assets["USD"] = usdList
	return assets
}

func getInitTransactions() []model.Transaction {
	pureToday, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	transactions := []model.Transaction{
		{
			Asset:  "USD",
			Date:   pureToday,
			Volume: 1000000.0,
		},
	}
	return transactions
}

func getDefaultUserConfig() model.UserConfig {
	return model.UserConfig{
		NotifyDaily:         true,
		NotifyAtTransaction: true,
	}
}

func getUserNames() []string {
	var names []string
	users := userRepository.GetAllUsers()
	for _, user := range users {
		names = append(names, user.Name)
	}
	return names
}
