package chart

import (
	"io"
	"os"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/udvarid/don-trade-golang/calculator"
	"github.com/udvarid/don-trade-golang/model"
)

// https://github.com/go-echarts/go-echarts

type klineData struct {
	date string
	data [4]float32
}

func BuildDetailedChart(candles []model.Candle, description string) {
	page := components.NewPage()

	period := 15
	multiplier := 2.0

	bollingerBands := calculator.CalculateBollingerBands(candles, period, multiplier)

	shortPeriod := 12
	longPeriod := 26
	signalPeriod := 9

	mediumPeriod := 19

	macd := calculator.CalculateMACD(candles, shortPeriod, longPeriod, signalPeriod)
	sma := calculator.CalculateSmaLines(candles, shortPeriod, mediumPeriod, longPeriod)

	var kd []klineData
	for _, candle := range candles {
		myDate := candle.Date.Format("2006/01/02")
		klineData := klineData{date: myDate, data: [4]float32{float32(candle.Open), float32(candle.Close), float32(candle.Low), float32(candle.High)}}
		kd = append(kd, klineData)
	}

	detailedChart := klineDetailed(kd[period:], description)
	boilingerChart := boilingerLineMulti(bollingerBands)
	detailedChart.Overlap(boilingerChart)
	macdChart := macdLineMulti(macd)
	smaChart := maLineMulti(sma)

	page.AddCharts(detailedChart, smaChart, macdChart)

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

func generateLineItems(numbers []float64) []opts.LineData {
	items := make([]opts.LineData, 0)
	for _, number := range numbers {
		items = append(items, opts.LineData{Value: number})
	}
	return items
}

func boilingerLineMulti(boilingerBands []model.BollingerBand) *charts.Line {
	line := charts.NewLine()
	var upBand []float64
	var centerBand []float64
	var downBand []float64
	var date []string
	for _, bollingerBand := range boilingerBands {
		upBand = append(upBand, bollingerBand.UpperBand)
		centerBand = append(centerBand, bollingerBand.CenterBand)
		downBand = append(downBand, bollingerBand.LowerBand)
		date = append(date, bollingerBand.Date.Format("2006/01/02"))
	}

	line.SetGlobalOptions(
		charts.WithXAxisOpts(opts.XAxis{SplitNumber: 20}),
		charts.WithYAxisOpts(opts.YAxis{Scale: opts.Bool(true)}),
	)

	line.SetXAxis(date).
		AddSeries("", generateLineItems(upBand),
			charts.WithLineStyleOpts(opts.LineStyle{Color: "#D3D3D3"}),
			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true), Symbol: "none"})).
		AddSeries("", generateLineItems(centerBand),
			charts.WithLineStyleOpts(opts.LineStyle{Color: "lightblue"}),
			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true), Symbol: "none"})).
		AddSeries("", generateLineItems(downBand),
			charts.WithLineStyleOpts(opts.LineStyle{Color: "#D3D3D3"}),
			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true), Symbol: "none"}))
	return line
}

func maLineMulti(maPoints []model.Ma) *charts.Line {
	line := charts.NewLine()
	var longLine []float64
	var mediumLine []float64
	var shortLine []float64
	var date []string
	for _, ma := range maPoints {
		longLine = append(longLine, ma.MaLong)
		mediumLine = append(mediumLine, ma.MaMedium)
		shortLine = append(shortLine, ma.MaShort)
		date = append(date, ma.Date.Format("2006/01/02"))
	}

	line.SetGlobalOptions(
		charts.WithXAxisOpts(opts.XAxis{SplitNumber: 20}),
		charts.WithYAxisOpts(opts.YAxis{Scale: opts.Bool(true)}),
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

	line.SetXAxis(date).
		AddSeries("Sma long", generateLineItems(longLine),
			charts.WithLineStyleOpts(opts.LineStyle{Color: "blue"}),
			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true), Symbol: "diamond", ShowSymbol: opts.Bool(false)})).
		AddSeries("Sma medium", generateLineItems(mediumLine),
			charts.WithLineStyleOpts(opts.LineStyle{Color: "green"}),
			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true), Symbol: "diamond", ShowSymbol: opts.Bool(false)})).
		AddSeries("Sma short", generateLineItems(shortLine),
			charts.WithLineStyleOpts(opts.LineStyle{Color: "orange"}),
			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true), Symbol: "diamond", ShowSymbol: opts.Bool(false)}))
	return line
}

func macdLineMulti(macdPoints []model.Macd) *charts.Line {
	line := charts.NewLine()
	var macdLine []float64
	var signalLine []float64
	var date []string
	for _, macd := range macdPoints {
		macdLine = append(macdLine, macd.Macd)
		signalLine = append(signalLine, macd.Signal)
		date = append(date, macd.Date.Format("2006/01/02"))
	}

	line.SetGlobalOptions(
		charts.WithXAxisOpts(opts.XAxis{SplitNumber: 20}),
		charts.WithYAxisOpts(opts.YAxis{Scale: opts.Bool(true)}),
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

	line.SetXAxis(date).
		AddSeries("MACD Line", generateLineItems(macdLine),
			charts.WithLineStyleOpts(opts.LineStyle{Color: "blue"}),
			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true), Symbol: "diamond", ShowSymbol: opts.Bool(false)})).
		AddSeries("Signal Line", generateLineItems(signalLine),
			charts.WithLineStyleOpts(opts.LineStyle{Color: "green"}),
			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true), Symbol: "diamond", ShowSymbol: opts.Bool(false)}))
	return line
}
