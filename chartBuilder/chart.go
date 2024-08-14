package chart

import (
	"io"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/udvarid/don-trade-golang/model"
)

// https://github.com/go-echarts/go-echarts

type klineData struct {
	date string
	data [4]float32
}

func BuildDetailedChart(candles []model.Candle, description string) {
	page := components.NewPage()

	var kd []klineData
	for _, candle := range candles {
		myDate := candle.Date.Format("2006/01/02")
		klineData := klineData{date: myDate, data: [4]float32{float32(candle.Open), float32(candle.Close), float32(candle.Low), float32(candle.High)}}
		kd = append(kd, klineData)
	}

	page.AddCharts(klineDetailed(kd, description))

	f, err := os.Create("html/kline-detailed-" + candles[0].Item + ".html")
	if err != nil {
		panic(err)

	}
	page.Render(io.MultiWriter(f))
}

func BuildSimpleCandleChart(candles []model.Candle, description string) {

	page := components.NewPage()

	var kd []klineData
	for _, candle := range candles {
		myDate := candle.Date.Format("2006/01/02")
		klineData := klineData{date: myDate, data: [4]float32{float32(candle.Open), float32(candle.Close), float32(candle.Low), float32(candle.High)}}
		kd = append(kd, klineData)
	}

	page.AddCharts(klineBase(kd, description))

	f, err := os.Create("html/kline-" + candles[0].Item + ".html")
	if err != nil {
		panic(err)

	}
	page.Render(io.MultiWriter(f))

}

func klineBase(kd []klineData, description string) *charts.Kline {
	kline := charts.NewKLine()

	x := make([]string, 0)
	y := make([]opts.KlineData, 0)
	for i := 0; i < len(kd); i++ {
		x = append(x, kd[i].date)
		y = append(y, opts.KlineData{Value: kd[i].data})
	}

	kline.SetGlobalOptions(
		charts.WithXAxisOpts(opts.XAxis{
			SplitNumber: 20,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: opts.Bool(true),
		}),
	)

	kline.SetXAxis(x).AddSeries(description, y).SetSeriesOptions(
		charts.WithItemStyleOpts(opts.ItemStyle{
			Color:        "#47b262",
			Color0:       "#eb5454",
			BorderColor:  "#47b262",
			BorderColor0: "#eb5454",
		}),
	)
	return kline
}

func klineDetailed(kd []klineData, description string) *charts.Kline {
	kline := charts.NewKLine()

	x := make([]string, 0)
	y := make([]opts.KlineData, 0)
	for i := 0; i < len(kd); i++ {
		x = append(x, kd[i].date)
		y = append(y, opts.KlineData{Value: kd[i].data})
	}

	kline.SetGlobalOptions(
		charts.WithXAxisOpts(opts.XAxis{
			SplitNumber: 20,
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: opts.Bool(true),
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "inside",
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Start:      50,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	kline.SetXAxis(x).AddSeries(description, y).SetSeriesOptions(
		charts.WithMarkPointNameTypeItemOpts(opts.MarkPointNameTypeItem{
			Name:      "highest value",
			Type:      "max",
			ValueDim:  "highest",
			ItemStyle: &opts.ItemStyle{Color: "#47b262"},
		}),
		charts.WithMarkPointNameTypeItemOpts(opts.MarkPointNameTypeItem{
			Name:      "lowest value",
			Type:      "min",
			ValueDim:  "lowest",
			ItemStyle: &opts.ItemStyle{Color: "#eb5454"},
		}),
		charts.WithMarkPointStyleOpts(opts.MarkPointStyle{
			Label: &opts.Label{
				Show: opts.Bool(true),
			},
		}),
		charts.WithItemStyleOpts(opts.ItemStyle{
			Color:        "#47b262",
			Color0:       "#eb5454",
			BorderColor:  "#47b262",
			BorderColor0: "#eb5454",
		}),
	)
	return kline
}
