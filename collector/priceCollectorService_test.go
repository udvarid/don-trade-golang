package collector

import (
	"fmt"
	"slices"
	"testing"
	"time"

	"github.com/udvarid/don-trade-golang/model"
)

/*
shouldDelete
*/

func getSummaries() []model.CandleSummary {
	var summaries []model.CandleSummary
	var candleSummary model.CandleSummary
	summary := make(map[string]model.CandleStatistic)
	var cs1 model.CandleStatistic
	var cs2 model.CandleStatistic
	cs1.Number = 10
	cs2.Number = 500
	summary["test1"] = cs1
	summary["test2"] = cs2
	candleSummary.Summary = summary
	summaries = append(summaries, candleSummary)
	return summaries
}

func TestGetDayParameter(t *testing.T) {
	summaries := getSummaries()
	var tests = []struct {
		s    []model.CandleSummary
		b    string
		want int
	}{
		{summaries, "test1", 700},
		{summaries, "test2", 15},
		{summaries, "test3", 700},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%s,%d", tt.b, tt.want)
		t.Run(testname, func(t *testing.T) {
			ans := getDayParameter(tt.s, tt.b, 15)
			if ans != tt.want {
				t.Errorf("got %d, want %d", ans, tt.want)
			}
		})
	}
}

func getCandles() []model.Candle {
	var candles []model.Candle
	var candleOne model.Candle
	candleOne.Item = "TEST1"
	candleOne.Date, _ = time.Parse("2006-01-02", "2024-01-20")
	candles = append(candles, candleOne)
	var candleTwo model.Candle
	candleTwo.Item = "TEST2"
	candleTwo.Date, _ = time.Parse("2006-01-02", "2024-01-21")
	candles = append(candles, candleTwo)
	var candleThree model.Candle
	candleThree.Item = "TEST2"
	candleThree.Date, _ = time.Parse("2006-01-02", "2024-01-22")
	candles = append(candles, candleThree)
	return candles
}

func getCandle(item string, date string) model.Candle {
	var candle model.Candle
	candle.Item = item
	candle.Date, _ = time.Parse("2006-01-02", date)
	return candle
}

func TestIsCandleNew(t *testing.T) {
	candles := getCandles()
	var tests = []struct {
		s    []model.Candle
		b    model.Candle
		want bool
	}{
		{candles, getCandle("TEST1", "2024-01-20"), false},
		{candles, getCandle("TEST1", "2024-01-21"), true},
		{candles, getCandle("TEST2", "2024-01-21"), false},
		{candles, getCandle("TEST2", "2024-01-20"), true},
		{candles, getCandle("TEST2", "2024-01-22"), false},
		{candles, getCandle("TEST2", "2024-01-23"), true},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%v,%t", tt.b, tt.want)
		t.Run(testname, func(t *testing.T) {
			ans := isCandleNew(&tt.s, &tt.b)
			if ans != tt.want {
				t.Errorf("got %t, want %t", ans, tt.want)
			}
		})
	}
}

func getItemMap() map[string][]model.Item {
	itemMap := make(map[string][]model.Item)

	var items []model.Item
	var item1 model.Item
	item1.Name = "Test1"
	items = append(items, item1)
	var item2 model.Item
	item2.Name = "Test2"
	items = append(items, item2)
	itemMap["TEST1"] = items

	var items2 []model.Item
	var item3 model.Item
	item3.Name = "Test3"
	items2 = append(items2, item3)
	var item4 model.Item
	item4.Name = "Test2"
	items2 = append(items2, item4)

	itemMap["TEST2"] = items2
	return itemMap
}

func TestGetItemsFromItemMap(t *testing.T) {
	answer := GetItemsFromItemMap(getItemMap())
	var keys []string
	for itemName := range answer {
		keys = append(keys, itemName)
	}
	if len(answer) != 3 || !slices.Contains(keys, "Test1") || !slices.Contains(keys, "Test2") || !slices.Contains(keys, "Test3") {
		t.Errorf("Not proper answer for getItemsFromItemMap")
	}
}

func TestShouldBeDeleted(t *testing.T) {
	itemNamesWithItem := GetItemsFromItemMap(getItemMap())
	var itemNames []string
	for itemName := range itemNamesWithItem {
		itemNames = append(itemNames, itemName)
	}
	myDate, _ := time.Parse("2006-01-02", "2023-01-01")

	var tests = []struct {
		c    model.Candle
		in   []string
		d    time.Time
		want bool
	}{
		{getCandle("Test1", "2024-02-20"), itemNames, myDate, false},
		{getCandle("Test2", "2023-02-20"), itemNames, myDate, false},
		{getCandle("Test4", "2024-02-20"), itemNames, myDate, true},
		{getCandle("Test1", "2022-02-20"), itemNames, myDate, true},
	}
	for _, tt := range tests {
		testname := fmt.Sprint(tt.c.Item)
		t.Run(testname, func(t *testing.T) {
			ans := shouldBeDeleted(&tt.c, tt.in, tt.d)
			if ans != tt.want {
				t.Errorf("got %t, want %t", ans, tt.want)
			}
		})
	}

}
