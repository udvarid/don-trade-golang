package orderManager

import (
	"fmt"
	"math"
	"slices"
	"time"

	"github.com/udvarid/don-trade-golang/communicator"
	"github.com/udvarid/don-trade-golang/model"
	"github.com/udvarid/don-trade-golang/repository/candleRepository"
	"github.com/udvarid/don-trade-golang/repository/orderRepository"
	"github.com/udvarid/don-trade-golang/repository/userRepository"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
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
	candles := candleRepository.GetAllCandles()
	lastCandles := getlastCandles(candleSummary, itemNames, candles)

	orderServed := true
	userAssetPairs := make(map[string]bool)
	pureToday, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	var completedOrders []model.CompletedOrderToMail
	p := message.NewPrinter(language.Hungarian)
	for orderServed {
		orderServed = false
		orders := orderRepository.GetAllOrders()
		for _, order := range getRelevantOrders(normal, user, itemNames, orders) {
			_, exists := userAssetPairs[order.UserID+"-"+order.Item]
			if exists {
				continue
			}
			user, _ := userRepository.FindUser(order.UserID)
			candle := lastCandles[order.Item]
			if order.Direction == "BUY" && order.Type == "MARKET" && getVolumen(user.Assets["USD"]) >= 0.0001 {
				price := candle.Open
				initUsd := getVolumen(user.Assets["USD"])
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
				transactionNegative.Volume = initVolume * price * -1

				handleTransaction(transactionPositive, transactionNegative, user.ID)
				if user.Config.NotifyAtTransaction {
					var completedOrder model.CompletedOrderToMail
					completedOrder.Id = user.ID
					completedOrder.Item = order.Item
					completedOrder.Type = order.Direction
					completedOrder.Price = fmt.Sprintf("%.2f", price)
					completedOrder.Volumen = p.Sprintf("%d", int(initVolume))
					completedOrder.Usd = p.Sprintf("%.d", int(initVolume*price))
					completedOrders = append(completedOrders, completedOrder)
				}

				orderRepository.DeleteOrder(order.ID)
				orderServed = true
				userAssetPairs[order.UserID+"-"+order.Item] = true
			}
			if order.Direction == "BUY" && order.Type == "LIMIT" && getVolumen(user.Assets["USD"]) >= 0.0001 && candle.Low <= order.LimitPrice {
				price := candle.Open
				if candle.Open > order.LimitPrice {
					price = order.LimitPrice
				}
				initUsd := getVolumen(user.Assets["USD"])
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
				transactionNegative.Volume = initVolume * price * -1

				handleTransaction(transactionPositive, transactionNegative, user.ID)
				if user.Config.NotifyAtTransaction {
					var completedOrder model.CompletedOrderToMail
					completedOrder.Id = user.ID
					completedOrder.Item = order.Item
					completedOrder.Type = order.Direction
					completedOrder.Price = fmt.Sprintf("%.2f", price)
					completedOrder.Volumen = p.Sprintf("%d", int(initVolume))
					completedOrder.Usd = p.Sprintf("%.d", int(initVolume*price))
					completedOrders = append(completedOrders, completedOrder)
				}

				orderRepository.DeleteOrder(order.ID)
				orderServed = true
				userAssetPairs[order.UserID+"-"+order.Item] = true
			}
			if order.Direction == "BUY" && order.Type == "STOP-LIMIT" && getVolumen(user.Assets["USD"]) >= 0.0001 && candle.High >= order.LimitPrice {
				price := order.LimitPrice
				if candle.Open > order.LimitPrice {
					price = candle.Open
				}
				initUsd := getVolumen(user.Assets["USD"])
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
				transactionNegative.Volume = initVolume * price * -1

				handleTransaction(transactionPositive, transactionNegative, user.ID)
				if user.Config.NotifyAtTransaction {
					var completedOrder model.CompletedOrderToMail
					completedOrder.Id = user.ID
					completedOrder.Item = order.Item
					completedOrder.Type = order.Direction
					completedOrder.Price = fmt.Sprintf("%.2f", price)
					completedOrder.Volumen = p.Sprintf("%d", int(initVolume))
					completedOrder.Usd = p.Sprintf("%.d", int(initVolume*price))
					completedOrders = append(completedOrders, completedOrder)
				}

				orderRepository.DeleteOrder(order.ID)
				orderServed = true
				userAssetPairs[order.UserID+"-"+order.Item] = true
			}
			if order.Direction == "SELL" && order.Type == "MARKET" && getVolumen(user.Assets[order.Item]) >= 0.0001 {
				price := candle.Open
				initVolume := getVolumen(user.Assets[order.Item])
				if order.NumberOfItems > 0.0001 && initVolume > order.NumberOfItems {
					initVolume = order.NumberOfItems
				}

				var transactionPositive model.Transaction
				transactionPositive.Asset = "USD"
				transactionPositive.Date = pureToday
				transactionPositive.Volume = price * initVolume

				var transactionNegative model.Transaction
				transactionNegative.Asset = order.Item
				transactionNegative.Date = pureToday
				transactionNegative.Volume = initVolume * -1

				handleTransaction(transactionPositive, transactionNegative, user.ID)
				if user.Config.NotifyAtTransaction {
					var completedOrder model.CompletedOrderToMail
					completedOrder.Id = user.ID
					completedOrder.Item = order.Item
					completedOrder.Type = order.Direction
					completedOrder.Price = fmt.Sprintf("%.2f", price)
					completedOrder.Volumen = p.Sprintf("%d", int(initVolume))
					completedOrder.Usd = p.Sprintf("%.d", int(initVolume*price))
					completedOrders = append(completedOrders, completedOrder)
				}

				orderRepository.DeleteOrder(order.ID)
				orderServed = true
				userAssetPairs[order.UserID+"-"+order.Item] = true
			}
			if order.Direction == "SELL" && order.Type == "LIMIT" && getVolumen(user.Assets[order.Item]) >= 0.0001 && candle.High >= order.LimitPrice {

				price := candle.Open
				if candle.Open < order.LimitPrice {
					price = order.LimitPrice
				}
				initVolume := getVolumen(user.Assets[order.Item])
				if order.NumberOfItems > 0.0001 && initVolume > order.NumberOfItems {
					initVolume = order.NumberOfItems
				}
				var transactionPositive model.Transaction
				transactionPositive.Asset = "USD"
				transactionPositive.Date = pureToday
				transactionPositive.Volume = price * initVolume

				var transactionNegative model.Transaction
				transactionNegative.Asset = order.Item
				transactionNegative.Date = pureToday
				transactionNegative.Volume = initVolume * -1

				handleTransaction(transactionPositive, transactionNegative, user.ID)
				if user.Config.NotifyAtTransaction {
					var completedOrder model.CompletedOrderToMail
					completedOrder.Id = user.ID
					completedOrder.Item = order.Item
					completedOrder.Type = order.Direction
					completedOrder.Price = fmt.Sprintf("%.2f", price)
					completedOrder.Volumen = p.Sprintf("%d", int(initVolume))
					completedOrder.Usd = p.Sprintf("%.d", int(initVolume*price))
					completedOrders = append(completedOrders, completedOrder)
				}

				orderRepository.DeleteOrder(order.ID)
				orderServed = true
				userAssetPairs[order.UserID+"-"+order.Item] = true
			}
			if order.Direction == "SELL" && order.Type == "STOP-LIMIT" && getVolumen(user.Assets[order.Item]) >= 0.0001 && candle.Low <= order.LimitPrice {

				price := order.LimitPrice
				if candle.Open < order.LimitPrice {
					price = candle.Open
				}
				initVolume := getVolumen(user.Assets[order.Item])
				if order.NumberOfItems > 0.0001 && initVolume > order.NumberOfItems {
					initVolume = order.NumberOfItems
				}
				var transactionPositive model.Transaction
				transactionPositive.Asset = "USD"
				transactionPositive.Date = pureToday
				transactionPositive.Volume = price * initVolume

				var transactionNegative model.Transaction
				transactionNegative.Asset = order.Item
				transactionNegative.Date = pureToday
				transactionNegative.Volume = initVolume * -1

				handleTransaction(transactionPositive, transactionNegative, user.ID)
				if user.Config.NotifyAtTransaction {
					var completedOrder model.CompletedOrderToMail
					completedOrder.Id = user.ID
					completedOrder.Item = order.Item
					completedOrder.Type = order.Direction
					completedOrder.Price = fmt.Sprintf("%.2f", price)
					completedOrder.Volumen = p.Sprintf("%d", int(initVolume))
					completedOrder.Usd = p.Sprintf("%.d", int(initVolume*price))
					completedOrders = append(completedOrders, completedOrder)
				}

				orderRepository.DeleteOrder(order.ID)
				orderServed = true
				userAssetPairs[order.UserID+"-"+order.Item] = true
			}
		}
	}

	orders := orderRepository.GetAllOrders()
	for _, order := range getRelevantOrders(normal, user, itemNames, orders) {
		order.ValidDays--
		if order.ValidDays <= 0 {
			orderRepository.DeleteOrder(order.ID)
		} else {
			orderRepository.UpdateOrder(order)
		}
	}
	if len(completedOrders) > 0 {
		completedOrdersByUser := make(map[string][]model.CompletedOrderToMail)
		for _, completedOrder := range completedOrders {
			completedOrdersByUser[completedOrder.Id] = append(completedOrdersByUser[completedOrder.Id], completedOrder)
		}
		for userId, orders := range completedOrdersByUser {
			communicator.SendMessageAboutOrders(userId, orders)
		}
	}

}

