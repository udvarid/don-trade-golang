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

func CalculateADX(candles []model.Candle, period int) []model.Adx {
	if len(candles) < period {
		return []model.Adx{} // Not enough data to calculate ADX
	}

	plusDM := make([]float64, len(candles))
	minusDM := make([]float64, len(candles))
	tr := make([]float64, len(candles))

	for i := 1; i < len(candles); i++ {
		highDiff := candles[i].High - candles[i-1].High
		lowDiff := candles[i-1].Low - candles[i].Low

		if highDiff > lowDiff && highDiff > 0 {
			plusDM[i] = highDiff
		} else {
			plusDM[i] = 0
		}

		if lowDiff > highDiff && lowDiff > 0 {
			minusDM[i] = lowDiff
		} else {
			minusDM[i] = 0
		}

		tr[i] = math.Max(candles[i].High-candles[i].Low, math.Max(math.Abs(candles[i].High-candles[i-1].Close), math.Abs(candles[i].Low-candles[i-1].Close)))
	}

	plusDI := CalculateEMA(plusDM, tr, period)
	minusDI := CalculateEMA(minusDM, tr, period)

	dx := make([]float64, len(plusDI)-period)
	for i := range dx {
		dx[i] = 100 * math.Abs(plusDI[i+period]-minusDI[i+period]) / (plusDI[i+period] + minusDI[i+period])
	}

	adx := CalculateEMA(dx, nil, period)

	adxDtos := make([]model.Adx, len(adx))
	for i := range adxDtos {
		adxDtos[i] = model.Adx{
			Item: candles[i+len(candles)-len(adx)].Item,
			Date: candles[i+len(candles)-len(adx)].Date,
			ADX:  adx[i],
			PDI:  plusDI[i+period],
			MDI:  minusDI[i+period],
		}
	}

	return adxDtos
}

func CalculateEMA(prices []float64, tr []float64, period int) []float64 {
	ema := make([]float64, len(prices))
	k := 2.0 / float64(period+1)

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
			if i == period-1 {
				ema[i] = prices[i] // Initial EMA if no TR is provided
			} else {
				ema[i] = (prices[i]-ema[i-1])*k + ema[i-1]
			}
		}
	}

	return ema
}

func CalculateMACD(candles []model.Candle, shortPeriod int, longPeriod int, signalPeriod int) []model.Macd {
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
			Item:   candles[i+len(candles)-len(macdDtos)].Item,
			Date:   candles[i+len(candles)-len(macdDtos)].Date,
			Macd:   macdLine[i+len(macdLine)-len(signalLine)],
			Signal: signalLine[i],
		}
	}

	return macdDtos
}

func CalculateRSI(candles []model.Candle, period int) []model.Rsi {
	if len(candles) < period {
		return []model.Rsi{} // Not enough data to calculate RSI
	}

	gains := make([]float64, len(candles))
	losses := make([]float64, len(candles))

	// Calculate initial gains and losses
	for i := 1; i < len(candles); i++ {
		change := candles[i].Close - candles[i-1].Close
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
	rsis := make([]model.Rsi, len(candles)-period)
	for i := period; i < len(candles); i++ {
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
			Item: candles[i].Item,
			Date: candles[i].Date,
			RSI:  rsi,
		}
	}

	return rsis
}

func CalculateOBV(candles []model.Candle) []model.Obv {
	if len(candles) == 0 {
		return []model.Obv{}
	}

	obvs := make([]model.Obv, len(candles))
	obvs[0] = model.Obv{
		Item: candles[0].Item,
		Date: candles[0].Date,
		Obv:  candles[0].Volume,
	}

	for i := 1; i < len(candles); i++ {
		var obv float64
		if candles[i].Close > candles[i-1].Close {
			obv = obvs[i-1].Obv + candles[i].Volume
		} else if candles[i].Close < candles[i-1].Close {
			obv = obvs[i-1].Obv - candles[i].Volume
		} else {
			obv = obvs[i-1].Obv
		}

		obvs[i] = model.Obv{
			Item: candles[i].Item,
			Date: candles[i].Date,
			Obv:  obv,
		}
	}

	return obvs
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

func CalculateTrend(data []float64) (slope, intercept float64) {
	n := float64(len(data))
	if n == 0 {
		return 0, 0 // Handle empty array
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
		return 0, 0 // Handle divide by zero case
	}
	slope = (n*sumXY - sumX*sumY) / denominator
	intercept = (sumY*sumX2 - sumX*sumXY) / denominator

	return slope, intercept
}

func candleToFloat(candles []model.Candle) []float64 {
	var result []float64
	for _, candle := range candles {
		result = append(result, candle.Close)
	}
	return result
}
