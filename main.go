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
	"github.com/udvarid/don-trade-golang/repository/sessionRepository"
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

	communicator.Init(&config)

	collector.CollectData(&config)

	authenticator.ClearOldSessions()

	cs := candleRepository.GetAllCandleSummaries()[0]
	if !cs.DailyStatusSent {
		fmt.Println("Sending daily status")
		userService.SendDailyStatus()
		cs.DailyStatusSent = true
		candleRepository.UpdateCandleSummary(cs)
	}

	activeSessions := sessionRepository.GetAllSessions()
	for _, session := range activeSessions {
		fmt.Println(session)
		if session.IsChecked {
			chart.BuildUserHistoryChart(userService.GetUserHistory(session.ID, 30), session.Session)
		}
	}

	controller.Init(&config)
}
