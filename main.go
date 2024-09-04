package main

import (
	"embed"
	"encoding/json"
	"flag"

	"github.com/udvarid/don-trade-golang/authenticator"
	chart "github.com/udvarid/don-trade-golang/chartBuilder"
	"github.com/udvarid/don-trade-golang/collector"
	"github.com/udvarid/don-trade-golang/communicator"
	"github.com/udvarid/don-trade-golang/controller"
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
	remoteAddress := flag.String("remote_address", "https://don-trade-golang.fly.dev/", "remote address of the application")
	flag.Parse()

	configFileInString, _ := f.ReadFile("resources/" + *configFile)
	json.Unmarshal([]byte(configFileInString), &config)

	config.Environment = *environment
	config.RemoteAddress = *remoteAddress
	repoUtil.Init()

	forceRefresh := false
	if forceRefresh {
		cs := candleRepository.GetAllCandleSummaries()[0]
		cs.Date = cs.Date.AddDate(0, 0, -1)
		candleRepository.UpdateCandleSummary(cs)
	}

	chart.DeleteHtml()

	collector.CollectData(&config)
	communicator.Init(&config)

	authenticator.ClearOldSessions()
	/*
		userRepository.DeleteUser("udvarid@hotmail.com")

		var donat model.User
		donat.ID = "udvarid@hotmail.com"
		donat.Name = "Udvari Donát"
		donat.Config = model.UserConfig{NotifyDaily: true, NotifyAtTransaction: true}
		assets := make(map[string]float64)
		assets["USD"] = 750000
		assets["NVDA"] = 1000
		assets["AMZN"] = 1000
		donat.Assets = assets

		pureToday, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))

		trs := donat.Transactions
		var tr0 model.Transaction
		tr0.Asset = "USD"
		tr0.Date = pureToday.AddDate(0, 0, -25)
		tr0.Volume = 1000000
		var tr1 model.Transaction
		tr1.Asset = "USD"
		tr1.Date = pureToday.AddDate(0, 0, -15)
		tr1.Volume = -100000
		var tr2 model.Transaction
		tr2.Asset = "NVDA"
		tr2.Date = pureToday.AddDate(0, 0, -15)
		tr2.Volume = 1000
		var tr3 model.Transaction
		tr3.Asset = "USD"
		tr3.Date = pureToday.AddDate(0, 0, -10)
		tr3.Volume = -150000
		var tr4 model.Transaction
		tr4.Asset = "AMZN"
		tr4.Date = pureToday.AddDate(0, 0, -10)
		tr4.Volume = 1000

		trs = append(trs, tr0)
		trs = append(trs, tr1)
		trs = append(trs, tr2)
		trs = append(trs, tr3)
		trs = append(trs, tr4)
		donat.Transactions = trs
		userRepository.AddUser(donat)
	*/

	/*
		var order1 model.Order
		order1.UserID = "udvarid@hotmail.com"
		order1.Item = "NVDA"
		order1.Direction = "BUY"
		order1.Type = "MARKET"
		order1.Usd = 10000.0
		order1.AllIn = false
		order1.ValidDays = 3
		orderService.AddOrder(order1)

		var order2 model.Order
		order2.UserID = "udvarid@hotmail.com"
		order2.Item = "AMZN"
		order2.Direction = "BUY"
		order2.Type = "LIMIT"
		order2.LimitPrice = 100.0
		order2.NumberOfItems = 100.0
		order2.AllIn = false
		order2.ValidDays = 3
		orderService.AddOrder(order2)

		var order3 model.Order
		order3.UserID = "udvarid@hotmail.com"
		order3.Item = "TSLA"
		order3.Direction = "SELL"
		order3.Type = "MARKET"
		order3.AllIn = true
		order3.ValidDays = 3
		orderService.AddOrder(order3)

		var order4 model.Order
		order4.UserID = "udvarid@hotmail.com"
		order4.Item = "OTP.BD"
		order4.Direction = "SELL"
		order4.Type = "LIMIT"
		order4.LimitPrice = 50.0
		order4.AllIn = true
		order4.ValidDays = 3
		orderService.AddOrder(order4)
	*/

	controller.Init(&config)
}
