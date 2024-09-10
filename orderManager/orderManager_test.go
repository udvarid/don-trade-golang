package orderManager

import (
	"testing"
	"time"

	"github.com/udvarid/don-trade-golang/model"
)

func TestGetRelevantOrders(t *testing.T) {
	testCases := []struct {
		name           string
		normal         bool
		user           string
		itemNames      []string
		orders         []model.Order
		expectedOrders []model.Order
	}{
		{
			name:      "Normal mode, all orders",
			normal:    true,
			user:      "user1",
			itemNames: []string{"item1", "item2", "item3"},
			orders: []model.Order{
				{UserID: "user1", Item: "item1"},
				{UserID: "user2", Item: "item2"},
				{UserID: "user1", Item: "item3"},
				{UserID: "user3", Item: "item1"},
			},
			expectedOrders: []model.Order{
				{UserID: "user1", Item: "item1"},
				{UserID: "user2", Item: "item2"},
				{UserID: "user1", Item: "item3"},
				{UserID: "user3", Item: "item1"},
			},
		},
		{
			name:      "User-specific orders",
			normal:    false,
			user:      "user1",
			itemNames: []string{"item1", "item3"},
			orders: []model.Order{
				{UserID: "user1", Item: "item1"},
				{UserID: "user2", Item: "item2"},
				{UserID: "user1", Item: "item3"},
				{UserID: "user3", Item: "item1"},
			},
			expectedOrders: []model.Order{
				{UserID: "user1", Item: "item1"},
				{UserID: "user1", Item: "item3"},
			},
		},
		{
			name:      "User-specific orders with item filter",
			normal:    false,
			user:      "user1",
			itemNames: []string{"item1"},
			orders: []model.Order{
				{UserID: "user1", Item: "item1"},
				{UserID: "user2", Item: "item2"},
				{UserID: "user1", Item: "item3"},
				{UserID: "user3", Item: "item1"},
			},
			expectedOrders: []model.Order{
				{UserID: "user1", Item: "item1"},
			},
		},
		{
			name:      "No matching orders",
			normal:    false,
			user:      "user4",
			itemNames: []string{"item1", "item2"},
			orders: []model.Order{
				{UserID: "user1", Item: "item1"},
				{UserID: "user2", Item: "item2"},
				{UserID: "user1", Item: "item3"},
				{UserID: "user3", Item: "item1"},
			},
			expectedOrders: []model.Order{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := getRelevantOrders(tc.normal, tc.user, tc.itemNames, tc.orders)
			if len(result) != len(tc.expectedOrders) {
				t.Errorf("expected %d results, got %d", len(tc.expectedOrders), len(result))
			}
			for i, expected := range tc.expectedOrders {
				foundResult := false
				var gotResult model.Order
				for _, gotResult = range result {
					if expected.Item == gotResult.Item && expected.UserID == gotResult.UserID {
						foundResult = true
						break
					}
				}
				if !foundResult || gotResult != expected {
					t.Errorf("expected result[%d] to be %+v, got %+v", i, expected, result[i])
				}
			}
		})
	}
}

func TestGetLastCandles(t *testing.T) {
	testCases := []struct {
		name            string
		candleSummary   model.CandleSummary
		itemNames       []string
		candles         []model.Candle
		expectedCandles map[string]model.Candle
	}{
		{
			name: "Matching candles",
			candleSummary: model.CandleSummary{
				Summary: map[string]model.CandleStatistic{
					"item1": {LastDate: time.Date(2023, 10, 3, 0, 0, 0, 0, time.UTC)},
					"item2": {LastDate: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC)},
				},
			},
			itemNames: []string{"item1", "item2"},
			candles: []model.Candle{
				{Item: "item1", Date: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)},
				{Item: "item2", Date: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC)},
				{Item: "item1", Date: time.Date(2023, 10, 3, 0, 0, 0, 0, time.UTC)},
				{Item: "item3", Date: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)},
			},
			expectedCandles: map[string]model.Candle{
				"item1": {Item: "item1", Date: time.Date(2023, 10, 3, 0, 0, 0, 0, time.UTC)},
				"item2": {Item: "item2", Date: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC)},
			},
		},
		{
			name: "No matching candles",
			candleSummary: model.CandleSummary{
				Summary: map[string]model.CandleStatistic{
					"item1": {LastDate: time.Date(2023, 10, 4, 0, 0, 0, 0, time.UTC)},
					"item2": {LastDate: time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC)},
				},
			},
			itemNames: []string{"item1", "item2"},
			candles: []model.Candle{
				{Item: "item1", Date: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)},
				{Item: "item2", Date: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC)},
				{Item: "item1", Date: time.Date(2023, 10, 3, 0, 0, 0, 0, time.UTC)},
				{Item: "item3", Date: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)},
			},
			expectedCandles: map[string]model.Candle{},
		},
		{
			name: "Partial matching candles",
			candleSummary: model.CandleSummary{
				Summary: map[string]model.CandleStatistic{
					"item1": {LastDate: time.Date(2023, 10, 3, 0, 0, 0, 0, time.UTC)},
					"item2": {LastDate: time.Date(2023, 10, 5, 0, 0, 0, 0, time.UTC)},
				},
			},
			itemNames: []string{"item1", "item2"},
			candles: []model.Candle{
				{Item: "item1", Date: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)},
				{Item: "item2", Date: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC)},
				{Item: "item1", Date: time.Date(2023, 10, 3, 0, 0, 0, 0, time.UTC)},
				{Item: "item3", Date: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC)},
			},
			expectedCandles: map[string]model.Candle{
				"item1": {Item: "item1", Date: time.Date(2023, 10, 3, 0, 0, 0, 0, time.UTC)},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := getlastCandles(tc.candleSummary, tc.itemNames, tc.candles)
			if len(result) != len(tc.expectedCandles) {
				t.Errorf("expected %d candles, got %d", len(tc.expectedCandles), len(result))
			}
			for item, expectedCandle := range tc.expectedCandles {
				resultCandle, exists := result[item]
				if !exists {
					t.Errorf("expected candle for item %s, but it was not found", item)
				} else if resultCandle != expectedCandle {
					t.Errorf("expected candle %+v, got %+v", expectedCandle, resultCandle)
				}
			}
		})
	}
}
