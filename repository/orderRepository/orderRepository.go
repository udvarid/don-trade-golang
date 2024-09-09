package orderRepository

import (
	"encoding/json"

	"github.com/udvarid/don-trade-golang/model"
	"github.com/udvarid/don-trade-golang/repository/repoUtil"
	bolt "go.etcd.io/bbolt"
)

func GetAllOrders() []model.Order {
	db := repoUtil.OpenDb()
	var result []model.Order
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Order"))

		b.ForEach(func(k, v []byte) error {
			var order model.Order
			json.Unmarshal([]byte(v), &order)
			result = append(result, order)
			return nil
		})
		return nil
	})
	defer db.Close()

	return result
}

func GetOrder(orderId int) model.Order {
	db := repoUtil.OpenDb()

	var result model.Order
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Order"))
		v := b.Get(repoUtil.Itob(orderId))
		if v != nil {
			var order model.Order
			json.Unmarshal([]byte(v), &order)
			result = order
		}
		return nil
	})
	defer db.Close()
	return result
}

func DeleteOrder(orderId int) {
	db := repoUtil.OpenDb()

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Order"))
		err := b.Delete(repoUtil.Itob(orderId))
		return err
	})
	defer db.Close()
}

func AddOrder(order model.Order) model.Order {
	db := repoUtil.OpenDb()
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Order"))
		id, _ := b.NextSequence()
		order.ID = int(id)
		buf, err := json.Marshal(order)
		if err != nil {
			return err
		}
		return b.Put(repoUtil.Itob(order.ID), buf)
	})

	defer db.Close()
	return order
}

func UpdateOrder(order model.Order) {
	db := repoUtil.OpenDb()
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Order"))
		buf, err := json.Marshal(order)
		if err != nil {
			return err
		}
		return b.Put(repoUtil.Itob(order.ID), buf)
	})
	defer db.Close()
}
