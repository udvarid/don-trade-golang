package user

import (
	"errors"
	"fmt"
	"log"
	"math"
	"slices"
	"sort"
	"time"

	"github.com/udvarid/don-trade-golang/collector"
	"github.com/udvarid/don-trade-golang/communicator"
	"github.com/udvarid/don-trade-golang/model"
	"github.com/udvarid/don-trade-golang/orderService"
	"github.com/udvarid/don-trade-golang/repository/candleRepository"
	"github.com/udvarid/don-trade-golang/repository/orderRepository"
	"github.com/udvarid/don-trade-golang/repository/userRepository"
	userstatistic "github.com/udvarid/don-trade-golang/userStatistic"
)

func ChangeName(id string, name string) {
	user, err := userRepository.FindUser(id)
	if err == nil {
		user.Name = name
		userRepository.UpdateUser(user)
	}
}

func ChangeNotify(id string, transaction bool, daily bool) {
	user, err := userRepository.FindUser(id)
	if err == nil {
		user.Config.NotifyAtTransaction = transaction
		user.Config.NotifyDaily = daily
		userRepository.UpdateUser(user)
	}
}

func DeleteUser(id string) {
	for _, order := range orderService.GetOrdersByUserId(id) {
		orderRepository.DeleteOrder(order.ID)
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
		newUser.Debts = make(map[string][]model.VolumeWithPrice)
		userRepository.AddUser(newUser)
	}
}

func GetTraders() []model.UserSummary {
	result := make([]model.UserSummary, 0)
	users := userRepository.GetAllUsers()
	for _, user := range users {
		var userSummary model.UserSummary
		userAssets := userstatistic.GetUserStatistic(user.ID, false).Assets
		userSummary.UserID = user.ID
		userSummary.UserName = user.Name
		userSummary.TraderSince = int(time.Since(user.Transactions[0].Date).Hours() / 24)
		total := userAssets[len(userAssets)-1].Value
		userSummary.Profit = total/1000000 - 1
		invested := 0.0
		for _, asset := range userAssets {
			if asset.Item != "Total" && asset.Item != "USD" {
				invested += asset.Value
			}
		}
		userSummary.Invested = invested / total
		userSummary.CreditLimit = userstatistic.CalculateCreditLimit(userAssets)
		result = append(result, userSummary)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Profit > result[j].Profit
	})

	return result
}

func SendDailyStatus() {
	users := userRepository.GetAllUsers()
	for _, user := range users {
		if user.Config.NotifyDaily {
			communicator.SendMessageAboutStatus(userstatistic.GetUserStatistic(user.ID, false))
		}
	}
}

func GetUser(id string) model.User {
	user, _ := userRepository.FindUser(id)
	return user
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

func GetUserHistory(id string, days int) *model.GroupOfHistoryElement {
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
				multiplier := 0.0
				if len(volumenElement.Items) > 0 {
					multiplier = volumenElement.Items[item]
				}
				historyElement.Items[item] = price * multiplier
			}

			result = append(result, historyElement)
		}
		periodStart = periodStart.AddDate(0, 0, 1)
	}

	var groupResult model.GroupOfHistoryElement
	groupResult.Group = result

	return &groupResult
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
