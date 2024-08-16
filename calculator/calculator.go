package calculator

import (
	"math"

	"github.com/udvarid/don-trade-golang/model"
)

func CalculateSMA(candles []model.Candle, period int) []float64 {
	sma := make([]float64, len(candles)-period+1)
	for i := range sma {
		sum := 0.0
		for j := i; j < i+period; j++ {
			sum += candles[j].Close
		}
		sma[i] = sum / float64(period)
	}
	return sma
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
	sma := CalculateSMA(candles, period)
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