func getVolumen(assets []model.VolumeWithPrice) float64 {
	var total float64
	for _, asset := range assets {
		total += asset.Volume
	}
	return total
}

func getRelevantOrders(normal bool, user string, itemNames []string, orders []model.Order) []model.Order {
	var ordersToServe []model.Order
	for _, order := range orders {
		if !normal && order.UserID != user || !slices.Contains(itemNames, order.Item) {
			continue
		}
		ordersToServe = append(ordersToServe, order)
	}
	return ordersToServe
}

func getlastCandles(candleSummary model.CandleSummary, itemNames []string, candles []model.Candle) map[string]model.Candle {
	lastCandles := make(map[string]model.Candle)
	for _, candle := range candles {
		if slices.Contains(itemNames, candle.Item) && candleSummary.Summary[candle.Item].LastDate == candle.Date {
			lastCandles[candle.Item] = candle
		}
	}

	return lastCandles
}

func handleTransaction(transactionPositive model.Transaction, transactionNegative model.Transaction, userId string) {
	user, _ := userRepository.FindUser(userId)
	user.Transactions = append(user.Transactions, transactionPositive)
	user.Transactions = append(user.Transactions, transactionNegative)
	if transactionPositive.Asset == "USD" {
		user.Assets["USD"][0].Volume += transactionPositive.Volume
		volumeSold := math.Abs(transactionNegative.Volume)
		oldPackages := user.Assets[transactionNegative.Asset]
		var newPackages []model.VolumeWithPrice
		for _, pack := range oldPackages {
			if volumeSold >= pack.Volume {
				volumeSold -= pack.Volume
			} else {
				newPackage := model.VolumeWithPrice{Volume: pack.Volume - volumeSold, Price: pack.Price}
				newPackages = append(newPackages, newPackage)
				volumeSold = 0
			}
		}
		user.Assets[transactionNegative.Asset] = newPackages
	} else {
		price := math.Abs(transactionNegative.Volume / transactionPositive.Volume)
		user.Assets["USD"][0].Volume += transactionNegative.Volume
		newPackage := model.VolumeWithPrice{Volume: transactionPositive.Volume, Price: price}
		user.Assets[transactionPositive.Asset] = append(user.Assets[transactionPositive.Asset], newPackage)
	}
	userRepository.UpdateUser(user)
}
