package calculator

import (
	"math"

	"github.com/udvarid/don-trade-golang/model"
)

func CalculateSMA(candles []float64, period int) []float64 {
	sma := make([]float64, len(candles)-period+1)
	for i := range sma {
		sum := 0.0
		for j := i; j < i+period; j++ {
			sum += candles[j]
		}
		sma[i] = sum / float64(period)
	}
	return sma
}

func CalculateEMA(prices []float64, period int) []float64 {
	ema := make([]float64, len(prices))
	multiplier := 2.0 / (float64(period) + 1.0)

	// Calculate the first EMA using the simple moving average (SMA)
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += prices[i]
	}
	ema[period-1] = sum / float64(period)

	// Calculate the rest of the EMA values
	for i := period; i < len(prices); i++ {
		ema[i] = ((prices[i] - ema[i-1]) * multiplier) + ema[i-1]
	}

	// Return only the EMA values starting from the period-1 index
	return ema[period-1:]
}

func CalculateMACD(candles []model.Candle, shortPeriod int, longPeriod int, signalPeriod int) []model.Macd {
	closePrices := candleToFloat(candles)
	shortEMA := CalculateEMA(closePrices, shortPeriod)
	longEMA := CalculateEMA(closePrices, longPeriod)

	// Calculate MACD line (shortEMA - longEMA)
	macdLine := make([]float64, len(longEMA))
	for i := range macdLine {
		macdLine[i] = shortEMA[i+len(shortEMA)-len(longEMA)] - longEMA[i]
	}

	// Calculate Signal line (EMA of MACD line)
	signalLine := CalculateEMA(macdLine, signalPeriod)

	// Calculate Histogram (MACD line - Signal line)
	macdDtos := make([]model.Macd, len(signalLine))
	for i := range macdDtos {
		macdDtos[i] = model.Macd{
			Item:   candles[i+len(candles)-len(macdDtos)].Item,
			Date:   candles[i+len(candles)-len(macdDtos)].Date,
			Macd:   macdLine[i+len(macdLine)-len(signalLine)],
			Signal: signalLine[i],
		}
	}

	return macdDtos
}

func CalculateStandardDeviation(candles []model.Candle, sma []float64, period int) []float64 {
	stdDev := make([]float64, len(sma))
	for i := range stdDev {
		sum := 0.0
		for j := i; j < i+period; j++ {
			sum += math.Pow(candles[j].Close-sma[i], 2)
		}
		stdDev[i] = math.Sqrt(sum / float64(period))
	}
	return stdDev
}

func CalculateBollingerBands(candles []model.Candle, period int, multiplier float64) []model.BollingerBand {
	sma := CalculateSMA(candleToFloat(candles), period)
	stdDev := CalculateStandardDeviation(candles, sma, period)

	bollingerBands := make([]model.BollingerBand, len(sma))
	for i := range bollingerBands {
		bollingerBands[i] = model.BollingerBand{
			Item:       candles[i+period-1].Item,
			Date:       candles[i+period-1].Date,
			UpperBand:  sma[i] + (stdDev[i] * multiplier),
			LowerBand:  sma[i] - (stdDev[i] * multiplier),
			CenterBand: sma[i],
		}
	}
	return bollingerBands
}

func CalculateSmaLines(candles []model.Candle, shortPeriod int, mediumPeriod int, longPeriod int) []model.Ma {
	closePrices := candleToFloat(candles)
	shortMa := CalculateSMA(closePrices, shortPeriod)
	mediumMa := CalculateSMA(closePrices, mediumPeriod)
	longMa := CalculateSMA(closePrices, longPeriod)

	maDtos := make([]model.Ma, len(longMa))
	for i := range maDtos {
		maDtos[i] = model.Ma{
			Item:     candles[i+longPeriod-1].Item,
			Date:     candles[i+longPeriod-1].Date,
			MaShort:  shortMa[i+longPeriod-shortPeriod],
			MaMedium: mediumMa[i+longPeriod-mediumPeriod],
			MaLong:   longMa[i],
		}
	}

	return maDtos
}

func candleToFloat(candles []model.Candle) []float64 {
	var result []float64
	for _, candle := range candles {
		result = append(result, candle.Close)
	}
	return result
}
