package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"

	"github.com/udvarid/don-trade-golang/authenticator"
	chart "github.com/udvarid/don-trade-golang/chartBuilder"
	"github.com/udvarid/don-trade-golang/collector"
	"github.com/udvarid/don-trade-golang/communicator"
	"github.com/udvarid/don-trade-golang/controller"
	"github.com/udvarid/don-trade-golang/model"
	"github.com/udvarid/don-trade-golang/repository/candleRepository"
	"github.com/udvarid/don-trade-golang/repository/repoUtil"
	userService "github.com/udvarid/don-trade-golang/user"
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
		donat, _ := userRepository.FindUser("udvarid@hotmail.com")
		assets := donat.Assets
		assets["USD"] = 900000
		assets["NVDA"] = 1000

		trs := donat.Transactions
		var tr1 model.Transaction
		tr1.Asset = "USD"
		tr1.Date = trs[0].Date.AddDate(0, 0, 5)
		tr1.Volume = -100000
		var tr2 model.Transaction
		tr2.Asset = "NVDA"
		tr2.Date = trs[0].Date.AddDate(0, 0, 5)
		tr2.Volume = 1000
		trs = append(trs, tr1)
		trs = append(trs, tr2)
		donat.Transactions = trs
		userRepository.UpdateUser(donat)*/

	result := userService.GetUserHistory("udvarid@hotmail.com", 60)
	for _, r := range result {
		fmt.Println(r.Date, r.Items["USD"], r.Items["NVDA"])
	}

	controller.Init(&config)
}
