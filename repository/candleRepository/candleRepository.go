package candleRepository

import (
	"encoding/json"

	"github.com/udvarid/don-trade-golang/model"
	"github.com/udvarid/don-trade-golang/repository/repoUtil"
	bolt "go.etcd.io/bbolt"
)

func GetAllCandles() []model.Candle {
	db := repoUtil.OpenDb()
	var result []model.Candle
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Candle"))

		b.ForEach(func(k, v []byte) error {
			var candle model.Candle
			json.Unmarshal([]byte(v), &candle)
			result = append(result, candle)
			return nil
		})
		return nil
	})
	defer db.Close()

	return result
}

func DeleteCandle(candleId int) {
	db := repoUtil.OpenDb()

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Candle"))
		err := b.Delete(repoUtil.Itob(candleId))
		return err
	})
	defer db.Close()
}

func AddCandle(candle *model.Candle) {
	db := repoUtil.OpenDb()
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Candle"))
		id, _ := b.NextSequence()
		candle.ID = int(id)
		buf, err := json.Marshal(candle)
		if err != nil {
			return err
		}
		return b.Put(repoUtil.Itob(candle.ID), buf)
	})

	defer db.Close()
}
