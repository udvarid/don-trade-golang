package chart

import (
	"io"
	"math"
	"os"
	"sort"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/udvarid/don-trade-golang/calculator"
	"github.com/udvarid/don-trade-golang/model"
)

// https://github.com/go-echarts/go-echarts

func rollOut() int {
	return 60
}

type klineData struct {
	date string
	data [4]float32
}

func BuildUserHistoryChart(history []model.HistoryElement, session string) {
	var date []string
	assets := make(map[string]bool)
	for _, element := range history {
		date = append(date, element.Date.Format("2006/01/02"))
		for item, value := range element.Items {
			if math.Abs(value) > 0.00001 {
				assets[item] = true
			}
		}
	}
	var allAssets [][]float64
	var assetNames []string
	for asset := range assets {
		assetNames = append(assetNames, asset)
		var assetValues []float64
		for _, element := range history {
			assetValues = append(assetValues, element.Items[asset]/1000)
		}
		allAssets = append(allAssets, assetValues)
	}
	bar := barChart(date, allAssets, assetNames)

	page := components.NewPage()
	page.AddCharts(bar)

	f, err := os.Create("html/kline-" + session + ".html")
	if err != nil {
		panic(err)
	}
	page.Render(io.MultiWriter(f))
	WaitUntilHtmlReady(session)
}

func BuildDetailedChart(candles []model.Candle) {
	page := components.NewPage()
	page2 := components.NewPage()

	period := 15
	multiplier := 2.0

	bollingerBands := calculator.CalculateBollingerBands(candles, period, multiplier)

	shortPeriod := 12
	longPeriod := 26
	signalPeriod := 9

	vwapShortPeriod := 10
	vwapMediumPeriod := 25
	vwapLongPeriod := 50

	macd := calculator.CalculateMACD(candles, shortPeriod, longPeriod, signalPeriod)
	sma := calculator.CalculateSmaLines(candles, vwapShortPeriod, vwapMediumPeriod, vwapLongPeriod)
	vwap := calculator.CalculateVwapLines(candles, vwapShortPeriod, vwapMediumPeriod, vwapLongPeriod)

	rsi := calculator.CalculateRSI(candles, 14)
	obv := calculator.CalculateOBV(candles)

	var kd []klineData
	for _, candle := range candles {
		myDate := candle.Date.Format("2006/01/02")
		klineData := klineData{date: myDate, data: [4]float32{float32(candle.Open), float32(candle.Close), float32(candle.Low), float32(candle.High)}}
		kd = append(kd, klineData)
	}

	detailedChart := klineDetailed(kd[period:])
	trendLineChart := trendLineOnCandle(candles[period:])
	boilingerChart := boilingerLineMulti(bollingerBands)
	detailedChart.Overlap(boilingerChart)
	detailedChart.Overlap(trendLineChart)
	macdChart := macdLineMulti(macd)
	smaChart := maLineMulti(sma)
	vwapChart := vwapLineMulti(vwap, candles[len(candles)-len(vwap):])

	rsiChart := rsiLine(rsi)
	obvChart := obvLine(obv)

	page.AddCharts(detailedChart, macdChart, rsiChart)

	f, err := os.Create("html/kline-detailed-" + candles[0].Item + ".html")
	if err != nil {
		panic(err)

	}
	page.Render(io.MultiWriter(f))

	page2.AddCharts(vwapChart, smaChart, obvChart)

	f2, err2 := os.Create("html/kline-detailed2-" + candles[0].Item + ".html")
	if err2 != nil {
		panic(err2)

	}
	page2.Render(io.MultiWriter(f2))
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

func trendLineOnCandle(candles []model.Candle) *charts.Line {
	var prices []float64
	var date []string
	for _, candle := range candles {
		prices = append(prices, candle.Close)
		date = append(date, candle.Date.Format("2006/01/02"))
	}
	trendPoints := calculator.GetTrendLine(prices, 20, 0.5)

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithXAxisOpts(opts.XAxis{SplitNumber: 20}),
		charts.WithYAxisOpts(opts.YAxis{Scale: opts.Bool(true)}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "inside",
			Start:      75,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Start:      75,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	line.SetXAxis(date[rollOut():]).
		AddSeries("", generateLineItemsForTrend(trendPoints[rollOut():]),
			charts.WithLineStyleOpts(opts.LineStyle{Color: "blue", Type: "dashed"}),
			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true), Symbol: "diamond", ShowSymbol: opts.Bool(false)}))
	return line

}

