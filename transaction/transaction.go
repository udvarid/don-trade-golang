package transaction

import (
	"github.com/udvarid/don-trade-golang/model"
	"github.com/udvarid/don-trade-golang/repository/userRepository"
)

func HandleTransaction(transactionPositive model.Transaction, transactionNegative model.Transaction, userId string) {
	user, _ := userRepository.FindUser(userId)
	user.Transactions = append(user.Transactions, transactionPositive)
	user.Transactions = append(user.Transactions, transactionNegative)
	user.Assets[transactionPositive.Asset] += transactionPositive.Volume
	user.Assets[transactionNegative.Asset] += transactionNegative.Volume
	userRepository.UpdateUser(user)
}
