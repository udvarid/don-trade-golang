package orderService

import (
	"testing"

	"github.com/udvarid/don-trade-golang/model"
)

func TestIsOrderValid(t *testing.T) {
	testCases := []struct {
		name           string
		orderInString  model.OrderInString
		itemNames      []string
		expectedResult bool
	}{
		{
			name: "Valid LIMIT order",
			orderInString: model.OrderInString{
				Type:          "LIMIT",
				LimitPrice:    "50.00",
				NumberOfItems: "10",
				Usd:           "",
				Item:          "item1",
				Direction:     "BUY",
				AllIn:         false,
			},
			itemNames:      []string{"item1", "item2"},
			expectedResult: true,
		},
		{
			name: "Invalid LIMIT order with empty LimitPrice",
			orderInString: model.OrderInString{
				Type:          "LIMIT",
				LimitPrice:    "",
				NumberOfItems: "10",
				Usd:           "",
				Item:          "item1",
				Direction:     "BUY",
				AllIn:         false,
			},
			itemNames:      []string{"item1", "item2"},
			expectedResult: false,
		},
		{
			name: "Invalid order with empty NumberOfItems and Usd",
			orderInString: model.OrderInString{
				Type:          "MARKET",
				LimitPrice:    "",
				NumberOfItems: "",
				Usd:           "",
				Item:          "item1",
				Direction:     "BUY",
				AllIn:         false,
			},
			itemNames:      []string{"item1", "item2"},
			expectedResult: false,
		},
		{
			name: "Invalid order with non-existent item",
			orderInString: model.OrderInString{
				Type:          "MARKET",
				LimitPrice:    "",
				NumberOfItems: "10",
				Usd:           "",
				Item:          "nonExistentItem",
				Direction:     "BUY",
				AllIn:         false,
			},
			itemNames:      []string{"item1", "item2"},
			expectedResult: false,
		},
		{
			name: "Valid SELL order with AllIn",
			orderInString: model.OrderInString{
				Type:          "MARKET",
				LimitPrice:    "",
				NumberOfItems: "",
				Usd:           "",
				Item:          "item1",
				Direction:     "SELL",
				AllIn:         true,
			},
			itemNames:      []string{"item1", "item2"},
			expectedResult: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := isOrderValid(&tc.orderInString, tc.itemNames)
			if result != tc.expectedResult {
				t.Errorf("expected %t, got %t", tc.expectedResult, result)
			}
		})
	}
}
