package calculator

import (
	"math"
	"testing"
	"time"

	"github.com/udvarid/don-trade-golang/model"
)

func TestCandleToFloat(t *testing.T) {
	tests := []struct {
		name     string
		candles  []model.Candle
		expected []float64
	}{
		{
			name:     "Empty slice",
			candles:  []model.Candle{},
			expected: []float64{},
		},
		{
			name: "Single candle",
			candles: []model.Candle{
				{Close: 100.0},
			},
			expected: []float64{100.0},
		},
		{
			name: "Multiple candles",
			candles: []model.Candle{
				{Close: 100.0},
				{Close: 200.0},
				{Close: 300.0},
			},
			expected: []float64{100.0, 200.0, 300.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := candleToFloat(tt.candles)
			if !sameSlice(result, tt.expected) {
				t.Errorf("candleToFloat() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCalculateSMA(t *testing.T) {
	tests := []struct {
		name     string
		prices   []float64
		period   int
		expected []float64
	}{
		{
			name:     "Empty prices",
			prices:   []float64{},
			period:   3,
			expected: []float64{},
		},
		{
			name:     "Single price",
			prices:   []float64{100.0},
			period:   1,
			expected: []float64{100.0},
		},
		{
			name:     "Period longer than prices",
			prices:   []float64{100.0, 200.0},
			period:   3,
			expected: []float64{},
		},
		{
			name:     "Normal case",
			prices:   []float64{100.0, 200.0, 300.0, 400.0, 500.0},
			period:   3,
			expected: []float64{200.0, 300.0, 400.0},
		},
		{
			name:     "Another normal case",
			prices:   []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			period:   2,
			expected: []float64{1.5, 2.5, 3.5, 4.5},
		},
		{
			name:     "Another normal case 2",
			prices:   []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			period:   4,
			expected: []float64{2.5, 3.5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateSMA(tt.prices, tt.period)
			if !sameSlice(result, tt.expected) {
				t.Errorf("CalculateSMA() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCalculateSmaLines(t *testing.T) {
	tests := []struct {
		name         string
		candles      []model.Candle
		shortPeriod  int
		mediumPeriod int
		longPeriod   int
		expected     []model.Ma
	}{
		{
			name:         "Empty candles",
			candles:      []model.Candle{},
			shortPeriod:  3,
			mediumPeriod: 5,
			longPeriod:   7,
			expected:     []model.Ma{},
		},
		{
			name: "Normal case",
			candles: []model.Candle{
				{Item: "A", Date: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC), Close: 100.0},
				{Item: "A", Date: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC), Close: 200.0},
				{Item: "A", Date: time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC), Close: 300.0},
				{Item: "A", Date: time.Date(2023, 1, 4, 0, 0, 0, 0, time.UTC), Close: 400.0},
				{Item: "A", Date: time.Date(2023, 1, 5, 0, 0, 0, 0, time.UTC), Close: 500.0},
				{Item: "A", Date: time.Date(2023, 1, 6, 0, 0, 0, 0, time.UTC), Close: 600.0},
				{Item: "A", Date: time.Date(2023, 1, 7, 0, 0, 0, 0, time.UTC), Close: 700.0},
			},
			shortPeriod:  3,
			mediumPeriod: 5,
			longPeriod:   7,
			expected: []model.Ma{
				{
					Item:     "A",
					Date:     time.Date(2023, 1, 7, 0, 0, 0, 0, time.UTC),
					MaShort:  600.0,
					MaMedium: 500.0,
					MaLong:   400.0,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateSmaLines(tt.candles, tt.shortPeriod, tt.mediumPeriod, tt.longPeriod)
			if !sameMA(result, tt.expected) {
				t.Errorf("CalculateSmaLines() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestCalculateStandardDeviation(t *testing.T) {
	tests := []struct {
		prices []float64
		ma     []float64
		period int
		want   []float64
	}{
		{
			prices: []float64{1, 2, 3, 4, 5},
			ma:     []float64{2, 3, 4},
			period: 3,
			want:   []float64{0.816496580927726, 0.816496580927726, 0.816496580927726},
		},
		{
			prices: []float64{10, 20, 30, 40, 50},
			ma:     []float64{20, 30, 40},
			period: 3,
			want:   []float64{8.16496580927726, 8.16496580927726, 8.16496580927726},
		},
		{
			prices: []float64{1, 1, 1, 1, 1},
			ma:     []float64{1, 1, 1},
			period: 3,
			want:   []float64{0, 0, 0},
		},
		{
			prices: []float64{1, 2, 3},
			ma:     []float64{2},
			period: 3,
			want:   []float64{0.816496580927726},
		},
		{
			prices: []float64{},
			ma:     []float64{},
			period: 3,
			want:   []float64{},
		},
		{
			prices: []float64{2, 3, 2, 4, 3, 4, 5},
			ma:     []float64{2.75, 3, 3.25, 4},
			period: 4,
			want:   []float64{0.82915619758885, 0.7071067811865476, 0.82915619758885, 0.7071067811865476},
		},
	}

	for _, tt := range tests {
		got := CalculateStandardDeviation(tt.prices, tt.ma, tt.period)
		if len(got) != len(tt.want) {
			t.Errorf("CalculateStandardDeviation(%v, %v, %d) = %v; want %v", tt.prices, tt.ma, tt.period, got, tt.want)
			continue
		}
		for i := range got {
			if math.Abs(got[i]-tt.want[i]) > 1e-9 {
				t.Errorf("CalculateStandardDeviation(%v, %v, %d) = %v; want %v", tt.prices, tt.ma, tt.period, got, tt.want)
				break
			}
		}
	}
}

func TestCalculateEMAWithNilTR(t *testing.T) {
	tests := []struct {
		prices   []float64
		period   int
		expected []float64
	}{
		{
			prices:   []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			period:   3,
			expected: []float64{0, 0, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			prices:   []float64{10, 20, 30, 40, 50},
			period:   2,
			expected: []float64{0, 15, 25, 35, 45},
		},
		{
			prices:   []float64{10, 20, 34, 40, 50, 65, 72},
			period:   3,
			expected: []float64{0, 0, 21.333333333333332, 30.666666666666664, 40.33333333333333, 52.666666666666664, 62.33333333333333},
		},
	}

	for _, test := range tests {
		result := CalculateEMA(test.prices, nil, test.period)
		if !sameSlice(result, test.expected) {
			t.Errorf("CalculateEMA(%v, nil, %d) = %v; want %v", test.prices, test.period, result, test.expected)
		}
	}
}

func TestCalculateOBV(t *testing.T) {
	candles := []model.Candle{
		{Item: "AAPL", Date: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), Close: 150, Volume: 1000},
		{Item: "AAPL", Date: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC), Close: 155, Volume: 1500},
		{Item: "AAPL", Date: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC), Close: 150, Volume: 1200},
		{Item: "AAPL", Date: time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC), Close: 150, Volume: 1400},
		{Item: "AAPL", Date: time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC), Close: 151, Volume: 1700},
	}

	expected := []model.Obv{
		{Item: "AAPL", Date: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC), Obv: 1000},
		{Item: "AAPL", Date: time.Date(2021, 1, 2, 0, 0, 0, 0, time.UTC), Obv: 2500},
		{Item: "AAPL", Date: time.Date(2021, 1, 3, 0, 0, 0, 0, time.UTC), Obv: 1300},
		{Item: "AAPL", Date: time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC), Obv: 1300},
		{Item: "AAPL", Date: time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC), Obv: 3000},
	}

	result := CalculateOBV(candles)

	if len(result) != len(expected) {
		t.Fatalf("expected length %d, got %d", len(expected), len(result))
	}

	for i := range result {
		if result[i] != expected[i] {
			t.Errorf("at index %d, expected %+v, got %+v", i, expected[i], result[i])
		}
	}
}

func TestCalculateRSI(t *testing.T) {

	// Test data
	candles := []model.Candle{
		createCandle("item1", "2023-01-01", 100),
		createCandle("item1", "2023-01-02", 105),
		createCandle("item1", "2023-01-03", 102),
		createCandle("item1", "2023-01-04", 108),
		createCandle("item1", "2023-01-05", 107),
		createCandle("item1", "2023-01-06", 111),
		createCandle("item1", "2023-01-07", 115),
		createCandle("item1", "2023-01-08", 113),
		createCandle("item1", "2023-01-09", 117),
		createCandle("item1", "2023-01-10", 120),
	}

	// Expected RSI values (calculated manually or using a reliable tool)
	expectedRSIs := []float64{
		78.947368,
		83.333333,
		73.732719,
		79.606440,
		83.140771,
	}

	// Calculate RSI
	period := 5
	rsis := CalculateRSI(candles, period)

	// Check the length of the result
	if len(rsis) != len(expectedRSIs) {
		t.Fatalf("expected length %d, got %d", len(expectedRSIs), len(rsis))
	}

	for i := range rsis {
		if math.Abs(rsis[i].RSI-expectedRSIs[i]) > 1e-2 {
			t.Errorf("at index %d, expected %f, got %f", i, expectedRSIs[i], rsis[i].RSI)
		}
	}

}

func TestCalculateTrend(t *testing.T) {
	tests := []struct {
		name              string
		data              []float64
		expectedSlope     float64
		expectedIntercept float64
		expectedRSquared  float64
	}{
		{
			name:              "Basic linear trend",
			data:              []float64{1, 2, 3, 4, 5},
			expectedSlope:     1,
			expectedIntercept: 1,
			expectedRSquared:  1,
		},
		{
			name:              "No trend",
			data:              []float64{1, 1, 1, 1, 1},
			expectedSlope:     0,
			expectedIntercept: 1,
			expectedRSquared:  1,
		},
		{
			name:              "Negative trend",
			data:              []float64{5, 4, 3, 2, 1},
			expectedSlope:     -1,
			expectedIntercept: 5,
			expectedRSquared:  1,
		},
		{
			name:              "Empty data",
			data:              []float64{},
			expectedSlope:     0,
			expectedIntercept: 0,
			expectedRSquared:  0,
		},
		{
			name:              "Relative strong trend",
			data:              []float64{1, 3, 3, 4, 5},
			expectedSlope:     0.9,
			expectedIntercept: 1.4,
			expectedRSquared:  0.9204545454545454,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slope, intercept, rSquared := CalculateTrend(tt.data)
			if math.Abs(slope-tt.expectedSlope) > 1e-2 ||
				math.Abs(intercept-tt.expectedIntercept) > 1e-2 ||
				math.Abs(rSquared-tt.expectedRSquared) > 1e-2 {
				t.Errorf("CalculateTrend() = (%v, %v, %v), want (%v, %v, %v)", slope, intercept, rSquared, tt.expectedSlope, tt.expectedIntercept, tt.expectedRSquared)
			}
		})
	}
}

func createCandle(item string, date string, close float64) model.Candle {
	parsedDate, _ := time.Parse("2006-01-02", date)
	return model.Candle{
		Item:  item,
		Date:  parsedDate,
		Close: close,
	}
}

func sameMA(a, b []model.Ma) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v.MaLong != b[i].MaLong || v.MaMedium != b[i].MaMedium || v.MaShort != b[i].MaShort {
			return false
		}
	}
	return true
}

func sameSlice(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