func klineDetailed(kd []klineData) *charts.Kline {
	kline := charts.NewKLine()

	x := make([]string, 0)
	y := make([]opts.KlineData, 0)
	for i := rollOut(); i < len(kd); i++ {
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
			Start:      75,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Start:      75,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	kline.SetXAxis(x).AddSeries("", y).SetSeriesOptions(
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

func generateLineItemsForTrend(numbers []model.TrendPoint) []opts.LineData {
	items := make([]opts.LineData, 0)
	for _, number := range numbers {
		if number.TrendFlag {
			items = append(items, opts.LineData{Value: number.TrendPoint})
		} else {
			items = append(items, opts.LineData{Value: nil})
		}
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

	line.SetXAxis(date[rollOut():]).
		AddSeries("", generateLineItems(upBand[rollOut():]),
			charts.WithLineStyleOpts(opts.LineStyle{Color: "#D3D3D3"}),
			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true), Symbol: "none"})).
		AddSeries("", generateLineItems(centerBand[rollOut():]),
			charts.WithLineStyleOpts(opts.LineStyle{Color: "lightblue"}),
			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true), Symbol: "none"})).
		AddSeries("", generateLineItems(downBand[rollOut():]),
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
			Start:      75,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Start:      75,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	line.ExtendYAxis(opts.YAxis{
		Type:  "value",
		Show:  opts.Bool(true),
		Scale: opts.Bool(true),
	})

	line.SetXAxis(date[rollOut():]).
		AddSeries("Sma long", generateLineItems(longLine[rollOut():]),
			charts.WithLineStyleOpts(opts.LineStyle{Color: "blue"}),
			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true), Symbol: "diamond", ShowSymbol: opts.Bool(false)})).
		AddSeries("Sma medium", generateLineItems(mediumLine[rollOut():]),
			charts.WithLineStyleOpts(opts.LineStyle{Color: "green"}),
			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true), Symbol: "diamond", ShowSymbol: opts.Bool(false)})).
		AddSeries("Sma short", generateLineItems(shortLine[rollOut():]),
			charts.WithLineStyleOpts(opts.LineStyle{Color: "orange"}),
			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true), Symbol: "diamond", ShowSymbol: opts.Bool(false)}))

	return line
}

func vwapLineMulti(maPoints []model.Ma, candles []model.Candle) *charts.Line {
	line := charts.NewLine()
	var longLine []float64
	var mediumLine []float64
	var shortLine []float64
	var date []string
	var volume []float64
	for i, ma := range maPoints {
		longLine = append(longLine, ma.MaLong)
		mediumLine = append(mediumLine, ma.MaMedium)
		shortLine = append(shortLine, ma.MaShort)
		volume = append(volume, candles[i].Volume/1000000)
		date = append(date, ma.Date.Format("2006/01/02"))
	}

	line.SetGlobalOptions(
		charts.WithXAxisOpts(opts.XAxis{SplitNumber: 20}),
		charts.WithYAxisOpts(opts.YAxis{Scale: opts.Bool(true)}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "inside",
			Start:      75,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Start:      75,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	line.ExtendYAxis(opts.YAxis{
		Type:  "value",
		Show:  opts.Bool(true),
		Scale: opts.Bool(true),
	})

	line.SetXAxis(date[rollOut():]).
		AddSeries("VWAP long", generateLineItems(longLine[rollOut():]),
			charts.WithLineStyleOpts(opts.LineStyle{Color: "blue"}),
			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true), Symbol: "diamond", ShowSymbol: opts.Bool(false)})).
		AddSeries("VWAP medium", generateLineItems(mediumLine[rollOut():]),
			charts.WithLineStyleOpts(opts.LineStyle{Color: "green"}),
			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true), Symbol: "diamond", ShowSymbol: opts.Bool(false)})).
		AddSeries("VWAP short", generateLineItems(shortLine[rollOut():]),
			charts.WithLineStyleOpts(opts.LineStyle{Color: "orange"}),
			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true), Symbol: "diamond", ShowSymbol: opts.Bool(false)})).
		AddSeries("Vol.bil.", generateLineItems(volume[rollOut():]),
			charts.WithLineStyleOpts(opts.LineStyle{Color: "red", Type: "dashed"}),
			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true), Symbol: "diamond", ShowSymbol: opts.Bool(false), YAxisIndex: 1}),
			charts.WithAreaStyleOpts(opts.AreaStyle{
				Opacity: 0.2,
			}))

	return line
}

