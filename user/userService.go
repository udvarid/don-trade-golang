package user

import (
	"time"

	"github.com/udvarid/don-trade-golang/model"
	"github.com/udvarid/don-trade-golang/repository/candleRepository"
	"github.com/udvarid/don-trade-golang/repository/userRepository"
)

func GetUser(id string) model.UserStatistic {
	user, err := userRepository.FindUser(id)
	if err != nil {
		var newUser model.User
		newUser.ID = id
		newUser.Name = GetRandomUniqueName(getUserNames())
		newUser.Config = getDefaultUserConfig()
		newUser.Transactions = getInitTransactions()
		newUser.Assets = getInitAssets()
		userRepository.AddUser(newUser)
		user = newUser
	}

	var userStatistic model.UserStatistic
	userStatistic.ID = user.ID
	userStatistic.Name = user.Name
	if len(user.Transactions) > 10 {
		userStatistic.Transactions = user.Transactions[len(user.Transactions)-10:]
	} else {
		userStatistic.Transactions = user.Transactions
	}
	candleSummary := candleRepository.GetAllCandleSummaries()[0]
	userStatistic.Assets = getAssetsWithValue(user.Assets, candleSummary)
	return userStatistic
}

func getAssetsWithValue(assets map[string]float64, candleSummary model.CandleSummary) []model.AssetWithValue {
	var result []model.AssetWithValue
	for asset, volume := range assets {
		price := 1.0
		if asset != "USD" {
			price = candleSummary.Summary[asset].LastPrice
		}
		value := price * volume
		result = append(result, model.AssetWithValue{
			Item:   asset,
			Volume: volume,
			Price:  price,
			Value:  value,
		})
	}
	return result
}

func getInitAssets() map[string]float64 {
	assets := make(map[string]float64)
	assets["USD"] = 1000000.0
	return assets
}

func getInitTransactions() []model.Transaction {
	pureToday, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	transactions := []model.Transaction{
		{
			Asset:  "USD",
			Date:   pureToday,
			Volume: 1000000.0,
		},
	}
	return transactions
}

func getDefaultUserConfig() model.UserConfig {
	return model.UserConfig{
		NotifyDaily:         true,
		NotifyAtTransaction: true,
	}
}

func getUserNames() []string {
	var names []string
	users := userRepository.GetAllUsers()
	for _, user := range users {
		names = append(names, user.Name)
	}
	return names
}
