package user

import (
	"testing"
	"time"

	"github.com/udvarid/don-trade-golang/model"
)

func TestGetAssetsWithValue(t *testing.T) {
	tests := []struct {
		name           string
		assets         map[string]float64
		candleSummary  model.CandleSummary
		expectedResult []model.AssetWithValue
	}{
		{
			name: "Single asset",
			assets: map[string]float64{
				"BTC": 2.0,
			},
			candleSummary: model.CandleSummary{
				ID:   1,
				Date: time.Now(),
				Summary: map[string]model.CandleStatistic{
					"BTC": {LastPrice: 50000.0},
				},
				Persisted: []string{"BTC"},
			},
			expectedResult: []model.AssetWithValue{
				{Item: "BTC", Volume: 2.0, Price: 50000.0, Value: 100000.0},
				{Item: "Total", Value: 100000.0},
			},
		},
		{
			name: "Multiple assets",
			assets: map[string]float64{
				"BTC": 1.0,
				"ETH": 5.0,
			},
			candleSummary: model.CandleSummary{
				ID:   2,
				Date: time.Now(),
				Summary: map[string]model.CandleStatistic{
					"BTC": {LastPrice: 50000.0},
					"ETH": {LastPrice: 2000.0},
				},
				Persisted: []string{"BTC", "ETH"},
			},
			expectedResult: []model.AssetWithValue{
				{Item: "BTC", Volume: 1.0, Price: 50000.0, Value: 50000.0},
				{Item: "ETH", Volume: 5.0, Price: 2000.0, Value: 10000.0},
				{Item: "Total", Value: 60000.0},
			},
		},
		{
			name:   "No assets",
			assets: map[string]float64{},
			candleSummary: model.CandleSummary{
				ID:        3,
				Date:      time.Now(),
				Summary:   map[string]model.CandleStatistic{},
				Persisted: []string{},
			},
			expectedResult: []model.AssetWithValue{{Item: "Total", Value: 0.0}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getAssetsWithValue(tt.assets, tt.candleSummary)
			if len(result) != len(tt.expectedResult) {
				t.Errorf("expected %d results, got %d", len(tt.expectedResult), len(result))
			}
			for i, expected := range tt.expectedResult {
				foundResult := false
				var gotResult model.AssetWithValue
				for _, gotResult = range result {
					if expected.Item == gotResult.Item {
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
