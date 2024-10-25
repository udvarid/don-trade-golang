package userstatistic

import (
	"math"
	"sort"

	"github.com/udvarid/don-trade-golang/model"
	"github.com/udvarid/don-trade-golang/repository/candleRepository"
	"github.com/udvarid/don-trade-golang/repository/userRepository"
)

func GetUserStatistic(id string, onlyTransactions bool) *model.UserStatistic {
	user, _ := userRepository.FindUser(id)

	var userStatistic model.UserStatistic
	userStatistic.ID = user.ID
	userStatistic.Name = user.Name
	transactionLimit := 25
	if onlyTransactions {
		if len(user.Transactions) > transactionLimit {
			userStatistic.Transactions = user.Transactions[len(user.Transactions)-transactionLimit:]
		} else {
			userStatistic.Transactions = user.Transactions
		}
	} else {
		candleSummary := candleRepository.GetAllCandleSummaries()[0]
		userStatistic.Assets = getAssetsWithValue(user.Assets, user.Debts, candleSummary)
		userStatistic.CreditLimit = CalculateCreditLimit(userStatistic.Assets)
	}
	return &userStatistic
}

func getAssetsWithValue(
	assets map[string][]model.VolumeWithPrice,
	debts map[string][]model.VolumeWithPrice,
	candleSummary model.CandleSummary) []model.AssetWithValue {
	assetsAndDebts := joinAssetsAndDebts(assets, debts)
	var result []model.AssetWithValue
	totalValue := 0.0
	for asset, volumes := range assetsAndDebts {
		if asset != "USD" {
			price := candleSummary.Summary[asset].LastPrice
			if len(volumes) > 0 {
				volume := 0.0
				bookValue := 0.0
				for _, volumeWithPrice := range volumes {
					volume += volumeWithPrice.Volume
					bookValue += volumeWithPrice.Volume * volumeWithPrice.Price
				}
				value := price * volume
				totalValue += value
				if math.Abs(value) >= 0.1 {
					result = append(result, model.AssetWithValue{
						Item:      asset,
						Volume:    volume,
						Price:     price,
						Value:     value,
						BookValue: bookValue,
					})
				}
			}
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Item < result[j].Item
	})

	usd := assetsAndDebts["USD"]
	if len(usd) > 0 {
		result = append(result, model.AssetWithValue{
			Item:   "USD",
			Volume: usd[0].Volume,
			Price:  1.0,
			Value:  usd[0].Volume,
		})
		totalValue += usd[0].Volume
	}
	result = append(result, model.AssetWithValue{
		Item:      "Total",
		Value:     totalValue,
		BookValue: 1000000,
	})
	return result
}

func CalculateCreditLimit(assets []model.AssetWithValue) float64 {
	total := 0.0
	debt := 0.0
	for _, asset := range assets {
		if asset.Item == "Total" {
			total = asset.Value
		} else if asset.Item != "USD" && asset.Value < 0 {
			debt += asset.Value
		}
	}
	return math.Abs(debt / total)
}

func joinAssetsAndDebts(assets map[string][]model.VolumeWithPrice, debts map[string][]model.VolumeWithPrice) map[string][]model.VolumeWithPrice {
	result := make(map[string][]model.VolumeWithPrice)
	for asset, volumes := range assets {
		if len(volumes) > 0 {
			result[asset] = volumes
		}
	}
	for asset, volumes := range debts {
		if len(volumes) > 0 {
			result[asset] = volumes
		}
	}
	return result
}
