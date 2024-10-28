package calculator

import (
	"math"

	"github.com/udvarid/don-trade-golang/model"
)

func CalculateSMA(prices []float64, period int) []float64 {
	if len(prices) < period || len(prices) == 0 {
		return []float64{}
	}
	sma := make([]float64, len(prices)-period+1)
	for i := range sma {
		sum := 0.0
		for j := i; j < i+period; j++ {
			sum += prices[j]
		}
		sma[i] = sum / float64(period)
	}
	return sma
}

func CalculateVwap(candles *model.GroupOfCandles, period int) []float64 {
	var vwap []float64

	if len(candles.Group) < period {
		return vwap
	}

	for i := 0; i <= len(candles.Group)-period; i++ {
		var cumulativePriceVolume float64
		var cumulativeVolume float64

		for j := i; j < i+period; j++ {
			typicalPrice := (candles.Group[j].High + candles.Group[j].Low + candles.Group[j].Close) / 3.0
			priceVolume := typicalPrice * candles.Group[j].Volume
			cumulativePriceVolume += priceVolume
			cumulativeVolume += candles.Group[j].Volume
		}

		divisor := cumulativeVolume
		if divisor == 0 {
			divisor = 1
		}
		vwapValue := cumulativePriceVolume / divisor
		vwap = append(vwap, vwapValue)
	}

	return vwap
}

func CalculateEMA(prices []float64, tr []float64, period int) []float64 {
	ema := make([]float64, len(prices))
	k := 2.0 / float64(period+1)

	sum := 0.0
	for i := 0; i < period; i++ {
		sum += prices[i]
	}
	ema[period-1] = sum / float64(period)

	for i := period - 1; i < len(prices); i++ {
		if tr != nil && tr[i] != 0 {
			di := (prices[i] / tr[i]) * 100 // Calculate the DI as a percentage
			if i == period-1 {
				ema[i] = di // Initial EMA is just the DI for the first period
			} else {
				ema[i] = (di-ema[i-1])*k + ema[i-1] // EMA calculation
			}
		} else if tr != nil && tr[i] == 0 {
			ema[i] = 0 // If TR is zero, DI should be zero
		} else {
			if i >= period {
				ema[i] = (prices[i]-ema[i-1])*k + ema[i-1]
			}
		}
	}

	return ema
}

func CalculateMACD(candles *model.GroupOfCandles, shortPeriod int, longPeriod int, signalPeriod int) []model.Macd {
	closePrices := candleToFloat(candles)
	shortEMA := CalculateEMA(closePrices, nil, shortPeriod)
	longEMA := CalculateEMA(closePrices, nil, longPeriod)

	// Calculate MACD line (shortEMA - longEMA)
	macdLine := make([]float64, len(longEMA))
	for i := range macdLine {
		macdLine[i] = shortEMA[i+len(shortEMA)-len(longEMA)] - longEMA[i]
	}

	// Calculate Signal line (EMA of MACD line)
	signalLine := CalculateEMA(macdLine, nil, signalPeriod)

	// Calculate Histogram (MACD line - Signal line)
	macdDtos := make([]model.Macd, len(signalLine))
	for i := range macdDtos {
		macdDtos[i] = model.Macd{
			Item:   candles.Group[i+len(candles.Group)-len(macdDtos)].Item,
			Date:   candles.Group[i+len(candles.Group)-len(macdDtos)].Date,
			Macd:   macdLine[i+len(macdLine)-len(signalLine)],
			Signal: signalLine[i],
		}
	}

	return macdDtos
}

func CalculateRSI(candles *model.GroupOfCandles, period int) []model.Rsi {
	if len(candles.Group) < period {
		return []model.Rsi{} // Not enough data to calculate RSI
	}

	gains := make([]float64, len(candles.Group))
	losses := make([]float64, len(candles.Group))

	// Calculate initial gains and losses
	for i := 1; i < len(candles.Group); i++ {
		change := candles.Group[i].Close - candles.Group[i-1].Close
		if change > 0 {
			gains[i] = change
		} else {
			losses[i] = -change
		}
	}

	// Calculate the initial average gain and loss
	var avgGain, avgLoss float64
	for i := 1; i <= period; i++ {
		avgGain += gains[i]
		avgLoss += losses[i]
	}
	avgGain /= float64(period)
	avgLoss /= float64(period)

	// Calculate RSI for each point in the series
	rsis := make([]model.Rsi, len(candles.Group)-period)
	for i := period; i < len(candles.Group); i++ {
		if i > period {
			// Update the average gain and loss with a rolling average
			avgGain = (avgGain*(float64(period-1)) + gains[i]) / float64(period)
			avgLoss = (avgLoss*(float64(period-1)) + losses[i]) / float64(period)
		}

		var rs, rsi float64
		if avgLoss == 0 {
			rsi = 100.0 // RSI is 100 if there are no losses
		} else {
			rs = avgGain / avgLoss
			rsi = 100.0 - (100.0 / (1.0 + rs))
		}

		rsis[i-period] = model.Rsi{
			Item: candles.Group[i].Item,
			Date: candles.Group[i].Date,
			RSI:  rsi,
		}
	}

	return rsis
}