func obvLine(obvPoints []model.Obv) *charts.Line {
	line := charts.NewLine()
	var obvLine []float64
	var date []string
	for _, obv := range obvPoints {
		obvLine = append(obvLine, obv.Obv)
		date = append(date, obv.Date.Format("2006/01/02"))
	}

	obvLinePuff := make([]float64, len(obvLine))
	copy(obvLinePuff, obvLine)
	sort.Float64Slice(obvLinePuff).Sort()
	max := math.Max(math.Abs(obvLinePuff[0]), math.Abs(obvLinePuff[len(obvLine)-1]))
	powerOfMax := math.Pow(10, math.Floor(math.Log10(max)))
	if powerOfMax == 0 {
		powerOfMax = 1
	}
	for i := range obvLine {
		obvLine[i] = obvLine[i] / powerOfMax
	}

	trendPoints := calculator.GetTrendLine(obvLine, 20, 0.5)

	line.SetGlobalOptions(
		charts.WithXAxisOpts(opts.XAxis{SplitNumber: 20}),
		charts.WithYAxisOpts(opts.YAxis{Scale: opts.Bool(true)}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "inside",
			Start:      75,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Start:      75,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	line.SetXAxis(date[rollOut():]).
		AddSeries("OBV", generateLineItems(obvLine[rollOut():]),
			charts.WithLineStyleOpts(opts.LineStyle{Color: "green"}),
			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true), Symbol: "diamond", ShowSymbol: opts.Bool(false)})).
		AddSeries("", generateLineItemsForTrend(trendPoints[rollOut():]),
			charts.WithLineStyleOpts(opts.LineStyle{Color: "blue", Type: "dashed"}),
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
			Start:      75,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Start:      75,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	line.SetXAxis(date[rollOut():]).
		AddSeries("MACD Line", generateLineItems(macdLine[rollOut():]),
			charts.WithLineStyleOpts(opts.LineStyle{Color: "blue"}),
			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true), Symbol: "diamond", ShowSymbol: opts.Bool(false)})).
		AddSeries("Signal Line", generateLineItems(signalLine[rollOut():]),
			charts.WithLineStyleOpts(opts.LineStyle{Color: "green"}),
			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true), Symbol: "diamond", ShowSymbol: opts.Bool(false)}))
	return line
}

func barChart(date []string, allAssets [][]float64, assets []string) *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "$K",
		}),
	)
	bar.SetXAxis(date)
	for i, assetArray := range allAssets {
		bar.AddSeries(assets[i], generateBarItems(assetArray)).
			SetSeriesOptions(charts.WithBarChartOpts(opts.BarChart{
				Stack: "stackA",
			}))
	}
	return bar
}

func generateBarItems(values []float64) []opts.BarData {
	items := make([]opts.BarData, 0)
	for _, value := range values {
		items = append(items, opts.BarData{Value: value})
	}
	return items
}

func rsiLine(rsiPoints []model.Rsi) *charts.Line {
	line := charts.NewLine()
	var rsiLine []float64
	var line70 []float64
	var line30 []float64
	var date []string
	for _, rsi := range rsiPoints {
		rsiLine = append(rsiLine, rsi.RSI)
		line70 = append(line70, 70)
		line30 = append(line30, 30)
		date = append(date, rsi.Date.Format("2006/01/02"))
	}

	line.SetGlobalOptions(
		charts.WithXAxisOpts(opts.XAxis{SplitNumber: 20}),
		charts.WithYAxisOpts(opts.YAxis{
			Scale: opts.Bool(true),
			Min:   "0",
			Max:   "100",
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Type:       "inside",
			Start:      75,
			End:        100,
			XAxisIndex: []int{0},
		}),
		charts.WithDataZoomOpts(opts.DataZoom{
			Start:      75,
			End:        100,
			XAxisIndex: []int{0},
		}),
	)

	line.SetXAxis(date[rollOut():]).
		AddSeries("RSI Line", generateLineItems(rsiLine[rollOut():]),
			charts.WithLineStyleOpts(opts.LineStyle{Color: "blue"}),
			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true), Symbol: "diamond", ShowSymbol: opts.Bool(false)})).
		AddSeries("RSI Line", generateLineItems(line70[rollOut():]),
			charts.WithLineStyleOpts(opts.LineStyle{Color: "red"}),
			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true), Symbol: "diamond", ShowSymbol: opts.Bool(false)})).
		AddSeries("RSI Line", generateLineItems(line30[rollOut():]),
			charts.WithLineStyleOpts(opts.LineStyle{Color: "green"}),
			charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true), Symbol: "diamond", ShowSymbol: opts.Bool(false)}))
	return line
}
