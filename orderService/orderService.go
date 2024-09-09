package orderService

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/udvarid/don-trade-golang/collector"
	"github.com/udvarid/don-trade-golang/model"
	"github.com/udvarid/don-trade-golang/repository/orderRepository"
)

func DeleteOrder(orderId int, userId string) {
	order := orderRepository.GetOrder(orderId)
	if order.UserID == userId {
		orderRepository.DeleteOrder(orderId)
	}
}

func GetOrdersByUserId(userId string) []model.Order {
	var orders []model.Order
	allOrders := orderRepository.GetAllOrders()
	for _, order := range allOrders {
		if order.UserID == userId {
			orders = append(orders, order)
		}
	}
	return orders
}

func ValidateAndAddOrder(orderInString model.OrderInString, userId string) {
	var order model.Order
	items := collector.GetItemsFromItemMap(collector.GetItems())
	var itemNames []string
	for item := range items {
		itemNames = append(itemNames, item)
	}
	if isOrderValid(orderInString, itemNames) {
		order.UserID = userId
		order.Item = orderInString.Item
		order.Direction = orderInString.Direction
		order.Type = orderInString.Type
		order.AllIn = orderInString.AllIn
		if orderInString.LimitPrice != "" {
			limitPrice, err := strconv.ParseFloat(orderInString.LimitPrice, 64)
			if err != nil {
				fmt.Println("Error converting LimitPrice:", err)
				return
			}
			order.LimitPrice = limitPrice
		}

		if orderInString.NumberOfItems != "" {
			numberOfItems, err := strconv.ParseFloat(orderInString.NumberOfItems, 64)
			if err != nil {
				fmt.Println("Error converting NumberOfItems:", err)
				return
			}
			order.NumberOfItems = numberOfItems
		}

		if orderInString.Usd != "" {
			usd, err := strconv.ParseFloat(orderInString.Usd, 64)
			if err != nil {
				fmt.Println("Error converting Usd:", err)
				return
			}
			order.Usd = usd
		}

		validDays, err := strconv.Atoi(orderInString.ValidDays)
		if err != nil {
			fmt.Println("Error converting ValidDays:", err)
			return
		}
		order.ValidDays = validDays

		AddOrder(order)
	}
}

func isOrderValid(orderInString model.OrderInString, itemNames []string) bool {
	if orderInString.Type == "LIMIT" && orderInString.LimitPrice == "" {
		return false
	}
	if orderInString.NumberOfItems == "" && orderInString.Usd == "" && !orderInString.AllIn {
		return false
	}
	if !slices.Contains(itemNames, orderInString.Item) {
		return false
	}
	if orderInString.Direction == "SELL" && orderInString.NumberOfItems == "" && !orderInString.AllIn {
		return false
	}
	if orderInString.Direction == "SELL" && orderInString.AllIn && orderInString.NumberOfItems != "" {
		return false
	}
	return true
}

func AddOrder(order model.Order) {
	orderRepository.AddOrder(order)
}