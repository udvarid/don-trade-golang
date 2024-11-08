package user

import (
	"testing"
	"time"

	"github.com/udvarid/don-trade-golang/model"
)

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
