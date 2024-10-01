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

func MakeClearOrder(userId string, item string) {
	items := collector.GetItemsFromItemMap(collector.GetItems())
	_, exist := items[item]
	if !exist {
		return
	}
	orders := GetOrdersByUserId(userId)
	for _, order := range orders {
		if order.Item == item {
			return
		}
	}
	var clearOrder model.Order
	clearOrder.AllIn = true
	clearOrder.Direction = "SELL"
	clearOrder.Item = item
	clearOrder.Type = "MARKET"
	clearOrder.UserID = userId
	clearOrder.ValidDays = 1
	AddOrder(clearOrder)
}

func ModifyOrder(userId string, orderModify model.OrderModifyInString) {
	orders := GetOrdersByUserId(userId)
	for _, order := range orders {
		if strconv.Itoa(order.ID) == orderModify.OrderId {
			validChange := false
			if order.Type == "LIMIT" || order.Type == "STOP-LIMIT" {
				limitPrice, err := strconv.ParseFloat(orderModify.LimitPrice, 64)
				if err != nil {
					fmt.Println("Error converting LimitPrice:", err)
					return
				}
				if limitPrice > 0.0 && limitPrice != order.LimitPrice {
					order.LimitPrice = limitPrice
					validChange = true
				}
			}

			validDays, err := strconv.Atoi(orderModify.ValidDays)
			if err != nil {
				fmt.Println("Error converting ValidDays:", err)
				return
			}
			if validDays > 0 && validDays != order.ValidDays {
				order.ValidDays = validDays
				validChange = true
			}

			if validChange {
				orderRepository.UpdateOrder(order)
			}
		}
	}
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
	if (orderInString.Type == "LIMIT" || orderInString.Type == "STOP-LIMIT") && orderInString.LimitPrice == "" {
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
