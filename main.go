package main

import (
	"embed"
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/udvarid/don-trade-golang/authenticator"
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

	deleteHtml()

	collector.CollectData(&config)
	communicator.Init(&config)

	authenticator.ClearOldSessions()

	controller.Init(&config)
}

func deleteHtml() {
	folderPath := "./html"
	filePattern := "kline-*.html"

	fullPattern := filepath.Join(folderPath, filePattern)

	files, err := filepath.Glob(fullPattern)
	if err != nil {
		log.Fatalf("Error finding files: %v", err)
	}

	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			log.Printf("Error deleting file %s: %v", file, err)
		}
	}
}
