package orderService

import (
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

func AddOrder(order model.Order) {
	orderRepository.AddOrder(order)
}
