package orderManager

import (
	"fmt"

	"github.com/udvarid/don-trade-golang/repository/candleRepository"
)

func ServeOrders(normal bool, user string) {
	candleSummary := candleRepository.GetAllCandleSummaries()[0]
	itemNames := candleSummary.Persisted
	if !normal {
		var tempItemNames []string
		for item := range candleSummary.Summary {
			tempItemNames = append(tempItemNames, item)
		}
		itemNames = tempItemNames
	}
	fmt.Println("Serving orders for", user, "with items:", itemNames)
}
