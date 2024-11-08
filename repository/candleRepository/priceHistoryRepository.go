package candleRepository

import (
	"encoding/json"

	"github.com/udvarid/don-trade-golang/model"
	"github.com/udvarid/don-trade-golang/repository/repoUtil"
	bolt "go.etcd.io/bbolt"
)

func GetAllPriceHistory() []model.GroupOfHistoryElement {
	db := repoUtil.OpenDb()
	var result []model.GroupOfHistoryElement
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("PriceHistory"))

		b.ForEach(func(k, v []byte) error {
			var priceHistory model.GroupOfHistoryElement
			json.Unmarshal([]byte(v), &priceHistory)
			result = append(result, priceHistory)
			return nil
		})
		return nil
	})
	defer db.Close()

	return result
}

func UpdatePriceHistory(priceHistory *model.GroupOfHistoryElement) {
	db := repoUtil.OpenDb()
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("PriceHistory"))
		buf, err := json.Marshal(priceHistory)
		if err != nil {
			return err
		}
		return b.Put(repoUtil.Itob(priceHistory.ID), buf)
	})
	defer db.Close()
}

func AddPriceHistory(priceHistory *model.GroupOfHistoryElement) {
	db := repoUtil.OpenDb()
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("PriceHistory"))
		id, _ := b.NextSequence()
		priceHistory.ID = int(id)
		buf, err := json.Marshal(priceHistory)
		if err != nil {
			return err
		}
		return b.Put(repoUtil.Itob(priceHistory.ID), buf)
	})

	defer db.Close()
}

func DeletePriceHistory(priceHistoryId int) {
	db := repoUtil.OpenDb()

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("PriceHistory"))
		err := b.Delete(repoUtil.Itob(priceHistoryId))
		return err
	})
	defer db.Close()
}