func CalculateOBV(candles *model.GroupOfCandles) []model.Obv {
	if len(candles.Group) == 0 {
		return []model.Obv{}
	}

	obvs := make([]model.Obv, len(candles.Group))
	obvs[0] = model.Obv{
		Item: candles.Group[0].Item,
		Date: candles.Group[0].Date,
		Obv:  candles.Group[0].Volume,
	}

	for i := 1; i < len(candles.Group); i++ {
		var obv float64
		if candles.Group[i].Close > candles.Group[i-1].Close {
			obv = obvs[i-1].Obv + candles.Group[i].Volume
		} else if candles.Group[i].Close < candles.Group[i-1].Close {
			obv = obvs[i-1].Obv - candles.Group[i].Volume
		} else {
			obv = obvs[i-1].Obv
		}

		obvs[i] = model.Obv{
			Item: candles.Group[i].Item,
			Date: candles.Group[i].Date,
			Obv:  obv,
		}
	}

	return obvs
}

func CalculateStandardDeviation(prices []float64, ma []float64, period int) []float64 {
	if len(prices) < period || len(prices) == 0 {
		return []float64{}
	}
	stdDev := make([]float64, len(ma))
	for i := range stdDev {
		sum := 0.0
		for j := i; j < i+period; j++ {
			sum += math.Pow(prices[j]-ma[i], 2)
		}
		stdDev[i] = math.Sqrt(sum / float64(period))
	}
	return stdDev
}

func CalculateBollingerBands(candles *model.GroupOfCandles, period int, multiplier float64) []model.BollingerBand {
	sma := CalculateSMA(candleToFloat(candles), period)
	stdDev := CalculateStandardDeviation(candleToFloat(candles), sma, period)

	bollingerBands := make([]model.BollingerBand, len(sma))
	for i := range bollingerBands {
		bollingerBands[i] = model.BollingerBand{
			Item:       candles.Group[i+period-1].Item,
			Date:       candles.Group[i+period-1].Date,
			UpperBand:  sma[i] + (stdDev[i] * multiplier),
			LowerBand:  sma[i] - (stdDev[i] * multiplier),
			CenterBand: sma[i],
		}
	}
	return bollingerBands
}

func CalculateSmaLines(candles *model.GroupOfCandles, shortPeriod int, mediumPeriod int, longPeriod int) []model.Ma {
	closePrices := candleToFloat(candles)
	shortMa := CalculateSMA(closePrices, shortPeriod)
	mediumMa := CalculateSMA(closePrices, mediumPeriod)
	longMa := CalculateSMA(closePrices, longPeriod)

	maDtos := make([]model.Ma, len(longMa))
	for i := range maDtos {
		maDtos[i] = model.Ma{
			Item:     candles.Group[i+longPeriod-1].Item,
			Date:     candles.Group[i+longPeriod-1].Date,
			MaShort:  shortMa[i+longPeriod-shortPeriod],
			MaMedium: mediumMa[i+longPeriod-mediumPeriod],
			MaLong:   longMa[i],
		}
	}

	return maDtos
}

func CalculateVwapLines(candles *model.GroupOfCandles, shortPeriod int, mediumPeriod int, longPeriod int) []model.Ma {
	shortVwap := CalculateVwap(candles, shortPeriod)
	mediumVwap := CalculateVwap(candles, mediumPeriod)
	longVwap := CalculateVwap(candles, longPeriod)

	vwapDtos := make([]model.Ma, len(longVwap))
	for i := range vwapDtos {
		vwapDtos[i] = model.Ma{
			Item:     candles.Group[i+longPeriod-1].Item,
			Date:     candles.Group[i+longPeriod-1].Date,
			MaShort:  shortVwap[i+longPeriod-shortPeriod],
			MaMedium: mediumVwap[i+longPeriod-mediumPeriod],
			MaLong:   longVwap[i],
		}
	}

	return vwapDtos
}

func CalculateTrend(data []float64) (slope, intercept, rSquared float64) {
	n := float64(len(data))
	if n == 0 {
		return 0, 0, 0 // Handle empty array
	}

	var sumX, sumY, sumXY, sumX2 float64

	for i := 0; i < int(n); i++ {
		x := float64(i)
		y := data[i]

		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}

	// Calculate the slope (m) and intercept (b)
	denominator := n*sumX2 - sumX*sumX
	if denominator == 0 {
		return 0, 0, 0 // Handle divide by zero case
	}
	slope = (n*sumXY - sumX*sumY) / denominator
	intercept = (sumY*sumX2 - sumX*sumXY) / denominator

	// Calculate R²
	var ssTotal, ssResidual float64
	meanY := sumY / n

	for i := 0; i < int(n); i++ {
		x := float64(i)
		y := data[i]
		yPredicted := slope*x + intercept

		ssTotal += (y - meanY) * (y - meanY)
		ssResidual += (y - yPredicted) * (y - yPredicted)
	}

	if ssTotal == 0 {
		return slope, intercept, 1 // Perfect fit case
	}

	rSquared = 1 - (ssResidual / ssTotal)

	return slope, intercept, rSquared
}

func GetTrendLine(dataLine []float64, period int, strenght float64) (trendLine []model.TrendPoint) {
	trendPoints := make([]model.TrendPoint, len(dataLine))
	for i := range trendPoints {
		trendPoints[i].TrendFlag = false
	}

	if len(dataLine) < period {
		return
	}

	for i := 0; i <= len(dataLine)-period; i++ {
		slope, intercept, rSquared := CalculateTrend(dataLine[i : i+period])
		if rSquared < strenght {
			continue
		}
		for j := i; j < i+period; j++ {
			trendPoints[j].TrendFlag = true
			trendPoints[j].TrendPoint = slope*float64(j-i) + intercept
		}
		i += period
	}

	return trendPoints
}

func candleToFloat(candles *model.GroupOfCandles) []float64 {
	var result []float64
	for _, candle := range candles.Group {
		result = append(result, candle.Close)
	}
	return result
}
