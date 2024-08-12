package main

import (
	"embed"
	"encoding/json"
	"flag"
	"sort"

	chart "github.com/udvarid/don-trade-golang/chartBuilder"
	"github.com/udvarid/don-trade-golang/collector"
	"github.com/udvarid/don-trade-golang/model"
	"github.com/udvarid/don-trade-golang/repository/candleRepository"
	"github.com/udvarid/don-trade-golang/repository/repoUtil"
)

var config = model.Configuration{}

//go:embed resources
var f embed.FS

func main() {
	configFile := flag.String("config", "conf.json", "the Json file contains the configurations")
	environment := flag.String("environment", "local", "where do we run the application, local or internet?")
	flag.Parse()

	configFileInString, _ := f.ReadFile("resources/" + *configFile)
	json.Unmarshal([]byte(configFileInString), &config)

	config.Environment = *environment
	repoUtil.Init()
	collector.CollectData(&config)

	var amzn []model.Candle
	for _, c := range candleRepository.GetAllCandles() {
		if c.Item == "AMZN" {
			amzn = append(amzn, c)
		}
	}

	sort.Slice(amzn, func(i, j int) bool {
		return amzn[i].Date.Before(amzn[j].Date)
	})

	chart.BuildSimpleCandleChart(amzn)
}
