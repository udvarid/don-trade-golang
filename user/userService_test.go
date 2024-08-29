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
func TestGetElementByDate(t *testing.T) {
	// Define test cases
	tests := []struct {
		name     string
		history  []model.HistoryElement
		date     time.Time
		expected model.HistoryElement
	}{
		{
			name: "Element found",
			history: []model.HistoryElement{
				{Date: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), Items: map[string]float64{"item1": 1.1, "item2": 2.2}},
				{Date: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC), Items: map[string]float64{"item3": 3.3, "item4": 4.4}},
			},
			date:     time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC),
			expected: model.HistoryElement{Date: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC), Items: map[string]float64{"item3": 3.3, "item4": 4.4}},
		},
		{
			name: "Element not found",
			history: []model.HistoryElement{
				{Date: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), Items: map[string]float64{"item1": 1.1, "item2": 2.2}},
				{Date: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC), Items: map[string]float64{"item3": 3.3, "item4": 4.4}},
			},
			date:     time.Date(2023, 10, 3, 0, 0, 0, 0, time.UTC),
			expected: model.HistoryElement{},
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getElementByDate(tt.history, tt.date)
			if !historyElementEqual(result, tt.expected) {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

// Helper function to compare two HistoryElement structs
func historyElementEqual(a, b model.HistoryElement) bool {
	if !a.Date.Equal(b.Date) {
		return false
	}
	if len(a.Items) != len(b.Items) {
		return false
	}
	for k, v := range a.Items {
		if bv, ok := b.Items[k]; !ok || v != bv {
			return false
		}
	}
	return true
}

func TestGetPriceHistory(t *testing.T) {
	// Define test cases
	tests := []struct {
		name      string
		candles   []model.Candle
		itemNames []string
		firstDate time.Time
		pureToday time.Time
		expected  []model.HistoryElement
	}{
		{
			name: "Basic case",
			candles: []model.Candle{
				{ID: 1, Item: "item1", Date: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), Open: 100, Close: 110, High: 115, Low: 95, Volume: 1000},
				{ID: 2, Item: "item2", Date: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), Open: 200, Close: 210, High: 215, Low: 195, Volume: 2000},
				{ID: 3, Item: "item1", Date: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC), Open: 110, Close: 120, High: 125, Low: 105, Volume: 1100},
				{ID: 4, Item: "item2", Date: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC), Open: 210, Close: 220, High: 225, Low: 205, Volume: 2100},
			},
			itemNames: []string{"item1", "item2"},
			firstDate: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			pureToday: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC),
			expected: []model.HistoryElement{
				{
					Date: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
					Items: map[string]float64{
						"item1": 110,
						"item2": 210,
					},
				},
				{
					Date: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC),
					Items: map[string]float64{
						"item1": 120,
						"item2": 220,
					},
				},
			},
		},
		{
			name: "Missing item on first date",
			candles: []model.Candle{
				{ID: 1, Item: "item1", Date: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), Open: 100, Close: 110, High: 115, Low: 95, Volume: 1000},
				{ID: 2, Item: "item2", Date: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC), Open: 200, Close: 210, High: 215, Low: 195, Volume: 2000},
			},
			itemNames: []string{"item1", "item2"},
			firstDate: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			pureToday: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC),
			expected: []model.HistoryElement{
				{
					Date: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
					Items: map[string]float64{
						"item1": 110,
					},
				},
				{
					Date: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC),
					Items: map[string]float64{
						"item1": 110,
						"item2": 210,
					},
				},
			},
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getPriceHistory(tt.candles, tt.itemNames, tt.firstDate, tt.pureToday)
			if !historyElementsEqual(result, tt.expected) {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetFirstDate(t *testing.T) {
	// Define test cases
	tests := []struct {
		name      string
		candles   []model.Candle
		itemNames []string
		pureToday time.Time
		expected  time.Time
		expectErr bool
	}{
		{
			name: "Basic case",
			candles: []model.Candle{
				{ID: 1, Item: "item1", Date: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), Open: 100, Close: 110, High: 115, Low: 95, Volume: 1000},
				{ID: 2, Item: "item2", Date: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), Open: 200, Close: 210, High: 215, Low: 195, Volume: 2000},
				{ID: 3, Item: "item1", Date: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC), Open: 110, Close: 120, High: 125, Low: 105, Volume: 1100},
				{ID: 4, Item: "item2", Date: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC), Open: 210, Close: 220, High: 225, Low: 205, Volume: 2100},
			},
			itemNames: []string{"item1", "item2"},
			pureToday: time.Date(2023, 10, 3, 0, 0, 0, 0, time.UTC),
			expected:  time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			expectErr: false,
		},
		{
			name: "Items not available on the same date",
			candles: []model.Candle{
				{ID: 1, Item: "item1", Date: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), Open: 100, Close: 110, High: 115, Low: 95, Volume: 1000},
				{ID: 2, Item: "item2", Date: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC), Open: 200, Close: 210, High: 215, Low: 195, Volume: 2000},
			},
			itemNames: []string{"item1", "item2"},
			pureToday: time.Date(2023, 10, 3, 0, 0, 0, 0, time.UTC),
			expected:  time.Time{},
			expectErr: true,
		},
		{
			name: "No matching items",
			candles: []model.Candle{
				{ID: 1, Item: "item1", Date: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC), Open: 100, Close: 110, High: 115, Low: 95, Volume: 1000},
				{ID: 2, Item: "item3", Date: time.Date(2023, 10, 2, 0, 0, 0, 0, time.UTC), Open: 200, Close: 210, High: 215, Low: 195, Volume: 2000},
			},
			itemNames: []string{"item1", "item2"},
			pureToday: time.Date(2023, 10, 3, 0, 0, 0, 0, time.UTC),
			expected:  time.Time{}, // No matching date
			expectErr: true,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := getFirstDate(tt.candles, tt.itemNames, tt.pureToday)
			if (err != nil) != tt.expectErr {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
			if (err != nil) == tt.expectErr && err == nil && !result.Equal(tt.expected) {
				t.Errorf("got %v, want %v", result, tt.expected)
			}
		})
	}
}

// Helper function to compare two slices of HistoryElement structs
func historyElementsEqual(a, b []model.HistoryElement) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !historyElementEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}
