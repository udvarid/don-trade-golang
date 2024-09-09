package orderManager

import (
	"fmt"
	"slices"
	"time"

	"github.com/udvarid/don-trade-golang/model"
	"github.com/udvarid/don-trade-golang/repository/candleRepository"
	"github.com/udvarid/don-trade-golang/repository/orderRepository"
	"github.com/udvarid/don-trade-golang/repository/userRepository"
	"github.com/udvarid/don-trade-golang/transaction"
)

func ServeOrders(normal bool, user string) {
	candleSummary := candleRepository.GetAllCandleSummaries()[0]
	itemNames := candleSummary.Persisted
	if !normal {
		var tempItemNames []string
		for item := range candleSummary.Summary {
			tempItemNames = append(tempItemNames, item)
		}
		itemNames = tempItemNames
	}
	fmt.Println("Serving orders for", user, "with items:", itemNames)
	lastCandles := getlastCandles(candleSummary, itemNames)

	orderServed := true
	userAssetPairs := make(map[string]bool)
	pureToday, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	for orderServed {
		orderServed = false
		for _, order := range getRelevantOrders(normal, user, itemNames) {
			_, exists := userAssetPairs[order.UserID+"-"+order.Item]
			if exists {
				continue
			}
			user, _ := userRepository.FindUser(order.UserID)
			candle := lastCandles[order.Item]
			if order.Direction == "BUY" && order.Type == "MARKET" && user.Assets["USD"] >= 0.0001 {
				price := candle.Open
				initUsd := user.Assets["USD"]
				if order.Usd > 0.0001 && initUsd > order.Usd {
					initUsd = order.Usd
				}
				initVolume := initUsd / price
				if order.NumberOfItems > 0.0001 && initVolume > order.NumberOfItems {
					initVolume = order.NumberOfItems
				}
				var transactionPositive model.Transaction
				transactionPositive.Asset = order.Item
				transactionPositive.Date = pureToday
				transactionPositive.Volume = initVolume

				var transactionNegative model.Transaction
				transactionNegative.Asset = "USD"
				transactionNegative.Date = pureToday
				transactionNegative.Volume = initVolume * price

				transaction.HandleTransaction(transactionPositive, transactionNegative, user.ID)

				orderRepository.DeleteOrder(order.ID)
				orderServed = true
				userAssetPairs[order.UserID+"-"+order.Item] = true
			}
			if order.Direction == "BUY" && order.Type == "LIMIT" && user.Assets["USD"] >= 0.0001 {
				/*
					- buy limit
						- ha a candle open ára alacsonyabb, mint a limit, akkor open ár, ha magasabb, de a 'lowest' alacsonyabb, akkor limit áron, amúgy nem vesz
						- vagy a megadott darabszámot, vagy a megadott usd-ért (ha all-in, akkor minden pénez számít), ami kevesebb értékű
				*/
				// user transaction
				// user summary change
				orderRepository.DeleteOrder(order.ID)
				orderServed = true
				userAssetPairs[order.UserID+"-"+order.Item] = true
			}
			if order.Direction == "SELL" && order.Type == "MARKET" && user.Assets[order.Item] >= 0.0001 {
				/*
					- sell market
						- a candle open árán ad el
						- vagy a megadott darabszámot, vagy az összeset, ha all-in csekkolva van
				*/
				// user transaction
				// user summary change
				orderRepository.DeleteOrder(order.ID)
				orderServed = true
				userAssetPairs[order.UserID+"-"+order.Item] = true
			}
			if order.Direction == "SELL" && order.Type == "LIMIT" && user.Assets[order.Item] >= 0.0001 {
				/*
					- sell limit
						- ha a candle open ár magasabb, mint a limit, akkor open ár, ha alacsonyabb, de a highest magasabb, akkor a limit áron, amúgy nem adja el
						- vagy a megadott darabszámot, vagy az összeset, ha all-in csekkolva van
				*/
				// user transaction
				// user summary change
				orderRepository.DeleteOrder(order.ID)
				orderServed = true
				userAssetPairs[order.UserID+"-"+order.Item] = true
			}
		}

		for _, order := range getRelevantOrders(normal, user, itemNames) {
			order.ValidDays--
			if order.ValidDays <= 0 {
				orderRepository.DeleteOrder(order.ID)
			} else {
				orderRepository.UpdateOrder(order)
			}
		}
	}

}

func getRelevantOrders(normal bool, user string, itemNames []string) []model.Order {
	var ordersToServe []model.Order
	for _, order := range orderRepository.GetAllOrders() {
		if !normal && order.UserID != user || !slices.Contains(itemNames, order.Item) {
			continue
		}
		ordersToServe = append(ordersToServe, order)
	}
	return ordersToServe
}

func getlastCandles(candleSummary model.CandleSummary, itemNames []string) map[string]model.Candle {
	lastCandles := make(map[string]model.Candle)
	candles := candleRepository.GetAllCandles()
	for _, candle := range candles {
		if slices.Contains(itemNames, candle.Item) && candleSummary.Summary[candle.Item].LastDate == candle.Date {
			lastCandles[candle.Item] = candle
		}
	}

	return lastCandles

}
