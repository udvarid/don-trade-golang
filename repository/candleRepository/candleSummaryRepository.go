package candleRepository

import (
	"encoding/json"

	"github.com/udvarid/don-trade-golang/model"
	"github.com/udvarid/don-trade-golang/repository/repoUtil"
	bolt "go.etcd.io/bbolt"
)

func GetAllCandleSummaries() []model.CandleSummary {
	db := repoUtil.OpenDb()
	var result []model.CandleSummary
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("CandleSummary"))

		b.ForEach(func(k, v []byte) error {
			var candleSummary model.CandleSummary
			json.Unmarshal([]byte(v), &candleSummary)
			result = append(result, candleSummary)
			return nil
		})
		return nil
	})
	defer db.Close()

	return result
}

func UpdateCandleSummary(candleSummary *model.CandleSummary) {
	db := repoUtil.OpenDb()
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("CandleSummary"))
		buf, err := json.Marshal(candleSummary)
		if err != nil {
			return err
		}
		return b.Put(repoUtil.Itob(candleSummary.ID), buf)
	})
	defer db.Close()
}

func AddCandleSummary(candleSummary *model.CandleSummary) {
	db := repoUtil.OpenDb()
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("CandleSummary"))
		id, _ := b.NextSequence()
		candleSummary.ID = int(id)
		buf, err := json.Marshal(candleSummary)
		if err != nil {
			return err
		}
		return b.Put(repoUtil.Itob(candleSummary.ID), buf)
	})

	defer db.Close()
}

func DeleteCandleSummary(candleSummaryId int) {
	db := repoUtil.OpenDb()

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("CandleSummary"))
		err := b.Delete(repoUtil.Itob(candleSummaryId))
		return err
	})
	defer db.Close()
}
